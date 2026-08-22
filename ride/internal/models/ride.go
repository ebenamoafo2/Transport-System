package models

import "time"

type RideStatus string

const (
	StatusPending   RideStatus = "pending"
	StatusActive    RideStatus = "active"
	StatusCompleted RideStatus = "completed"
)

type Assignment struct {
	ID        string
	VehicleID string
	RouteID   string
	StartsAt  time.Time
	Status    RideStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
