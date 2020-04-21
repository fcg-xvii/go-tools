package jsonmap_test

import (
	"fmt"

	"github.com/fcg-xvii/go-tools/json/jsonmap"
)

func Example_basic() {
	source := map[string]interface{}{
		"v_int":    10,
		"v_string": "string value",
	}

	m := jsonmap.FromMap(source)

	// Defned int key
	fmt.Println(m.Int("v_int", 0))

	// Not defined int key, default output
	fmt.Println(m.Int("v_undefined", 0))
}

// Output:
// 10
// 0
