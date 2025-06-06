package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var (
	Main_config *viper.Viper
	DB_PATH,
	DB_FILE,
	LABEL_PATH,
	LOG_FILE string

	JSON_PATHS []string
)

func Load_config() *viper.Viper {
	viper_config := viper.New()
	viper_config.SetConfigName("config") // name of config file (without extension)
	viper_config.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
	viper_config.AddConfigPath(".")      // optionally look for config in the working directory
	// viper_config.AddConfigPath("/etc/appname/")  // path to look for the config file in
	// viper.AddConfigPath("$HOME/.config/qc_data_entry") // call multiple times to add many search paths
	viper_config.AddConfigPath("$HOME/.config/qc_data_entry") // call multiple times to add many search paths
	err := viper_config.ReadInConfig()                        // Find and read the config file
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// Config file not found; ignore error if desired

		viper_config.Set("db_path", ".")
		viper_config.Set("label_path", ".")
		viper_config.Set("log_file", "./qc_data_entry.log")
		viper_config.Set("json_paths", []string{"."})

		log.Println(viper_config.WriteConfigAs("config.toml"))
	} else if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	DB_FILE = viper_config.GetString("db_path") + "/qc.sqlite3"
	LABEL_PATH = viper_config.GetString("label_path")
	JSON_PATHS = viper_config.GetStringSlice("json_paths")
	LOG_FILE = viper_config.GetString("log_file")

	return viper_config
}
