<template>
  <div class="log-item" :class="`log-level-${log.level}`">
    <div class="log-content">
      <span class="log-timestamp">{{ formatTimestamp(log.timestamp) }}</span>
      <span class="log-level" :class="`level-${log.level}`">{{ log.level.toUpperCase() }}</span>
      <span class="log-source">[{{ log.source }}]</span>
      <span class="log-message" v-html="highlightText(log.message, searchKeyword)"></span>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { format } from 'date-fns';
import type { TrainingApi } from '#/types/api';

interface Props {
  log: TrainingApi.TrainingLog;
  searchKeyword?: string;
}

const props = defineProps<Props>();

/**
 * 格式化时间戳
 */
const formatTimestamp = (timestamp: string) => {
  return format(new Date(timestamp), 'HH:mm:ss.SSS');
};

/**
 * 高亮搜索关键词
 */
const highlightText = (text: string, keyword?: string) => {
  if (!keyword) return text;
  
  const regex = new RegExp(`(${keyword})`, 'gi');
  return text.replace(regex, '<mark class="highlight">$1</mark>');
};
</script>

<style lang="scss" scoped>
.log-item {
  padding: 4px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  word-wrap: break-word;
  transition: background-color 0.2s ease;

  &:hover {
    background-color: rgba(255, 255, 255, 0.05);
  }

  .log-content {
    display: flex;
    align-items: baseline;
    font-family: 'JetBrains Mono', 'Consolas', monospace;
    font-size: 13px;
    line-height: 1.5;

    .log-timestamp {
      color: #888;
      margin-right: 8px;
      font-size: 11px;
      min-width: 80px;
    }

    .log-level {
      margin-right: 8px;
      padding: 2px 6px;
      border-radius: 3px;
      font-size: 10px;
      font-weight: 600;
      min-width: 48px;
      text-align: center;

      &.level-error {
        background-color: #ff4d4f;
        color: #fff;
      }

      &.level-warn {
        background-color: #faad14;
        color: #fff;
      }

      &.level-info {
        background-color: #1890ff;
        color: #fff;
      }

      &.level-debug {
        background-color: #52c41a;
        color: #fff;
      }
    }

    .log-source {
      color: #1890ff;
      margin-right: 8px;
      font-size: 11px;
      min-width: 100px;
    }

    .log-message {
      flex: 1;
      color: #ffffff;

      :deep(.highlight) {
        background-color: #faad14;
        color: #000;
        padding: 1px 2px;
        border-radius: 2px;
      }
    }
  }

  // 不同级别的背景色
  &.log-level-error {
    background-color: rgba(255, 77, 79, 0.1);
    border-left: 3px solid #ff4d4f;
  }

  &.log-level-warn {
    background-color: rgba(250, 173, 20, 0.1);
    border-left: 3px solid #faad14;
  }

  &.log-level-info {
    background-color: rgba(24, 144, 255, 0.1);
    border-left: 3px solid #1890ff;
  }

  &.log-level-debug {
    background-color: rgba(82, 196, 26, 0.1);
    border-left: 3px solid #52c41a;
  }
}
</style>