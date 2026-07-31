<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { usePromptsStore } from '@/stores/prompts'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const store = usePromptsStore()
const router = useRouter()

const key = ref('')
const description = ref('')
const content = ref('')
const creating = ref(false)

onMounted(() => store.fetchPrompts())

async function create() {
  if (!key.value.trim() || !content.value.trim()) return
  creating.value = true
  await store.createPrompt(key.value.trim(), description.value.trim(), content.value, 'initial version')
  creating.value = false
  key.value = ''
  description.value = ''
  content.value = ''
}

async function remove(id: string) {
  if (!confirm(t('prompts.card.delete.confirm'))) return
  await store.deletePrompt(id)
}

function open(id: string) {
  router.push(`/prompts/${id}`)
}

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <h2 class="text-sm font-semibold">{{ t('prompts.title') }}</h2>

    <div v-if="store.error" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.error }}
      <button @click="store.fetchPrompts" class="ml-2 underline">{{ t('error.retry') }}</button>
    </div>

    <div class="border border-border rounded-lg p-4 space-y-3 max-w-lg">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('prompts.new') }}</h3>
      <input v-model="key" :placeholder="t('prompts.new.key')" class="w-full text-sm border border-border rounded px-2 py-1.5" />
      <input v-model="description" :placeholder="t('prompts.new.description')" class="w-full text-sm border border-border rounded px-2 py-1.5" />
      <textarea
        v-model="content"
        :placeholder="t('prompts.new.content')"
        rows="4"
        class="w-full text-sm border border-border rounded px-2 py-1.5 font-mono"
      />
      <button
        @click="create"
        :disabled="creating || !key.trim() || !content.trim()"
        class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white disabled:opacity-40"
      >
        {{ t('prompts.new.create') }}
      </button>
    </div>

    <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>
    <div v-else-if="store.prompts.length === 0" class="text-sm text-muted-foreground">{{ t('prompts.empty') }}</div>

    <div v-else class="space-y-1.5">
      <div
        v-for="p in store.prompts"
        :key="p.id"
        class="flex items-center gap-3 border border-border rounded-lg px-4 py-2.5 border-l-4"
        style="border-left-color: #7c5cd4"
      >
        <div class="flex-1 cursor-pointer" @click="open(p.id)">
          <div class="text-sm font-medium">{{ p.key }}</div>
          <div class="text-xs text-muted-foreground">
            <span v-if="p.description">{{ p.description }} · </span>
            {{ fmt(p.updated_at) }}
          </div>
        </div>
        <button @click="open(p.id)" class="text-xs px-2 py-1 border border-border rounded hover:bg-gray-50">
          {{ t('prompts.card.open') }}
        </button>
        <button @click="remove(p.id)" class="text-xs text-red-600 hover:underline">{{ t('prompts.card.delete') }}</button>
      </div>
    </div>
  </div>
</template>
