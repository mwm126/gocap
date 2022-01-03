package watt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseProjects(t *testing.T) {
	text := `"3880f9d8def44f759324b8881d8fc736","ML_Demo"`

	got := parseProjects(text)

	want := []string{"ML_Demo"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch: %s", diff)
	}
}

func TestParseInstances(t *testing.T) {
	text := `"ID","Name","Status","Networks","Image","Flavor"
	"1224e86d-f1b6-414d-ae3b-5475ebd291e4","hopeitworks","ERROR","","Windows 10","ml.gpu.1.bigmem"
	"8b65a1f3-9706-4653-ba8e-57b4307e5ad8","UbuntuEvan","SHUTOFF","ML_Demo_Network=172.16.1.104","Ubuntu-20.04-desktop.base","ml.tiny" `

	got := parseInstances(text)

	want := []inst{{
		"1224e86d-f1b6-414d-ae3b-5475ebd291e4", "hopeitworks", "ERROR"},
		{"8b65a1f3-9706-4653-ba8e-57b4307e5ad8", "UbuntuEvan", "SHUTOFF"}}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch: %s", diff)
	}
}
