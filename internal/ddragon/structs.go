package ddragon

type championsWrapper struct {
	Data map[string]*Champion `json:"data"`
}

type Champion struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type RuneTree struct {
	UID   string      `json:"uid"`
	Name  string      `json:"name"`
	Slots []*RuneSlot `json:"slots"`
}

type RuneSlot struct {
	Runes []*Rune `json:"runes"`
}

type Rune struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	ShortDesc string `json:"shortDesc"`
	LongDesc  string `json:"longDesc"`
}
