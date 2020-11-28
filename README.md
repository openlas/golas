# golas

**Parse geophysical well logs (.las files) in Golang**

- *Supported LAS Versions*:
  - [2.0](http://www.cwls.org/wp-content/uploads/2014/09/LAS_20_Update_Jan2014.pdf)

## Details

- Contains both a lexer & parser
- Self contained
  - Uses stdlib only, no third party libraries
- Easily marshal `las` data to `json`/`yaml`
- Fast, memory conscience, limited allocations
- We step through .las files one character (rune) at a time
- We fill a buffer with necessary data instead of indexing into an already loaded string
  - Buffer is reset after each line

---

## Examples

The following example uses [this](/samples/unwrapped.las) .las file as input

```golang
package main

import (
	"os"
	"github.com/openlas/golas"
)

func main() {
	lasReader, _ := os.Open("samples/unwrapped.las")
	las := golas.Parse(lasReader)
	prettyPrintStructAsJSON(las)
}

func prettyPrintStructAsJSON(v interface{}) {
	if j, e := json.MarshalIndent(v, "", "    "); e != nil {
		fmt.Printf("Error : %s \n", e.Error())
	} else {
		fmt.Printf("%s\n", string(j))
	}
}
```

Which gives us...

```
{
    "Sections": [
        {
            "Name": "Version Information",
            "Lines": [
                {
                    "Mnem": "VERS",
                    "Units": "",
                    "Data": "2.0",
                    "Description": "CWLS LOG ASCII STANDARD -VERSION 2.0"
                },
                {
                    "Mnem": "WRAP",
                    "Units": "",
                    "Data": "NO",
                    "Description": "ONE LINE PER DEPTH STEP"
                }
            ],
            "Comments": null
        },
        {
            "Name": "Well Information",
            "Lines": [
                {
                    "Mnem": "WELL",
                    "Units": "",
                    "Data": "NORVEHC MGSU 1 MITSUE 01-01",
                    "Description": "Well_name    - WELL"
                },
                {
                    "Mnem": "LOC",
                    "Units": "",
                    "Data": "00/01-01-073-05W5/0",
                    "Description": "Location     - LOCATION"
                },
                {
                    "Mnem": "UWI",
                    "Units": "",
                    "Data": "00/01-01-073-05W5/0",
                    "Description": "Uwi          - UNIQUE WELL ID"
                },
                {
                    "Mnem": "ENTR",
                    "Units": "",
                    "Data": "JOHN",
                    "Description": "Entered      - ENTERED BY"
                },
                {
                    "Mnem": "SRVC",
                    "Units": "",
                    "Data": "REGREBMULHCS",
                    "Description": "Scn          - SERVICE COMPANY"
                },
                {
                    "Mnem": "DATE",
                    "Units": "",
                    "Data": "01 JAN 70",
                    "Description": "Date         - LOG DATE"
                },
                {
                    "Mnem": "STRT",
                    "Units": "M",
                    "Data": "390",
                    "Description": "top_depth    - START DEPTH"
                },
                {
                    "Mnem": "STOP",
                    "Units": "M",
                    "Data": "650",
                    "Description": "bot_depth    - STOP DEPTH"
                },
                {
                    "Mnem": "STEP",
                    "Units": "M",
                    "Data": "0.25",
                    "Description": "increment    - STEP LENGTH"
                },
                {
                    "Mnem": "NULL",
                    "Units": "",
                    "Data": "-999.2500",
                    "Description": "NULL Value"
                }
            ],
            "Comments": [
                "#MNEM.UNIT           DATA                    DESCRIPTION OF MNEMONIC",
                "#---------    -------------------            -------------------------------",
                "# Generated from Intellog Unique Number\tCW_0099_0099/WELL/0099"
            ]
        },
        {
            "Name": "Curve Information",
            "Lines": [
                {
                    "Mnem": "DEPT",
                    "Units": "M",
                    "Data": "00 001 00 00",
                    "Description": "DEPTH        - DEPTH"
                },
                {
                    "Mnem": "DPHI",
                    "Units": "V/V",
                    "Data": "00 890 00 00",
                    "Description": "PHID         - DENSITY POROSITY (SANDSTONE)"
                },
                {
                    "Mnem": "NPHI",
                    "Units": "V/V",
                    "Data": "00 330 00 00",
                    "Description": "PHIN         - NEUTRON POROSITY (SANDSTONE)"
                },
                {
                    "Mnem": "GR",
                    "Units": "API",
                    "Data": "00 310 00 00",
                    "Description": "GR           - GAMMA RAY"
                },
                {
                    "Mnem": "CALI",
                    "Units": "MM",
                    "Data": "00 280 01 00",
                    "Description": "CAL          - CALIPER"
                },
                {
                    "Mnem": "ILD",
                    "Units": "OHMM",
                    "Data": "00 120 00 00",
                    "Description": "RESD         - DEEP RESISTIVITY (DIL)"
                }
            ],
            "Comments": [
                "#MNEM.UNIT       ERCB CURVE CODE    CURVE DESCRIPTION",
                "#-----------   ------------------   ----------------------------------"
            ]
        },
        {
            "Name": "Parameter Information",
            "Lines": [
                {
                    "Mnem": "GL",
                    "Units": "M",
                    "Data": "583.3",
                    "Description": "gl           - GROUND LEVEL ELEVATION"
                },
                {
                    "Mnem": "EREF",
                    "Units": "M",
                    "Data": "589",
                    "Description": "kb           - ELEVATION OF DEPTH REFERENCE"
                },
                {
                    "Mnem": "DATM",
                    "Units": "M",
                    "Data": "583.3",
                    "Description": "datum        - DATUM ELEVATION"
                },
                {
                    "Mnem": "TDD",
                    "Units": "M",
                    "Data": "733.4",
                    "Description": "tdd          - TOTAL DEPTH DRILLER"
                },
                {
                    "Mnem": "RUN",
                    "Units": "",
                    "Data": "ONE",
                    "Description": "Run          - RUN NUMBER"
                },
                {
                    "Mnem": "ENG",
                    "Units": "",
                    "Data": "SIMMONS",
                    "Description": "Engineer     - RECORDING ENGINEER"
                },
                {
                    "Mnem": "WIT",
                    "Units": "",
                    "Data": "SANK",
                    "Description": "Witness      - WITNESSED BY"
                },
                {
                    "Mnem": "BASE",
                    "Units": "",
                    "Data": "S.L.",
                    "Description": "Branch       - HOME BASE OF LOGGING UNIT"
                },
                {
                    "Mnem": "MUD",
                    "Units": "",
                    "Data": "GEL CHEM",
                    "Description": "Mud_type     - MUD TYPE"
                },
                {
                    "Mnem": "MATR",
                    "Units": "",
                    "Data": "SANDSTONE",
                    "Description": "Logunit      - NEUTRON MATRIX"
                },
                {
                    "Mnem": "TMAX",
                    "Units": "C",
                    "Data": "41",
                    "Description": "BHT          - MAXIMUM RECORDED TEMPERATURE"
                },
                {
                    "Mnem": "BHTD",
                    "Units": "M",
                    "Data": "733.8",
                    "Description": "BHTDEP       - MAXIMUM RECORDED TEMPERATURE"
                },
                {
                    "Mnem": "RMT",
                    "Units": "C",
                    "Data": "17",
                    "Description": "MDTP         - TEMPERATURE OF MUD"
                },
                {
                    "Mnem": "MUDD",
                    "Units": "KG/M",
                    "Data": "1100",
                    "Description": "MWT          - MUD DENSITY"
                },
                {
                    "Mnem": "NEUT",
                    "Units": "",
                    "Data": "1",
                    "Description": "NEUTRON      - NEUTRON TYPE"
                },
                {
                    "Mnem": "RESI",
                    "Units": "",
                    "Data": "0",
                    "Description": "RESIST       - RESISTIVITY TYPE"
                },
                {
                    "Mnem": "RM",
                    "Units": "OHMM",
                    "Data": "2.62",
                    "Description": "RM           - RESISTIVITY OF MUD"
                },
                {
                    "Mnem": "RMC",
                    "Units": "OHMM",
                    "Data": "0",
                    "Description": "RMC          - RESISTIVITY OF MUD CAKE"
                },
                {
                    "Mnem": "RMF",
                    "Units": "OHMM",
                    "Data": "1.02",
                    "Description": "RMF          - RESISTIVITY OF MUD FILTRATE"
                },
                {
                    "Mnem": "SUFT",
                    "Units": "C",
                    "Data": "0",
                    "Description": "SUFT         - SURFACE TEMPERATURE"
                }
            ],
            "Comments": [
                "#MNEM.UNIT           DATA             DESCRIPTION OF MNEMONIC",
                "#---------         -----------     ------------------------------"
            ]
        },
        {
            "Name": "~My Custom Section",
            "Lines": [
                {
                    "Mnem": "MNEM_VAL",
                    "Units": "UNIT_VAL",
                    "Data": "DATA_VAL",
                    "Description": "DESCRIPTION_VAL"
                }
            ],
            "Comments": null
        }
    ]
}
```
