package datatypes

type Tokentype int

const (
	Number Tokentype = iota
	Add
	Sub
	Mult
	Div
	Group
)
