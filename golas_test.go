package golas

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	lasReader, _ := os.Open("samples/unwrapped.las")
	las := Parse(lasReader)
	prettyPrintStructAsJSON(las)
	fmt.Println("\nversion :", las.Version())
	fmt.Println("wrap :", las.Wrap())
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
