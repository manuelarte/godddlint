package simple

//godddlint:valueObject
//godddlint:disable:VOX001
type ImmutableDisableStructLevel struct {
	x, y int
}

func (p ImmutableDisableStructLevel) X() int {
	return p.x
}

//godddlint:valueObject
type ImmutableDisableFieldLevel struct {
	//godddlint:disable:VOX001
	x, y int
}

func NewImmutableDisableFieldLevel() (ImmutableDisableFieldLevel, error) {
	return ImmutableDisableFieldLevel{}, nil
}
