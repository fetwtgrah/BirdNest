package configs

import "github.com/spf13/viper"

var Conf *Config

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	Conf = &Config{}
	if err := viper.Unmarshal(Conf); err != nil {
		return err
	}
	return nil
}
