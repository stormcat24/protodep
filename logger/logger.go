package logger

import "github.com/fatih/color"

func Info(format string, a ...interface{}) {
	color.Green("[INFO] " + format, a)
}

func Error(format string, a ...interface{}) {
	color.Red("[ERROR] " + format, a)
}