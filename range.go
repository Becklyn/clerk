package clerk

type Range struct {
	SkipValue int
	TakeValue int
}

func NewRange() *Range {
	return &Range{
		SkipValue: 0,
		TakeValue: 0,
	}
}

func (r *Range) Skip(skip int) *Range {
	r.SkipValue = skip
	return r
}

func (r *Range) Take(take int) *Range {
	r.TakeValue = take
	return r
}
