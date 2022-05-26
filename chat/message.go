package main

import (
	"time"
)

// messageは1つのメッセージを表す
type message struct {
	Name    string
	Message string
	When    time.Time
}
