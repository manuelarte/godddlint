package simple

import "fmt"

type (
	Name    string
	Surname string
)

//godddlint:entity
//godddlint:disable:E001
type MethodReceiverDisableStructLevel struct {
	id      int // want `E003: Prefer custom domain types to primitives`
	name    Name
	surname Surname
}

func (u MethodReceiverDisableStructLevel) FullName() string {
	return fmt.Sprintf("%s %s", u.name, u.surname)
}

//godddlint:entity
type MethodReceiverDisableMethodLevel struct {
	id      int // want `E003: Prefer custom domain types to primitives`
	name    Name
	surname Surname
}

//godddlint:disable:E001
func (u MethodReceiverDisableMethodLevel) FullName() string {
	return fmt.Sprintf("%s %s", u.name, u.surname)
}
