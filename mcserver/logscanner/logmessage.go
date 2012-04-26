package logscanner

import (
    "time"
)

// Type LogLevel represents a log level.
type LogLevel uint

// LogLevel constants
const (
    _ LogLevel = iota

    Info
    Warning
    Severe
)

// Type LogMessage represents a log message.
type LogMessage struct {
    message string
    level   LogLevel
    date    time.Time
}

// Function Message returns the actual content of the message.
func (msg *LogMessage) Message() (message string) {
    return msg.message
}

// Function Level returns the log level of the message.
func (msg *LogMessage) Level() (level LogLevel) {
    return msg.level
}

// Function Date returns the date/time at which the message was produced.
func (msg *LogMessage) Date() (date time.Time) {
    return msg.date
}
