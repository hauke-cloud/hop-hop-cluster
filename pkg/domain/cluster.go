package domain

type Status int

const (
	Waiting Status = iota
	Error
	Initializing
	Initialized
	Joining
	Joined
)

func (s Status) String() string {
	return [...]string{"Waiting", "Error", "Initializing", "Initialized", "Joining", "Joined"}[s]
}

type Cluster struct {
	Name      string `validate:"-" yaml:"name" json:"name"`
	IPAddress string `validate:"-" yaml:"ip_address" json:"ip_address"`
	Priority  int    `validate:"-" yaml:"priority" json:"priority"`
	Status    Status `validate:"-" yaml:"status" json:"status"`
}
