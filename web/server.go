package web

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/LiterallyEthical/r3conwhal3/pkg/logger"
)

var (
	myLogger = logger.GetLogger()
)

type Image struct {
	Name string
	URL  string
}

type PageData struct {
	Images       []Image
	HasPrev      bool
	HasNext      bool
	PrevPage     int
	NextPage     int
	CurrentPage  int
	TotalPages   int
	PageNumbers  []int
	ShowEllipsis bool
}

func trimSuffix(suffix, s string) string {
	return strings.TrimSuffix(s, suffix)
}

func indexSequence(start, end int) []int {
	seq := make([]int, end-start)
	for i := range seq {
		seq[i] = start + i
	}
	return seq
}

func inc(i int) int {
	return i + 1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// resolvePath resolves the given relative path to an absolute path based on the current file location
func resolvePath(relPath string) (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", os.ErrInvalid
	}
	baseDir := filepath.Dir(filename)
	return filepath.Join(baseDir, relPath), nil
}

func StartServer(imageDir string) error {
	// Resolve the path to the templates directory
	templatePath, err := resolvePath("templates/index.html")
	if err != nil {
		return fmt.Errorf("Failed to resolve template path: %v", err)
	}

	tmpl := template.New("index.html").Funcs(template.FuncMap{
		"trimSuffix":    trimSuffix,
		"indexSequence": indexSequence,
		"inc":           inc,
	})
	tmpl = template.Must(tmpl.ParseFiles(templatePath))

	imagesPerPage := 10

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imageDir))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Handle panics
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				myLogger.Error("Recovered from panic:", r)
			}
		}()

		// Set no-cache headers
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		var images []Image
		files, err := os.ReadDir(imageDir)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			myLogger.Error("Error reading directory:", err)
			return
		}

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		start := (page - 1) * imagesPerPage
		end := start + imagesPerPage
		if end > len(files) {
			end = len(files)
		}

		if start >= len(files) {
			http.Error(w, "Page out of range:", http.StatusBadRequest)
			myLogger.Error("Page out of range: requested page", page)
		}

		for _, file := range files[start:end] {
			if !file.IsDir() {
				fileName := file.Name()
				// Remove the file extension
				url := strings.TrimSuffix(fileName, ".png")
				// Replace the first dash with "://"
				protocolIndex := strings.Index(url, "-")
				if protocolIndex != -1 {
					url = url[:protocolIndex] + "://" + url[protocolIndex+1:]
				}
				// Replace remaining dashes with dots
				// url := strings.Replace(url, "-", ".", -1)
				images = append(images, Image{Name: fileName, URL: url})
			}
		}

		totalPages := (len(files) + imagesPerPage - 1) / imagesPerPage

		// Create pagination numbers with ellipsis logic
		pageNumbers := []int{}
		showEllipsis := totalPages > 5

		if totalPages <= 5 {
			pageNumbers = indexSequence(1, totalPages+1)
		} else {
			pageNumbers = append(pageNumbers, 1)
			if page > 4 {
				pageNumbers = append(pageNumbers, -1) // -1 will indicate ellipsis
			}
			start := max(2, min(totalPages-2, page-1))
			end := min(totalPages-1, page+1)
			pageNumbers = append(pageNumbers, indexSequence(start, end+1)...)
			if page < totalPages-3 {
				pageNumbers = append(pageNumbers, -1) // -1 will indicate ellipsis
			}
			pageNumbers = append(pageNumbers, totalPages)
		}

		data := PageData{
			Images:       images,
			HasPrev:      page > 1,
			HasNext:      end < len(files),
			PrevPage:     page - 1,
			NextPage:     page + 1,
			CurrentPage:  page,
			TotalPages:   totalPages,
			PageNumbers:  pageNumbers,
			ShowEllipsis: showEllipsis,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			myLogger.Error("Error executing template:", err)
		}
	})

	myLogger.Info("R3conwhal3 Web Galery running on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return fmt.Errorf("ListenAndServe: %v", err)
	}

	return nil
}
