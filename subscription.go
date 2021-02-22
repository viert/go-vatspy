package vatspy

import (
	"fmt"
	"strings"

	"github.com/op/go-logging"
	"github.com/viert/go-vatspy/dynamic"
	"github.com/viert/go-vatspy/static"
)

// UpdateType is a update type enum
type UpdateType int

// UpdateType enum definition
const (
	ObjectAdd UpdateType = iota
	ObjectModify
	ObjectRemove
)

var (
	updateTypeNames = map[UpdateType]string{
		ObjectAdd:    "add",
		ObjectModify: "modify",
		ObjectRemove: "remove",
	}
	log = logging.MustGetLogger("vatspy")
)

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

// Subscription is a subscriber descriptor
// Use Updates() to get updates channel
// Use Stop() to unsubscribe
type Subscription struct {
	subID          uint64
	state          *state
	updates        chan Update
	controlledOnly bool
	filters        []UpdateFilter
}

func (s *Subscription) processStatic(data *static.Data) {
	// in case channel is already closed
	if s.updates == nil {
		return
	}

	// process countries
	for _, vsCountry := range data.Countries {
		country := Country{
			Country: vsCountry,
		}

		if existing, found := s.state.countries[country.Prefix]; found {
			if !country.equals(&existing) {
				if s.sendUpdate(Update{ObjectModify, &country}) {
					s.state.countries[country.Prefix] = country
				}
			} else {
				if s.sendUpdate(Update{ObjectAdd, &country}) {
					s.state.countries[country.Prefix] = country
				}
			}
		}
	}

	for _, country := range s.state.countries {
		if data.FindCountryByPrefix(country.Prefix) == nil {
			if s.sendUpdate(Update{ObjectRemove, &country}) {
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
				if airport.IsEmpty() && s.controlledOnly {
					// This code should never run but just in case.
					if s.sendUpdate(Update{ObjectRemove, &airport}) {
						delete(s.state.airports, airport.ICAO)
					}
				} else {
					if s.sendUpdate(Update{ObjectModify, &airport}) {
						s.state.airports[airport.ICAO] = airport
					}
				}
			}
		} else {
			if !airport.IsEmpty() || !s.controlledOnly {
				if s.sendUpdate(Update{ObjectAdd, &airport}) {
					s.state.airports[airport.ICAO] = airport
				}
			}
		}
	}

	for _, airport := range s.state.airports {
		if data.FindAirportByICAO(airport.ICAO) == nil {
			if s.sendUpdate(Update{ObjectRemove, &airport}) {
				delete(s.state.airports, airport.ICAO)
			}
		}
	}
}

