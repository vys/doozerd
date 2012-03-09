// doozerp is a persistence client for doozerd.  It is implemented as ausual
// doozer client that connects to a cluster, monitors I/O and writesmutations
// to pesistent medium.  To signal its users, doozerp maintainsa clone of
// the namespae in /ctl/persistence/<n>/.  A mutation in the mirroredtree
// signals the succesful logging of the associated mutation to disk.
package main

import (
	"flag"
	"fmt"
	"github.com/ha/doozer"
	"github.com/ha/doozerd/persistence"
	"os"
)

var (
	buri = flag.String("b", "", "the DzNS uri")
	f    = flag.Bool("f", false, "try fix a broken journal.")
	j    = flag.String("j", "journal", "file to log mutations")
	r    = flag.Bool("r", false, "restore from file")
	uri  = flag.String("a", "doozer:?ca=127.0.0.1:8046", "the address to bind to")
	v    = flag.Bool("v", false, "print each mutation on stdout")
)

var (
	conn    *doozer.Conn
	id      = 0
	journal *persistence.Journal
	notify  = make(chan *mutation)
	store   = make(chan *mutation)
)

func usage() {
	errln("usage: doozerp [options]")
	flag.PrintDefaults()
	os.Exit(1)
}

func errln(err string) {
	fmt.Fprintln(os.Stderr, "doozerp: "+err)
}

func exit(err error) {
	errln(err.Error())
	os.Exit(2)
}

func dial() {
	var err error
	conn, err = doozer.DialUri(*uri, *buri)
	if err != nil {
		exit(err)
	}
	getid()
}

func getid() {
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var err error
	journal, err = persistence.NewJournal(*j)
	if err != nil {
		exit(err)
	}

	dial()
	if *r {
		restore()
	}
	go Store()
	go Notify()
	monitor()
}
