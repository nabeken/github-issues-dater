package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var updatedRegxp = regexp.MustCompile(`([[:digit:])])+(w|d|h)`)

const githubURL = "https://github.com"

func Bind(r *mux.Router) {
	r.HandleFunc("/{user}/{repo}/issues", HandleGet).Methods("GET")
}

func HandleGet(w http.ResponseWriter, req *http.Request) {
	v := req.URL.Query()

	if v.Get("q") == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	fields := strings.Fields(v.Get("q"))
	q := make([]string, len(fields))
	for i, qv := range strings.Fields(v.Get("q")) {
		if strings.HasPrefix(qv, "updated:") {
			q[i] = ConvertRelativeDate(time.Now(), qv)
		} else {
			q[i] = qv
		}
	}

	v.Set("q", strings.Join(q, " "))

	http.Redirect(
		w,
		req,
		fmt.Sprintf("%s%s?%s", githubURL, req.URL.Path, v.Encode()),
		http.StatusFound,
	)
}

func ConvertRelativeDate(now time.Time, rd string) string {
	// return as-is
	if !strings.HasPrefix(rd, "updated:within:") {
		return rd
	}

	matches := updatedRegxp.FindStringSubmatch(rd[len("updated:within:"):])
	if len(matches) < 3 {
		return rd
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return rd
	}
	unit := matches[2]

	switch unit {
	case "w":
		year, month, day := now.AddDate(0, 0, -7*num).Date()
		return fmt.Sprintf("updated:>=%d-%0.2d-%0.2d", year, month, day)
	}
	return rd
}

func main() {
	r := mux.NewRouter()
	Bind(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	http.ListenAndServe(":"+port, r)
}
