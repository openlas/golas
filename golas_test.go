package golas

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	lasReader, _ := os.Open("samples/1.las")
	_ = Parse(lasReader)
	// prettyPrintStructAsJSON(las)
}

func BenchmarkTest(b *testing.B) {
	lasReader, _ := os.Open("samples/unwrapped.las")
	Parse(lasReader)
}

func prettyPrintStructAsJSON(v interface{}) {
	if j, e := json.MarshalIndent(v, "", "    "); e != nil {
		fmt.Printf("Error : %s \n", e.Error())
	} else {
		fmt.Printf("%s\n", string(j))
	}
}
