package deku

var symbolValuePlaceholder = "??"
var symbolEdgeTraversal = "->"

type DQLizer interface {
	ToDQL() (query string, args []interface{}, err error)
}
