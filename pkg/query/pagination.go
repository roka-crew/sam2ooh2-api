package query

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type Page struct {
	Limit  int `json:"limit" form:"limit"`   // 한 페이지당 개수
	Offset int `json:"offset" form:"offset"` // 조회 시작점
}

func NewPage(limit, offset int) Page {
	if limit <= 0 {
		limit = DefaultPageSize
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}
	if offset < 0 {
		offset = 0
	}

	return Page{
		Limit:  limit,
		Offset: offset,
	}
}

func (p Page) GetLimit() int {
	if p.Limit <= 0 {
		return DefaultPageSize
	}
	if p.Limit > MaxPageSize {
		return MaxPageSize
	}
	return p.Limit
}

func (p Page) GetOffset() int {
	if p.Offset < 0 {
		return 0
	}
	return p.Offset
}
