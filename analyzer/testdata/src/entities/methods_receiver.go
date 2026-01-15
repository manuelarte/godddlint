package simple

import "fmt"

type (
	Name    string
	Surname string
)

//godddlint:entity
type User struct {
	Id      int // want `E003: Prefer custom domain types to primitives`
	Name    Name
	Surname Surname
}

func (u User) FullName() string { // want `E001: Entity's method not using pointer receiver`
	return fmt.Sprintf("%s %s", u.Name, u.Surname)
}
