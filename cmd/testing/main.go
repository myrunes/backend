package main

import (
	"fmt"

	"github.com/myrunes/backend/pkg/ddragon"
)

func main() {
	d, _ := ddragon.Fetch("latest")
	fmt.Printf(d.Runes[0].Slots[1].Runes[1].UID)
}
