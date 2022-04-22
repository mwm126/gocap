package watt

import (
	"testing"

	// "aeolustec.com/capclient/cap"
	"github.com/stretchr/testify/assert"
)

type FakeInstanceLister struct {
	instances map[string][]Instance
}

func (lister FakeInstanceLister) find_instances() map[string][]Instance {
	return lister.instances
}

func TestEmptyInstanceTab(t *testing.T) {
	insttab := NewInstanceTab(FakeInstanceLister{})
	assert.Equal(t, len(insttab.inst_table), 0)
}

func TestListInstanceTab(t *testing.T) {
	insttab := NewInstanceTab(FakeInstanceLister{
		map[string][]Instance{
			"my_project": []Instance{
				Instance{},
			},
		},
	})

	assert.Equal(t, len(insttab.inst_table), 1)
}

func TestFilterInstanceTab(t *testing.T) {
	insttab := NewInstanceTab(FakeInstanceLister{
		map[string][]Instance{
			"my_project": []Instance{
				Instance{"12345", "my_instance", "RUNNING"},
			},
			"other_project": []Instance{
				Instance{"99999", "other_instance", "OFF"},
			},
		},
	})
	assert.Equal(t, len(insttab.inst_table), 2)

	insttab.filterEntry.SetText("my_")

	assert.Equal(t, len(insttab.inst_table), 1)
}
