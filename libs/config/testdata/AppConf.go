package testdata

type AppConf struct {
	DataConf `yaml:"app"`
}

type DataConf struct {
	Data string `yaml:"data" envconfig:"DATA_APP"`
	Host string `yaml:"host" default:"localhost"`
}
