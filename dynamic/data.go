package dynamic

import "encoding/json"

func newData(raw []byte) (*Data, error) {
	var data Data
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}

	data.facilityMap = make(map[int]*Facility)
	for i, facility := range data.Facilities {
		data.facilityMap[facility.ID] = &data.Facilities[i]
	}

	data.ctrlCallsignMap = make(map[string]*Controller)
	for i, ctrl := range data.Controllers {
		data.ctrlCallsignMap[ctrl.Callsign] = &data.Controllers[i]
	}

	for i, atis := range data.ATIS {
		data.ctrlCallsignMap[atis.Callsign] = &data.ATIS[i]
	}

	return &data, nil
}

// FindController finds a controller by callsign
func (d *Data) FindController(cs string) *Controller {
	return d.ctrlCallsignMap[cs]
}
