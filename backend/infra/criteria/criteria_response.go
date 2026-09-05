package criteria

type CResponse struct {
	ActivePage int `json:"activePage"`
	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`
}

func (cr CResponse) New(ActivePage, TotalPages, TotalItems int) CResponse {
	return CResponse{ActivePage: ActivePage, TotalItems: TotalItems, TotalPages: TotalPages}
}

// NewWithCriteria
// Faz a inicialização de Active page, Total pages e Total items
func (cr CResponse) NewWithCriteria(c Criteria) CResponse {
	return CResponse{ActivePage: c.ActivePage, TotalItems: c.TotalItems, TotalPages: c.TotalPages}
}
