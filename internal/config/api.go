package config

import (
	"os"
	"strconv"
)

type api int

const APIPort api = 9999

func (a api) ToString() string {
	return strconv.Itoa(int(a))
}

func (a api) Parse() string {
	return ":" + a.ToString()
}

var APISrvStatusURL = os.Getenv("API_SRV_STATUS_URL")
