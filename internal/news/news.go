package news

import (
	"telegram/internal/datastructures"
	"telegram/internal/fasttext"
)

type NewsArticleList struct {
	Articles []string `json:"articles"`
}

// var Stage2Resuls
func DetectNews(documents []datastructures.Document) (NewsArticleList, []datastructures.Document) {
	var results NewsArticleList
	// Benchmark here, should we do lang-lang or mixed mode
	fasttext.LoadModel("IsNewsModelEN", "models/isnews_model.ftz")
	fasttext.LoadModel("IsNewsModelRU", "models/isnews_ru_model.ftz")

	for i := range documents {
		if len(documents[i].FullText) < 10 {
			continue
		}

		if documents[i].Language == "EN" {
			tmp, _ := fasttext.Predict("IsNewsModelEN", documents[i].FullText, 1)
			// TODO TUNE A THRESHOLD HERE
			if _, ok := tmp["__label__good"]; ok {
				documents[i].IsNews = true
				results.Articles = append(results.Articles, documents[i].FileName)
			}

		} else if documents[i].Language == "RU" {
			tmp, _ := fasttext.Predict("IsNewsModelRU", documents[i].FullText, 1)
			if _, ok := tmp["__label__good"]; ok {
				documents[i].IsNews = true
				results.Articles = append(results.Articles, documents[i].FileName)
			}
		} else {
			documents[i].IsNews = false
		}
	}
	return results, documents
}
