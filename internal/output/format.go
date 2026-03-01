package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/rescoot/sunshine-cli/internal/api"
)

var JSONOutput bool

func PrintJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func PrintScooterList(scooters []api.Scooter) {
	if JSONOutput {
		PrintJSON(scooters)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tVIN\tSTATE\tBATTERY\tONLINE")
	for _, s := range scooters {
		battery := batteryDisplay(s.Batteries)
		online := "no"
		if s.Online {
			online = "yes"
		}
		state := s.State
		if state == "" {
			state = "-"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			s.ID, s.Name, s.VIN, state, battery, online)
	}
	w.Flush()
}

func PrintScooterDetail(s *api.Scooter) {
	if JSONOutput {
		PrintJSON(s)
		return
	}

	kv := [][]string{
		{"Name", s.Name},
		{"VIN", s.VIN},
		{"Model", modelName(s.Model)},
		{"Color", fmt.Sprintf("%s (%s)", s.Color, s.ColorHex)},
		{"State", nonEmpty(s.State)},
		{"Online", boolStr(s.Online)},
	}

	if s.Batteries != nil {
		if s.Batteries.Battery0 != nil && s.Batteries.Battery0.Level != nil {
			kv = append(kv, []string{"Battery 0", fmt.Sprintf("%d%% (%s)", *s.Batteries.Battery0.Level, nonEmpty(s.Batteries.Battery0.State))})
		}
		if s.Batteries.Battery1 != nil && s.Batteries.Battery1.Level != nil {
			kv = append(kv, []string{"Battery 1", fmt.Sprintf("%d%% (%s)", *s.Batteries.Battery1.Level, nonEmpty(s.Batteries.Battery1.State))})
		}
	}

	if s.Telemetry != nil && s.Telemetry.GPS != nil {
		gps := s.Telemetry.GPS
		if gps.Lat != nil && gps.Lng != nil {
			kv = append(kv, []string{"Location", fmt.Sprintf("%.4f, %.4f", *gps.Lat, *gps.Lng)})
		}
	}

	if s.Speed != nil {
		kv = append(kv, []string{"Speed", fmt.Sprintf("%.0f km/h", *s.Speed)})
	}
	if s.Odometer != nil {
		kv = append(kv, []string{"Odometer", fmt.Sprintf("%.1f km", *s.Odometer/1000)})
	}

	kv = append(kv, []string{"Alarm", s.AlarmStateHumanized})

	if s.LastSeenAt != nil {
		kv = append(kv, []string{"Last Seen", timeAgo(*s.LastSeenAt)})
	}

	if s.RadioGagaVersion != "" {
		kv = append(kv, []string{"Radio Gaga", s.RadioGagaVersion})
	}

	printKeyValue(kv)
}

func PrintScooterStatus(s *api.Scooter) {
	if JSONOutput {
		PrintJSON(s)
		return
	}

	// Header
	fmt.Printf("%s (%s)\n", s.Name, s.VIN)
	if m := modelName(s.Model); m != "-" {
		fmt.Printf("%s\n", m)
	}
	fmt.Println()

	// Core status
	kv := [][]string{
		{"State", nonEmpty(s.State)},
		{"Online", boolStr(s.Online)},
	}
	if s.LastSeenAt != nil {
		kv = append(kv, []string{"Last Seen", timeAgo(*s.LastSeenAt)})
	}
	if s.LicensePlate != "" {
		kv = append(kv, []string{"Plate", s.LicensePlate})
	}
	kv = append(kv, []string{"Alarm", s.AlarmStateHumanized})
	if s.RadioGagaVersion != "" {
		kv = append(kv, []string{"Radio Gaga", s.RadioGagaVersion})
	}
	printKeyValue(kv)

	t := s.Telemetry
	if t == nil {
		return
	}

	// Vehicle state
	if vs := t.VehicleState; len(vs) > 0 {
		printSection("Vehicle", mapToKV(vs,
			"state", "kickstand", "seatbox", "blinkers",
			"handlebar_lock", "handlebar_position",
			"brake_left", "brake_right",
			"main_power", "seatbox_button", "horn_button", "blinker_switch",
		))
	}

	// Engine
	if eng := t.Engine; len(eng) > 0 {
		printSection("Engine", mapToKVFormatted(eng, map[string]func(interface{}) string{
			"speed":       func(v interface{}) string { return fmtInt(v) + " km/h" },
			"odometer":    func(v interface{}) string { return fmtFloat(v, 1000) + " km" },
			"temperature": func(v interface{}) string { return fmtInt(v) + " °C" },
			"motor_voltage": func(v interface{}) string {
				if f, ok := toFloat(v); ok {
					return fmt.Sprintf("%.1f V", f/1000)
				}
				return fmtVal(v)
			},
			"motor_current": func(v interface{}) string {
				if f, ok := toFloat(v); ok {
					return fmt.Sprintf("%.1f A", f/1000)
				}
				return fmtVal(v)
			},
		}, "speed", "odometer", "temperature", "state",
			"motor_voltage", "motor_current", "motor_rpm",
			"kers_state", "kers_reason_off", "throttle_state", "fw_version"))
	}

	// Batteries
	printBatterySection("Battery 0", t.Battery0)
	printBatterySection("Battery 1", t.Battery1)
	if aux := t.AuxBattery; len(aux) > 0 {
		printSection("Aux Battery", mapToKVFormatted(aux, map[string]func(interface{}) string{
			"level": func(v interface{}) string { return fmtInt(v) + "%" },
			"voltage": func(v interface{}) string {
				if f, ok := toFloat(v); ok {
					return fmt.Sprintf("%.1f V", f/1000)
				}
				return fmtVal(v)
			},
		}, "level", "voltage", "charge_status"))
	}
	if cbb := t.CBBBattery; len(cbb) > 0 {
		// CBB fuel gauge (MAX17301) reports in µA/µWh
		printSection("CBB Battery", mapToKVFormatted(cbb, map[string]func(interface{}) string{
			"level": func(v interface{}) string { return fmtInt(v) + "%" },
			"current": func(v interface{}) string {
				if f, ok := toFloat(v); ok {
					return fmt.Sprintf("%.1f mA", f/1000)
				}
				return fmtVal(v)
			},
			"temperature": func(v interface{}) string { return fmtInt(v) + " °C" },
			"soh": func(v interface{}) string { return fmtInt(v) + "%" },
			"full_capacity": func(v interface{}) string {
				if f, ok := toFloat(v); ok {
					return fmt.Sprintf("%.1f Wh", f/1000000)
				}
				return fmtVal(v)
			},
			"remaining_capacity": func(v interface{}) string {
				if f, ok := toFloat(v); ok {
					return fmt.Sprintf("%.1f Wh", f/1000000)
				}
				return fmtVal(v)
			},
			"time_to_empty": func(v interface{}) string {
				if f, ok := toFloat(v); ok {
					return formatDuration(int(f))
				}
				return fmtVal(v)
			},
			"time_to_full": func(v interface{}) string {
				if f, ok := toFloat(v); ok {
					return formatDuration(int(f))
				}
				return fmtVal(v)
			},
		},
			"level", "current", "temperature", "soh",
			"charge_status", "cycle_count",
			"full_capacity", "remaining_capacity",
			"time_to_empty", "time_to_full",
			"serial_number", "part_number", "unique_id",
		))
	}

	// GPS
	if gps := t.GPS; gps != nil {
		var gpsKV [][]string
		if gps.Lat != nil && gps.Lng != nil {
			gpsKV = append(gpsKV, []string{"Position", fmt.Sprintf("%.6f, %.6f", *gps.Lat, *gps.Lng)})
		}
		if gps.Altitude != nil {
			gpsKV = append(gpsKV, []string{"Altitude", fmt.Sprintf("%.0f m", *gps.Altitude)})
		}
		if gps.Speed != nil {
			gpsKV = append(gpsKV, []string{"GPS Speed", fmt.Sprintf("%.0f km/h", *gps.Speed)})
		}
		if gps.Course != nil {
			gpsKV = append(gpsKV, []string{"Course", fmt.Sprintf("%.0f°", *gps.Course)})
		}
		if len(gpsKV) > 0 {
			printSection("GPS", gpsKV)
		}
	}

	// Connectivity
	if conn := t.Connectivity; len(conn) > 0 {
		printSection("Connectivity", mapToKV(conn,
			"modem_state", "internet_status", "cloud_status",
			"access_tech", "signal_quality", "ip_address", "imei"))
	}

	// System
	if sys := t.System; len(sys) > 0 {
		printSection("System", mapToKV(sys,
			"mdb_version", "dbc_version", "nrf_fw_version",
			"engine_fw_version", "environment",
			"mdb_sn", "dbc_sn", "mdb_flavor", "dbc_flavor"))
	}

	// Power
	if pwr := t.Power; len(pwr) > 0 {
		printSection("Power", mapToKV(pwr,
			"state", "mux_input", "wakeup_source",
			"hibernate_level", "nrf_reset_count", "nrf_reset_reason"))
	}
}

func printSection(title string, kv [][]string) {
	if len(kv) == 0 {
		return
	}
	fmt.Printf("\n%s\n", title)
	printKeyValue(kv)
}

// mapToKV extracts keys from a map in the given order, skipping nil/empty values.
func mapToKV(m map[string]interface{}, keys ...string) [][]string {
	return mapToKVFormatted(m, nil, keys...)
}

func mapToKVFormatted(m map[string]interface{}, formatters map[string]func(interface{}) string, keys ...string) [][]string {
	var kv [][]string
	for _, k := range keys {
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}
		var display string
		if fn, hasFmt := formatters[k]; hasFmt {
			display = fn(v)
		} else {
			display = fmtVal(v)
		}
		if display == "" || display == "0" {
			continue
		}
		label := humanizeKey(k)
		kv = append(kv, []string{label, display})
	}
	return kv
}

