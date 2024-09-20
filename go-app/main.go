// go-app/main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SensorData struct {
	ID        uint      `gorm:"primaryKey"`
	Timestamp time.Time `gorm:"index"`
	Value     float64
}

func main() {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Échec de la connexion à la base de données : %v", err)
	}

	// Migration du modèle
	err = db.AutoMigrate(&SensorData{})
	if err != nil {
		log.Fatalf("Échec de la migration : %v", err)
	}

	// Création de l'hypertable TimescaleDB
	db.Exec("SELECT create_hypertable('sensor_data', 'timestamp', if_not_exists => TRUE);")

	// Routes HTTP pour les opérations CRUD
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		// Création d'une nouvelle entrée
		data := SensorData{
			Timestamp: time.Now(),
			Value:     42.0,
		}
		result := db.Create(&data)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Données créées avec l'ID : %d\n", data.ID)
	})

	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		// Lecture des entrées
		var data []SensorData
		result := db.Find(&data)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		for _, d := range data {
			fmt.Fprintf(w, "ID: %d, Timestamp: %s, Value: %f\n", d.ID, d.Timestamp, d.Value)
		}
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		// Mise à jour d'une entrée
		var data SensorData
		result := db.First(&data)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		data.Value = 84.0
		db.Save(&data)
		fmt.Fprintf(w, "Données mises à jour pour l'ID : %d\n", data.ID)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		// Suppression d'une entrée
		var data SensorData
		result := db.First(&data)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		db.Delete(&data)
		fmt.Fprintf(w, "Données supprimées pour l'ID : %d\n", data.ID)
	})

	log.Println("Serveur démarré sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
