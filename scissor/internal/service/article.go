package service

import (
	"context"
	"errors"

	"github.com/gorexlv/cabinet/scissor/internal/domain"
	"github.com/gorexlv/cabinet/scissor/internal/repository"
	"github.com/gorexlv/cabinet/scissor/pkg/ent"
)

type ArticleService struct {
	repo *repository.ArticleRepository
}

func NewArticleService(repo *repository.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) Create(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	// 检查URL是否已存在
	existing, err := s.repo.FindByURL(ctx, article.URL)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("文章URL已存在")
	}

	// 创建文章
	entArticle, err := s.repo.Create(ctx, &ent.Article{
		Title:       article.Title,
		Content:     article.Content,
		URL:         article.URL,
		Author:      article.Author,
		Source:      article.Source,
		Summary:     article.Summary,
		Tags:        article.Tags,
		PublishedAt: article.PublishedAt,
		Edges: ent.ArticleEdges{
			User: &ent.User{
				ID: int(article.UserID),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &domain.Article{
		ID:          uint(entArticle.ID),
		Title:       entArticle.Title,
		Content:     entArticle.Content,
		URL:         entArticle.URL,
		Author:      entArticle.Author,
		Source:      entArticle.Source,
		Summary:     entArticle.Summary,
		Tags:        entArticle.Tags,
		PublishedAt: entArticle.PublishedAt,
		UserID:      uint(entArticle.Edges.User.ID),
		CreatedAt:   entArticle.CreatedAt,
		UpdatedAt:   entArticle.UpdatedAt,
	}, nil
}

func (s *ArticleService) GetByID(ctx context.Context, id uint) (*domain.Article, error) {
	article, err := s.repo.FindByID(ctx, int(id))
	if err != nil {
		return nil, err
	}

	return &domain.Article{
		ID:          uint(article.ID),
		Title:       article.Title,
		Content:     article.Content,
		URL:         article.URL,
		Author:      article.Author,
		Source:      article.Source,
		Summary:     article.Summary,
		Tags:        article.Tags,
		PublishedAt: article.PublishedAt,
		UserID:      uint(article.Edges.User.ID),
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}, nil
}

func (s *ArticleService) GetByURL(ctx context.Context, url string) (*domain.Article, error) {
	article, err := s.repo.FindByURL(ctx, url)
	if err != nil {
		return nil, err
	}

	return &domain.Article{
		ID:          uint(article.ID),
		Title:       article.Title,
		Content:     article.Content,
		URL:         article.URL,
		Author:      article.Author,
		Source:      article.Source,
		Summary:     article.Summary,
		Tags:        article.Tags,
		PublishedAt: article.PublishedAt,
		UserID:      uint(article.Edges.User.ID),
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}, nil
}

func (s *ArticleService) List(ctx context.Context, page, pageSize int) ([]*domain.Article, error) {
	articles, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Article, len(articles))
	for i, article := range articles {
		result[i] = &domain.Article{
			ID:          uint(article.ID),
			Title:       article.Title,
			Content:     article.Content,
			URL:         article.URL,
			Author:      article.Author,
			Source:      article.Source,
			Summary:     article.Summary,
			Tags:        article.Tags,
			PublishedAt: article.PublishedAt,
			UserID:      uint(article.Edges.User.ID),
			CreatedAt:   article.CreatedAt,
			UpdatedAt:   article.UpdatedAt,
		}
	}

	return result, nil
}

func (s *ArticleService) Search(ctx context.Context, keyword string) ([]*domain.Article, error) {
	articles, err := s.repo.Search(ctx, keyword)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Article, len(articles))
	for i, article := range articles {
		result[i] = &domain.Article{
			ID:          uint(article.ID),
			Title:       article.Title,
			Content:     article.Content,
			URL:         article.URL,
			Author:      article.Author,
			Source:      article.Source,
			Summary:     article.Summary,
			Tags:        article.Tags,
			PublishedAt: article.PublishedAt,
			UserID:      uint(article.Edges.User.ID),
			CreatedAt:   article.CreatedAt,
			UpdatedAt:   article.UpdatedAt,
		}
	}

	return result, nil
}

func (s *ArticleService) GetByUserID(ctx context.Context, userID uint) ([]*domain.Article, error) {
	articles, err := s.repo.FindByUserID(ctx, int(userID))
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Article, len(articles))
	for i, article := range articles {
		result[i] = &domain.Article{
			ID:          uint(article.ID),
			Title:       article.Title,
			Content:     article.Content,
			URL:         article.URL,
			Author:      article.Author,
			Source:      article.Source,
			Summary:     article.Summary,
			Tags:        article.Tags,
			PublishedAt: article.PublishedAt,
			UserID:      uint(article.Edges.User.ID),
			CreatedAt:   article.CreatedAt,
			UpdatedAt:   article.UpdatedAt,
		}
	}

	return result, nil
}

func (s *ArticleService) Update(ctx context.Context, id uint, article *domain.Article) (*domain.Article, error) {
	// 检查文章是否存在
	existing, err := s.repo.FindByID(ctx, int(id))
	if err != nil {
		return nil, err
	}

	// 如果URL发生变化，检查新URL是否已存在
	if existing.URL != article.URL {
		urlExists, err := s.repo.FindByURL(ctx, article.URL)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
		if urlExists != nil {
			return nil, errors.New("文章URL已存在")
		}
	}

	// 更新文章
	entArticle, err := s.repo.Update(ctx, int(id), &ent.Article{
		Title:       article.Title,
		Content:     article.Content,
		URL:         article.URL,
		Author:      article.Author,
		Source:      article.Source,
		Summary:     article.Summary,
		Tags:        article.Tags,
		PublishedAt: article.PublishedAt,
		Edges: ent.ArticleEdges{
			User: &ent.User{
				ID: int(article.UserID),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &domain.Article{
		ID:          uint(entArticle.ID),
		Title:       entArticle.Title,
		Content:     entArticle.Content,
		URL:         entArticle.URL,
		Author:      entArticle.Author,
		Source:      entArticle.Source,
		Summary:     entArticle.Summary,
		Tags:        entArticle.Tags,
		PublishedAt: entArticle.PublishedAt,
		UserID:      uint(entArticle.Edges.User.ID),
		CreatedAt:   entArticle.CreatedAt,
		UpdatedAt:   entArticle.UpdatedAt,
	}, nil
}

func (s *ArticleService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, int(id))
}
