package main

import (
	"time"
)

func getTime() int64 {
	time := (time.Now()).Unix()
	return time
}
