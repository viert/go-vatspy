package vatspy

func newData() *Data {
	return &Data{
		Countries:        make([]Country, 0),
		Airports:         make([]Airport, 0),
		FIRs:             make([]FIR, 0),
		UIRs:             make([]UIR, 0),
		countryNameIdx:   make(map[string][]*Country),
		countryPrefixIdx: make(map[string]*Country),
		airportICAOIdx:   make(map[string]*Airport),
		airportIATAIdx:   make(map[string]*Airport),
		firIDIdx:         make(map[string]*FIR),
		uirIDIdx:         make(map[string]*UIR),
	}
}

// VATSpy data public URLs
const (
	VATSpyDataPublicURL    = "https://github.com/vatsimnetwork/vatspy-data-project/raw/master/VATSpy.dat"
	FIRBoundariesPublicURL = "https://github.com/vatsimnetwork/vatspy-data-project/raw/master/FIRBoundaries.data"
)

// FindCountriesByName searches for Country objects with a given name
//
// A Country object actually represents a single country prefix,
// so there might be multiple countries with a same name but different prefixes
func (d *Data) FindCountriesByName(name string) []*Country {
	return d.countryNameIdx[name]
}

// FindCountryByPrefix searches for a single country with a given region prefix
func (d *Data) FindCountryByPrefix(prefix string) *Country {
	return d.countryPrefixIdx[prefix]
}

// FindAirportByICAO searches for an airport by its ICAO code
func (d *Data) FindAirportByICAO(icao string) *Airport {
	return d.airportICAOIdx[icao]
}

// FindAirportByIATA searches for an airport by its IATA code
func (d *Data) FindAirportByIATA(iata string) *Airport {
	return d.airportIATAIdx[iata]
}

// FindAirport searches for an airport by its ICAO code and IATA code in that order
func (d *Data) FindAirport(id string) *Airport {
	aip := d.FindAirportByICAO(id)
	if aip == nil {
		aip = d.FindAirportByIATA(id)
	}
	return aip
}

// FindFIR searches for a given FIR by its ID
func (d *Data) FindFIR(id string) *FIR {
	return d.firIDIdx[id]
}

// FindUIR searches for a given UIR by its ID
func (d *Data) FindUIR(id string) *UIR {
	return d.uirIDIdx[id]
}

// FindUIRFIRs searches for UIR's children FIRs
func (d *Data) FindUIRFIRs(id string) []*FIR {
	uir := d.FindUIR(id)
	if uir == nil {
		return nil
	}
	firs := make([]*FIR, 0)
	for _, firID := range uir.FIRIDs {
		fir := d.FindFIR(firID)
		if fir != nil {
			firs = append(firs, fir)
		}
	}
	return firs
}

// AirportICAOCodes returns ICAO codes of all the airports in the database
func (d *Data) AirportICAOCodes() []string {
	codes := make([]string, 0)
	for code := range d.airportICAOIdx {
		codes = append(codes, code)
	}
	return codes
}

// AirportIATACodes returns IATA codes of all the airports in the database
func (d *Data) AirportIATACodes() []string {
	codes := make([]string, 0)
	for code := range d.airportIATAIdx {
		codes = append(codes, code)
	}
	return codes
}
