package simple

//godddlint:valueObject
type PointWithoutConstructor struct { // want `VOX001: Constructor for Value Object not found`
	x, y int
}

func (p PointWithoutConstructor) X() int {
	return p.x
}

//godddlint:valueObject
type PointWithConstructor struct {
	x, y int
}

func NewPointWithConstructor() (PointWithConstructor, error) {
	return PointWithConstructor{}, nil
}

func MustPointWithConstructor() PointWithConstructor {
	return PointWithConstructor{}
}
