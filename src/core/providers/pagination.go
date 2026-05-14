package providers

type Pagination struct {
	Limit  int
	Offset int
}

const (
	DefaultPageLimit = 50
	MaxPageLimit     = 200
)

func (p Pagination) Normalize() Pagination {
	if p.Limit <= 0 {
		p.Limit = DefaultPageLimit
	}
	if p.Limit > MaxPageLimit {
		p.Limit = MaxPageLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	return p
}
