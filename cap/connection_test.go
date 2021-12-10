//go:build integration
// +build integration

package cap

import (
	"github.com/google/go-cmp/cmp"
	"net"
	"testing"
)

type StubYubikey struct{}

func (yk *StubYubikey) FindSerial() (int32, error) {
	return 0, nil
}

func (yk *StubYubikey) ChallengeResponse(chal [6]byte) ([16]byte, error) {
	return [16]byte{}, nil
}

func (yk *StubYubikey) ChallengeResponseHMAC(chal SHADigest) ([20]byte, error) {
	return [20]byte{}, nil
}

func NewFakeKnocker() *PortKnocker {
	var fake_yk StubYubikey
	var entropy [32]byte
	return NewPortKnocker(&fake_yk, entropy)
}

func DisabledTestCapConnection(t *testing.T) {
	username := "testusername"
	password := "testpassword"
	ext_ip := net.IPv4(11, 22, 33, 44)
	server := net.IPv4(55, 66, 77, 88)

	fake_kckr := NewFakeKnocker()
	conn_man := NewCapConnectionManager(fake_kckr)
	ch := make(chan string)
	err := conn_man.Connect(
		username,
		password,
		ext_ip,
		server,
		123,
		func(pwc PasswordChecker) {},
		ch,
	)
	if err != nil {
		t.Error("failed to make cap connection:", err)
	}

	t.Run("Test connection username", func(t *testing.T) {
		want := "0estusername"
		got := conn_man.connection.username
		if want != got {
			t.Errorf("Did not set connection username: want %s but got %s", want, got)
		}
	})
	t.Run("Test connection password", func(t *testing.T) {
		want := "0estpassword"
		got := conn_man.connection.password
		if want != got {
			t.Errorf("Did not set connection password: want %s but got %s", want, got)
		}
	})

}

