package be

import (
	"log"
	"mime"
	"os"
	"strings"

	"github.com/nu7hatch/gouuid"
)

var logger = log.New(os.Stdout, "[be] ", 0)

func UUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}

func MimeTypeFromFileName(name string) string {
	lindex := strings.LastIndex(name, ".")
	if lindex == -1 || lindex == len(name)-1 {
		return ""
	}
	return mime.TypeByExtension(name[lindex:])
}
