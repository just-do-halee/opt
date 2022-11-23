package internal

import (
	"log"
	"os"
)

var OptLogger = log.New(os.Stderr, "[OPT] ", log.LstdFlags|log.Llongfile)

//go:inline
func OptFatal(v ...any) {
	OptLogger.Fatal(v...)
}
