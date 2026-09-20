<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Edit3, Filter, Plus, RefreshCw, Search, Trash2, X } from 'lucide-vue-next'
import WorkspaceHeader from '../components/WorkspaceHeader.vue'
import { api, type Passage } from '../services/api'

const router = useRouter()
const passages = ref<Passage[]>([])
const listLoading = ref(false)
const filterTag = ref('')
const activeTag = ref('')
const error = ref('')
const success = ref('')

function getPassageTags(passage: Passage) {
  return (passage.tags ?? []).map((tag) => tag.tagName).filter(Boolean)
}

function formatDate(date?: string) {
  if (!date) return ''
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }).format(new Date(date))
}

async function loadPassages(tag = activeTag.value) {
  listLoading.value = true
  error.value = ''

  try {
    const normalizedTag = tag.trim()
    const result = normalizedTag
      ? await api.getPassagesByTag(normalizedTag)
      : await api.getAllPassages()
    activeTag.value = normalizedTag
    filterTag.value = normalizedTag
    passages.value = result.data ?? []
  } catch (err) {
    error.value = err instanceof Error ? err.message : '文章加载失败'
  } finally {
    listLoading.value = false
  }
}

async function filterByTag(tag = filterTag.value) {
  await loadPassages(tag)
}

async function clearTagFilter() {
  await loadPassages('')
}

async function editPassage(passage: Passage) {
  await router.push({ name: 'publish', query: { edit: String(passage.ID) } })
}

async function removePassage(passage: Passage) {
  if (!window.confirm('确定要删除这篇文章吗？')) return

  error.value = ''
  success.value = ''

  try {
    const result = await api.deletePassage(passage.ID)
    passages.value = passages.value.filter((item) => item.ID !== passage.ID)
    success.value = result.msg
  } catch (err) {
    error.value = err instanceof Error ? err.message : '删除失败，请稍后重试'
  }
}

onMounted(() => loadPassages())
</script>

<template>
  <div class="workspace">
    <WorkspaceHeader />

    <main class="home-page">
      <div class="home-heading">
        <div>
          <p class="eyebrow">文章管理</p>
          <h1>{{ activeTag ? `标签：${activeTag}` : '我的文章' }}</h1>
          <p>查看、筛选和管理你发布的内容。</p>
        </div>
        <RouterLink class="primary-button new-passage-button" to="/publish">
          <Plus :size="17" />
          发布文章
        </RouterLink>
      </div>

      <form class="tag-filter" @submit.prevent="filterByTag()">
        <div class="tag-filter-input">
          <Search :size="17" />
          <input
            v-model="filterTag"
            type="search"
            maxlength="20"
            placeholder="按标签筛选，例如：Vue"
            aria-label="按标签筛选"
          />
        </div>
        <button class="secondary-button filter-button" type="submit" :disabled="listLoading">
          <Filter :size="16" />
          筛选
        </button>
        <button
          v-if="activeTag"
          class="icon-button"
          type="button"
          title="清除标签筛选"
          @click="clearTagFilter"
        >
          <X :size="18" />
        </button>
        <button
          class="icon-button"
          type="button"
          title="刷新文章"
          :disabled="listLoading"
          @click="loadPassages()"
        >
          <RefreshCw :size="18" :class="{ spinning: listLoading }" />
        </button>
      </form>

      <p v-if="success" class="page-message form-message">{{ success }}</p>
      <p v-if="error" class="page-message form-error" role="alert">{{ error }}</p>

      <div v-if="listLoading" class="empty-state">正在加载文章...</div>
      <div v-else-if="passages.length === 0" class="empty-state">
        {{ activeTag ? `没有找到标签“${activeTag}”下的文章。` : '还没有文章，写下第一篇吧。' }}
      </div>
      <div v-else class="passage-list">
        <article v-for="passage in passages" :key="passage.ID" class="passage-item">
          <div class="passage-content">
            <p>{{ passage.content }}</p>
            <div class="passage-meta">
              <span>{{ formatDate(passage.CreatedAt) }}</span>
              <button
                v-for="tag in getPassageTags(passage)"
                :key="`${passage.ID}-${tag}`"
                class="mini-tag"
                type="button"
                :title="`筛选标签：${tag}`"
                @click="filterByTag(tag)"
              >
                {{ tag }}
              </button>
            </div>
          </div>
          <div class="passage-actions">
            <button
              class="icon-button"
              type="button"
              title="编辑文章"
              @click="editPassage(passage)"
            >
              <Edit3 :size="17" />
            </button>
            <button
              class="icon-button danger-button"
              type="button"
              title="删除文章"
              @click="removePassage(passage)"
            >
              <Trash2 :size="17" />
            </button>
          </div>
        </article>
      </div>
    </main>
  </div>
</template>
