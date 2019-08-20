package app

type Command interface {
	Execute() error
}