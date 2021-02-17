package vatspy

import (
	"github.com/viert/go-vatspy/dynamic"
	"github.com/viert/go-vatspy/static"
)

type (
	// Controller is a VatSim controller
	Controller struct {
		dynamic.Controller
	}

	// ControllerSet is a set of VatSim controllers attached to an airport
	ControllerSet struct {
		Delivery *Controller
		Ground   *Controller
		Tower    *Controller
		ATIS     *Controller
	}

	// Airport is a VatSim airport
	Airport struct {
		static.Airport
		Controllers ControllerSet
	}

	state struct {
		Airports map[string]Airport
	}
)

func (c *Controller) equals(other *Controller) bool {
	if c == nil {
		return other == nil
	}
	if other == nil {
		return false
	}

	if c.Cid == other.Cid &&
		c.Name == other.Name &&
		c.Callsign == other.Callsign &&
		c.Frequency == other.Frequency &&
		c.Facility == other.Facility &&
		c.Rating == other.Rating &&
		c.Server == other.Server &&
		c.VisualRange == other.VisualRange &&
		c.AtisCode == other.AtisCode &&
		c.LastUpdated == other.LastUpdated &&
		c.LogonTime == other.LogonTime {
		if len(c.TextAtis) == len(other.TextAtis) {
			for i := 0; i < len(c.TextAtis); i++ {
				if c.TextAtis[i] != other.TextAtis[i] {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (cs *ControllerSet) equals(other *ControllerSet) bool {
	if cs == nil {
		return other == nil
	}
	if other == nil {
		return false
	}

	return cs.Delivery.equals(other.Delivery) &&
		cs.Ground.equals(other.Ground) &&
		cs.Tower.equals(other.Tower) &&
		cs.ATIS.equals(other.ATIS)
}

func (a *Airport) equals(other *Airport) bool {
	if a == nil {
		return other == nil
	}
	if other == nil {
		return false
	}

	return a.ICAO == other.ICAO &&
		a.IATA == other.IATA &&
		a.Name == other.Name &&
		a.FIRID == other.FIRID &&
		a.Pseudo == other.Pseudo &&
		a.Position.Lat == other.Position.Lat &&
		a.Position.Lng == other.Position.Lng &&
		a.Controllers == other.Controllers
}

// IsEmpty returns true if the airport has no controllers online
func (a *Airport) IsEmpty() bool {
	c := a.Controllers
	return c.ATIS == nil && c.Delivery == nil && c.Ground == nil && c.Tower == nil
}

func newStateData() *state {
	return &state{
		Airports: make(map[string]Airport),
	}
}
