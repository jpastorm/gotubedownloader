import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  StartDownload, CancelDownload, RetryDownload, RemoveDownload,
  ClearCompleted, CancelAllDownloads, GetQueue, AnalyzeURL,
  AnalyzePlaylistURL, StartBatchDownload, HideWindow
} from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'

export const useDownloadStore = defineStore('downloads', () => {
  const queue = ref([])
  const analyzing = ref(false)
  const analysisResult = ref(null)
  const analysisError = ref('')
  const analyzeGen = ref(0) // generation counter to discard stale analysis results

  const activeCount = computed(() => queue.value.filter(i => i.status === 'downloading' || i.status === 'queued').length)
  const completedCount = computed(() => queue.value.filter(i => i.status === 'completed').length)

  function init() {
    EventsOn('queue:update', (item) => {
      const idx = queue.value.findIndex(i => i.id === item.id)
      if (idx >= 0) {
        queue.value[idx] = { ...queue.value[idx], ...item }
      } else {
        queue.value.push(item)
      }
    })
    refreshQueue()
  }

  async function refreshQueue() {
    try {
      const items = await GetQueue()
      queue.value = items || []
    } catch (e) {
      console.error('Failed to refresh queue:', e)
    }
  }

  async function analyze(url) {
    const myGen = ++analyzeGen.value
    analyzing.value = true
    analysisResult.value = null
    analysisError.value = ''
    try {
      const result = await AnalyzeURL(url)
      if (myGen !== analyzeGen.value) return null // stale, user already did instant download
      analysisResult.value = result
      return result
    } catch (e) {
      if (myGen !== analyzeGen.value) return null // stale
      analysisError.value = e.toString()
      return null
    } finally {
      if (myGen === analyzeGen.value) {
        analyzing.value = false
      }
    }
  }

  function cancelAnalysis() {
    analyzeGen.value++ // invalidate any pending analysis
    analyzing.value = false
    analysisResult.value = null
    analysisError.value = ''
  }

  async function download(url, mode, quality) {
    try {
      const id = await StartDownload(url, mode || 'video', quality || 'high')
      refreshQueue() // don't await — let it update in background
      return id
    } catch (e) {
      throw e
    }
  }

  async function analyzePlaylist(url) {
    const myGen = ++analyzeGen.value
    analyzing.value = true
    analysisResult.value = null
    analysisError.value = ''
    try {
      const result = await AnalyzePlaylistURL(url)
      if (myGen !== analyzeGen.value) return null
      analysisResult.value = result
      return result
    } catch (e) {
      if (myGen !== analyzeGen.value) return null
      analysisError.value = e.toString()
      return null
    } finally {
      if (myGen === analyzeGen.value) {
        analyzing.value = false
      }
    }
  }

  async function batchDownload(urls, mode, quality) {
    try {
      const ids = await StartBatchDownload(urls, mode, quality || 'high')
      refreshQueue()
      return ids
    } catch (e) {
      throw e
    }
  }

  async function cancel(id) { await CancelDownload(id) }
  async function retry(id) { await RetryDownload(id); await refreshQueue() }
  async function remove(id) { await RemoveDownload(id); await refreshQueue() }
  async function clearCompleted() { await ClearCompleted(); await refreshQueue() }
  async function cancelAll() { await CancelAllDownloads() }

  return {
    queue, analyzing, analysisResult, analysisError,
    activeCount, completedCount,
    init, refreshQueue, analyze, analyzePlaylist, cancelAnalysis, download, batchDownload,
    cancel, retry, remove, clearCompleted, cancelAll
  }
})
