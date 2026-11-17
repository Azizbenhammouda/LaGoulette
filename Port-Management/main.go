package main

import (
	"PortManagement/vessel"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("POSTGRES_CONFIG")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error; err != nil {
		log.Fatal("failed to enable pgcrypto:", err)
	}
	if err := db.AutoMigrate(&vessel.Vessel{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}
	vesselRepo := vessel.NewVesselRepository(db)
	vesselService := vessel.NewVesselService(vesselRepo)
	vesselHandler := vessel.NewVesselHandler(vesselService)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/vessels", vesselHandler.Register)
	mux.HandleFunc("GET /api/vessels", vesselHandler.List)
	mux.HandleFunc("GET /api/vessels/{id}", vesselHandler.GetVessel)
	mux.HandleFunc("PATCH /api/vessels/{id}", vesselHandler.Update)
	mux.HandleFunc("DELETE /api/vessels/{id}", vesselHandler.Delete)
}
