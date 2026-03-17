package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/conmeo200/Golang-V1/internal/dto"
)

type WebHandler struct {
    // We can inject services here later if we want the web to talk to DB
}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

type NewsItem struct {
	Title    string
	Category string
	Excerpt  string
	Author   string
	Date     string
}

type NewsPageData struct {
	Title    string
	Heading  string
	NewsList []NewsItem
}

func (h *WebHandler) NewsPage(w http.ResponseWriter, r *http.Request) {
	// 1. Parsing the template
	tmpl, err := template.ParseFiles("web/template/news.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		dto.RespondWithError(w, dto.ErrInternal)
		return
	}

	// 2. Preparing Mock Data
	data := NewsPageData{
		Title:   "HeeloNews - Latest Updates",
		Heading: "Latest Insights & Updates",
		NewsList: []NewsItem{
			{
				Title:    "Go 1.25 Released with Context Improvements",
				Category: "Technology",
				Excerpt:  "The latest version of Go introduces new ways to handle request scoping and context cancellation effectively across massive microservices. Engineers at major tech companies are already reporting up to 30% reduction in memory leaks and dangling goroutines. \"This is the paradigm shift we've been waiting for,\" says a lead developer from the Go Team. The release also brings subtle compiler optimizations that speed up build times for large applications.",
				Author:   "Alice Wonderland",
				Date:     "Oct 15, 2026",
			},
			{
				Title:    "Why Docker is still relevant in 2026",
				Category: "DevOps",
				Excerpt:  "Despite the rise of serverless and new orchestration tools, containerization remains the backbone of modern cloud native applications. The community continues to find ways to strip down image sizes and improve security profiles. Our correspondent explores how the old whale learned some new tricks.",
				Author:   "Bob Builder",
				Date:     "Oct 12, 2026",
			},
			{
				Title:    "Building Aesthetic UIs with Pure CSS",
				Category: "Design",
				Excerpt:  "You don't always need a heavy CSS framework. Discover how pure CSS variables and flexbox can create stunning, responsive interfaces. We dive into the revival of traditional newspaper layouts using modern CSS Grid.",
				Author:   "Charlie Chaplin",
				Date:     "Oct 10, 2026",
			},
			{
				Title:    "The Return of the Monolith",
				Category: "Architecture",
				Excerpt:  "Some startups are ditching microservices and returning to majestic monoliths to reduce operational overhead. Is this a step backwards or a pragmatic evolution? A detailed report on the industry's shifting tide.",
				Author:   "Diana Prince",
				Date:     "Oct 08, 2026",
			},
			{
				Title:    "Stock Markets Rally on AI News",
				Category: "Finance",
				Excerpt:  "Tech stocks saw an unprecedented surge today following major breakthroughs in generative AI capabilities. Analysts are calling it the 'Golden Age of Silicon'. Could this momentum last through the quarter?",
				Author:   "Edison Teller",
				Date:     "Oct 16, 2026",
			},
		},
	}

	// 3. Executing the template with data
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		dto.RespondWithError(w, dto.NewAppError(http.StatusInternalServerError, "Error rendering page", "RENDER_ERROR"))
	}
}
