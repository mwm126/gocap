package login

import (
	"testing"
)

func TestParseServicesJson(t *testing.T) {
	InitServices(nil)
	services, err := FindServices()
	if err != nil {
		t.Fatal("failed to find services:", err)
	}

	joule_svc := services[0]
	if joule_svc.Name != "joule" {
		t.Error("incorrect name")
	}
	if joule_svc.CapPort != 62201 {
		t.Error("incorrect port")
	}
	if joule_svc.Networks["external"].CapServerAddress != "204.154.139.11" {
		t.Error("incorrect address")
	}

}
