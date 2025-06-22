<script lang="ts" setup>
import { ref, h } from 'vue';
import { Table, Card, Button, Space, Modal, Form, Input, InputNumber, Switch, Tabs, Select } from 'ant-design-vue';

defineOptions({ name: 'TrainingConfig' });

// Volcano相关配置
const volcanoConfig = ref([
  {
    id: '1',
    name: 'default',
    description: '默认调度策略',
    schedulerName: 'volcano',
    queueName: 'default',
    priority: 'Medium',
    preemptible: true,
    minAvailable: 1,
    status: '启用',
  },
  {
    id: '2',
    name: 'high-priority',
    description: '高优先级调度策略',
    schedulerName: 'volcano',
    queueName: 'high-priority',
    priority: 'High',
    preemptible: false,
    minAvailable: 2,
    status: '启用',
  },
  {
    id: '3',
    name: 'batch-training',
    description: '批量训练调度策略',
    schedulerName: 'volcano',
    queueName: 'batch',
    priority: 'Low',
    preemptible: true,
    minAvailable: 1,
    status: '禁用',
  },
]);

// 资源配额配置
const quotaConfig = ref([
  {
    id: '1',
    name: '标准用户',
    description: '标准用户资源配额',
    maxGPUs: 4,
    maxCPUs: 16,
    maxMemory: '64Gi',
    maxStorage: '100Gi',
    status: '启用',
  },
  {
    id: '2',
    name: 'VIP用户',
    description: 'VIP用户资源配额',
    maxGPUs: 8,
    maxCPUs: 32,
    maxMemory: '128Gi',
    maxStorage: '200Gi',
    status: '启用',
  },
  {
    id: '3',
    name: '管理员',
    description: '管理员资源配额',
    maxGPUs: 16,
    maxCPUs: 64,
    maxMemory: '256Gi',
    maxStorage: '500Gi',
    status: '启用',
  },
]);

// 集群配置
const clusterConfig = ref([
  {
    id: '1',
    name: '主集群',
    description: '主要训练集群',
    nodeCount: 10,
    totalGPUs: 40,
    totalCPUs: 320,
    totalMemory: '1024Gi',
    kubeConfig: '/etc/kubernetes/admin.conf',
    status: '运行中',
  },
  {
    id: '2',
    name: '备用集群',
    description: '备用训练集群',
    nodeCount: 5,
    totalGPUs: 20,
    totalCPUs: 160,
    totalMemory: '512Gi',
    kubeConfig: '/etc/kubernetes/backup.conf',
    status: '运行中',
  },
]);

// Volcano列
const volcanoColumns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '策略名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '调度器',
    dataIndex: 'schedulerName',
    key: 'schedulerName',
  },
  {
    title: '队列名称',
    dataIndex: 'queueName',
    key: 'queueName',
  },
  {
    title: '优先级',
    dataIndex: 'priority',
    key: 'priority',
  },
  {
    title: '可抢占',
    dataIndex: 'preemptible',
    key: 'preemptible',
    customRender: ({ text }) => (text ? '是' : '否'),
  },
  {
    title: '最小可用资源',
    dataIndex: 'minAvailable',
    key: 'minAvailable',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
  },
  {
    title: '操作',
    key: 'action',
    customRender: () => {
      return h(Space, { size: 'middle' }, [
        h('a', {}, '编辑'),
        h('a', {}, '删除')
      ]);
    },
  },
];

// 配额列
const quotaColumns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '配额名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '最大GPU数量',
    dataIndex: 'maxGPUs',
    key: 'maxGPUs',
  },
  {
    title: '最大CPU数量',
    dataIndex: 'maxCPUs',
    key: 'maxCPUs',
  },
  {
    title: '最大内存',
    dataIndex: 'maxMemory',
    key: 'maxMemory',
  },
  {
    title: '最大存储',
    dataIndex: 'maxStorage',
    key: 'maxStorage',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
  },
  {
    title: '操作',
    key: 'action',
    customRender: () => {
      return h(Space, { size: 'middle' }, [
        h('a', {}, '编辑'),
        h('a', {}, '删除')
      ]);
    },
  },
];

// 集群列
const clusterColumns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '集群名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '节点数量',
    dataIndex: 'nodeCount',
    key: 'nodeCount',
  },
  {
    title: '总GPU数量',
    dataIndex: 'totalGPUs',
    key: 'totalGPUs',
  },
  {
    title: '总CPU数量',
    dataIndex: 'totalCPUs',
    key: 'totalCPUs',
  },
  {
    title: '总内存',
    dataIndex: 'totalMemory',
    key: 'totalMemory',
  },
  {
    title: 'Kube配置',
    dataIndex: 'kubeConfig',
    key: 'kubeConfig',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
  },
  {
    title: '操作',
    key: 'action',
    customRender: () => {
      return h(Space, { size: 'middle' }, [
        h('a', {}, '编辑'),
        h('a', {}, '详情')
      ]);
    },
  },
];

