<script lang="ts" setup>
import { ref, onMounted, h } from 'vue';
import {
  Card,
  Table,
  Tag,
  Space,
  Progress,
  Button,
  Tabs,
  Row,
  Col,
  Statistic,
  Timeline,
  Divider,
} from 'ant-design-vue';

defineOptions({ name: 'TaskMonitor' });

// 任务统计信息
const taskStats = ref({
  totalTasks: 28,
  runningTasks: 12,
  waitingTasks: 8,
  completedTasks: 6,
  failedTasks: 2,
  avgDuration: '3小时24分钟',
  avgGpuUsage: 78,
  avgMemoryUsage: 62,
});

// 任务列表数据
const tasksData = ref([
  {
    id: '1',
    name: 'BERT微调任务1',
    type: '训练任务',
    status: '运行中',
    progress: 45,
    gpuUsage: 85,
    memoryUsage: 72,
    startTime: '2023-11-15 08:30',
    duration: '3小时25分钟',
    node: 'node-01',
    user: '张三',
  },
  {
    id: '2',
    name: 'ResNet图像分类',
    type: '训练任务',
    status: '运行中',
    progress: 78,
    gpuUsage: 92,
    memoryUsage: 80,
    startTime: '2023-11-15 06:15',
    duration: '5小时40分钟',
    node: 'node-02',
    user: '李四',
  },
  {
    id: '3',
    name: '数据预处理任务',
    type: '数据处理',
    status: '等待中',
    progress: 0,
    gpuUsage: 0,
    memoryUsage: 0,
    startTime: '-',
    duration: '-',
    node: '-',
    user: '王五',
  },
  {
    id: '4',
    name: 'GPT模型训练',
    type: '训练任务',
    status: '已完成',
    progress: 100,
    gpuUsage: 0,
    memoryUsage: 0,
    startTime: '2023-11-14 12:00',
    duration: '18小时30分钟',
    node: 'node-03',
    user: '张三',
  },
  {
    id: '5',
    name: '模型评估任务',
    type: '评估任务',
    status: '失败',
    progress: 65,
    gpuUsage: 0,
    memoryUsage: 0,
    startTime: '2023-11-15 01:20',
    duration: '1小时15分钟',
    node: 'node-01',
    user: '李四',
  },
]);

// 任务表格列配置
const columns = [
  {
    title: '任务ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '任务名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '任务类型',
    dataIndex: 'type',
    key: 'type',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      const color =
        text === '运行中'
          ? 'blue'
          : text === '等待中'
            ? 'orange'
            : text === '已完成'
              ? 'green'
              : text === '失败'
                ? 'red'
                : 'default';
      return h(Tag, { color }, () => text);
    },
  },
  {
    title: '进度',
    dataIndex: 'progress',
    key: 'progress',
    customRender: ({ text }) => {
      return h(Progress, { percent: text });
    },
  },
  {
    title: 'GPU使用率',
    dataIndex: 'gpuUsage',
    key: 'gpuUsage',
    customRender: ({ text }) => {
      if (text === 0) return '-';
      const color = text > 90 ? 'red' : text > 70 ? 'orange' : 'green';
      return h(Progress, { percent: text, size: 'small', strokeColor: color });
    },
  },
  {
    title: '内存使用率',
    dataIndex: 'memoryUsage',
    key: 'memoryUsage',
    customRender: ({ text }) => {
      if (text === 0) return '-';
      const color = text > 90 ? 'red' : text > 70 ? 'orange' : 'green';
      return h(Progress, { percent: text, size: 'small', strokeColor: color });
    },
  },
  {
    title: '开始时间',
    dataIndex: 'startTime',
    key: 'startTime',
  },
  {
    title: '运行时长',
    dataIndex: 'duration',
    key: 'duration',
  },
  {
    title: '运行节点',
    dataIndex: 'node',
    key: 'node',
  },
  {
    title: '创建用户',
    dataIndex: 'user',
    key: 'user',
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }) =>
      h(
        Space,
        { size: 'middle' },
        [
          h('a', { onClick: () => handleViewDetails(record) }, '详情'),
          record.status === '运行中' ? h('a', {}, '暂停') : null,
          record.status === '等待中' ? h('a', {}, '取消') : null,
          record.status === '失败' ? h('a', {}, '重试') : null,
        ].filter(Boolean),
      ),
  },
];

const activeKey = ref('1');
const selectedTask = ref(null);

const handleViewDetails = (task) => {
  selectedTask.value = task;
};

// 假设的任务日志数据
const taskLogs = ref([
  {
    time: '08:30:15',
    content: '任务启动',
  },
  {
    time: '08:31:05',
    content: '加载数据集',
  },
  {
    time: '08:35:22',
    content: '开始训练',
  },
  {
    time: '09:15:43',
    content: '完成第一轮训练，损失: 0.567',
  },
  {
    time: '10:00:12',
    content: '完成第二轮训练，损失: 0.432',
  },
  {
    time: '10:45:35',
    content: '完成第三轮训练，损失: 0.385',
  },
  {
    time: '11:30:08',
    content: '完成第四轮训练，损失: 0.341',
  },
  {
    time: '11:55:24',
    content: '当前训练中，损失: 0.328',
  },
]);

