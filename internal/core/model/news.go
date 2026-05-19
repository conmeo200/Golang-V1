package model

import (
	"time"

	"gorm.io/gorm"
)

// ArticleStatus defines the enum for article status
type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"
	ArticleStatusPending   ArticleStatus = "pending"
	ArticleStatusPublished ArticleStatus = "published"
	ArticleStatusArchived  ArticleStatus = "archived"
	ArticleStatusHidden    ArticleStatus = "hidden"
)


// NewsUser represents the author or staff
type NewsUser struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Username        string         `gorm:"unique;not null;size:50" json:"username"`
	Email           string         `gorm:"unique;not null;size:255" json:"email"`
	PasswordHash    string         `gorm:"not null" json:"-"`
	DisplayName     string         `gorm:"size:100" json:"display_name"`
	Bio             string         `json:"bio"`
	AvatarURL       string         `gorm:"size:255" json:"avatar_url"`
	Role            string         `gorm:"size:20;default:'author'" json:"role"`
	SecurityVersion int            `gorm:"default:1" json:"-"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	
	Articles        []Article      `gorm:"foreignKey:AuthorID" json:"articles,omitempty"`
}

// Category represents news category with hierarchy
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ParentID    *uint          `gorm:"index" json:"parent_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Slug        string         `gorm:"unique;not null;size:120" json:"slug"`
	Description string         `json:"description"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Articles    []Article      `gorm:"many2many:article_categories;" json:"articles,omitempty"`
}

// Article represents the core news post
type Article struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	AuthorID         uint           `gorm:"not null" json:"author_id"`
	Slug             string         `gorm:"unique;not null;size:255" json:"slug"`
	Status           ArticleStatus  `gorm:"type:varchar(20);default:'draft'" json:"status"`
	ImageFeatureURL  string         `gorm:"size:255" json:"image_feature_url"`
	PublishedAt      *time.Time     `gorm:"index" json:"published_at"`

	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	Author           NewsUser       `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Translations     []ArticleTrans `gorm:"foreignKey:ArticleID" json:"translations,omitempty"`
	Categories       []Category     `gorm:"many2many:article_categories;" json:"categories,omitempty"`
	Tags             []Tag          `gorm:"many2many:article_tags;" json:"tags,omitempty"`
	Stats            ArticleStats   `gorm:"foreignKey:ArticleID" json:"stats,omitempty"`
	Comments         []NewsComment  `gorm:"foreignKey:ArticleID" json:"comments,omitempty"`
}

// ArticleTrans handles multi-language content
type ArticleTrans struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	ArticleID       uint   `gorm:"not null;uniqueIndex:idx_art_lang" json:"article_id"`
	LanguageCode    string `gorm:"size:5;not null;uniqueIndex:idx_art_lang" json:"language_code"`
	Title           string `gorm:"size:255;not null" json:"title"`
	Excerpt         string `json:"excerpt"`
	Content         string `gorm:"type:text;not null" json:"content"`
	MetaTitle       string `gorm:"size:255" json:"meta_title"`
	MetaDescription string `json:"meta_description"`
}

// Tag represents many-to-many labels
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null;size:50" json:"name"`
	Slug      string    `gorm:"unique;not null;size:70" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	Articles  []Article `gorm:"many2many:article_tags;" json:"articles,omitempty"`
}

// NewsComment handles nested discussions
type NewsComment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ArticleID uint           `gorm:"not null;index" json:"article_id"`
	UserID    *uint          `gorm:"index" json:"user_id"`
	ParentID  *uint          `gorm:"index" json:"parent_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Status    string         `gorm:"size:20;default:'approved'" json:"status"`
	IPAddress string         `gorm:"size:45" json:"ip_address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User      *NewsUser      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Replies   []NewsComment  `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// ArticleStats keeps track of high-frequency counters
type ArticleStats struct {
	ArticleID    uint      `gorm:"primaryKey" json:"article_id"`
	ViewCount    int64     `gorm:"default:0" json:"view_count"`
	LikeCount    int64     `gorm:"default:0" json:"like_count"`
	ShareCount   int64     `gorm:"default:0" json:"share_count"`
	CommentCount int64     `gorm:"default:0" json:"comment_count"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ArticleViewLog for analytics tracking
type ArticleViewLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"index" json:"article_id"`
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	ViewedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"viewed_at"`
}

// ArticleVersion for edit history
type ArticleVersion struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ArticleID  uint      `gorm:"not null;index" json:"article_id"`
	AuthorID   uint      `json:"author_id"`
	Title      string    `gorm:"size:255" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	VersionNote string   `json:"version_note"`
	CreatedAt  time.Time `json:"created_at"`
}
