package main

import (
	"fmt"

	"github.com/myrunes/myrunes/pkg/random"
)

func main() {
	fmt.Println(random.String(16, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
}
