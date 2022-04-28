package watt

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeInstanceLister struct {
	instances map[string][]Instance
}

func (lister FakeInstanceLister) find_instances() map[string][]Instance {
	return lister.instances
}

type FakeSpiceClient struct {
	called_with string
}

func NewFakeSpiceClient() SpiceClient {
	return &FakeSpiceClient{}
}

func (client *FakeSpiceClient) connect(inst Instance) (uint, error) {
	return 0, nil
}

func (spice *FakeSpiceClient) RunSpice(localPort uint) {
	cmd, _ := SpiceCmd(localPort)
	spice.called_with = cmd.String()
}

func TestEmptyInstanceTab(t *testing.T) {
	insttab := NewInstanceTab(FakeInstanceLister{}, NewFakeSpiceClient())
	assert.Equal(t, 0, len(insttab.instances))
}

func TestListInstanceTab(t *testing.T) {
	insttab := NewInstanceTab(FakeInstanceLister{
		map[string][]Instance{
			"my_project": []Instance{
				Instance{},
			},
		},
	}, NewFakeSpiceClient())
	insttab.refresh()

	assert.Equal(t, 1, len(insttab.instances))
}

func TestFilterInstanceTab(t *testing.T) {
	insttab := NewInstanceTab(FakeInstanceLister{
		map[string][]Instance{
			"my_project": []Instance{
				Instance{"my_project", "12345", "my_instance", "RUNNING"},
			},
			"other_project": []Instance{
				Instance{"other_project", "99999", "other_instance", "OFF"},
			},
		},
	}, NewFakeSpiceClient())
	insttab.refresh()
	assert.Equal(t, 2, len(insttab.filtered_instances))
	assert.Equal(t, 2, insttab.list.Length())

	insttab.filterEntry.SetText("my_")

	assert.Equal(t, 1, len(insttab.filtered_instances))
	assert.Equal(t, 1, insttab.list.Length())
}

func TestConnectInstanceTab(t *testing.T) {
	fake_client := NewFakeSpiceClient()
	insttab := NewInstanceTab(FakeInstanceLister{
		map[string][]Instance{
			"my_project": []Instance{
				Instance{"my_project", "12345", "my_instance", "RUNNING"},
			},
			"other_project": []Instance{
				Instance{"other_project", "99999", "other_instance", "OFF"},
			},
		},
	}, fake_client)
	insttab.refresh()
	assert.Equal(t, 2, len(insttab.filtered_instances))
	assert.Equal(t, 2, insttab.list.Length())

	insttab.buttons[1].OnTapped()

	assert.Equal(t,
		fake_client.(*FakeSpiceClient).called_with,
		exec.Command("spicy", "-h", "127.0.0.1", "-p", "0", "-s", "0").String())
}
