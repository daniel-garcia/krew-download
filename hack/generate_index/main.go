package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {

	pluginFile := ".krew.yaml"
	bs, err := ioutil.ReadFile(pluginFile)
	if err != nil {
		log.Printf("could not open %s: %s", pluginFile, err)
		os.Exit(1)
	}

	var p Plugin
	if err := yaml.Unmarshal(bs, &p); err != nil {
		log.Printf("could not unmarshal plugin: %s", err)
		os.Exit(1)
	}
	for i, _ := range p.Spec.Platforms {
		parts := strings.Split(p.Spec.Platforms[i].URI, "/")
		filename := "dist/" + parts[len(parts)-1]

		file, err := os.Open(filename)
		if err != nil {
			log.Printf("could not find artifact '%s': %s", filename, err)
			os.Exit(1)
		}
		defer file.Close()

		// Calculate the SHA256 checksum
		out, err := exec.Command("sha256sum", filename).Output()
		if err != nil {
			log.Printf("could not compute hash: %s", err)
			os.Exit(1)
		}
		checksum := strings.Fields(strings.TrimSpace(string(out)))[0]

		// Print the checksum
		p.Spec.Platforms[i].Sha256 = checksum
	}

	bs, err = yaml.Marshal(p)
	if err != nil {
		log.Printf("could not marshal plugin: %s", err)
		os.Exit(1)
	}
	if err := os.WriteFile(pluginFile, bs, 0644); err != nil {
		log.Printf("could not write plugin: %s", err)
		os.Exit(1)
	}
}

type Plugin struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Version   string `yaml:"version"`
		Platforms []struct {
			Selector struct {
				MatchLabels struct {
					Os   string `yaml:"os"`
					Arch string `yaml:"arch"`
				} `yaml:"matchLabels"`
			} `yaml:"selector"`
			URI    string `yaml:"uri"`
			Sha256 string `yaml:"sha256"`
			Files  []struct {
				From string `yaml:"from"`
				To   string `yaml:"to"`
			} `yaml:"files"`
			Bin string `yaml:"bin"`
		} `yaml:"platforms"`
		ShortDescription string `yaml:"shortDescription"`
		Homepage         string `yaml:"homepage"`
		Caveats          string `yaml:"caveats"`
		Description      string `yaml:"description"`
	} `yaml:"spec"`
}
