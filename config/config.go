package config

import (
	"fmt"
	"log"
	"os"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/spf13/viper"
)

var (
	Main_config *viper.Viper
	DB_PATH,
	DB_FILE,
	LABEL_PATH,
	COA_TEMPLATE_PATH,
	COA_FILEPATH,
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

	// Set defaults
	viper_config.SetDefault("db_path", ".")
	viper_config.SetDefault("label_path", ".")
	viper_config.SetDefault("coa_template_path", ".")
	viper_config.SetDefault("coa_filepath", ".")
	viper_config.SetDefault("log_file", fmt.Sprintf("./%s.log", appname))
	viper_config.SetDefault("json_paths", []string{"."})
	viper_config.SetDefault("font_size", GUI.BASE_FONT_SIZE)

	err = viper_config.ReadInConfig() // Find and read the config file
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// Config file not found; ignore error if desired

		os.MkdirAll(config_path, 660)
		log.Println(viper_config.WriteConfigAs(config_file))
	} else if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	DB_FILE = viper_config.GetString("db_path") + "/qc.sqlite3"
	LABEL_PATH = viper_config.GetString("label_path")
	JSON_PATHS = viper_config.GetStringSlice("json_paths")
	COA_TEMPLATE_PATH = viper_config.GetString("coa_template_path")
	COA_FILEPATH = viper_config.GetString("coa_filepath")
	LOG_FILE = viper_config.GetString("log_file")
	GUI.BASE_FONT_SIZE = viper_config.GetInt("font_size")

	return viper_config
}

func Write_config(viper_config *viper.Viper) {
	viper_config.WriteConfigAs(viper_config.ConfigFileUsed())
}
