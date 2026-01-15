package simple

import "fmt"

type (
	Name    string
	Surname string
)

//godddlint:entity
type User struct {
	id      int // want `E003: Prefer custom domain types to primitives`
	name    Name
	surname Surname
}

func (u User) FullName() string { // want `E001: Entity's method not using pointer receiver`
	return fmt.Sprintf("%s %s", u.name, u.surname)
}
