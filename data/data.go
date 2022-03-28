package data

const (
	Viewing     int8 = 1
	Downloading      = 2
)

type PresentableCodecast struct {
	Title           string
	Description     string
	Picture         string
	PublicationDate string
	IsViewable      bool
	IsDownloadable  bool
}
