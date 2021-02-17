package dynamic

type (

	// General is a VATSIM data header
	General struct {
		Version          int    `json:"version"`
		Reload           int    `json:"reload"`
		Update           string `json:"update"`
		UpdateTimestamp  string `json:"update_timestamp"`
		ConnectedClients int    `json:"connected_clients"`
		UniqueUsers      int    `json:"unique_users"`
	}

	// FlightPlan is a VATSIM flight plan
	FlightPlan struct {
		FlightRules string `json:"flight_rules"`
		Aircraft    string `json:"aircraft"`
		Departure   string `json:"departure"`
		Arrival     string `json:"arrival"`
		Alternate   string `json:"alternate"`
		CruiseTas   string `json:"cruise_tas"`
		Altitude    string `json:"altitude"`
		Deptime     string `json:"deptime"`
		EnrouteTime string `json:"enroute_time"`
		FuelTime    string `json:"fuel_time"`
		Remarks     string `json:"remarks"`
		Route       string `json:"route"`
	}

	// Pilot is a VATSIM pilot
	Pilot struct {
		Cid         int         `json:"cid"`
		Name        string      `json:"name"`
		Callsign    string      `json:"callsign"`
		Server      string      `json:"server"`
		PilotRating int         `json:"pilot_rating"`
		Latitude    float64     `json:"latitude"`
		Longitude   float64     `json:"longitude"`
		Altitude    int         `json:"altitude"`
		Groundspeed int         `json:"groundspeed"`
		Transponder string      `json:"transponder"`
		Heading     int         `json:"heading"`
		QnhIHg      float64     `json:"qnh_i_hg"`
		QnhMb       int         `json:"qnh_mb"`
		FlightPlan  *FlightPlan `json:"flight_plan"`
		LogonTime   string      `json:"logon_time"`
		LastUpdated string      `json:"last_updated"`
	}

	// Facility is a VATSIM controller facility descriptor
	Facility struct {
		ID    int    `json:"id"`
		Short string `json:"short"`
		Long  string `json:"long"`
	}

	// Controller is a VATSIM controller
	Controller struct {
		Cid         int      `json:"cid"`
		Name        string   `json:"name"`
		Callsign    string   `json:"callsign"`
		Frequency   string   `json:"frequency"`
		Facility    int      `json:"facility"`
		Rating      int      `json:"rating"`
		Server      string   `json:"server"`
		VisualRange int      `json:"visual_range"`
		AtisCode    string   `json:"atis_code,omitempty"`
		TextAtis    []string `json:"text_atis"`
		LastUpdated string   `json:"last_updated"`
		LogonTime   string   `json:"logon_time"`
	}

	// Data represents all the dynamic data with helper methods and index maps
	Data struct {
		General         General      `json:"general"`
		Pilots          []Pilot      `json:"pilots"`
		Controllers     []Controller `json:"controllers"`
		ATIS            []Controller `json:"atis"`
		Facilities      []Facility   `json:"facilities"`
		facilityMap     map[int]*Facility
		ctrlCallsignMap map[string]*Controller
	}
)
