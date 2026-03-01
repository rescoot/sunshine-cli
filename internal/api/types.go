package api

import (
	"strconv"
	"strings"
	"time"
)

type Scooter struct {
	ID                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	VIN                  string                 `json:"vin"`
	Color                string                 `json:"color"`
	ColorHex             string                 `json:"color_hex"`
	Online               bool                   `json:"online"`
	State                string                 `json:"state"`
	Speed                *float64               `json:"speed"`
	Odometer             *float64               `json:"odometer"`
	LastSeenAt           *time.Time             `json:"last_seen_at"`
	RadioGagaVersion     string                 `json:"radio_gaga_version"`
	LicensePlate         string                 `json:"license_plate"`
	Model                map[string]interface{} `json:"model"`
	AlarmState           string                 `json:"alarm_state"`
	AlarmStateHumanized  string                 `json:"alarm_state_humanized"`
	AlarmTriggered       bool                   `json:"alarm_triggered"`
	AlarmStateUpdatedAt  *time.Time             `json:"alarm_state_updated_at"`
	Location             *Location              `json:"location"`
	Batteries            *Batteries             `json:"batteries"`
	Telemetry            *Telemetry             `json:"telemetry"`
	Kickstand            string                 `json:"kickstand"`
	Seatbox              string                 `json:"seatbox"`
	Blinkers             string                 `json:"blinkers"`
	HandlebarStatus      map[string]interface{} `json:"handlebar_status"`
	ConnectivityStatus   map[string]interface{} `json:"connectivity_status"`
}

type Batteries struct {
	Battery0 *Battery `json:"battery0"`
	Battery1 *Battery `json:"battery1"`
	Aux      *Battery `json:"aux"`
	CBB      *Battery `json:"cbb"`
}

type Battery struct {
	Present      *bool   `json:"present"`
	Level        *int    `json:"level"`
	Voltage      *int    `json:"voltage"`
	State        string  `json:"state"`
	SOH          *int    `json:"soh"`
	CycleCount   *int    `json:"cycle_count"`
	SerialNumber string  `json:"serial_number"`
}

type Telemetry struct {
	Timestamp    *time.Time             `json:"timestamp"`
	VehicleState map[string]interface{} `json:"vehicle_state"`
	Engine       map[string]interface{} `json:"engine"`
	Battery0     map[string]interface{} `json:"battery0"`
	Battery1     map[string]interface{} `json:"battery1"`
	GPS          *GPS                   `json:"gps"`
	Connectivity map[string]interface{} `json:"connectivity"`
}

type GPS struct {
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
	Altitude *float64 `json:"altitude"`
	Speed    *float64 `json:"speed"`
	Course   *float64 `json:"course"`
}

type Trip struct {
	ID              int        `json:"id"`
	StartedAt       *time.Time `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at"`
	Distance        *float64   `json:"distance"`
	AvgSpeed        *float64   `json:"avg_speed"`
	DurationSeconds *int       `json:"duration_seconds"`
	StartLocation   *Location  `json:"start_location"`
	EndLocation     *Location  `json:"end_location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Destination struct {
	Latitude  *StringFloat `json:"latitude"`
	Longitude *StringFloat `json:"longitude"`
	Address   string       `json:"address"`
}

// StringFloat handles JSON values that may be either a number or a string-encoded number.
type StringFloat float64

func (sf *StringFloat) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "null" || s == "" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*sf = StringFloat(f)
	return nil
}

type CommandResponse struct {
	Status  string `json:"status"`
	Queued  bool   `json:"queued"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
