/**
 * 用户体验和交互优化工具
 */
import { ref, reactive, computed, nextTick } from 'vue';
import { message, notification, Modal } from 'ant-design-vue';

// 加载状态管理
export interface LoadingState {
  global: boolean;
  components: Record<string, boolean>;
}

export function useLoadingState() {
  const loadingState = reactive<LoadingState>({
    global: false,
    components: {}
  });

  const setGlobalLoading = (loading: boolean) => {
    loadingState.global = loading;
  };

  const setComponentLoading = (component: string, loading: boolean) => {
    loadingState.components[component] = loading;
  };

  const isLoading = computed(() => {
    return loadingState.global || Object.values(loadingState.components).some(Boolean);
  });

  return {
    loadingState,
    setGlobalLoading,
    setComponentLoading,
    isLoading
  };
}

// 通知管理
export interface NotificationConfig {
  type: 'success' | 'info' | 'warning' | 'error';
  title: string;
  message: string;
  duration?: number;
  showProgress?: boolean;
  actions?: Array<{
    label: string;
    action: () => void;
    type?: 'primary' | 'default';
  }>;
}

export function useNotification() {
  const notifications = ref<Array<NotificationConfig & { id: string }>>([]);

  const showNotification = (config: NotificationConfig) => {
    const id = `notification-${Date.now()}-${Math.random()}`;
    
    notifications.value.push({ ...config, id });

    notification[config.type]({
      message: config.title,
      description: config.message,
      duration: config.duration || 4.5,
      key: id,
      onClick: () => {
        removeNotification(id);
      }
    });

    // 自动移除
    if (config.duration !== 0) {
      setTimeout(() => {
        removeNotification(id);
      }, (config.duration || 4.5) * 1000);
    }

    return id;
  };

  const removeNotification = (id: string) => {
    const index = notifications.value.findIndex(n => n.id === id);
    if (index > -1) {
      notifications.value.splice(index, 1);
    }
    notification.destroy();
  };

  const clearAllNotifications = () => {
    notifications.value = [];
    notification.destroy();
  };

  // 预设通知类型
  const showSuccess = (title: string, message: string, duration?: number) => {
    return showNotification({ type: 'success', title, message, duration });
  };

  const showError = (title: string, message: string, duration = 0) => {
    return showNotification({ type: 'error', title, message, duration });
  };

  const showWarning = (title: string, message: string, duration?: number) => {
    return showNotification({ type: 'warning', title, message, duration });
  };

  const showInfo = (title: string, message: string, duration?: number) => {
    return showNotification({ type: 'info', title, message, duration });
  };

  return {
    notifications,
    showNotification,
    removeNotification,
    clearAllNotifications,
    showSuccess,
    showError,
    showWarning,
    showInfo
  };
}

// 键盘快捷键
export interface ShortcutConfig {
  key: string;
  ctrl?: boolean;
  alt?: boolean;
  shift?: boolean;
  description: string;
  action: () => void;
  global?: boolean;
}

export function useKeyboardShortcuts() {
  const shortcuts = ref<ShortcutConfig[]>([]);
  const isEnabled = ref(true);

  const addShortcut = (config: ShortcutConfig) => {
    shortcuts.value.push(config);
  };

  const removeShortcut = (key: string) => {
    const index = shortcuts.value.findIndex(s => s.key === key);
    if (index > -1) {
      shortcuts.value.splice(index, 1);
    }
  };

  const handleKeyDown = (event: KeyboardEvent) => {
    if (!isEnabled.value) return;

    const pressedKey = event.key.toLowerCase();
    const hasCtrl = event.ctrlKey || event.metaKey;
    const hasAlt = event.altKey;
    const hasShift = event.shiftKey;

    const matchingShortcut = shortcuts.value.find(shortcut => {
      return (
        shortcut.key.toLowerCase() === pressedKey &&
        !!shortcut.ctrl === hasCtrl &&
        !!shortcut.alt === hasAlt &&
        !!shortcut.shift === hasShift
      );
    });

    if (matchingShortcut) {
      event.preventDefault();
      matchingShortcut.action();
    }
  };

  const enable = () => {
    isEnabled.value = true;
    document.addEventListener('keydown', handleKeyDown);
  };

  const disable = () => {
    isEnabled.value = false;
    document.removeEventListener('keydown', handleKeyDown);
  };

  // 默认快捷键
  const setupDefaultShortcuts = () => {
    addShortcut({
      key: 'f5',
      description: '刷新页面',
      action: () => {
        window.location.reload();
      }
    });

    addShortcut({
      key: 's',
      ctrl: true,
      description: '保存',
      action: () => {
        message.info('保存快捷键触发');
      }
    });

    addShortcut({
      key: '/',
      description: '显示快捷键帮助',
      action: () => {
        showShortcutHelp();
      }
    });
  };

  const showShortcutHelp = () => {
    const content = shortcuts.value.map(shortcut => {
      const keys = [];
      if (shortcut.ctrl) keys.push('Ctrl');
      if (shortcut.alt) keys.push('Alt');
      if (shortcut.shift) keys.push('Shift');
      keys.push(shortcut.key.toUpperCase());
      
      return `${keys.join(' + ')}: ${shortcut.description}`;
    }).join('\n');

    Modal.info({
      title: '键盘快捷键',
      content: content.split('\n').map(line => `<div>${line}</div>`).join(''),
      width: 400
    });
  };

  return {
    shortcuts,
    addShortcut,
    removeShortcut,
    enable,
    disable,
    setupDefaultShortcuts,
    showShortcutHelp
  };
}

