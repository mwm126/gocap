package cap

import (
	"encoding/hex"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeYubikey struct{}

func (yk *FakeYubikey) FindSerial() (int32, error) {
	return 5417533, nil
}

func (yk *FakeYubikey) challengeResponse(chal [16]byte) ([16]byte, error) {
	var resp [16]byte
	r, _ := hex.DecodeString("2ee619bc248bcefbe09e733d2cdda3be")
	copy(resp[:], r)
	return resp, nil
}

func (yk *FakeYubikey) challengeResponseHMAC(chal SHADigest) ([16]byte, error) {
	var hmac [16]byte
	h, _ := hex.DecodeString("2ee619bc248bcefbe09e733d2cdda3be")
	copy(hmac[:], h)
	return hmac, nil
}

func TestPacketFactory(t *testing.T) {
	entropyBuf, _ := hex.DecodeString(
		"280c6d2ea06d20db231dc93bf16db0b7308016d782c7dfe7c969c08cf68cc984",
	)

	var entropy [32]byte
	copy(entropy[:], entropyBuf)

	pk := &PortKnocker{&FakeYubikey{}, [32]byte(entropy)}
	timestamp := int32(1627324072)
	auth_addr := net.ParseIP("74.109.234.77").To4()
	ssh_addr := net.ParseIP("74.109.234.77").To4()
	server_addr := net.ParseIP("104.154.139.11").To4()

	pkt, _ := pk.makePacket(
		"mmeredith",
		"xUZv!jA&TgHTkw#!3$bUVcDXxW3sY",
		timestamp,
		auth_addr,
		ssh_addr,
		server_addr,
	)
	hexstring := hex.EncodeToString(pkt)

	assert.Equal(
		t,
		len(
			"a8fefe603daa5200df3b4956df8baa7ac19756127f2069d03590162f38bbab498f0b61715302396ae47b3c81b204933f1f89e0745e501fe1935a2c1e5940d17e42dddbc5511371d131c02d0c4db351436fdc783f70a036ceb5b97a5b982233859bd58dff266435111e0e85e9dfd61ce138bec63fb4dc96d5f5ec1b404027200d55267cc82d9576bf30807e03400efcf31385347b7727d8209dd4028c05527266f98c7412f47080e8387fb0d4fb444720519409019bdc58ec7c120c40408d1b483eb38d5276473f9bb1bd0879304fc425729bf3126a5716c6a8df5fa20721b034056d69dc70d97cb8f8c57322abe3d8be8ae7442c538661d4",
		),
		len(hexstring),
	)
	assert.Equal(
		t,
		"a8fefe603daa5200df3b4956df8baa7ac19756127f2069d03590162f38bbab498f0b61715302396ae47b3c81b204933f1f89e0745e501fe1935a2c1e5940d17e42dddbc5511371d131c02d0c4db351436fdc783f70a036ceb5b97a5b982233859bd58dff266435111e0e85e9dfd61ce138bec63fb4dc96d5f5ec1b404027200d55267cc82d9576bf30807e03400efcf31385347b7727d8209dd4028c05527266f98c7412f47080e8387fb0d4fb444720519409019bdc58ec7c120c40408d1b483eb38d5276473f9bb1bd0879304fc425729bf3126a5716c6a8df5fa20721b034056d69dc70d97cb8f8c57322abe3d8be8ae7442c538661d4",
		hexstring,
	)
}

func TestNtp(t *testing.T) {
	time, err := getNtpTime()
	assert.Equal(t, err, nil)
	assert.Less(t, int32(1631907374), time)
	assert.Less(t, time, int32(2100000000)) // Good until 2036
}
