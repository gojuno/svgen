# Summary
Svgen search for the constants of types which are the aliases for int or string types, i.e.:

```go
package db

type RideStatus string

const(
  RideStatusInitiated RideStatus = "initiated"
  RideStatusComplete  RideStatus = "complete"
)
``` 

Using found types and constants it generates sql.Scanner and driver.Valuer implementations for these types:

```go
/*
DO NOT EDIT! This code was generated automatically using github.com/gojuno/svgen v1.0
*/
package db

import (
	"database/sql/driver"
	"fmt"
)

func (t *RideStatus) Scan(i interface{}) error {
	var vv RideStatus
	switch v := i.(type) {
	case nil:
		return nil
	case []byte:
		vv = RideStatus(v)
	case string:
		vv = RideStatus(v)
	default:
		return fmt.Errorf("can't scan %T into %T", v, t)
	}

	switch vv {
	case RideStatusInitiated:
	case RideStatusComplete:
	default:
		return fmt.Errorf("invalid value of type RideStatus: %v", *t)
	}

	*t = vv

	return nil
}

func (t RideStatus) Value() (driver.Value, error) {
	if t == "" {
		return nil, nil
	}
	switch t {
	case RideStatusInitiated:
	case RideStatusComplete:
	default:
		return nil, fmt.Errorf("invalid value of type RideStatus: %v", t)
	}

	return string(t), nil
}
``` 

# Supported command line flags:
```
  -i string
    	import path of the package containing type declarations
  -o string
    	output file name (default "scanners_valuers.go")
```

# Usage of svgen in go:generate instruction:
```go
//go:generate svgen -i your.package/name -o scanners_valuers_generated.go
```
