package simple

//godddlint:entity
//godddlint:disable:E005
type DisableExportedFieldsStructLevel struct {
	Id      UserID
	Name    Name
	Surname Surname
}

//godddlint:entity
type DisableExportedFieldsFieldLevel struct {
	Id UserID // want `E005: Entity's field is exported`
	//godddlint:disable:E005
	Name    Name
	Surname Surname // want `E005: Entity's field is exported`
}
