package language

import (
	"telegram/internal/datastructures"
	"telegram/internal/fasttext"
)

type LanguageArticlesList struct {
	LangCode     string   `json:"lang_codes"`
	ArticlesList []string `json:"articles"`
}

func firstKey(m map[string]float32) string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys[0]
}

func getBiggestLang(m map[string]float32) string {
	keys := make([]string, 0, len(m))
	var max float32 = 0.0
	langCode := ""
	for key, value := range m {
		if value > max {
			max = value
			langCode = key
		}
		keys = append(keys, key)
	}
	return langCode
}

// DetectLanguages Detects languages in channel
func DetectLanguages(c []datastructures.Document) ([]LanguageArticlesList, []datastructures.Document) {
	var tmp map[string]float32
	var nulstruct datastructures.Document
	var Stage1Result = []LanguageArticlesList{
		LanguageArticlesList{
			LangCode: "en",
		},
		LanguageArticlesList{
			LangCode: "ru",
		},
	}

	fasttext.LoadModel("LanguageModel", "models/lid.176.ftz")
	for i := 0; i < len(c); i++ {
		finalResult := make(map[string]float32, 2)
		tmp, _ = fasttext.Predict("LanguageModel", c[i].FullText, 1)
		finalResult[firstKey(tmp)] += tmp[firstKey(tmp)]
		tmp, _ = fasttext.Predict("LanguageModel", c[i].Header, 1)
		finalResult[firstKey(tmp)] += tmp[firstKey(tmp)]
		tmp, _ = fasttext.Predict("LanguageModel", c[i].H1, 1)
		finalResult[firstKey(tmp)] += tmp[firstKey(tmp)]
		tmp, _ = fasttext.Predict("LanguageModel", c[i].H2, 1)
		finalResult[firstKey(tmp)] += tmp[firstKey(tmp)]

		articleLanguage := getBiggestLang(finalResult)
		if articleLanguage == "__label__en" {
			c[i].Language = "EN"
			Stage1Result[0].ArticlesList = append(Stage1Result[0].ArticlesList, c[i].FileName)
		} else if articleLanguage == "__label__ru" {
			c[i].Language = "RU"
			Stage1Result[1].ArticlesList = append(Stage1Result[1].ArticlesList, c[i].FileName)
		} else {
			c[i] = c[len(c)-1]      // O(1) element removal
			c[len(c)-1] = nulstruct // We don't need non-en or non-ru docs
			c = c[:len(c)-1]        //
		}
	}
	return Stage1Result, c
}
