package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for true {
		// Prompt
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		s, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			_ = fmt.Errorf(err.Error())
		}
		s = strings.TrimSpace(s)

		// All commands are unknown to us
		if true {
			fmt.Printf("%s: command not found\n", s)
		}
	}
}
