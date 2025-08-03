import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

// 扩展dayjs插件
dayjs.extend(relativeTime);

/**
 * 日期时间格式化工具函数
 */

// 格式化日期时间
export const formatDateTime = (date: string | Date | undefined | null, format = 'YYYY-MM-DD HH:mm:ss'): string => {
  if (!date) return '';
  return dayjs(date).format(format);
};

// 格式化日期
export const formatDate = (date: string | Date | undefined | null, format = 'YYYY-MM-DD'): string => {
  if (!date) return '';
  return dayjs(date).format(format);
};

// 格式化时间
export const formatTime = (date: string | Date | undefined | null, format = 'HH:mm:ss'): string => {
  if (!date) return '';
  return dayjs(date).format(format);
};

// 相对时间（多久前）
export const formatRelativeTime = (date: string | Date | undefined | null): string => {
  if (!date) return '';
  return dayjs(date).fromNow();
};

// 是否为今天
export const isToday = (date: string | Date): boolean => {
  return dayjs(date).isSame(dayjs(), 'day');
};

// 是否为本周
export const isThisWeek = (date: string | Date): boolean => {
  return dayjs(date).isSame(dayjs(), 'week');
};

// 是否为本月
export const isThisMonth = (date: string | Date): boolean => {
  return dayjs(date).isSame(dayjs(), 'month');
};

// 计算时间差（天数）
export const getDaysDiff = (date1: string | Date, date2?: string | Date): number => {
  const target = date2 ? dayjs(date2) : dayjs();
  return target.diff(dayjs(date1), 'day');
};

// 计算时间差（小时）
export const getHoursDiff = (date1: string | Date, date2?: string | Date): number => {
  const target = date2 ? dayjs(date2) : dayjs();
  return target.diff(dayjs(date1), 'hour');
};

// 格式化时长（秒转换为可读格式）
export const formatDuration = (seconds: number): string => {
  if (seconds < 60) {
    return `${seconds}秒`;
  } else if (seconds < 3600) {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return remainingSeconds > 0 ? `${minutes}分${remainingSeconds}秒` : `${minutes}分钟`;
  } else if (seconds < 86400) {
    const hours = Math.floor(seconds / 3600);
    const remainingMinutes = Math.floor((seconds % 3600) / 60);
    return remainingMinutes > 0 ? `${hours}小时${remainingMinutes}分钟` : `${hours}小时`;
  } else {
    const days = Math.floor(seconds / 86400);
    const remainingHours = Math.floor((seconds % 86400) / 3600);
    return remainingHours > 0 ? `${days}天${remainingHours}小时` : `${days}天`;
  }
};

// 获取当前时间戳
export const getCurrentTimestamp = (): string => {
  return dayjs().format('YYYY-MM-DD HH:mm:ss');
};

// 格式化文件大小
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`;
};