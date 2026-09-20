<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, ArrowRight, Mail } from 'lucide-vue-next'
import BrandMark from '../components/BrandMark.vue'
import { api } from '../services/api'

const router = useRouter()
const name = ref('')
const email = ref('')
const password = ref('')
const code = ref('')
const loading = ref(false)
const sendingCode = ref(false)
const countdown = ref(0)
const error = ref('')
const message = ref('')
let timer: number | undefined

async function sendCode() {
  error.value = ''
  message.value = ''

  if (!name.value.trim() || !email.value.trim()) {
    error.value = '请先填写昵称和邮箱'
    return
  }

  sendingCode.value = true

  try {
    const result = await api.sendCode(name.value.trim(), email.value.trim())
    message.value = result.msg
    countdown.value = 60
    timer = window.setInterval(() => {
      countdown.value -= 1
      if (countdown.value <= 0 && timer) {
        window.clearInterval(timer)
      }
    }, 1000)
  } catch (err) {
    error.value = err instanceof Error ? err.message : '验证码发送失败'
  } finally {
    sendingCode.value = false
  }
}

async function submit() {
  error.value = ''
  message.value = ''
  loading.value = true

  try {
    await api.register(
      name.value.trim(),
      email.value.trim(),
      password.value,
      code.value.trim(),
    )
    await router.push({ name: 'login', query: { registered: '1' } })
  } catch (err) {
    error.value = err instanceof Error ? err.message : '注册失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

onBeforeUnmount(() => {
  if (timer) window.clearInterval(timer)
})
</script>

<template>
  <main class="auth-page">
    <header class="auth-header">
      <BrandMark />
      <p>
        已有账号？
        <RouterLink to="/login">登录</RouterLink>
      </p>
    </header>

    <section class="auth-panel" aria-labelledby="register-title">
      <RouterLink class="back-link" to="/login">
        <ArrowLeft :size="16" />
        返回登录
      </RouterLink>

      <div class="auth-intro">
        <p class="eyebrow">创建账号</p>
        <h1 id="register-title">加入 BirdNest</h1>
        <p>只需要一分钟。</p>
      </div>

      <form class="form-stack" @submit.prevent="submit">
        <label class="field">
          <span>昵称</span>
          <input
            v-model="name"
            type="text"
            name="name"
            autocomplete="nickname"
            placeholder="你的昵称"
            required
          />
        </label>

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
          <input
            v-model="password"
            type="password"
            name="password"
            autocomplete="new-password"
            minlength="6"
            placeholder="至少 6 位"
            required
          />
        </label>

        <label class="field">
          <span>邮箱验证码</span>
          <span class="code-row">
            <input
              v-model="code"
              type="text"
              name="code"
              inputmode="numeric"
              autocomplete="one-time-code"
              placeholder="输入验证码"
              required
            />
            <button
              class="secondary-button code-button"
              type="button"
              :disabled="sendingCode || countdown > 0"
              @click="sendCode"
            >
              <Mail v-if="countdown === 0" :size="17" />
              {{ countdown > 0 ? `${countdown}s` : sendingCode ? '发送中' : '获取验证码' }}
            </button>
          </span>
        </label>

        <p v-if="message" class="form-message">{{ message }}</p>
        <p v-if="error" class="form-error" role="alert">{{ error }}</p>

        <button class="primary-button" type="submit" :disabled="loading">
          <span>{{ loading ? '注册中...' : '创建账号' }}</span>
          <ArrowRight v-if="!loading" :size="18" />
        </button>
      </form>
    </section>
  </main>
</template>
