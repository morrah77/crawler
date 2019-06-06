package common

type IParser interface {
	Init([]byte)
	Parse() [][]byte
	ParseNext() ([]byte, bool)
}
