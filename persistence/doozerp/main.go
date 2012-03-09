// doozerp is a persistence client for doozerd.  It is implemented as
// a regular doozer client that connects to a cluster, monitors I/O and
// writes mutations to pesistent medium.  To signal its users, doozerp
// maintains a clone of the namespace in /ctl/persistence/<n>/.  A mutation
// in the mirrored tree signals the succesful logging of the associated
// mutation to disk.
package main

import (
	"flag"
	"fmt"
	"github.com/ha/doozer"
	"github.com/ha/doozerd/persistence"
	"os"
	"strconv"
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
	conn    *doozer.Conn           // connection to the cluster.
	id      = 0                    // client id.
	journal *persistence.Journal   // journal to log to.
	notify  = make(chan *mutation) // ack. write operation.
	store   = make(chan *mutation) // save to disk.
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

// dial connects to the server.
func dial() {
	var err error
	conn, err = doozer.DialUri(*uri, *buri)
	if err != nil {
		exit(err)
	}
	setid()
}

// setid determines and sets the client id.
func setid() {
	body, rev, err := conn.Get("/ctl/persistence/id", nil)
	id, _ = strconv.Atoi(string(body))
	id++
	body = []byte(fmt.Sprint(id))
	_, err = conn.Set("/ctl/persistence/id", rev, body)
	if err != nil {
		exit(err)
	}
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
