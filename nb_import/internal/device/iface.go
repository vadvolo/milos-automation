package device

type Interface struct {
	Name         string `json:"name"`
	ShortName    string `json:"shortname"`
	Description  string `json:"description"`
	Neighbor     string `json:"neighbor"`
	NeighborPort string `json:"neighbor_port"`
}
