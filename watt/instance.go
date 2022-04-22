package watt

import (
	"fmt"
	"log"
	"strings"

	"aeolustec.com/capclient/cap"
)

type Instance struct {
	UUID  string
	Name  string
	State string
}

func find_instances(conn *cap.Connection) (map[string][]Instance, error) {
	instmap := make(map[string][]Instance)
	projects, err := get_projects(conn)
	if err != nil {
		log.Println("Could not get projects: ", err)
		return instmap, err
	}
	log.Printf("Refreshing:   found %d projects.\n", len(projects))
	for _, project := range projects {
		instances, err := get_instances(conn, project)
		if err != nil {
			log.Printf("Could not refresh instances for project %s: %s", project, err)
			continue
		}
		instmap[project] = instances
	}
	return instmap, nil
}

func filter_instances(instmap map[string][]Instance, txt string) [][]string {
	insttab_filtered := make([][]string, 0)
	for proj, insts := range instmap {
		// instmap_filtered[proj] = make([]Instance, 0)
		for _, inst := range insts {
			txt = strings.ToLower(txt)
			for _, field := range []string{proj, inst.UUID, inst.Name, inst.State} {
				if strings.Contains(strings.ToLower(field), txt) {
					insttab_filtered = append(insttab_filtered, []string{proj, inst.UUID, inst.Name, inst.State})

					break
				}
			}
		}
	}
	return insttab_filtered
}

func get_projects(conn *cap.Connection) ([]string, error) {
	cmd := env(conn, "openstack project list -f csv", "")
	output, err := conn.GetClient().CleanExec(cmd)
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

func env(conn *cap.Connection, cmd string, project_name string) string {
	envvars := map[string]string{
		"OS_AUTH_URL":             "http://192.168.101.182:5000/v3",
		"OS_IDENTITY_API_VERSION": "3",
		"OS_PASSWORD":             conn.GetPassword(),
		"OS_PROJECT_DOMAIN_NAME":  AD_HOSTNAME,
		"OS_PROJECT_NAME":         project_name,
		"OS_USERNAME":             conn.GetUsername(),
		"OS_USER_DOMAIN_NAME":     AD_HOSTNAME,
	}
	for name, value := range envvars {
		cmd = fmt.Sprintf("%s='%s' %s", name, value, cmd)
	}
	return cmd
}

func get_instances(conn *cap.Connection, project_name string) ([]Instance, error) {
	cmd := "openstack server list -f csv"
	cmd = env(conn, cmd, project_name)
	output, err := conn.GetClient().CleanExec(cmd)
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
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		fields := strings.Split(trimmed, ",")
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
