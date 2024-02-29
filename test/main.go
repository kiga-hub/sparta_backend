package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdout)

	for scanner.Scan() { 
		line := scanner.Text()
		fmt.Println("Read from stdout:", line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
	}
}
