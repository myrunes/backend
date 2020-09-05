// Package ddragon provides bindings to the
// Riots data dragon CDN.
package ddragon

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// endpoint definitions
const (
	epVersions  = "https://ddragon.leagueoflegends.com/api/versions.json"
	epChampions = "https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json"
	epRunes     = "https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/runesReforged.json"
)

// Fetch collects version, champion and rune
// information from the Datadragon API and
// wraps them into a DDragon object returned.
func Fetch(version string) (d *DDragon, err error) {
	d = new(DDragon)

	if d.Version, err = GetVersion(version); err != nil {
		return
	}

	if d.Champions, err = GetChampions(d.Version); err != nil {
		return
	}

	if d.Runes, err = GetRunes(d.Version); err != nil {
		return
	}

	return
}

// GetVersions returns an array of valid LoL patch
// version strings.
func GetVersions() (res []string, err error) {
	err = getJSON(epVersions, &res)
	return
}

// GetVersion validates the given version v
// against the array of valid versions collected
// from the API. If v is empty or equals "latest",
// the most recent version string will be returned.
// If the given version string is invalid, an
// error will be returned.
func GetVersion(v string) (string, error) {
	versions, err := GetVersions()
	if err != nil {
		return "", err
	}

	if v == "" || v == "latest" {
		return versions[0], nil
	}

	for _, ver := range versions {
		if ver == v {
			return v, nil
		}
	}

	return "", fmt.Errorf("invalid version")
}

// GetChampions returns an array of Champion objects
// collected from the datadragon API.
func GetChampions(v string) ([]*Champion, error) {
	res := new(championsWrapper)
	err := getJSON(fmt.Sprintf(epChampions, v), res)
	if err != nil {
		return nil, err
	}

	champs := res.Data
	fChamps := make([]*Champion, len(champs))
	var i int
	for _, c := range champs {
		c.UID = championUIDFormatter(c.Name)
		fChamps[i] = c
		i++
	}

	return fChamps, nil
}

// GetRunes returns an array of RuneTree objects
// collected from the datadragon API.
func GetRunes(v string) (res []*RuneTree, err error) {
	err = getJSON(fmt.Sprintf(epRunes, v), &res)

	for _, tree := range res {
		tree.UID = runeTreeUIDFormatter(tree.Name)
		for _, slot := range tree.Slots {
			for _, r := range slot.Runes {
				r.UID = runeUIDFormatter(r.Name)
			}
		}
	}

	return
}

// getJSON executes a GET request on the given URL
// and tries to decode the JSON response body
// into the given object reference v.
func getJSON(url string, v interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("status code was %d", res.StatusCode)
	}

	dec := json.NewDecoder(res.Body)
	return dec.Decode(v)
}
