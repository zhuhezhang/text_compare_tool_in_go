<template>
  <div class="resizable-layout" :style="{ '--font-size': fontSize + 'px' }" ref="layoutContainer">
    <!-- 顶部区域 -->
    <div class="top-container" :style="{ height: topHeight + '%' }">
      <!-- 左侧子容器 -->
      <div class="left-sub-container" :style="{ width: leftWidth + '%' }">
        <div class="line-numbers" ref="leftLineNumbers">
          <div v-for="n in leftLineCount" :key="n" class="line-number">{{ n }}</div>
        </div>
        <textarea v-model="textarea1Content" class="textarea textarea1" placeholder="请输入对比文本..."
          @scroll="syncScroll('left')" @input="updateLineCount('left')" @click="setSearchScope('left')"></textarea>
      </div>
      <!-- 垂直拖拽条（调整左右比例） -->
      <div class="drag-handle vertical" :style="{ left: leftWidth + '%' }" @mousedown="startDrag('left-right')"
        @dblclick="resetWidthOrHeight('left-right')"></div>
      <!-- 右侧子容器 -->
      <div class="right-sub-container" :style="{ width: (100 - leftWidth) + '%' }">
        <div class="line-numbers" ref="rightLineNumbers">
          <div v-for="n in rightLineCount" :key="n" class="line-number">{{ n }}</div>
        </div>
        <textarea v-model="textarea2Content" class="textarea textarea2" placeholder="请输入对比文本..."
          @scroll="syncScroll('right')" @input="updateLineCount('right')" @click="setSearchScope('right')"></textarea>
      </div>
    </div>
    <!-- 水平拖拽条（调整上下比例） -->
    <div class="drag-handle horizontal" @mousedown="startDrag('top-bottom')"
      @dblclick="resetWidthOrHeight('top-bottom')"></div>
    <!-- 底部容器 -->
    <div class="bottom-container" :style="{ height: (100 - topHeight) + '%' }">
      <div class="line-numbers" ref="bottomLineNumbers">
        <div v-for="n in bottomLineCount" :key="n" class="line-number">{{ n }}</div>
      </div>
      <div class="result-content" ref="resultContent" @scroll="syncScroll('bottom')" v-html="cmpResult"></div>
    </div>
  </div>
  <!-- 悬浮搜索按钮 -->
  <div class="floating-btn-container search-btn-container" @click="openSearchBox">
    <div class="floating-btn">
      <div class="search-icon"></div>
    </div>
  </div>
  <!-- 悬浮运行按钮 -->
  <div class="floating-btn-container compare-btn-container" @click="startCompare">
    <div class="floating-btn">
      <div class="triangle"></div>
    </div>
  </div>


  <!-- 可拖拽搜索替换窗口 -->
  <div v-if="showSearchBox" class="search-floating-window" :style="{
    left: windowPosition.x + 'px',
    top: windowPosition.y + 'px',
    opacity: windowOpacity,
    transform: `scale(${windowScale})`
  }" ref="searchWindow">
    <!-- 窗口标题栏（拖拽区域） -->
    <div class="window-header" @mousedown="startWindowDrag">
      <div class="window-title">搜索与替换</div>
      <div class="window-controls">
        <button class="window-btn close-btn" @click="closeSearchBox" title="关闭">×</button>
      </div>
    </div>
    <!-- 窗口内容区域 -->
    <div class="window-content">
      <!-- 搜索区域 -->
      <div class="search-input-group">
        <div class="search-group">
          <input type="text" v-model="searchText" placeholder="搜索内容" @keyup.enter="searchNext" ref="searchInput"
            class="search-replace-input">
          <button @click="performSearch" class="search-replace-btn search-btn">搜索</button>
        </div>
        <div class="search-options">
          <div class="options-left">
            <label class="search-option">
              <Tooltip text="区分大小写" position="right">
                <input type="checkbox" v-model="caseSensitive"> Aa
              </Tooltip>
            </label>
            <label class="search-option">
              <Tooltip text="全字匹配" position="right">
                <input type="checkbox" v-model="wholeWord"> ab
              </Tooltip>
            </label>
            <label class="search-option">
              <Tooltip text="使用正则表达式" position="right">
                <input type="checkbox" v-model="useRegex"> .*
              </Tooltip>
            </label>
          </div>
          <div class="options-right">
            <div class="search-info" v-if="matches.length > 0">
              第 {{ currentMatchIndex + 1 }} / {{ matches.length }}
            </div>
            <div class="search-nav">
              <button @click="searchPrevious" class="search-replace-btn"
                :disabled="matches.length === 0"><span>◀</span>上一个</button>
              <button @click="searchNext" class="search-replace-btn"
                :disabled="matches.length === 0">下一个<span>▶</span></button>
            </div>
          </div>
        </div>
      </div>
      <!-- 替换区域 -->
      <div class="replace-input-group" v-if="showReplace">
        <input type="text" v-model="replaceText" placeholder="替换为" class="search-replace-input replace-input">
        <button @click="replaceAll" class="search-replace-btn replace-all-btn"
          :disabled="matches.length == 0">全部替换</button>
        <button @click="replaceCurrent" class="search-replace-btn" :disabled="matches.length === 0">替换</button>
      </div>
    </div>
  </div>
