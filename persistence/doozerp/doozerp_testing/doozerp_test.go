// Package doozerp_testing tests the doozerp command.
package doozerp_testing

import (
	"errors"
	"github.com/ha/doozer"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"
)

type Cluster struct {
	conn    *doozer.Conn
	doozerd *exec.Cmd
	doozerp *exec.Cmd
	j       string
}

func NewCluster(t *testing.T, doozerpArgs ...string) *Cluster {
	doozerd := exec.Command("doozerd", "-l=127.0.0.1:19999", "-w=false")
	err := doozerd.Start()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	var doozerdIsDead bool
	go func() {
		doozerd.Wait()
		doozerdIsDead = true
	}()
	if doozerdIsDead {
		t.Fatal(errors.New("doozerd died prematurely"))
	}
	
	conn, err := doozer.DialUri("doozer:?ca=127.0.0.1:19999", "")
	if err != nil {
		doozerd.Process.Kill()
		t.Fatal(err)
	}
	
	f, err := ioutil.TempFile("", "j")
	if err != nil {
		conn.Close()
		doozerd.Process.Kill()
		t.Fatal(err)
	}
	j := f.Name()
	f.Close()
	args := append(doozerpArgs, "-a=doozer:?ca=127.0.0.1:19999")
	args = append(args, "-j="+j)
	doozerp := exec.Command("doozerp", args...)
	err = doozerp.Start()
	if err != nil {
		conn.Close()
		doozerd.Process.Kill()
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	var doozerpIsDead bool
	go func() {
		doozerp.Wait()
		doozerpIsDead = true
	}()
	if doozerpIsDead {
		t.Fatal(errors.New("doozerp died prematurely"))
	}
	
	return &Cluster{conn: conn, doozerd: doozerd, doozerp: doozerp, j: j}
}

func (c *Cluster) Close() {
	c.conn.Close()
	c.doozerd.Process.Kill()
	c.doozerp.Wait()
	os.Remove(c.j)
}


func TestNewCluster(t *testing.T) {
	c := NewCluster(t, "")
	defer c.Close()
}
