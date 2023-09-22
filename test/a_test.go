package test

import (
	"strings"
	"testing"
)

func TestXxx(t *testing.T) {
	str := "vmess://ew0KICAidiI6ICIyIiwNCiAgInBzIjogIue+juWbvV9UZWxlZ3JhbUBreHN3YSIsDQogICJhZGQiOiAiMTU2LjIyNS42Ny4xMDEiLA0KICAicG9ydCI6ICI1Mjc5MiIsDQogICJpZCI6ICI0MTgwNDhhZi1hMjkzLTRiOTktOWIwYy05OGNhMzU4MGRkMjQiLA0KICAiYWlkIjogIjY0IiwNCiAgInNjeSI6ICJhdXRvIiwNCiAgIm5ldCI6ICJ0Y3AiLA0KICAidHlwZSI6ICJub25lIiwNCiAgImhvc3QiOiAiIiwNCiAgInBhdGgiOiAiIiwNCiAgInRscyI6ICIiLA0KICAic25pIjogIiIsDQogICJhbHBuIjogIiINCn0="
	t.Log(strings.HasPrefix(str, "vmess://"))

}
