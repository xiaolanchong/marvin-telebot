package bot

import(
	"log"
	"io/ioutil"
	
	"gopkg.in/yaml.v2"
)

type MultipleChoiceQuestion struct {
	Text    string    `json:"text"`
	Audio   string    `json:"audio,omitempty"`
	Choice  []string  `json:"choice"`
	Answer  string    `json:"answer"`
	Comment string    `json:"comment,omitempty"`
}

type MultipleChoiceLoadRoot struct {
	Title    string                     `json:"title"`
	Question []MultipleChoiceQuestion   `json:"question"`
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