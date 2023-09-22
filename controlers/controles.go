package controlers

import (
	"encoding/json"

	clash "github.com/sanzhang007/webgin/protocol"
)

type ServerAndPortAndPassword struct {
	Server   string
	Port     string
	Password string
}

func GetPotocolServerAndPortAndPassword(link string) ([]byte, error) {

	c, err := clash.ClashParse(&link)
	if err != nil {
		return nil, nil
	}
	b, err := json.Marshal(c)
	if err != nil {
		return nil, nil
	}
	var s ServerAndPortAndPassword
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, nil
	}
	return json.Marshal(s)
}