func (s *Subscription) processDynamic(dynamicData *dynamic.Data, staticData *static.Data) {
	// in case channel is already closed
	if s.updates == nil {
		return
	}

	// process controllers
	for _, vsController := range dynamicData.Controllers {
		if vsController.Facility >= 2 && vsController.Facility <= 5 {
			controller := AirportController{
				Controller: vsController,
			}

			tokens := strings.Split(controller.Callsign, "_")
			prefix := tokens[0]
			vsAirport := staticData.FindAirport(prefix)
			if vsAirport == nil {
				postfix := tokens[len(tokens)-1]
				if postfix != "OBS" && postfix != "SUP" {
					log.Debugf("can't find airport named %s, the controller is %v", prefix, controller)
				}
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
				if airportModified && s.sendUpdate(Update{ObjectModify, &airport}) {
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
				if s.sendUpdate(Update{ObjectModify, &airport}) {
					s.state.airports[airport.ICAO] = airport
				}
			}
		} else if vsController.Facility == 6 {
			// CTR
			tokens := strings.Split(vsController.Callsign, "_")
			prefix := tokens[0]
			fir := staticData.FindFIR(prefix)
			if fir == nil {
				postfix := tokens[len(tokens)-1]
				if postfix != "OBS" && postfix != "SUP" {
					log.Debugf("can't find FIR named %s, the controller is %v", prefix, vsController)
				}
				continue
			}
			radar := Radar{
				Controller: vsController,
				Boundaries: fir.Boundaries,
			}

			controlName := "Centre"
			countryPrefix := fir.ID[:2]
			country := staticData.FindCountryByPrefix(countryPrefix)
			if country != nil && country.ControlCustomName != "" {
				controlName = country.ControlCustomName
			}
			radar.HumanReadableName = fmt.Sprintf("%s %s", fir.Name, controlName)

			if existing, found := s.state.radars[radar.Callsign]; found {
				if !existing.equals(&radar) {
					if s.sendUpdate(Update{ObjectModify, &radar}) {
						s.state.radars[radar.Callsign] = radar
					}
				}
			} else {
				if s.sendUpdate(Update{ObjectAdd, &radar}) {
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

		tokens := strings.Split(atis.Callsign, "_")
		prefix := tokens[0]
		vsAirport := staticData.FindAirport(prefix)
		if vsAirport == nil {
			postfix := tokens[len(tokens)-1]
			if postfix != "OBS" && postfix != "SUP" {
				log.Debugf("can't find airport named %s, the controller is %v", prefix, atis)
			}
			continue
		}
		if airport, found := s.state.airports[vsAirport.ICAO]; found {
			existing := airport.Controllers.ATIS
			if !existing.equals(&atis) {
				airport.Controllers.ATIS = &atis
				if s.sendUpdate(Update{ObjectModify, &airport}) {
					s.state.airports[vsAirport.ICAO] = airport
				}
			}
		} else {
			airport = Airport{
				Airport: *vsAirport,
			}
			airport.Controllers.ATIS = &atis
			if s.sendUpdate(Update{ObjectModify, &airport}) {
				s.state.airports[airport.ICAO] = airport
			}
		}
	}

	// Removing controllers
	for key, airport := range s.state.airports {
		// a readonly copy to keep changed values
		current := s.state.airports[key]

		var ctrl *AirportController
		ctrl = airport.Controllers.ATIS
		if ctrl != nil {
			if dct := dynamicData.FindController(ctrl.Callsign); dct == nil {
				airport.Controllers.ATIS = nil
			}
		}
		ctrl = airport.Controllers.Delivery
		if ctrl != nil {
			if dct := dynamicData.FindController(ctrl.Callsign); dct == nil {
				airport.Controllers.Delivery = nil
			}
		}
		ctrl = airport.Controllers.Ground
		if ctrl != nil {
			if dct := dynamicData.FindController(ctrl.Callsign); dct == nil {
				airport.Controllers.Ground = nil
			}
		}
		ctrl = airport.Controllers.Tower
		if ctrl != nil {
			if dct := dynamicData.FindController(ctrl.Callsign); dct == nil {
				airport.Controllers.Tower = nil
			}
		}
		ctrl = airport.Controllers.Approach
		if ctrl != nil {
			if dct := dynamicData.FindController(ctrl.Callsign); dct == nil {
				airport.Controllers.Approach = nil
			}
		}

		if !current.equals(&airport) {
			if airport.IsEmpty() && s.controlledOnly {
				if s.sendUpdate(Update{ObjectRemove, &current}) {
					delete(s.state.airports, key)
				}
			} else {
				if s.sendUpdate(Update{ObjectModify, &airport}) {
					s.state.airports[key] = airport
				}
			}
		}
	}

	for callsign, radar := range s.state.radars {
		if ctrl := dynamicData.FindController(callsign); ctrl == nil {
			if s.sendUpdate(Update{ObjectRemove, &radar}) {
				delete(s.state.radars, callsign)
			}
		}
	}
}

func (s *Subscription) sendUpdate(update Update) bool {

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

// Updates returns a readonly updates channel
func (s *Subscription) Updates() <-chan Update {
	return s.updates
}
