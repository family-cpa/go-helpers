package fingerprint

import (
	"github.com/mileusna/useragent"
	"strings"
)

type UserAgent struct {
	Browser    string
	Device     string
	Os         string
	OsVersion  string
	DeviceType int
	Bot        bool
}

func decodeUserAgent(ua string) *UserAgent {
	agent := useragent.Parse(ua)

	return &UserAgent{
		Browser:    strings.ToLower(agent.Name),
		Device:     strings.ToLower(agent.Device),
		Os:         strings.ToLower(agent.OS),
		OsVersion:  strings.ToLower(agent.OSVersion),
		DeviceType: typeFromAgent(&agent),
		Bot:        agent.Bot,
	}
}

func typeFromAgent(agent *useragent.UserAgent) int {
	if agent.Desktop {
		return 1
	} else if agent.Tablet {
		return 2
	} else if agent.Mobile {
		return 3
	}
	return 0
}
