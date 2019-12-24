package main

import (
	"fmt"

	"github.com/myrunes/backend/internal/ddragon"
)

func main() {
	d, _ := ddragon.Poll("latest")
	fmt.Printf(d.Runes[0].Slots[1].Runes[1].UID)
}
