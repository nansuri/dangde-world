import { computed, onMounted, onUnmounted, ref } from 'vue'

export function useViewportProfile() {
  const width = ref(typeof window === 'undefined' ? 1024 : window.innerWidth)
  const height = ref(typeof window === 'undefined' ? 768 : window.innerHeight)

  const orientation = computed(() => (width.value >= height.value ? 'landscape' : 'portrait'))
  const isTabletWidth = computed(() => width.value >= 768 && width.value <= 1366)
  const isLandscape = computed(() => orientation.value === 'landscape')
  const isPortrait = computed(() => orientation.value === 'portrait')
  const isTabletLandscape = computed(() => isTabletWidth.value && isLandscape.value)
  const isTabletPortrait = computed(() => isTabletWidth.value && isPortrait.value)

  function syncViewport() {
    width.value = window.innerWidth
    height.value = window.innerHeight
  }

  onMounted(() => {
    syncViewport()
    window.addEventListener('resize', syncViewport)
    window.visualViewport?.addEventListener('resize', syncViewport)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', syncViewport)
    window.visualViewport?.removeEventListener('resize', syncViewport)
  })

  return {
    width,
    height,
    orientation,
    isLandscape,
    isPortrait,
    isTabletWidth,
    isTabletLandscape,
    isTabletPortrait,
  }
}
