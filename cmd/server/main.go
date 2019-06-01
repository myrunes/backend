package main

import (
	"github.com/zekroTJA/lol-runes/internal/database"
	"github.com/zekroTJA/lol-runes/internal/logger"
)

func main() {
	logger.Setup(`%{color}▶  %{level:.4s} %{id:03d}%{color:reset} %{message}`, 5)

	db := new(database.MongoDB)
}
