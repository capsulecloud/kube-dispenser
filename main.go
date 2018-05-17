package main

import (
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	var dir string

	flag.StringVar(&dir, "dir", "./assets", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", clustersGet)
	r.HandleFunc("/clusters", clustersPost).Methods("POST")
	r.HandleFunc("/clusters/{id}", clusterDelete).Methods("POST")
	r.HandleFunc("/clusters/{id}", clusterGet)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(dir))))

	fmt.Println("Listening on :80")
	http.ListenAndServe(":80", r)
}

// GET /
func clustersGet(w http.ResponseWriter, r *http.Request) {
	out, _ := exec.Command("docker", "ps", "-a", "-f", "ancestor=capsulecloud/kube-factory", "--format", "<td><a href=\"/clusters/{{.Names}}\">{{.Names}}</a></td><td>{{.CreatedAt}}</td><td>{{.Status}}</td><td><form action=\"/clusters/{{.Names}}\" method=\"post\" class=\"ui center aligned\"><button class=\"ui red button\" type=\"submit\">Delete</button></form></td>").Output()
	print(string(out))
	executeTemplate(w, "template/list.html", map[string]interface{}{"data": template.HTML(strings.Replace(string(out), "\n", "<br/>", -1))})
}

// POST /clusters
func clustersPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	clusterID := r.Form.Get("cluster-id")
	if clusterID == "" {
		clusterID = randomString(8)
	}
	fmt.Println(clusterID)
	out, _ := exec.Command("/root/create.sh", clusterID).Output()
	print(string(out))
	http.Redirect(w, r, "/", http.StatusFound)
}

// GET /clusters/:id
func clusterGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["id"]
	fmt.Println(clusterID)
	out, _ := exec.Command("docker", "logs", clusterID).Output()
	executeTemplate(w, "template/show.html", map[string]interface{}{"data": template.HTML(strings.Replace(string(out), "\n", "<br/>", -1))})
}

// DELETE /clusters/:id
func clusterDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["id"]
	fmt.Println(clusterID)
	out, _ := exec.Command("/root/delete.sh", clusterID).Output()
	print(string(out))
	http.Redirect(w, r, "/", http.StatusFound)
}

func executeTemplate(w http.ResponseWriter, name string, data interface{}) {
	t, err := template.ParseFiles("template/layout.html", name)
	checkError(err)
	err = t.Execute(w, data)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
