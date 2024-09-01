package dto

import (
	"encoding/json"
)

type BookParkingPayload struct {
	Me           string   `json:"me"`
	ParkingID    string   `json:"parkingId"`
	UID          string   `json:"uid"`
	SpotID       string   `json:"spotId"`
	Day          int      `json:"day"`
	AddShifts    []string `json:"addShifts"`
	RemoveShifts []string `json:"removeShifts"`
	Version      string   `json:"v"`
}

func NewBookingPayload(daysSinceEpoch int, spotID, parkingID string) *BookParkingPayload {
	const (
		me      = "sTFF1cJSeyZObYPvamECYqwB0OE2"
		uid     = "sTFF1cJSeyZObYPvamECYqwB0OE2"
		version = "xKpQV"
	)

	return &BookParkingPayload{
		Me:           me,
		ParkingID:    parkingID,
		UID:          uid,
		SpotID:       spotID,
		Day:          daysSinceEpoch,
		AddShifts:    []string{"08001700"},
		RemoveShifts: []string{},
		Version:      version,
	}
}

func (p *BookParkingPayload) ToJson() ([]byte, error) {
	jsonData, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
