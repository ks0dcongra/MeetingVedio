package database

import (
	"fmt"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	UserName string = "postgres"
	Password string = "postgres"
	Addr     string = "127.0.0.1"
	Port     int    = 5432
	Database string = "postgres"
)

var dsn = fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
	UserName,
	Password,
	Addr,
	Port,
	Database,
)

var DB *gorm.DB

func DBinit() {
	for {
		var err error
		println(dsn)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		fmt.Println("Trying to connect database...")
		fmt.Println("DB Error===>", err)
		time.Sleep(3 * time.Second)
	}

	fmt.Println("Database connected!")
}
