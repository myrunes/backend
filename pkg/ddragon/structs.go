package ddragon

// DDragon wraps the current LoL patch version and
// information about champions and runes collected
// from Riot's Datadragon API.
type DDragon struct {
	Version   string      `json:"version"`
	Champions []*Champion `json:"champions"`
	Runes     []*RuneTree `json:"runes"`
}

// Champion describes a champion object.
type Champion struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

// RuneTree describes a rune tree and
// contains the rune slots for this tree.
type RuneTree struct {
	UID   string      `json:"uid"`
	Name  string      `json:"name"`
	Slots []*RuneSlot `json:"slots"`
}

// RuneSlot wraps a row of runes.
type RuneSlot struct {
	Runes []*Rune `json:"runes"`
}

// Rune describes the properties of
// a rune in a rune tree row.
type Rune struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	ShortDesc string `json:"shortDesc"`
	LongDesc  string `json:"longDesc"`
}

// championsWrapper describes the response
// model of the champions ddragon API
// endpoint response.
type championsWrapper struct {
	Data map[string]*Champion `json:"data"`
}
