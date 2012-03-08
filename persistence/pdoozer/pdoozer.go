// pdoozer is a persistence client for doozerd.  It is implemented as a
// usual doozer client that connects to a cluster, monitors I/O and writes
// mutations to persistent medium.  To signal its users, pdoozerd maintains
// a clone of the namespace in /ctl/pdoozer/<n>/.  A mutation in the mirrored
// tree signals the successful logging of the associated mutation to disk.
package main

func main() {
}
