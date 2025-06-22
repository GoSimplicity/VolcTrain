<script lang="ts" setup>
import { ref } from 'vue';
import {
  Card,
  Table,
  Tag,
  Space,
  Button,
  Modal,
  Form,
  Input,
  Select,
  Divider,
  Tabs,
  Radio,
  Tooltip,
  message,
} from 'ant-design-vue';
import { SearchOutlined, CloudDownloadOutlined, StarOutlined, StarFilled } from '@ant-design/icons-vue';

defineOptions({ name: 'ModelRegistry' });

// 模型仓库数据
const registryModels = ref([
  {
    id: '1',
    name: 'BERT-base-chinese',
    type: 'NLP',
    framework: 'PyTorch',
    description: '中文BERT预训练模型，适用于多种中文NLP任务',
    version: '1.2.0',
    size: '375MB',
    publisher: 'HuggingFace',
    publishDate: '2023-05-15',
    downloads: 12453,
    stars: 245,
    license: 'Apache-2.0',
    isStarred: true,
  },
  {
    id: '2',
    name: 'ResNet50-ImageNet',
    type: '图像分类',
    framework: 'TensorFlow',
    description: '基于ImageNet训练的ResNet50模型，提供高精度图像分类能力',
    version: '2.3.1',
    size: '98MB',
    publisher: 'TensorFlow',
    publishDate: '2023-04-20',
    downloads: 8765,
    stars: 189,
    license: 'MIT',
    isStarred: false,
  },
  {
    id: '3',
    name: 'YOLOv5',
    type: '目标检测',
    framework: 'PyTorch',
    description: 'YOLOv5目标检测模型，速度快、精度高',
    version: '6.1.0',
    size: '87MB',
    publisher: 'Ultralytics',
    publishDate: '2023-06-10',
    downloads: 15678,
    stars: 312,
    license: 'GPL-3.0',
    isStarred: true,
  },
  {
    id: '4',
    name: 'GPT-2-small',
    type: 'NLP',
    framework: 'PyTorch',
    description: 'OpenAI的GPT-2小型模型，用于文本生成',
    version: '1.0.1',
    size: '510MB',
    publisher: 'HuggingFace',
    publishDate: '2023-03-05',
    downloads: 9876,
    stars: 210,
    license: 'MIT',
    isStarred: false,
  },
  {
    id: '5',
    name: 'EfficientNet-B0',
    type: '图像分类',
    framework: 'TensorFlow',
    description: 'EfficientNet-B0图像分类模型，在保持高精度的同时体积更小',
    version: '2.0.0',
    size: '29MB',
    publisher: 'Google',
    publishDate: '2023-02-18',
    downloads: 7654,
    stars: 156,
    license: 'Apache-2.0',
    isStarred: false,
  },
  {
    id: '6',
    name: 'Wav2Vec2',
    type: '语音识别',
    framework: 'PyTorch',
    description: '用于语音识别的自监督学习模型',
    version: '2.0.1',
    size: '1.2GB',
    publisher: 'Facebook AI',
    publishDate: '2023-07-01',
    downloads: 5432,
    stars: 178,
    license: 'MIT',
    isStarred: false,
  },
  {
    id: '7',
    name: 'MobileNetV3',
    type: '图像分类',
    framework: 'TensorFlow',
    description: '移动设备优化的图像分类模型',
    version: '1.0.0',
    size: '23MB',
    publisher: 'Google',
    publishDate: '2023-01-10',
    downloads: 8901,
    stars: 145,
    license: 'Apache-2.0',
    isStarred: false,
  },
]);

