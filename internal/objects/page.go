package objects

import (
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/myrunes/internal/static"
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

var Champs = []string{
	"aatrox", "ahri", "akali", "alistar", "amumu", "anivia", "annie", "aphelios", "ashe",
	"aurelion-sol", "azir", "bard", "blitzcrank", "brand", "braum", "caitlyn",
	"camille", "cassiopeia", "chogath", "corki", "darius", "dr-mundo", "draven",
	"ekko", "elise", "evelynn", "ezreal", "fiddlesticks", "fiora", "fizz", "galio",
	"gangplank", "garen", "gnar", "gragas", "graves", "hecarim", "heimerdinger",
	"illaoi", "irelia", "ivern", "janna", "jarvan-iv", "jax", "jayce", "jhin",
	"jinx", "kaisa", "kalista", "karma", "karthus", "kassadin", "kayn", "kennen",
	"khazix", "kindred", "kled", "kogmaw", "leblanc", "lee-sin", "leona", "lissandra",
	"lucian", "lulu", "lux", "malphite", "malzahar", "maokai", "master-yi",
	"miss-fortune", "mordekaiser", "morgana", "nami", "nasus", "nautilus", "neeko",
	"nidalee", "nocturne", "nunu", "olaf", "orianna", "ornn", "pantheon", "poppy",
	"pyke", "qiyana", "quinn", "rakan", "rammus", "reksai", "renekton", "rengar", "riven",
	"rumble", "ryze", "sejuani", "senna", "shaco", "shen", "shyvana", "singed", "sion", "sivir",
	"skarner", "sona", "swain", "sylas", "syndra", "tahm-kench", "taliyah", "talon",
	"taric", "teemo", "thresh", "tristana", "trundle", "tryndamere", "twisted-fate",
	"twitch", "udyr", "urgot", "varus", "vayne", "veigar", "velkoz", "vi", "viktor",
	"vladimir", "volibear", "warwick", "wukong", "xayah", "xerath", "xin-zhao", "yasuo",
	"yuumi", "zac", "zed", "ziggs", "zilean", "zoe", "zyra", "diana", "katarina",
	"kayle", "soraka", "yorick",
}

var RuneTrees = []string{
	"precission", "domination", "sorcery", "resolve", "inspiration",
}

var RunesPrimary = map[string][][]string{
	"precission": [][]string{
		[]string{"press-the-attack", "lethal-tempo", "fleet-footwork", "conqueror"},
		[]string{"overheal", "triumph", "presence-of-mind"},
		[]string{"legend-alacrity", "legend-tenacity", "legend-bloodline"},
		[]string{"coup-de-grace", "cut-down", "last-stand"},
	},
	"domination": [][]string{
		[]string{"electrocute", "predator", "dark-harvest", "hail-of-blades"},
		[]string{"cheap-shot", "taste-of-blood", "sudden-impact"},
		[]string{"zombie-ward", "ghost-poro", "eyeball-collection"},
		[]string{"ravenous-hunter", "ingenious-hunter", "relentless-hunter", "ultimate-hunter"},
	},
	"sorcery": [][]string{
		[]string{"summon-aery", "arcane-comet", "phase-rush"},
		[]string{"nullifying-orb", "manaflow-band", "nimbus-cloak"},
		[]string{"transcendence", "celerity", "absolute-focus"},
		[]string{"scorch", "waterwalking", "gathering-storm"},
	},
	"resolve": [][]string{
		[]string{"grasp-of-the-undying", "aftershock", "guardian"},
		[]string{"demolish", "font-of-life", "shield-bash"},
		[]string{"conditioning", "second-wind", "bone-plating"},
		[]string{"overgrowth", "revitalize", "unflinching"},
	},
	"inspiration": [][]string{
		[]string{"unsealed-spellbook", "glacial-augment", "prototype-omnistone"},
		[]string{"hextech-flashtraption", "magical-footwear", "perfect-timing"},
		[]string{"futures-market", "minion-dematerializer", "biscuit-delivery"},
		[]string{"cosmic-insight", "approach-velocity", "time-warp-tonic"},
	},
}

var RunesSecondary = map[string][][]string{
	"precission": [][]string{
		[]string{"overheal", "triumph", "presence-of-mind"},
		[]string{"legend-alacrity", "legend-tenacity", "legend-bloodline"},
		[]string{"coup-de-grace", "cut-down", "last-stand"},
	},
	"domination": [][]string{
		[]string{"cheap-shot", "taste-of-blood", "sudden-impact"},
		[]string{"zombie-ward", "ghost-poro", "eyeball-collection"},
		[]string{"ravenous-hunter", "ingenious-hunter", "relentless-hunter", "ultimate-hunter"},
	},
	"sorcery": [][]string{
		[]string{"nullifying-orb", "manaflow-band", "nimbus-cloak"},
		[]string{"transcendence", "celerity", "absolute-focus"},
		[]string{"scorch", "waterwalking", "gathering-storm"},
	},
	"resolve": [][]string{
		[]string{"demolish", "font-of-life", "shield-bash"},
		[]string{"conditioning", "second-wind", "bone-plating"},
		[]string{"overgrowth", "revitalize", "unflinching"},
	},
	"inspiration": [][]string{
		[]string{"hextech-flashtraption", "magical-footwear", "perfect-timing"},
		[]string{"futures-market", "minion-dematerializer", "biscuit-delivery"},
		[]string{"cosmic-insight", "approach-velocity", "time-warp-tonic"},
	},
}

var PerksPool = [][]string{
	[]string{"diamond", "axe", "time"},
	[]string{"diamond", "shield", "circle"},
	[]string{"heart", "shield", "circle"},
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
	if p.Title == "" {
		return errInvalidTitle
	}

	primaryTree, ok := RunesPrimary[p.Primary.Tree]
	if !ok {
		return errInvalidTree
	}

	secondaryTree, ok := RunesSecondary[p.Secondary.Tree]
	if !ok {
		return errInvalidTree
	}

	if p.Secondary.Tree == p.Primary.Tree {
		return errInvalidTree
	}

	for i, row := range p.Primary.Rows {
		var exists bool
		for _, r := range primaryTree[i] {
			if r == row {
				exists = true
			}
		}
		if !exists {
			return errInvalidPriRune
		}
	}

	sec := 0
	for _, row := range secondaryTree {
		for _, ru := range row {
			var ex bool
			for _, r := range p.Secondary.Rows {
				if r == ru {
					ex = true
				}
			}
			if ex {
				sec++
				break
			}
		}
	}
	if sec != 2 {
		return errInvalidSecRune
	}

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

	champMap := map[string]interface{}{}
	for _, champ := range p.Champions {
		var exists bool
		for _, c := range Champs {
			if champ == c {
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
