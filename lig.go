package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"lig/interpreter"
	"lig/parser"
	"lig/scanner"
	"os"
)

var global *interpreter.Environment = interpreter.NewEnvironment(nil)

func main() {
	// Flag processing
	modePtr := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()
	debug := *modePtr
	if flag.NArg() > 1 {
		fmt.Println("Usage: lox <file_path>")
		os.Exit(1)
	} else if flag.NArg() == 1 {
		path := flag.Args()[0]
		filePtr, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			filePtr.Close()
			os.Exit(1)
		}
		defer filePtr.Close()

		fileReader := bufio.NewReader(filePtr)

		src, readErr := io.ReadAll(fileReader) // Can be made better with unsafe(no modification ensured)
		if readErr != nil {
			fmt.Println(readErr)
			filePtr.Close()
			os.Exit(1)
		}

		if debug {
			runWithDebug(src)
		} else {
			run(src)
		}
	} else {

		reader := bufio.NewReader(os.Stdin)

		for {
			// Take input
			fmt.Print("> ")
			input, readErr := reader.ReadString('\n')
			if readErr != nil {
				fmt.Println("ReadError: %w", readErr)
				os.Exit(1)
			}

			if debug {
				runWithDebug([]byte(input))
			} else {
				run([]byte(input))
			}
		}
	}
}

func runWithDebug(src []byte) {
	// Scan it to token list
	srcScanner := scanner.New(src)
	tokenSlice, scanErr := srcScanner.ScanTokens()
	var scanError *scanner.ScanError

	if scanErr != nil {
		for _, err := range scanErr {
			if errors.As(err, &scanError) {
				fmt.Println("[ line", scanError.CurLine, "]", "ScanError:", scanError)
			} else {
				fmt.Println(err)
			}
		}

		fmt.Println("Until then, scanned this:")

		for i, value := range tokenSlice {
			fmt.Printf("%v: %+v ", i, value)
		}
		fmt.Println()
		return //panic-- More like, the token list is wrong, so don't parse it!
	}

	fmt.Printf("Scan Result: ")
	for _, value := range tokenSlice {
		fmt.Printf("%+v ", value)
	}
	fmt.Println()

	// Parse to AST
	parser := parser.New(tokenSlice)
	stmtSlice, parseErr := parser.Parse()

	if parseErr != nil {
		fmt.Println("[ line", parseErr.Token.Line, "]", "ParseError:", parseErr)

		fmt.Println("Until then, parsed this:")

		fmt.Printf("%+v\n", stmtSlice)
		return // panic
	}

	fmt.Printf("Parse Result: ")
	fmt.Printf("%+v\n", stmtSlice)

	// Interpret AST to value
	interp := interpreter.NewInterpreter(stmtSlice, global)
	interpErr := interp.Interpret()
	if interpErr != nil {
		fmt.Println("[ line", interpErr.ErrToken.Line, "]", "RuntimeError:", interpErr)
		return // Panic
	}
}

func run(src []byte) {
	// Scan it to token list
	srcScanner := scanner.New(src)
	tokenSlice, scanErr := srcScanner.ScanTokens()
	var scanError *scanner.ScanError

	if scanErr != nil {
		for _, err := range scanErr {
			if errors.As(err, &scanError) {
				fmt.Println("[ line", scanError.CurLine, "]", "ScanError:", scanError)
			} else {
				fmt.Println(err)
			}
		}
		return //panic-- More like, the token list is wrong, so don't parse it!
	}

	// Parse to AST
	parser := parser.New(tokenSlice)
	stmtSlice, parseErr := parser.Parse()

	if parseErr != nil {
		fmt.Println("[ line", parseErr.Token.Line, "]", "ParseError:", parseErr)
		return // panic
	}

	// Interpret AST to value
	interp := interpreter.NewInterpreter(stmtSlice, global)
	interpErr := interp.Interpret()
	if interpErr != nil {
		fmt.Println("[ line", interpErr.ErrToken.Line, "]", "RuntimeError:", interpErr)
		return // Panic
	}
}
