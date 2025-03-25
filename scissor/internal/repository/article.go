package repository

import (
	"context"
	"errors"

	"github.com/gorexlv/cabinet/scissor/pkg/ent"
	"github.com/gorexlv/cabinet/scissor/pkg/ent/article"
	"github.com/gorexlv/cabinet/scissor/pkg/ent/user"
)

var ErrNotFound = errors.New("记录未找到")

type ArticleRepository struct {
	client *ent.Client
}

func NewArticleRepository(client *ent.Client) *ArticleRepository {
	return &ArticleRepository{client: client}
}

func (r *ArticleRepository) Create(ctx context.Context, article *ent.Article) (*ent.Article, error) {
	return r.client.Article.Create().
		SetTitle(article.Title).
		SetContent(article.Content).
		SetURL(article.URL).
		SetAuthor(article.Author).
		SetSource(article.Source).
		SetSummary(article.Summary).
		SetTags(article.Tags).
		SetPublishedAt(article.PublishedAt).
		SetUserID(article.Edges.User.ID).
		Save(ctx)
}

func (r *ArticleRepository) FindByID(ctx context.Context, id int) (*ent.Article, error) {
	article, err := r.client.Article.Get(ctx, uint(id))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return article, nil
}

func (r *ArticleRepository) FindByURL(ctx context.Context, url string) (*ent.Article, error) {
	article, err := r.client.Article.Query().
		Where(article.URL(url)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return article, nil
}

func (r *ArticleRepository) FindByUserID(ctx context.Context, userID int) ([]*ent.Article, error) {
	return r.client.Article.Query().
		Where(article.HasUserWith(user.ID(userID))).
		All(ctx)
}

func (r *ArticleRepository) List(ctx context.Context, page, pageSize int) ([]*ent.Article, error) {
	offset := (page - 1) * pageSize
	return r.client.Article.Query().
		Offset(offset).
		Limit(pageSize).
		Order(ent.Desc(article.FieldPublishedAt)).
		All(ctx)
}

func (r *ArticleRepository) Search(ctx context.Context, keyword string) ([]*ent.Article, error) {
	return r.client.Article.Query().
		Where(
			article.Or(
				article.TitleContains(keyword),
				article.ContentContains(keyword),
				article.AuthorContains(keyword),
			),
		).
		Order(ent.Desc(article.FieldPublishedAt)).
		All(ctx)
}

func (r *ArticleRepository) Update(ctx context.Context, id int, article *ent.Article) (*ent.Article, error) {
	return r.client.Article.UpdateOneID(uint(id)).
		SetTitle(article.Title).
		SetContent(article.Content).
		SetURL(article.URL).
		SetAuthor(article.Author).
		SetSource(article.Source).
		SetSummary(article.Summary).
		SetTags(article.Tags).
		SetPublishedAt(article.PublishedAt).
		SetUserID(article.Edges.User.ID).
		Save(ctx)
}

func (r *ArticleRepository) Delete(ctx context.Context, id int) error {
	return r.client.Article.DeleteOneID(uint(id)).Exec(ctx)
}