const activeKey = ref('volcano');

// 模态框相关状态
const visible = ref(false);
const formRef = ref();
const currentTab = ref('volcano');
const volcanoFormState = ref({
  name: '',
  description: '',
  schedulerName: 'volcano',
  queueName: '',
  priority: 'Medium',
  preemptible: true,
  minAvailable: 1,
});

const showModal = (tab) => {
  currentTab.value = tab;
  visible.value = true;
};

const handleCancel = () => {
  visible.value = false;
};

const handleOk = () => {
  formRef.value
    .validate()
    .then(() => {
      if (currentTab.value === 'volcano') {
        volcanoConfig.value.push({
          id: String(volcanoConfig.value.length + 1),
          name: volcanoFormState.value.name,
          description: volcanoFormState.value.description,
          schedulerName: volcanoFormState.value.schedulerName,
          queueName: volcanoFormState.value.queueName,
          priority: volcanoFormState.value.priority,
          preemptible: volcanoFormState.value.preemptible,
          minAvailable: volcanoFormState.value.minAvailable,
          status: '启用',
        });
        volcanoFormState.value = {
          name: '',
          description: '',
          schedulerName: 'volcano',
          queueName: '',
          priority: 'Medium',
          preemptible: true,
          minAvailable: 1,
        };
      }
      visible.value = false;
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};
</script>

<template>
  <div class="training-config-container">
    <Card :bordered="false">
      <Tabs v-model:activeKey="activeKey">
        <Tabs.TabPane key="volcano" tab="Volcano调度配置">
          <div class="table-operations">
            <Button type="primary" @click="showModal('volcano')">新增调度配置</Button>
          </div>
          <Table
            :columns="volcanoColumns"
            :dataSource="volcanoConfig"
            rowKey="id"
            :pagination="{ pageSize: 10 }"
          />
        </Tabs.TabPane>
        
        <Tabs.TabPane key="quota" tab="资源配额配置">
          <div class="table-operations">
            <Button type="primary" @click="showModal('quota')">新增配额配置</Button>
          </div>
          <Table
            :columns="quotaColumns"
            :dataSource="quotaConfig"
            rowKey="id"
            :pagination="{ pageSize: 10 }"
          />
        </Tabs.TabPane>
        
        <Tabs.TabPane key="cluster" tab="集群配置">
          <div class="table-operations">
            <Button type="primary" @click="showModal('cluster')">新增集群配置</Button>
          </div>
          <Table
            :columns="clusterColumns"
            :dataSource="clusterConfig"
            rowKey="id"
            :pagination="{ pageSize: 10 }"
          />
        </Tabs.TabPane>
      </Tabs>
    </Card>

    <!-- Volcano配置模态框 -->
    <Modal
      v-if="currentTab === 'volcano'"
      v-model:visible="visible"
      title="新增Volcano调度配置"
      @ok="handleOk"
      @cancel="handleCancel"
      :maskClosable="false"
      width="700px"
    >
      <Form ref="formRef" :model="volcanoFormState" layout="vertical">
        <Form.Item
          name="name"
          label="策略名称"
          :rules="[{ required: true, message: '请输入策略名称' }]"
        >
          <Input v-model:value="volcanoFormState.name" placeholder="请输入策略名称" />
        </Form.Item>
        <Form.Item name="description" label="描述">
          <Input.TextArea v-model:value="volcanoFormState.description" placeholder="请输入描述" />
        </Form.Item>
        <Form.Item name="schedulerName" label="调度器名称">
          <Input v-model:value="volcanoFormState.schedulerName" disabled />
        </Form.Item>
        <Form.Item
          name="queueName"
          label="队列名称"
          :rules="[{ required: true, message: '请输入队列名称' }]"
        >
          <Input v-model:value="volcanoFormState.queueName" placeholder="请输入队列名称" />
        </Form.Item>
        <Form.Item name="priority" label="优先级">
          <Select v-model:value="volcanoFormState.priority">
            <Select.Option value="High">High</Select.Option>
            <Select.Option value="Medium">Medium</Select.Option>
            <Select.Option value="Low">Low</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item name="preemptible" label="可抢占">
          <Switch v-model:checked="volcanoFormState.preemptible" />
        </Form.Item>
        <Form.Item name="minAvailable" label="最小可用资源">
          <InputNumber v-model:value="volcanoFormState.minAvailable" :min="1" />
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.training-config-container {
  padding: 0;
}
.table-operations {
  margin-bottom: 16px;
}
</style>
