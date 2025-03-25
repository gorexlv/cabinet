<template>
  <div class="max-w-3xl mx-auto">
    <div class="bg-white shadow rounded-lg">
      <!-- 个人信息 -->
      <div class="px-4 py-5 sm:px-6">
        <h3 class="text-lg leading-6 font-medium text-gray-900">
          个人信息
        </h3>
      </div>
      <div class="border-t border-gray-200 px-4 py-5 sm:px-6">
        <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
          <div class="sm:col-span-3">
            <label
              for="username"
              class="block text-sm font-medium text-gray-700"
            >
              用户名
            </label>
            <div class="mt-1">
              <input
                type="text"
                name="username"
                id="username"
                v-model="form.username"
                class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
              />
            </div>
          </div>

          <div class="sm:col-span-3">
            <label
              for="email"
              class="block text-sm font-medium text-gray-700"
            >
              邮箱
            </label>
            <div class="mt-1">
              <input
                type="email"
                name="email"
                id="email"
                v-model="form.email"
                class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
              />
            </div>
          </div>

          <div class="sm:col-span-6">
            <label
              for="avatar"
              class="block text-sm font-medium text-gray-700"
            >
              头像
            </label>
            <div class="mt-1 flex items-center">
              <img
                :src="form.avatar || '/default-avatar.png'"
                class="h-12 w-12 rounded-full"
                alt="用户头像"
              />
              <input
                type="file"
                name="avatar"
                id="avatar"
                accept="image/*"
                class="ml-4"
                @change="handleAvatarChange"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 修改密码 -->
      <div class="px-4 py-5 sm:px-6">
        <h3 class="text-lg leading-6 font-medium text-gray-900">
          修改密码
        </h3>
      </div>
      <div class="border-t border-gray-200 px-4 py-5 sm:px-6">
        <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
          <div class="sm:col-span-3">
            <label
              for="current_password"
              class="block text-sm font-medium text-gray-700"
            >
              当前密码
            </label>
            <div class="mt-1">
              <input
                type="password"
                name="current_password"
                id="current_password"
                v-model="form.current_password"
                class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
              />
            </div>
          </div>

          <div class="sm:col-span-3">
            <label
              for="new_password"
              class="block text-sm font-medium text-gray-700"
            >
              新密码
            </label>
            <div class="mt-1">
              <input
                type="password"
                name="new_password"
                id="new_password"
                v-model="form.new_password"
                class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
              />
            </div>
          </div>

          <div class="sm:col-span-3">
            <label
              for="confirm_password"
              class="block text-sm font-medium text-gray-700"
            >
              确认新密码
            </label>
            <div class="mt-1">
              <input
                type="password"
                name="confirm_password"
                id="confirm_password"
                v-model="form.confirm_password"
                class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
              />
            </div>
          </div>
        </div>
      </div>

      <div class="px-4 py-3 bg-gray-50 text-right sm:px-6">
        <button
          type="button"
          @click="handleSubmit"
          :disabled="loading"
          class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          {{ loading ? '保存中...' : '保存' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const loading = ref(false)

const form = ref({
  username: '',
  email: '',
  avatar: '',
  current_password: '',
  new_password: '',
  confirm_password: '',
})

const handleAvatarChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      form.value.avatar = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

const handleSubmit = async () => {
  if (form.value.new_password && form.value.new_password !== form.value.confirm_password) {
    alert('两次输入的密码不一致')
    return
  }

  loading.value = true
  try {
    await userStore.updateProfile({
      username: form.value.username,
      email: form.value.email,
      avatar: form.value.avatar,
      current_password: form.value.current_password,
      new_password: form.value.new_password,
    })
    alert('保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    alert('保存失败')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const user = await userStore.fetchUser()
  if (user) {
    form.value = {
      ...form.value,
      username: user.username,
      email: user.email,
      avatar: user.avatar,
    }
  }
})
</script> 