// 防抖和节流
export function useDebounce<T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: NodeJS.Timeout;

  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => fn(...args), delay);
  };
}

export function useThrottle<T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let lastCall = 0;

  return (...args: Parameters<T>) => {
    const now = Date.now();
    if (now - lastCall >= delay) {
      lastCall = now;
      fn(...args);
    }
  };
}

// 错误处理
export interface ErrorConfig {
  title?: string;
  description?: string;
  showNotification?: boolean;
  showModal?: boolean;
  logToConsole?: boolean;
  reportToServer?: boolean;
}

export function useErrorHandler() {
  const handleError = (error: Error | string, config: ErrorConfig = {}) => {
    const errorMessage = typeof error === 'string' ? error : error.message;
    const stackTrace = typeof error === 'string' ? undefined : error.stack;

    const {
      title = '发生错误',
      description = errorMessage,
      showNotification = true,
      showModal = false,
      logToConsole = true,
      reportToServer = false
    } = config;

    // 记录到控制台
    if (logToConsole) {
      console.error('应用错误:', {
        message: errorMessage,
        stack: stackTrace,
        config
      });
    }

    // 显示通知
    if (showNotification) {
      notification.error({
        message: title,
        description,
        duration: 0
      });
    }

    // 显示模态框
    if (showModal) {
      Modal.error({
        title,
        content: description,
        width: 500
      });
    }

    // 上报到服务器
    if (reportToServer) {
      // 这里可以集成错误上报服务
      console.log('上报错误到服务器:', {
        message: errorMessage,
        stack: stackTrace,
        userAgent: navigator.userAgent,
        url: window.location.href,
        timestamp: new Date().toISOString()
      });
    }
  };

  const handleAsyncError = async (
    asyncFn: () => Promise<any>,
    config: ErrorConfig = {}
  ) => {
    try {
      return await asyncFn();
    } catch (error) {
      handleError(error as Error, config);
      throw error;
    }
  };

  return {
    handleError,
    handleAsyncError
  };
}

// 本地存储管理
export function useLocalStorage() {
  const setItem = (key: string, value: any) => {
    try {
      const serializedValue = JSON.stringify(value);
      localStorage.setItem(key, serializedValue);
    } catch (error) {
      console.error('保存到本地存储失败:', error);
    }
  };

  const getItem = <T = any>(key: string, defaultValue?: T): T | null => {
    try {
      const item = localStorage.getItem(key);
      if (item === null) {
        return defaultValue || null;
      }
      return JSON.parse(item);
    } catch (error) {
      console.error('从本地存储读取失败:', error);
      return defaultValue || null;
    }
  };

  const removeItem = (key: string) => {
    try {
      localStorage.removeItem(key);
    } catch (error) {
      console.error('从本地存储删除失败:', error);
    }
  };

  const clear = () => {
    try {
      localStorage.clear();
    } catch (error) {
      console.error('清空本地存储失败:', error);
    }
  };

  return {
    setItem,
    getItem,
    removeItem,
    clear
  };
}

// 用户偏好设置
export interface UserPreferences {
  theme: 'light' | 'dark' | 'auto';
  language: string;
  autoRefresh: boolean;
  refreshInterval: number;
  notifications: {
    enabled: boolean;
    types: string[];
  };
  layout: {
    sidebarCollapsed: boolean;
    density: 'compact' | 'default' | 'comfortable';
  };
}

