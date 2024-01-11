package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() || filepath.Ext(path) != ".xml" {
			return nil
		}

		*files = append(*files, path)
		return nil
	}
}

func fingerprint(fingerprints []FingerprintDB, text string) []map[string]string {
	var matches []map[string]string

	for _, fdb := range fingerprints {
		match := fdb.MatchFirst(text)
		if match.Matched {
			matches = append(matches, match.Values)
			// j, _ := json.Marshal(match.Values)
			// fmt.Printf("%s\n", j)
		}
	}

	return matches
}

type Results struct {
	TcpScanResult int `json:"tcp_scan_result"`
}

type HostData struct {
	Port    int               `json:"port"`
	Results struct{ Results } `json:"results"`
	Banner  string            `json:"banner"`
	Recog   []map[string]string
}

// type NewHostData struct {
// 	Port    int
// 	Results Results
// 	Banner  string
// 	  []map[string]string
// }

type NewHosts map[string][]*HostData

// type ProcessedHost struct {
// 	Ip []map[string][]NewHostData
// }

func main() {
	var files []string
	if len(os.Args) < 2 {
		log.Fatalf("missing: xml directory")
	}

	err := filepath.Walk(os.Args[1], visit(&files))
	if err != nil {
		log.Fatal(err)
	}

	var fingerprints []FingerprintDB
	for _, file := range files {
		fdb, err := LoadFingerprintDBFromFile(file)
		if err != nil {
			log.Fatalf("error loading fingerprints from %s: %s", file, err)
		}
		fingerprints = append(fingerprints, fdb)
	}

	var hosts NewHosts
	byteValue, _ := os.ReadFile(os.Args[2])
	err = json.Unmarshal(byteValue, &hosts)
	if err == nil {
		for _, data := range hosts {
			for _, port := range data {
				port.Recog = fingerprint(fingerprints, port.Banner)
			}
		}
	}

	j, err := json.MarshalIndent(hosts, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile("result.json", j, 0644)
}
