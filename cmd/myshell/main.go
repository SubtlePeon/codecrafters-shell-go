package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	for true {
		// Prompt
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		s, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		s = strings.TrimSuffix(s, "\n")

		cmd := strings.Split(s, " ")

		for i := len(cmd) - 1; i >= 0; i-- {
			if cmd[i] == "" {
				cmd = append(cmd[:i], cmd[i+1:]...)
			}
		}

		// All commands are unknown to us
		if len(cmd) == 0 {
			// Just pass
		} else if cmd[0] == "exit" {
			handle_exit(cmd[1:])
		} else if cmd[0] == "echo" {
			s = strings.TrimPrefix(s, "echo")
			s = strings.TrimPrefix(s, " ")
			fmt.Println(s)
		} else {
			fmt.Printf("%s: command not found\n", s)
		}
	}
}

// Handles the `exit` command. The arguments after "exit" should be split and
// passed in.
func handle_exit(cmd []string) {
	if len(cmd) == 0 {
		os.Exit(0)
	} else if len(cmd) == 1 {
		num, err := strconv.ParseInt(cmd[0], 10, 32)
		if err != nil {
			fmt.Printf("exit: could not parse exit code\n%v\n", err)
			return
		}
		os.Exit(int(num))
	} else {
		fmt.Println("exit: too many arguments")
	}
}
