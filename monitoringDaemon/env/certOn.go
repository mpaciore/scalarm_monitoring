// +build !certOff

package env

import (
	"net/http"
)

var Client = &http.Client{}
const CertOff = false
