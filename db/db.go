package db

import (
	"customer-api/models"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *sql.DB
var GormDB *gorm.DB

func InitDB() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve environment variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	// Validate required environment variables
	if user == "" || password == "" || dbname == "" || host == "" || port == "" {
		log.Fatal("Missing required database environment variables")
	}

	// Construct the connection string
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	GormDB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening gorm database: %v", err)
	}

	if err = runMigrations(GormDB); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	log.Println("Successfully connected to the database!")
}

func runMigrations(gdb *gorm.DB) error {
	m := gormigrate.New(gdb, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202605161920_create_customers_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(`
					CREATE TABLE IF NOT EXISTS customers (
						id SERIAL PRIMARY KEY,
						first_name VARCHAR(100) NOT NULL,
						last_name VARCHAR(100) NOT NULL,
						email VARCHAR(255) NOT NULL UNIQUE
					);
				`).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Customer{})
			},
		},
		{
			ID: "202605161921_seed_customers",
			Migrate: func(tx *gorm.DB) error {
				customers := []models.Customer{
					{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
					{FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"},
					{FirstName: "Alice", LastName: "Johnson", Email: "alice.johnson@example.com"},
					{FirstName: "Bob", LastName: "Williams", Email: "bob.williams@example.com"},
					{FirstName: "Carol", LastName: "Brown", Email: "carol.brown@example.com"},
				}

				for _, customer := range customers {
					if err := tx.Where(models.Customer{Email: customer.Email}).FirstOrCreate(&customer).Error; err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				emails := []string{
					"john.doe@example.com",
					"jane.smith@example.com",
					"alice.johnson@example.com",
					"bob.williams@example.com",
					"carol.brown@example.com",
				}
				return tx.Where("email IN ?", emails).Delete(&models.Customer{}).Error
			},
		},
	})

	return m.Migrate()
}
