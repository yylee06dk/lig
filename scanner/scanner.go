package scanner

import (
	"fmt"
	"lig/datatypes"
	"strconv"
)

func ScanTokens(source string) ([]datatypes.Token, error) {
    var res []datatypes.Token
    var temp datatypes.Token
    var start int
    var err error
    var current int

    for !isAtEnd(source, current) {
    	temp, start, err = scanToken(source, current)
    	if (err != nil) {
    		break
    	}
    	res = append(res, temp)
    	current = start
    }

    return res, err
}

type ScanError struct {
	Source string
	Pos int
	Msg string
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("In position %v in source [%s], error occured: %s\n", e.Pos, e.Source, e.Msg)
}

func isAtEnd(source string, current int) bool  {
	return len(source) <= current
}


func scanToken(source string, start int) (datatypes.Token, int, error) {
	var res datatypes.Token
	current := start

	c, current := advance(source, current)

	switch c {
		case '+':
			res = datatypes.Token{datatypes.Add, 0}
		case '-':
			res = datatypes.Token{datatypes.Sub, 0}
		case '*':
			res = datatypes.Token{datatypes.Mult, 0}
		case '/':
			res = datatypes.Token{datatypes.Div, 0}
		default:
			if isDigit(c) {
				res, current = number(source, current)
			} else {
				return res, -1, &ScanError{source, current, "Unexpected character"}
			}
	}

	return res, skipWhitespace(source, current), nil
}

func skipWhitespace(source string, start int) int {
	current := start
	for !isAtEnd(source, current) &&  isWhitespace(source, current) {
		current += 1
	}
	return current
}

func isWhitespace(source string, current int) bool {
	return (source[current] == ' ' || source[current] == '\n' || source[current] == '\t')
}

func number(source string, start int) (datatypes.Token, int) {
	current := start
	
	for !isAtEnd(source, current) {
		if isDigit(peek(source, current)) {
			_, current = advance(source, current)
		} else { break }
	}

	num, _ := strconv.Atoi(source[start-1:current])

	return datatypes.Token{datatypes.Number, num}, current
}

func peek(source string, current int) (byte) {
	return source[current]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func advance(source string, current int) (byte, int) {
	return source[current], current+1
}
