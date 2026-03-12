import { defineStore } from 'pinia'
import { GetSettings, SaveSettings, SelectFolder } from '../../wailsjs/go/main/App'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    settings: {
      downloadDir: '',
      maxConcurrent: 2,
      fragmentConcurrent: 4,
      defaultMode: 'video',
      subtitleLang: 'en',
      speedLimit: '',
      proxy: '',
      skipDuplicates: true,
      continueDownloads: true,
      embedMetadata: false,
      downloadSubtitles: false,
      embedSubtitles: false,
      notifyOnComplete: true,
    },
    loaded: false,
  }),
  actions: {
    async load() {
      try {
        const s = await GetSettings()
        this.settings = { ...this.settings, ...s }
        this.loaded = true
      } catch (e) {
        console.error('Failed to load settings:', e)
      }
    },
    async save() {
      try {
        await SaveSettings(this.settings)
      } catch (e) {
        console.error('Failed to save settings:', e)
        throw e
      }
    },
    async selectFolder() {
      try {
        const dir = await SelectFolder()
        if (dir) {
          this.settings.downloadDir = dir
          await this.save()
        }
        return dir
      } catch (e) {
        console.error('Failed to select folder:', e)
        return null
      }
    }
  }
})
