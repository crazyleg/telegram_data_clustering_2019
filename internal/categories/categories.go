package categories

import (
	"telegram/internal/datastructures"
	"telegram/internal/fasttext"
)

type NewsCategory struct {
	Category string   `json:"category"`
	Articles []string `json:"articles"`
}

var results = []NewsCategory{
	NewsCategory{
		Category: "society",
	},
	NewsCategory{
		Category: "economy",
	},
	NewsCategory{
		Category: "technology",
	},
	NewsCategory{
		Category: "sports",
	},
	NewsCategory{
		Category: "entertainment",
	},
	NewsCategory{
		Category: "science",
	},
	NewsCategory{
		Category: "other",
	},
}

func firstKey(m map[string]float32) string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys[0]
}

var labeltoTGlabel = map[string]string{
	"__label__society":       "society",
	"__label__economy":       "economy",
	"__label__technology":    "technology",
	"__label__sports":        "sports",
	"__label__entertainment": "entertainment",
	"__label__science":       "science",
	"__label__unknown":       "other",
	"__label__other":         "other",
}

var labeltoclass = map[string]int{
	"__label__society":       0,
	"__label__economy":       1,
	"__label__technology":    2,
	"__label__sports":        3,
	"__label__entertainment": 4,
	"__label__science":       5,
	"__label__unknown":       6,
	"__label__other":         6,
}

//AssignCategories Step 3 of competitons
func AssignCategories(documents []datastructures.Document) ([]NewsCategory, []datastructures.Document) {
	fasttext.LoadModel("CatModel", "models/allcats_model_big.bin")

	for i := range documents {
		if documents[i].IsNews == false {
			continue
		}

		if documents[i].Language == "RU" || documents[i].Language == "EN" {
			text := documents[i].H1 + documents[i].FullText
			tmp, _ := fasttext.Predict("CatModel", text, 1)
			var label = firstKey(tmp)
			documents[i].Category = labeltoTGlabel[label]
			results[labeltoclass[label]].Articles = append(results[labeltoclass[label]].Articles, documents[i].FileName)
		} else {
			documents[i].Category = ""
		}
	}
	return results, documents
}
