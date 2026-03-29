# Text Compare (Wails + Vue3)

This is a desktop app template using [Wails](https://wails.io), Go + Vue3. It compares two texts and supports highlighting diff, synchronized scrolling, drag-resizable split layout, search and replace.

## Project layout
- `main.go`: Wails bootstrap, binds `App`.
- `app.go`: backend diff logic, `CompareText` using [`github.com/sergi/go-diff/diffmatchpatch`](go.mod).
- `frontend/src/App.vue`: frontend UI and behavior.
- `frontend/src/components/Tooltip.vue`: tooltip component.
- `frontend/src/style.css`: style sheet.
- `frontend/wailsjs/go/main/App.js` + `App.d.ts`: generated Wails JS bindings.
- `wails.json`: Wails configuration.

## Dev
1. Install Wails:
2. Start in dev mode:
   - `wails dev`

## Features
- Text comparison via [`CompareText`](frontend/wailsjs/go/main/App.js)
- Search/replace with case-sensitive, whole-word, regex options
- Next/Prev match, highlight + scroll to match
- Sync scrolling for both text areas + result
- Drag split layout (horizontal/vertical)
- Ctrl/Cmd + wheel for font scaling