package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("expecting first argument to be a command")
	}

	switch os.Args[1] {
	case "db:codegen":
		dbCodeGen()
	case "db:seed":
		dbSeed()
	default:
		log.Fatalf("unknown command %q", os.Args[1])
	}
}
