package main

import (
	"fmt"
	"github.com/ha/doozer"
	"github.com/ha/doozerd/store"
)

func encode(ev doozer.Event) (mut string) {
	switch {
	case ev.IsSet():
		mut = store.MustEncodeSet(ev.Path, string(ev.Body), ev.Rev)
	case ev.IsDel():
		mut = store.MustEncodeDel(ev.Path, ev.Rev)
	}
	return
}

func monitor() {
	var (
		ev  doozer.Event
		err error
		rev int64 = -1
	)
	for {
		ev, err = conn.Wait("/**", rev)
		if err != nil {
			exit(err)
		}
		rev = ev.Rev + 1
		fmt.Println(encode(ev))
	}
}
