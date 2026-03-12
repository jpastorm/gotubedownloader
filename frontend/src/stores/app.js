import { defineStore } from 'pinia'
import { CheckBootstrap, InstallYtDlp } from '../../wailsjs/go/main/App'

export const useAppStore = defineStore('app', {
  state: () => ({
    bootstrapReady: false,
    ytdlpReady: false,
    ffmpegReady: false,
    ytdlpPath: '',
    ffmpegPath: '',
    bootstrapError: '',
    installing: false,
  }),
  actions: {
    async checkBootstrap() {
      try {
        const status = await CheckBootstrap()
        this.ytdlpReady = status.ytDlpReady
        this.ffmpegReady = status.ffmpegReady
        this.ytdlpPath = status.ytDlpPath
        this.ffmpegPath = status.ffmpegPath
        this.bootstrapReady = status.ytDlpReady
        this.bootstrapError = status.error || ''
      } catch (e) {
        this.bootstrapError = e.toString()
      }
    },
    async installYtDlp() {
      this.installing = true
      this.bootstrapError = ''
      try {
        const status = await InstallYtDlp()
        this.ytdlpReady = status.ytDlpReady
        this.ytdlpPath = status.ytDlpPath
        this.ffmpegReady = status.ffmpegReady
        this.ffmpegPath = status.ffmpegPath
        this.bootstrapReady = status.ytDlpReady
        if (status.error) {
          this.bootstrapError = status.error
        }
      } catch (e) {
        this.bootstrapError = e.toString()
      } finally {
        this.installing = false
      }
    }
  }
})
