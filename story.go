package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

//Chapter is the struct model of each story that is mapped to an arch
type Chapter struct {
	Title   string    `json:"title"`
	Stories []string  `json:"story"`
	Options []Options `json:"options"`
}

//Options is a struct type that is to hold the options at the end of every chapter
type Options struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

//Story is a map of string and chapter to map the arc to the actual story
type Story map[string]Chapter

//JSONStory deserialises json into a map of type story
func JSONStory(r io.Reader) (Story, error) {
	var story Story
	err := json.NewDecoder(r).Decode(&story)
	if err != nil {
		return nil, err
	}
	return story, nil
}

var defaultHandlerTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Choose Your Own adventure</title>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{range .Stories}}
    <p>{{.}}</p>
   {{end}}
    <ul>
		{{range .Options}}
		<li><a href="/{{.Arc}}">{{.Text}}</a></li>
		{{end}}
    </ul>
</body>
</html>`

//NewHandler takes in a story and returns a handler
func NewHandler(s Story) http.Handler {
	return &handler{s}
}

type handler struct {
	s Story
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	path = path[1:]
	if path == "" || path == "/" {
		path = "intro"
	}
	if chapter, ok := h.s[path]; ok {
		err := tmpl.Execute(w, chapter)
		if err != nil {
			log.Println(err)
		}
		return
	}
	http.Error(w, "Chapter NOT FOUND", http.StatusNotFound)

}
