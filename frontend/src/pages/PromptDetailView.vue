<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePromptsStore } from '@/stores/prompts'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const store = usePromptsStore()
const route = useRoute()
const router = useRouter()

const promptId = computed(() => route.params.id as string)

const newContent = ref('')
const newMessage = ref('')
const saving = ref(false)

const diffFrom = ref('')
const diffTo = ref('')

function load() {
  store.fetchDetail(promptId.value)
  diffFrom.value = ''
  diffTo.value = ''
}

onMounted(load)
watch(promptId, load)

watch(
  () => store.detailPrompt,
  (p) => {
    if (p) {
      const current = store.versions.find((v) => v.id === p.current_version_id)
      newContent.value = current?.content ?? ''
    }
  },
)

async function saveVersion() {
  if (!newContent.value.trim()) return
  saving.value = true
  await store.createVersion(promptId.value, newContent.value, newMessage.value.trim())
  saving.value = false
  newMessage.value = ''
}

async function rollbackTo(versionId: string) {
  if (!confirm(t('detail.history.rollback.confirm'))) return
  await store.rollback(promptId.value, versionId)
}

async function runDiff() {
  if (!diffFrom.value || !diffTo.value) return
  await store.diff(diffFrom.value, diffTo.value)
}

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}

function qualityFor(versionId: string) {
  return store.quality.find((q) => q.version_id === versionId)
}
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <button @click="router.push('/prompts')" class="text-xs text-muted-foreground hover:underline">
      &larr; {{ t('detail.back') }}
    </button>

    <div v-if="store.detailError" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.detailError }}
      <button @click="load" class="ml-2 underline">{{ t('error.retry') }}</button>
    </div>

    <div v-if="store.detailLoading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>

    <template v-else-if="store.detailPrompt">
      <div>
        <h2 class="text-sm font-semibold">{{ store.detailPrompt.key }}</h2>
        <p v-if="store.detailPrompt.description" class="text-xs text-muted-foreground mt-0.5">
          {{ store.detailPrompt.description }}
        </p>
      </div>

      <section class="border border-border rounded-lg p-4 space-y-3 max-w-2xl">
        <h3 class="text-xs font-semibold text-muted-foreground">{{ t('detail.newVersion') }}</h3>
        <textarea
          v-model="newContent"
          :placeholder="t('detail.newVersion.content')"
          rows="6"
          class="w-full text-sm border border-border rounded px-2 py-1.5 font-mono"
        />
        <input
          v-model="newMessage"
          :placeholder="t('detail.newVersion.message')"
          class="w-full text-sm border border-border rounded px-2 py-1.5"
        />
        <button
          @click="saveVersion"
          :disabled="saving || !newContent.trim()"
          class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white disabled:opacity-40"
        >
          {{ t('detail.newVersion.save') }}
        </button>
      </section>

      <section class="space-y-2">
        <h3 class="text-xs font-semibold text-muted-foreground">{{ t('detail.history') }}</h3>
        <div class="space-y-1.5">
          <div
            v-for="v in store.versions"
            :key="v.id"
            class="flex items-center gap-3 border border-border rounded-lg px-4 py-2.5 border-l-4"
            :style="{ borderLeftColor: v.id === store.detailPrompt.current_version_id ? '#7c5cd4' : '#e5e5e5' }"
          >
            <div class="flex-1">
              <div class="text-sm font-medium flex items-center gap-2">
                v{{ v.version_no }}
                <span v-if="v.id === store.detailPrompt.current_version_id" class="text-xs px-1.5 py-0.5 rounded bg-gray-900 text-white">
                  {{ t('detail.history.current') }}
                </span>
                <span v-if="qualityFor(v.id)" class="text-xs text-muted-foreground">
                  · avg {{ qualityFor(v.id)!.avg_score.toFixed(2) }} ({{ t('detail.quality.rated', { n: qualityFor(v.id)!.rating_count }) }})
                </span>
              </div>
              <div class="text-xs text-muted-foreground">
                <span v-if="v.message">{{ v.message }} · </span>{{ fmt(v.created_at) }}
              </div>
            </div>
            <button
              v-if="v.id !== store.detailPrompt.current_version_id"
              @click="rollbackTo(v.id)"
              class="text-xs px-2 py-1 border border-border rounded hover:bg-gray-50"
            >
              {{ t('detail.history.rollback') }}
            </button>
          </div>
        </div>
      </section>

      <section class="space-y-2 max-w-2xl">
        <h3 class="text-xs font-semibold text-muted-foreground">{{ t('detail.diff.title') }}</h3>
        <div class="flex gap-2 items-center">
          <select v-model="diffFrom" class="text-sm border border-border rounded px-2 py-1.5 flex-1">
            <option value="" disabled>{{ t('detail.diff.from') }}</option>
            <option v-for="v in store.versions" :key="v.id" :value="v.id">v{{ v.version_no }}</option>
          </select>
          <select v-model="diffTo" class="text-sm border border-border rounded px-2 py-1.5 flex-1">
            <option value="" disabled>{{ t('detail.diff.to') }}</option>
            <option v-for="v in store.versions" :key="v.id" :value="v.id">v{{ v.version_no }}</option>
          </select>
          <button @click="runDiff" :disabled="!diffFrom || !diffTo" class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white disabled:opacity-40">
            {{ t('detail.diff.title') }}
          </button>
        </div>

        <div v-if="store.diffError" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
          {{ t('error.prefix') }}{{ store.diffError }}
        </div>
        <div v-else-if="store.diffLoading" class="text-xs text-muted-foreground">{{ t('loading') }}</div>
        <div v-else-if="!store.diffResult" class="text-xs text-muted-foreground">{{ t('detail.diff.empty') }}</div>
        <pre v-else class="text-xs font-mono border border-border rounded-lg p-3 overflow-x-auto whitespace-pre-wrap"><span
          v-for="(line, i) in store.diffResult.lines"
          :key="i"
          :class="{
            'bg-green-50 text-green-800 block': line.type === 'add',
            'bg-red-50 text-red-800 block': line.type === 'remove',
            'text-muted-foreground block': line.type === 'equal',
          }"
        >{{ line.type === 'add' ? '+' : line.type === 'remove' ? '-' : ' ' }} {{ line.text }}</span></pre>
      </section>

      <section class="space-y-2 max-w-2xl">
        <h3 class="text-xs font-semibold text-muted-foreground">{{ t('detail.quality.title') }}</h3>
        <div v-if="store.quality.length === 0" class="text-xs text-muted-foreground">{{ t('detail.quality.empty') }}</div>
        <div v-else class="flex items-end gap-2 h-24">
          <div v-for="q in store.quality" :key="q.version_id" class="flex-1 flex flex-col items-center gap-1">
            <div
              class="w-full rounded-t bg-[#7c5cd4]"
              :style="{ height: Math.max(4, q.avg_score * 20) + 'px' }"
              :title="`v${q.version_no}: ${q.avg_score.toFixed(2)} (${q.rating_count})`"
            />
            <span class="text-[10px] text-muted-foreground">v{{ q.version_no }}</span>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>
