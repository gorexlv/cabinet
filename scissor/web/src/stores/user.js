import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || null,
    loading: false,
    error: null
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    currentUser: (state) => state.user
  },

  actions: {
    async login(username, password) {
      try {
        this.loading = true
        this.error = null
        const response = await fetch('/api/users/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ username, password })
        })
        const data = await response.json()
        if (response.ok) {
          this.token = data.token
          this.user = data.user
          localStorage.setItem('token', data.token)
        } else {
          this.error = data.error || '登录失败'
        }
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    },

    async wxLogin(wxData) {
      try {
        this.loading = true
        this.error = null
        const response = await fetch('/api/users/wx-login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(wxData)
        })
        const data = await response.json()
        if (response.ok) {
          this.token = data.token
          this.user = data.user
          localStorage.setItem('token', data.token)
        } else {
          this.error = data.error || '微信登录失败'
        }
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    },

    logout() {
      this.user = null
      this.token = null
      localStorage.removeItem('token')
    },

    async fetchUserInfo() {
      if (!this.token) return

      try {
        this.loading = true
        this.error = null
        const response = await fetch('/api/users/me', {
          headers: {
            'Authorization': `Bearer ${this.token}`
          }
        })
        const data = await response.json()
        if (response.ok) {
          this.user = data
        } else {
          this.error = data.error || '获取用户信息失败'
        }
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    },

    async updateProfile(data) {
      try {
        const response = await fetch('/api/users/me', {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(data)
        })
        const updatedUser = await response.json()
        this.user = updatedUser
        return updatedUser
      } catch (error) {
        this.error = error.message
        throw error
      }
    }
  }
}) 