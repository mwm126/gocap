// Go Implementation of the Cloaked Access Protocol (CAP) client
// CAP was developed by Aeolus Technologies, Inc.
// (C)opyright 2013, Aeolus Technologies, Inc.  All rights reserved.

package cap

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math/rand"
	"time"
)

// Knocker send port knock UDP packet
type Knocker interface {
	Knock(username, password string)
}

// PortKnocker for actual Knocker implementation
type PortKnocker struct {
	username string
	password string
	yubikey  Yubikey
}

// Shared secret seed values for initial MAC test
func getSeedHeaderList() [32]byte {
	return [32]byte{
		0xEF,
		0xA8,
		0xFE,
		0x36,
		0xEB,
		0x80,
		0x02,
		0x5C,
		0x0F,
		0xFD,
		0x09,
		0x1A,
		0xA9,
		0x1C,
		0x50,
		0xF8,
		0x3E,
		0xEB,
		0x52,
		0x74,
		0x9C,
		0x56,
		0xA4,
		0x44,
		0x7B,
		0x31,
		0x6C,
		0x1A,
		0xE5,
		0xBC,
		0xF7,
		0x5D,
	}
}

func getSeedFooterList() [32]byte {
	return [32]byte{
		0x42,
		0x3C,
		0x05,
		0xF2,
		0xC4,
		0x9B,
		0x8C,
		0x3E,
		0x79,
		0x16,
		0xBA,
		0xD2,
		0x54,
		0xD7,
		0x92,
		0x48,
		0xC2,
		0x55,
		0xBA,
		0x8C,
		0xE2,
		0xE5,
		0xE5,
		0xD2,
		0x1B,
		0x1A,
		0x1C,
		0xBC,
		0x49,
		0xAB,
		0x28,
		0x18,
	}
}

// class PacketSender:
//     def __init__(
//         self, configuration: Configuration, login_info: LoginInfo, yk: YubikeyInterface
//     ) -> None:
//         self._conf = configuration
//         self._login_info = login_info
//         self._yk = yk

func (sk *PortKnocker) makePacket(uname, pword string) []byte {
	//         Prepare data for CAP packet

	time.Sleep(1 * time.Second)
	//         serial = self._yk.find_serial()
	//         serverAddressTxt = self._conf.serverAddressTxt
	//         authAddressTxt = self._conf.externalAddress

	timeOffset := 1234
	log.Println("timeOffset: ", timeOffset)
	timestamp := int(time.Now().Unix()) + timeOffset
	timestamp = 1627324072
	log.Println("timestamp: ", timestamp)

	plainBlock := make([]byte, 240)
	// buf := new(bytes.Buffer)

	entropy := make([]byte, 32)
	rand.Read(entropy)
	OTP := getOTP(sk.yubikey, entropy)

	initVec := makeSHADigest(entropy, OTP[:])
	challenge, response := getChallengeResponse(sk.yubikey, OTP, entropy)
	serial := 1234

	// createPacket(serial int, initVec, chal, resp [32]byte, plainBlock []byte, timestamp int) []byte {
	var chal [32]byte
	chal = challenge
	var resp [32]byte
	resp = response

	pack := sk.createPacket(serial, initVec, chal, resp, plainBlock, timestamp)

	// sk.createPacket(serial,
	// 	make([]byte, 32), make([]byte, 32), make([]byte, 32),
	// 	plainBlock, timestamp)
	// pack := sk.createPacket(sk.yubikey.findSerial(), initVec, challenge, response, plainBlock, timestamp)

	//         addresses = Addresses(serverAddressTxt, authAddressTxt, self._conf.externalAddress)
	//         plainBlock = ChecksumBlock(addresses, OTP, self._login_info.username, entropy).plainBlock()

	log.Println(
		"Sending CAP packet to ",
		//             self._conf.serverAddressTxt,
		":",
		//             self._conf.serverPort,
		" with ",
	//             self._login_info.username,
	)

	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, pack)
	if err != nil {
		log.Fatal(err)
	}
	timestamp = 1627324072
	err = binary.Write(buf, binary.LittleEndian, int64(timestamp))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("HERE IS THE BUFFFFFF", buf)
	return buf.Bytes()
}

func (sk *PortKnocker) Knock(uname, pword string) {
	//         Send off CAP UDP packet
	//         sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
	// binary.Write(buf, binary.LittleEndian, plainBlock)
	// binary.Write(buf, binary.LittleEndian, timestamp)
	//         sock.sendto(
	//             pack.packet(plainBlock, timestamp), (self._conf.serverAddressTxt, self._conf.serverPort)
	//         )
}

