package configure

type Configure struct {
	Project Project `json:"project" yaml:"project"`
}

type Project struct {
	Registries []*Registry `json:"registries" yaml:"registries"`
	Services   []*Service  `json:"services" yaml:"services"`
}
