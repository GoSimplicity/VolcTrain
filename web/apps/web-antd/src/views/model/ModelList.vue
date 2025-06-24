<script lang="ts" setup>
import { ref } from 'vue';
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
  Divider,
  Tooltip,
} from 'ant-design-vue';

defineOptions({ name: 'ModelList' });

// 模型列表数据
const dataSource = ref([
  {
    id: '1',
    name: 'BERT-base-chinese',
    type: 'NLP',
    framework: 'PyTorch',
    description: '中文BERT预训练模型',
    size: '375MB',
    creator: '张三',
    createTime: '2023-10-05',
    updateTime: '2023-11-10',
    version: '1.0.3',
    status: '已发布',
    shared: true,
  },
  {
    id: '2',
    name: 'ResNet50-ImageNet',
    type: '图像分类',
    framework: 'TensorFlow',
    description: '基于ImageNet的ResNet50预训练模型',
    size: '98MB',
    creator: '李四',
    createTime: '2023-09-15',
    updateTime: '2023-11-12',
    version: '2.1.0',
    status: '已发布',
    shared: true,
  },
  {
    id: '3',
    name: 'GPT-2-small-zh',
    type: 'NLP',
    framework: 'PyTorch',
    description: '中文GPT-2小模型',
    size: '510MB',
    creator: '王五',
    createTime: '2023-11-01',
    updateTime: '2023-11-14',
    version: '0.2.5',
    status: '训练中',
    shared: false,
  },
  {
    id: '4',
    name: 'YOLOv5-Custom',
    type: '目标检测',
    framework: 'PyTorch',
    description: '自定义YOLOv5目标检测模型',
    size: '87MB',
    creator: '赵六',
    createTime: '2023-10-20',
    updateTime: '2023-11-05',
    version: '1.5.0',
    status: '已发布',
    shared: false,
  },
  {
    id: '5',
    name: 'LSTM-Sentiment',
    type: 'NLP',
    framework: 'Keras',
    description: '情感分析LSTM模型',
    size: '45MB',
    creator: '张三',
    createTime: '2023-11-08',
    updateTime: '2023-11-15',
    version: '0.9.1',
    status: '测试中',
    shared: true,
  },
]);

// 表格列配置
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '模型名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
    customRender: ({ text }) => {
      const color = 
        text === 'NLP' ? 'blue' : 
        text === '图像分类' ? 'green' : 
        text === '目标检测' ? 'purple' : 'default';
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '框架',
    dataIndex: 'framework',
    key: 'framework',
    customRender: ({ text }) => {
      const color = 
        text === 'PyTorch' ? 'volcano' : 
        text === 'TensorFlow' ? 'geekblue' : 
        text === 'Keras' ? 'magenta' : 'default';
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
    customRender: ({ text }) => (
      <Tooltip placement="topLeft" title={text}>
        <div style="width: 150px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
          {text}
        </div>
      </Tooltip>
    ),
  },
  {
    title: '大小',
    dataIndex: 'size',
    key: 'size',
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    key: 'createTime',
  },
  {
    title: '更新时间',
    dataIndex: 'updateTime',
    key: 'updateTime',
  },
  {
    title: '版本',
    dataIndex: 'version',
    key: 'version',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      const color = 
        text === '已发布' ? 'green' : 
        text === '训练中' ? 'blue' : 
        text === '测试中' ? 'orange' : 'default';
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '共享',
    dataIndex: 'shared',
    key: 'shared',
    customRender: ({ text }) => {
      return text ? <Tag color="green">是</Tag> : <Tag>否</Tag>;
    },
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }) => (
      <Space size="middle">
        <a>查看</a>
        {record.status === '已发布' && <a>部署</a>}
        {record.status === '已发布' && <a>下载</a>}
        <a>删除</a>
      </Space>
    ),
  },
];

// 表单相关
const visible = ref(false);
const formRef = ref();
const formState = ref({
  name: '',
  type: '',
  framework: '',
  description: '',
  shared: false,
});

// 类型选项
const typeOptions = [
  { label: 'NLP', value: 'NLP' },
  { label: '图像分类', value: '图像分类' },
  { label: '目标检测', value: '目标检测' },
  { label: '语音识别', value: '语音识别' },
  { label: '推荐系统', value: '推荐系统' },
];

