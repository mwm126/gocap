package watt

import (
	"testing"

	// "aeolustec.com/capclient/cap"
	"github.com/stretchr/testify/assert"
)

func TestConnInstanceTab(t *testing.T) {
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

	// insttab.list.Objects.SetText("my_")

	// assert.Equal(t, 1, len(insttab.filtered_instances))
	// assert.Equal(t, 1, insttab.list.Length())
}
