<template>
  <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-medium">微信登录</h3>
        <button @click="close" class="text-gray-500 hover:text-gray-700">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="text-center">
        <div v-if="loading" class="flex justify-center items-center h-48">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
        </div>
        <div v-else-if="qrcodeUrl" class="space-y-4">
          <img :src="qrcodeUrl" alt="微信登录二维码" class="mx-auto w-48 h-48" />
          <p class="text-gray-600">请使用微信扫描二维码登录</p>
        </div>
        <div v-else class="text-red-500">
          获取二维码失败，请重试
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  }
})

const emit = defineEmits(['close'])
const userStore = useUserStore()
const router = useRouter()

const loading = ref(true)
const qrcodeUrl = ref('')
let checkTimer = null

const getQrCode = async () => {
  try {
    loading.value = true
    const response = await fetch('/api/users/wx-qrcode')
    const data = await response.json()
    if (data.qrcode_url) {
      qrcodeUrl.value = data.qrcode_url
      startCheckLogin()
    }
  } catch (error) {
    console.error('获取二维码失败:', error)
  } finally {
    loading.value = false
  }
}

const startCheckLogin = () => {
  checkTimer = setInterval(async () => {
    try {
      const response = await fetch('/api/users/wx-check-login')
      const data = await response.json()
      if (data.status === 'success') {
        clearInterval(checkTimer)
        await userStore.wxLogin(data.data)
        router.push('/')
        close()
      }
    } catch (error) {
      console.error('检查登录状态失败:', error)
    }
  }, 2000)
}

const close = () => {
  if (checkTimer) {
    clearInterval(checkTimer)
  }
  emit('close')
}

onMounted(() => {
  if (props.show) {
    getQrCode()
  }
})

onUnmounted(() => {
  if (checkTimer) {
    clearInterval(checkTimer)
  }
})
</script> 