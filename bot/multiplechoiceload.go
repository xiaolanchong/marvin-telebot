package bot

import(
	"log"
	"io/ioutil"
	
	"gopkg.in/yaml.v2"
)

type MultipleChoiceQuestion struct {
	Text    string    `json:"text"`
	Choice  []string  `json:"choice"`
	Answer  string    `json:"answer"`
	Comment string    `json:"comment"`
}

type MultipleChoiceSection struct {
	Title    string                     `json:"title"`
	Question []MultipleChoiceQuestion   `json:"question"`
}

type MultipleChoiceLoadRoot struct {
	Description string                  `json:"description"`
	Section		[]MultipleChoiceSection `json:"section"`
}

func LoadMultipleChoiceTest(filename string) (MultipleChoiceLoadRoot, error) {
    fileData, errFile := ioutil.ReadFile(filename)
    if errFile != nil {
        log.Printf("error: %v", errFile)
		return MultipleChoiceLoadRoot{}, errFile
    }
	var data MultipleChoiceLoadRoot
	errUnmarshal := yaml.Unmarshal(fileData, &data)
    if errUnmarshal != nil {
        log.Printf("error: %v", errUnmarshal)
		return MultipleChoiceLoadRoot{}, errUnmarshal
    }
	return data, nil
}