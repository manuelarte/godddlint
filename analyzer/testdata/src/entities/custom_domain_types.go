package simple

import "fmt"

//godddlint:entity
type User2 struct {
	id      int    // want `E003: Prefer custom domain types to primitives`
	name    string // want `E003: Prefer custom domain types to primitives`
	surname string // want `E003: Prefer custom domain types to primitives`
}

func (u User2) FullName() string { // want `E001: Entity's method not using pointer receiver`
	return fmt.Sprintf("%s %s", u.name, u.surname)
}
