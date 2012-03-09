package main

import (
	"github.com/ha/doozer"
)

type mutation struct {
	mut string
	ev  doozer.Event
}

func Store() {
	for m := range store {
		err := journal.WriteMutation(m.mut)
		if err != nil {
			exit(err)
		}
		notify <- m
	}
}
