package config

import (
	"fmt"
	"log"
	"os"

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

func Load_config(appname string) *viper.Viper {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	config_name := fmt.Sprintf("config_%s.toml", appname)
	config_path := fmt.Sprintf("%s/.config/%s", home, appname)
	config_file := fmt.Sprintf("%s/%s", config_path, config_name)
	viper_config := viper.New()
	viper_config.SetConfigName(config_name) // name of config file (without extension)
	viper_config.SetConfigType("toml")      // REQUIRED if the config file does not have the extension in the name
	viper_config.AddConfigPath(".")         // optionally look for config in the working directory
	// viper_config.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper_config.AddConfigPath(config_path) // call multiple times to add many search paths
	err = viper_config.ReadInConfig()       // Find and read the config file
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// Config file not found; ignore error if desired

		viper_config.Set("db_path", ".")
		viper_config.Set("label_path", ".")
		viper_config.Set("log_file", fmt.Sprintf("./%s.log", appname))
		viper_config.Set("json_paths", []string{"."})

		os.MkdirAll(config_path, 660)
		log.Println(viper_config.WriteConfigAs(config_file))
	} else if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	DB_FILE = viper_config.GetString("db_path") + "/qc.sqlite3"
	LABEL_PATH = viper_config.GetString("label_path")
	JSON_PATHS = viper_config.GetStringSlice("json_paths")
	LOG_FILE = viper_config.GetString("log_file")

	return viper_config
}
