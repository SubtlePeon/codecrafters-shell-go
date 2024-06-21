package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var builtin_commands = []string{
	"echo",
	"exit",
	"type",
}

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
		} else if cmd[0] == "type" {
			handle_type(cmd[1:])
		} else {
			fmt.Printf("%s: command not found\n", s)
		}
	}
}

// Handles the `exit` command. The arguments after "exit" should be split and
// passed in.
func handle_exit(args []string) {
	if len(args) == 0 {
		os.Exit(0)
	} else if len(args) == 1 {
		num, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			fmt.Printf("exit: could not parse exit code\n%v\n", err)
			return
		}
		os.Exit(int(num))
	} else {
		fmt.Println("exit: too many arguments")
	}
}

func handle_type(args []string) {
	if len(args) == 0 {
		fmt.Println("type: too few arguments")
	} else if len(args) > 1 {
		fmt.Println("type: too many arguments")
	} else {
		_, found := slices.BinarySearch(builtin_commands, args[0]) 
		if found {
			fmt.Printf("%s is a shell builtin\n", args[0])
		} else {
			// Currently nothing else other than shell builtins
			fmt.Printf("%s: not found\n", args[0])
		}
	}
}
