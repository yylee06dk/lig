package main

import (
	"bufio"
	"fmt"
	"os"
	"lig/scanner"
	"lig/parser"
	"lig/interpreter"
)
	
func main() {
	if len(os.Args) != 1 {
		fmt.Println("Working on taking files yet... Please use the REPL")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	var srcScanner *scanner.Scanner


	for {
		// Take input
		fmt.Print("> ")
		input, readErr := reader.ReadString('\n')
		if readErr != nil {
			fmt.Println("ReadError: %w", readErr)
			os.Exit(1)
		}

		// Scan it to token list
		srcScanner = scanner.New(input)
		tokenSlice, scanErr := srcScanner.ScanTokens()

		if scanErr != nil {
			fmt.Println(scanErr)
			
			fmt.Println("Until then, scanned this:")

			for i, value := range tokenSlice {
				fmt.Printf("%v: %+v ", i, value)
			}
			fmt.Println()
			continue
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
			continue
		}

		fmt.Printf("Parse Result: ")
		fmt.Printf("%+v\n", expr)

		// Interpret AST to value
		resVal, interpErr := interpreter.Interpret(expr)
		if interpErr != nil {
			fmt.Println(interpErr)
			continue
		}
		fmt.Printf("Interpret Result: %v\n", resVal)
	}
}
