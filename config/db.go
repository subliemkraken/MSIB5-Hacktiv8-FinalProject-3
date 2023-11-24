package config

import (
	"FinalProject3/model"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var User = model.User{}
var Category = model.Category{}
var Task = model.Task{}
var err error

func seedAdmin() {
	admin := &model.User{
		Full_name: "admin",
		Email:     "admin@email.com",
		Password:  "admin123",
		Role:      "admin",
	}

	admin.HashPassword()

	DB.Create(admin)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	log.Println("Admin's account has been successfully seeded.")

}

func StartDB() {
	godotenv.Load(".env")
	dsn := "host=" + os.Getenv("PGHOST") + " user=" + os.Getenv("PGUSER") + " password=" + os.Getenv("PGPASSWORD") + " dbname=" + os.Getenv("PGDATABASE") + " port=" + os.Getenv("PGPORT") + " sslmode=disable TimeZone=Asia/Jakarta"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&Category, &User, &Task)

	if DB.Migrator().HasTable(&User) {
		var userCount int64
		if err := DB.Model(&User).Count(&userCount).Error; err != nil {
			log.Fatalf("Error checking user's table: %v", err)
		}

		if userCount == 0 {
			seedAdmin()
			if err != nil {
				log.Fatalf("Failed to seed admin: %v", err)
			}
		}
	}

}
