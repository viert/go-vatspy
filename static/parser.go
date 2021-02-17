package static

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type parseStaet int

const (
	stateReadCategory parseStaet = iota
	stateReadCountries
	stateReadAirports
	stateReadFIRs
	stateReadUIRs

	coordErrorTemplate = "invalid lat/lng value '%s' on line %d"
)

func parse(dataRaw []byte, boundariesRaw []byte) (*Data, error) {
	bds, err := parseBoundaries(boundariesRaw)
	if err != nil {
		return nil, err
	}

	data, err := parseData(dataRaw, bds)
	if err != nil {
		return nil, err
	}

	makeIndexes(data)

	return data, nil
}

func parseData(data []byte, boundaries map[string]Boundaries) (*Data, error) {
	results := newData()

	state := stateReadCategory
	sc := bufio.NewScanner(bytes.NewReader(data))
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := strings.TrimSpace(sc.Text())

		if len(line) == 0 || line[0] == ';' {
			// skip comments and empty lines
			continue
		}

		tokens := strings.Split(line, "|")

	redecide:
		switch state {
		case stateReadCategory:
			if line[0] == '[' {
				cat := line[1 : len(line)-1]
				cat = strings.ToLower(cat)
				switch cat {
				case "countries":
					state = stateReadCountries
				case "airports":
					state = stateReadAirports
				case "firs":
					state = stateReadFIRs
				case "uirs":
					state = stateReadUIRs
				}
			}
		case stateReadCountries:
			if len(tokens) != 3 {
				state = stateReadCategory
				goto redecide
			}
			country := Country{tokens[0], tokens[1], tokens[2]}
			results.Countries = append(results.Countries, country)
		case stateReadAirports:
			if len(tokens) != 7 {
				state = stateReadCategory
				goto redecide
			}
			lat, err := strconv.ParseFloat(tokens[2], 64)
			if err != nil {
				return nil, fmt.Errorf(coordErrorTemplate, tokens[2], lineNum)
			}
			lng, err := strconv.ParseFloat(tokens[3], 64)
			if err != nil {
				return nil, fmt.Errorf(coordErrorTemplate, tokens[3], lineNum)
			}

			airport := Airport{
				ICAO:     tokens[0],
				Name:     tokens[1],
				Position: Point{Lat: lat, Lng: lng},
				IATA:     tokens[4],
				FIRID:    tokens[5],
				Pseudo:   tokens[6] == "1",
			}
			results.Airports = append(results.Airports, airport)
		case stateReadFIRs:
			if len(tokens) != 4 {
				state = stateReadCategory
				goto redecide
			}

			fir := FIR{
				ID:       tokens[0],
				Name:     tokens[1],
				Prefix:   tokens[2],
				ParentID: tokens[3],
			}

			if bnds, found := boundaries[fir.ID]; found {
				fir.Boundaries = bnds
			}
			results.FIRs = append(results.FIRs, fir)

		case stateReadUIRs:
			if len(tokens) != 3 {
				state = stateReadCategory
				goto redecide
			}

			firIDs := strings.Split(tokens[2], ",")

			uir := UIR{
				ID:     tokens[0],
				Name:   tokens[1],
				FIRIDs: firIDs,
			}
			results.UIRs = append(results.UIRs, uir)
		}
	}

	return results, nil
}

func parseBoundaries(data []byte) (map[string]Boundaries, error) {

	var points []Point
	var icao string
	var current Boundaries

	boundaries := make(map[string]Boundaries)

	sc := bufio.NewScanner(bytes.NewReader(data))
	pointsLeft := 0
	pointIdx := 0
	lineNum := 0

	for sc.Scan() {
		lineNum++

		line := strings.TrimSpace(sc.Text())
		if len(line) == 0 {
			continue
		}

		tokens := strings.Split(line, "|")
		if pointsLeft <= 0 {
			// the line is a header
			if len(tokens) != 10 {
				return nil, fmt.Errorf("invalid header '%s' on line %d", line, lineNum)
			}

			icao = tokens[0]
			pCount, err := strconv.ParseInt(tokens[3], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid point count value '%s' on line %d", tokens[3], lineNum)
			}
			pointsLeft = int(pCount)

			minLat, err := strconv.ParseFloat(tokens[4], 64)
			if err != nil {
				return nil, fmt.Errorf(coordErrorTemplate, tokens[4], lineNum)
			}
			minLng, err := strconv.ParseFloat(tokens[5], 64)
			if err != nil {
				return nil, fmt.Errorf(coordErrorTemplate, tokens[5], lineNum)
			}
			maxLat, err := strconv.ParseFloat(tokens[6], 64)
			if err != nil {
				return nil, fmt.Errorf(coordErrorTemplate, tokens[6], lineNum)
			}
			maxLng, err := strconv.ParseFloat(tokens[7], 64)
			if err != nil {
				return nil, fmt.Errorf(coordErrorTemplate, tokens[7], lineNum)
			}
			cntLat, err := strconv.ParseFloat(tokens[8], 64)
			if err != nil {
				return nil, fmt.Errorf(coordErrorTemplate, tokens[8], lineNum)
			}
			cntLng, err := strconv.ParseFloat(tokens[9], 64)
			if err != nil {
				return nil, fmt.Errorf(coordErrorTemplate, tokens[9], lineNum)
			}

			current = Boundaries{
				Oceanic:   tokens[1] == "1",
				Extension: tokens[2] == "1",
				Min:       Point{Lat: minLat, Lng: minLng},
				Max:       Point{Lat: maxLat, Lng: maxLng},
				Center:    Point{Lat: cntLat, Lng: cntLng},
			}
			points = make([]Point, pointsLeft)
			pointIdx = 0
		} else {
			// parse point
			if len(tokens) != 2 {
				return nil, fmt.Errorf("invalid point '%s' on line %d", line, lineNum)
			}
			lat, err := strconv.ParseFloat(tokens[0], 64)
			if err != nil {
				return nil, fmt.Errorf(coordErrorTemplate, tokens[0], lineNum)
			}
			lng, err := strconv.ParseFloat(tokens[0], 64)
			if err != nil {
				return nil, fmt.Errorf(coordErrorTemplate, tokens[0], lineNum)
			}
			points[pointIdx] = Point{Lat: lat, Lng: lng}
			pointIdx++
			pointsLeft--

			if pointsLeft == 0 {
				current.Points = points
				boundaries[icao] = current
			}
		}
	}
	return boundaries, nil
}

func makeIndexes(data *Data) {
	for i := range data.Countries {
		country := &data.Countries[i]
		data.countryPrefixIdx[country.Prefix] = country

		if _, found := data.countryNameIdx[country.Name]; !found {
			data.countryNameIdx[country.Name] = make([]*Country, 0)
		}

		cnIndex := data.countryNameIdx[country.Name]
		cnIndex = append(cnIndex, country)
		data.countryNameIdx[country.Name] = cnIndex
	}

	for i := range data.Airports {
		airport := &data.Airports[i]
		data.airportIATAIdx[airport.IATA] = airport
		data.airportICAOIdx[airport.ICAO] = airport
	}

	for i := range data.FIRs {
		fir := &data.FIRs[i]
		data.firIDIdx[fir.ID] = fir
	}

	for i := range data.UIRs {
		uir := &data.UIRs[i]
		data.uirIDIdx[uir.ID] = uir
	}
}