func printBatterySection(title string, m map[string]interface{}) {
	if len(m) == 0 {
		return
	}
	// Skip if not present
	if p, ok := m["present"]; ok {
		if present, isBool := p.(bool); isBool && !present {
			return
		}
	}
	printSection(title, mapToKVFormatted(m, map[string]func(interface{}) string{
		"level": func(v interface{}) string { return fmtInt(v) + "%" },
		"voltage": func(v interface{}) string {
			if f, ok := toFloat(v); ok {
				return fmt.Sprintf("%.1f V", f/1000)
			}
			return fmtVal(v)
		},
		"current": func(v interface{}) string {
			if f, ok := toFloat(v); ok {
				return fmt.Sprintf("%.2f A", f/1000)
			}
			return fmtVal(v)
		},
		"temps": func(v interface{}) string {
			if arr, ok := v.([]interface{}); ok {
				parts := make([]string, len(arr))
				for i, t := range arr {
					parts[i] = fmt.Sprintf("%.0f", toFloatOr(t, 0))
				}
				return strings.Join(parts, "/") + " °C"
			}
			return fmtVal(v)
		},
		"soh": func(v interface{}) string { return fmtInt(v) + "%" },
	},
		"level", "state", "voltage", "current",
		"temp_state", "temps", "soh", "cycle_count",
		"serial_number", "manufacturing_date", "fw_version"))
}

