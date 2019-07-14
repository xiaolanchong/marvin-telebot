
package bot

import (
	"log"
	"io/ioutil"
	
	"gopkg.in/yaml.v2"
)

type Unit struct {
	Title    string `yaml:"title"`
	Exercise []struct {
		Description string `yaml:"description"`
		Question    []struct {
			Text    string `yaml:"text"`
			Answer  string `yaml:"answer"`
			Show    string `yaml:"show,omitempty"`
		} `yaml:"question"`
	} `yaml:"exercise"`
} 

type unitDataImport struct {
	Units []Unit `yaml:"unit"`
}

type Units = []Unit

func LoadUnits(filename string) (Units, error) {
	fileData, errFile := ioutil.ReadFile(filename)
    if errFile != nil {
        log.Printf("error: %v", errFile)
		return Units{}, errFile
    }
	var data unitDataImport
	errUnmarshal := yaml.Unmarshal(fileData, &data)
    if errUnmarshal != nil {
        log.Printf("error: %v", errUnmarshal)
		return Units{}, errUnmarshal
    }
	return data.Units, nil
}