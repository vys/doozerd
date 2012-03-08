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
	uri  = flag.String("a", "", "the address to bind to")
	buri = flag.String("b", "", "the DzNS uri")
	journal = flag.String("j", "", "file to log mutations")
)

func usage() {
	fmt.Fprintln(os.Stderr, "pdoozer: usage: pdoozer [options] -j journal")
	flag.PrintDefaults()
	os.Exit(1)
}

var conn *doozer.Conn

func main() {
	flag.Usage = usage
	flag.Parse()
	if *uri == "" && *buri == "" {
		fmt.Fprintln(os.Stderr, "pdoozer: either -a or -b must be specified.")
		usage()
	}
	if *journal == "" {
		fmt.Fprintln(os.Stderr, "pdoozer: must use -j to journal file.")
		usage()
	}
}
