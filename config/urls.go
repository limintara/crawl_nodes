package config

import (
	"encoding/json"
	"os"
)

type ConfType struct {
	Urls     []Url
	FileName string
	Proxy    string
}
type Url struct {
	Url  string
	Reg1 string
	Reg2 string
}

func LoadConfig(file string) (ConfType, error) {
	b, err := os.ReadFile(file)
	var urls ConfType
	if err != nil {
		return ConfType{}, nil
	}
	err2 := json.Unmarshal(b, &urls)
	if err2 != nil {
		return ConfType{}, err2
	}
	return urls, nil
}
