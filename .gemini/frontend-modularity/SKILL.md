---
name: frontend-modularity
description: Enforce modular frontend architecture. Use when creating or refactoring UI components, styles, or logic to ensure reusability, modularity, and slim files.
---

# Frontend Modularity

This skill provides a set of architectural rules and workflows to maintain a clean, maintainable, and modular frontend codebase.

## Core Mandates

1.  **Components First**:
    *   Always use components. If a UI element is used more than once, or if a component is becoming too complex, extract it into a reusable component.
    *   Prefer small, single-purpose components over large, monolithic ones.
    *   Locate reusable components in `src/shared/ui/` or `src/components/`.

2.  **Modular Structure**:
    *   **NEVER create "One Big File"**. If a file (CSS, JS, or Vue) exceeds ~200 lines, consider it a candidate for refactoring.
    *   Split logic into hooks/composables (e.g., `src/shared/lib/` or `src/features/*/hooks/`).
    *   Split styles into component-specific styles or theme modules.

3.  **Slim Files**:
    *   Each file should have a single responsibility.
    *   Keep templates, scripts, and styles focused. If a Vue file has more than 50 lines of CSS, move the CSS to a scoped block or a separate module if it's shared.

## CSS Refactoring Workflow

When refactoring large CSS files (like `main.css`):

1.  **Identify Component Styles**: Extract styles that are specific to a single component and move them into that component's `<style scoped>` block.
2.  **Identify Feature Styles**: Group styles related to a specific feature (e.g., `dashboard`, `auth`) and move them into feature-specific CSS files or component blocks.
3.  **Global Theme & Utilities**: Keep only the following in `main.css`:
    *   CSS Variable definitions (`:root`).
    *   Base element resets (`html`, `body`, `*`).
    *   Highly generic utility classes (e.g., `.sr-only`, `.flex-center`).
    *   Global layout primitives that cannot be componentized yet.
4.  **Consistency**: Ensure that extracted styles still use the global CSS variables for colors, spacing, and shadows.

## Example: Splitting a Large Component

**Before (Monolith):**
`BigPage.vue` (500 lines) contains header, sidebar, main content, and 200 lines of CSS.

**After (Modular):**
*   `BigPage.vue`: Orchestrates small components.
*   `PageHeader.vue`: Just the header.
*   `PageSidebar.vue`: Just the sidebar.
*   `usePageData.js`: Composable for fetching data.
*   Styles are scoped within each component.
