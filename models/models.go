package models

import (
	"github.com/google/uuid"
	"time"
)

type PVZ struct {
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

type Reception struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PVZID    uuid.UUID `json:"pvzId"`
	Status   string    `json:"status"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionID uuid.UUID `json:"receptionId"`
}
type Error struct {
	Message string `json:"message"`
}

type AddProductRequest struct {
	Type  string    `json:"type"`
	PVZID uuid.UUID `json:"pvzId"`
}

type PVZWithReceptions struct {
	PVZ        PVZ         `json:"pvz"`
	Receptions []Reception `json:"receptions"`
}

type PVZFilterParams struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}
