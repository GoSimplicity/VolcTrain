<script lang="ts" setup>
import { ref } from 'vue';
import {
  Table,
  Card,
  Tag,
  Progress,
  Space,
  Tooltip,
  Button,
  Statistic,
  Row,
  Col,
} from 'ant-design-vue';

defineOptions({ name: 'GPUList' });

// GPU统计信息
const gpuStats = ref({
  totalGPUs: 48,
  activeGPUs: 42,
  avgUsage: 68,
  availableGPUs: 15,
});

// GPU列表数据
const dataSource = ref([
  {
    id: '1',
    name: 'GPU-001',
    type: 'NVIDIA A100',
    node: 'node-01',
    memory: '80GB',
    usage: 95,
    memoryUsage: 85,
    temperature: 72,
    powerUsage: '250W / 300W',
    status: '使用中',
    user: '张三',
    task: 'BERT微调任务1',
    startTime: '2023-11-15 08:30',
    runTime: '3小时25分钟',
  },
  {
    id: '2',
    name: 'GPU-002',
    type: 'NVIDIA A100',
    node: 'node-01',
    memory: '80GB',
    usage: 88,
    memoryUsage: 92,
    temperature: 75,
    powerUsage: '260W / 300W',
    status: '使用中',
    user: '李四',
    task: 'ResNet图像分类',
    startTime: '2023-11-15 06:15',
    runTime: '5小时40分钟',
  },
  {
    id: '3',
    name: 'GPU-003',
    type: 'NVIDIA A100',
    node: 'node-01',
    memory: '80GB',
    usage: 0,
    memoryUsage: 0,
    temperature: 45,
    powerUsage: '30W / 300W',
    status: '空闲',
    user: '-',
    task: '-',
    startTime: '-',
    runTime: '-',
  },
  {
    id: '4',
    name: 'GPU-004',
    type: 'NVIDIA A100',
    node: 'node-01',
    memory: '80GB',
    usage: 0,
    memoryUsage: 0,
    temperature: 42,
    powerUsage: '25W / 300W',
    status: '空闲',
    user: '-',
    task: '-',
    startTime: '-',
    runTime: '-',
  },
  {
    id: '5',
    name: 'GPU-005',
    type: 'NVIDIA V100',
    node: 'node-02',
    memory: '32GB',
    usage: 76,
    memoryUsage: 68,
    temperature: 68,
    powerUsage: '180W / 250W',
    status: '使用中',
    user: '王五',
    task: '数据预处理任务',
    startTime: '2023-11-15 10:45',
    runTime: '1小时10分钟',
  },
  {
    id: '6',
    name: 'GPU-006',
    type: 'NVIDIA V100',
    node: 'node-02',
    memory: '32GB',
    usage: 0,
    memoryUsage: 0,
    temperature: 0,
    powerUsage: '0W / 250W',
    status: '故障',
    user: '-',
    task: '-',
    startTime: '-',
    runTime: '-',
  },
]);

// 表格列配置
const columns = [
  {
    title: 'GPU ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: 'GPU名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: 'GPU型号',
    dataIndex: 'type',
    key: 'type',
  },
  {
    title: '所在节点',
    dataIndex: 'node',
    key: 'node',
  },
  {
    title: '显存容量',
    dataIndex: 'memory',
    key: 'memory',
  },
  {
    title: 'GPU使用率',
    dataIndex: 'usage',
    key: 'usage',
    customRender: ({ text }) => {
      if (text === 0 && record.status === '故障') return <Tag color="red">故障</Tag>;
      const color = text > 90 ? 'red' : text > 70 ? 'orange' : 'green';
      return <Progress percent={text} size="small" strokeColor={color} />;
    },
  },
  {
    title: '显存使用率',
    dataIndex: 'memoryUsage',
    key: 'memoryUsage',
    customRender: ({ text, record }) => {
      if (text === 0 && record.status === '故障') return <Tag color="red">故障</Tag>;
      const color = text > 90 ? 'red' : text > 70 ? 'orange' : 'green';
      return <Progress percent={text} size="small" strokeColor={color} />;
    },
  },
  {
    title: '温度',
    dataIndex: 'temperature',
    key: 'temperature',
    customRender: ({ text, record }) => {
      if (text === 0 && record.status === '故障') return <Tag color="red">故障</Tag>;
      let color = 'green';
      if (text > 80) {
        color = 'red';
      } else if (text > 70) {
        color = 'orange';
      }
      return <Tooltip title={`${text}°C`}><Tag color={color}>{text}°C</Tag></Tooltip>;
    },
  },
  {
    title: '功耗',
    dataIndex: 'powerUsage',
    key: 'powerUsage',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      let color = 'green';
      if (text === '使用中') {
        color = 'blue';
      } else if (text === '故障') {
        color = 'red';
      }
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '使用用户',
    dataIndex: 'user',
    key: 'user',
  },
  {
    title: '任务名称',
    dataIndex: 'task',
    key: 'task',
  },
  {
    title: '开始时间',
    dataIndex: 'startTime',
    key: 'startTime',
  },
  {
    title: '运行时长',
    dataIndex: 'runTime',
    key: 'runTime',
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }) => (
      <Space size="middle">
        <a>详情</a>
        {record.status === '使用中' && <a>释放</a>}
        {record.status === '故障' && <a>维修</a>}
      </Space>
    ),
  },
]);
</script>

<template>
  <div class="gpu-list-container">
    <Row :gutter="16" class="stats-row">
      <Col :span="6">
        <Card>
          <Statistic title="GPU总数" :value="gpuStats.totalGPUs" />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic title="活跃GPU" :value="gpuStats.activeGPUs" />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic title="平均使用率" :value="gpuStats.avgUsage" suffix="%" />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic title="可用GPU数" :value="gpuStats.availableGPUs" />
        </Card>
      </Col>
    </Row>

    <Card title="GPU资源列表">
      <template #extra>
        <Space>
          <Button>刷新</Button>
          <Button type="primary">申请GPU资源</Button>
        </Space>
      </template>
      <Table 
        :columns="columns" 
        :dataSource="dataSource" 
        rowKey="id"
        :pagination="{ pageSize: 10 }"
      />
    </Card>
  </div>
</template>

<style scoped>
.gpu-list-container {
  padding: 0;
}
.stats-row {
  margin-bottom: 16px;
}
</style>
