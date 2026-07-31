<script setup lang="ts">
import { useI18n } from '@/i18n'

const { t } = useI18n()
</script>

<template>
  <div class="p-6 space-y-6 max-w-xl overflow-y-auto h-full">
    <div>
      <h2 class="text-sm font-semibold">{{ t('help.title') }}</h2>
      <p class="text-xs text-muted-foreground mt-1">{{ t('help.intro') }}</p>
    </div>

    <section class="space-y-2">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('help.what.title') }}</h3>
      <p class="text-sm">{{ t('help.what.body') }}</p>
    </section>

    <section class="space-y-2">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('help.start.title') }}</h3>
      <ol class="space-y-1.5 text-sm list-decimal list-inside">
        <li>{{ t('help.start.1') }}</li>
        <li>{{ t('help.start.2') }}</li>
        <li>{{ t('help.start.3') }}</li>
        <li>{{ t('help.start.4') }}</li>
      </ol>
    </section>

    <section class="space-y-2">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('help.stuck.title') }}</h3>
      <ul class="space-y-1.5 text-sm list-disc list-inside text-muted-foreground">
        <li><span class="text-foreground">{{ t('help.stuck.1') }}</span></li>
        <li><span class="text-foreground">{{ t('help.stuck.2') }}</span></li>
        <li><span class="text-foreground">{{ t('help.stuck.3') }}</span></li>
      </ul>
    </section>
  </div>
</template>
