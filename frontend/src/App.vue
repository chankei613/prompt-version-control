<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterView, RouterLink, useRoute } from 'vue-router'
import { EventsOn } from '../wailsjs/runtime/runtime'
import { useI18n } from '@/i18n'

const route = useRoute()
const { t, toggleLocale } = useI18n()

// 実行エラー通知（バックグラウンド処理のエラーを表示）
const execError = ref<string | null>(null)
let errorTimer: ReturnType<typeof setTimeout> | null = null

onMounted(() => {
  // バックグラウンド処理のエラーイベントをリッスン
  EventsOn('execution:error', (data: { error: string }) => {
    execError.value = data.error
    if (errorTimer) clearTimeout(errorTimer)
    errorTimer = setTimeout(() => { execError.value = null }, 6000)
  })
})

function isActive(prefix: string) {
  return route.path.startsWith(prefix)
}

const navItems = [
  {
    to: '/prompts',
    key: 'nav.prompts',
    icon: `<line x1="6" x2="6" y1="3" y2="15"/><circle cx="18" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><path d="M18 9a9 9 0 0 1-9 9"/>`,
  },
  {
    to: '/help',
    key: 'nav.help',
    icon: `<circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/>`,
  },
  {
    to: '/settings',
    key: 'nav.settings',
    icon: `<path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/>`,
  },
]
</script>

<template>
  <div class="flex h-screen bg-background text-foreground overflow-hidden">
    <!-- Sidebar -->
    <aside class="w-52 border-r border-border flex flex-col shrink-0 bg-background">
      <!-- Header: TitleBarHiddenInset のトラフィックライト分の余白（pt-16）が必要 -->
      <div class="px-4 pb-4 pt-16 border-b border-border" style="-webkit-app-region: drag">
        <div class="flex items-start justify-between">
          <div class="flex items-center gap-2">
            <svg class="w-5 h-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="#7c5cd4" stroke-width="1.75">
              <line x1="6" x2="6" y1="3" y2="15"/><circle cx="18" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><path d="M18 9a9 9 0 0 1-9 9"/>
            </svg>
            <div>
              <h1 class="text-sm font-semibold tracking-tight text-foreground">Prompt Version Control</h1>
              <p class="text-xs text-muted-foreground">{{ t('app.subtitle') }}</p>
            </div>
          </div>
          <!-- 言語切替ボタン -->
          <button
            @click="toggleLocale"
            style="-webkit-app-region: no-drag"
            class="text-xs text-gray-400 hover:text-gray-700 px-1.5 py-0.5 rounded border border-gray-200 hover:border-gray-400 transition-colors shrink-0 mt-0.5"
          >{{ t('lang.toggle') }}</button>
        </div>
      </div>

      <nav class="flex-1 p-2 space-y-0.5">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-2.5 px-3 py-2 rounded-md text-sm transition-colors"
          :class="isActive(item.to)
            ? 'bg-gray-900 text-white'
            : 'text-gray-500 hover:bg-gray-100 hover:text-gray-900'"
        >
          <svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" v-html="item.icon" />
          <span>{{ t(item.key) }}</span>
        </RouterLink>
      </nav>
    </aside>

    <!-- Main content -->
    <main class="flex-1 overflow-auto bg-gray-50/50 flex flex-col">
      <!-- エラーバナー（バックグラウンド処理のエラー通知） -->
      <div v-if="execError"
        class="mx-4 mt-3 px-4 py-2.5 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700 flex items-center justify-between gap-3 shrink-0">
        <span>{{ execError }}</span>
        <button @click="execError = null" class="text-red-400 hover:text-red-600 shrink-0">
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      <RouterView class="flex-1" />
    </main>
  </div>
</template>
