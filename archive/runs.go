package archive

type ValkRank string

const (
	S   ValkRank = "S"
	S1           = "S1"
	S2           = "S2"
	S3           = "S3"
	SS           = "SS"
	SS1          = "SS1"
	SS2          = "SS2"
	SS3          = "SS3"
	SSS          = "SSS"
)

type ElfRank string

const (
	Elf1Star ElfRank = "1*"
	Elf2Star         = "2*"
	Elf3Star         = "3*"
	Elf4Star         = "4*"
)

type Mode string

const (
	Memorial Mode = "MA"
	Abyss         = "Abyss"
)

type Run struct {
	Boss  string
	Mode  Mode
	Score int
	Team  Team
}

type Runs []Run

type Team struct {
	Valk1  Valk
	Valk2  Valk
	Valk3  Valk
	Elf    *Elf
	Keys   *Keys
	Emblem *Emblem
}

type Valk struct {
	Name    string
	Rank    ValkRank
	Gear    Gear
	Comment string
}

type Elf struct {
	Name string
	Rank ElfRank
}

type Keys struct {
}

type Emblem struct {
	Name string
}

type Gear struct {
	Stigma1 Stigma
	Stigma2 Stigma
	Stigma3 Stigma
	Weapon  Weapon
	Comment string
}

type Stigma struct {
	Name   string
	Affix1 string
	Affix2 string
}

type Weapon struct {
	Name  string
	Level int
}
