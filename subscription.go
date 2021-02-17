package vatspy

import "github.com/viert/go-vatspy/static"

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

type subscription struct {
	state   *state
	updates chan Update
}

func (s *subscription) processStatic(data *static.Data) {
	for _, vsAirport := range data.Airports {
		airport := Airport{
			Airport: vsAirport,
		}

		if existing, found := s.state.Airports[airport.ICAO]; found {
			airport.Controllers.ATIS = existing.Controllers.ATIS
			airport.Controllers.Delivery = existing.Controllers.Delivery
			airport.Controllers.Ground = existing.Controllers.Ground
			airport.Controllers.Tower = existing.Controllers.Tower
			if !airport.equals(&existing) {
				s.updates <- Update{objectModify, &airport}
				s.state.Airports[airport.ICAO] = airport
			}
		} else {
			s.updates <- Update{objectAdd, &airport}
			s.state.Airports[airport.ICAO] = airport
		}
	}
	for _, airport := range s.state.Airports {
		if data.FindAirportByICAO(airport.ICAO) == nil {
			s.updates <- Update{objectRemove, &airport}
			delete(s.state.Airports, airport.ICAO)
		}
	}
}
