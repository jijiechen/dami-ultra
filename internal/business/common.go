package business

type Pager struct {
	Page     int32 `json:"page,omitempty" form:"page,omitempty"`
	PageSize int32 `json:"page_size,omitempty" form:"page_size,omitempty"`
	Total    int32 `json:"total,omitempty" form:"total,omitempty"`
}

func GetPager(pager *Pager) *Pager {
	if pager == nil {
		pager = &Pager{
			Page:     1,
			PageSize: 10,
		}
	}
	return pager
}
