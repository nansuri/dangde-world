# DangDe! World iPad Safari Compatibility Plan

## Objective

Make the frontend feel native on iPad Safari with touch-first interaction, better use of the available screen, and safer behavior around Safari viewport quirks.

## Primary Issues Identified

- The current app shell assumes a desktop sidebar layout and wastes usable width on tablet.
- The app uses `100vh` patterns without accounting for Safari dynamic viewport changes.
- Safe-area insets for notches, home indicators, and browser chrome are not handled.
- Touch targets and spacing are acceptable on desktop but should be more deliberate for tablet use.
- Navigation should remain reachable and comfortable in both portrait and landscape orientations.

## Implementation Approach

### 1. Safari viewport and fullscreen support

- Add `viewport-fit=cover` to the document meta viewport.
- Add Apple web app meta tags for a better installed-web-app experience.
- Use `100svh` and `100dvh` patterns with `100vh` fallback.
- Respect `env(safe-area-inset-*)` in main layout containers.

### 2. iPad-first responsive shell

- Keep desktop sidebar behavior on large screens.
- Switch to a full-width top navigation shell for tablet widths.
- Make role navigation horizontally scrollable with large touch targets.
- Reduce wasted empty columns and allow content panels to fill the screen.

### 3. Touch interaction improvements

- Increase minimum button and field height.
- Add touch-friendly tap behavior and reduce accidental zoom triggers.
- Improve section spacing and card layouts for finger interaction.

### 4. Page-specific adjustments

- Login page should collapse earlier for tablet and prioritize the login action area.
- Dashboard grids should shift from 4-column desktop layouts to 2-column tablet layouts.
- Admin editing cards should stack cleanly on tablet.

## Acceptance Criteria

- Works cleanly in iPad Safari portrait and landscape.
- Navigation remains easily tappable without precision cursor behavior.
- No clipped content due to Safari browser chrome or safe-area issues.
- Main dashboard surfaces use most of the available screen width.
