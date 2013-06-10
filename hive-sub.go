package obcore

import (
	"io/ioutil"

	uio "github.com/metaleap/go-util/io"
)

type HiveSub struct {
	name []string
	root *HiveRoot

	Paths struct {
		Client, ClientPub string
	}
}

func (me *HiveSub) init(root *HiveRoot, name string) {
	me.root, me.name = root, []string{name}
	p := &me.Paths
	p.Client, p.ClientPub = me.Path("client"), me.Path("client", "pub")
}

func (me *HiveSub) FileExists(relPath ...string) bool {
	return uio.FileExists(me.Path(relPath...))
}

func (me *HiveSub) FilePath(relPath ...string) (filePath string) {
	if me.FileExists(relPath...) {
		filePath = me.Path(relPath...)
	}
	return
}

func (me *HiveSub) Path(relPath ...string) string {
	return me.root.Path(append(me.name, relPath...)...)
}

func (me *HiveSub) ReadFile(relPath ...string) (data []byte, err error) {
	data, err = ioutil.ReadFile(me.Path(relPath...))
	return
}
