package ddragon

import "strings"

var uidExceptions = map[string]string{
	"Nunu \u0026 Willump": "nunu",
}

func championUIDFormatter(name string) string {
	if u, ok := uidExceptions[name]; ok {
		return u
	}

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
