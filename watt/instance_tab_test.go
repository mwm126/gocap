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
	assert.Equal(t, len(insttab.instances), 0)
}

func TestListInstanceTab(t *testing.T) {
	insttab := NewInstanceTab(FakeInstanceLister{
		map[string][]Instance{
			"my_project": []Instance{
				Instance{},
			},
		},
	})

	assert.Equal(t, len(insttab.instances), 1)
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
	})
	assert.Equal(t, 2, len(insttab.filtered_instances))
	assert.Equal(t, 2, insttab.list.Length())

	insttab.filterEntry.SetText("my_")

	assert.Equal(t, 1, len(insttab.filtered_instances))
	assert.Equal(t, 1, insttab.list.Length())
}
