package scanner

import (
	"fmt"
	dt "lig/datatypes"
	"strconv"
)

var dummy dt.Token = dt.Token{dt.Error, 0}

type Scanner struct {
	Src string 
	cur int
}

func New(source string) *Scanner {
	return &Scanner{source, 0}
}

func (s *Scanner) ScanTokens() ([]dt.Token, error) {
    var res []dt.Token
    var temp dt.Token
    var err error

    for !s.isAtEnd() {
    	temp, err = s.scanToken()
    	if (err != nil) {
    		return res, fmt.Errorf("ScanError: %w", err)
    	}
    	res = append(res, temp)
    }

    // input line terminated without leftover stuffs
    res = append(res, dt.Token{dt.EOF, 0})

    return res, nil
}

type ScanError struct {
	Source string // Bad design? if dealing with source code not repl-like source line.
	Pos int // Character level position
	Msg string // Error message
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("In position %v in source [%s], error occured: %s\n", e.Pos, e.Source, e.Msg)
	// This is excessive. Deal such errors at top level.
}

func (s *Scanner) isAtEnd() bool  {
	return len(s.Src) <= s.cur
}


// When error retval is not used
func (s *Scanner) scanToken() (dt.Token, error) {
	var res dt.Token

	c := s.advance()

	switch c {
		case '+':
			if s.match('+') {
				res = dt.Token{dt.AddAdd, 0}
				break
			}
			res = dt.Token{dt.Add, 0}

		case '-':
			res = dt.Token{dt.Sub, 0}

		case '*':
			res = dt.Token{dt.Mult, 0}

		case '/':
			res = dt.Token{dt.Div, 0}

		case '!':
			if s.match('=') {
				res = dt.Token{dt.BangEqual, 0}
				break
			}
			res = dt.Token {dt.Bang, 0}

		case '>':
			if s.match('=') {
				res = dt.Token{dt.GreaterEqual, 0}
				break
			}
			res = dt.Token {dt.Greater, 0}

		case '<':
			if s.match('=') {
				res = dt.Token{dt.LessEqual, 0}
				break
			}
			res = dt.Token {dt.Less, 0}

		case '=':
			if s.match('=') {
				res = dt.Token {dt.EqualEqual, 0}
				break
			}
			res = dt.Token {dt.Equal, 0}

		case '&':
			if s.match('&') {
				res = dt.Token{dt.And, 0}
				break
			}
			return dummy, &ScanError{s.Src, s.cur, fmt.Sprintf("Character %v cannot be used alone", '&')}

		case '|':
			if s.match('|') {
				res = dt.Token{dt.Or, 0}
				break
			}
			return dummy, &ScanError{s.Src, s.cur, fmt.Sprintf("Character %v cannot be used alone", '|')}

		case '"':
			var strErr error
			res, strErr = s.string()
			if strErr != nil {
				return dummy, strErr
			}
 
		default:
			if isDigit(c) {
				var numErr error
				res, numErr = s.number()
				if numErr != nil {
					return res, numErr
				}
			} else if isAlpha(c) {
				res = s.identifier()
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

func (s *Scanner) number() (dt.Token, error) {
	start := s.cur-1
	for !s.isAtEnd() && isDigit(s.peek()) {
		_ = s.advance()
	}

	num, atoiErr := strconv.Atoi(s.Src[start:s.cur])

	if(atoiErr != nil) {
		return dummy, fmt.Errorf("Failed to parse string to int: %w", atoiErr)
	}

	return dt.Token{dt.Number, num}, nil
}

func (s *Scanner) string() (dt.Token, error) {
	start := s.cur // right after "
	for !s.isAtEnd() && (s.peek() != '"') {
		_ = s.advance()
	}

	stringVal := s.Src[start:s.cur]

	if s.isAtEnd() {
		return dummy, &ScanError{s.Src, s.cur, fmt.Sprintf("Unterminated string: %s", stringVal)}
	}

	return dt.Token{dt.String, stringVal}, nil
}

func (s *Scanner) identifier() dt.Token {
	start := s.cur-1
	for !s.isAtEnd() && isAlphaNumeric(s.peek()) {
		_ = s.advance()
	}

	name := s.Src[start:s.cur]

	return dt.Token{dt.Identifier, name}
}

func (s *Scanner) peek() byte {
	return s.Src[s.cur]
}

func (s *Scanner) advance() byte {
	s.cur += 1
	return s.Src[s.cur-1]
}

// Consumes byte if matches
func (s *Scanner) match(expected byte) bool {
	if s.peek() == expected {
		_ = s.advance()
		return true
	}
	return false
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isAlphaNumeric(c byte) bool {
	return isDigit(c) || isAlpha(c)
}
