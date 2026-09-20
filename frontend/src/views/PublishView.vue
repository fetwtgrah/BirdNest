<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Check, Plus, Send, X } from 'lucide-vue-next'
import WorkspaceHeader from '../components/WorkspaceHeader.vue'
import { api, auth } from '../services/api'

const route = useRoute()
const router = useRouter()
const content = ref('')
const tagInput = ref('')
const tags = ref<string[]>([])
const loading = ref(false)
const loadingPassage = ref(false)
const error = ref('')
const success = ref('')
const editingId = ref<number | null>(null)
const editor = ref<HTMLTextAreaElement | null>(null)

const contentLength = computed(() => content.value.length)
const isEditing = computed(() => editingId.value !== null)
const canSubmit = computed(
  () =>
    content.value.trim().length > 0 &&
    contentLength.value <= 5000 &&
    !loading.value &&
    !loadingPassage.value,
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

async function loadPassageForEdit(id: number) {
  loadingPassage.value = true
  error.value = ''

  try {
    const result = await api.getPassage(id)
    editingId.value = result.data.ID
    content.value = result.data.content
    tags.value = (result.data.tags ?? []).map((tag) => tag.tagName).filter(Boolean)
    await nextTick()
    editor.value?.focus()
  } catch (err) {
    error.value = err instanceof Error ? err.message : '文章加载失败'
  } finally {
    loadingPassage.value = false
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

    if (isEditing.value) {
      await router.push('/home')
      return
    }

    content.value = ''
    tags.value = []
    await nextTick()
    editor.value?.focus()
  } catch (err) {
    const message = err instanceof Error ? err.message : '保存失败，请稍后重试'
    error.value = message

    if (message.includes('权限')) {
      auth.clear()
      await router.push('/login')
    }
  } finally {
    loading.value = false
  }
}

async function cancelEdit() {
  await router.push('/home')
}

onMounted(() => {
  const editId = Number(route.query.edit)
  if (Number.isInteger(editId) && editId > 0) {
    loadPassageForEdit(editId)
  }
})
</script>

<template>
  <div class="workspace">
    <WorkspaceHeader />

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
            :placeholder="loadingPassage ? '正在加载文章...' : '开始写作...'"
            :disabled="loadingPassage"
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
              :disabled="tags.length >= 6 || loadingPassage"
              @keydown="handleTagKeydown"
              @blur="addTag"
            />
            <button
              class="icon-button add-tag-button"
              type="button"
              title="添加标签"
              :disabled="!tagInput.trim() || tags.length >= 6 || loadingPassage"
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

          <div class="editor-actions">
            <button
              v-if="isEditing"
              class="secondary-button cancel-edit-button"
              type="button"
              @click="cancelEdit"
            >
              取消编辑
            </button>
            <button class="primary-button publish-button" type="submit" :disabled="!canSubmit">
              <span>{{ loading ? '保存中...' : isEditing ? '保存修改' : '发布文章' }}</span>
              <Send v-if="!loading" :size="17" />
            </button>
          </div>
        </div>
      </form>
    </main>
  </div>
</template>
