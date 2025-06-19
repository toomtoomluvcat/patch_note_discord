package connect

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Postgres() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Fail to load .env")
	}

	host := os.Getenv("Host")
	username := os.Getenv("Username")
	database := os.Getenv("Database")
	port := os.Getenv("Port")
	password := os.Getenv("Password")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		host, username, password, database, port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// if err := DB.AutoMigrate(&schema.Account{}, &schema.Calendar{}, &schema.PassKey{}); err != nil {
	// 	log.Fatal("AutoMigrate error:", err)
	// }

	log.Println("Connected to PostgreSQL and migrated schema.")
}
