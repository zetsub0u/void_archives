package store

type Filter struct {
	Creator *string
	Boss    *string
	Score   *int
}

type FilterOption func(*Filter)

func (f *Filter) applyOpts(opts ...FilterOption) {
	for _, opt := range opts {
		opt(f)
	}
}

func ByBoss(name string) FilterOption {
	return func(filter *Filter) { filter.Boss = &name }
}

func ByScore(score int) FilterOption {
	return func(filter *Filter) { filter.Score = &score }
}

func ByCreator(name string) FilterOption {
	return func(filter *Filter) { filter.Creator = &name }
}
