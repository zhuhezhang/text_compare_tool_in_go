# 文本对比（Wails + Vue3）

本项目是一个基于 [Wails](https://wails.io)、Go + Vue3 的桌面应用示例，实现了两段文本对比、高亮、同步滚动、可拖拽布局、搜索与替换等功能。

## 项目结构
- `main.go`：Wails 应用入口，绑定后端 `App`。
- `app.go`：后端逻辑，`CompareText` 使用 [`github.com/sergi/go-diff/diffmatchpatch`](go.mod) 实现文本对比。
- `frontend/src/App.vue`：前端主界面逻辑及交互。
- `frontend/src/components/Tooltip.vue`：悬浮提示组件。
- `frontend/src/style.css`：样式定义。
- `frontend/wailsjs/go/main/App.js` + `App.d.ts`：后端接口自动生成。
- `wails.json`：Wails 项目配置。

## 开发（Dev）
1. 安装Wails：
2. 启动 Wails 开发模式：
   - `wails dev`

## 功能说明
- 文本对比：`frontend/src/App.vue` 调用 [`CompareText`](frontend/wailsjs/go/main/App.js)  
- 搜索替换：支持大小写、全词、正则、上下一个查询、全部替换。
- 同步滚动：左右及结果区域同步滚动（`syncScroll`）。
- 拖拽分割：左右比例与上下比例可拖拽调整。
- 文本缩放：Ctrl/Cmd + 滚轮缩放字体。