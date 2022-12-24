package hashref

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/NodyHub/hashref/pkg/util"
)

type Config struct {
	Publisher     string            `json:"HASHREF_PUBLISHER"`
	DefaultMeta   map[string]string `json:"HASHREF_DEFAULT_META"`
	HashrefServer string            `json:"HASHREF_SERVER"`
}

// LoadConfig loads the configuration from the provided file path
func LoadConfig(cfgFileName string) Config {
	log.Println("Try to load hashref configuration json")

	// figure out home dir and use that as cfg path
	if len(cfgFileName) == 0 {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Println("Could not figure out USER_HOME")
			return getDefaultConfig()
		}
		cfgFileName = filepath.Join(dirname, ".hashref")
	}

	// create empty default cfg
	cfg := NewConfig()

	// load file
	if cfgFileName != "" {
		cfgFile, err := ioutil.ReadFile(cfgFileName)
		if err != nil {
			log.Printf("ERROR: %v\n", err) // Cannot read config: %v\n", CLI.Config)
			log.Println("Loading default config")
		} else {
			err := json.Unmarshal([]byte(cfgFile), &cfg)
			if err != nil {
				log.Printf("%v\n", err)
				log.Println("Loading default config")
			} else {
				log.Printf("Loading config %v successfull!\n", cfgFileName)
			}
		}
	}

	// finalize
	return cfg
}

// NewConfig returns a Config object with default values
func NewConfig() Config {
	return getDefaultConfig()
}

// LoadEnvValues loads Config object values based on the json field
// names from the environment
func (c *Config) LoadEnvValues() {
	log.Println("Check env for configuration")
	for _, key := range GetJsonFields() {
		if value := os.Getenv(key); value != "" {
			log.Printf("Found '%v' in env\n", key)
			util.SetValueInStructByJsonKey(c, key, value)
		}
	}

}

// getDefaultConfig returns a Config object with default values
func getDefaultConfig() Config {
	return Config{
		Publisher:     "anonymous",
		DefaultMeta:   map[string]string{},
		HashrefServer: "http://127.0.0.1:8080",
	}
}

// GetJsonFields returns a slice of strings with json field names
// for a Config object
func GetJsonFields() (fields []string) {
	cfg := getDefaultConfig()
	bCfg, _ := json.Marshal(cfg)
	keyMap := map[string]string{}
	json.Unmarshal(bCfg, &keyMap)
	for key := range keyMap {
		fields = append(fields, key)
	}
	return fields
}