</template>

<script>
import { CompareText } from "../wailsjs/go/main/App"
import Tooltip from './components/Tooltip.vue';

export default {
  data() {
    return {
      //主界面相关
      topHeight: 50,
      leftWidth: 50,
      isDragging: false,
      dragType: '',
      textarea1Content: '',
      textarea2Content: '',
      cmpResult: '',
      leftLineCount: 1,
      rightLineCount: 1,
      bottomLineCount: 1,
      isScrolling: false,
      fontSize: 16,
      minFontSize: 10,
      maxFontSize: 30,

      //搜索替换窗口UI界面相关
      showSearchBox: false,
      windowPosition: { x: 100, y: 100 },
      windowOpacity: 0,
      windowScale: 0.95,
      isDraggingWindow: false,
      dragStartPos: { x: 0, y: 0 },
      windowStartPos: { x: 0, y: 0 },
      windowSize: { width: 450, height: 189 },

      //搜索替换逻辑相关
      searchText: '',
      replaceText: '',
      caseSensitive: false,
      wholeWord: false,
      useRegex: false,
      showReplace: true,
      searchScope: 'left',
      matches: [],
      currentMatchIndex: -1,
    };
  },
  mounted() {
    this.updateTopLineCounts();
    this.setupWheelZoom();
    this.setupKeyboardShortcuts();
  },
  beforeDestroy() {
    this.cleanupWheelZoom();
    this.cleanupWindowListeners();
  },
  methods: {
    /**
     * 更新左右文本区域的行数
     */
    updateTopLineCounts() {
      this.updateLineCount('left');
      this.updateLineCount('right');
    },

    /**
     * 添加滚轮监听事件
     */
    setupWheelZoom() {
      this.handleWheelRef = this.handleWheel.bind(this);//将 handleWheel 方法绑定到当前组件实例，保存为引用，以便后续可以正确移除事件监听器
      this.$refs.layoutContainer.addEventListener('wheel', this.handleWheelRef, { passive: false });//主界面监听滚轮事件；添加处理函数；阻止滚轮默认行为
    },

    /**
     * 监听键盘事件，Ctrl + F或Command + F 打开替换搜索框
     */
    setupKeyboardShortcuts() {
      document.addEventListener('keydown', (e) => {
        if ((e.ctrlKey || e.metaKey) && e.key === 'f') {
          e.preventDefault();// 阻止默认行为
          this.openSearchBox();
        }
      });
    },

    /**
     * 清理滚轮监听事件
     */
    cleanupWheelZoom() {
      if (this.$refs.layoutContainer && this.handleWheelRef) {
        this.$refs.layoutContainer.removeEventListener('wheel', this.handleWheelRef);
      }
    },

    /**
     * 处理滚轮监听事件：允许用户通过 Ctrl+鼠标滚轮（或 Cmd+鼠标滚轮）来调整整个界面的字体大小
     */
    handleWheel(event) {
      if (event.ctrlKey || event.metaKey) {
        event.preventDefault();
        event.stopPropagation();

        if (event.deltaY < 0) {
          this.fontSize = Math.min(this.maxFontSize, this.fontSize + 2);
        } else {
          this.fontSize = Math.max(this.minFontSize, this.fontSize - 2);
        }
      }
    },

    /**
     * 鼠标按住拖动调整左右比例/上下比例
     * @param type {string} left-right：调整左右占比；top-bottom：调整上下占比
     */
    startDrag(type) {
      this.dragType = type;
      this.isDragging = true;
      document.body.style.userSelect = 'none';//禁用文本选择，避免拖拽时选中文本

      const handleDrag = (e) => {//定义拖拽处理函数，鼠标移动时调用
        if (!this.isDragging) return;
        const layout = document.querySelector('.resizable-layout');
        const rect = layout.getBoundingClientRect();//获取布局容器的位置和尺寸信息

        if (this.dragType === 'left-right') {
          const topRect = document.querySelector('.top-container').getBoundingClientRect();
          const offsetX = e.clientX - topRect.left;//鼠标相对于顶部容器的水平位置
          const newLeftPercent = (offsetX / topRect.width) * 100;
          if (newLeftPercent > 10 && newLeftPercent < 90) {//左侧容器宽度必须在 10%-90% 之间，防止完全遮挡
            this.leftWidth = newLeftPercent;
          }
        } else if (this.dragType === 'top-bottom') {
          const offsetY = e.clientY - rect.top;
          const newTopPercent = (offsetY / rect.height) * 100;

          if (newTopPercent > 20 && newTopPercent < 80) {
            this.topHeight = newTopPercent;
          }
        }
      };

      const stopDrag = () => {//定义停止拖拽函数，鼠标释放时调用
        this.isDragging = false;
        this.dragType = '';
        document.body.style.userSelect = '';//恢复文本选择功能
        document.removeEventListener('mousemove', handleDrag);//移除事件监听器，防止内存泄漏
        document.removeEventListener('mouseup', stopDrag);
      };

      document.addEventListener('mousemove', handleDrag);
      document.addEventListener('mouseup', stopDrag);
    },

    /**
     * 恢复左右占比和上下之间的占比
     * @param type {string} left-right：恢复左右占比；top-bottom：恢复上下占比
     */
    resetWidthOrHeight(type) {
      if (type === 'left-right') {
        this.leftWidth = 50;
      } else if (type === 'top-bottom') {
        this.topHeight = 50;
      }
    },

    /**
     * 根据换行符个数更新左/右文本行数
     * @param side {string} left：左侧；right：右侧
     */
    updateLineCount(side) {
      let lineCount = 1;

      if (side === 'left') {
        lineCount = Math.max(1, this.textarea1Content.split('\n').length);
        this.leftLineCount = lineCount;
      } else if (side === 'right') {
        lineCount = Math.max(1, this.textarea2Content.split('\n').length);
        this.rightLineCount = lineCount;
      }
    },

    /**
     * 实现三个编辑区域同步滚动
     * @param source {string} left：左侧；right：右侧；bottom：下侧
     */
    syncScroll(source) {
      if (this.isScrolling) return;//防止滚动循环，如果已经在处理滚动事件，直接返回，防止无限递归
      this.isScrolling = true;

      setTimeout(() => {//使用 setTimeout 延迟 10ms 执行同步操作，确保浏览器完成当前的滚动渲染后再进行同步，避免同步过程中的性能问题
        const leftTextarea = document.querySelector('.textarea1');
        const rightTextarea = document.querySelector('.textarea2');
        const resultContent = this.$refs.resultContent;
        const leftLineNumbers = this.$refs.leftLineNumbers;
        const rightLineNumbers = this.$refs.rightLineNumbers;
        const bottomLineNumbers = this.$refs.bottomLineNumbers;

        if (source === 'left') {//滚动源是左侧区域
          if (rightTextarea) {
            rightTextarea.scrollTop = leftTextarea.scrollTop;
            rightTextarea.scrollLeft = leftTextarea.scrollLeft;
          }
          if (resultContent) {
            resultContent.scrollTop = leftTextarea.scrollTop;
            resultContent.scrollLeft = leftTextarea.scrollLeft;
          }
          if (leftLineNumbers) leftLineNumbers.scrollTop = leftTextarea.scrollTop;
          if (rightLineNumbers) rightLineNumbers.scrollTop = leftTextarea.scrollTop;
          if (bottomLineNumbers) bottomLineNumbers.scrollTop = leftTextarea.scrollTop;
        }
        else if (source === 'right') {//滚动源是右侧区域
          if (leftTextarea) {
            leftTextarea.scrollTop = rightTextarea.scrollTop;
            leftTextarea.scrollLeft = rightTextarea.scrollLeft;
          }
          if (resultContent) {
            resultContent.scrollTop = rightTextarea.scrollTop;
            resultContent.scrollLeft = rightTextarea.scrollLeft;
          }
          if (leftLineNumbers) leftLineNumbers.scrollTop = rightTextarea.scrollTop;
          if (rightLineNumbers) rightLineNumbers.scrollTop = rightTextarea.scrollTop;
          if (bottomLineNumbers) bottomLineNumbers.scrollTop = rightTextarea.scrollTop;
        }
        else if (source === 'bottom') {//滚动源是下侧区域
          if (leftTextarea) {
            leftTextarea.scrollTop = resultContent.scrollTop;
            leftTextarea.scrollLeft = resultContent.scrollLeft;
          }
          if (rightTextarea) {
            rightTextarea.scrollTop = resultContent.scrollTop;
            rightTextarea.scrollLeft = resultContent.scrollLeft;
          }
          if (leftLineNumbers) leftLineNumbers.scrollTop = resultContent.scrollTop;
          if (rightLineNumbers) rightLineNumbers.scrollTop = resultContent.scrollTop;
          if (bottomLineNumbers) bottomLineNumbers.scrollTop = resultContent.scrollTop;
        }

        this.isScrolling = false;
      }, 10);
    },


    /**
     * 调用后端的go执行文本对比
     */
    async startCompare() {
      CompareText(this.textarea1Content, this.textarea2Content).then((result) => {
        this.cmpResult = result;
        this.bottomLineCount = Math.max(this.textarea1Content.split('\n').length, this.textarea2Content.split('\n').length);
      });
    },


    /**
     * 打开替换搜索界面
     */
    openSearchBox() {
      if (this.showSearchBox) {
        this.closeSearchBox();
        return;
      }
      this.showSearchBox = true;
      this.windowOpacity = 0;
      this.windowScale = 0.5;
      const windowWidth = this.windowSize.width;// 设置窗口初始位置
      const windowHeight = this.windowSize.height;
      this.windowPosition = {
        x: Math.max(20, (window.innerWidth - windowWidth) / 2),
        y: Math.max(20, (window.innerHeight - windowHeight) / 4)
      };

      this.$nextTick(() => { // 确保DOM更新完成后执行
        setTimeout(() => {// 实现淡入和缩放动画
          this.windowOpacity = 1;
          this.windowScale = 1;
        }, 10);

        if (this.$refs.searchInput) {// 搜索框加载完成聚集到搜索框
          this.$refs.searchInput.focus();
        }
      });
    },

    /**
     * 关闭搜索替换界面
     */
    closeSearchBox() {
      this.windowOpacity = 0;
      this.windowScale = 0.5;
      this.showSearchBox = false;
    },

    /**
     * 设置搜索替换区域
     * @param {string} side left：搜索区域为左；right：搜索区域为右
     */
    setSearchScope(side) {
      if (side === 'left') {
        this.searchScope = 'left';
      } else if (side === 'right') {
        this.searchScope = 'right';
      }
    },

    /**
     * 实现可拖拽悬浮窗口功能的核心函数，允许用户点击并拖动搜索窗口的标题栏来移动窗口位置
     * @param event 
     */
    startWindowDrag(event) {
      event.preventDefault();
      this.isDraggingWindow = true;
      this.dragStartPos = {//记录鼠标按下时的初始位置（屏幕坐标）
        x: event.clientX,
        y: event.clientY
      };
      this.windowStartPos = {//记录窗口的初始位置（相对于页面左上角）
        x: this.windowPosition.x,
        y: this.windowPosition.y
      };

      const handleDrag = (e) => {//定义拖拽处理函数
        if (!this.isDraggingWindow) return;
        const dx = e.clientX - this.dragStartPos.x;//计算鼠标从按下点到当前位置的移动距离
        const dy = e.clientY - this.dragStartPos.y;
        const maxX = window.innerWidth - (this.$refs.searchWindow?.offsetWidth || 450);//计算新位置，确保窗口不会移出屏幕。允许的最大 X 坐标 = 屏幕宽度 - 窗口宽度（默认450px）。使用可选链 ?. 防止 searchWindow 未加载时报错
        const maxY = window.innerHeight - 40;//允许的最大 Y 坐标 = 屏幕高度 - 40px（保留底部空间）
        let newX = this.windowStartPos.x + dx;//计算新位置并限制边界
        let newY = this.windowStartPos.y + dy;
        newX = Math.max(0, Math.min(maxX, newX));
        newY = Math.max(0, Math.min(maxY, newY));
        this.windowPosition = { x: newX, y: newY };
      };

      const stopDrag = () => {//定义停止拖拽函数
        this.isDraggingWindow = false;
        document.removeEventListener('mousemove', handleDrag);
        document.removeEventListener('mouseup', stopDrag);
      };

      document.addEventListener('mousemove', handleDrag);
      document.addEventListener('mouseup', stopDrag);
    },

    /**
     * 执行文本搜索核心
     */
    performSearch() {
      if (!this.searchText.trim()) {// 空则退出搜索
        this.matches = [];
        this.currentMatchIndex = -1;
        return;
      }
      let content = this.getContentByScope();//获取目标搜索的文本，非空则继续
      if (!content) return;

      let searchPattern;
      try {
        if (this.useRegex) {//使用正则表达式
          const flags = this.caseSensitive ? 'g' : 'gi';//g（全局匹配）加上可选的i（不区分大小写）
          searchPattern = new RegExp(this.searchText, flags);
        } else {//不使用正则表达式
          let escapedText = this.searchText.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');//转义正则表达式中的特殊字符
          if (this.wholeWord) {
            escapedText = `\\b${escapedText}\\b`;//如果启用全字匹配，添加\b（单词边界）标记
          }
          const flags = this.caseSensitive ? 'g' : 'gi';
          searchPattern = new RegExp(escapedText, flags);
        }
      } catch (e) {
        console.error('无效的正则表达式:', e);
        return;
      }

      const matches = [];//存储匹配结果
      let match;
      while ((match = searchPattern.exec(content)) !== null) {//使用exec()方法循环查找所有匹配
        matches.push({
          index: match.index,
          length: match[0].length,
          text: match[0]
        });
      }
      this.matches = matches;
      this.currentMatchIndex = matches.length > 0 ? 0 : -1;
      this.hightAndScrollToMatch();
    },

    /**
     * 根据搜索区域（左/右对比文本）返回搜索的文本
     * @returns {string} 搜索文本
     */
    getContentByScope() {
      switch (this.searchScope) {
        case 'left':
          return this.textarea1Content;
        case 'right':
          return this.textarea2Content;
        default:
          return '';
      }
    },

    /**
     * 获取目标（左/右对比文本textarea）DOM
     * @returns 目标DOM
     */
    getTargetElement() {
      switch (this.searchScope) {
        case 'left':
          return document.querySelector('.textarea1');
        case 'right':
          return document.querySelector('.textarea2');
        default:
          return null;
      }
    },

    /**
     * 高亮显示匹配文本名，并将用户界面滚动到当前匹配的搜索项
     */
    hightAndScrollToMatch() {
      if (this.matches.length === 0 || this.currentMatchIndex < 0) return;
      const targetElement = this.getTargetElement();
      if (!targetElement) return;

      const match = this.matches[this.currentMatchIndex];

      targetElement.focus();//1. 设置选中范围并聚焦
      targetElement.setSelectionRange(match.index, match.index + match.length);
      targetElement.selectionStart = match.index;//2. 使用这个简单的技巧：将光标移到匹配项末尾，然后移回开头，这通常会触发浏览器的自动滚动
      targetElement.selectionEnd = match.index;
      targetElement.setSelectionRange(match.index, match.index + match.length);//3. 然后选中整个匹配项
      this.adjustScrollPosition(targetElement, match.index, match.index + match.length);//4. 手动调整滚动位置以确保完全可见
    },

    /**
     * 调整滚动位置
     */
    adjustScrollPosition(textarea, start, end) {
      const value = textarea.value;//获取文本内容
      const textBeforeStart = value.substring(0, start);//计算匹配项所在的行
      const lines = textBeforeStart.split('\n');
      const lineNumber = lines.length - 1;
      const style = getComputedStyle(textarea);//获取样式
      const lineHeight = parseInt(style.lineHeight) || 20;//默认行高
      const targetLineTop = lineNumber * lineHeight;//计算目标垂直位置（让匹配项所在行在中间）
      const targetScrollTop = targetLineTop - textarea.clientHeight / 2;

      textarea.scrollTop = Math.max(0, targetScrollTop);// 应用垂直滚动
      const currentLineText = lines[lines.length - 1];//水平滚动处理，获取当前行的文本
      const canvas = document.createElement('canvas');//使用canvas测量文本宽度
      const ctx = canvas.getContext('2d');
      ctx.font = style.font;

      const textUpToMatch = currentLineText.substring(0, start - (value.substring(0, start).lastIndexOf('\n') + 1 || 0));//测量当前行到匹配项开始位置的宽度
      const widthUpToMatch = ctx.measureText(textUpToMatch).width;

      const matchText = value.substring(start, end);//测量匹配项宽度
      const matchWidth = ctx.measureText(matchText).width;

      const matchEndPosition = widthUpToMatch + matchWidth;//计算匹配项的结束位置

      const currentScrollLeft = textarea.scrollLeft;//获取当前的滚动位置
      const visibleStart = currentScrollLeft;
      const visibleEnd = currentScrollLeft + textarea.clientWidth;

      if (widthUpToMatch < visibleStart || matchEndPosition > visibleEnd) {//如果匹配项不在可见区域内，调整滚动
        const targetScrollLeft = widthUpToMatch - textarea.clientWidth / 2 + matchWidth / 2;// 、让匹配项居中
        textarea.scrollLeft = Math.max(0, targetScrollLeft);
      }
    },

    /**
     * 搜索下一项
     */
    searchNext() {
      if (this.matches.length === 0) {
        this.performSearch();
        return;
      }
      this.currentMatchIndex = (this.currentMatchIndex + 1) % this.matches.length;
      this.hightAndScrollToMatch();
    },

    /**
     * 搜索上一项
     */
    searchPrevious() {
      if (this.matches.length === 0) {
        this.performSearch();
        return;
      }
      this.currentMatchIndex = this.currentMatchIndex <= 0 ? this.matches.length - 1 : this.currentMatchIndex - 1;
      this.hightAndScrollToMatch();
    },

    /**
     * 替换当前选中的匹配项后重新搜索
     */
    replaceCurrent() {
      if (this.matches.length === 0 || this.currentMatchIndex < 0) return;
      const match = this.matches[this.currentMatchIndex];
      let content = this.getContentByScope();
      const before = content.substring(0, match.index);
      const after = content.substring(match.index + match.length);
      content = before + this.replaceText + after;
      this.updateContentByScope(content);
      this.performSearch();
    },

    /**
     * 替换当前匹配所有项
     */
    replaceAll() {
      if (this.matches.length === 0) return;

      let content = this.getContentByScope();
      let searchPattern;
      try {
        if (this.useRegex) {
          const flags = this.caseSensitive ? 'g' : 'gi';
          searchPattern = new RegExp(this.searchText, flags);
        } else {
          let escapedText = this.searchText.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
          if (this.wholeWord) {
            escapedText = `\\b${escapedText}\\b`;
          }
          const flags = this.caseSensitive ? 'g' : 'gi';
          searchPattern = new RegExp(escapedText, flags);
        }
      } catch (e) {
        console.error('无效的正则表达式:', e);
        return;
      }

      content = content.replace(searchPattern, this.replaceText);
      this.updateContentByScope(content);
      this.matches = [];
      this.currentMatchIndex = -1;
    },

    /**
     * 根据搜索区域左/右更新其文本内容
     * @param content left/right
     */
    updateContentByScope(content) {
      switch (this.searchScope) {
        case 'left':
          this.textarea1Content = content;
          break;
        case 'right':
          this.textarea2Content = content;
          break;
      }
    },
  },

  watch: {
    textarea1Content() {
      this.updateLineCount('left');
    },
    textarea2Content() {
      this.updateLineCount('right');
    },
    cmpResult() {
      this.updateLineCount('bottom');
    },
    caseSensitive() {
      this.performSearch();
    },
    wholeWord() {
      this.performSearch();
    },
    useRegex() {
      this.performSearch();
    },
    searchScope() {
      this.matches = [];
      this.currentMatchIndex = -1;
    }
  },
  components: {
    Tooltip,
  }
};
</script>