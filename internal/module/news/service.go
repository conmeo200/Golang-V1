package news

import (
	"context"
	//"fmt"

	"github.com/conmeo200/Golang-V1/internal/core/model"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type NewsRepository interface {
	// Articles
	CreateArticle(ctx context.Context, article *model.Article) error
	GetArticleByID(ctx context.Context, id uint) (*model.Article, error)
	ListArticles(ctx context.Context, filter map[string]interface{}) ([]model.Article, error)
	UpdateArticle(ctx context.Context, article *model.Article) error
	
	// Categories
	CreateCategory(ctx context.Context, category *model.Category) error
	ListCategories(ctx context.Context) ([]model.Category, error)
	
	// Tags
	GetOrCreateTags(ctx context.Context, tagNames []string) ([]model.Tag, error)
	
	WithTx(tx *gorm.DB) NewsRepository
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) WithTx(tx *gorm.DB) NewsRepository {
	return &newsRepository{db: tx}
}

func (r *newsRepository) CreateArticle(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *newsRepository) GetArticleByID(ctx context.Context, id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Translations").
		Preload("Categories").
		Preload("Tags").
		Preload("Stats").
		First(&article, id).Error
	return &article, err
}

func (r *newsRepository) ListArticles(ctx context.Context, filter map[string]interface{}) ([]model.Article, error) {
	var articles []model.Article
	query := r.db.WithContext(ctx).Preload("Author").Preload("Categories").Order("created_at DESC")
	
	if status, ok := filter["status"]; ok {
		query = query.Where("status = ?", status)
	}
	
	err := query.Find(&articles).Error
	return articles, err
}

func (r *newsRepository) UpdateArticle(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *newsRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *newsRepository) ListCategories(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func (r *newsRepository) GetOrCreateTags(ctx context.Context, tagNames []string) ([]model.Tag, error) {
	var tags []model.Tag
	for _, name := range tagNames {
		var tag model.Tag
		tagNameSlug := slug.Make(name)
		err := r.db.WithContext(ctx).FirstOrCreate(&tag, model.Tag{Name: name, Slug: tagNameSlug}).Error
		if err == nil {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

// Service Implementation
type NewsService struct {
	repo NewsRepository
}

func NewNewsService(repo NewsRepository) *NewsService {
	return &NewsService{repo: repo}
}

func (s *NewsService) CreateArticle(ctx context.Context, authorID uint, title, content, lang string, catIDs []uint, tagNames []string) (*model.Article, error) {
	articleSlug := slug.Make(title)
	
	tags, _ := s.repo.GetOrCreateTags(ctx, tagNames)
	
	var categories []model.Category
	for _, id := range catIDs {
		categories = append(categories, model.Category{ID: id})
	}

	article := &model.Article{
		AuthorID: authorID,
		Slug:     articleSlug,
		Status:   model.ArticleStatusDraft,
		Categories: categories,

		Tags:       tags,
		Translations: []model.ArticleTrans{
			{
				LanguageCode: lang,
				Title:        title,
				Content:      content,
			},
		},
		Stats: model.ArticleStats{
			ViewCount: 0,
		},
	}

	err := s.repo.CreateArticle(ctx, article)
	return article, err
}

func (s *NewsService) GetAllArticles(ctx context.Context) ([]model.Article, error) {
	return s.repo.ListArticles(ctx, nil)
}

func (s *NewsService) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *NewsService) CreateCategory(ctx context.Context, name string) error {
	catSlug := slug.Make(name)
	category := &model.Category{
		Name: name,
		Slug: catSlug,
	}
	return s.repo.CreateCategory(ctx, category)
}
