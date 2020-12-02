package golas

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestHelloWorld(t *testing.T) {
	s := time.Now()
	lasReader, _ := os.Open("samples/unwrapped.las")
	p := Parse(lasReader)
	fmt.Printf("took %ss\n", time.Since(s))
	prettyPrintStructAsJSON(p)
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
