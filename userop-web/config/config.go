package config

type ServerConfig struct {
	Name          string         `mapstructure:"name" json:"name"`
	Host          string         `mapstructure:"host" json:"host"`
	Port          int            `mapstructure:"port" json:"port"`
	Tags          []string       `mapstructure:"tags" json:"tags"`
	UserOpSrvInfo GoodsSrvConfig `mapstructure:"userop_srv" json:"userop_srv"`
	JWTInfo       JWTConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulInfo    ConsulConfig   `mapstructure:"consul" json:"consul"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
