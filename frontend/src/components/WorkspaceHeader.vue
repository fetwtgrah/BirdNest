<script setup lang="ts">
import { FileText, LogOut, PenLine } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import BrandMark from './BrandMark.vue'
import { auth } from '../services/api'

const router = useRouter()

async function logout() {
  auth.clear()
  await router.push('/login')
}
</script>

<template>
  <header class="workspace-header">
    <BrandMark />

    <div class="workspace-actions">
      <nav class="main-nav" aria-label="主要导航">
        <RouterLink to="/home">
          <FileText :size="17" />
          <span>文章</span>
        </RouterLink>
        <RouterLink to="/publish">
          <PenLine :size="17" />
          <span>发布</span>
        </RouterLink>
      </nav>
      <button class="icon-button" type="button" title="退出登录" @click="logout">
        <LogOut :size="19" />
      </button>
    </div>
  </header>
</template>
