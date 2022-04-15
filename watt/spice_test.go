package watt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetSpiceCmd(t *testing.T) {
	got, err := get_spice_cmd(12345)
	if err != nil {
		t.Errorf("Error in get_spice_cmd(%d)", 12345)
	}

	want := "env -u LD_LIBRARY_PATH /usr/bin/remote-viewer spice://127.0.0.1:12345"
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch: %s", diff)
	}
}
