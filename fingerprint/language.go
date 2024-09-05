package fingerprint

import (
	"sort"
	"strconv"
	"strings"
)

type language struct {
	name    string
	quality float64
}

type languageSlice []language

func (ls languageSlice) SortByQuality() {
	sort.Sort(ls)
}

func (ls languageSlice) Len() int {
	return len(ls)
}

func (ls languageSlice) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func (ls languageSlice) Less(i, j int) bool {
	return ls[i].quality > ls[j].quality
}

func parseAcceptLanguage(languages string) []string {
	preferredLanguages := strings.Split(languages, ",")
	preferredLanguagesLen := len(preferredLanguages)

	langsCap := preferredLanguagesLen
	langs := make(languageSlice, 0, langsCap)

	for i, rawPreferredLanguage := range preferredLanguages {
		preferredLanguage := strings.Replace(strings.ToLower(strings.TrimSpace(rawPreferredLanguage)), "_", "-", 0)
		if preferredLanguage == "" {
			continue
		}

		parts := strings.SplitN(preferredLanguage, ";", 2)

		lang := language{parts[0], 0}
		if len(parts) == 2 {
			q := parts[1]

			if strings.HasPrefix(q, "q=") {
				q = strings.SplitN(q, "=", 2)[1]
				var err error
				if lang.quality, err = strconv.ParseFloat(q, 64); err != nil {
					lang.quality = 1
				}
			}
		}

		if lang.quality == 0 {
			lang.quality = float64(preferredLanguagesLen - i)
		}

		langs = append(langs, lang)

	}

	langs.SortByQuality()

	langString := make([]string, 0, len(langs))
	for _, lang := range langs {
		langString = append(langString, lang.name)
	}

	if len(langString) < 1 {
		langString = append(langString, "")
	}

	return langString
}
