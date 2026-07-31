import { defineStore } from 'pinia'
import {
  ListPrompts,
  GetPrompt,
  CreatePrompt,
  DeletePrompt,
  CreateVersion,
  ListVersions,
  Rollback,
  Diff,
  QualityTrend,
} from '../../wailsjs/go/main/App'
import { db, api } from '../../wailsjs/go/models'

export const usePromptsStore = defineStore('prompts', {
  state: () => ({
    prompts: [] as db.PromptDoc[],
    loading: false,
    error: null as string | null,

    detailPrompt: null as db.PromptDoc | null,
    versions: [] as db.PromptVersion[],
    quality: [] as api.VersionQuality[],
    detailLoading: false,
    detailError: null as string | null,

    diffResult: null as api.DiffResult | null,
    diffLoading: false,
    diffError: null as string | null,
  }),
  actions: {
    async fetchPrompts() {
      this.loading = true
      this.error = null
      try {
        this.prompts = (await ListPrompts()) ?? []
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async createPrompt(key: string, description: string, content: string, message: string) {
      this.error = null
      try {
        await CreatePrompt(key, description, content, message)
        await this.fetchPrompts()
      } catch (e) {
        this.error = String(e)
      }
    },
    async deletePrompt(id: string) {
      this.error = null
      try {
        await DeletePrompt(id)
        await this.fetchPrompts()
      } catch (e) {
        this.error = String(e)
      }
    },
    async fetchDetail(promptId: string) {
      this.detailLoading = true
      this.detailError = null
      this.diffResult = null
      try {
        const [prompt, versions, quality] = await Promise.all([
          GetPrompt(promptId),
          ListVersions(promptId),
          QualityTrend(promptId),
        ])
        this.detailPrompt = prompt
        this.versions = (versions ?? []).slice().sort((a, b) => b.version_no - a.version_no)
        this.quality = quality ?? []
      } catch (e) {
        this.detailError = String(e)
      } finally {
        this.detailLoading = false
      }
    },
    async createVersion(promptId: string, content: string, message: string) {
      this.detailError = null
      try {
        await CreateVersion(promptId, content, message)
        await this.fetchDetail(promptId)
      } catch (e) {
        this.detailError = String(e)
      }
    },
    async rollback(promptId: string, versionId: string) {
      this.detailError = null
      try {
        await Rollback(promptId, versionId)
        await this.fetchDetail(promptId)
      } catch (e) {
        this.detailError = String(e)
      }
    },
    async diff(fromId: string, toId: string) {
      this.diffLoading = true
      this.diffError = null
      try {
        this.diffResult = await Diff(fromId, toId)
      } catch (e) {
        this.diffError = String(e)
      } finally {
        this.diffLoading = false
      }
    },
  },
})
