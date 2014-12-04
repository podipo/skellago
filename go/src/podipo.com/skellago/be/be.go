package be

import (
	"log"
	"os"
)

var VERSION = "0.1.0"

var logger = log.New(os.Stdout, "[be] ", 0)
