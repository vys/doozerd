package main

import (
	"fmt"
)

func Notify() {
	prefix := "/ctl/persistence/" + fmt.Sprint(id)

	for m := range notify {
		ev := m.ev
		switch {
		case ev.IsSet():
			body := []byte(fmt.Sprint(ev.Rev))
			rev, err := conn.Rev()
			if err != nil {
				exit(err)
			}
			_, err = conn.Set(prefix+ev.Path, rev, body)
			if err != nil {
				exit(err)
			}
		case ev.IsDel():
			rev, err := conn.Rev()
			if err != nil {
				exit(err)
			}
			err = conn.Del(prefix+ev.Path, rev)
			if err != nil {
				exit(err)
			}
		}
	}
}
