<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-gray-900">最近浏览</h1>
      <div class="flex gap-2">
        <button
          @click="clearHistory"
          class="px-4 py-2 text-sm text-gray-600 hover:text-gray-900"
        >
          清空历史
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="article in recentArticles"
        :key="article.id"
        class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow"
      >
        <div class="flex justify-between items-start">
          <div class="flex-1">
            <h3 class="text-lg font-medium text-gray-900">
              <a
                :href="article.url"
                target="_blank"
                class="hover:text-indigo-600"
              >
                {{ article.title }}
              </a>
            </h3>
            <div class="mt-2 flex items-center text-sm text-gray-500">
              <span>{{ formatDate(article.created_at) }}</span>
              <span class="mx-2">·</span>
              <span>{{ article.content.length }} 字</span>
            </div>
            <div class="mt-2 flex flex-wrap gap-2">
              <span
                v-for="tag in article.tags"
                :key="tag"
                class="px-2 py-1 text-xs font-medium bg-indigo-100 text-indigo-800 rounded-full"
              >
                {{ tag }}
              </span>
            </div>
          </div>
          <div class="flex gap-2">
            <button
              @click="handleFavorite(article)"
              class="px-3 py-1 text-sm text-indigo-600 hover:text-indigo-800"
            >
              <svg
                class="h-5 w-5"
                :class="{ 'text-red-500': article.is_favorite }"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
                />
              </svg>
            </button>
            <button
              @click="handleDelete(article)"
              class="px-3 py-1 text-sm text-red-600 hover:text-red-800"
            >
              <svg
                class="h-5 w-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { formatDate } from '@/utils/date'
import { useArticleStore } from '@/stores/article'

const articleStore = useArticleStore()
const recentArticles = ref([])

const fetchRecentArticles = async () => {
  try {
    const response = await articleStore.getRecentArticles()
    recentArticles.value = response
  } catch (error) {
    console.error('获取最近文章失败:', error)
  }
}

const handleFavorite = async (article) => {
  try {
    await articleStore.toggleFavorite(article.id)
    article.is_favorite = !article.is_favorite
  } catch (error) {
    console.error('收藏操作失败:', error)
  }
}

const handleDelete = async (article) => {
  if (confirm('确定要删除这篇文章吗？')) {
    try {
      await articleStore.deleteArticle(article.id)
      recentArticles.value = recentArticles.value.filter(
        (a) => a.id !== article.id
      )
    } catch (error) {
      console.error('删除文章失败:', error)
    }
  }
}

const clearHistory = async () => {
  if (confirm('确定要清空浏览历史吗？')) {
    try {
      await articleStore.clearHistory()
      recentArticles.value = []
    } catch (error) {
      console.error('清空历史失败:', error)
    }
  }
}

onMounted(() => {
  fetchRecentArticles()
})
</script> 