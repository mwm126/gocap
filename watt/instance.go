package watt

import (
	"fmt"
	"image/color"
	"log"
	"strings"
	"time"

	"aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type InstanceTab struct {
	TabItem    *container.TabItem
	table      *widget.Table
	instances  map[string][]Instance
	inst_table [][]string
	connection *cap.Connection
	closed     bool
}

type Instance struct {
	UUID  string
	Name  string
	State string
}

func (t *InstanceTab) Close() {
	t.closed = true
}

func NewInstanceTab(conn *cap.Connection) *InstanceTab {
	t := InstanceTab{
		TabItem:    nil,
		table:      nil,
		instances:  make(map[string][]Instance),
		connection: conn,
		closed:     false,
	}
	t.table = widget.NewTable(
		func() (int, int) {
			num_rows := 1 // for header
			for _, insts := range t.instances {
				num_rows += len(insts)
			}
			return num_rows, 4
		},
		func() fyne.CanvasObject {
			obj := canvas.NewText("lorem ipsum", theme.PrimaryColorNamed("yellow"))
			return obj
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row == 0 {
				// give me some header
				txt := map[int]string{
					0: "Project",
					1: "UUID",
					2: "Name",
					3: "State",
				}[i.Col]
				o.(*canvas.Text).Text = txt
				o.(*canvas.Text).TextStyle.Italic = true
				o.(*canvas.Text).Color = theme.PrimaryColorNamed("gray")
				return
			}
			o.(*canvas.Text).Color = color.White
			o.(*canvas.Text).TextStyle.Italic = false
			o.(*canvas.Text).Text = t.inst_table[i.Row-1][i.Col]
		})
	t.table.SetColumnWidth(0, 200)
	t.table.SetColumnWidth(1, 500)
	t.table.SetColumnWidth(2, 200)
	t.table.SetColumnWidth(3, 200)
	t.table.Resize(fyne.NewSize(1000, 1000))

	go func() {
		for !t.closed {
			t.refresh()
		}
	}()

	scroll := container.NewScroll(t.table)
	t.TabItem = container.NewTabItem("Instances", scroll)
	return &t
}

func (t *InstanceTab) refresh() {
	projects, err := t.get_projects()
	if err != nil {
		log.Println("Could not get projects: ", err)
		return
	}
	instmap := make(map[string][]Instance)
	for _, project := range projects {
		instances, err := t.get_instances(project)
		if err != nil {
			log.Printf("Could not refresh instances for project %s: %s", project, err)
			continue
		}
		instmap[project] = instances
	}

	insttab := make([][]string, 0)
	for proj, insts := range t.instances {
		for _, inst := range insts {
			insttab = append(insttab, []string{proj, inst.UUID, inst.Name, inst.State})
		}
	}
	t.instances = instmap
	t.inst_table = insttab
	t.table.Refresh()
	time.Sleep(123 * time.Second) // TODO: configure refresh interval
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
