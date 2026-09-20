<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Check, Edit3, LogOut, Plus, RefreshCw, Send, Trash2, X } from 'lucide-vue-next'
import BrandMark from '../components/BrandMark.vue'
import { api, auth, type Passage } from '../services/api'

const router = useRouter()
const content = ref('')
const tagInput = ref('')
const tags = ref<string[]>([])
const loading = ref(false)
const error = ref('')
const success = ref('')
const passages = ref<Passage[]>([])
const listLoading = ref(false)
const editingId = ref<number | null>(null)
const editor = ref<HTMLTextAreaElement | null>(null)

const contentLength = computed(() => content.value.length)
const isEditing = computed(() => editingId.value !== null)
const canSubmit = computed(
  () => content.value.trim().length > 0 && contentLength.value <= 5000 && !loading.value,
)

function addTag() {
  const candidates = tagInput.value
    .split(/[,，]/)
    .map((tag) => tag.trim())
    .filter(Boolean)

  for (const tag of candidates) {
    if (tags.value.length >= 6) break
    if (tag.length <= 20 && !tags.value.includes(tag)) tags.value.push(tag)
  }

  tagInput.value = ''
}

function handleTagKeydown(event: KeyboardEvent) {
  if (event.isComposing) return
  if (event.key === 'Enter' || event.key === ',') {
    event.preventDefault()
    addTag()
  }
}

function removeTag(index: number) {
  tags.value.splice(index, 1)
}

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

async function loadPassages() {
  listLoading.value = true

  try {
    const result = await api.getAllPassages()
    passages.value = result.data ?? []
  } catch (err) {
    error.value = err instanceof Error ? err.message : '文章加载失败'
  } finally {
    listLoading.value = false
  }
}

function startEdit(passage: Passage) {
  editingId.value = passage.ID
  content.value = passage.content
  tags.value = getPassageTags(passage)
  tagInput.value = ''
  error.value = ''
  success.value = ''
  window.scrollTo({ top: 0, behavior: 'smooth' })
  nextTick(() => editor.value?.focus())
}

function cancelEdit() {
  editingId.value = null
  content.value = ''
  tags.value = []
  tagInput.value = ''
  error.value = ''
  success.value = ''
}

async function removePassage(passage: Passage) {
  if (!window.confirm('确定要删除这篇文章吗？')) return

  error.value = ''
  success.value = ''

  try {
    const result = await api.deletePassage(passage.ID)
    passages.value = passages.value.filter((item) => item.ID !== passage.ID)
    if (editingId.value === passage.ID) cancelEdit()
    success.value = result.msg
  } catch (err) {
    error.value = err instanceof Error ? err.message : '删除失败，请稍后重试'
  }
}

async function submit() {
  if (!canSubmit.value) return

  addTag()
  error.value = ''
  success.value = ''
  loading.value = true

  try {
    const result = isEditing.value
      ? await api.updatePassage(editingId.value!, content.value.trim(), tags.value)
      : await api.publish(content.value.trim(), tags.value)
    success.value = result.msg
    content.value = ''
    tags.value = []
    editingId.value = null
    await loadPassages()
    await nextTick()
    editor.value?.focus()
  } catch (err) {
    const message = err instanceof Error ? err.message : '发布失败，请稍后重试'
    error.value = message

    if (message.includes('权限')) {
      auth.clear()
      await router.push('/login')
    }
  } finally {
    loading.value = false
  }
}

async function logout() {
  auth.clear()
  await router.push('/login')
}

onMounted(loadPassages)
</script>

<template>
  <div class="workspace">
    <header class="workspace-header">
      <BrandMark />
      <div class="workspace-actions">
        <span class="status-dot">已登录</span>
        <button class="icon-button" type="button" title="退出登录" @click="logout">
          <LogOut :size="19" />
        </button>
      </div>
    </header>

    <main class="publish-page">
      <div class="page-heading">
        <p class="eyebrow">{{ isEditing ? '编辑文章' : '新文章' }}</p>
        <h1>{{ isEditing ? '修改你的文章' : '分享此刻的想法' }}</h1>
        <p>{{ isEditing ? '更新内容后保存修改。' : '写点值得留下的内容。' }}</p>
      </div>

      <form class="editor" @submit.prevent="submit">
        <label class="editor-field">
          <span class="editor-label">正文</span>
          <textarea
            ref="editor"
            v-model="content"
            maxlength="5000"
            placeholder="开始写作..."
            autofocus
            required
          />
        </label>

        <div class="editor-divider" />

        <div class="tag-field">
          <div class="editor-label-row">
            <span class="editor-label">标签</span>
            <span>最多 6 个</span>
          </div>

          <div v-if="tags.length" class="tag-list">
            <span v-for="(tag, index) in tags" :key="tag" class="tag-chip">
              {{ tag }}
              <button type="button" :title="`移除 ${tag}`" @click="removeTag(index)">
                <X :size="14" />
              </button>
            </span>
          </div>

          <div class="tag-input-row">
            <input
              v-model="tagInput"
              type="text"
              maxlength="20"
              placeholder="输入标签后按回车"
              :disabled="tags.length >= 6"
              @keydown="handleTagKeydown"
              @blur="addTag"
            />
            <button
              class="icon-button add-tag-button"
              type="button"
              title="添加标签"
              :disabled="!tagInput.trim() || tags.length >= 6"
              @mousedown.prevent
              @click="addTag"
            >
              <Plus :size="18" />
            </button>
          </div>
        </div>

        <div class="editor-footer">
          <div class="editor-feedback">
            <p v-if="success" class="form-message">
              <Check :size="16" />
              {{ success }}
            </p>
            <p v-else-if="error" class="form-error" role="alert">{{ error }}</p>
            <span v-else>{{ contentLength }} / 5000</span>
          </div>

          <button class="primary-button publish-button" type="submit" :disabled="!canSubmit">
            <span>{{ loading ? '保存中...' : isEditing ? '保存修改' : '发布文章' }}</span>
            <Send v-if="!loading" :size="17" />
          </button>
          <button
            v-if="isEditing"
            class="secondary-button cancel-edit-button"
            type="button"
            @click="cancelEdit"
          >
            取消编辑
          </button>
        </div>
      </form>

      <section class="passage-section" aria-labelledby="passages-title">
        <div class="section-heading">
          <div>
            <p class="eyebrow">内容管理</p>
            <h2 id="passages-title">我的文章</h2>
          </div>
          <button
            class="icon-button"
            type="button"
            title="刷新文章"
            :disabled="listLoading"
            @click="loadPassages"
          >
            <RefreshCw :size="18" :class="{ spinning: listLoading }" />
          </button>
        </div>

        <div v-if="listLoading" class="empty-state">正在加载文章...</div>
        <div v-else-if="passages.length === 0" class="empty-state">还没有文章，写下第一篇吧。</div>
        <div v-else class="passage-list">
          <article v-for="passage in passages" :key="passage.ID" class="passage-item">
            <div class="passage-content">
              <p>{{ passage.content }}</p>
              <div class="passage-meta">
                <span>{{ formatDate(passage.CreatedAt) }}</span>
                <span
                  v-for="tag in getPassageTags(passage)"
                  :key="`${passage.ID}-${tag}`"
                  class="mini-tag"
                >
                  {{ tag }}
                </span>
              </div>
            </div>
            <div class="passage-actions">
              <button
                class="icon-button"
                type="button"
                title="编辑文章"
                @click="startEdit(passage)"
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
      </section>
    </main>
  </div>
</template>
