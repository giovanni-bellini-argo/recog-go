package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
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

type PortData struct {
	Port    int                 `json:"port"`
	Results map[string]int      `json:"results"`
	Banner  string              `json:"banner"`
	Recog   []map[string]string `json:"recog"`
}

type HostData struct {
	Ports     []*PortData `json:"ports"`
	Mac       string      `json:"mac"`
	Discovery string      `json:"discovery"`
}

type NewHosts map[string]*HostData

func main() {
	/////////// Args ///////////
	var args struct {
		InputFile  string `arg:"-i,--input,required" help:"Input json file path"`
		OutputFile string `arg:"-o,--output" default:"./result.json" help:"Output json file path"`
		XMLFolder  string `arg:"-x,--xml" default:"./xml" help:"folder containing the XML files"`
	}
	arg.MustParse(&args)

	/////////// Load XML Files ///////////
	var files []string
	err := filepath.Walk(args.XMLFolder, visit(&files))
	if err != nil {
		log.Fatal(err)
	}

	/////////// Load Fingerprints ///////////
	var fingerprints []FingerprintDB
	for _, file := range files {
		fdb, err := LoadFingerprintDBFromFile(file)
		if err != nil {
			log.Fatalf("error loading fingerprints from %s: %s", file, err)
		}
		fingerprints = append(fingerprints, fdb)
	}

	/////////// Load Input File Struct ///////////
	var hosts NewHosts
	byteValue, _ := os.ReadFile(args.InputFile)
	err = json.Unmarshal(byteValue, &hosts)
	if err == nil {
		for _, data := range hosts {
			for _, port := range data.Ports {
				// adds fingerprints to loaded struct
				port.Recog = fingerprint(fingerprints, port.Banner)
			}
		}
	}

	/////////// Write Results To File ///////////
	j, err := json.Marshal(hosts)
	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile(args.OutputFile, j, 0644)
}
