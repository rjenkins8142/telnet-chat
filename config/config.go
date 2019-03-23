package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ParseConfig handles initialization of the config variables. Panics if any
// error occurs while parsing the config file.
func ParseConfig(filepath string) {
	viper.SetConfigFile(filepath)
	viper.AddConfigPath(".")

	// Map to new style automatic variables
	viper.SetEnvPrefix("tchat")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Unable to read config: [%s]", err))
	}
}
