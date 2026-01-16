package simple

//godddlint:valueObject
//godddlint:disable:VO001
type NonPointerReceiverDisabledStructLevel struct {
	x, y int
}

func NewNonPointerReceiverDisabledStructLevel() NonPointerReceiverDisabledStructLevel {
	return NonPointerReceiverDisabledStructLevel{}
}

func (p *NonPointerReceiverDisabledStructLevel) X() int {
	return p.x
}

//godddlint:valueObject
type NonPointerReceiverDisabledMethodLevel struct {
	x, y int
}

func NewNonPointerReceiverDisabledMethodLevel() NonPointerReceiverDisabledMethodLevel {
	return NonPointerReceiverDisabledMethodLevel{}
}

//godddlint:disable:VO001
func (p *NonPointerReceiverDisabledMethodLevel) X() int {
	return p.x
}

func (p *NonPointerReceiverDisabledMethodLevel) Y() int { // want `VO001: Value Object's method using a pointer receiver`
	return p.y
}
