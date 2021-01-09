# Edifact - specify and read arbitrary edifact document formats

A Golang module to specify edifact document formats and to read from `io.Reader` into user-defined data structures. Inspired by `json.Unmarshal` and  `xml.Unmarshal`.

## Installation
Install with `go get`
```bash
go get github.com/shogg/edifact
```
## Document specification
Document specifications are written as `go` source code: `Msg` defines a message, `S` a segment and `SG` a segment group (a loop). `"UNA"` is a segment tag. `C` and `M` stand for conditional or mandatory. The last number specifies maximal repetitions.

The notation is derived from here (example):
https://service.unece.org/trade/untdid/d96a/trmd/desadv_s.htm
```go
import "github.com/shogg/edifact/spec"

var desadv = spec.Msg("DESADV",
	spec.S("UNA", spec.C, 1),
	spec.S("UNB", spec.C, 1),
	spec.S("UNH", spec.M, 1),
	spec.S("BGM", spec.M, 1),
	spec.S("DTM", spec.C, 10),
	spec.S("ALI", spec.C, 5),
	spec.S("MEA", spec.C, 5),
	spec.S("MOA", spec.C, 5),
	spec.SG("SG1", spec.C, 10,
		spec.S("RFF", spec.M, 1),
		spec.S("DTM", spec.C, 1),
	),
[..]
```
A specification must be registered. This can be done anywhere. For instance in the main function or in an init function:
```go
func init() {
	spec.Add(desadv)
}
```
## User defined data structures
Currently only `struct` is supported (and array, slice and pointer of struct). The tag of a struct field specifies where data is located in an edifact document. All annotated information must match to make segment data suitable (segment groups, segment tag, fixed data in elements). A question mark `?` annotates the position of data.

Arrays and slices in a struct can specify a segment group path like this `edifact:"SG10/SG17"`. Each time a segment group repeats in the edifact data a new array/slice element is inserted.
```go
type Message struct {
	Date       time.Time `edifact:"DTM+17|18"`
	DeliveryNr string    `edifact:"SG1/RFF+VN|DQ+?"`
	OrderNr    int       `edifact:"SG1/RFF+ON+?"`
	Items      []Item    `edifact:"SG10/SG17"`
}
type Item struct {
	ItemNr      int    `edifact:"SG10/SG17/LIN+?"`
	Description string `edifact:"SG10/SG17/LIN+++?"`
	Quantity    int    `edifact:"SG10/SG17/QTY+12:?"`
}
```
## edifact.Unmarshaller
Out of the box only simple data types (`string`, `int`, `bool` ... ) and `time.Time` are parseable. A user defined type can implement `edifact.Unmarshaller` to provide its own parsing.
```go
// Unmarshaller interface for custom data type parsing.
type Unmarshaller interface {
	UnmarshalEdifact(data []byte) error
}
```
## Read edifact documents
The top level API is `edifact.Unmarshal`. You provide an `io.Reader` and a pointer to your data structure to be filled. Here is a complete example. It uses the built-in document specification `"DESADV"` so no doc spec here:
```go
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/shogg/edifact"
)

var ediMessage = `
UNA:+.? '
UNB+UNOC:3+sender+receiver+060620:0931+1++1234567'
UNH+1+DESADV:D:96A:UN'
BGM+220+B10001'
DTM+17:20060620:102'
RFF+ON+1'
RFF+DQ+2'
NAD+BY+++name+street+city++23436+xx'
CPS+'
LIN+1++product A:SA'
QTY+12:10'
LIN+2++product B:SA'
QTY+12:20'
CNT+2:1'
UNT+9+1'`

type Message struct {
	Date       time.Time `edifact:"DTM+17|18"`
	DeliveryNr string    `edifact:"SG1/RFF+VN|DQ+?"`
	OrderNr    int       `edifact:"SG1/RFF+ON+?"`
	Items      []Item    `edifact:"SG10/SG17"`
}
type Item struct {
	ItemNr      int    `edifact:"SG10/SG17/LIN+?"`
	Description string `edifact:"SG10/SG17/LIN+++?"`
	Quantity    int    `edifact:"SG10/SG17/QTY+12:?"`
}

func main() {
	document := strings.NewReader(ediMessage)
	var messages []*Message
	if err := edifact.Unmarshal(document, &messages); err != nil {
		panic(err)
	}
	fmt.Println("number of messages", len(messages))
}
```
## Issues
* meta characters are fixed to `UNA:+.? '`
* UNA, UNB don't get evaluated
* release (escape) character `?` is not fully implemented
* `time.Time` can only parse DTM formats 102, 201, 203
* only ORDERS and DESADV are included atm
* multiple message format versions are not supported atm

Contributions are welcome.
