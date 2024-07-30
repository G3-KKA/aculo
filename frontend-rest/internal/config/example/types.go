package EXAMPLE_DO_NOT_COPY

// sole purpose is to make example.go errorless
// do not copy this file or example.go
// config package containt everything inside
type overrideContainer struct {
	name  string
	value any
}

type flagSetter func()
type elseSetter func() error
