package definitions

//XXX this'll write to csv on creation
type RemMachine struct {
	Name   string `json:"name" yaml:"name" toml:"name"`
	Driver string `json:"driver" yaml:"driver" toml:"driver"` //will need its own strcut for different types of drivers...
	Region string `json:"region" yaml:"region" toml:"region"` //will need its own strcut for different types of drivers...
	IP     string `json:"ip" yaml:"ip" toml:"ip"`
}

func BlankRemMachine() *RemMachine {
	return &RemMachine{}
}

//RemMachines []string `mapstructure:"machines" json:"machines,omitempty" yaml:"machines,omitempty" toml:"machines,omitempty"`
