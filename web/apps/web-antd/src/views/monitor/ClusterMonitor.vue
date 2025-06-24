<script lang="ts" setup>
import { ref, h } from 'vue';
import {
  Card,
  Row,
  Col,
  Statistic,
  Table,
  Progress,
  Space,
  Divider,
  Tabs,
  Tag,
} from 'ant-design-vue';

defineOptions({ name: 'ClusterMonitor' });

// 集群统计信息
const clusterStats = ref({
  totalNodes: 12,
  activeNodes: 10,
  totalGPUs: 48,
  activeGPUs: 42,
  cpuUsage: 68,
  memoryUsage: 72,
  networkIn: '256.7 MB/s',
  networkOut: '78.3 MB/s',
});

// 节点列表数据
const nodesData = ref([
  {
    id: '1',
    name: 'node-01',
    status: '在线',
    ip: '192.168.1.101',
    cpu: 75,
    memory: 68,
    gpu: 90,
    gpuCount: 4,
    gpuType: 'NVIDIA A100',
    uptime: '30天12小时',
  },
  {
    id: '2',
    name: 'node-02',
    status: '在线',
    ip: '192.168.1.102',
    cpu: 45,
    memory: 60,
    gpu: 82,
    gpuCount: 4,
    gpuType: 'NVIDIA A100',
    uptime: '15天8小时',
  },
  {
    id: '3',
    name: 'node-03',
    status: '在线',
    ip: '192.168.1.103',
    cpu: 30,
    memory: 42,
    gpu: 50,
    gpuCount: 4,
    gpuType: 'NVIDIA V100',
    uptime: '45天3小时',
  },
  {
    id: '4',
    name: 'node-04',
    status: '离线',
    ip: '192.168.1.104',
    cpu: 0,
    memory: 0,
    gpu: 0,
    gpuCount: 4,
    gpuType: 'NVIDIA V100',
    uptime: '-',
  },
]);

// 节点表格列配置
const columns = [
  {
    title: '节点ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '节点名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      return h(
        Tag,
        { color: text === '在线' ? 'success' : 'error' },
        () => text
      );
    },
  },
  {
    title: 'IP地址',
    dataIndex: 'ip',
    key: 'ip',
  },
  {
    title: 'CPU使用率',
    dataIndex: 'cpu',
    key: 'cpu',
    customRender: ({ text }) => {
      const color = text > 90 ? 'red' : text > 70 ? 'orange' : 'green';
      return h(Progress, { percent: text, size: "small", strokeColor: color });
    },
  },
  {
    title: '内存使用率',
    dataIndex: 'memory',
    key: 'memory',
    customRender: ({ text }) => {
      const color = text > 90 ? 'red' : text > 70 ? 'orange' : 'green';
      return h(Progress, { percent: text, size: "small", strokeColor: color });
    },
  },
  {
    title: 'GPU使用率',
    dataIndex: 'gpu',
    key: 'gpu',
    customRender: ({ text }) => {
      const color = text > 90 ? 'red' : text > 70 ? 'orange' : 'green';
      return h(Progress, { percent: text, size: "small", strokeColor: color });
    },
  },
  {
    title: 'GPU数量',
    dataIndex: 'gpuCount',
    key: 'gpuCount',
  },
  {
    title: 'GPU类型',
    dataIndex: 'gpuType',
    key: 'gpuType',
  },
  {
    title: '运行时间',
    dataIndex: 'uptime',
    key: 'uptime',
  },
  {
    title: '操作',
    key: 'action',
    customRender: () => {
      return h(Space, { size: "middle" }, [
        h('a', {}, '详情'),
        h('a', {}, '重启')
      ]);
    },
  },
];

const activeKey = ref('1');
</script>

<template>
  <div class="cluster-monitor-container">
    <Row :gutter="16" class="stats-row">
      <Col :span="6">
        <Card>
          <Statistic 
            title="节点总数" 
            :value="clusterStats.totalNodes" 
            :suffix="'/ ' + clusterStats.activeNodes + ' 活跃'"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic 
            title="GPU总数" 
            :value="clusterStats.totalGPUs" 
            :suffix="'/ ' + clusterStats.activeGPUs + ' 活跃'"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic title="CPU使用率" :value="clusterStats.cpuUsage" suffix="%" />
          <Progress :percent="clusterStats.cpuUsage" :showInfo="false" />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic title="内存使用率" :value="clusterStats.memoryUsage" suffix="%" />
          <Progress :percent="clusterStats.memoryUsage" :showInfo="false" />
        </Card>
      </Col>
    </Row>

    <Row :gutter="16" class="stats-row">
      <Col :span="12">
        <Card>
          <Statistic title="网络入站流量" :value="clusterStats.networkIn" />
        </Card>
      </Col>
      <Col :span="12">
        <Card>
          <Statistic title="网络出站流量" :value="clusterStats.networkOut" />
        </Card>
      </Col>
    </Row>

    <Divider />

    <Card title="集群节点监控">
      <Tabs v-model:activeKey="activeKey">
        <Tabs.TabPane key="1" tab="所有节点">
          <Table :columns="columns" :dataSource="nodesData" rowKey="id" />
        </Tabs.TabPane>
        <Tabs.TabPane key="2" tab="在线节点">
          <Table :columns="columns" :dataSource="nodesData.filter(node => node.status === '在线')" rowKey="id" />
        </Tabs.TabPane>
        <Tabs.TabPane key="3" tab="离线节点">
          <Table :columns="columns" :dataSource="nodesData.filter(node => node.status === '离线')" rowKey="id" />
        </Tabs.TabPane>
      </Tabs>
    </Card>
  </div>
</template>

<style scoped>
.cluster-monitor-container {
  padding: 0;
}
.stats-row {
  margin-bottom: 16px;
}
.text-success {
  color: #52c41a;
}
.text-danger {
  color: #f5222d;
}
</style>
