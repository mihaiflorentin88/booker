package config

type Credentials struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	GoogleAPIKey string `json:"google_api_key"`
}

type StandBy struct {
	CanStandby bool   `json:"can_standby"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

type ParkingSpot struct {
	ParkingID string `json:"parking_id"`
	SpotID    string `json:"spot_id"`
}

type Schedule struct {
	Date   string `json:"date"`
	Offset int    `json:"offset"`
}

type Parking struct {
	StandBy     StandBy       `json:"standby"`
	Credentials Credentials   `json:"credentials"`
	ParkingSpot []ParkingSpot `json:"parkingSpot"`
	Schedule    Schedule      `json:"schedule"`
}

func NewEmptyParking() Parking {
	return Parking{
		StandBy:     StandBy{},
		Credentials: Credentials{},
		Schedule:    Schedule{},
		ParkingSpot: []ParkingSpot{},
	}
}