// 我的模型仓库数据
const myRegistryModels = ref([
  {
    id: '1',
    name: 'LSTM-Sentiment',
    type: 'NLP',
    framework: 'Keras',
    description: '情感分析LSTM模型',
    version: '0.9.1',
    size: '45MB',
    publisher: '当前用户',
    publishDate: '2023-11-15',
    downloads: 12,
    stars: 3,
    license: 'MIT',
    isPublic: true,
    status: '已发布',
  },
  {
    id: '2',
    name: 'YOLOv5-Custom',
    type: '目标检测',
    framework: 'PyTorch',
    description: '自定义YOLOv5目标检测模型',
    version: '1.5.0',
    size: '87MB',
    publisher: '当前用户',
    publishDate: '2023-11-05',
    downloads: 8,
    stars: 1,
    license: 'Apache-2.0',
    isPublic: true,
    status: '已发布',
  },
  {
    id: '3',
    name: 'ResNet50-Custom',
    type: '图像分类',
    framework: 'PyTorch',
    description: '自定义ResNet50模型',
    version: '0.5.0',
    size: '95MB',
    publisher: '当前用户',
    publishDate: '-',
    downloads: 0,
    stars: 0,
    license: 'MIT',
    isPublic: false,
    status: '未发布',
  },
]);

// 表格列配置 - 公共模型仓库
const publicColumns = [
  {
    title: '模型名称',
    dataIndex: 'name',
    key: 'name',
    customRender: ({ text, record }) => (
      <div>
        <div>{text}</div>
        <div style="font-size: 12px; color: #999;">
          {record.publisher} | {record.license}
        </div>
      </div>
    ),
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
    customRender: ({ text }) => {
      const color = 
        text === 'NLP' ? 'blue' : 
        text === '图像分类' ? 'green' : 
        text === '目标检测' ? 'purple' : 
        text === '语音识别' ? 'orange' :
        'default';
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
        text === 'Keras' ? 'magenta' : 
        'default';
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
    customRender: ({ text }) => (
      <Tooltip placement="topLeft" title={text}>
        <div style="width: 200px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
          {text}
        </div>
      </Tooltip>
    ),
  },
  {
    title: '版本',
    dataIndex: 'version',
    key: 'version',
  },
  {
    title: '大小',
    dataIndex: 'size',
    key: 'size',
  },
  {
    title: '发布日期',
    dataIndex: 'publishDate',
    key: 'publishDate',
  },
  {
    title: '下载量',
    dataIndex: 'downloads',
    key: 'downloads',
    sorter: (a, b) => a.downloads - b.downloads,
  },
  {
    title: '收藏',
    dataIndex: 'stars',
    key: 'stars',
    sorter: (a, b) => a.stars - b.stars,
    customRender: ({ text, record }) => (
      <Space>
        <span>{text}</span>
        {record.isStarred ? 
          <StarFilled style="color: #faad14; cursor: pointer;" onClick={() => toggleStar(record)} /> : 
          <StarOutlined style="cursor: pointer;" onClick={() => toggleStar(record)} />}
      </Space>
    ),
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }) => (
      <Space size="middle">
        <a>详情</a>
        <a onClick={() => showImportModal(record)}><CloudDownloadOutlined /> 导入</a>
      </Space>
    ),
  },
];

// 表格列配置 - 我的模型仓库
const myColumns = [
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
        text === '目标检测' ? 'purple' : 
        text === '语音识别' ? 'orange' :
        'default';
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '框架',
    dataIndex: 'framework',
    key: 'framework',
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
    customRender: ({ text }) => (
      <Tooltip placement="topLeft" title={text}>
        <div style="width: 200px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
          {text}
        </div>
      </Tooltip>
    ),
  },
  {
    title: '版本',
    dataIndex: 'version',
    key: 'version',
  },
  {
    title: '大小',
    dataIndex: 'size',
    key: 'size',
  },
  {
    title: '可见性',
    dataIndex: 'isPublic',
    key: 'isPublic',
    customRender: ({ text }) => {
      return text ? <Tag color="green">公开</Tag> : <Tag>私有</Tag>;
    },
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      const color = 
        text === '已发布' ? 'green' : 
        text === '未发布' ? 'orange' : 
        'default';
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '下载量',
    dataIndex: 'downloads',
    key: 'downloads',
  },
  {
    title: '收藏',
    dataIndex: 'stars',
    key: 'stars',
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }) => (
      <Space size="middle">
        <a>详情</a>
        {record.status === '未发布' ? 
          <a onClick={() => handlePublish(record)}>发布</a> : 
          <a onClick={() => handleUnpublish(record)}>取消发布</a>}
        <a>删除</a>
      </Space>
    ),
  },
];

// 导入模型弹窗
const importVisible = ref(false);
const importFormRef = ref();
const importModel = ref(null);
const importFormState = ref({
  name: '',
  shared: false,
});

