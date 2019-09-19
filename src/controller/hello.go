package controller

import (
	"fmt"
	"os"
)

func SayHello() (interface{}, error) {
	msg := os.Getenv("MSG")
	return fmt.Sprintf("Hello from %s", msg), nil
}