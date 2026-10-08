package work

import (
	"context"
	"encoding/json"
	"math"
)

func (s *Service) Contracts(ctx context.Context, limit, offset int) (json.RawMessage, error) {
	return s.Repo.List(ctx, s.Repo.Pool, query11, limit, offset, s.Now().In(Bangkok).Format("2006-01-02"))
}
func (s *Service) Contract(ctx context.Context, id string) (json.RawMessage, error) {
	if e := validID(id, "tor_id"); e != nil {
		return nil, e
	}
	return one(ctx, s.Repo.Pool, query12, id)
}
func (s *Service) ContractState(ctx context.Context, id, state string) error {
	if e := validID(id, "tor_id"); e != nil {
		return e
	}
	switch state {
	case "registered", "active", "complete", "cancelled":
	default:
		return invalid("status", "invalid contract status")
	}
	tag, e := s.Repo.Exec(ctx, s.Repo.Pool, query13, id, state)
	if e == nil && tag.RowsAffected() != 1 {
		return conflict("contract not found")
	}
	return e
}

type LocationInput struct {
	Name      string   `json:"location_name"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func (s *Service) Location(ctx context.Context, id string, in LocationInput) (string, error) {
	if e := text(in.Name, "location_name", 255); e != nil {
		return "", e
	}
	if (in.Latitude == nil) != (in.Longitude == nil) {
		return "", invalid("coordinates", "latitude and longitude must be provided together")
	}
	if in.Latitude != nil {
		a, b := *in.Latitude, *in.Longitude
		if math.IsNaN(a) || math.IsNaN(b) || math.IsInf(a, 0) || math.IsInf(b, 0) || a < -90 || a > 90 || b < -180 || b > 180 {
			return "", invalid("coordinates", "invalid coordinates")
		}
	}
	if id != "" {
		if e := validID(id, "location_id"); e != nil {
			return "", e
		}
	}
	e := s.Repo.Transaction(ctx, func(q Query) error {
		if e := lock(ctx, q, "catalog:location"); e != nil {
			return e
		}
		var duplicate bool
		if e := s.Repo.Row(ctx, q, query14, in.Name, nullableID(id)).Scan(&duplicate); e != nil {
			return e
		}
		if duplicate {
			return conflict("location name already exists")
		}
		if id == "" {
			return s.Repo.Row(ctx, q, query15, in.Name, in.Address, in.Latitude, in.Longitude).Scan(&id)
		}
		tag, e := s.Repo.Exec(ctx, q, query16, id, in.Name, in.Address, in.Latitude, in.Longitude)
		if e == nil && tag.RowsAffected() != 1 {
			return conflict("location not found")
		}
		return e
	})
	return id, e
}
func nullableID(id string) any {
	if id == "" {
		return nil
	}
	return id
}

type EquipmentInput struct {
	Name   string `json:"equipment_name"`
	Active *bool  `json:"is_active"`
}

func (s *Service) Equipment(ctx context.Context, id string, in EquipmentInput) (string, error) {
	if e := text(in.Name, "equipment_name", 255); e != nil {
		return "", e
	}
	if id != "" {
		if e := validID(id, "equipment_id"); e != nil {
			return "", e
		}
	}
	e := s.Repo.Transaction(ctx, func(q Query) error {
		if e := lock(ctx, q, "catalog:equipment"); e != nil {
			return e
		}
		var duplicate bool
		if e := s.Repo.Row(ctx, q, query17, in.Name, nullableID(id)).Scan(&duplicate); e != nil {
			return e
		}
		if duplicate {
			return conflict("equipment name already exists")
		}
		if id == "" {
			return s.Repo.Row(ctx, q, query18, in.Name, in.Active).Scan(&id)
		}
		tag, e := s.Repo.Exec(ctx, q, query19, id, in.Name, in.Active)
		if e == nil && tag.RowsAffected() != 1 {
			return conflict("equipment not found")
		}
		return e
	})
	return id, e
}
func (s *Service) Equipments(ctx context.Context, limit, offset int) (json.RawMessage, error) {
	return s.Repo.List(ctx, s.Repo.Pool, query20, limit, offset)
}