// 发布模型弹窗
const publishVisible = ref(false);
const publishFormRef = ref();
const publishingModel = ref(null);
const publishFormState = ref({
  isPublic: true,
  license: '',
  description: '',
});

// 搜索表单
const searchForm = ref({
  type: '',
  framework: '',
  sort: 'downloads',
});

// 标签页
const activeKey = ref('public');

// 显示导入模型弹窗
const showImportModal = (model) => {
  importModel.value = model;
  importFormState.value = {
    name: `${model.name}-imported`,
    shared: false,
  };
  importVisible.value = true;
};

// 显示发布模型弹窗
const handlePublish = (model) => {
  publishingModel.value = model;
  publishFormState.value = {
    isPublic: true,
    license: 'MIT',
    description: model.description,
  };
  publishVisible.value = true;
};

// 取消发布
const handleUnpublish = (model) => {
  Modal.confirm({
    title: '确认取消发布',
    content: `确定要取消发布模型 ${model.name} 吗？取消发布后，该模型将不再对其他用户可见。`,
    onOk() {
      // 更新模型状态
      const index = myRegistryModels.value.findIndex(m => m.id === model.id);
      if (index !== -1) {
        myRegistryModels.value[index].status = '未发布';
        myRegistryModels.value[index].isPublic = false;
        message.success('已取消发布');
      }
    },
  });
};

// 收藏/取消收藏
const toggleStar = (model) => {
  model.isStarred = !model.isStarred;
  if (model.isStarred) {
    model.stars++;
    message.success(`已收藏 ${model.name}`);
  } else {
    model.stars--;
    message.info(`已取消收藏 ${model.name}`);
  }
};

// 处理搜索
const handleSearch = () => {
  console.log('搜索条件:', searchForm.value);
  // 实际应用中这里会根据条件过滤数据
};

