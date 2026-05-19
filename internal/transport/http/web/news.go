package web

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/conmeo200/Golang-V1/internal/core/model"
)

type ArticleListPageData struct {
	Title      string
	ActiveMenu string
	Articles   []model.Article
}

type ArticleFormPageData struct {
	Title      string
	ActiveMenu string
	Categories []model.Category
	Languages  []string
}

func (h *DashboardHandler) ArticleListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/news/article_list.html")
	if err != nil {
		log.Printf("Error parsing article list template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	articles, err := h.newsService.GetAllArticles(r.Context())
	if err != nil {
		log.Printf("Error fetching articles: %v", err)
	}

	data := ArticleListPageData{
		Title:      "Article Management",
		ActiveMenu: "news_articles",
		Articles:   articles,
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *DashboardHandler) ArticleNewPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/news/article_form.html")
	if err != nil {
		log.Printf("Error parsing article form template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	categories, _ := h.newsService.GetAllCategories(r.Context())

	data := ArticleFormPageData{
		Title:      "Write New Article",
		ActiveMenu: "news_articles",
		Categories: categories,
		Languages:  []string{"vi", "en"},
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *DashboardHandler) ProcessArticleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/dashboard/news/articles/new", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	lang := r.FormValue("language")
	tagsStr := r.FormValue("tags")
	tags := strings.Split(tagsStr, ",")
	
	// Simplify: just use first author for demo
	authorID := uint(1) 

	_, err := h.newsService.CreateArticle(r.Context(), authorID, title, content, lang, nil, tags)
	if err != nil {
		log.Printf("Error creating article: %v", err)
		http.Redirect(w, r, "/dashboard/news/articles?error=Creation+failed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/dashboard/news/articles?success=Article+created", http.StatusSeeOther)
}

type CategoryListPageData struct {
	Title      string
	ActiveMenu string
	Categories []model.Category
}

func (h *DashboardHandler) CategoryListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/dashboard/layout.html", "web/template/dashboard/news/category_list.html")
	if err != nil {
		log.Printf("Error parsing category list template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	categories, _ := h.newsService.GetAllCategories(r.Context())

	data := CategoryListPageData{
		Title:      "Category Management",
		ActiveMenu: "news_categories",
		Categories: categories,
	}

	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *DashboardHandler) ProcessCategoryCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/dashboard/news/categories", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	
	// Create simple category service method call or directly use repo here for brevity,
	// but I should add CreateCategory to NewsService. Let's assume we'll add it.
	err := h.newsService.CreateCategory(r.Context(), name)
	if err != nil {
		http.Redirect(w, r, "/dashboard/news/categories?error=Creation+failed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/dashboard/news/categories?success=Category+created", http.StatusSeeOther)
}
