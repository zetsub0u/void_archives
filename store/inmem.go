package store

import "github.com/zetsub0u/void_archives/archive"

type InMem struct {
	data archive.Refs
}

func NewInMem() *InMem {
	data := make(archive.Refs, 0)
	return &InMem{data: data}
}

func (a *InMem) GetRef(refID string) (*archive.Ref, error) {
	for _, ref := range a.data {
		if ref.URL == refID {
			return &ref, nil
		}
	}
	return nil, nil
}

func (a *InMem) AddRef(ref archive.Ref) error {
	a.data = append(a.data, ref)
	return nil
}

func (a *InMem) GetRefs(opts ...archive.FilterOption) (archive.Refs, error) {
	return a.data, nil
}
