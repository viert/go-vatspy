package static

type (
	// Point is a map point
	Point struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	// Boundaries object
	Boundaries struct {
		IsOceanic   bool    `json:"is_oceanic"`
		IsExtension bool    `json:"is_extension"`
		Min         Point   `json:"min"`
		Max         Point   `json:"max"`
		Center      Point   `json:"center"`
		Points      []Point `json:"points"`
	}

	// Country object
	Country struct {
		Name              string `json:"name"`
		Prefix            string `json:"prefix"`
		ControlCustomName string `json:"control_custom_name"`
	}

	// Airport object
	Airport struct {
		ICAO     string `json:"icao"`
		Name     string `json:"name"`
		Position Point  `json:"position"`
		IATA     string `json:"iata"`
		FIRID    string `json:"fir_id"`
		IsPseudo bool   `json:"is_pseudo"`
	}

	// FIR object
	FIR struct {
		ID         string     `json:"id"`
		Name       string     `json:"name"`
		Prefix     string     `json:"prefix"`
		ParentID   string     `json:"parent_id"`
		Boundaries Boundaries `json:"boundaries"`
	}

	// UIR object
	UIR struct {
		ID     string   `json:"id"`
		Name   string   `json:"name"`
		FIRIDs []string `json:"fir_ids"`
	}

	// Data holds all the objects and indexes
	// as well as methods to search for the data
	Data struct {
		Countries        []Country `json:"countries"`
		Airports         []Airport `json:"airports"`
		FIRs             []FIR     `json:"firs"`
		UIRs             []UIR     `json:"uirs"`
		countryNameIdx   map[string][]*Country
		countryPrefixIdx map[string]*Country
		airportICAOIdx   map[string]*Airport
		airportIATAIdx   map[string]*Airport
		firIDIdx         map[string]*FIR
		firPrefixIdx     map[string]*FIR
		uirIDIdx         map[string]*UIR
	}
)
