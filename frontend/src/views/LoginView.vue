<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowRight, Eye, EyeOff } from 'lucide-vue-next'
import BrandMark from '../components/BrandMark.vue'
import { api, auth } from '../services/api'

const route = useRoute()
const router = useRouter()
const email = ref('')
const password = ref('')
const showPassword = ref(false)
const loading = ref(false)
const error = ref('')

const registered = computed(() => route.query.registered === '1')

async function submit() {
  error.value = ''
  loading.value = true

  try {
    const result = await api.login(email.value.trim(), password.value)
    auth.save(result.token)

    const redirect =
      typeof route.query.redirect === 'string' ? route.query.redirect : '/publish'
    await router.push(redirect)
  } catch (err) {
    error.value = err instanceof Error ? err.message : '登录失败，请稍后重试'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="auth-page">
    <header class="auth-header">
      <BrandMark />
      <p>
        还没有账号？
        <RouterLink to="/register">注册</RouterLink>
      </p>
    </header>

    <section class="auth-panel" aria-labelledby="login-title">
      <div class="auth-intro">
        <p class="eyebrow">欢迎回来</p>
        <h1 id="login-title">登录 BirdNest</h1>
        <p>继续记录你的想法。</p>
      </div>

      <p v-if="registered" class="notice notice--success">注册成功，请登录。</p>

      <form class="form-stack" @submit.prevent="submit">
        <label class="field">
          <span>邮箱</span>
          <input
            v-model="email"
            type="email"
            name="email"
            autocomplete="email"
            placeholder="name@example.com"
            required
          />
        </label>

        <label class="field">
          <span>密码</span>
          <span class="input-with-action">
            <input
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              name="password"
              autocomplete="current-password"
              placeholder="输入密码"
              required
            />
            <button
              class="icon-button input-action"
              type="button"
              :title="showPassword ? '隐藏密码' : '显示密码'"
              @click="showPassword = !showPassword"
            >
              <EyeOff v-if="showPassword" :size="18" />
              <Eye v-else :size="18" />
            </button>
          </span>
        </label>

        <p v-if="error" class="form-error" role="alert">{{ error }}</p>

        <button class="primary-button" type="submit" :disabled="loading">
          <span>{{ loading ? '登录中...' : '登录' }}</span>
          <ArrowRight v-if="!loading" :size="18" />
        </button>
      </form>
    </section>
  </main>
</template>
