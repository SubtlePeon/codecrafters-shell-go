package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// Builtin commands. Needs to be updated with every new command. Must be sorted.
// Cannot be const unfortunately.
var builtin_commands = [...]string{
	"echo",
	"exit",
	"pwd",
	"type",
}

func main() {
	for true {
		// Prompt
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		s, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			// End shell on end of file
			if err.Error() == "EOF" {
				// Add extra newline for formatting reasons
				fmt.Println()
				os.Exit(0)
			}

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
		} else if cmd[0] == "pwd" {
			if len(cmd) > 1 {
				fmt.Fprintln(os.Stderr, "pwd: too many arguments")
				continue
			}
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: pwd: %v\n", err)
				continue
			}
			fmt.Println(cwd)
		} else if cmd[0] == "type" {
			handle_type(cmd[1:])
		} else if cmd_abspath := find_executable(cmd[0]); cmd_abspath != "" {
			command := exec.Command(cmd[0], cmd[1:]...)
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			command.Run()
		} else {
			fmt.Fprintf(os.Stderr, "%s: command not found\n", s)
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

// Handles the `type` command.
func handle_type(args []string) {
	if len(args) == 0 {
		fmt.Println("type: too few arguments")
	} else if len(args) > 1 {
		fmt.Println("type: too many arguments")
	} else {
		_, found := slices.BinarySearch(builtin_commands[:], args[0])
		if found {
			fmt.Printf("%s is a shell builtin\n", args[0])
		} else if abspath := find_executable(args[0]); abspath != "" {
			fmt.Printf("%s is %s\n", args[0], abspath)
		} else {
			fmt.Printf("%s: not found\n", args[0])
		}
	}
}

// Finds an executable with the `PATH` environmental variable. If the
// executable cannot be found, returns an empty string. Otherwise, returns
// the absolute path to the executable.
// - Finds the first match
func find_executable(command string) string {
	env_path := strings.Split(os.Getenv("PATH"), ":")
	for _, walk_path := range env_path {
		if !filepath.IsAbs(walk_path) {
			new_path, err := filepath.Abs(walk_path)
			if err != nil {
				// Something went wrong, skip this one
				fmt.Fprintf(
					os.Stderr,
					"Couldn't convert '%s' to abspath",
					walk_path,
				)
				continue
			}
			walk_path = new_path
		}

		file_info, err := os.Stat(walk_path)
		if err != nil {
			// Probably invalid path, check next path
			continue
		}

		// Only check directories
		if !file_info.IsDir() {
			continue
		}

		files, err := os.ReadDir(walk_path)
		// Ignore errors
		if err != nil {
		}
		for _, file := range files {
			if file.Name() == command {
				// We found it, return early.
				return filepath.Join(walk_path, file.Name())
			}
		}
	}
	return ""
}
