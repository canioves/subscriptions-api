package logger

import (
	"log"
	"os"
)

var (
	Info  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime).Printf
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime).Printf
	Fatal = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime).Fatalf
)
