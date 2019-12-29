package datastructures

// Document .. struct, produced by parser
type Document struct {
	Path     string
	FileName string
	IsNews   bool
	Category string

	VectorH1       [55]float32
	VectorFullText [55]float32

	Header         string
	H1             string
	H2             string
	Author         string
	Time           string
	FullText       string
	WebsiteAddress string
	WebsiteName    string
	Language       string
}
