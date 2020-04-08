package rcontext

const (
	ContextPath = "%s/contexts"
	CurrentCtx  = "Current -> "
	DefaultCtx  = "default"
)

type ContextHolder struct {
	Current string   `json:"current_context"`
	All     []string `json:"contexts"`
}

//go:generate $GOPATH/bin/moq -out setter_moq.go . Setter
type Setter interface {
	Set(ctx string) (ContextHolder, error)
}

//go:generate $GOPATH/bin/moq -out finder_moq.go . Finder
type Finder interface {
	Find() (ContextHolder, error)
}

//go:generate $GOPATH/bin/moq -out remover_moq.go . Remover
type Remover interface {
	Remove(ctx string) (ContextHolder, error)
}

//go:generate $GOPATH/bin/moq -out find_remover_moq.go . FindRemover
type FindRemover interface {
	Finder
	Remover
}

//go:generate $GOPATH/bin/moq -out find_setter_moq.go . FindSetter
type FindSetter interface {
	Finder
	Setter
}
