// pdoozer is a persistence client for doozerd.  It is implemented as a
// usual doozer client that connects to a cluster, monitors I/O and writes
// mutations to persistent medium.  To signal its users, pdoozerd maintains
// a clone of the namespace in /ctl/pdoozer/<n>/.  A mutation in the mirrored
// tree signals the successful logging of the associated mutation to disk.
package main

import (
	"flag"
	"fmt"
	"github.com/ha/doozer"
	"os"
)

var (
	uri     = flag.String("a", "doozer:?ca=127.0.0.1:8046", "the address to bind to")
	buri    = flag.String("b", "", "the DzNS uri")
	journal = flag.String("j", "journal", "file to log mutations")
)

func usage() {
	errln("usage: pdoozer [options]")
	flag.PrintDefaults()
	os.Exit(1)
}

func errln(err string) {
	fmt.Fprintln(os.Stderr, "pdoozer: "+err)
}

func exit(err string) {
	errln(err)
	os.Exit(2)
}

func errexit(err error) {
	exit(err.Error())
}

// connection to the cluster.
var conn *doozer.Conn

func main() {
	flag.Usage = usage
	flag.Parse()
	if *uri == "" && *buri == "" {
		errln("pdoozer: either -a or -b must be specified.")
		usage()
	}
	if *journal == "" {
		errln("pdoozer: must use -j to journal file.")
		usage()
	}

	var err error
	if *buri != "" {
		exit("pdoozer: -b not implemented")
	} else {
		conn, err = doozer.Dial(*uri)
		if err != nil {
			errexit(err)
		}
	}
}
