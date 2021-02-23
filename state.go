package vatspy

import (
	"github.com/viert/go-vatspy/dynamic"
	"github.com/viert/go-vatspy/static"
)

type (
	// FIR is a VatSim static FIR with helper methods
	FIR struct {
		static.FIR
	}

	// Radar is a VatSim controller controlling a region
	Radar struct {
		dynamic.Controller
		FIRs              []*FIR `json:"firs"`
		HumanReadableName string `json:"human_readable_name"`
	}

	// AirportController is a VatSim controller controlling an airport facility
	AirportController struct {
		dynamic.Controller
	}

	// AirportControllerSet is a set of VatSim controllers attached to an airport
	AirportControllerSet struct {
		Approach *AirportController `json:"APP"`
		Delivery *AirportController `json:"DEL"`
		Ground   *AirportController `json:"GND"`
		Tower    *AirportController `json:"TWR"`
		ATIS     *AirportController `json:"ATIS"`
	}

	// Airport is a VatSim airport
	Airport struct {
		static.Airport
		Controllers AirportControllerSet `json:"controllers"`
	}

	// Country is a VatSim country
	Country struct {
		static.Country
	}

	state struct {
		airports  map[string]Airport
		countries map[string]Country
		radars    map[string]Radar
	}
)

func (c *Country) equals(other *Country) bool {
	if c == nil {
		return other == nil
	}
	if other == nil {
		return false
	}
	return c.Name == other.Name && c.Prefix == other.Prefix && c.ControlCustomName == other.ControlCustomName
}

func (c *FIR) equals(other *FIR) bool {
	if c == nil {
		return other == nil
	}
	if other == nil {
		return false
	}

	if c.ID == other.ID &&
		c.Prefix == other.Prefix &&
		c.ParentID == other.ParentID &&
		c.Name == other.Name {
		if len(c.Boundaries.Points) == len(other.Boundaries.Points) {
			for i := 0; i < len(c.Boundaries.Points); i++ {
				if c.Boundaries.Points[i].Lat != other.Boundaries.Points[i].Lat ||
					c.Boundaries.Points[i].Lng != other.Boundaries.Points[i].Lng {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (c *Radar) equals(other *Radar) bool {
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
		c.LogonTime == other.LogonTime {
		if len(c.TextAtis) == len(other.TextAtis) {
			for i := 0; i < len(c.TextAtis); i++ {
				if c.TextAtis[i] != other.TextAtis[i] {
					return false
				}
			}
		} else {
			return false
		}

		if len(c.FIRs) == len(other.FIRs) {
			for i := 0; i < len(c.FIRs); i++ {
				if !c.FIRs[i].equals(other.FIRs[i]) {
					return false
				}
				return true
			}
		}
	}
	return false
}

func (c *AirportController) equals(other *AirportController) bool {
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

func (cs *AirportControllerSet) equals(other *AirportControllerSet) bool {
	if cs == nil {
		return other == nil
	}
	if other == nil {
		return false
	}

	return cs.Delivery.equals(other.Delivery) &&
		cs.Ground.equals(other.Ground) &&
		cs.Tower.equals(other.Tower) &&
		cs.ATIS.equals(other.ATIS) &&
		cs.Approach.equals(other.Approach)
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
		a.IsPseudo == other.IsPseudo &&
		a.Position.Lat == other.Position.Lat &&
		a.Position.Lng == other.Position.Lng &&
		a.Controllers == other.Controllers
}

// IsEmpty returns true if the airport has no controllers online
func (a *Airport) IsEmpty() bool {
	c := a.Controllers
	return c.ATIS == nil && c.Delivery == nil && c.Ground == nil && c.Tower == nil
}

func (a *Airport) listControllers() []*AirportController {
	return []*AirportController{
		a.Controllers.ATIS,
		a.Controllers.Delivery,
		a.Controllers.Ground,
		a.Controllers.Tower,
		a.Controllers.Approach,
	}
}

func newStateData() *state {
	return &state{
		airports:  make(map[string]Airport),
		countries: make(map[string]Country),
		radars:    make(map[string]Radar),
	}
}