export function useUserPreferences() {
  const storage = useLocalStorage();
  const PREFERENCES_KEY = 'user-preferences';

  const defaultPreferences: UserPreferences = {
    theme: 'light',
    language: 'zh-CN',
    autoRefresh: true,
    refreshInterval: 30,
    notifications: {
      enabled: true,
      types: ['error', 'warning', 'success']
    },
    layout: {
      sidebarCollapsed: false,
      density: 'default'
    }
  };

  const preferences = ref<UserPreferences>(
    storage.getItem(PREFERENCES_KEY, defaultPreferences) ?? defaultPreferences
  );

  const updatePreferences = (updates: Partial<UserPreferences>) => {
    preferences.value = { ...preferences.value, ...updates };
    storage.setItem(PREFERENCES_KEY, preferences.value);
  };

  const resetPreferences = () => {
    preferences.value = { ...defaultPreferences };
    storage.setItem(PREFERENCES_KEY, preferences.value);
  };

  return {
    preferences,
    updatePreferences,
    resetPreferences
  };
}

// 无障碍功能
export function useAccessibility() {
  const isHighContrast = ref(false);
  const isLargeText = ref(false);
  const isReducedMotion = ref(false);

  const toggleHighContrast = () => {
    isHighContrast.value = !isHighContrast.value;
    document.body.classList.toggle('high-contrast', isHighContrast.value);
  };

  const toggleLargeText = () => {
    isLargeText.value = !isLargeText.value;
    document.body.classList.toggle('large-text', isLargeText.value);
  };

  const toggleReducedMotion = () => {
    isReducedMotion.value = !isReducedMotion.value;
    document.body.classList.toggle('reduced-motion', isReducedMotion.value);
  };

  const announceToScreenReader = (message: string) => {
    const announcement = document.createElement('div');
    announcement.setAttribute('aria-live', 'polite');
    announcement.setAttribute('aria-atomic', 'true');
    announcement.setAttribute('class', 'sr-only');
    announcement.textContent = message;
    
    document.body.appendChild(announcement);
    
    setTimeout(() => {
      document.body.removeChild(announcement);
    }, 1000);
  };

  return {
    isHighContrast,
    isLargeText,
    isReducedMotion,
    toggleHighContrast,
    toggleLargeText,
    toggleReducedMotion,
    announceToScreenReader
  };
}

// 响应式设计辅助
export function useResponsive() {
  const breakpoints = {
    xs: 480,
    sm: 576,
    md: 768,
    lg: 992,
    xl: 1200,
    xxl: 1600
  };

  const currentBreakpoint = ref('lg');
  const isMobile = computed(() => ['xs', 'sm'].includes(currentBreakpoint.value));
  const isTablet = computed(() => currentBreakpoint.value === 'md');
  const isDesktop = computed(() => ['lg', 'xl', 'xxl'].includes(currentBreakpoint.value));

  const updateBreakpoint = () => {
    const width = window.innerWidth;
    
    if (width < breakpoints.xs) {
      currentBreakpoint.value = 'xs';
    } else if (width < breakpoints.sm) {
      currentBreakpoint.value = 'sm';
    } else if (width < breakpoints.md) {
      currentBreakpoint.value = 'md';
    } else if (width < breakpoints.lg) {
      currentBreakpoint.value = 'lg';
    } else if (width < breakpoints.xl) {
      currentBreakpoint.value = 'xl';
    } else {
      currentBreakpoint.value = 'xxl';
    }
  };

  // 初始化
  updateBreakpoint();
  window.addEventListener('resize', updateBreakpoint);

  return {
    currentBreakpoint,
    isMobile,
    isTablet,
    isDesktop,
    breakpoints
  };
}

// 表单验证辅助
export interface ValidationRule {
  required?: boolean;
  min?: number;
  max?: number;
  pattern?: RegExp;
  custom?: (value: any) => string | null;
  message?: string;
}

