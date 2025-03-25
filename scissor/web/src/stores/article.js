import axios from 'axios'
import { defineStore } from 'pinia'

export const useArticleStore = defineStore('article', {
  state: () => ({
    articles: [],
    tags: [],
    loading: false,
    error: null,
  }),

  actions: {
    async getArticles(params) {
      this.loading = true
      try {
        const response = await axios.get('/api/articles', { params })
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async getRecentArticles() {
      try {
        const response = await axios.get('/api/articles/recent')
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async getFavoriteArticles() {
      try {
        const response = await axios.get('/api/articles/favorites')
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async getTags() {
      try {
        const response = await axios.get('/api/articles/tags')
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async updateArticle(id, data) {
      try {
        const response = await axios.put(`/api/articles/${id}`, data)
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async deleteArticle(id) {
      try {
        await axios.delete(`/api/articles/${id}`)
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async toggleFavorite(id) {
      try {
        const response = await axios.post(`/api/articles/${id}/favorite`)
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async clearHistory() {
      try {
        await axios.delete('/api/articles/history')
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async clearFavorites() {
      try {
        await axios.delete('/api/articles/favorites')
      } catch (error) {
        this.error = error.message
        throw error
      }
    },

    async searchArticles(query) {
      try {
        const response = await axios.get('/api/articles/search', {
          params: { q: query }
        })
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      }
    }
  },
}) 