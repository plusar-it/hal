package config

// Config hold the configuration settings
type Config struct {
	NetworkAddress         string
	FromAddress            string
	NoOfProducers          int
	BatchSize              int
	TriggersSourceFilePath string
}

// Loads the hardcoded configuration, can be extended to be loaded from an external file or so
func LoadConfig() *Config {
	return &Config{
		NetworkAddress:         "https://mainnet.infura.io/v3/17ed7fe26d014e5b9be7dfff5368c69d",
		FromAddress:            "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		NoOfProducers:          3,
		BatchSize:              3,
		TriggersSourceFilePath: "data/balance_triggers.json",
	}
}
