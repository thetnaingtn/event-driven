package command

import "github.com/ThreeDotsLabs/watermill/components/cqrs"

var (
	marshaler = cqrs.JSONMarshaler{
		GenerateName: cqrs.StructName,
	}
)
