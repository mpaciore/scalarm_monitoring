// +build certOff

package env

import (
	"net/http"
	"crypto/tls"
)

var tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

var Client = &http.Client{Transport: tr}
const CertOff = true
