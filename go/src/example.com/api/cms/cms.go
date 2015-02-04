package cms

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[cms] ", 0)
