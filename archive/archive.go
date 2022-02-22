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
	GetRef(refID string) (*Ref, error)
	GetRefs(opts ...FilterOption) (Refs, error)
	AddRef(ref Ref) error
}
