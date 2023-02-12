package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/brer/brer-cli/cmd"
)

func main() {
	var err error = nil

	if len(os.Args) < 2 {
		err = errors.New("expected command")
	} else {
		switch os.Args[1] {
		case "publish":
			err = cmd.Publish()
		case "trigger":
			err = cmd.Trigger()
		default:
			err = errors.New("invalid command")
		}
	}

	if err != nil {
		fmt.Println("error found", err)
		os.Exit(1)
	}
}
