package models

import (
	"github.com/google/uuid"
	"time"
)

type AddProductRequest struct {
	Type  string    `json:"type"`
	PVZID uuid.UUID `json:"pvzId"`
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
