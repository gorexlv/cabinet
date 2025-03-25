package domain

import (
	"time"
)

type Article struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	URL         string    `json:"url"`
	Author      string    `json:"author"`
	Source      string    `json:"source"`
	Summary     string    `json:"summary"`
	Tags        []string  `json:"tags"`
	PublishedAt time.Time `json:"published_at"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateArticleRequest struct {
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	URL         string    `json:"url"`
	Author      string    `json:"author"`
	Source      string    `json:"source"`
	Summary     string    `json:"summary"`
	Tags        []string  `json:"tags"`
	PublishedAt time.Time `json:"published_at"`
	UserID      uint      `json:"user_id"`
}

type ArticleResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	URL         string    `json:"url"`
	Author      string    `json:"author"`
	Source      string    `json:"source"`
	Summary     string    `json:"summary"`
	Tags        []string  `json:"tags"`
	PublishedAt time.Time `json:"published_at"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
