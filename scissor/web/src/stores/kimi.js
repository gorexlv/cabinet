import axios from 'axios'
import { defineStore } from 'pinia'

export const useKimiStore = defineStore('kimi', {
  state: () => ({
    loading: false,
    error: null,
  }),

  actions: {
    async summarize(articleId) {
      this.loading = true
      try {
        const response = await axios.post(`/api/articles/${articleId}/summarize`)
        return response.data
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },
  },
}) 