package simple

import "fmt"

//godddlint:entity
type User2 struct {
	Id      int    // want `E003: Prefer custom types to primitives`
	Name    string // want `E003: Prefer custom types to primitives`
	Surname string // want `E003: Prefer custom types to primitives`
}

func (u User2) FullName() string { // want `E001: Entity's method not using pointer receiver`
	return fmt.Sprintf("%s %s", u.Name, u.Surname)
}
