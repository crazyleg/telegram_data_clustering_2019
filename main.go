package main


func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func main() {

	var text string = ""
	start := time.Now()
	LoadModel("fssdf", "models/lid.176.ftz")
	//read files
	files, err := FilePathWalkDir("/Users/oleksandrsavsunenko/Downloads/DataClusteringSample0107")
	if err != nil {
		panic(err)
	}
	//fmt.Println(files)

	//open file and parse it
	for i := 0; i < 100000; i++ {
		file, err := os.Open(files[i]) // For read access.
		if err != nil {
			log.Fatal(err)
		}

		doc, err := html.Parse(file)
		if err != nil {
			// ...
		}

		var f func(*html.Node)
		f = func(n *html.Node) {
			//fmt.Println(n.Data)
			if n.Type == html.TextNode && n.Parent.Data == "p" {
				text += n.Data + " "
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)

			}
		}

		f(doc)
		var language map[string]float32
		language, _ = Predict("fssdf", text, 1)
		if _, found := language["__label__en"]; found {
			Stage1Result[0].ArticlesList = append(Stage1Result[0].ArticlesList, filepath.Base(files[i]))
		} else if _, found := language["__label__ru"]; found {
			Stage1Result[1].ArticlesList = append(Stage1Result[1].ArticlesList, filepath.Base(files[i]))
		}
		text = ""
	}
	elapsed := time.Since(start)
	log.Printf("Processing %s", elapsed)
	file, _ := json.MarshalIndent(Stage1Result, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)
}
