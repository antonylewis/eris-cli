package definitions

type RemoteDefinition struct {
	Name     string   `json:"name" yaml:"name" toml:"name"`
	Nodes    int      `json:"nodes" yaml:"nodes" toml:"nodes"`
	Driver   string   `json:"driver" yaml:"driver" toml:"driver"`
	Regions  []string `mapstructure:"regions" json:"regions,omitempty" yaml:"regions,omitempty" toml:"regions,omitempty"`
	Services []string `mapstructure:"services" json:"services,omitempty" yaml:"services,omitempty" toml:"services,omitempty"`
	Machines []string `mapstructure:"machines" json:"machines,omitempty" yaml:"machines,omitempty" toml:"machines,omitempty"` //should this be a pointer?

	Maintainer *Maintainer `json:"maintainer,omitempty" yaml:"maintainer,omitempty" toml:"maintainer,omitempty"`
}

func BlankRemoteDefinition() *RemoteDefinition {
	return &RemoteDefinition{
		Maintainer: BlankMaintainer(),
	}
}
