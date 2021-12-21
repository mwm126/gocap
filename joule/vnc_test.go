package joule

import (
	"testing"

	"aeolustec.com/capclient/cap"
	"fyne.io/fyne/v2/test"
	"github.com/google/go-cmp/cmp"
)

type SpyRunner struct {
	calls []string
}

func (r *SpyRunner) Run(conn *cap.Connection, otp, display string) {
	r.calls = append(r.calls, VncCmd(otp, display))
}

func _TestVncTab(t *testing.T) {
	a := test.NewApp()

	t.Run("Vnc Refresh Sessions", func(t *testing.T) {
		conn := NewFakeVncConnection(t)
		vncTab := newVncTab(a, conn, &SpyRunner{})

		want := 0
		got := vncTab.sessions.Length()
		if want != got {
			t.Error("Initially # of sessions should be 0 but was: ", vncTab.sessions.Length())
		}

		test.Tap(vncTab.refresh_btn)

		want = 1
		got = vncTab.sessions.Length()
		if want != got {
			t.Error("After refresh # of sessions should be 1 but was: ", vncTab.sessions.Length())
		}
	})

	t.Run("Vnc New Session", func(t *testing.T) {
		var conn cap.Connection
		// conn.sessions = []cap.Session{init_session}
		vncTab := newVncTab(a, &conn, &SpyRunner{})

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

	default_rezs := []string{
		"800x600",
		"1024x768",
		"1280x1024",
		"1600x1200",
	}

	t.Run("Test Preset Resolution", func(t *testing.T) {
		conn := NewFakeVncConnection(t)
		vncTab := newVncTab(a, conn, &SpyRunner{})
		w := test.NewWindow(nil)
		vsf := vncTab.NewVncSessionForm(w, default_rezs)
		last_index := len(vsf.preset_select.Options) - 1
		vsf.preset_select.SetSelectedIndex(last_index)

		vsf.Form.OnSubmit()

		want := []cap.Session{
			{
				Username:      "the_user",
				DisplayNumber: ":123",
				Geometry:      "3840x2160",
				DateCreated:   "Aug03",
				HostAddress:   "localhost",
				HostPort:      "5905"}}
		got, err := conn.FindSessions()
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})

	t.Run("Test Custom Resolution", func(t *testing.T) {
		conn := NewFakeVncConnection(t)
		vncTab := newVncTab(a, conn, &SpyRunner{})
		w := test.NewWindow(nil)
		vsf := vncTab.NewVncSessionForm(w, default_rezs)
		vsf.xres_entry.SetText("999")
		vsf.yres_entry.SetText("555")

		vsf.Form.OnSubmit()

		want := []cap.Session{
			{
				Username:      "the_user",
				DisplayNumber: ":123",
				Geometry:      "3840x2160",
				DateCreated:   "Aug03",
				HostAddress:   "localhost",
				HostPort:      "5905"}}

		got, err := conn.FindSessions()
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})

}

func _TestVncCmd(t *testing.T) {
	cmd := VncCmd("abc", "xyz")

	want := "echo abc | env -u LD_LIBRARY_PATH vncviewer_HPCEE -highqual -autopass 127.0.0.1::10055 &"
	got := cmd
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch: %s", diff)
	}
}

func TestVncConnect(t *testing.T) {
	a := test.NewApp()
	conn := NewFakeVncConnection(t)
	vncTab := newVncTab(a, conn, &SpyRunner{})
	w := test.NewWindow(nil)
	vsf := vncTab.NewVncSessionForm(w, make([]string, 0))

	vncTab.refresh()
	vncTab.List.Refresh()

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

	want := "echo test_get_shared_otp | env -u LD_LIBRARY_PATH vncviewer_HPCEE -highqual -autopass 127.0.0.1::10055 &"

	got := vncTab.VncRunner.(*SpyRunner).calls[0]
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch: %s", diff)
	}
}
