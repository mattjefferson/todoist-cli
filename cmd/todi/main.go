// Command todi runs the CLI.
package main

import (
	"os"

	"github.com/mattjefferson/todi/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
