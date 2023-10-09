package api

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Language struct {
	Name    string
	Quality float32
}

type Translation map[string]string

var translations = map[string]Translation{
	"uk": {
		"Already logged in":             "Ви вже увійшли",
		"Username or password is empty": "Ім'я користувача або пароль порожні",
		"Invalid credentials":           "Невірні облікові дані",
		"Successfully logged in":        "Успішний вхід",
	},
}

func parseAcceptLanguage(acceptLanguage string) (langs []string) {
	// split by comma
	parts := strings.Split(acceptLanguage, ",")
	langTypes := make([]Language, len(parts))
	for _, part := range parts {
		// split by semicolon
		langParts := strings.Split(part, ";")
		// get language name
		langName := langParts[0]
		// get quality
		quality := float32(1.0)
		if len(langParts) > 1 {
			quality = 0.0
			qualityParts := strings.Split(langParts[1], "=")
			q, err := strconv.ParseFloat(qualityParts[1], 32)
			if err == nil {
				quality = float32(q)
			}
		}
		langTypes = append(langTypes, Language{
			Name:    langName,
			Quality: quality,
		})
	}

	sort.Slice(langTypes, func(i, j int) bool {
		return langTypes[i].Quality > langTypes[j].Quality
	})

	for _, langType := range langTypes {
		langs = append(langs, langType.Name)
	}

	return langs
}

// L10n is used to get the translation of a given key (and also detect user language from header)
func L10n(key string, r *http.Request) (translation string) {
	// try to get 'Accept-Language' header and parse it
	acceptLanguage := r.Header.Get("Accept-Language")
	langs := parseAcceptLanguage(acceptLanguage)

	// try to get the translation for the given key in the user's language
	for _, lang := range langs {
		if val, ok := translations[lang][key]; ok {
			return val
		}
	}

	// if we can't find the translation, we return the key itself
	return key
}
