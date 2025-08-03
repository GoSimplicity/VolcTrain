<template>
  <div class="schedule-calendar">
    <!-- 日历头部控制 -->
    <div class="calendar-header">
      <div class="header-left">
        <Space>
          <Button @click="previousWeek">
            <LeftOutlined />
          </Button>
          <Button @click="nextWeek">
            <RightOutlined />
          </Button>
          <Button @click="goToToday">今天</Button>
        </Space>
        <span class="date-range">{{ dateRangeText }}</span>
      </div>
      
      <div class="header-right">
        <Space>
          <Select v-model:value="viewMode" style="width: 120px">
            <Select.Option value="week">周视图</Select.Option>
            <Select.Option value="day">日视图</Select.Option>
          </Select>
          <Select v-model:value="filterCluster" style="width: 150px" allow-clear>
            <Select.Option value="">全部集群</Select.Option>
            <Select.Option value="cluster-01">Main Cluster</Select.Option>
            <Select.Option value="cluster-02">Training Cluster</Select.Option>
          </Select>
        </Space>
      </div>
    </div>

    <!-- 日历网格 -->
    <div class="calendar-grid">
      <!-- 时间轴 -->
      <div class="time-axis">
        <div class="time-header"></div>
        <div 
          v-for="hour in hours" 
          :key="hour" 
          class="time-slot"
        >
          {{ hour }}:00
        </div>
      </div>

      <!-- 日期列 -->
      <div 
        v-for="(date, dateIndex) in currentDates" 
        :key="date.format('YYYY-MM-DD')"
        class="date-column"
      >
        <!-- 日期头部 -->
        <div class="date-header">
          <div class="date-day">{{ date.format('MM-DD') }}</div>
          <div class="date-weekday">{{ getWeekdayText(date.day()) }}</div>
        </div>

        <!-- 时间槽 -->
        <div class="time-slots">
          <div 
            v-for="hour in hours" 
            :key="hour"
            class="time-slot"
            :class="{ 'current-hour': isCurrentHour(date, hour) }"
          >
            <!-- 调度事件 -->
            <div
              v-for="schedule in getSchedulesForDateTime(date, hour)"
              :key="schedule.id"
              class="schedule-event"
              :class="[
                `status-${schedule.status}`,
                { 'multi-hour': getScheduleDuration(schedule) > 1 }
              ]"
              :style="getEventStyle(schedule)"
              @click="showScheduleDetail(schedule)"
            >
              <div class="event-title">{{ schedule.name }}</div>
              <div class="event-info">
                <span class="event-user">{{ schedule.userName }}</span>
                <span class="event-gpu">{{ schedule.gpuCount }}GPU</span>
              </div>
              <div class="event-time">
                {{ formatTime(schedule.startTime) }} - {{ formatTime(schedule.endTime) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图例 -->
    <div class="calendar-legend">
      <Space>
        <span>状态图例：</span>
        <Tag color="processing">等待中</Tag>
        <Tag color="success">运行中</Tag>
        <Tag color="default">已完成</Tag>
        <Tag color="error">失败</Tag>
        <Tag color="warning">已取消</Tag>
      </Space>
    </div>

    <!-- 调度详情模态框 -->
    <Modal
      v-model:open="detailModalVisible"
      title="调度详情"
      width="600px"
      :footer="null"
    >
      <div v-if="selectedSchedule" class="schedule-detail">
        <Descriptions :column="2" bordered>
          <Descriptions.Item label="任务名称">
            {{ selectedSchedule.name }}
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag :color="getScheduleStatusColor(selectedSchedule.status)">
              {{ getScheduleStatusLabel(selectedSchedule.status) }}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="创建者">
            {{ selectedSchedule.userName }}
          </Descriptions.Item>
          <Descriptions.Item label="集群">
            {{ selectedSchedule.clusterName }}
          </Descriptions.Item>
          <Descriptions.Item label="GPU类型">
            {{ selectedSchedule.gpuType }}
          </Descriptions.Item>
          <Descriptions.Item label="GPU数量">
            {{ selectedSchedule.gpuCount }}
          </Descriptions.Item>
          <Descriptions.Item label="优先级">
            <Tag :color="getPriorityColor(selectedSchedule.priority)">
              {{ selectedSchedule.priority }}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="策略">
            {{ getStrategyText(selectedSchedule.strategy) }}
          </Descriptions.Item>
          <Descriptions.Item label="开始时间" :span="2">
            {{ selectedSchedule.startTime }}
          </Descriptions.Item>
          <Descriptions.Item label="结束时间" :span="2">
            {{ selectedSchedule.endTime }}
          </Descriptions.Item>
          <Descriptions.Item label="持续时间" :span="2">
            {{ formatDuration(selectedSchedule.duration) }}
          </Descriptions.Item>
        </Descriptions>
      </div>
    </Modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import {
  Space,
  Button,
  Select,
  Tag,
  Modal,
  Descriptions,
} from 'ant-design-vue';
import {
  LeftOutlined,
  RightOutlined,
} from '@ant-design/icons-vue';
import dayjs, { type Dayjs } from 'dayjs';
import { formatDuration } from '#/utils/date';

const props = defineProps<{
  schedules: any[];
}>();

// 响应式数据
const currentDate = ref(dayjs());
const viewMode = ref('week');
const filterCluster = ref('');
const detailModalVisible = ref(false);
const selectedSchedule = ref<any>(null);

// 计算属性
const currentDates = computed(() => {
  if (viewMode.value === 'week') {
    const startOfWeek = currentDate.value.startOf('week');
    return Array.from({ length: 7 }, (_, i) => startOfWeek.add(i, 'day'));
  } else {
    return [currentDate.value];
  }
});

const dateRangeText = computed(() => {
  if (viewMode.value === 'week') {
    const start = currentDates.value[0];
    const end = currentDates.value[6];
    return `${start.format('YYYY年MM月DD日')} - ${end.format('YYYY年MM月DD日')}`;
  } else {
    return currentDate.value.format('YYYY年MM月DD日');
  }
});

const filteredSchedules = computed(() => {
  let schedules = props.schedules || [];
  
  if (filterCluster.value) {
    schedules = schedules.filter(s => s.clusterId === filterCluster.value);
  }
  
  return schedules;
});

// 小时数组 (0-23)
const hours = Array.from({ length: 24 }, (_, i) => i);

// 工具方法
const getWeekdayText = (day: number) => {
  const weekdays = ['日', '一', '二', '三', '四', '五', '六'];
  return weekdays[day];
};

const isCurrentHour = (date: Dayjs, hour: number) => {
  const now = dayjs();
  return date.isSame(now, 'day') && hour === now.hour();
};

const getSchedulesForDateTime = (date: Dayjs, hour: number) => {
  return filteredSchedules.value.filter(schedule => {
    const startTime = dayjs(schedule.startTime);
    const endTime = dayjs(schedule.endTime);
    const currentDateTime = date.hour(hour);
    
    return currentDateTime.isBetween(startTime, endTime, 'hour', '[)');
  });
};

const getScheduleDuration = (schedule: any) => {
  const startTime = dayjs(schedule.startTime);
  const endTime = dayjs(schedule.endTime);
  return endTime.diff(startTime, 'hour');
};

const getEventStyle = (schedule: any) => {
  const startTime = dayjs(schedule.startTime);
  const endTime = dayjs(schedule.endTime);
  const duration = endTime.diff(startTime, 'hour');
  
  return {
    height: `${Math.max(duration * 40 - 2, 38)}px`, // 每小时40px高度
    zIndex: 10,
  };
};

const formatTime = (timeStr: string) => {
  return dayjs(timeStr).format('HH:mm');
};

const getScheduleStatusColor = (status: string) => {
  const colors = {
    pending: 'processing',
    running: 'success',
    completed: 'default',
    failed: 'error',
    cancelled: 'warning',
  };
  return colors[status as keyof typeof colors] || 'default';
};

const getScheduleStatusLabel = (status: string) => {
  const labels = {
    pending: '等待中',
    running: '运行中',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消',
  };
  return labels[status as keyof typeof labels] || status;
};

const getPriorityColor = (priority: number) => {
  if (priority >= 8) return 'red';
  if (priority >= 6) return 'orange';
  if (priority >= 4) return 'blue';
  return 'default';
};

const getStrategyText = (strategy: string) => {
  const strategies = {
    best_fit: '最佳适配',
    balanced: '负载均衡',
    gpu_optimized: 'GPU优化',
    round_robin: '轮询',
    anti_affinity: '反亲和性',
    memory_optimized: '内存优化',
  };
  return strategies[strategy as keyof typeof strategies] || strategy;
};

// 事件处理
const previousWeek = () => {
  currentDate.value = currentDate.value.subtract(1, viewMode.value);
};

const nextWeek = () => {
  currentDate.value = currentDate.value.add(1, viewMode.value);
};

const goToToday = () => {
  currentDate.value = dayjs();
};

const showScheduleDetail = (schedule: any) => {
  selectedSchedule.value = schedule;
  detailModalVisible.value = true;
};

// 初始化
onMounted(() => {
  // 可以在这里加载更多数据
});
</script>

<style scoped lang="scss">
.schedule-calendar {
  width: 100%;
  background: white;
  border-radius: 8px;
  overflow: hidden;
}

.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  background: #fafafa;
  
  .header-left {
    display: flex;
    align-items: center;
    gap: 16px;
    
    .date-range {
      font-weight: 600;
      font-size: 16px;
      color: #1890ff;
    }
  }
}

.calendar-grid {
  display: flex;
  min-height: 600px;
  overflow-x: auto;
}

.time-axis {
  flex-shrink: 0;
  width: 80px;
  border-right: 1px solid #f0f0f0;
  
  .time-header {
    height: 60px;
    border-bottom: 1px solid #f0f0f0;
    background: #fafafa;
  }
  
  .time-slot {
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    color: #666;
    border-bottom: 1px solid #f5f5f5;
  }
}

.date-column {
  flex: 1;
  min-width: 150px;
  border-right: 1px solid #f0f0f0;
  
  &:last-child {
    border-right: none;
  }
  
  .date-header {
    height: 60px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border-bottom: 1px solid #f0f0f0;
    background: #fafafa;
    
    .date-day {
      font-weight: 600;
      font-size: 14px;
      color: #333;
    }
    
    .date-weekday {
      font-size: 12px;
      color: #666;
      margin-top: 2px;
    }
  }
  
  .time-slots {
    position: relative;
    
    .time-slot {
      height: 40px;
      border-bottom: 1px solid #f5f5f5;
      position: relative;
      
      &.current-hour {
        background-color: rgba(24, 144, 255, 0.05);
      }
    }
  }
}

.schedule-event {
  position: absolute;
  left: 2px;
  right: 2px;
  top: 1px;
  padding: 4px 6px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 11px;
  line-height: 1.2;
  transition: all 0.2s;
  
  &:hover {
    transform: scale(1.02);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  }
  
  &.status-pending {
    background: #e6f7ff;
    border: 1px solid #91d5ff;
    color: #1890ff;
  }
  
  &.status-running {
    background: #f6ffed;
    border: 1px solid #b7eb8f;
    color: #52c41a;
  }
  
  &.status-completed {
    background: #f5f5f5;
    border: 1px solid #d9d9d9;
    color: #666;
  }
  
  &.status-failed {
    background: #fff2f0;
    border: 1px solid #ffccc7;
    color: #ff4d4f;
  }
  
  &.status-cancelled {
    background: #fffbe6;
    border: 1px solid #ffe58f;
    color: #faad14;
  }
  
  .event-title {
    font-weight: 500;
    margin-bottom: 2px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .event-info {
    display: flex;
    justify-content: space-between;
    margin-bottom: 2px;
    
    .event-user,
    .event-gpu {
      font-size: 10px;
      opacity: 0.8;
    }
  }
  
  .event-time {
    font-size: 10px;
    opacity: 0.7;
  }
}

.calendar-legend {
  padding: 12px 20px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
}

.schedule-detail {
  .ant-descriptions {
    margin-top: 16px;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .calendar-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
    
    .header-left,
    .header-right {
      justify-content: center;
    }
  }
  
  .date-column {
    min-width: 120px;
  }
  
  .time-axis {
    width: 60px;
  }
  
  .schedule-event {
    font-size: 10px;
    padding: 2px 4px;
  }
}
</style>