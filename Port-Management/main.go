package main

import (
	"PortManagement/berth"
	"PortManagement/vessel"
	"fmt"
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
	if err := db.AutoMigrate(&berth.Berth{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}
	vesselRepo := vessel.NewVesselRepository(db)
	vesselService := vessel.NewVesselService(vesselRepo)
	vesselHandler := vessel.NewVesselHandler(vesselService)
	berthRepo := berth.NewBerthRepository(db)
	berthService := berth.NewBerthService(berthRepo)
	berthHandler := berth.NewBerthHandler(berthService)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/vessels", vesselHandler.Register)
	mux.HandleFunc("GET /api/vessels", vesselHandler.List)
	mux.HandleFunc("GET /api/vessels/{id}", vesselHandler.GetVessel)
	mux.HandleFunc("PATCH /api/vessels/{id}", vesselHandler.Update)
	mux.HandleFunc("DELETE /api/vessels/{id}", vesselHandler.Delete)
	mux.HandleFunc("GET /api/berths", berthHandler.List)
	mux.HandleFunc("GET /api/berths/available", berthHandler.ListAvailable)
	mux.HandleFunc("POST /api/berth-reservations", berthHandler.ReserveBerth)
	mux.HandleFunc("GET /api/berth-reservations/{id}", berthHandler.GetReserve)
	fmt.Println("Server Running..")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
