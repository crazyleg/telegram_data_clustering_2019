package main

import (
	"encoding/json"
	"fmt"
	"os"
	"telegram/internal/categories"
	"telegram/internal/datastructures"
	"telegram/internal/docparser"
	"telegram/internal/language"
	"telegram/internal/news"
	"telegram/internal/threads"
)

func worker(id int, jobs <-chan datastructures.Document, results chan<- datastructures.Document) {
	for j := range jobs {
		j.AddLanguage()
		results <- j.AddLanguage()

	}
}

func main() {
	var documents []datastructures.Document

	if len(os.Args) < 2 {
		fmt.Println("expected 'languages', 'news', 'categories', 'threads' or 'top' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "languages":
		documents = docparser.ParseDoc(os.Args[2], documents)
		jobsChannels := make(chan datastructures.Document, 1)
		resultsChannels := make(chan datastructures.Document, 1)

		for i := 0; i < 25; i++ {
			go worker(i, jobsChannels, resultsChannels)
		}

		for i := 0; i < len(documents); i++ {
			jobsChannels <- documents[i]
		}
		close(jobsChannels)

		var results []datastructures.Document
		for i := 0; i < len(documents); i++ {
			results = append(results, <-resultsChannels)
		}

		// results, _ := language.DetectLanguages(documents)
		json.MarshalIndent(results, "", "   ")
		// os.Stdout.Write(json_result)

	case "news":
		documents = docparser.ParseDoc(os.Args[2], documents)
		_, documents = language.DetectLanguages(documents)
		results, _ := news.DetectNews(documents)

		json_result, _ := json.MarshalIndent(results, "", "   ")
		os.Stdout.Write(json_result)

	case "categories":
		documents = docparser.ParseDoc(os.Args[2], documents)

		_, documents = language.DetectLanguages(documents)
		_, documents = news.DetectNews(documents)
		results, _ := categories.AssignCategories(documents)
		json_result, _ := json.MarshalIndent(results, "", "   ")
		os.Stdout.Write(json_result)

	case "threads":
		documents = docparser.ParseDoc(os.Args[2], documents)
		_, documents = language.DetectLanguages(documents)
		_, documents = news.DetectNews(documents)
		_, documents := categories.AssignCategories(documents)
		documents = threads.PrepareVectors(documents)
		clusters := threads.Clusters(documents)

		results := threads.ProcessClusters(clusters)
		json_result, _ := json.MarshalIndent(results, "", "   ")
		os.Stdout.Write(json_result)

	default:
		fmt.Println("expected subcommand")
		os.Exit(1)
	}
}
