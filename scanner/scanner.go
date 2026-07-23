package scanner

import (
	"fmt"
	dt "lig/datatypes"
	"strconv"
)

var errToken dt.Token = dt.Token{Type:dt.Error}
var skipToken dt.Token = dt.Token{Type:dt.Skip}

type Scanner struct {
	Src []byte
	cur int
	curLine int
}

func New(source []byte) *Scanner {
	return &Scanner{source, 0, 1}
}

func (s *Scanner) ScanTokens() ([]dt.Token, []error) {
  var res []dt.Token
  var errSlice[]error = nil

  for !s.isAtEnd() {
  	temp, scanErr := s.scanToken()
  	if scanErr != nil {
  		errSlice = append(errSlice, scanErr)
   		continue
   	}
   	if temp.Type == dt.EOF {
   		break
   	}
   	if temp.Type == dt.Skip {
   		continue
   	}
   	res = append(res, temp)
  }

  // input line terminated without leftover stuffs
  res = append(res, dt.Token{Type:dt.EOF, Line: s.curLine})
	return res, errSlice // err can be nil
}

type ScanError struct {
	CurLine int
	Msg string // Error message
}

func (e *ScanError) Error() string {
	return e.Msg
}

func (s *Scanner) isAtEnd() bool  {
	return len(s.Src) <= s.cur
}


// When error retval is not used
func (s *Scanner) scanToken() (dt.Token, error) {
	var res dt.Token
	var c byte

	s.skipWhitespace()
	if !s.isAtEnd() {
		c = s.advance()
	} else {
		return dt.Token{Type:dt.EOF, Line: s.curLine}, nil
	}

	switch c {
		case '+':
			if s.match('+') {
				res = dt.Token{Type:dt.AddAdd, Line: s.curLine}
				break
			}
			res = dt.Token{Type:dt.Add, Line: s.curLine}

		case '-':
			res = dt.Token{Type:dt.Sub, Line: s.curLine}

		case '*':
			res = dt.Token{Type:dt.Mult, Line: s.curLine}

		case '/':
			if s.match('/') {
				s.skipTilNewline()
				res = skipToken
			}
			res = dt.Token{Type:dt.Div, Line: s.curLine}

		case '!':
			if s.match('=') {
				res = dt.Token{Type:dt.BangEqual, Line: s.curLine}
				break
			}
			res = dt.Token{Type:dt.Bang, Line: s.curLine}

		case '>':
			if s.match('=') {
				res = dt.Token{Type:dt.GreaterEqual, Line: s.curLine}
				break
			}
			res = dt.Token{Type:dt.Greater, Line: s.curLine}

		case '<':
			if s.match('=') {
				res = dt.Token{Type:dt.LessEqual, Line: s.curLine}
				break
			}
			res = dt.Token{Type:dt.Less, Line: s.curLine}

		case '=':
			if s.match('=') {
				res = dt.Token{Type:dt.EqualEqual, Line: s.curLine}
				break
			}
			res = dt.Token{Type:dt.Equal, Line: s.curLine}

		case '&':
			if s.match('&') {
				res = dt.Token{Type:dt.And, Line: s.curLine}
				break
			}
			return errToken, &ScanError{s.curLine, fmt.Sprintf("Character %v cannot be used alone", '&')}

		case '|':
			if s.match('|') {
				res = dt.Token{Type:dt.Or, Value:0, Line: s.curLine}
				break
			}
			return errToken, &ScanError{s.curLine, fmt.Sprintf("Character %v cannot be used alone", '|')}

		case '"':
			var strErr error
			res, strErr = s.string()
			if strErr != nil {
				return errToken, strErr
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
				return skipToken, &ScanError{s.curLine, fmt.Sprintf("Unexpected character: %q", c)}
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
	if s.peek() == '\n' { s.curLine += 1 }
	return (s.peek() == ' ' || s.peek() == '\n' || s.peek() == '\t')
}

func (s *Scanner) skipTilNewline() {
	for !s.isAtEnd() && s.peek() != '\n' { s.cur += 1 }
	if s.peek() == '\n' {
		s.cur += 1 // Now looking at next line!
		s.curLine += 1
	}
}

func (s *Scanner) number() (dt.Token, error) {
	start := s.cur-1
	for !s.isAtEnd() && isDigit(s.peek()) {
		_ = s.advance()
	}

	num, atoiErr := strconv.Atoi(string(s.Src[start:s.cur]))

	if(atoiErr != nil) {
		return errToken, &ScanError{s.curLine, fmt.Sprintf("Failed to parse string to int: %s", atoiErr.Error())} // Wrap it up for line info
	}

	return dt.Token{Type:dt.Number, Value:num, Line: s.curLine}, nil
}

func (s *Scanner) string() (dt.Token, error) {
	start := s.cur // right after "

	for !s.isAtEnd() && (s.peek() != '"') {
		if s.peek() == '\n' { s.curLine += 1 }
		_ = s.advance()
	}

	stringVal := string(s.Src[start:s.cur])

	if s.isAtEnd() {
		return errToken, &ScanError{s.curLine, fmt.Sprintf("Unterminated string: %s", stringVal)}
	}
	_ = s.advance()

	return dt.Token{Type:dt.String, Value:stringVal, Line: s.curLine}, nil
}

func (s *Scanner) identifier() dt.Token {
	start := s.cur-1
	for !s.isAtEnd() && isAlphaNumeric(s.peek()) {
		_ = s.advance()
	}

	name := string(s.Src[start:s.cur])

	if tokenType, exists := dt.Keywords[name]; exists {
		return dt.Token{Type:tokenType, Name:name, Line: s.curLine}
	}
	return dt.Token{Type:dt.Identifier, Name:name, Line: s.curLine}
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
