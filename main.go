package main

import (
	"fmt"
	"os"
	"quartz/dagbench.io/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
