package main

import cmd "github.com/mrclmr/icm/cmd/icm"

var version = "dev"

func main() {
	cmd.Execute(version)
}
