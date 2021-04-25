package deku

type DQLizer interface {
	ToDQL() (query string, args []interface{}, err error)
}
