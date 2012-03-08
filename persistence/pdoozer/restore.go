package main

import (
	"errors"
	"github.com/ha/doozer"
	"io"
	"strconv"
	"strings"
)

const (
	_ = 1 << iota
	_
	set
	del
)

var ErrBadMutation = errors.New("bad mutation")

func restore() {
	for {
		m, err := journal.ReadMutation()
		if err == io.EOF {
			return
		} else if err != nil {
			exit(err)
		}
		err = apply(m)
		if err != nil {
			exit(err)
		}
	}
}

func apply(mut string) error {
	ev, err := decode(mut)
	if err != nil {
		return err
	}
	rev, err := conn.Rev()
	if err != nil {
		return err
	}
	switch {
	case ev.IsSet():
		_, err = conn.Set(ev.Path, rev, ev.Body)
	case ev.IsDel():
		err = conn.Del(ev.Path, rev)
	}
	return err
}

// decode decodes a mutation into an equivalent doozer.Event,
// from store/store.go:/decode.
func decode(mut string) (ev doozer.Event, err error) {
	cm := strings.SplitN(mut, ":", 2)

	if len(cm) != 2 {
		err = ErrBadMutation
		return
	}

	ev.Rev, err = strconv.ParseInt(cm[0], 10, 64)
	if err != nil {
		return
	}

	kv := strings.SplitN(cm[1], "=", 2)

	ev.Path = kv[0]
	switch len(kv) {
	case 1:
		ev.Flag = del
	case 2:
		ev.Body = []byte(kv[1])
		ev.Flag = set
	}
	return
}
