package cap

import (
	"encoding/hex"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeYubikey struct{}

func (yk *FakeYubikey) FindSerial() (int32, error) {
	return 5417533, nil
}

func (yk *FakeYubikey) challengeResponse(chal [6]byte) ([16]byte, error) {
	if hex.EncodeToString(chal[:]) != "d459c24da2f9" {
		log.Fatal("FakeYubikey expects hardcoded challenge...", "d459c24da2f9")
	}
	var resp [16]byte
	r, _ := hex.DecodeString("9e7244e281d1e3b93f1005ba138b8a04")
	copy(resp[:], r)
	return resp, nil
}

func (yk *FakeYubikey) challengeResponseHMAC(chal SHADigest) ([20]byte, error) {
	if hex.EncodeToString(
		chal[:],
	) != "72542b8786762da3178a035eb5f2fcef2d020dd18be729f6f67fa46ee134d5c7" {
		log.Fatal(
			"FakeYubikey expects hardcoded HMAC challenge...",
			"72542b8786762da3178a035eb5f2fcef2d020dd18be729f6f67fa46ee134d5c7",
		)
	}
	var hmac [20]byte
	h, _ := hex.DecodeString("50d30849c0e623f20665267b02fd37f4528f8cf2")
	copy(hmac[:], h)
	return hmac, nil
}

func TestPacketFactory(t *testing.T) {
	entropyBuf, _ := hex.DecodeString(
		"280c6d2ea06d20db231dc93bf16db0b7308016d782c7dfe7c969c08cf68cc984",
	)

	var entropyBufArray [32]byte
	copy(entropyBufArray[:], entropyBuf)
	pk := &PortKnocker{&FakeYubikey{}, entropyBufArray}
	timestamp := int32(1627324072)
	auth_addr := net.ParseIP("74.109.234.77").To4()
	ssh_addr := net.ParseIP("74.109.234.77").To4()
	server_addr := net.ParseIP("204.154.139.11").To4()

	pkt, _ := pk.makePacket(
		"mmeredith",
		timestamp,
		auth_addr,
		ssh_addr,
		server_addr,
	)

	hexstring := hex.EncodeToString(pkt)
	assert.Equal(
		t,
		"a8fefe603daa52002fbf19ccd4d14352a1436dd80a098b9672542b8786762da3178a035eb5f2fcef2d020dd18be729f6f67fa46ee134d5c79b529d8aad7a6067e9a2806af85c2af2a711321150dcbf1d9836fad448c684a8dbca402d15da7c80116078e77c38eecc4a94a83cd244a07258662ce1e04b6a29e0cd6937fb70d7059db5221bc891393b43aa55b2a452e39e5d490b4f27cc0d64ccc974932ce1979e64449e4d4d2d9e3bdd0da91d668039f5b1b6dc6a3ab411f216d9599373226cbc711f184c6f18b97a90b0d31231b4c580822237a8b3204575efc3c5356ac33bc7c16e8456aeb64c45ef6933ca5845a66619ef4338f477d474",
		hexstring,
	)
}

func TestNtp(t *testing.T) {
	time, err := getNtpTime()
	assert.Equal(t, err, nil)
	assert.Less(t, int32(1631907374), time)
	assert.Less(t, time, int32(2100000000)) // Good until 2036
}
