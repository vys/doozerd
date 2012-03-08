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
		_ = m // TODO(aram): user persistence to save it.
		notify <- m
	}
}
