package objects

import (
	"errors"
	"time"

	"github.com/myrunes/backend/internal/ddragon"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/static"
)

var pageIDCluster, _ = snowflake.NewNode(static.ClusterIDPages)

var (
	ErrInvalidChamp = errors.New("invalid champion")

	errInvalidTree    = errors.New("invalid tree")
	errInvalidPriRune = errors.New("invalid primary rune")
	errInvalidSecRune = errors.New("invalid secondary rune")
	errInvalidPerk    = errors.New("invalid perk")
	errInvalidTitle   = errors.New("invalid title")
)

var PerksPool = [][]string{
	{"diamond", "axe", "time"},
	{"diamond", "shield", "circle"},
	{"heart", "shield", "circle"},
}

type Page struct {
	UID       snowflake.ID   `json:"uid"`
	Owner     snowflake.ID   `json:"owner"`
	Title     string         `json:"title"`
	Created   time.Time      `json:"created"`
	Edited    time.Time      `json:"edited"`
	Champions []string       `json:"champions"`
	Primary   *PrimaryTree   `json:"primary"`
	Secondary *SecondaryTree `json:"secondary"`
	Perks     *Perks         `json:"perks"`
}

type PrimaryTree struct {
	Tree string    `json:"tree"`
	Rows [4]string `json:"rows"`
}

type SecondaryTree struct {
	Tree string    `json:"tree"`
	Rows [2]string `json:"rows"`
}

type Perks struct {
	Rows [3]string `json:"rows"`
}

func NewEmptyPage() *Page {
	return &Page{
		Champions: make([]string, 0),
		Primary: &PrimaryTree{
			Rows: [4]string{},
		},
		Secondary: &SecondaryTree{
			Rows: [2]string{},
		},
		Perks: &Perks{
			Rows: [3]string{},
		},
	}
}

func (p *Page) Validate() error {
	// Check for Title
	if p.Title == "" || len(p.Title) > 1024 {
		return errInvalidTitle
	}

	// Check if primary and secondary tree are the same,
	// which is not allowed
	if p.Secondary.Tree == p.Primary.Tree {
		return errInvalidTree
	}

	// Get Primary and Secondary Tree Objects From ddragon
	// instance by rune tree UIDs
	var primaryTree, secondaryTree *ddragon.RuneTree
	for _, tree := range ddragon.DDragonInstance.Runes {
		if tree.UID == p.Primary.Tree {
			primaryTree = tree
		} else if tree.UID == p.Secondary.Tree {
			secondaryTree = tree
		}
	}

	// If no passing trees could be matched, this
	// is an invalid tree request
	if primaryTree == nil || secondaryTree == nil {
		return errInvalidTree
	}

	// Check if more rune rows are passed as possible
	// rune slots are available
	if len(p.Primary.Rows) > len(primaryTree.Slots) {
		return errInvalidTree
	}

	// Check if primary selected runes exist by UID
	for i, row := range p.Primary.Rows {
		var exists bool
		for _, r := range primaryTree.Slots[i].Runes {
			if r.UID == row {
				exists = true
			}
		}
		if !exists {
			return errInvalidPriRune
		}
	}

	// Check if secondary selected runes exist by
	// UID and check if count is equal 2, else this
	// rune tree is invalid.
	sec := 0
	for _, row := range secondaryTree.Slots {
		for _, ru := range row.Runes {
			var exists bool
			for _, r := range p.Secondary.Rows {
				if r == ru.UID {
					exists = true
				}
			}
			if exists {
				sec++
				break
			}
		}
	}
	if sec != 2 {
		return errInvalidSecRune
	}

	// Check perks
	for i, row := range p.Perks.Rows {
		var exists bool
		for _, p := range PerksPool[i] {
			if row == p {
				exists = true
			}
		}
		if !exists {
			return errInvalidPerk
		}
	}

	// Check if listed champions exists by their
	// champion UIDs
	champMap := map[string]interface{}{}
	for _, champ := range p.Champions {
		var exists bool
		for _, c := range ddragon.DDragonInstance.Champions {
			if champ == c.UID {
				exists = true
			}
		}
		if !exists {
			return ErrInvalidChamp
		}

		champMap[champ] = nil
	}

	champs := make([]string, len(champMap))
	i := 0
	for k := range champMap {
		champs[i] = k
		i++
	}

	p.Champions = champs

	return nil
}

func (p *Page) FinalizeCreate(owner snowflake.ID) {
	now := time.Now()
	p.UID = pageIDCluster.Generate()
	p.Owner = owner
	p.Created = now
	p.Edited = now
}

func (p *Page) Update(newPage *Page) {
	p.Edited = time.Now()
	p.Title = newPage.Title
	p.Champions = newPage.Champions
	p.Perks = newPage.Perks
	p.Primary = newPage.Primary
	p.Secondary = newPage.Secondary
}
