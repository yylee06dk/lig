package datatypes

var Keywords map[string]Tokentype = map[string]Tokentype{
	"true":  True,
	"false": False,
	"print": Print,
	"var":   Var,
	"if":    If,
	"else":  Else,
	"while": While,
	"for":   For,
}
