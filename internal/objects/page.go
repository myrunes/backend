package objects

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var Champs = []string{
	"aatrox", "ahri", "akali", "alistar", "amumu", "anivia", "annie", "ashe",
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
	"pyke", "quinn", "rakan", "rammus", "reksai", "renekton", "rengar", "riven",
	"rumble", "ryze", "sejuani", "shaco", "shen", "shyvana", "singed", "sion", "sivir",
	"skarner", "sona", "swain", "sylas", "syndra", "tahm-kench", "taliyah", "talon",
	"taric", "teemo", "thresh", "tristana", "trundle", "tryndamere", "twisted-fate",
	"twitch", "udyr", "urgot", "varus", "vayne", "veigar", "velkoz", "vi", "viktor",
	"vladimir", "volibear", "warwick", "wukong", "xayah", "xerath", "xin-zhao", "yasuo",
	"yuumi", "zac", "zed", "ziggs", "zilean", "zoe", "zyra", "diana", "katarina",
	"kayle", "soraka", "yorick",
}

var RunesPrinary = map[string][][]string{
	"domination": [][]string{
		[]string{"electrocute", "predator", "dark-harvest", "hail-of-blades"},
		[]string{"cheap-shot", "taste-of-blood", "sudden-impact"},
		[]string{"zombie-ward", "ghost-poro", "eyeball-collection"},
		[]string{"ravenous-hunter", "ingenious-hunter", "relentless-hunter", "ultimate-hunter"},
	},
	"inspiration": [][]string{
		[]string{"unsealed-spellbook", "glacial-augment", "kleptomancy"},
		[]string{"hextech-flashtraption", "magical-footwear", "perfect-timing"},
		[]string{"futures-market", "minion-dematerializer", "biscuit-delivery"},
		[]string{"cosmic-insight", "approach-velocity", "time-warp-tonic"},
	},
	"precission": [][]string{
		[]string{"press-the-attack", "lethal-tempo", "fleet-footwork", "conqueror"},
		[]string{"overheal", "triumph", "presence-of-mind"},
		[]string{"legend-alacrity", "legend-tenacity", "legend-bloodline"},
		[]string{"coup-de-grace", "cut-down", "last-stand"},
	},
	"resolve": [][]string{
		[]string{"grasp-of-the-undying", "aftershock", "guardian"},
		[]string{"demolish", "font-of-life", "shield-bash"},
		[]string{"conditioning", "second-wind", "bone-plating"},
		[]string{"overgrowth", "revitalize", "unflinching"},
	},
	"sorcery": [][]string{
		[]string{"summon-aery", "arcane-comet", "phase-rush"},
		[]string{"nullifying-orb", "manaflow-band", "nimbus-cloak"},
		[]string{"transcendence", "celerity", "absolute-focus"},
		[]string{"scorch", "waterwalking", "gathering-storm"},
	},
}

var RunesSecondary = map[string][][]string{
	"domination": [][]string{
		[]string{"cheap-shot", "taste-of-blood", "sudden-impact"},
		[]string{"zombie-ward", "ghost-poro", "eyeball-collection"},
		[]string{"ravenous-hunter", "ingenious-hunter", "relentless-hunter", "ultimate-hunter"},
	},
	"inspiration": [][]string{
		[]string{"hextech-flashtraption", "magical-footwear", "perfect-timing"},
		[]string{"futures-market", "minion-dematerializer", "biscuit-delivery"},
		[]string{"cosmic-insight", "approach-velocity", "time-warp-tonic"},
	},
	"precission": [][]string{
		[]string{"overheal", "triumph", "presence-of-mind"},
		[]string{"legend-alacrity", "legend-tenacity", "legend-bloodline"},
		[]string{"coup-de-grace", "cut-down", "last-stand"},
	},
	"resolve": [][]string{
		[]string{"demolish", "font-of-life", "shield-bash"},
		[]string{"conditioning", "second-wind", "bone-plating"},
		[]string{"overgrowth", "revitalize", "unflinching"},
	},
	"sorcery": [][]string{
		[]string{"nullifying-orb", "manaflow-band", "nimbus-cloak"},
		[]string{"transcendence", "celerity", "absolute-focus"},
		[]string{"scorch", "waterwalking", "gathering-storm"},
	},
}

var PerksPool = [][]string{
	[]string{"diamond", "axe", "time"},
	[]string{"diamond", "shield", "circle"},
	[]string{"heart", "shield", "circle"},
}

type Page struct {
	UID       snowflake.ID   `json:"uid"`
	Title     string         `json:"Title"`
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
	Rows [3]string `json:"rows"`
}

type Perks struct {
	Rows [3]string `json:"rows"`
}
