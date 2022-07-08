package database

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

const (
	dsnFormat = "%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

type DBEnv struct {
	User     string
	Password string
	Host     string
	Name     string
}

func SetupDB() error {
	var env DBEnv
	err := envconfig.Process("db", &env)
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf(dsnFormat, env.User, env.Password, env.Host, env.Name)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

func GetDBCli() *gorm.DB {
	return db
}
