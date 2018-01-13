package main

import(

)

/*
TYPE : Page
struct for use with HTTP/TEMPLATE to display web pages.  Webpages internal data is stored here.
*/
type Page struct {
	Title string
	Body  string
}

func loadPage(title string) (*Page, error){
	return &Page{Title: title, Body: "blank..."}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/" + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

var templates = template.Must(template.ParseFiles("./html/index.html", "./html/view.html", "./html/editsmall.html","./html/results.html", "./html/login.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(index|editsmall|view|results)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
  http.HandleFunc("/", makeHandler(viewHandler))
	http.Handle("/css/",http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/fonts/",http.StripPrefix("/fonts/", http.FileServer(http.Dir("./fonts"))))
	http.Handle("/js/",http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/img/",http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	http.ListenAndServeTLS(":8085", "cert.pem", "privkey.pem", nil)
}