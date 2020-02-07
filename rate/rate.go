package rate

type Rate struct {
	USD float64 `json:"USD"`
	MXN float64 `json:"MXN"`
}

type Service interface {
	LatestRate() (*Provider, error)
}

type Provider struct {
	Rate    float64 `json:"value"`
	Updated string  `json:"last_updated"`
	Name    string  `json:"name"`
}
