package simple

import "errors"

type UserID int
type Address string

//godddlint:entity
//godddlint:disable:E004
type CustomDomainErrorsIgnoreStructLevel struct {
	id        UserID
	addresses []Address
}

func (u *CustomDomainErrorsIgnoreStructLevel) AddAddress(na Address) error {
	if len(u.addresses) >= 2 {
		return errors.New("maximum number of addresses reached")
	}
	u.addresses = append(u.addresses, na)

	return nil
}

//godddlint:entity
type CustomDomainErrorsIgnoreMethodLevel struct {
	id        UserID
	addresses []Address
}

//godddlint:disable:E004
func (u *CustomDomainErrorsIgnoreMethodLevel) AddAddress(na Address) error {
	if len(u.addresses) >= 2 {
		return errors.New("maximum number of addresses reached")
	}
	u.addresses = append(u.addresses, na)

	return nil
}
