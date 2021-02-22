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
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}

	Config struct {
		Tunnels []Tunnel `mapstructure:"tunnels"`
	}
)
