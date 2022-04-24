package shows

type Show struct {
	Id       *int64 `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	TotalEps *int64 `json:"totalEps,omitempty"`
	Type     string `json:"type,omitempty"`
}

type ShowFilter struct {
	Name     string `json:"name,omitempty"`
	TotalEps *int64 `json:"totalEps,string,omitempty"`
	Type     string `json:"type,omitempty"`
}
