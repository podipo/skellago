package be

import (
	"log"
	"os"

	"github.com/nu7hatch/gouuid"
)

var logger = log.New(os.Stdout, "[be] ", 0)

func UUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}
