package main

import (
	"bufio"
	"fmt"
	"os"
	"lig/scanner"
)
	
func main() {
	if len(os.Args) != 1 {
		fmt.Println("Working on taking files yet... Please use the REPL")
		return
	}

	reader := bufio.NewReader(os.Stdin)


	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')

		tokenSlice, err := scanner.ScanTokens(input[:len(input)-1])

		if err != nil {
			fmt.Println(err)
			
			fmt.Println("Until then, scanned this:")
		}


		for i, value := range tokenSlice {
			fmt.Printf("%v: %+v ", i, value)
		}
		fmt.Println()
	}
}
