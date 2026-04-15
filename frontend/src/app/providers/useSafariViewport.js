import { onMounted, onUnmounted } from 'vue'

export function useSafariViewport() {
  let teardown = () => {}

  onMounted(() => {
    const root = document.documentElement
    const body = document.body
    const focusSelector = 'input, textarea, select, [contenteditable="true"]'

    function syncViewportFlags() {
      const width = window.innerWidth
      const height = window.innerHeight
      const orientation = width >= height ? 'landscape' : 'portrait'
      const isTablet = width >= 768 && width <= 1366
      const visualHeight = window.visualViewport?.height || height
      const keyboardOpen = height-visualHeight > 120

      root.dataset.orientation = orientation
      root.dataset.deviceClass = isTablet ? 'tablet' : width < 768 ? 'phone' : 'desktop'
      body.classList.toggle('keyboard-open', keyboardOpen)
    }

    function onFocusIn(event) {
      const target = event.target
      if (!(target instanceof HTMLElement) || !target.matches(focusSelector)) {
        return
      }

      window.setTimeout(() => {
        target.scrollIntoView({ block: 'center', inline: 'nearest', behavior: 'smooth' })
      }, 160)
    }

    function onTouchEnd(event) {
      const target = event.target
      if (!(target instanceof HTMLElement)) {
        return
      }

      const clickable = target.closest('button, [role="button"], .nav-pill, .primary-button, .ghost-button')
      if (clickable instanceof HTMLElement) {
        clickable.blur()
      }
    }

    syncViewportFlags()
    window.addEventListener('resize', syncViewportFlags)
    window.visualViewport?.addEventListener('resize', syncViewportFlags)
    document.addEventListener('focusin', onFocusIn)
    document.addEventListener('touchend', onTouchEnd, { passive: true })

    teardown = () => {
      window.removeEventListener('resize', syncViewportFlags)
      window.visualViewport?.removeEventListener('resize', syncViewportFlags)
      document.removeEventListener('focusin', onFocusIn)
      document.removeEventListener('touchend', onTouchEnd)
    }
  })

  onUnmounted(() => teardown())
}
