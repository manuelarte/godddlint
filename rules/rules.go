package rules

type (
	RuleMetadata struct {
		Code, Name, URL string
	}
)

//nolint:gochecknoglobals // enum of rules
var (
	PointerReceivers = RuleMetadata{
		Code: "E001",
		Name: "Pointer Receivers",
	}
	CustomTypesOverPrimitives = RuleMetadata{
		Code: "E003",
		Name: "Custom Types Over Primitives",
	}
	CustomDomainErrors = RuleMetadata{
		Code: "E004",
		Name: "Custom Domain Errors",
	}
	UnexportedFields = RuleMetadata{
		Code: "E005",
		Name: "Unexported Fields",
	}

	NonPointerReceivers = RuleMetadata{
		Code: "VO001",
		Name: "Non Pointer Receivers",
	}
	Immutable = RuleMetadata{
		Code: "VOX001",
		Name: "Immutable",
	}
	DefensiveCopy = RuleMetadata{
		Code: "VOX002",
		Name: "Maps/Slices Not Defensive Copied",
	}
)
