package event

import (
	"regexp"
	"strings"
)

var goodNameReg = regexp.MustCompile(`^[a-zA-Z][\w-.*]*$`)

func goodName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		panic("event: the event name cannot be empty")
	}

	if !goodNameReg.MatchString(name) {
		panic(`event: the event name is invalid, must match regex '^[a-zA-Z][\w-.]*$'`)
	}

	return name
}
