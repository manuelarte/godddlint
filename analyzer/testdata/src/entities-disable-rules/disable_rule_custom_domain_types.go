package simple

//godddlint:entity
//godddlint:disable:E003
type CustomDomainTypesIgnoreStructLevel struct {
	id      int
	name    string
	surname string
}

//godddlint:entity
type CustomDomainTypesIgnoreFieldLevel struct {
	//godddlint:disable:E003
	id      int
	name    string // want `E003: Prefer custom domain types to primitives`
	surname string // want `E003: Prefer custom domain types to primitives`
}
