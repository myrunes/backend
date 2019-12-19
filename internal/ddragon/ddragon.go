package ddragon

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	epVersions  = "https://ddragon.leagueoflegends.com/api/versions.json"
	epChampions = "https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json"
	epRunes     = "https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/runesReforged.json"
)

type DDragon struct {
	Version   string      `json:"version"`
	Champions []*Champion `json:"champions"`
	Runes     []*RuneTree `json:"runes"`
}

func Poll(version string) (d *DDragon, err error) {
	d = new(DDragon)

	if d.Version, err = getVersion(version); err != nil {
		return
	}

	if d.Champions, err = getChampions(d.Version); err != nil {
		return
	}

	if d.Runes, err = getRunes(d.Version); err != nil {
		return
	}

	return
}

func getJSON(url string, d interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("status code was %d", res.StatusCode)
	}

	dec := json.NewDecoder(res.Body)
	return dec.Decode(d)
}

func getVersions() (res []string, err error) {
	err = getJSON(epVersions, &res)
	return
}

func getVersion(v string) (string, error) {
	versions, err := getVersions()
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

func getChampions(v string) ([]*Champion, error) {
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

func getRunes(v string) (res []*RuneTree, err error) {
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
