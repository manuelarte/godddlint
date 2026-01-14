package simple

//godddlint:valueObject
type Point struct {
	x, y int
}

func (p *Point) X() int {
	return p.x
}

func (p *Point) Y() int {
	return p.y
}

//godddlint:valueObject
type Username string

func (u *Username) String() string {
	return string(*u)
}

type NormalStruct struct {
	Name string
}

func (n *NormalStruct) ChangeName(newName string) {
	n.Name = newName
}

type (
	//godddlint:valueObject
	Point2 struct {
		x, y int
	}

	//godddlint:valueObject
	Username2 string

	NormalStruct2 struct {
		Name string
	}
)

func (p *Point2) X() int {
	return p.x
}

func (p *Point2) Y() int {
	return p.y
}

func (u *Username2) String() string {
	return string(*u)
}

func (n *NormalStruct2) ChangeName(newName string) {
	n.Name = newName
}
