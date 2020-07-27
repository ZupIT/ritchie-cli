package formula

type CreateBuilderManager struct {
	Creator
	LocalBuilder
}

func NewCreateBuilder(creator Creator, builder LocalBuilder) CreateBuilder {
	return CreateBuilderManager{creator, builder}
}
