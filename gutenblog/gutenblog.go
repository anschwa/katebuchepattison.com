package gutenblog

import (
	"bufio"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

func New(layoutTmpl, indexTmpl, postDir, outDir, webRoot string) *Blog {
	return &Blog{
		LayoutTmpl: layoutTmpl,
		IndexTmpl:  indexTmpl,
		PostDir:    postDir,
		OutDir:     outDir,
		WebRoot:    webRoot,
	}
}

type Blog struct {
	LayoutTmpl string
	IndexTmpl  string
	PostDir    string
	OutDir     string
	WebRoot    string // Might be the same as OutDir
}

func (b *Blog) Generate() error {
	// Get post metadata
	posts, err := getPosts(b.PostDir)
	if err != nil {
		return fmt.Errorf("error getting posts: %w", err)
	}
	archive := makeArchive(posts)

	// Generate blog index.html
	if err := createDir(b.OutDir); err != nil {
		return fmt.Errorf("error creating %s: %w", b.OutDir, err)
	}

	indexPath := path.Join(b.OutDir, "index.html")
	w, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("error creating %s: %w", indexPath, err)
	}
	defer w.Close()

	tmpl := template.Must(template.ParseFiles(b.LayoutTmpl, b.IndexTmpl))
	blogData := tmplData{
		Archive:       archive,
		DocumentTitle: "",
	}

	if err := tmpl.ExecuteTemplate(w, "layout", blogData); err != nil {
		w.Close()
		return fmt.Errorf("error generating %s: %w", indexPath, err)
	}

	// Generate blog posts
	for _, p := range posts {
		postDir := path.Join(b.OutDir, p.Date.Format("2006/01/02"), p.slug)

		// Create new directory for each post
		if err := createDir(postDir); err != nil {
			return fmt.Errorf("error creating %s: %w", postDir, err)
		}

		// Generate HTML from templates and write to file
		postPath := path.Join(postDir, "index.html")
		w, err = os.Create(postPath)
		if err != nil {
			return fmt.Errorf("error creating %s: %w", postPath, err)
		}

		postData := tmplData{
			DocumentTitle: p.Title,
			Archive:       archive,
		}

		postTmpl := path.Join(b.PostDir, p.filename)
		tmpl := template.Must(template.ParseFiles(postTmpl, b.LayoutTmpl))
		if err := tmpl.ExecuteTemplate(w, "layout", postData); err != nil {
			w.Close()
			return fmt.Errorf("error generating %s: %w", postPath, err)
		}
	}

	return nil
}

func (b *Blog) Serve(port string) {
	fs := http.FileServer(http.Dir(b.WebRoot))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s\t%s", r.Method, r.URL)

		// Regenerate the blog on with each request
		if err := b.Generate(); err != nil {
			log.Printf("Error generating blog: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fs.ServeHTTP(w, r)
	})

	// Adapted from:
	// - https://pkg.go.dev/net/http#ServeMux
	// - https://pkg.go.dev/net/http#Server.Shutdown
	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: mux,
	}

	idleConns := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
		close(idleConns)
	}()

	log.Printf("Starting server on: %s [%s]", srv.Addr, b.OutDir)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}

	<-idleConns
}

type tmplData struct {
	DocumentTitle string
	Archive       blogArchive
}

type blogArchive []blogMonth

type blogMonth struct {
	Date  Date
	Posts []blogPost
}

func (m *blogMonth) Title() string {
	return m.Date.Format("January 2006")
}

type blogPost struct {
	Title string
	Date  Date

	filename string
	slug     string // slug is the filename of the post on-disk and what will show up in the URL
}

func (p *blogPost) URL() string {
	return path.Join(p.Date.Format("2006/01/02"), p.slug, "/index.html")
}

func NewBlogPost(filename string, title string, date time.Time) blogPost {
	return blogPost{
		Title:    title,
		Date:     Date{Time: date},
		filename: filename,
		slug:     strings.TrimSuffix(filename, ".html.tmpl"),
	}
}

type Date struct {
	time.Time
}

func NewDate(year int, month time.Month, day int) Date {
	return Date{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func (d *Date) ISO() string {
	return d.Format("2006-01-02")
}

func (d *Date) Short() string {
	return d.Format("Jan _2")
}

func (d *Date) Suffix() string {
	switch d.Day() {
	case 1, 21, 31:
		return "st"
	case 2, 22:
		return "nd"
	case 3, 23:
		return "rd"
	default:
		return "th"
	}
}

func getFilenames(path string) ([]string, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", path, err)
	}

	files := make([]string, 0, len(dir))
	for _, f := range dir {
		if f.IsDir() {
			log.Printf("[WARN] %s is a directory, skipping...", f.Name())
			continue
		}

		files = append(files, f.Name())
	}

	return files, nil
}

func getMetadata(r io.Reader, reTitle, reDate *regexp.Regexp) (title, date string, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if match := reTitle.FindStringSubmatch(line); len(match) > 0 {
			title = match[1]
		} else if match := reDate.FindStringSubmatch(line); len(match) > 0 {
			date = match[1]
		}

		if title != "" && date != "" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return title, date, fmt.Errorf("error scanning metadata: %w", err)
	}

	return title, date, nil
}

func getPosts(postDir string) ([]blogPost, error) {
	fnames, err := getFilenames(postDir)
	if err != nil {
		return nil, fmt.Errorf("error getting filenames from %s: %w", postDir, err)
	}

	reTitle := regexp.MustCompile(`<h1>(.*)</h1>`)
	reDate := regexp.MustCompile(`<time datetime="(\d{4}-\d{2}-\d{2})"`)

	posts := make([]blogPost, 0, len(fnames))
	for _, filename := range fnames {
		f, err := os.Open(path.Join(postDir, filename))
		if err != nil {
			return nil, fmt.Errorf("error opening %s: %w", filename, err)
		}
		defer f.Close()

		title, date, err := getMetadata(f, reTitle, reDate)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("error getting metadata for %s: %w", filename, err)
		}

		pubDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, fmt.Errorf("error parsing %q as date in %s: %w", date, filename, err)
		}

		posts = append(posts, NewBlogPost(filename, title, pubDate))
	}

	return posts, nil
}

func makeArchive(posts []blogPost) blogArchive {
	// Group all the posts by month
	monthMap := make(map[time.Time][]blogPost)
	for _, p := range posts {
		// Normalize all dates to YYYY-MM: truncate day, time, etc.
		t := p.Date.Time
		m := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

		_, ok := monthMap[m]
		if !ok {
			monthMap[m] = []blogPost{}
		}

		monthMap[m] = append(monthMap[m], p)
	}

	// Sort monthMap by keys
	months := make([]time.Time, 0, len(monthMap))
	for t := range monthMap {
		months = append(months, t)
	}
	sort.SliceStable(months, func(i, j int) bool {
		return months[i].Before(months[j])
	})

	// Now build the sorted archive
	archive := make(blogArchive, 0, len(monthMap))
	for _, m := range months {
		items := monthMap[m]
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Date.Before(items[j].Date.Time)
		})

		archive = append(archive, blogMonth{
			Date:  Date{m},
			Posts: items,
		})
	}

	return archive
}

// createDir is a wrapper around os.MkdirAll and os.Chmod to achieve
// the same results as issuing "mkdir -p ..." from the command line
func createDir(dir string) error {
	if err := os.MkdirAll((dir), os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	// We need to update the directory permissions because we
	// might lose the executable bit after `umask` is applied
	if err := os.Chmod(dir, 0755); err != nil {
		return fmt.Errorf("error setting permissions on %s: %w", dir, err)
	}

	return nil
}
