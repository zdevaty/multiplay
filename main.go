package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Production struct {
	XMLName xml.Name `xml:"Production"`
	CueList CueList  `xml:"CueList"`
}

type CueList struct {
	XMLName xml.Name `xml:"CueList"`
	Cues    []Cue    `xml:"Cue"`
}

type Cue struct {
	XMLName     xml.Name `xml:"Cue"`
	UID         string   `xml:"UID"`
	Type        int      `xml:"Type"`
	Enabled     int      `xml:"Enabled"`
	Q           string   `xml:"Q"`
	Description string   `xml:"Description"`
	Msgs        []Msg    `xml:"Msg"`
}

type Msg struct {
	Name    string `xml:"Name,attr"`
	Command int    `xml:"Command,attr"`
	Channel int    `xml:"Channel,attr"`
	Data1   int    `xml:"Data1,attr"`
	Data2   int    `xml:"Data2,attr"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_mpp_file>")
		return
	}
	filePath := os.Args[1]

	xmlFile, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}

	var production Production
	err = xml.Unmarshal(xmlFile, &production)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return
	}

	for _, cue := range production.CueList.Cues {
		var mics []string
		for _, msg := range cue.Msgs {
			if msg.Data2 == 1 {
				num, err := parseMicName(msg.Name)
				if err == nil {
					mics = append(mics, strconv.Itoa(num))
				}
			}
		}
		if len(mics) > 0 {
			fmt.Printf("%s | %s, Mics On: [%s]\n", cue.Q, cue.Description, strings.Join(mics, ", "))
		} else {
			fmt.Printf("%s | %s, Mics On: none\n", cue.Q, cue.Description)
		}
	}
}

func parseMicName(name string) (int, error) {
	const prefix = "Mic"
	if len(name) < len(prefix) || name[:len(prefix)] != prefix {
		return 0, errors.New("invalid mic name format: must start with 'Mic'")
	}
	numStr := name[len(prefix):]
	if numStr == "" {
		return 0, errors.New("invalid mic name format: no number after 'Mic'")
	}
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, errors.New("invalid mic name format: non-numeric characters after 'Mic'")
	}
	return num, nil
}
