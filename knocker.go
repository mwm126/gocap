package main

import (
// "fmt"
)

type Knocker interface {
	Knock(uname string, pword string)
}

type PortKnocker struct {
}

func (sk *PortKnocker) Knock(uname string, pword string) {
}

func (sk *PortKnocker) packet() string {
	return "TODO"
}
