package archive

type Ref struct {
	URL      string
	Creator  string
	Runs     Runs
	Parsed   bool
	Verified bool
}

type Refs []Ref

type Archive interface {
	GetRefs() Refs
	AddRef(ref Ref) error
}
