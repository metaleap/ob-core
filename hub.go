package obcore

func RootHub() *Hub {
	return GetHub("/")
}

func GetHub(path string) (hub *Hub) {
	return
}

type Hub struct {
	parent *Hub
}

func newHub() (me *Hub) {
	me = &Hub{}
	return
}

func (me *Hub) Parent() (parent *Hub) {
	return me.parent
}
