<script lang="ts" setup>
import { ref, h } from 'vue';
import {
  Table,
  Card,
  Button,
  Space,
  Tag,
  Modal,
  Form,
  Input,
  Select,
} from 'ant-design-vue';

defineOptions({ name: 'TrainingQueue' });

const dataSource = ref([
  {
    id: '1',
    name: 'BERT微调任务',
    dataset: 'NLP语料集',
    model: 'BERT-base',
    gpu: 'NVIDIA A100',
    status: '等待中',
    priority: '高',
    createTime: '2023-11-10 10:30:00',
    user: '张三',
    estimatedTime: '2小时',
  },
  {
    id: '2',
    name: '图像分类训练',
    dataset: '图像数据集',
    model: 'ResNet50',
    gpu: 'NVIDIA V100',
    status: '运行中',
    priority: '中',
    createTime: '2023-11-10 09:15:00',
    user: '李四',
    estimatedTime: '5小时',
  },
]);

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
    title: '数据集',
    dataIndex: 'dataset',
    key: 'dataset',
  },
  {
    title: '模型',
    dataIndex: 'model',
    key: 'model',
  },
  {
    title: 'GPU资源',
    dataIndex: 'gpu',
    key: 'gpu',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      const color =
        text === '运行中' ? 'green' : text === '等待中' ? 'orange' : 'default';
      return h(Tag, { color }, () => text);
    },
  },
  {
    title: '优先级',
    dataIndex: 'priority',
    key: 'priority',
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    key: 'createTime',
  },
  {
    title: '创建用户',
    dataIndex: 'user',
    key: 'user',
  },
  {
    title: '预计运行时长',
    dataIndex: 'estimatedTime',
    key: 'estimatedTime',
  },
  {
    title: '操作',
    key: 'action',
    customRender: () => {
      return h(Space, { size: 'middle' }, [
        h('a', {}, '查看详情'),
        h('a', {}, '暂停'),
        h('a', {}, '删除')
      ]);
    },
  },
];

const visible = ref(false);
const formRef = ref();
const formState = ref({
  name: '',
  dataset: '',
  model: '',
  gpuResource: '',
  priority: '中',
});

const showModal = () => {
  visible.value = true;
};

const handleCancel = () => {
  visible.value = false;
};

const handleOk = () => {
  formRef.value
    .validate()
    .then(() => {
      dataSource.value.push({
        id: String(dataSource.value.length + 1),
        name: formState.value.name,
        dataset: formState.value.dataset,
        model: formState.value.model,
        gpu: formState.value.gpuResource,
        status: '等待中',
        priority: formState.value.priority,
        createTime: new Date().toLocaleString(),
        user: '当前用户',
        estimatedTime: '预估中',
      });
      visible.value = false;
      formState.value = {
        name: '',
        dataset: '',
        model: '',
        gpuResource: '',
        priority: '中',
      };
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};
</script>

<template>
  <div class="training-queue-container">
    <Card title="训练队列管理" :bordered="false">
      <template #extra>
        <Button type="primary" @click="showModal">新增训练任务</Button>
      </template>
      <Table :columns="columns" :dataSource="dataSource" rowKey="id" />
    </Card>

    <Modal
      v-model:visible="visible"
      title="新增训练任务"
      @ok="handleOk"
      @cancel="handleCancel"
      :maskClosable="false"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="name"
          label="任务名称"
          :rules="[{ required: true, message: '请输入任务名称' }]"
        >
          <Input v-model:value="formState.name" placeholder="请输入任务名称" />
        </Form.Item>
        <Form.Item
          name="dataset"
          label="数据集"
          :rules="[{ required: true, message: '请选择数据集' }]"
        >
          <Select v-model:value="formState.dataset" placeholder="请选择数据集">
            <Select.Option value="NLP语料集">NLP语料集</Select.Option>
            <Select.Option value="图像数据集">图像数据集</Select.Option>
            <Select.Option value="音频数据集">音频数据集</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="model"
          label="模型"
          :rules="[{ required: true, message: '请选择模型' }]"
        >
          <Select v-model:value="formState.model" placeholder="请选择模型">
            <Select.Option value="BERT-base">BERT-base</Select.Option>
            <Select.Option value="ResNet50">ResNet50</Select.Option>
            <Select.Option value="GPT-2">GPT-2</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="gpuResource"
          label="GPU资源"
          :rules="[{ required: true, message: '请选择GPU资源' }]"
        >
          <Select
            v-model:value="formState.gpuResource"
            placeholder="请选择GPU资源"
          >
            <Select.Option value="NVIDIA A100">NVIDIA A100</Select.Option>
            <Select.Option value="NVIDIA V100">NVIDIA V100</Select.Option>
            <Select.Option value="NVIDIA T4">NVIDIA T4</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item name="priority" label="优先级">
          <Select v-model:value="formState.priority">
            <Select.Option value="高">高</Select.Option>
            <Select.Option value="中">中</Select.Option>
            <Select.Option value="低">低</Select.Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.training-queue-container {
  padding: 0;
}
</style>
