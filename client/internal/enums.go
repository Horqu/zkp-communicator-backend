package internal

type AppView int

const (
	ViewMain AppView = iota
	ViewRegister
	ViewLogin
	ViewResolver
	ViewLogged
	ViewLoading
	ViewError
)
