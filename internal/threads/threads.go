package threads

import (
	"errors"
	"math"
	"strings"
	dbscan "telegram/internal/DBSCAN"
	"telegram/internal/datastructures"
	"telegram/internal/fasttext"
)

var langToModelName = map[string]string{
	"EN": "CatModelEN",
	"RU": "CatModelRU",
}

type Clusterable struct {
	filename   string
	title      string
	positionH1 []float64
	positionFT []float64
}

func Cosine(a []float64, b []float64) (cosine float64, err error) {
	count := 0
	length_a := len(a)
	length_b := len(b)
	if length_a > length_b {
		count = length_a
	} else {
		count = length_b
	}
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < count; k++ {
		if k >= length_a {
			s2 += math.Pow(b[k], 2)
			continue
		}
		if k >= length_b {
			s1 += math.Pow(a[k], 2)
			continue
		}
		sumA += a[k] * b[k]
		s1 += math.Pow(a[k], 2)
		s2 += math.Pow(b[k], 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0, errors.New("Vectors should not be null (all zeros)")
	}
	return sumA / (math.Sqrt(s1) * math.Sqrt(s2)), nil
}

func L2Distance(a []float64, b []float64) (cosine float64, err error) {
	var result float64 = 0
	for i := range a {
		result += math.Pow((a[i] - b[i]), 2.0)
	}

	return math.Sqrt(result), nil
}

func (s Clusterable) Distance(c interface{}) float64 {

	// distance, _ := Cosine(c.(Clusterable).positionH1, s.positionH1)
	// distanceFT, _ := Cosine(c.(Clusterable).positionFT, s.positionFT)
	distance, _ := L2Distance(c.(Clusterable).positionH1, s.positionH1)
	return distance
}

func (s Clusterable) GetID() string {
	return string(s.title)
}

func PrepareVectors(documents []datastructures.Document) []datastructures.Document {
	fasttext.LoadModel("CatModel1", "models/allcats_model_big.bin")
	for i := range documents {
		if documents[i].IsNews == false {
			continue
		}

		if documents[i].Language == "RU" || documents[i].Language == "EN" {
			words := strings.Split(documents[i].H1, " ")
			var tmpH1 [55]float32
			for w := range words {
				tmp, _ := fasttext.GetVector("CatModel1", words[w])
				for i := 0; i < 55; i++ {
					tmpH1[i] += float32(tmp[i]) / float32(len(words))
				}
			}
			documents[i].VectorH1 = tmpH1

			words = strings.Split(documents[i].FullText, " ")
			var tmpFullText [55]float32
			for w := range words {
				tmp, _ := fasttext.GetVector("CatModel1", words[w])
				for i := 0; i < 55; i++ {
					tmpFullText[i] += float32(tmp[i]) / float32(len(words))
				}
			}
			documents[i].VectorFullText = tmpFullText

		} else {
			documents[i].Category = ""
		}
	}
	return documents
}

type Thread struct {
	Name     string   `json:"title"`
	Articles []string `json:"articles"`
}

var results []Thread

func Clusters(documents []datastructures.Document) []dbscan.Cluster {

	var vectorsToAnalyze []dbscan.Clusterable
	for i := range documents {
		if documents[i].IsNews == false {
			continue
		}

		if documents[i].Language == "EN" || documents[i].Language == "RU" {
			var tmpH1 = make([]float64, len(documents[i].VectorH1))
			var tmpFT = make([]float64, len(documents[i].VectorFullText))
			for q := range documents[i].VectorH1 {
				tmpH1[q] = float64(documents[i].VectorH1[q])
				tmpFT[q] = float64(documents[i].VectorFullText[q])
			}
			vectorsToAnalyze = append(vectorsToAnalyze, Clusterable{
				filename:   documents[i].FileName,
				title:      documents[i].H1,
				positionH1: tmpH1,
				positionFT: tmpFT})
		}
	}
	//TODO: CHECK THAT THEY ARE SAME CATEGORIES
	return dbscan.Clusterize(vectorsToAnalyze, 3, 0.0001)
}

func ProcessClusters(clusters []dbscan.Cluster) []Thread {
	for i := range clusters {
		var articles []string

		for j := range clusters[i] {
			articles = append(articles, clusters[i][j].(Clusterable).filename)
		}
		results = append(results, Thread{
			Name:     clusters[i][1].(Clusterable).title,
			Articles: articles,
		})
	}

	return results
}
