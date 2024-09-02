package config

import (
	"strings"
)

type Config struct {
	Parking Parking `json:"parking"`
}

func NewEmptyConfig() *Config {
	return &Config{Parking: NewEmptyParking()}
}

func (c *Config) IsValid() (bool, string) {
	missingFields := ""
	if c.Parking.Schedule.Date == "" {
		missingFields = missingFields + "Parking.Schedule.Date, "
	}
	if c.Parking.Credentials.Username == "" {
		missingFields = missingFields + "Parking.Credentials.Username, "
	}
	if c.Parking.Credentials.Password == "" {
		missingFields = missingFields + "Parking.Credentials.Password, "
	}
	if c.Parking.Credentials.GoogleAPIKey == "" {
		missingFields = missingFields + "Parking.Credentials.GoogleAPIKey, "
	}
	if len(c.Parking.ParkingSpot) < 1 {
		missingFields = missingFields + "Parking.Spot, "
	}
	if c.Parking.StandBy.CanStandby {
		if c.Parking.StandBy.StartTime == "" {
			missingFields = missingFields + "Parking.StandBy.StartTime, "
		}
		if c.Parking.StandBy.EndTime == "" {
			missingFields = missingFields + "Parking.StandBy.EndTime, "
		}
	}
	if len(missingFields) > 0 {
		missingFields = strings.TrimSuffix(missingFields, ", ")
		return false, missingFields
	}
	return true, ""
}