// 处理导入
const handleImportOk = () => {
  importFormRef.value
    .validate()
    .then(() => {
      // 模拟导入成功
      message.success(`成功导入模型 ${importModel.value.name}`);
      importVisible.value = false;
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};

// 处理发布
const handlePublishOk = () => {
  publishFormRef.value
    .validate()
    .then(() => {
      // 更新模型状态
      const index = myRegistryModels.value.findIndex(m => m.id === publishingModel.value.id);
      if (index !== -1) {
        myRegistryModels.value[index].status = '已发布';
        myRegistryModels.value[index].isPublic = publishFormState.value.isPublic;
        myRegistryModels.value[index].license = publishFormState.value.license;
        myRegistryModels.value[index].publishDate = new Date().toLocaleDateString();
        message.success('模型发布成功');
      }
      
      publishVisible.value = false;
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};

// 取消弹窗
const handleCancel = (type) => {
  if (type === 'import') {
    importVisible.value = false;
  } else if (type === 'publish') {
    publishVisible.value = false;
  }
};

// License选项
const licenseOptions = [
  { label: 'MIT', value: 'MIT' },
  { label: 'Apache-2.0', value: 'Apache-2.0' },
  { label: 'GPL-3.0', value: 'GPL-3.0' },
  { label: 'BSD-3-Clause', value: 'BSD-3-Clause' },
];
</script>

<template>
  <div class="model-registry-container">
    <Tabs v-model:activeKey="activeKey">
      <Tabs.TabPane key="public" tab="公共模型仓库">
        <Card title="公共模型仓库">
          <!-- 搜索区域 -->
          <div class="search-area">
            <Form layout="inline" :model="searchForm">
              <Form.Item name="keywords" label="">
                <Input placeholder="搜索模型" prefix={<SearchOutlined />} style="width: 250px" />
              </Form.Item>
              <Form.Item name="type" label="">
                <Select 
                  v-model:value="searchForm.type" 
                  placeholder="模型类型" 
                  style="width: 120px" 
                  allowClear
                >
                  <Select.Option value="NLP">NLP</Select.Option>
                  <Select.Option value="图像分类">图像分类</Select.Option>
                  <Select.Option value="目标检测">目标检测</Select.Option>
                  <Select.Option value="语音识别">语音识别</Select.Option>
                </Select>
              </Form.Item>
              <Form.Item name="framework" label="">
                <Select 
                  v-model:value="searchForm.framework" 
                  placeholder="框架" 
                  style="width: 120px" 
                  allowClear
                >
                  <Select.Option value="PyTorch">PyTorch</Select.Option>
                  <Select.Option value="TensorFlow">TensorFlow</Select.Option>
                  <Select.Option value="Keras">Keras</Select.Option>
                </Select>
              </Form.Item>
              <Form.Item name="sort" label="排序">
                <Radio.Group v-model:value="searchForm.sort">
                  <Radio.Button value="downloads">下载量</Radio.Button>
                  <Radio.Button value="stars">收藏数</Radio.Button>
                  <Radio.Button value="date">日期</Radio.Button>
                </Radio.Group>
              </Form.Item>
              <Form.Item>
                <Button type="primary" @click="handleSearch">搜索</Button>
              </Form.Item>
            </Form>
          </div>

          <Divider style="margin: 12px 0" />

          <Table :columns="publicColumns" :dataSource="registryModels" rowKey="id" />
        </Card>
      </Tabs.TabPane>

      <Tabs.TabPane key="my" tab="我的模型仓库">
        <Card title="我的模型仓库">
          <template #extra>
            <Button type="primary">上传新模型</Button>
          </template>
          <Table :columns="myColumns" :dataSource="myRegistryModels" rowKey="id" />
        </Card>
      </Tabs.TabPane>
    </Tabs>

    <!-- 导入模型弹窗 -->
    <Modal
      v-model:visible="importVisible"
      title="导入模型"
      @ok="handleImportOk"
      @cancel="() => handleCancel('import')"
      :maskClosable="false"
    >
      <div v-if="importModel" class="import-model-info">
        <h3>{{ importModel.name }}</h3>
        <p><Tag :color="
          importModel.type === 'NLP' ? 'blue' : 
          importModel.type === '图像分类' ? 'green' : 
          importModel.type === '目标检测' ? 'purple' : 
          'default'">{{ importModel.type }}</Tag> | {{ importModel.framework }}</p>
        <p>{{ importModel.description }}</p>
        <Divider />
      </div>

      <Form ref="importFormRef" :model="importFormState" layout="vertical">
        <Form.Item
          name="name"
          label="导入后的名称"
          :rules="[{ required: true, message: '请输入模型名称' }]"
        >
          <Input v-model:value="importFormState.name" placeholder="请输入模型名称" />
        </Form.Item>
        <Form.Item name="shared" label="是否共享">
          <Select v-model:value="importFormState.shared">
            <Select.Option :value="true">是</Select.Option>
            <Select.Option :value="false">否</Select.Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>

    <!-- 发布模型弹窗 -->
    <Modal
      v-model:visible="publishVisible"
      title="发布模型"
      @ok="handlePublishOk"
      @cancel="() => handleCancel('publish')"
      :maskClosable="false"
    >
      <div v-if="publishingModel" class="publish-model-info">
        <h3>{{ publishingModel.name }} ({{ publishingModel.version }})</h3>
        <Divider />
      </div>

      <Form ref="publishFormRef" :model="publishFormState" layout="vertical">
        <Form.Item
          name="isPublic"
          label="发布类型"
          :rules="[{ required: true, message: '请选择发布类型' }]"
        >
          <Radio.Group v-model:value="publishFormState.isPublic">
            <Radio :value="true">公开</Radio>
            <Radio :value="false">私有</Radio>
          </Radio.Group>
        </Form.Item>
        <Form.Item
          name="license"
          label="许可证"
          :rules="[{ required: true, message: '请选择许可证' }]"
        >
          <Select v-model:value="publishFormState.license" placeholder="请选择许可证">
            <Select.Option
              v-for="option in licenseOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="description"
          label="描述"
          :rules="[{ required: true, message: '请输入模型描述' }]"
        >
          <Input.TextArea
            v-model:value="publishFormState.description"
            placeholder="请输入模型描述"
            rows={4}
          />
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.model-registry-container {
  padding: 0;
}
.search-area {
  margin-bottom: 16px;
}
.import-model-info, .publish-model-info {
  margin-bottom: 16px;
}
</style>