onMounted(() => {
  // 模拟选中第一个任务
  selectedTask.value = tasksData.value[0];
});
</script>

<template>
  <div class="task-monitor-container">
    <Row :gutter="16" class="stats-row">
      <Col :span="4">
        <Card>
          <Statistic title="总任务数" :value="taskStats.totalTasks" />
        </Card>
      </Col>
      <Col :span="5">
        <Card>
          <Statistic title="运行中任务" :value="taskStats.runningTasks" />
        </Card>
      </Col>
      <Col :span="5">
        <Card>
          <Statistic title="等待中任务" :value="taskStats.waitingTasks" />
        </Card>
      </Col>
      <Col :span="5">
        <Card>
          <Statistic title="已完成任务" :value="taskStats.completedTasks" />
        </Card>
      </Col>
      <Col :span="5">
        <Card>
          <Statistic title="失败任务" :value="taskStats.failedTasks" />
        </Card>
      </Col>
    </Row>

    <Card title="任务监控">
      <Tabs v-model:activeKey="activeKey">
        <Tabs.TabPane key="1" tab="所有任务">
          <Table
            :columns="columns"
            :dataSource="tasksData"
            rowKey="id"
            :pagination="{ pageSize: 10 }"
          />
        </Tabs.TabPane>
        <Tabs.TabPane key="2" tab="运行中">
          <Table
            :columns="columns"
            :dataSource="tasksData.filter((task) => task.status === '运行中')"
            rowKey="id"
            :pagination="{ pageSize: 10 }"
          />
        </Tabs.TabPane>
        <Tabs.TabPane key="3" tab="等待中">
          <Table
            :columns="columns"
            :dataSource="tasksData.filter((task) => task.status === '等待中')"
            rowKey="id"
            :pagination="{ pageSize: 10 }"
          />
        </Tabs.TabPane>
        <Tabs.TabPane key="4" tab="已完成">
          <Table
            :columns="columns"
            :dataSource="tasksData.filter((task) => task.status === '已完成')"
            rowKey="id"
            :pagination="{ pageSize: 10 }"
          />
        </Tabs.TabPane>
        <Tabs.TabPane key="5" tab="失败">
          <Table
            :columns="columns"
            :dataSource="tasksData.filter((task) => task.status === '失败')"
            rowKey="id"
            :pagination="{ pageSize: 10 }"
          />
        </Tabs.TabPane>
      </Tabs>
    </Card>

    <Divider />

    <Card title="任务详情" v-if="selectedTask">
      <Row :gutter="16">
        <Col :span="8">
          <Card title="基本信息" size="small">
            <p><strong>任务名称:</strong> {{ selectedTask.name }}</p>
            <p><strong>任务类型:</strong> {{ selectedTask.type }}</p>
            <p><strong>状态:</strong> {{ selectedTask.status }}</p>
            <p><strong>开始时间:</strong> {{ selectedTask.startTime }}</p>
            <p><strong>运行时长:</strong> {{ selectedTask.duration }}</p>
            <p><strong>运行节点:</strong> {{ selectedTask.node }}</p>
            <p><strong>创建用户:</strong> {{ selectedTask.user }}</p>
          </Card>
        </Col>
        <Col :span="8">
          <Card title="资源使用" size="small">
            <p><strong>进度:</strong></p>
            <Progress :percent="selectedTask.progress" status="active" />
            <p><strong>GPU使用率:</strong></p>
            <Progress
              :percent="selectedTask.gpuUsage"
              :status="selectedTask.status === '运行中' ? 'active' : 'normal'"
              :strokeColor="selectedTask.gpuUsage > 90 ? '#f5222d' : '#1890ff'"
            />
            <p><strong>内存使用率:</strong></p>
            <Progress
              :percent="selectedTask.memoryUsage"
              :status="selectedTask.status === '运行中' ? 'active' : 'normal'"
              :strokeColor="
                selectedTask.memoryUsage > 90 ? '#f5222d' : '#1890ff'
              "
            />
          </Card>
        </Col>
        <Col :span="8">
          <Card title="任务日志" size="small">
            <Timeline>
              <Timeline.Item v-for="(log, index) in taskLogs" :key="index">
                <p>{{ log.time }} - {{ log.content }}</p>
              </Timeline.Item>
            </Timeline>
          </Card>
        </Col>
      </Row>
      <div style="margin-top: 16px; text-align: center">
        <Space>
          <Button type="primary" v-if="selectedTask.status === '运行中'"
            >暂停任务</Button
          >
          <Button
            type="primary"
            danger
            v-if="
              selectedTask.status === '运行中' ||
              selectedTask.status === '等待中'
            "
            >终止任务</Button
          >
          <Button type="primary" v-if="selectedTask.status === '失败'"
            >重试任务</Button
          >
          <Button>查看详细日志</Button>
        </Space>
      </div>
    </Card>
  </div>
</template>

<style scoped>
.task-monitor-container {
  padding: 0;
}
.stats-row {
  margin-bottom: 16px;
}
</style>
