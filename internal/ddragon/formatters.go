package ddragon

import "strings"

func championUIDFormatter(name string) string {
	name = strings.ToLower(name)
	name = strings.Replace(name, " ", "-", -1)
	name = strings.Replace(name, "'", "", -1)
	name = strings.Replace(name, ".", "", -1)

	return name
}

func runeTreeUIDFormatter(name string) string {
	return strings.ToLower(name)
}

func runeUIDFormatter(name string) string {
	name = strings.ToLower(name)
	name = strings.Replace(name, " ", "-", -1)

	return name
}
