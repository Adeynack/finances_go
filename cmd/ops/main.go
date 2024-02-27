package main

import (
	"log"
	"os"

	"github.com/adeynack/finances/app/appenv"
)

func main() {
	appenv.Init()

	if len(os.Args) < 2 {
		log.Fatalln("expecting first argument to be a command")
	}

	switch os.Args[1] {
	case "db:codegen":
		dbCodeGen()
	default:
		log.Fatalf("unknown command %q", os.Args[1])
	}
}
