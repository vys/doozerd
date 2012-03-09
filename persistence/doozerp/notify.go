package main

import (
	"fmt"
)

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
