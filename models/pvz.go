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

type PVZWithReceptions struct {
	PVZ        PVZ         `json:"pvz"`
	Receptions []Reception `json:"receptions"`
	Products   []Product   `json:"products"`
}

type PVZFilterParams struct {
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}
