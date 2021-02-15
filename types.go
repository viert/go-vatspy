package vatspy

type (
	// Point is a map point
	Point struct {
		Lat float64
		Lng float64
	}

	// Boundaries object
	Boundaries struct {
		Oceanic   bool
		Extension bool
		Min       Point
		Max       Point
		Center    Point
		Points    []Point
	}

	// Country object
	Country struct {
		Name              string
		Prefix            string
		ControlCustomName string
	}

	// Airport object
	Airport struct {
		ICAO     string
		Name     string
		Position Point
		IATA     string
		FIRID    string
		Pseudo   bool
	}

	// FIR object
	FIR struct {
		ID         string
		Name       string
		Prefix     string
		ParentID   string
		Boundaries Boundaries
	}

	// UIR object
	UIR struct {
		ID     string
		Name   string
		FIRIDs []string
	}

	// Data holds all the objects and indexes
	// as well as methods to search for the data
	Data struct {
		Countries        []Country
		Airports         []Airport
		FIRs             []FIR
		UIRs             []UIR
		countryNameIdx   map[string][]*Country
		countryPrefixIdx map[string]*Country
		airportICAOIdx   map[string]*Airport
		airportIATAIdx   map[string]*Airport
		firIDIdx         map[string]*FIR
		uirIDIdx         map[string]*UIR
	}
)
