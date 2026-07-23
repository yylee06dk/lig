package main

import (
	"errors"
	"io"
	"bufio"
	"fmt"
	"os"
	"lig/scanner"
	"lig/parser"
	"lig/interpreter"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox <file_path>")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		path := os.Args[1]
		filePtr, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fileReader := bufio.NewReader(filePtr)

		src, readErr := io.ReadAll(fileReader) // Can be made better with unsafe(no modification ensured)
		if readErr != nil {
			fmt.Println(readErr)
			os.Exit(1)
		}

		runWithDebug(src)
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

			runWithDebug([]byte(input))
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
				fmt.Println("[ line", scanError.CurLine, "]", scanError)
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
	expr, parseErr := parser.Parse()
	
	if parseErr != nil {
		fmt.Println(parseErr)
		
		fmt.Println("Until then, parsed this:")
	
		fmt.Printf("%+v\n", expr)
		return // panic
	}
	
	fmt.Printf("Parse Result: ")
	fmt.Printf("%+v\n", expr)
	
	// Interpret AST to value
	resVal, interpErr := interpreter.Interpret(expr)
	if interpErr != nil {
		fmt.Println(interpErr)
		return // Panic
	}
	if resVal == nil { fmt.Println() }
	fmt.Printf("Interpret Result: %v\n", resVal)
}
