package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
    DB *gorm.DB
}

func ConnectDB() *gorm.DB {
    var db *gorm.DB
    var err error

    err = godotenv.Load()
    if err != nil {
        log.Printf("Failed to load the database url, %v", err)
        return nil
    }

    dbs := os.Getenv("DATABASE_URL")
    db, err = gorm.Open( postgres.New(postgres.Config{
          DSN: dbs,
          PreferSimpleProtocol: true,
      }), &gorm.Config{
          PrepareStmt: true,
          SkipDefaultTransaction: true,
      })

    if err != nil {
        log.Printf("failed to connect, %v", err)
        return nil
    }

    // migrate(db)

    return db
}
