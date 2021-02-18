## go-vatspy

A golang library for working with VATSim/VATSpy data.

The library tracks two types of data: static data from VATSpy repository https://github.com/vatsimnetwork/vatspy-data-project and VATSim dynamic data from VatSim JSON3 API https://data.vatsim.net/v3/vatsim-data.json

Although you are able to use data loaders on their on importing go-vatspy/static or go-vatspy/dynamic accordingly it's recommended to use `vatspy.Provider` which merges the data from both sources and attaches airport controllers to the corresponding airports etc.

Pilots data is not tracked yet by `Provider` and only available when using `dynamic.Fetch` independently. This is going to change in the future.

### Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/viert/go-vatspy"
)

func main() {
	p, err := vatspy.New(time.Hour*24, time.Minute)
	if err != nil {
		panic(err)
	}

	sub := p.Subscribe(100)
	defer p.Stop()

	for u := range sub.Updates() {
		if arpt, ok := u.Object.(*vatspy.Airport); ok {
			fmt.Println(u.Type.String(), "airport", arpt)
		}
		if radar, ok := u.Object.(*vatspy.Radar); ok {
			fmt.Println(u.Type.String(), "radar", radar)
		}
	}
}
```
