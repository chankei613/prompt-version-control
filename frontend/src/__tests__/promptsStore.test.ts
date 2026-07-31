import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { usePromptsStore } from '@/stores/prompts'
import { ListPrompts, GetPrompt, ListVersions, QualityTrend } from '../../wailsjs/go/main/App'
import { db } from '../../wailsjs/go/models'

describe('prompts store error handling', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(ListPrompts).mockReset()
    vi.mocked(GetPrompt).mockReset()
    vi.mocked(ListVersions).mockReset()
    vi.mocked(QualityTrend).mockReset()
  })

  it('captures a failed fetchPrompts() as store.error and clears loading', async () => {
    vi.mocked(ListPrompts).mockRejectedValueOnce(new Error('network down'))
    const store = usePromptsStore()

    await store.fetchPrompts()

    expect(store.loading).toBe(false)
    expect(store.error).toContain('network down')
  })

  it('clears the previous error on a successful retry', async () => {
    vi.mocked(ListPrompts).mockRejectedValueOnce(new Error('network down'))
    const store = usePromptsStore()
    await store.fetchPrompts()
    expect(store.error).not.toBeNull()

    vi.mocked(ListPrompts).mockResolvedValueOnce([])
    await store.fetchPrompts()

    expect(store.error).toBeNull()
  })

  it('fetchDetail sorts versions descending by version_no', async () => {
    vi.mocked(GetPrompt).mockResolvedValueOnce(
      db.PromptDoc.createFrom({ id: 'p1', key: 'k', description: '', current_version_id: 'v2' }),
    )
    vi.mocked(ListVersions).mockResolvedValueOnce([
      db.PromptVersion.createFrom({ id: 'v1', prompt_doc_id: 'p1', version_no: 1, content: 'a', message: '' }),
      db.PromptVersion.createFrom({ id: 'v2', prompt_doc_id: 'p1', version_no: 2, content: 'b', message: '' }),
    ])
    vi.mocked(QualityTrend).mockResolvedValueOnce([])

    const store = usePromptsStore()
    await store.fetchDetail('p1')

    expect(store.detailError).toBeNull()
    expect(store.versions.map((v) => v.version_no)).toEqual([2, 1])
  })
})
