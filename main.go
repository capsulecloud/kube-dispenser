package main

import (
	"flag"
	"fmt"
	"html/template"
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

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(dir))))

	fmt.Println("Listening on :80")
	http.ListenAndServe(":80", r)
}

// GET /
func clustersGet(w http.ResponseWriter, r *http.Request) {
	out, _ := exec.Command("docker", "ps", "-a", "-f", "ancestor=capsulecloud/kube-factory", "--format", "<td>{{.Names}}</td><td>{{.CreatedAt}}</td><td>{{.Status}}</td><td><form action=\"/clusters/{{.Names}}\" method=\"post\" class=\"ui center aligned\"><button class=\"ui red button\" type=\"submit\">Delete</button></form></td>").Output()
	print(string(out))
	executeTemplate(w, "template/list.html", map[string]interface{}{"test": template.HTML(strings.Replace(string(out), "\n", "<br/>", -1))})
}

// POST /clusters
func clustersPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	clusterID := r.Form.Get("cluster-id")
	fmt.Println(clusterID)
	out, _ := exec.Command("/root/create.sh", clusterID).Output()
	print(string(out))
	http.Redirect(w, r, "/", http.StatusFound)
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
