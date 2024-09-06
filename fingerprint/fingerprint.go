package fingerprint

import "strings"

type Fingerprint struct {
	RawUserAgent    string
	UserAgent       *UserAgent
	AudienceNetwork bool
	Referer         string
	AppName         string
	IP              string
	AcceptLanguage  string
	AcceptEncoding  []string
	AcceptAny       bool
}

func Decode(headers map[string][]string) *Fingerprint {
	var result Fingerprint

	for header, values := range headers {
		h := strings.ToLower(header)

		if h == "user-agent" {
			result.RawUserAgent = values[0]
			result.UserAgent = decodeUserAgent(values[0])
		} else if h == "accept-language" {
			result.AcceptLanguage = parseAcceptLanguage(values[0])[0]
		} else if h == "accept-encoding" {
			result.AcceptEncoding = strings.Split(strings.ReplaceAll(values[0], " ", ""), ",")
		} else if h == "referer" {
			if values[0] == "fbapp://350685531728/unknown" {
				result.AudienceNetwork = true
			}
			result.Referer = values[0]
		} else if h == "x-requested-with" {
			result.AppName = values[0]
		} else if h == "accept" {
			if values[0] == "*/*" {
				result.AcceptAny = true
			}
		} else if h == "x-real-ip" {
			result.IP = values[0]
		}
	}

	return &result
}
