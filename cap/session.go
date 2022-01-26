package cap

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/data/binding"
)

func parseSessions(username, text string) []Session {
	sessions := make([]Session, 0, 10)
	for _, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		session, err := parseVncLine(line)
		if err != nil {
			continue
		}
		if session.Username == username {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

type Session struct {
	Username      string
	DisplayNumber string
	Geometry      string
	DateCreated   string
	HostPort      uint
}

func (s *Session) Label() string {
	return fmt.Sprintf(
		"Session %s - %s - %s",
		s.DisplayNumber,
		s.Geometry,
		s.DateCreated,
	)
}

func (s Session) AddListener(listener binding.DataListener) {
}

func (s Session) RemoveListener(listener binding.DataListener) {
}

func parseVncLine(line string) (Session, error) {
	var session Session
	fields := strings.Fields(line)
	if len(fields) < 15 {
		return session, errors.New("Parse error")
	}
	username := fields[15][1 : len(fields[15])-1]
	port, err := strconv.Atoi(get_field(fields, "-rfbport"))
	if err != nil {
		return Session{}, err
	}
	session = Session{
		Username:      username,
		DisplayNumber: fields[11],
		Geometry:      get_field(fields, "-geometry"),
		DateCreated:   fields[8],
		HostPort:      uint(port),
	}
	return session, nil
}

func get_field(fields []string, fieldname string) string {
	for ii, field := range fields {
		if field == fieldname {
			return fields[ii+1]
		}
	}
	return ""
}
