package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	fd := int(os.Stdin.Fd())
	tw, th, _ := terminal.GetSize(fd)
	fmt.Printf("%s, %s", tw, th)
}