export function useFormValidation() {
  const validateField = (value: any, rules: ValidationRule[]): string | null => {
    for (const rule of rules) {
      if (rule.required && (value === null || value === undefined || value === '')) {
        return rule.message || '此字段为必填项';
      }

      if (rule.min && value.length < rule.min) {
        return rule.message || `最少需要${rule.min}个字符`;
      }

      if (rule.max && value.length > rule.max) {
        return rule.message || `最多允许${rule.max}个字符`;
      }

      if (rule.pattern && !rule.pattern.test(value)) {
        return rule.message || '格式不正确';
      }

      if (rule.custom) {
        const customResult = rule.custom(value);
        if (customResult) {
          return customResult;
        }
      }
    }

    return null;
  };

  const validateForm = (formData: Record<string, any>, rules: Record<string, ValidationRule[]>) => {
    const errors: Record<string, string> = {};
    
    for (const [field, fieldRules] of Object.entries(rules)) {
      const error = validateField(formData[field], fieldRules);
      if (error) {
        errors[field] = error;
      }
    }

    return {
      isValid: Object.keys(errors).length === 0,
      errors
    };
  };

  // 常用验证规则
  const commonRules = {
    required: { required: true },
    email: {
      pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
      message: '请输入有效的邮箱地址'
    },
    phone: {
      pattern: /^1[3-9]\d{9}$/,
      message: '请输入有效的手机号码'
    },
    password: {
      min: 6,
      message: '密码至少需要6个字符'
    }
  };

  return {
    validateField,
    validateForm,
    commonRules
  };
}

// 动画辅助
export function useAnimation() {
  const fadeIn = (element: HTMLElement, duration = 300) => {
    element.style.opacity = '0';
    element.style.transition = `opacity ${duration}ms ease-in-out`;
    
    nextTick(() => {
      element.style.opacity = '1';
    });
  };

  const fadeOut = (element: HTMLElement, duration = 300) => {
    element.style.transition = `opacity ${duration}ms ease-in-out`;
    element.style.opacity = '0';
    
    setTimeout(() => {
      element.style.display = 'none';
    }, duration);
  };

  const slideIn = (element: HTMLElement, direction = 'left', duration = 300) => {
    const translateMap = {
      left: 'translateX(-100%)',
      right: 'translateX(100%)',
      up: 'translateY(-100%)',
      down: 'translateY(100%)'
    };

    element.style.transform = translateMap[direction as keyof typeof translateMap];
    element.style.transition = `transform ${duration}ms ease-in-out`;
    
    nextTick(() => {
      element.style.transform = 'translate(0)';
    });
  };

  const bounce = (element: HTMLElement) => {
    element.style.animation = 'bounce 0.6s ease-in-out';
    
    setTimeout(() => {
      element.style.animation = '';
    }, 600);
  };

  return {
    fadeIn,
    fadeOut,
    slideIn,
    bounce
  };
}

// 性能监控
export function usePerformanceMonitor() {
  const metrics = ref({
    loadTime: 0,
    renderTime: 0,
    apiCalls: 0,
    errorCount: 0
  });

  const startTime = Date.now();

  const recordLoadTime = () => {
    metrics.value.loadTime = Date.now() - startTime;
  };

  const recordRenderTime = (componentName: string, renderTime: number) => {
    console.log(`组件 ${componentName} 渲染时间: ${renderTime}ms`);
  };

  const recordApiCall = () => {
    metrics.value.apiCalls++;
  };

  const recordError = () => {
    metrics.value.errorCount++;
  };

  const getPerformanceReport = () => {
    return {
      ...metrics.value,
      timestamp: new Date().toISOString(),
      userAgent: navigator.userAgent,
      url: window.location.href
    };
  };

  return {
    metrics,
    recordLoadTime,
    recordRenderTime,
    recordApiCall,
    recordError,
    getPerformanceReport
  };
}

// 用户行为追踪
export interface UserAction {
  type: string;
  target: string;
  timestamp: number;
  data?: any;
}

export function useUserTracking() {
  const actions = ref<UserAction[]>([]);

  const trackAction = (type: string, target: string, data?: any) => {
    const action: UserAction = {
      type,
      target,
      timestamp: Date.now(),
      data
    };

    actions.value.push(action);

    // 限制存储的动作数量
    if (actions.value.length > 1000) {
      actions.value = actions.value.slice(-500);
    }

    console.log('用户动作:', action);
  };

  const trackClick = (target: string, data?: any) => {
    trackAction('click', target, data);
  };

  const trackView = (target: string, data?: any) => {
    trackAction('view', target, data);
  };

  const trackSearch = (query: string, results?: number) => {
    trackAction('search', 'search-box', { query, results });
  };

  const getActionHistory = (type?: string, limit = 100) => {
    let filtered = actions.value;
    
    if (type) {
      filtered = filtered.filter(action => action.type === type);
    }

    return filtered.slice(-limit);
  };

  return {
    actions,
    trackAction,
    trackClick,
    trackView,
    trackSearch,
    getActionHistory
  };
}

