package config

import (
	"fmt"
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
	RETAIN_FILE_NAME,
	RETAIN_WORKSHEET_NAME,
	PRODUCTION_SCHEDULE_FILE_NAME,
	PRODUCTION_SCHEDULE_WORKSHEET_NAME,
	LOG_FILE string

	JSON_PATHS []string
)

func load_config_base(appname string) (viper_config *viper.Viper, config_path, config_file string) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	config_name := fmt.Sprintf("config_%s.toml", appname)
	config_path = fmt.Sprintf("%s/.config/%s", home, appname)
	config_file = fmt.Sprintf("%s/%s", config_path, config_name)
	viper_config = viper.New()
	viper_config.SetConfigName(config_name) // name of config file (without extension)
	viper_config.SetConfigType("toml")      // REQUIRED if the config file does not have the extension in the name
	viper_config.AddConfigPath(".")         // optionally look for config in the working directory
	// viper_config.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper_config.AddConfigPath(config_path) // call multiple times to add many search paths
	return viper_config, config_path, config_file
}

func set_config_defaults(appname string, viper_config *viper.Viper) {
	// Set defaults
	viper_config.SetDefault("db_path", ".")
	viper_config.SetDefault("label_path", ".")
	viper_config.SetDefault("coa_template_path", ".")
	viper_config.SetDefault("coa_filepath", ".")
	viper_config.SetDefault("log_file", fmt.Sprintf("./%s.log", appname))
	viper_config.SetDefault("json_paths", []string{"."})
	viper_config.SetDefault("font_size", GUI.BASE_FONT_SIZE)
}

func set_config_defaults_entry(appname string, viper_config *viper.Viper) {
	// Set defaults
	set_config_defaults(appname, viper_config)
	viper_config.SetDefault("retain_file_name", "RETAIN-SAMPLE-TRACKING.xlsx")
	viper_config.SetDefault("retain_worksheet_name", "Sheet1")
}

func set_config_defaults_inbound(appname string, viper_config *viper.Viper) {
	// Set defaults
	set_config_defaults(appname, viper_config)
	viper_config.SetDefault("production_schedule_file_name", "PRODUCTION-SCHEDULE.xlsx")
	viper_config.SetDefault("production_schedule_worksheet_name", "Sheet1")
}

func read_or_create_config(viper_config *viper.Viper, config_path, config_file string) {
	err := viper_config.ReadInConfig() // Find and read the config file
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// Config file not found; ignore error if desired
		os.MkdirAll(config_path, 660)
		err = viper_config.WriteConfigAs(config_file)
		if err != nil { // Handle errors writing the config file
			panic(fmt.Errorf("fatal: error writing config file [config.read_or_create_config]: %w", err))
		}
	} else if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal: error reading config file [config.read_or_create_config]: %w", err))
	}
}

func set_config_globals(viper_config *viper.Viper) {
	DB_FILE = viper_config.GetString("db_path") + "/qc.sqlite3"
	LABEL_PATH = viper_config.GetString("label_path")
	JSON_PATHS = viper_config.GetStringSlice("json_paths")
	COA_TEMPLATE_PATH = viper_config.GetString("coa_template_path")
	COA_FILEPATH = viper_config.GetString("coa_filepath")
	LOG_FILE = viper_config.GetString("log_file")
	GUI.BASE_FONT_SIZE = viper_config.GetInt("font_size")
}

func set_config_globals_entry(viper_config *viper.Viper) {
	set_config_globals(viper_config)
	RETAIN_FILE_NAME = viper_config.GetString("retain_file_name")
	RETAIN_WORKSHEET_NAME = viper_config.GetString("retain_worksheet_name")
}

func set_config_globals_inbound(viper_config *viper.Viper) {
	set_config_globals(viper_config)
	PRODUCTION_SCHEDULE_FILE_NAME = viper_config.GetString("production_schedule_file_name")
	PRODUCTION_SCHEDULE_WORKSHEET_NAME = viper_config.GetString("production_schedule_worksheet_name")
}

func Load_config_viewer(appname string) *viper.Viper {
	viper_config, config_path, config_file := load_config_base(appname)
	set_config_defaults(appname, viper_config)
	read_or_create_config(viper_config, config_path, config_file)
	set_config_globals(viper_config)
	return viper_config
}

func Load_config_entry(appname string) *viper.Viper {
	viper_config, config_path, config_file := load_config_base(appname)
	set_config_defaults_entry(appname, viper_config)
	read_or_create_config(viper_config, config_path, config_file)
	set_config_globals_entry(viper_config)
	return viper_config
}

func Load_config_inbound(appname string) *viper.Viper {
	viper_config, config_path, config_file := load_config_base(appname)
	set_config_defaults_inbound(appname, viper_config)
	read_or_create_config(viper_config, config_path, config_file)
	set_config_globals_inbound(viper_config)
	return viper_config
}

func Write_config(viper_config *viper.Viper) {
	viper_config.WriteConfigAs(viper_config.ConfigFileUsed())
}
