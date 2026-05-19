-- ==========================================================================
-- High-Performance News Website Schema (Production Ready)
-- Database: PostgreSQL
-- Author: Senior Database Architect
-- ==========================================================================

-- 1. Extensions & Custom Types
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'article_status') THEN
        CREATE TYPE article_status AS ENUM ('draft', 'pending', 'published', 'archived', 'hidden');
    END IF;
END $$;

-- 2. Users Table (Authors/Staff)
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    bio TEXT,
    avatar_url VARCHAR(255),
    role VARCHAR(20) DEFAULT 'author',
    security_version INT DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 3. Categories Table (Hierarchical)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    parent_id INT REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(120) UNIQUE NOT NULL,
    description TEXT,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 4. Articles Table (Core Metadata)
CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY,
    author_id BIGINT NOT NULL REFERENCES users(id),
    slug VARCHAR(255) UNIQUE NOT NULL,
    status article_status DEFAULT 'draft',
    image_feature_url VARCHAR(255),
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT check_published_date CHECK (status != 'published' OR published_at IS NOT NULL)
);

-- 5. Article Translations (Localization)
CREATE TABLE IF NOT EXISTS article_translations (
    id BIGSERIAL PRIMARY KEY,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    language_code VARCHAR(5) NOT NULL,
    title VARCHAR(255) NOT NULL,
    excerpt TEXT,
    content TEXT NOT NULL,
    meta_title VARCHAR(255),
    meta_description TEXT,
    UNIQUE (article_id, language_code)
);

-- 6. Tags Table
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    slug VARCHAR(70) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 7. Junction Tables (Many-to-Many)
CREATE TABLE IF NOT EXISTS article_categories (
    article_id BIGINT REFERENCES articles(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, category_id)
);

CREATE TABLE IF NOT EXISTS article_tags (
    article_id BIGINT REFERENCES articles(id) ON DELETE CASCADE,
    tag_id INT REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, tag_id)
);

-- 8. Comments Table (Nested/Recursive)
CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    parent_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'approved',
    ip_address INET,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 9. Article Stats (1:1 Aggregated Table for Performance)
CREATE TABLE IF NOT EXISTS article_stats (
    article_id BIGINT PRIMARY KEY REFERENCES articles(id) ON DELETE CASCADE,
    view_count BIGINT DEFAULT 0,
    like_count BIGINT DEFAULT 0,
    share_count BIGINT DEFAULT 0,
    comment_count BIGINT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 10. Article View Logs (For Analytics)
CREATE TABLE IF NOT EXISTS article_views (
    id BIGSERIAL PRIMARY KEY,
    article_id BIGINT REFERENCES articles(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT,
    viewed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 11. Article Versions (History)
CREATE TABLE IF NOT EXISTS article_versions (
    id BIGSERIAL PRIMARY KEY,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    author_id BIGINT REFERENCES users(id),
    title VARCHAR(255),
    content TEXT,
    version_note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================================================
-- INDEXES FOR HIGH PERFORMANCE
-- ==========================================================================

-- SEO Search Optimization
CREATE INDEX IF NOT EXISTS idx_articles_slug ON articles(slug) WHERE deleted_at IS NULL;

-- List Listing Optimization (Home/Category Pages)
CREATE INDEX IF NOT EXISTS idx_articles_published ON articles(status, published_at DESC) 
WHERE status = 'published' AND deleted_at IS NULL;

-- Many-to-Many Optimization
CREATE INDEX IF NOT EXISTS idx_article_categories_cat_id ON article_categories(category_id);
CREATE INDEX IF NOT EXISTS idx_article_tags_tag_id ON article_tags(tag_id);

-- Hierarchy & Timeline Optimization
CREATE INDEX IF NOT EXISTS idx_comments_article_hierarchy ON comments(article_id, parent_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories(parent_id);

-- User Retrieval
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;
