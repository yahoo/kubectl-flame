package api

type ProgrammingLanguage string

const (
	Java   ProgrammingLanguage = "java"
	Go     ProgrammingLanguage = "go"
	Python ProgrammingLanguage = "python"
	Ruby   ProgrammingLanguage = "ruby"
	Node   ProgrammingLanguage = "node"
)

var (
	supportedLangs = []ProgrammingLanguage{Java, Go, Python, Ruby, Node}
)

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