func getOTP(yk Yubikey, entropy []byte) [32]byte {
	// Get Yubico OTP response from client.connection.yubikey
	yubicoChal := makeSHADigest([]byte("yubicoChal"), entropy)
	time.Sleep(1 * time.Second)
	log.Println("lenggggggggggggggth  ", len(yubicoChal))
	log.Println("LENGGGGGGGGGGGGGGTH  ", yubicoChal)
	var qwer [32]byte
	return qwer
	asdf := yk.challengeResponses(yubicoChal)
	log.Println("asdffffffffffff  ", asdf)
	return asdf
}

func getChallengeResponse(yk Yubikey, OTP [32]byte, entropy []byte) ([32]byte, [32]byte) {
	// Build challenge using entropy and OTP so it is unique
	challenge := makeSHADigest([]byte("SHA1-HMACChallenge"), OTP[:], entropy)
	// Get HMAC-SHA1 response from client.connection.yubikey
	time.Sleep(1 * time.Second)
	var qwer [32]byte
	return challenge, qwer
	response := yk.challengeResponseHMAC(challenge)
	//         return challenge, response
	return challenge, response
}

type PacketFactory struct {
	serial    []byte
	initVec   []byte
	challenge []byte
	response  []byte
	//     def __init__(self, serial: int, init_vec: bytes, chal: bytes, resp: bytes) -> None:
	//         self.serial = serial
	//         self.init_vec = init_vec
	//         self.challenge = chal
	//         self.response = resp
}

func (sk *PortKnocker) createPacket(serial int, initVec, chal, resp [32]byte, plainBlock []byte, timestamp int) []byte {
	return make([]byte, 234)
}

//     def packet(self, plainBlock: bytes, timestamp: int) -> bytes:
//         """Assemble packet"""
//         MACBlock = struct.pack(
//             "<ll16s32s160s",
//             timestamp,
//             self.serial,
//             self.init_vec,
//             self.challenge,
//             self.encryptedBlock(plainBlock),
//         )
//         sha = SHA256.new()
//         sha.update(getHeaderSecret())
//         sha.update(MACBlock)
//         sha.update(getFooterSecret())
//         MAC = sha.digest()
//         return MACBlock + MAC

func encryptedBlock(plainBlock []byte) []byte {
	//         """Encrypt payload with aes256key"""
	//         encBlock = b""
	//         cipher = AES.new(self._aes256key(), AES.MODE_CBC, self.init_vec)
	//         for i in range(0, 10):
	//             chunk = plainBlock[16 * i : 16 * (i + 1)]
	//             encBlock += cipher.encrypt(chunk)
	return make([]byte, 32)
}

//     def _aes256key(self) -> bytes:
//         """Build AES256 encryption key using response+challenge"""
//         return makeSHADigest(self.response, self.challenge)

// def getHeaderSecret() -> bytes:
//     """Derive shared secret from seed values and generator"""
//     headerSecret = bytearray()
//     generator = 0x55
//     for byte in SEED_HEADER_LIST:
//         generator = (generator * byte + 0x12) % 256
//         headerSecret.append(generator ^ byte)
//     return headerSecret

// def getFooterSecret() -> bytes:
//     """Derive shared secret from seed values and generator"""
//     footerSecret = bytearray()
//     generator = 0x18
//     for byte in SEED_FOOTER_LIST:
//         generator = (generator * byte + 0xE2) % 256
//         footerSecret.append(generator ^ byte)
//     return footerSecret

// class ChecksumBlock:
//     def __init__(self, addr: "Addresses", OTP: bytes, user: str, ent: bytes) -> None:
//         self.OTP = OTP
//         self.user = user
//         self.entropy = ent
//         self.addresses = addr

func getChecksumBlock() []byte {
	var buf bytes.Buffer

	//             self.addresses.auth(),
	//             self.addresses.ssh(),
	//             self.addresses.server(),
	//             self.OTP,
	//             self.user.encode("utf-8"),
	//             self.entropy,

	// buf.Write()
	return buf.Bytes()
}

func plainBlock() []byte {
	var buf bytes.Buffer
	checksumBlock := getChecksumBlock()
	chksum := makeSHADigest(checksumBlock)
	buf.Write(checksumBlock)
	buf.Write(chksum[:])

	return buf.Bytes()

}

func makeSHADigest(args ...[]byte) [32]byte {
	var buf bytes.Buffer

	for _, arg := range args {
		buf.Write(arg)
	}
	sum := sha256.Sum256(buf.Bytes())
	return sum
}
