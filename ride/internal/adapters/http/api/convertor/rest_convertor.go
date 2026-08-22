package convertor

import (
	"github.com/ebenamoafo2/transport/ride/internal/adapters/http/api"
	"github.com/ebenamoafo2/transport/ride/internal/models"
)

// AssignmentToDomain converts an api.Assignment to a models.Assignment.
func AssignmentToDomain(r api.Assignment) models.Assignment {
	return models.Assignment{
		ID:        *r.Metadata.Id,
		VehicleID: r.VehicleId,
		RouteID:   r.RouteId,
		StartsAt:  r.StartsAt,
		Status:    models.RideStatus(r.Status),
	}
}

// AssignmentFromDomain converts a models.Assignment to an api.Assignment.
func AssignmentFromDomain(r models.Assignment) api.Assignment {
	return api.Assignment{
		Metadata: &api.EntityMetadata{
			Id: &r.ID,
		},
		VehicleId: r.VehicleID,
		RouteId:   r.RouteID,
		StartsAt:  r.StartsAt,
		Status:    api.AssignmentStatus(r.Status),
	}
}
