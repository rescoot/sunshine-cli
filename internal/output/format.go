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
			dist = fmt.Sprintf("%.1f km", *t.Distance)
		}
		dur := "-"
		if t.DurationSeconds != nil {
			dur = formatDuration(*t.DurationSeconds)
		}
		speed := "-"
		if t.AvgSpeed != nil {
			speed = fmt.Sprintf("%.0f km/h", *t.AvgSpeed)
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", t.ID, date, dist, dur, speed)
	}
	w.Flush()
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
