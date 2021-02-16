## go-vatspy

A golang library to fetch, parse and work with VATSpy data like countries, airports, FIRs and so on.
The original date for the library is taken from https://github.com/vatsimnetwork/vatspy-data-project by default.

For caching and performance purposes you can provide your own data from a local file.

### Example

```go
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