// 框架选项
const frameworkOptions = [
  { label: 'PyTorch', value: 'PyTorch' },
  { label: 'TensorFlow', value: 'TensorFlow' },
  { label: 'Keras', value: 'Keras' },
  { label: 'MXNet', value: 'MXNet' },
  { label: 'ONNX', value: 'ONNX' },
];

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
      // 模拟添加一条新数据
      dataSource.value.push({
        id: String(dataSource.value.length + 1),
        name: formState.value.name,
        type: formState.value.type,
        framework: formState.value.framework,
        description: formState.value.description,
        size: '0KB',
        creator: '当前用户',
        createTime: new Date().toLocaleDateString(),
        updateTime: new Date().toLocaleDateString(),
        version: '0.1.0',
        status: '创建中',
        shared: formState.value.shared,
      });
      
      visible.value = false;
      formState.value = {
        name: '',
        type: '',
        framework: '',
        description: '',
        shared: false,
      };
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};

// 搜索表单
const searchForm = ref({
  name: '',
  type: '',
  framework: '',
  status: '',
});

const handleSearch = () => {
  // 实际场景中，这里应该根据搜索条件过滤数据
  console.log('搜索条件:', searchForm.value);
};

const resetSearch = () => {
  searchForm.value = {
    name: '',
    type: '',
    framework: '',
    status: '',
  };
};
</script>

<template>
  <div class="model-list-container">
    <Card title="模型管理">
      <template #extra>
        <Space>
          <Button @click="showModal" type="primary">新建模型</Button>
          <Button>导入模型</Button>
        </Space>
      </template>

      <!-- 搜索区域 -->
      <div class="search-area">
        <Form layout="inline" :model="searchForm">
          <Form.Item name="name" label="模型名称">
            <Input v-model:value="searchForm.name" placeholder="请输入模型名称" />
          </Form.Item>
          <Form.Item name="type" label="模型类型">
            <Select v-model:value="searchForm.type" placeholder="请选择类型" style="width: 120px" allowClear>
              <Select.Option value="NLP">NLP</Select.Option>
              <Select.Option value="图像分类">图像分类</Select.Option>
              <Select.Option value="目标检测">目标检测</Select.Option>
              <Select.Option value="语音识别">语音识别</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="framework" label="框架">
            <Select v-model:value="searchForm.framework" placeholder="请选择框架" style="width: 120px" allowClear>
              <Select.Option value="PyTorch">PyTorch</Select.Option>
              <Select.Option value="TensorFlow">TensorFlow</Select.Option>
              <Select.Option value="Keras">Keras</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select v-model:value="searchForm.status" placeholder="请选择状态" style="width: 120px" allowClear>
              <Select.Option value="已发布">已发布</Select.Option>
              <Select.Option value="训练中">训练中</Select.Option>
              <Select.Option value="测试中">测试中</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" @click="handleSearch">搜索</Button>
            <Button style="margin-left: 8px" @click="resetSearch">重置</Button>
          </Form.Item>
        </Form>
      </div>

      <Divider style="margin: 12px 0" />

      <Table :columns="columns" :dataSource="dataSource" rowKey="id" />
    </Card>

    <Modal
      v-model:visible="visible"
      title="新建模型"
      @ok="handleOk"
      @cancel="handleCancel"
      :maskClosable="false"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="name"
          label="模型名称"
          :rules="[{ required: true, message: '请输入模型名称' }]"
        >
          <Input v-model:value="formState.name" placeholder="请输入模型名称" />
        </Form.Item>
        <Form.Item
          name="type"
          label="模型类型"
          :rules="[{ required: true, message: '请选择模型类型' }]"
        >
          <Select v-model:value="formState.type" placeholder="请选择模型类型">
            <Select.Option
              v-for="option in typeOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="framework"
          label="框架"
          :rules="[{ required: true, message: '请选择框架' }]"
        >
          <Select v-model:value="formState.framework" placeholder="请选择框架">
            <Select.Option
              v-for="option in frameworkOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item name="description" label="描述">
          <Input.TextArea
            v-model:value="formState.description"
            placeholder="请输入模型描述"
          />
        </Form.Item>
        <Form.Item name="shared" label="是否共享">
          <Select v-model:value="formState.shared">
            <Select.Option :value="true">是</Select.Option>
            <Select.Option :value="false">否</Select.Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.model-list-container {
  padding: 0;
}
.search-area {
  margin-bottom: 16px;
}
</style>
