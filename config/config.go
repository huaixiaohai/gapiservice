package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

const (
	Develop    = "develop"
	Release    = "release"
	Production = "production"
)

var (
	file string
	C    *Config
)

type Config struct {
	RunEnv         string
	ServerName     string
	Log            *Log
	Http           *Http
	Gorm           *Gorm
	MySQL          *MySQL
	Cron           *Cron
	CookieUserUUID string
}

type Cron struct {
	GetLuckListJob   string
	RefreshCookieJob string
}

type Log struct {
	Level    uint32
	Output   string
	FilePath string
}

type Http struct {
	Addr     string
	CertFile string
	KeyFile  string
}

// Gorm gorm配置参数
type Gorm struct {
	Debug             bool
	DBType            string
	MaxLifetime       int
	MaxOpenConns      int
	MaxIdleConns      int
	EnableAutoMigrate bool
}

// MySQL mysql配置参数
type MySQL struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	Parameters string
}

// DSN 数据库连接串
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

func init() {
	flag.StringVar(&file, "c", "./config/config.yaml", "config file path")
	flag.Parse()
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		println(err)
		panic(err)
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		println(err)
		panic(err)
	}
}
