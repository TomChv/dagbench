package main

import "quartz/dagbenchmark.io/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}