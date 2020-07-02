package formula

type CreateBuilderManager struct {
	Creator
	Builder
}

func NewCreateBuilder(creator Creator, builder Builder) CreateBuilderManager {
	return CreateBuilderManager{creator, builder}
}
