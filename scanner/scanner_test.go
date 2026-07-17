package scanner

import (
  "github.com/google/go-cmp/cmp"
  "testing"
  dt "lig/datatypes"
)

func TestScanTokens(t *testing.T) {
  tests := []struct{
    name string
    source string
    expected []dt.Token
    wantErr bool
    errMsg string
  }{
    {
      name: "Stage1: Basic Arithmetic 1",
      source: "1+2",
      expected: []dt.Token{
        dt.Token{dt.Number, 1},
        dt.Token{dt.Add, 0},
        dt.Token{dt.Number, 2},
        dt.Token{dt.EOF, 0},
      },
      wantErr: false,
    },
    {
      name: "Stage1: Basic Arithmetic 2",
      source:"12*13+40/19-123",
      expected: []dt.Token{
        dt.Token{dt.Number, 12},
        dt.Token{dt.Mult, 0},
        dt.Token{dt.Number, 13},
        dt.Token{dt.Add, 0},
        dt.Token{dt.Number, 40},
        dt.Token{dt.Div, 0},
        dt.Token{dt.Number, 19},
        dt.Token{dt.Sub, 0},
        dt.Token{dt.Number, 123},
        dt.Token{dt.EOF, 0},
      },
      wantErr: false,
    },
  }

  

  for _, testCase := range tests {
    t.Run(testCase.name, func(t *testing.T) {
      scanner := New(testCase.source)
      res, _ := scanner.ScanTokens() // Currently no testcase giving error is written

      if !testCase.wantErr{
        if diff := cmp.Diff(res, testCase.expected); diff != "" {
          t.Errorf("Token list mismatch: (-want +got):\n%s", diff)
        }
      }
    })
  }
}

func FuzzScanTokens(f *testing.F) {
  f.Add("1+2")
  f.Add("123//12312/4++*")

  f.Fuzz(func(t *testing.T, source string) {
    scanner := New(source)
    _, err := scanner.ScanTokens()
    if err != nil { // Error occured?
      t.Errorf("%s", err.Error())
    }
  })
}