const ps_output = `
  8048  7536 12.8  2.1 2131820 2111176 ?     R    Aug02 21562:34 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :2 -desktop TurboVNC: login03:2 (pstrakey) -auth /nfs/home/3/pstrakey/.Xauthority -dontdisconnect -geometry 1852x1000 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5902 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8103 10480  0.8  1.2 1701664 1203956 ?     S    Aug02 1395:33 /usr/bin/Xvnc :3 -desktop TurboVNC: login03:3 (dietikej) -auth /nfs/home/3/dietikej/.Xauthority -dontdisconnect -geometry 1920x1080 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5903 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8339 11226  0.0  0.0  38200 16732 ?        S    Sep29   0:19 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :19 -desktop TurboVNC: login03:19 (holcombp) -auth /nfs/home/3/holcombp/.Xauthority -dontdisconnect -geometry 1280x1024 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5919 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8173 11988  5.9  0.2 298968 278716 ?       S    Nov09 1499:34 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :25 -desktop TurboVNC: login03:25 (nkonan) -auth /nfs/home/3/nkonan/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5925 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8134 15223  0.1  0.1 213512 179844 ?       S    Sep21 124:30 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :23 -desktop TurboVNC: login03:23 (vanessed) -auth /nfs/home/3/vanessed/.Xauthority -dontdisconnect -geometry 1280x1024 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5923 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8172 22086  0.0  0.2 290328 268776 ?       S    Aug09 108:15 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :16 -desktop TurboVNC: login03:16 (liuy) -auth /nfs/home/3/liuy/.Xauthority -dontdisconnect -geometry 1920x1200 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5916 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8172 22254  0.3  0.1 171008 150392 ?       S    Aug09 622:13 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :17 -desktop TurboVNC: login03:17 (liuy) -auth /nfs/home/3/liuy/.Xauthority -dontdisconnect -geometry 1920x1200 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5917 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8172 22762  0.0  0.1 156428 135844 ?       S    Sep08  32:58 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :15 -desktop TurboVNC: login03:15 (liuy) -auth /nfs/home/3/liuy/.Xauthority -dontdisconnect -geometry 1920x1200 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5915 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8310 23162  0.0  0.1 182488 161676 ?       S    Aug04   6:56 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :11 -desktop TurboVNC: login03:11 (wux) -auth /nfs/home/3/wux/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5911 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
   12012 24304  0.0  0.0 107704 86708 ?        S    Aug03  59:59 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :10 -desktop TurboVNC: login03:10 (wingop) -auth /nfs/home/3/wingop/.Xauthority -dontdisconnect -geometry 1920x1080 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5910 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8103 24970  0.2  0.1 190416 159008 ?       S    Nov17  32:42 /usr/bin/Xvnc :27 -desktop TurboVNC: login03:27 (dietikej) -auth /nfs/home/3/dietikej/.Xauthority -dontdisconnect -geometry 3440x1440 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5927 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
   10103 25055  0.0  0.0  51124 29444 ?        S    Oct14   0:13 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :20 -desktop TurboVNC: login03:20 (harbertw) -auth /nfs/home/3/harbertw/.Xauthority -dontdisconnect -geometry 1920x1200 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5920 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8173 25384  0.0  0.3 328236 304876 ?       S    Oct04  59:53 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :4 -desktop TurboVNC: login03:4 (nkonan) -auth /nfs/home/3/nkonan/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5904 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8172 25791  0.0  0.0  57344 36404 ?        S    Oct11   1:14 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :18 -desktop TurboVNC: login03:18 (liuy) -auth /nfs/home/3/liuy/.Xauthority -dontdisconnect -geometry 1920x1200 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5918 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:20 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :5 -desktop TurboVNC: login03:5 (mmeredith) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8048 28589  0.5  0.5 573512 553068 ?       S    Aug03 836:45 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :6 -desktop TurboVNC: login03:6 (pstrakey) -auth /nfs/home/3/pstrakey/.Xauthority -dontdisconnect -geometry 1852x1000 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5906 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8328 29085  0.3  0.9 968064 927944 ?       S    Aug03 576:44 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :7 -desktop TurboVNC: login03:7 (nandit) -auth /nfs/home/3/nandit/.Xauthority -dontdisconnect -geometry 1280x1024 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5907 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    9004 29368  0.0  0.0  79388 57804 ?        S    Aug07   7:02 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :13 -desktop TurboVNC: login03:13 (roya) -auth /nfs/home/3/roya/.Xauthority -dontdisconnect -geometry 1368x768 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5913 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8310 30177  0.2  0.6 633008 610196 ?       S    Aug05 330:01 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :8 -desktop TurboVNC: login03:8 (wux) -auth /nfs/home/3/wux/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5908 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8227 30238  0.0  0.0  70724 49284 ?        S    Sep10   1:07 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :22 -desktop TurboVNC: login03:22 (mmeredith) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5922 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8151 32488  0.1  0.1 124656 104396 ?       S    Aug12 170:30 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :21 -desktop TurboVNC: login03:21 (jweber) -auth /nfs/home/3/jweber/.Xauthority -dontdisconnect -geometry 1920x1080 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5921 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8173 33458 24.4  0.3 394972 375172 ?       S    Nov12 5043:34 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :26 -desktop TurboVNC: login03:26 (nkonan) -auth /nfs/home/3/nkonan/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5926 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
   10102 36361  0.0  0.0  47256 25756 ?        S    Oct29   0:09 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :24 -desktop TurboVNC: login03:24 (kangm) -auth /nfs/home/3/kangm/.Xauthority -dontdisconnect -geometry 1920x1080 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5924 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    9004 36622  0.2  0.4 495600 475208 ?       S    Aug06 462:43 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :12 -desktop TurboVNC: login03:12 (roya) -auth /nfs/home/3/roya/.Xauthority -dontdisconnect -geometry 1368x768 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5912 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    8328 39465  0.1  0.0  78548 48608 ?        S    Sep02 185:59 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :1 -desktop TurboVNC: login03:1 (nandit) -auth /nfs/home/3/nandit/.Xauthority -dontdisconnect -geometry 1280x1024 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5901 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
    9004 40930  0.0  0.1 146620 125996 ?       S    Aug03  84:37 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :9 -desktop TurboVNC: login03:9 (roya) -auth /nfs/home/3/roya/.Xauthority -dontdisconnect -geometry 1680x1050 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5909 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
`

func TestParseVncProcesses(t *testing.T) {
	sessions := parseSessions(
		"mmeredith",
		`8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:20 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :5 -desktop TurboVNC: login03:5 (not_mark) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:20 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :5 -desktop TurboVNC: login03:5 (mmeredith) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:20 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :5 -desktop TurboVNC: login03:5 (meredithm) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1`,
	)

	want := Session{
		"mmeredith",
		":5",
		"3840x2160",
		"Aug03",
		"localhost",
		"5905",
	}
	got := sessions[0]
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch: %s", diff)
	}
}
