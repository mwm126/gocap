package client

import (
	"testing"
)

func TestParseServicesJson(t *testing.T) {
	services, err := FindServices()
	if err != nil {
		t.Fatal("failed to find services:", err)
	}

	if services[0].Name != "joule" {
		t.Error("incorrect")
	}

}
