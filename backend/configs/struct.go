package configs

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Email    EmailConfig    `mapstructure:"email"`
}
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
	Sslmode  string `mapstructure:"sslmode"`
}
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
}
type EmailConfig struct {
	Smtp_host     string `mapstructure:"smtp_host"`
	Smtp_port     int    `mapstructure:"smtp_port"`
	Smtp_user     string `mapstructure:"smtp_user"`
	Smtp_password string `mapstructure:"smtp_password"`
}
