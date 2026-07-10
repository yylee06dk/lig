package main

import (
	"bufio"
	"fmt"
	"os"
	"lig/scanner"
	"lig/parser"
)
	
func main() {
	if len(os.Args) != 1 {
		fmt.Println("Working on taking files yet... Please use the REPL")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	var srcScanner *scanner.Scanner


	for {
		fmt.Print("> ")
		input, readErr := reader.ReadString('\n')
		if readErr != nil {
			fmt.Println("ReadError: %w", readErr)
			os.Exit(1)
		}

		srcScanner = scanner.New(input)
		tokenSlice, scanErr := srcScanner.ScanTokens()

		if scanErr != nil {
			fmt.Println(scanErr)
			
			fmt.Println("Until then, scanned this:")

			for i, value := range tokenSlice {
				fmt.Printf("%v: %+v ", i, value)
			}
			fmt.Println()
			os.Exit(1)
		}

		for _, value := range tokenSlice {
			fmt.Printf("%+v ", value)
		}
		fmt.Println()

		expr := parser.Parse(tokenSlice)

		fmt.Printf("%+v]\n", expr)
	}
}
