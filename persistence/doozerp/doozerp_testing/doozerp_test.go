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
	conn          *doozer.Conn
	doozerd       *exec.Cmd
	doozerp       *exec.Cmd
	j             string
	doozerdIsDead bool
	doozerpIsDead bool
}

func NewCluster(t *testing.T, doozerpArgs ...string) (c *Cluster) {
	c = new(Cluster)
	c.doozerd = exec.Command("doozerd", "-l=127.0.0.1:19999", "-w=false")
	err := c.doozerd.Start()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		c.doozerd.Wait()
		c.doozerdIsDead = true
	}()
	time.Sleep(100 * time.Millisecond)
	if c.doozerdIsDead {
		t.Fatal(errors.New("doozerd died prematurely"))
	}

	c.conn, err = doozer.DialUri("doozer:?ca=127.0.0.1:19999", "")
	if err != nil {
		c.doozerd.Process.Kill()
		t.Fatal(err)
	}

	f, err := ioutil.TempFile("", "j")
	if err != nil {
		c.conn.Close()
		c.doozerd.Process.Kill()
		t.Fatal(err)
	}
	c.j = f.Name()
	f.Close()
	args := append(doozerpArgs, "-a=doozer:?ca=127.0.0.1:19999")
	args = append(args, "-j="+c.j)
	c.doozerp = exec.Command("doozerp", args...)
	err = c.doozerp.Start()
	if err != nil {
		c.conn.Close()
		c.doozerd.Process.Kill()
		t.Fatal(err)
	}
	go func() {
		c.doozerp.Wait()
		c.doozerpIsDead = true
	}()
	time.Sleep(100 * time.Millisecond)
	if c.doozerpIsDead {
		t.Fatal(errors.New("doozerp died prematurely"))
	}

	return
}

func (c *Cluster) Close() {
	c.conn.Close()
	c.doozerd.Process.Kill()
	if !c.doozerpIsDead {
		c.doozerp.Wait()
	}
	os.Remove(c.j)
}

func TestNewCluster(t *testing.T) {
	c := NewCluster(t)
	defer c.Close()
}
