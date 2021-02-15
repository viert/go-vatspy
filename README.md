## go-vatspy

A golang library to fetch, parse and work with VATSpy data like countries, airports, FIRs and so on.

### Example

```(golang)
package main

import (
	"fmt"

	"github.com/viert/go-vatspy"
)

func main() {
	data, err := vatspy.Fetch(
		vatspy.VATSpyDataPublicURL,
		vatspy.FIRBoundariesPublicURL,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(data.FindAirport("KJFK"))
}
```
