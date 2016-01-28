package main

import (
	"os"

	"github.com/go-kit/kit/log"
)

var Logger log.Logger

func init() {
	Logger = log.NewLogfmtLogger(os.Stderr)
}
