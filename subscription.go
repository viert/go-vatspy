package vatspy

import (
	"fmt"
	"strings"

	"github.com/viert/go-vatspy/dynamic"
	"github.com/viert/go-vatspy/static"
)

const (
	objectAdd UpdateType = iota
	objectModify
	objectRemove
)

var (
	updateTypeNames = map[UpdateType]string{
		objectAdd:    "add",
		objectModify: "modify",
		objectRemove: "remove",
	}
)

// UpdateType is a update type enum
type UpdateType int

func (ut UpdateType) String() string {
	return updateTypeNames[ut]
}

// Update is an update of a state object
type Update struct {
	Type   UpdateType
	Object interface{}
}

// UpdateFilter is a callback for filtering out particular objects
type UpdateFilter func(interface{}) bool

type subscription struct {
	state   *state
	updates chan Update
	filters []UpdateFilter
}

func (s *subscription) processStatic(data *static.Data) {
	// process countries
	for _, vsCountry := range data.Countries {
		country := Country{
			Country: vsCountry,
		}

		if existing, found := s.state.countries[country.Prefix]; found {
			if !country.equals(&existing) {
				if s.sendUpdate(Update{objectModify, &country}) {
					s.state.countries[country.Prefix] = country
				}
			} else {
				if s.sendUpdate(Update{objectAdd, &country}) {
					s.state.countries[country.Prefix] = country
				}
			}
		}
	}

	for _, country := range s.state.countries {
		if data.FindCountryByPrefix(country.Prefix) == nil {
			if s.sendUpdate(Update{objectRemove, &country}) {
				delete(s.state.countries, country.Prefix)
			}
		}
	}

	// process airports
	for _, vsAirport := range data.Airports {
		airport := Airport{
			Airport: vsAirport,
		}

		if existing, found := s.state.airports[airport.ICAO]; found {
			airport.Controllers.ATIS = existing.Controllers.ATIS
			airport.Controllers.Delivery = existing.Controllers.Delivery
			airport.Controllers.Ground = existing.Controllers.Ground
			airport.Controllers.Tower = existing.Controllers.Tower
			if !airport.equals(&existing) {
				if s.sendUpdate(Update{objectModify, &airport}) {
					s.state.airports[airport.ICAO] = airport
				}
			}
		} else {
			if s.sendUpdate(Update{objectAdd, &airport}) {
				s.state.airports[airport.ICAO] = airport
			}
		}
	}
	for _, airport := range s.state.airports {
		if data.FindAirportByICAO(airport.ICAO) == nil {
			if s.sendUpdate(Update{objectRemove, &airport}) {
				delete(s.state.airports, airport.ICAO)
			}
		}
	}
}

func (s *subscription) processDynamic(dynamicData *dynamic.Data, staticData *static.Data) {
	// process controllers
	for _, vsController := range dynamicData.Controllers {
		if vsController.Facility >= 2 && vsController.Facility <= 5 {
			controller := AirportController{
				Controller: vsController,
			}

			prefix := strings.Split(controller.Callsign, "_")[0]
			vsAirport := staticData.FindAirport(prefix)
			if vsAirport == nil {
				fmt.Printf("can't find airport named %s, the controller is %v\n", prefix, controller)
				continue
			}

			if airport, found := s.state.airports[vsAirport.ICAO]; found {
				var existing *AirportController
				airportModified := false

				switch controller.Facility {
				case 2:
					existing = airport.Controllers.Delivery
					if !existing.equals(&controller) {
						airport.Controllers.Delivery = &controller
						airportModified = true
					}
				case 3:
					existing = airport.Controllers.Ground
					if !existing.equals(&controller) {
						airport.Controllers.Ground = &controller
						airportModified = true
					}
				case 4:
					existing = airport.Controllers.Tower
					if !existing.equals(&controller) {
						airport.Controllers.Tower = &controller
						airportModified = true
					}
				case 5:
					existing = airport.Controllers.Approach
					if !existing.equals(&controller) {
						airport.Controllers.Approach = &controller
						airportModified = true
					}
				}
				if airportModified && s.sendUpdate(Update{objectModify, &airport}) {
					s.state.airports[airport.ICAO] = airport
				}
			} else {
				airport = Airport{
					Airport: *vsAirport,
				}
				switch controller.Facility {
				case 2:
					airport.Controllers.Delivery = &controller
				case 3:
					airport.Controllers.Ground = &controller
				case 4:
					airport.Controllers.Tower = &controller
				case 5:
					airport.Controllers.Approach = &controller
				}
				if s.sendUpdate(Update{objectModify, &airport}) {
					s.state.airports[airport.ICAO] = airport
				}
			}
		} else if vsController.Facility == 6 {
			// CTR
			prefix := strings.Split(vsController.Callsign, "_")[0]
			fir := staticData.FindFIR(prefix)
			if fir == nil {
				fmt.Printf("can't find FIR named %s, the controller is %v\n", prefix, vsController)
				continue
			}
			radar := Radar{
				Controller: vsController,
				Boundaries: fir.Boundaries,
			}

			if existing, found := s.state.radars[radar.Callsign]; found {
				if !existing.equals(&radar) {
					if s.sendUpdate(Update{objectModify, &radar}) {
						s.state.radars[radar.Callsign] = radar
					}
				}
			} else {
				if s.sendUpdate(Update{objectAdd, &radar}) {
					s.state.radars[radar.Callsign] = radar
				}
			}
		}
	}

	// process ATIS stations
	for _, vsATIS := range dynamicData.ATIS {
		atis := AirportController{
			Controller: vsATIS,
		}

		prefix := strings.Split(atis.Callsign, "_")[0]
		vsAirport := staticData.FindAirport(prefix)
		if vsAirport == nil {
			fmt.Printf("can't find airport named %s, the controller is %v\n", prefix, atis)
			continue
		}
		if airport, found := s.state.airports[vsAirport.ICAO]; found {
			existing := airport.Controllers.ATIS
			if !existing.equals(&atis) {
				airport.Controllers.ATIS = &atis
				if s.sendUpdate(Update{objectModify, &airport}) {
					s.state.airports[vsAirport.ICAO] = airport
				}
			}
		} else {
			airport = Airport{
				Airport: *vsAirport,
			}
			airport.Controllers.ATIS = &atis
			if s.sendUpdate(Update{objectModify, &airport}) {
				s.state.airports[airport.ICAO] = airport
			}
		}
	}
}

func (s *subscription) sendUpdate(update Update) bool {
	// apply filters before sending anything
	for _, filter := range s.filters {
		// if a filter returns false, do not send anything,
		// but pretend we just did
		if !filter(update.Object) {
			return true
		}
	}

	// nonblocking send, drops the message if the channel buffer is full
	select {
	case s.updates <- update:
		return true
	default:
		return false
	}
}
