package model

type Config struct {
	Token         string       `yaml:"Token"`
	HandleDomain  DomainConfig `yaml:"HandleDomain"`
	ResolveDomain DomainConfig `yaml:"ResolveDomain"`
}

type DomainConfig struct {
	RootDomain string `yaml:"RootDomain"`
	SubDomain  string `yaml:"SubDomain"`
}

type DomainInfo struct {
	Status  Status   `json:"status"`
	Records []Record `json:"records"`
}

type Status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Record struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	TTL    string `json:"ttl"`
	Value  string `json:"value"`
	LineID string `json:"line_id"`
}
