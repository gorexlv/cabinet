<template>
  <div class="space-y-4">
    <!-- 过滤条件 -->
    <div class="flex flex-wrap gap-4 p-4 bg-white rounded-lg shadow">
      <div class="flex-1 min-w-[200px]">
        <label class="block text-sm font-medium text-gray-700"
          >关键词搜索</label
        >
        <input
          v-model="filters.keyword"
          type="text"
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          placeholder="搜索标题、内容或作者"
        />
      </div>
      <div class="flex-1 min-w-[200px]">
        <label class="block text-sm font-medium text-gray-700">标签过滤</label>
        <select
          v-model="filters.tag"
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
        >
          <option value="">全部标签</option>
          <option v-for="tag in tags" :key="tag" :value="tag">{{ tag }}</option>
        </select>
      </div>
      <div class="flex-1 min-w-[200px]">
        <label class="block text-sm font-medium text-gray-700">时间范围</label>
        <div class="flex gap-2">
          <input
            v-model="filters.startDate"
            type="date"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          />
          <input
            v-model="filters.endDate"
            type="date"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          />
        </div>
      </div>
      <div class="flex items-end">
        <button
          @click="applyFilters"
          class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
        >
          应用过滤
        </button>
      </div>
    </div>

    <!-- 文章列表 -->
    <div class="space-y-4">
      <div
        v-for="article in articles"
        :key="article.id"
        class="p-4 bg-white rounded-lg shadow hover:shadow-md transition-shadow"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <h3
              @click="goToDetail(article.id)"
              class="text-lg font-medium text-indigo-600 hover:text-indigo-800 cursor-pointer"
            >
              {{ article.title }}
            </h3>
            <div class="mt-2 text-sm text-gray-500">
              <span>作者: {{ article.author }}</span>
              <span class="mx-2">|</span>
              <span>发布时间: {{ formatDate(article.published_at) }}</span>
              <span class="mx-2">|</span>
              <span>字数: {{ article.content.length }}</span>
            </div>
            <div class="mt-2 flex flex-wrap gap-2">
              <span
                v-for="tag in article.tags"
                :key="tag"
                class="px-2 py-1 text-sm bg-indigo-100 text-indigo-800 rounded-full"
              >
                {{ tag }}
              </span>
            </div>
          </div>
          <div class="flex items-center space-x-2">
            <button
              @click="handleSummarize(article)"
              class="p-2 text-gray-500 hover:text-indigo-600"
              :disabled="article.summarizing"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 10V3L4 14h7v7l9-11h-7z"
                />
              </svg>
            </button>
            <button
              @click="handleTag(article)"
              class="p-2 text-gray-500 hover:text-indigo-600"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
                />
              </svg>
            </button>
            <button
              @click="handleDelete(article)"
              class="p-2 text-gray-500 hover:text-red-600"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
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

    <!-- 分页 -->
    <div class="flex justify-center mt-4">
      <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
        <button
          v-for="page in totalPages"
          :key="page"
          @click="currentPage = page"
          :class="[
            'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
            currentPage === page
              ? 'z-10 bg-indigo-50 border-indigo-500 text-indigo-600'
              : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50',
          ]"
        >
          {{ page }}
        </button>
      </nav>
    </div>

    <!-- 打标模态框 -->
    <div
      v-if="showTagModal"
      class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center"
    >
      <div class="bg-white rounded-lg p-6 max-w-md w-full">
        <h3 class="text-lg font-medium text-gray-900 mb-4">为文章打标</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700">标签</label>
            <div class="mt-1 flex flex-wrap gap-2">
              <span
                v-for="tag in selectedArticle.tags"
                :key="tag"
                class="px-2 py-1 text-xs font-medium bg-indigo-100 text-indigo-800 rounded-full flex items-center"
              >
                {{ tag }}
                <button
                  @click="removeTag(tag)"
                  class="ml-1 text-indigo-600 hover:text-indigo-800"
                >
                  ×
                </button>
              </span>
            </div>
            <div class="mt-2 flex gap-2">
              <input
                v-model="newTag"
                type="text"
                class="flex-1 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                placeholder="输入新标签"
                @keyup.enter="addTag"
              />
              <button
                @click="addTag"
                class="px-3 py-1 text-sm text-indigo-600 hover:text-indigo-800"
              >
                添加
              </button>
            </div>
          </div>
          <div class="flex justify-end gap-2">
            <button
              @click="showTagModal = false"
              class="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-800"
            >
              取消
            </button>
            <button
              @click="saveTags"
              class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700"
            >
              保存
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
import { useKimiStore } from '@/stores/kimi'
import { useRouter } from 'vue-router'

const router = useRouter()
const articleStore = useArticleStore()
const kimiStore = useKimiStore()

const articles = ref([])
const tags = ref([])
const currentPage = ref(1)
const totalPages = ref(1)
const showTagModal = ref(false)
const selectedArticle = ref(null)
const newTag = ref('')

const filters = ref({
  keyword: '',
  tag: '',
  startDate: '',
  endDate: '',
})

// 获取文章列表
const fetchArticles = async () => {
  try {
    const response = await articleStore.getArticles({
      page: currentPage.value,
      ...filters.value,
    })
    articles.value = response.data
    totalPages.value = response.total_pages
  } catch (error) {
    console.error('获取文章列表失败:', error)
  }
}

// 获取所有标签
const fetchTags = async () => {
  try {
    const response = await articleStore.getTags()
    tags.value = response
  } catch (error) {
    console.error('获取标签列表失败:', error)
  }
}

// 应用过滤条件
const applyFilters = () => {
  currentPage.value = 1
  fetchArticles()
}

// 处理文章摘要
const handleSummarize = async (article) => {
  try {
    const response = await kimiStore.summarize(article.id)
    article.summary = response.summary
    article.tags = response.tags
    await articleStore.updateArticle(article.id, {
      summary: response.summary,
      tags: response.tags,
    })
  } catch (error) {
    console.error('生成摘要失败:', error)
  }
}

// 处理文章打标
const handleTag = (article) => {
  selectedArticle.value = { ...article }
  showTagModal.value = true
}

// 添加标签
const addTag = () => {
  if (newTag.value && !selectedArticle.value.tags.includes(newTag.value)) {
    selectedArticle.value.tags.push(newTag.value)
  }
  newTag.value = ''
}

// 移除标签
const removeTag = (tag) => {
  selectedArticle.value.tags = selectedArticle.value.tags.filter((t) => t !== tag)
}

// 保存标签
const saveTags = async () => {
  try {
    await articleStore.updateArticle(selectedArticle.value.id, {
      tags: selectedArticle.value.tags,
    })
    showTagModal.value = false
    fetchArticles()
  } catch (error) {
    console.error('保存标签失败:', error)
  }
}

// 处理文章删除
const handleDelete = async (article) => {
  if (confirm('确定要删除这篇文章吗？')) {
    try {
      await articleStore.deleteArticle(article.id)
      fetchArticles()
    } catch (error) {
      console.error('删除文章失败:', error)
    }
  }
}

const goToDetail = (id) => {
  router.push(`/articles/${id}`)
}

onMounted(() => {
  fetchArticles()
  fetchTags()
})
</script>
