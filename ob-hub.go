package obcore

type ObHub struct {
	parent *ObHub
}

func (me *ObHub) Parent() (parent *ObHub) {
	return nil
}
