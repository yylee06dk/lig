package scanner

import (
	"fmt"
	"lig/datatypes"
	"strconv"
)

type Scanner struct {
	Src string
	cur int
}

func New(source string) *Scanner {
	return &Scanner{source, 0}
}

func (s *Scanner) ScanTokens() ([]datatypes.Token, error) {
    var res []datatypes.Token
    var temp datatypes.Token
    var err error

    for !s.isAtEnd() {
    	temp, err = s.scanToken()
    	if (err != nil) {
    		break
    	}
    	res = append(res, temp)
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

func (s *Scanner) isAtEnd() bool  {
	return len(s.Src) <= s.cur
}


func (s *Scanner) scanToken() (datatypes.Token, error) {
	var res datatypes.Token
	var numErr error

	c := s.advance()

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
				res, numErr = s.number()
				if numErr != nil {
					return res, fmt.Errorf("ScanError: %w", numErr)
				}
			} else {
				return res, &ScanError{s.Src, s.cur, "Unexpected character"}
			}
	}

	s.skipWhitespace()
	return res, nil
}

func (s *Scanner) skipWhitespace() {
	for !s.isAtEnd() &&  s.isWhitespace() {
		s.cur += 1
	}
}

func (s *Scanner) isWhitespace() bool {
	return (s.Src[s.cur] == ' ' || s.Src[s.cur] == '\n' || s.Src[s.cur] == '\t')
}

func (s *Scanner) number() (datatypes.Token, error) {
	start := s.cur-1
	for !s.isAtEnd() && isDigit(s.peek()) {
		s.cur += 1
	}

	num, atoiErr := strconv.Atoi(s.Src[start:s.cur])

	if(atoiErr != nil) {
		return datatypes.Token{datatypes.Number, 0}, fmt.Errorf("Failed to parse string to int: %w", atoiErr)
	}

	return datatypes.Token{datatypes.Number, num}, nil
}

func (s *Scanner) peek() byte {
	return s.Src[s.cur]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) advance() byte {
	s.cur += 1
	return s.Src[s.cur-1]
}
