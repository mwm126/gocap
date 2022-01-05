package joule

import (
	"os/exec"
	"testing"

	"aeolustec.com/capclient/cap"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/google/go-cmp/cmp"
)

type SpyRunner struct {
	calls []exec.Cmd
}

type TestPortFinder struct{}

func (tpf TestPortFinder) FindPort() (int, error) {
	return 54321, nil
}

func TestVncTab(t *testing.T) {
	a := test.NewApp()
	var co fyne.CanvasObject
	w := test.NewWindow(co)

	t.Run("Vnc Refresh Sessions", func(t *testing.T) {
		conn := NewFakeVncConnection(t, map[string]string{
			"hostname": "the_hostname",
			`ping -c 1 the_hostname| grep PING|awk '{print $3}'| sed "s/(//"|sed "s/)//"`: "1.2.3.4",
			`id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'`:                             "the_uid",
			"ps auxnww|grep Xvnc|grep -v grep":                                            `8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:234 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :123 -desktop TurboVNC: login03:5 (the_user) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1`,
		})

		vncTab := newVncTab(a, w, conn, &SpyRunner{}, &TestPortFinder{})

		want := 1
		got := vncTab.sessions.Length()
		if want != got {
			t.Error("# sessions should be 1 but was: ", vncTab.sessions.Length())
		}
	})

	t.Run("Vnc New Session", func(t *testing.T) {
		var conn cap.Connection
		// conn.sessions = []cap.Session{init_session}
		vncTab := newVncTab(a, w, &conn, &SpyRunner{}, &TestPortFinder{})

		want := 0
		got := vncTab.sessions.Length()
		if want != got {
			t.Error("Initially # of sessions should be 0 but was: ", vncTab.sessions.Length())
		}

		test.Tap(vncTab.new_btn) // Should show new session dialog (without error)
	})
}

func TestNewSessionDialog(t *testing.T) {
	a := test.NewApp()
	var co fyne.CanvasObject
	w := test.NewWindow(co)

	default_rezs := []string{
		"800x600",
		"1024x768",
		"1280x1024",
		"1600x1200",
	}

	t.Run("Preset Resolution", func(t *testing.T) {
		conn := NewFakeVncConnection(t, map[string]string{
			"hostname": "the_hostname",
			`ping -c 1 the_hostname| grep PING|awk '{print $3}'| sed "s/(//"|sed "s/)//"`: "1.2.3.4",
			`id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'`:                             "the_uid",
			"ps auxnww|grep Xvnc|grep -v grep":                                            `8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:234 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :123 -desktop TurboVNC: login03:5 (the_user) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1`,
			"vncserver -geometry 1600x1200 -otp -novncauth -nohttpd": `Desktop 'TurboVNC: login03.super:22 (mmeredith)' started on display login03.super:22

			One-Time Password authentication enabled.  Generating initial OTP ...
				Full control one-time password: 22256714
			Run 'vncpasswd -o' from within the TurboVNC session or
			'vncpasswd -o -display login03.super:22' from within this shell
			to generate additional OTPs
			Starting applications specified in /nfs/home/3/mmeredith/.vnc/xstartup.turbovnc
			Log file is /nfs/home/3/mmeredith/.vnc/login03.super:22.log`,
		})
		vncTab := newVncTab(a, w, conn, &SpyRunner{}, &TestPortFinder{})
		vsf := vncTab.NewVncSessionForm(test.NewWindow(nil), default_rezs)
		last_index := len(vsf.preset_select.Options) - 1
		vsf.preset_select.SetSelectedIndex(last_index)

		vsf.Form.OnSubmit()

		want := []cap.Session{
			{
				Username:      "the_user",
				DisplayNumber: ":123",
				Geometry:      "3840x2160",
				DateCreated:   "Aug03",
				HostPort:      5905}}
		got, err := conn.FindSessions()
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})

	t.Run("Custom Resolution", func(t *testing.T) {

		conn := NewFakeVncConnection(t, map[string]string{
			"hostname": "the_hostname",
			`ping -c 1 the_hostname| grep PING|awk '{print $3}'| sed "s/(//"|sed "s/)//"`: "1.2.3.4",
			`id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'`:                             "the_uid",
			"ps auxnww|grep Xvnc|grep -v grep":                                            `8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:234 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :123 -desktop TurboVNC: login03:5 (the_user) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1`,
			"vncserver -geometry 999x555 -otp -novncauth -nohttpd": `Desktop 'TurboVNC: login03.super:22 (mmeredith)' started on display login03.super:22

			One-Time Password authentication enabled.  Generating initial OTP ...
				Full control one-time password: 22256714
			Run 'vncpasswd -o' from within the TurboVNC session or
			'vncpasswd -o -display login03.super:22' from within this shell
			to generate additional OTPs
			Starting applications specified in /nfs/home/3/mmeredith/.vnc/xstartup.turbovnc
			Log file is /nfs/home/3/mmeredith/.vnc/login03.super:22.log`,
		})

		vncTab := newVncTab(a, w, conn, &SpyRunner{}, &TestPortFinder{})
		vsf := vncTab.NewVncSessionForm(test.NewWindow(nil), default_rezs)
		vsf.xres_entry.SetText("999")
		vsf.yres_entry.SetText("555")

		vsf.Form.OnSubmit()

		want := []cap.Session{
			{
				Username:      "the_user",
				DisplayNumber: ":123",
				Geometry:      "3840x2160",
				DateCreated:   "Aug03",
				HostPort:      5905}}

		got, err := conn.FindSessions()
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})
}

func TestVncConnect(t *testing.T) {
	var co fyne.CanvasObject
	w := test.NewWindow(co)

	conn := NewFakeVncConnection(t, map[string]string{
		"hostname": "the_hostname",
		`ping -c 1 the_hostname| grep PING|awk '{print $3}'| sed "s/(//"|sed "s/)//"`: "1.2.3.4",
		`id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'`:                             "the_uid",
		"ps auxnww|grep Xvnc|grep -v grep":                                            `8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:234 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :123 -desktop TurboVNC: login03:5 (the_user) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1`,
		"vncpasswd -o -display 1.2.3.4:123":                                           "Full control one-time password: 17760704",
	})
	vncTab := newVncTab(test.NewApp(), w, conn, &SpyRunner{}, &TestPortFinder{})
	vsf := vncTab.NewVncSessionForm(test.NewWindow(nil), make([]string, 0))

	vncTab.refresh()

	last_index := len(vsf.preset_select.Options) - 1
	vsf.preset_select.SetSelectedIndex(last_index)
	id2, err := vncTab.sessions.GetItem(0)
	if err != nil {
		t.Error(err)
	}
	id := id2.(cap.Session)
	obj := vncTab.list_items[id].Objects[0]
	connect_btn := obj.(*ItemButton)

	test.Tap(connect_btn)

	want := "/path/to/vncviewer 127.0.0.1::54321 -Password=17760704"

	got := vncTab.VncRunner.(*SpyRunner).calls[0].String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch: %s", diff)
	}
}

func TestVncCmd(t *testing.T) {
	cmd := VncCmd("/path/to/vncviewer", "xyz", 123)

	want := "/path/to/vncviewer 127.0.0.1::123 -Password=xyz"
	got := cmd.String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch: %s", diff)
	}
}

func (r *SpyRunner) RunVnc(otp, display string, port int) {
	r.calls = append(r.calls, *VncCmd("/path/to/vncviewer", otp, port))
}
