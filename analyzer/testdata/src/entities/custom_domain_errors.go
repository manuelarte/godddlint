package simple

import "errors"

type Repository interface {
	SaveUser(u *UserDomainError) error
}

type UserID int
type Address string

//godddlint:entity
type UserDomainError struct {
	id        UserID
	addresses []Address
}

func (u *UserDomainError) AddAddress(na Address) error {
	if len(u.addresses) >= 2 {
		return errors.New("maximum number of addresses reached") // want `E004: Use Custom Domain Errors`
	}
	u.addresses = append(u.addresses, na)

	return nil
}

func (u *UserDomainError) SaveUser(repo Repository) error {
	return repo.SaveUser(u) // want `E004: Use Custom Domain Errors`
}
