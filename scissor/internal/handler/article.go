package handler

import (
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful/v3"
	"github.com/gorexlv/cabinet/scissor/internal/domain"
	"github.com/gorexlv/cabinet/scissor/internal/service"
)

type ArticleHandler struct {
	articleService *service.ArticleService
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

func (h *ArticleHandler) Register(ws *restful.WebService) {
	ws.Route(ws.POST("/articles").To(h.Create).
		Doc("创建文章").
		Reads(domain.CreateArticleRequest{}).
		Returns(200, "OK", domain.Article{}).
		Returns(400, "Bad Request", nil))

	ws.Route(ws.GET("/articles/{id}").To(h.GetByID).
		Doc("获取文章").
		Param(ws.PathParameter("id", "文章ID").DataType("integer")).
		Returns(200, "OK", domain.Article{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.GET("/articles").To(h.List).
		Doc("获取文章列表").
		Param(ws.QueryParameter("offset", "偏移量").DataType("integer").DefaultValue("0")).
		Param(ws.QueryParameter("limit", "限制数量").DataType("integer").DefaultValue("10")).
		Returns(200, "OK", []domain.Article{}))

	ws.Route(ws.GET("/articles/search").To(h.Search).
		Doc("搜索文章").
		Param(ws.QueryParameter("keyword", "搜索关键词")).
		Returns(200, "OK", []domain.Article{}))

	ws.Route(ws.GET("/users/{userId}/articles").To(h.GetByUserID).
		Doc("获取用户文章").
		Param(ws.PathParameter("userId", "用户ID").DataType("integer")).
		Returns(200, "OK", []domain.Article{}))
}

func (h *ArticleHandler) Create(req *restful.Request, resp *restful.Response) {
	var createReq domain.CreateArticleRequest
	if err := req.ReadEntity(&createReq); err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{
			"error": "无效的请求数据",
		})
		return
	}

	article := &domain.Article{
		Title:       createReq.Title,
		Content:     createReq.Content,
		URL:         createReq.URL,
		Author:      createReq.Author,
		Source:      createReq.Source,
		Summary:     createReq.Summary,
		Tags:        createReq.Tags,
		PublishedAt: createReq.PublishedAt,
		UserID:      createReq.UserID,
	}

	createdArticle, err := h.articleService.Create(req.Request.Context(), article)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	resp.WriteHeaderAndEntity(http.StatusCreated, createdArticle)
}

func (h *ArticleHandler) GetByID(req *restful.Request, resp *restful.Response) {
	id, err := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{
			"error": "无效的文章ID",
		})
		return
	}

	article, err := h.articleService.GetByID(req.Request.Context(), uint(id))
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusNotFound, map[string]string{
			"error": "文章不存在",
		})
		return
	}

	resp.WriteEntity(article)
}

func (h *ArticleHandler) List(req *restful.Request, resp *restful.Response) {
	offset, _ := strconv.Atoi(req.QueryParameter("offset"))
	limit, _ := strconv.Atoi(req.QueryParameter("limit"))

	articles, err := h.articleService.List(req.Request.Context(), offset, limit)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	resp.WriteEntity(articles)
}

func (h *ArticleHandler) Search(req *restful.Request, resp *restful.Response) {
	keyword := req.QueryParameter("keyword")
	if keyword == "" {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{
			"error": "搜索关键词不能为空",
		})
		return
	}

	articles, err := h.articleService.Search(req.Request.Context(), keyword)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	resp.WriteEntity(articles)
}

func (h *ArticleHandler) GetByUserID(req *restful.Request, resp *restful.Response) {
	userID, err := strconv.ParseUint(req.PathParameter("userId"), 10, 32)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{
			"error": "无效的用户ID",
		})
		return
	}

	articles, err := h.articleService.GetByUserID(req.Request.Context(), uint(userID))
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	resp.WriteEntity(articles)
}

func (h *ArticleHandler) GetByURL(req *restful.Request, resp *restful.Response) {
	url := req.QueryParameter("url")
	if url == "" {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{
			"error": "URL不能为空",
		})
		return
	}

	article, err := h.articleService.GetByURL(req.Request.Context(), url)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusNotFound, map[string]string{
			"error": "文章不存在",
		})
		return
	}

	resp.WriteEntity(article)
}

func (h *ArticleHandler) Update(req *restful.Request, resp *restful.Response) {
	id, err := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{
			"error": "无效的文章ID",
		})
		return
	}

	var updateReq domain.CreateArticleRequest
	if err := req.ReadEntity(&updateReq); err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{
			"error": "无效的请求数据",
		})
		return
	}

	article := &domain.Article{
		Title:       updateReq.Title,
		Content:     updateReq.Content,
		URL:         updateReq.URL,
		Author:      updateReq.Author,
		Source:      updateReq.Source,
		Summary:     updateReq.Summary,
		Tags:        updateReq.Tags,
		PublishedAt: updateReq.PublishedAt,
		UserID:      updateReq.UserID,
	}

	updatedArticle, err := h.articleService.Update(req.Request.Context(), uint(id), article)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	resp.WriteEntity(updatedArticle)
}

func (h *ArticleHandler) Delete(req *restful.Request, resp *restful.Response) {
	id, err := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{
			"error": "无效的文章ID",
		})
		return
	}

	if err := h.articleService.Delete(req.Request.Context(), uint(id)); err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}
