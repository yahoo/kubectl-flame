package api

type ProgrammingLanguage string

const (
	Java ProgrammingLanguage = "java"
	Go   ProgrammingLanguage = "go"
)

var (
	supportedLangs      = []ProgrammingLanguage{Java, Go}
	requiresProcessName = []ProgrammingLanguage{Go}
)

func (l ProgrammingLanguage) IsPgrepRequired() bool {
	return containsLang(l, requiresProcessName)
}

func AvailableLanguages() []ProgrammingLanguage {
	return supportedLangs
}

func IsSupportedLanguage(lang string) bool {
	return containsLang(ProgrammingLanguage(lang), AvailableLanguages())
}

func containsLang(l ProgrammingLanguage, langs []ProgrammingLanguage) bool {
	for _, current := range langs {
		if l == current {
			return true
		}
	}

	return false
}
