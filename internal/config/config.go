package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	Env            string `mapstructure:"env"`
	Name           string `mapstructure:"name"`
	Port           string `mapstructure:"port"`
	DB             DB     `mapstructure:"database"`
	BcryptSalt     int    `mapstructure:"bcrypt_salt"`
	JWTSecretKey   string `mapstructure:"jwt_secret_key"`
	AdminSecretKey string `mapstructure:"admin_secret_key"`
}

type DB struct {
	MaxIdleCons    int  `mapstructure:"maxIdleCons"`
	MaxOpenCons    int  `mapstructure:"maxOpenCons"`
	ConMaxIdleTime int  `mapstructure:"conMaxIdleTime"`
	ConMaxLifetime int  `mapstructure:"conMaxLifeTime"`
	Replica        PSQL `mapstructure:"replica"`
	Master         PSQL `mapstructure:"master"`
}

type PSQL struct {
	DBName   string `mapstructure:"dbName"`
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Schema   string `mapstructure:"schema"`
	User     string `mapstructure:"user"`
	Debug    bool   `mapstructure:"debug"`
}

func (p PSQL) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=%s",
		p.Host, strconv.Itoa(p.Port), p.User, p.Password, p.DBName, p.Schema, "disable")

	//return fmt.Sprintf("postgres://%s:%s@%s/$%s?sslmode=verify-full",
	//p.User, p.Password, p.Host, p.DBName)
}

func Load() *Config {
	v := viper.New()

	v.SetConfigType("yaml")
	v.AddConfigPath("./config/")
	v.SetConfigName("config.local")

	if err := v.ReadInConfig(); err != nil {
		log.Panic("error read config")
	}

	var conf Config
	if err := v.Unmarshal(&conf); err != nil {
		log.Panic("error read config")
	}

	return &conf
}

func InitializeDB(conf *DB) (master *sql.DB, replica *sql.DB) {
	var err error
	ctx := context.Background()
	master, err = sql.Open("postgres", conf.Master.ConnectionString())
	if err != nil {
		log.Fatal(ctx, "Can't connect to master DB %+v", err)
	}

	log.Println("Successfully connect to master DB")

	master.SetMaxIdleConns(conf.MaxIdleCons)
	master.SetMaxOpenConns(conf.MaxOpenCons)
	master.SetConnMaxLifetime(time.Duration(conf.ConMaxLifetime) * time.Millisecond)
	master.SetConnMaxIdleTime(time.Duration(conf.ConMaxIdleTime) * time.Millisecond)

	replica, err = sql.Open("postgres", conf.Replica.ConnectionString())
	if err != nil {
		log.Fatal(ctx, "Can't connect to replica DB %+v", err)
	}

	replica.SetMaxIdleConns(conf.MaxIdleCons)
	replica.SetMaxOpenConns(conf.MaxOpenCons)
	replica.SetConnMaxLifetime(time.Duration(conf.ConMaxLifetime) * time.Millisecond)
	replica.SetConnMaxIdleTime(time.Duration(conf.ConMaxIdleTime) * time.Millisecond)

	log.Println("Successfully connect to replica DB")

	return
}
