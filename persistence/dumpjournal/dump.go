package main

import (
	"flag"
	"fmt"
	"github.com/ha/doozerd/persistence"
	"io"
	"os"
	"path"
)

var (
	f        = flag.Bool("f", false, "try to fix a broken journal")
	progname = path.Base(os.Args[0])
)

func usage() {
	errln("usage: " + progname + " [options] journal")
	flag.PrintDefaults()
	os.Exit(1)
}

func errln(s string) {
	fmt.Fprintln(os.Stderr, progname+": "+s)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}

	file := flag.Arg(0)
	j, err := persistence.NewJournal(file)
	if err != nil {
		panic(err)
	}
	for {
		m, err := j.ReadMutation()
		if err == io.EOF {
			break
		}
		if err != nil {
			errln("bad journal file: " + err.Error())
			if *f {
				err = j.Fsck()
				if err != nil {
					errln("can't fix journal")
					os.Exit(2)
				}
				errln("journal successfully fixed")
				continue
			} else {
				return
			}
		}
		fmt.Println(m)
	}
}
