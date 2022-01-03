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
	TabItem   *container.TabItem
	table     *widget.Table
	instances map[string][]inst
	client    cap.Client
}

type inst struct {
	UUID  string
	Name  string
	State string
}

func NewInstanceTab(client cap.Client) *InstanceTab {
	t := InstanceTab{}
	t.client = client
	t.table = widget.NewTable(
		func() (int, int) {
			return len(t.instances), 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("lorem ipsum")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			insts := make([][]string, 0)
			for proj, instss := range t.instances {
				for _, instt := range instss {

					insts = append(insts, []string{proj, instt.UUID, instt.Name, instt.State})
				}
			}
			o.(*widget.Label).SetText(insts[i.Row][i.Col])
		})
	refresh := widget.NewButton("Refresh", t.refresh)
	box := container.NewVBox(refresh, t.table)
	t.TabItem = container.NewTabItem("Instances", box)
	return &t
}

func (t *InstanceTab) refresh() {
	projects, err := get_projects(t.client)
	if err != nil {
		log.Println("Could not get projects:", err)
		return
	}
	t.instances = make(map[string][]inst)
	for _, project := range projects {
		instances, err := get_instances(t.client, project)
		if err != nil {
			log.Println("Could not refresh instances:", err)
			return
		}
		t.instances[project] = instances
	}
}

func get_projects(client cap.Client) ([]string, error) {
	cmd := "openstack project list -f csv"
	output, err := client.CleanExec(cmd)
	if err != nil {
		log.Println("Command ", cmd, " had error ", err)
		return make([]string, 0), err
	}
	return parseProjects(output), nil
}

func parseProjects(text string) []string {
	items := []string{}
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Split(line, ",")
		if len(fields) < 2 {
			continue
		}
		items = append(items, strings.Trim(fields[1], " \""))
	}
	return items
}

func get_instances(client cap.Client, project_name string) ([]inst, error) {
	cmd := fmt.Sprintf("env OS_PROJECT_NAME=%s openstack server list -f csv", project_name)
	log.Println("cmd:", cmd)
	output, err := client.CleanExec(cmd)
	log.Println("output: ", output)
	if err != nil {
		log.Println("Command ", cmd, " had error ", err)
		return make([]inst, 0), err
	}
	return parseInstances(output), nil
}

func parseInstances(text string) []inst {
	var instances []inst
	// skip first (header) line:  "Name     Status"
	for _, line := range strings.Split(text, "\n")[1:] {
		fields := strings.Split(strings.TrimSpace(line), ",")
		uuid := strings.Trim(fields[0], "\"")
		name := strings.Trim(fields[1], "\"")
		state := strings.Trim(fields[2], "\"")
		instances = append(instances, inst{uuid, name, state})
	}
	log.Println("GOT:", instances)
	return instances
}
