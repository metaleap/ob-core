package obcore

import (
	"io/ioutil"

	"github.com/go-utils/ufs"
)

//	Represents either the dist or the cust directory in a hive directory
type HiveSub struct {
	name []string
	root *HiveRoot

	//	Paths to some well-known HiveSub directories
	Paths struct {
		//	{hive}/{sub}/client
		Client string

		//	{hive}/{sub}/client/pub
		ClientPub string

		//	{hive}/{sub}/pkg
		Pkg string
	}
}

func (me *HiveSub) init(root *HiveRoot, name string) {
	me.root, me.name = root, []string{name}
	p := &me.Paths
	p.Client, p.ClientPub, p.Pkg = me.Path("client"), me.Path("client", "pub"), me.Path("pkg")
}

//	Returns whether the specified {hive}/{sub}-relative directory exists
func (me *HiveSub) DirExists(relPath ...string) bool {
	return ufs.DirExists(me.Path(relPath...))
}

//	Returns a {hive}/{sub}-joined representation of the specified
//	{hive}/{sub}-relative path, if it represents an existing directory
func (me *HiveSub) DirPath(relPath ...string) (dirPath string) {
	if me.DirExists(relPath...) {
		dirPath = me.Path(relPath...)
	}
	return
}

//	Returns whether the specified {hive}/{sub}-relative file exists
func (me *HiveSub) FileExists(relPath ...string) bool {
	return ufs.FileExists(me.Path(relPath...))
}

//	Returns a {hive}/{sub}-joined representation of the specified
//	{hive}/{sub}-relative path, if it represents an existing file
func (me *HiveSub) FilePath(relPath ...string) (filePath string) {
	if me.FileExists(relPath...) {
		filePath = me.Path(relPath...)
	}
	return
}

//	Returns a {hive}/{sub}-joined representation of the specified
//	{hive}/{sub}-relative path
func (me *HiveSub) Path(relPath ...string) string {
	return me.root.Path(append(me.name, relPath...)...)
}

//	Reads the file at the specified {hive}/{sub}-relative path
func (me *HiveSub) ReadFile(relPath ...string) (data []byte, err error) {
	data, err = ioutil.ReadFile(me.Path(relPath...))
	return
}
