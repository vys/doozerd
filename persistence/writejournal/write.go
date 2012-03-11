package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ha/doozerd/persistence"
	"os"
)

var progname = path.Base(os.Args[0])

func usage() {
	errln("usage: " + progname + " journal")
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
		errln(err)
		os.Exit(2)
	}
	r := bufio.NewReader(os.Stdin)
	for {
		line, prefix, err := r.ReadLine()
		if err != nil {
			return
		}
		for prefix {
			var l []byte
			l, prefix, _ = r.ReadLine()
			if err != nil {
				return
			}
			line = append(line, l...)
		}
		err = j.WriteMutation(string(line))
		if err != nil {
			errln(err)
			os.Exit(4)
		}
	}
}
