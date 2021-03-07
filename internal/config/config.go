package config

type (
	Tunnel struct {
		Name       string `mapstructure:"name"`
		SSHAddr    string `mapstructure:"sshAddr"`
		RemoteAddr string `mapstructure:"remoteAddr"`
		LocalAddr  string `mapstructure:"localAddr"`
		Auth       Auth   `mapstructure:"auth"`
	}

	Auth struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		KeyPath  string `mapstructure:"keyPath"`
	}

	Config struct {
		Tunnels []Tunnel `mapstructure:"tunnels"`
	}
)
