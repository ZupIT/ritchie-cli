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

type Setter interface {
	Set(ctx string) (ContextHolder, error)
}

type CtxFinder interface {
	Find() (ContextHolder, error)
}

type Remover interface {
	Remove(ctx string) (ContextHolder, error)
}

type FindRemover interface {
	CtxFinder
	Remover
}

type FindSetter interface {
	CtxFinder
	Setter
}
