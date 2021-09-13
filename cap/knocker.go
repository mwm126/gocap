// Go Implementation of the Cloaked Access Protocol (CAP) client
// CAP was developed by Aeolus Technologies, Inc.
// (C)opyright 2013, Aeolus Technologies, Inc.  All rights reserved.

package cap

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"github.com/beevik/ntp"
	"log"
	"net"
	"time"
)

type OneTimePassword [16]byte
type SHADigest [32]byte

// Knocker send port knock UDP packet
type Knocker interface {
	Knock(username, password string, network net.IP)
}

// PortKnocker for actual Knocker implementation
type PortKnocker struct {
	yubikey Yubikey
	entropy [32]byte
}

func NewPortKnocker(yk Yubikey, ent [32]byte) PortKnocker {
	return PortKnocker{yk, ent}
}

func (sk *PortKnocker) Knock(uname, pword string, addr net.IP) {
	log.Println("Sending CAP packet...")
	time.Sleep(1 * time.Second)
	response, err := ntp.Query("pool.ntp.org")
	if err != nil {
		log.Printf("Some NTP error %v", err)
		return
	}
	timestamp := int32(time.Now().Add(response.ClockOffset).Unix())

	packet, err := sk.makePacket(uname, pword, timestamp)
	if err != nil {
		log.Printf("Could not make CAP packet: %v", err)
		return
	}

	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		log.Printf("Some error %v", err)
		return
	}

	buf := bytes.Buffer{}
	binary.Write(&buf, binary.LittleEndian, timestamp)
	timestampBytes := buf.Bytes()

	conn.Write(packet)
	conn.Write(timestampBytes)
	conn.Close()
}

func (sk *PortKnocker) makePacket(uname, pword string, timestamp int32) ([]byte, error) {
	OTP := getOTP(sk.yubikey, sk.entropy[:])

	var initVec [16]byte
	digest := makeSHADigest(sk.entropy[:], OTP[:])
	copy(initVec[:], digest[:16])
	challenge, response, err := getChallengeResponse(sk.yubikey, OTP, sk.entropy[:])
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
	copy(auth[12:], net.ParseIP("74.109.234.77").To4())
	copy(ssh[12:], net.ParseIP("74.109.234.77").To4())
	copy(server[12:], net.ParseIP("104.154.139.11").To4())

	var buf bytes.Buffer
	buf.Write(auth[:])
	buf.Write(ssh[:])
	buf.Write(server[:])
	buf.Write(OTP[:])
	buf.Write(user[:])
	buf.Write(sk.entropy[:])
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
	iv := []byte(initVec[:])
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainBlock)

	var trimmedCiphertext [160]byte
	copy(trimmedCiphertext[:], ciphertext[16:])

	serial, err := sk.yubikey.findSerial()
	if err != nil {
		log.Println("Unable to get yubikey serial number: ", err)
		return nil, err
	}

	buf = bytes.Buffer{}
	binary.Write(&buf, binary.LittleEndian, timestamp)
	binary.Write(&buf, binary.LittleEndian, serial)
	binary.Write(&buf, binary.LittleEndian, initVec)
	binary.Write(&buf, binary.LittleEndian, challenge)
	binary.Write(&buf, binary.LittleEndian, trimmedCiphertext)

	header, _ := hex.DecodeString(
		"823220d0df9234263797c5d0c5fee27ab087f86e76f82efe0bb386cc65ae879f",
	)
	macBlock := buf.Bytes()
	footer, _ := hex.DecodeString(
		"50266198ce6bae2069546cbcae0f80ba847598f674f5d7343f90e6c6e56dfa8a",
	)
	digest = makeSHADigest(header, macBlock, footer)

	buf = bytes.Buffer{}
	binary.Write(&buf, binary.LittleEndian, macBlock)
	binary.Write(&buf, binary.LittleEndian, digest)
	return buf.Bytes(), nil
}

func getOTP(yk Yubikey, entropy []byte) OneTimePassword {
	// Get Yubico OTP response from client.connection.yubikey
	digest := makeSHADigest([]byte("yubicoChal"), entropy)
	var yubicoChal [16]byte
	copy(yubicoChal[:], digest[:16])
	time.Sleep(1 * time.Second)
	return yk.challengeResponse(yubicoChal)
}

func getChallengeResponse(
	yk Yubikey,
	OTP OneTimePassword,
	entropy []byte,
) (SHADigest, [16]byte, error) {
	// Build challenge using entropy and OTP so it is unique
	challenge := makeSHADigest([]byte("SHA1-HMACChallenge"), OTP[:], entropy)
	// Get HMAC-SHA1 response from client.connection.yubikey
	time.Sleep(1 * time.Second)
	response, err := yk.challengeResponseHMAC(challenge)
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
