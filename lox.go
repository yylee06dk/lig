package main

import (
	"bufio"
	"fmt"
	"os"
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

		fmt.Printf("Input: %s\n", input)
	}
}
