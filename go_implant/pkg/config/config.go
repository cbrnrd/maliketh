package config

import (
	"fmt"
	"maliketh/pkg/models"
)

const DEBUG = true

// !!!!!!!!!! CHANGE THESE !!!!!!!!! //
const C2_DOMAIN = "localhost"
const C2_PORT = 80
const C2_USE_TLS = false
const C2_REGISTER_PASSWORD = "SWh5bHhGOENYQWF1TW9KR3VTb0YwVkVWbDRud1RFaHc="

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! //

const INITIAL_SLEEP = 180
const REGISTER_USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)"

var CurrentProfile *models.MalleableProfile

func GetC2Url() string {
	if C2_USE_TLS {
		return fmt.Sprintf("https://%s:%d", C2_DOMAIN, C2_PORT)
	} else {
		return fmt.Sprintf("http://%s:%d", C2_DOMAIN, C2_PORT)
	}
}
