package app

import (
    "net/http"
    "html/template"

    "appengine"
    "appengine/urlfetch"

    "code.google.com/p/codereviews-extra.rietveld/rietveld"
)

func init() {
    http.HandleFunc("/", rootHandler)
}

var (
	rootTmpl = template.Must(template.ParseFiles(
		"templates/base.html", "templates/home.html"))
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)

	issues, err := getRietveldIssues(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ctx.Infof("Got %d issues", len(issues))

	data := struct {
		PageTitle string
		Issues    *[]rietveld.Issue
		Count     int
	}{ "Sample issues list", &issues, len(issues) }

	if err := rootTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getRietveldIssues(client *http.Client) (issues []rietveld.Issue, 
	                                         err error) {
	list, err := rietveld.Search(client)
	if err != nil {
		return nil, err
	}

	issues = list.Issues
	return
}
