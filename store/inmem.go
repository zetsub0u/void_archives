package store

import "github.com/zetsub0u/void_archives/archive"

type InMem struct {
	data archive.Refs
}

func NewInMem() *InMem {
	data := make(archive.Refs, 0)
	return &InMem{data: data}
}

func (a *InMem) GetRefs() archive.Refs {
	return a.data
}

func (a *InMem) AddRef(ref archive.Ref) error {
	a.data = append(a.data, ref)
	return nil
}

func (a *InMem) GetRef(ops ...FilterOption) (archive.Refs, error) {
	f := Filter{}
	f.applyOpts()
	return nil, nil
}
