package main

import (
	"fmt"
	"github.com/ha/doozer"
	st "github.com/ha/doozerd/store"
	"regexp"
)

var exclude = regexp.MustCompile("^/ctl")

func encode(ev doozer.Event) (mut string) {
	switch {
	case ev.IsSet():
		mut = st.MustEncodeSet(ev.Path, string(ev.Body), ev.Rev)
	case ev.IsDel():
		mut = st.MustEncodeDel(ev.Path, ev.Rev)
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
		// BUG(aram): what are these, NOPs?
		if ev.Flag == 0 {
			continue
		}

		if exclude.MatchString(ev.Path) {
			continue
		}

		mut := encode(ev)
		if *v {
			fmt.Println(mut)
		}

		store <- &mutation{ev: ev, mut: mut}
	}
}

func Notify() {
	prefix := "/ctl/persistence/" + fmt.Sprint(id)

	for m := range notify {
		ev := m.ev
		rev, err := conn.Rev()
		if err != nil {
			exit(err)
		}
		switch {
		case ev.IsSet():
			body := []byte(fmt.Sprint(ev.Rev))
			_, err = conn.Set(prefix+ev.Path, rev, body)
		case ev.IsDel():
			err = conn.Del(prefix+ev.Path, rev)
		}
		if err != nil {
			exit(err)
		}
	}
}
