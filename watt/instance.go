package watt

import (
	"fmt"
	"log"
	"strings"

	"aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type InstanceTab struct {
	TabItem    *container.TabItem
	table      *widget.Table
	instances  map[string][]Instance
	connection *cap.Connection
}

type Instance struct {
	UUID  string
	Name  string
	State string
}

func NewInstanceTab(conn *cap.Connection) *InstanceTab {
	t := InstanceTab{}
	t.connection = conn
	t.table = widget.NewTable(
		func() (int, int) {
			num_rows := 0
			for _, insts := range t.instances {
				num_rows += len(insts)
			}
			return num_rows, 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("lorem ipsum")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			inst_table := make([][]string, 0)
			for proj, insts := range t.instances {
				for _, inst := range insts {
					inst_table = append(inst_table, []string{proj, inst.UUID, inst.Name, inst.State})
				}
			}
			o.(*widget.Label).SetText(inst_table[i.Row][i.Col])
		})
	refresh := widget.NewButton("Refresh", t.refresh)
	box := container.NewVBox(refresh, t.table)
	t.TabItem = container.NewTabItem("Instances", box)
	return &t
}

func (t *InstanceTab) refresh() {
	projects, err := t.get_projects()
	if err != nil {
		log.Println("Could not get projects:", err)
		return
	}
	t.instances = make(map[string][]Instance)
	for _, project := range projects {
		instances, err := t.get_instances(project)
		if err != nil {
			log.Println("Could not refresh instances:", err)
			return
		}
		t.instances[project] = instances
	}
	for name, val := range t.instances {
		log.Println("...............", name, val)
	}
	t.table.Refresh()
}

func (t *InstanceTab) get_projects() ([]string, error) {
	cmd := t.env("openstack project list -f csv", "")
	output, err := t.connection.GetClient().CleanExec(cmd)
	if err != nil {
		log.Println("Command ", cmd, " had error ", err)
		return make([]string, 0), err
	}
	return parseProjects(output), nil
}

func parseProjects(text string) []string {
	items := []string{}
	// skip first (header) line:  "Name     Status"
	for _, line := range strings.Split(text, "\n")[1:] {
		fields := strings.Split(line, ",")
		if len(fields) < 2 {
			continue
		}
		items = append(items, strings.Trim(fields[1], " \""))
	}
	return items
}

const AD_HOSTNAME = "ad.science"

func (t *InstanceTab) env(cmd string, project_name string) string {
	envvars := map[string]string{
		"OS_AUTH_URL":             "http://192.168.101.182:5000/v3",
		"OS_IDENTITY_API_VERSION": "3",
		"OS_PASSWORD":             t.connection.GetPassword(),
		"OS_PROJECT_DOMAIN_NAME":  AD_HOSTNAME,
		"OS_PROJECT_NAME":         project_name,
		"OS_USERNAME":             t.connection.GetUsername(),
		"OS_USER_DOMAIN_NAME":     AD_HOSTNAME,
	}
	for name, value := range envvars {
		cmd = fmt.Sprintf("%s='%s' %s", name, value, cmd)
	}
	return cmd
}

func (t *InstanceTab) get_instances(project_name string) ([]Instance, error) {
	cmd := fmt.Sprintf("openstack server list -f csv")
	cmd = t.env(cmd, project_name)
	output, err := t.connection.GetClient().CleanExec(cmd)
	if err != nil {
		log.Println("Command ", cmd, " had error ", err)
		return make([]Instance, 0), err
	}
	return parseInstances(output), nil
}

func parseInstances(text string) []Instance {
	var instances []Instance
	// skip first (header) line:  "Name     Status"
	for _, line := range strings.Split(text, "\n")[1:] {
		if strings.TrimSpace(text) == "" {
			continue
		}
		fields := strings.Split(strings.TrimSpace(line), ",")
		if len(fields) < 3 {
			log.Println("Skipping instance line:", line)
			continue
		}
		uuid := strings.Trim(fields[0], "\"")
		name := strings.Trim(fields[1], "\"")
		state := strings.Trim(fields[2], "\"")
		instances = append(instances, Instance{uuid, name, state})
	}
	return instances
}