func humanizeKey(key string) string {
	key = strings.ReplaceAll(key, "_", " ")
	// Title case
	words := strings.Fields(key)
	for i, w := range words {
		if len(w) <= 3 && i > 0 {
			words[i] = strings.ToUpper(w) // acronyms: soh, fw, ip, sn
		} else {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func fmtVal(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case bool:
		return boolStr(val)
	case float64:
		if val == float64(int(val)) {
			return fmt.Sprintf("%d", int(val))
		}
		return fmt.Sprintf("%.2f", val)
	case []interface{}:
		parts := make([]string, len(val))
		for i, item := range val {
			parts[i] = fmtVal(item)
		}
		return strings.Join(parts, ", ")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func fmtInt(v interface{}) string {
	if f, ok := toFloat(v); ok {
		return fmt.Sprintf("%d", int(f))
	}
	return fmtVal(v)
}

func fmtFloat(v interface{}, divisor float64) string {
	if f, ok := toFloat(v); ok {
		return fmt.Sprintf("%.1f", f/divisor)
	}
	return fmtVal(v)
}

func toFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	default:
		return 0, false
	}
}

func toFloatOr(v interface{}, fallback float64) float64 {
	if f, ok := toFloat(v); ok {
		return f
	}
	return fallback
}

func PrintCommandResponse(resp *api.CommandResponse, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if JSONOutput {
		PrintJSON(resp)
		return
	}

	if resp.Status == "success" {
		fmt.Println("OK")
	} else {
		msg := resp.Message
		if msg == "" {
			msg = resp.Error
		}
		fmt.Fprintf(os.Stderr, "Failed: %s\n", msg)
		os.Exit(1)
	}
}

func PrintTripList(trips []api.Trip) {
	if JSONOutput {
		PrintJSON(trips)
		return
	}

	if len(trips) == 0 {
		fmt.Println("No trips found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tDATE\tDISTANCE\tDURATION\tAVG SPEED")
	for _, t := range trips {
		date := "-"
		if t.StartedAt != nil {
			date = t.StartedAt.Local().Format("2006-01-02 15:04")
		}
		dist := "-"
		if t.Distance != nil {
			dist = fmt.Sprintf("%.1f km", float64(*t.Distance)/1000)
		}
		dur := "-"
		if t.Duration != nil {
			dur = formatDuration(*t.Duration)
		}
		speed := "-"
		if t.AvgSpeed != nil {
			speed = fmt.Sprintf("%.0f km/h", float64(*t.AvgSpeed))
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", t.ID, date, dist, dur, speed)
	}
	w.Flush()
}

func PrintTripDetail(t *api.Trip) {
	if JSONOutput {
		PrintJSON(t)
		return
	}

	kv := [][]string{
		{"ID", fmt.Sprintf("%d", t.ID)},
	}
	if t.StartedAt != nil {
		kv = append(kv, []string{"Started", t.StartedAt.Local().Format("2006-01-02 15:04:05")})
	}
	if t.EndedAt != nil {
		kv = append(kv, []string{"Ended", t.EndedAt.Local().Format("2006-01-02 15:04:05")})
	}
	if t.Duration != nil {
		kv = append(kv, []string{"Duration", formatDuration(*t.Duration)})
	}
	if t.Distance != nil {
		kv = append(kv, []string{"Distance", fmt.Sprintf("%.1f km", float64(*t.Distance)/1000)})
	}
	if t.AvgSpeed != nil {
		kv = append(kv, []string{"Avg Speed", fmt.Sprintf("%.0f km/h", float64(*t.AvgSpeed))})
	}
	if t.StartLocation != nil {
		kv = append(kv, []string{"Start", fmt.Sprintf("%.4f, %.4f", t.StartLocation.Lat, t.StartLocation.Lng)})
	}
	if t.EndLocation != nil {
		kv = append(kv, []string{"End", fmt.Sprintf("%.4f, %.4f", t.EndLocation.Lat, t.EndLocation.Lng)})
	}
	printKeyValue(kv)
}

func PrintDestination(dest *api.Destination) {
	if JSONOutput {
		PrintJSON(dest)
		return
	}

	if dest.Latitude == nil || dest.Longitude == nil {
		fmt.Println("No destination set.")
		return
	}

	kv := [][]string{
		{"Latitude", fmt.Sprintf("%.6f", *dest.Latitude)},
		{"Longitude", fmt.Sprintf("%.6f", *dest.Longitude)},
	}
	if dest.Address != "" {
		kv = append(kv, []string{"Address", dest.Address})
	}
	printKeyValue(kv)
}

// Helpers

func printKeyValue(pairs [][]string) {
	maxKey := 0
	for _, p := range pairs {
		if len(p[0]) > maxKey {
			maxKey = len(p[0])
		}
	}
	for _, p := range pairs {
		fmt.Printf("%-*s  %s\n", maxKey+1, p[0]+":", p[1])
	}
}

func batteryDisplay(b *api.Batteries) string {
	if b == nil {
		return "-"
	}
	parts := []string{}
	if b.Battery0 != nil && b.Battery0.Level != nil {
		parts = append(parts, fmt.Sprintf("%d%%", *b.Battery0.Level))
	}
	if b.Battery1 != nil && b.Battery1.Level != nil {
		parts = append(parts, fmt.Sprintf("%d%%", *b.Battery1.Level))
	}
	if len(parts) == 0 {
		return "-"
	}
	return strings.Join(parts, "/")
}

func boolStr(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func modelName(m map[string]interface{}) string {
	if m == nil {
		return "-"
	}
	if name, ok := m["full_name"].(string); ok {
		return name
	}
	if name, ok := m["model_name"].(string); ok {
		return name
	}
	return "-"
}

func nonEmpty(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	}
}

func formatDuration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
