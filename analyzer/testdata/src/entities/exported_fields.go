package simple

import "fmt"

//godddlint:entity
type UserExported struct {
	Id      UserID // want `E005: Entity's field is exported`
	Name    Name // want `E005: Entity's field is exported`
	Surname Surname // want `E005: Entity's field is exported`
}

func (u UserExported) UserExported() string { // want `E001: Entity's method not using pointer receiver`
	return fmt.Sprintf("%s %s", u.Name, u.Surname)
}
