package main

import (
	"fmt"
	"os"
	"testing"
)

func TestM(t *testing.T) {
	fmt.Println(os.Getenv("USER"))
}
