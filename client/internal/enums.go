package internal

type AppView int

const (
	ViewMain AppView = iota
	ViewResolver
	ViewLoading
	ViewError
)
