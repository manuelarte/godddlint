package valueobject

var _ Rule = new(NonPointerReceivers)

type (
	RuleMetadata struct {
		Name, Description, URL string
	}

	Rule interface {
		Metadata() RuleMetadata
		Apply(definition Definition)
		x()
	}

	NonPointerReceivers struct{}
)

func (r NonPointerReceivers) Metadata() RuleMetadata {
	return RuleMetadata{
		Name:        "VO001",
		Description: "Non Pointer Receivers",
	}
}

func (r NonPointerReceivers) Apply(d Definition) {
	// TODO implement me
}

func (r NonPointerReceivers) x() {}
