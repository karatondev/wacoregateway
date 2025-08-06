//go:build !mock
// +build !mock

package amqpx

import "log"

var Debug bool

func debug(args ...interface{}) {
	if !Debug {
		return
	}
	log.Print(args...)
}

func debugf(format string, args ...interface{}) {
	if !Debug {
		return
	}
	log.Printf(format, args...)
}
