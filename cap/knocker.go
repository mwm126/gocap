// Go Implementation of the Cloaked Access Protocol (CAP) client
// CAP was developed by Aeolus Technologies, Inc.
// (C)opyright 2013, Aeolus Technologies, Inc.  All rights reserved.

package cap

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/beevik/ntp"
)

type OneTimePassword [16]byte
type SHADigest [32]byte

// Knocker send port knock UDP packet
type Knocker struct {
	Callbacks        []func(bool)
	Entropy          [32]byte
	Yubikey          Yubikey
	delay            uint
	yubikeyAvailable bool
}

func NewKnocker(yk Yubikey, delay uint) *Knocker {
	var entropy [32]byte
	_, err := rand.Read(entropy[:])
	if err != nil {
		log.Fatal("Unable to get entropy to send CAP packet")
	}

	return &Knocker{make([]func(bool), 0), entropy, yk, delay, false}
}

func (sk *Knocker) YubikeyAvailable() bool {
	return sk.yubikeyAvailable
}

func (sk *Knocker) AddCallback(cb func(bool)) {
	sk.Callbacks = append(sk.Callbacks, cb)
}

func (sk *Knocker) StartMonitor() {
	if sk.delay == 0 {
		return
	}
	go func() {
		for {
			for i := uint(0); i < sk.delay; i++ {
				time.Sleep(time.Second)
				serial, err := sk.Yubikey.FindSerial()
				sk.yubikeyAvailable = err == nil && serial > 0
				for _, cb := range sk.Callbacks {
					cb(sk.yubikeyAvailable)
				}
				if sk.yubikeyAvailable {
					break
				}
			}
		}
	}()
}

func (sk Knocker) Knock(uname string, ext_addr, server_addr net.IP, port uint) error {
	log.Println("Sending CAP packet...")
	time.Sleep(1 * time.Second)
	timestamp, err := getNtpTime()
	if err != nil {
		log.Printf("Unable to get NTP time:  %v", err)
		log.Printf("Warning: going to use local time, without checking for NTP offset")
		timestamp = (int32)(time.Now().Unix())
	}

	auth_addr := ext_addr
	ssh_addr := ext_addr
	packet, err := sk.makePacket(uname, timestamp, auth_addr, ssh_addr, server_addr)
	if err != nil {
		log.Printf("Could not make CAP packet:  %v", err)
		return err
	}

	addrPort := fmt.Sprintf("%s:%d", server_addr, port)
	conn, err := net.Dial("udp", addrPort)
	if err != nil {
		log.Printf("Unable to connect to CAP server:  %v", err)
		return err
	}

	_, err = conn.Write(packet)
	if err != nil {
		return err
	}
	conn.Close()
	time.Sleep(1 * time.Second)
	return nil
}

func getNtpTime() (int32, error) {
	timestamp, err := ntp.Time("0.pool.ntp.org")
	return int32(timestamp.Unix()), err
}

func (sk Knocker) makePacket(
	uname string,
	timestamp int32,
	auth_addr, ssh_addr, server_addr net.IP,
) ([]byte, error) {
	OTP, err := getOTP(sk.Yubikey, sk.Entropy[:])
	if err != nil {
		log.Println("could not get OTP", err)
		return nil, err
	}

	var initVec [16]byte
	digest := makeSHADigest(sk.Entropy[:], OTP[:])
	copy(initVec[:], digest[:16])
	challenge, response, err := getChallengeResponse(sk.Yubikey, OTP, sk.Entropy[:])
	if err != nil {
		log.Println("could not get challenge response", err)
		return nil, err
	}
	aeskey := makeSHADigest(response[:], challenge[:])

	var user [32]byte
	var auth [16]byte
	var ssh [16]byte
	var server [16]byte
	copy(user[:], []byte(uname))
	copy(auth[12:], auth_addr.To4())
	copy(ssh[12:], ssh_addr.To4())
	copy(server[12:], server_addr.To4())

	buf := new(bytes.Buffer)
	buf.Write(auth[:])
	buf.Write(ssh[:])
	buf.Write(server[:])
	buf.Write(OTP[:])
	buf.Write(user[:])
	buf.Write(sk.Entropy[:])
	plainBlock := plainBlockWithChecksum(buf.Bytes())

	key := []byte(aeskey[:])
	if len(plainBlock)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plainBlock))
	mode := cipher.NewCBCEncrypter(block, initVec[:])
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainBlock)

	var trimmedCiphertext [160]byte
	copy(trimmedCiphertext[:], ciphertext[16:])

	serial, err := sk.Yubikey.FindSerial()
	if err != nil {
		log.Println("Unable to get yubikey serial number: ", err)
		return nil, err
	}

	buf = new(bytes.Buffer)
	for _, field := range []interface{}{timestamp, serial, initVec, challenge, trimmedCiphertext} {
		err := binary.Write(buf, binary.LittleEndian, field)
		if err != nil {
			return nil, err
		}
	}

	header, _ := hex.DecodeString(
		"823220d0df9234263797c5d0c5fee27ab087f86e76f82efe0bb386cc65ae879f",
	)
	macBlock := buf.Bytes()
	footer, _ := hex.DecodeString(
		"50266198ce6bae2069546cbcae0f80ba847598f674f5d7343f90e6c6e56dfa8a",
	)
	digest = makeSHADigest(header, macBlock, footer)

	buf = new(bytes.Buffer)
	for _, field := range []interface{}{macBlock, digest} {
		err := binary.Write(buf, binary.LittleEndian, field)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func getOTP(yk Yubikey, entropy []byte) (OneTimePassword, error) {
	// Get Yubico OTP response from client.connection.yubikey
	digest := makeSHADigest([]byte("yubicoChal"), entropy)
	var yubicoChal [6]byte
	copy(yubicoChal[:], digest[:16])
	time.Sleep(1 * time.Second)
	response, err := yk.ChallengeResponse(yubicoChal)
	if err != nil {
		log.Println("Unable to get OTP", err)
		return OneTimePassword{}, err
	}
	return response, nil
}

func getChallengeResponse(
	yk Yubikey,
	OTP OneTimePassword,
	entropy []byte,
) (SHADigest, [20]byte, error) {
	// Build challenge using entropy and OTP so it is unique
	challenge := makeSHADigest([]byte("SHA1-HMACChallenge"), OTP[:], entropy)
	// Get HMAC-SHA1 response from client.connection.yubikey
	time.Sleep(1 * time.Second)
	response, err := yk.ChallengeResponseHMAC(challenge)
	if err != nil {
		log.Printf("Error getting challenge response %v", err)
		return challenge, response, err
	}
	return challenge, response, nil
}

func plainBlockWithChecksum(plainBlock []byte) []byte {
	var buf bytes.Buffer
	chksum := makeSHADigest(plainBlock)
	buf.Write(plainBlock)
	buf.Write(chksum[:])

	return buf.Bytes()

}

func makeSHADigest(args ...[]byte) SHADigest {
	var buf bytes.Buffer

	for _, arg := range args {
		buf.Write(arg)
	}
	sum := sha256.Sum256(buf.Bytes())
	return sum
}
