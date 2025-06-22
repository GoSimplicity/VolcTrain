<script lang="ts" setup>
import { ref, reactive } from 'vue';
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
  Descriptions,
  Timeline,
  Tabs,
  Collapse,
  Tooltip,
} from 'ant-design-vue';

defineOptions({ name: 'ModelVersion' });

// 模型数据
const modelOptions = [
  { label: 'BERT-base-chinese', value: 'BERT-base-chinese' },
  { label: 'ResNet50-ImageNet', value: 'ResNet50-ImageNet' },
  { label: 'GPT-2-small-zh', value: 'GPT-2-small-zh' },
  { label: 'YOLOv5-Custom', value: 'YOLOv5-Custom' },
  { label: 'LSTM-Sentiment', value: 'LSTM-Sentiment' },
];

// 选中的模型
const selectedModel = ref('BERT-base-chinese');

// 版本列表数据
const versionList = reactive({
  'BERT-base-chinese': [
    {
      id: '1',
      version: '1.0.3',
      description: '修复中文分词错误',
      creator: '张三',
      createTime: '2023-11-10 15:30',
      size: '375MB',
      status: '已发布',
      isCurrent: true,
      metrics: {
        accuracy: '92.5%',
        precision: '93.1%',
        recall: '91.8%',
        f1: '92.4%',
      },
      trainParams: {
        batchSize: 32,
        epochs: 5,
        learningRate: '2e-5',
        optimizer: 'Adam',
      },
      changelog: '修复了中文分词错误，优化了长文本处理性能',
    },
    {
      id: '2',
      version: '1.0.2',
      description: '优化模型性能',
      creator: '张三',
      createTime: '2023-10-25 09:45',
      size: '374MB',
      status: '存档',
      isCurrent: false,
      metrics: {
        accuracy: '91.8%',
        precision: '92.3%',
        recall: '91.5%',
        f1: '91.9%',
      },
      trainParams: {
        batchSize: 32,
        epochs: 5,
        learningRate: '2e-5',
        optimizer: 'Adam',
      },
      changelog: '优化了模型推理性能，减少了内存占用',
    },
    {
      id: '3',
      version: '1.0.1',
      description: '初始版本',
      creator: '张三',
      createTime: '2023-10-05 14:20',
      size: '380MB',
      status: '存档',
      isCurrent: false,
      metrics: {
        accuracy: '90.2%',
        precision: '91.0%',
        recall: '89.5%',
        f1: '90.2%',
      },
      trainParams: {
        batchSize: 32,
        epochs: 4,
        learningRate: '3e-5',
        optimizer: 'Adam',
      },
      changelog: '初始版本发布',
    },
  ],
  'ResNet50-ImageNet': [
    {
      id: '1',
      version: '2.1.0',
      description: '增加数据增强支持',
      creator: '李四',
      createTime: '2023-11-12 10:15',
      size: '98MB',
      status: '已发布',
      isCurrent: true,
      metrics: {
        accuracy: '76.5%',
        precision: '77.2%',
        recall: '75.8%',
        f1: '76.5%',
      },
      trainParams: {
        batchSize: 64,
        epochs: 90,
        learningRate: '0.1',
        optimizer: 'SGD',
      },
      changelog: '增加了更多数据增强支持，提高了模型泛化能力',
    },
    {
      id: '2',
      version: '2.0.0',
      description: '重新训练模型结构',
      creator: '李四',
      createTime: '2023-10-30 15:40',
      size: '97MB',
      status: '存档',
      isCurrent: false,
      metrics: {
        accuracy: '75.1%',
        precision: '76.0%',
        recall: '74.5%',
        f1: '75.2%',
      },
      trainParams: {
        batchSize: 64,
        epochs: 90,
        learningRate: '0.1',
        optimizer: 'SGD',
      },
      changelog: '使用新的ResNet50架构重新训练，提高了模型性能',
    },
  ],
  'GPT-2-small-zh': [
    {
      id: '1',
      version: '0.2.5',
      description: '中文语料扩充',
      creator: '王五',
      createTime: '2023-11-14 16:30',
      size: '510MB',
      status: '训练中',
      isCurrent: true,
      metrics: {
        perplexity: '32.1',
      },
      trainParams: {
        batchSize: 16,
        epochs: 3,
        learningRate: '5e-5',
        optimizer: 'AdamW',
      },
      changelog: '扩充中文语料库，优化模型生成能力',
    },
  ],
  'YOLOv5-Custom': [
    {
      id: '1',
      version: '1.5.0',
      description: '检测新目标类型',
      creator: '赵六',
      createTime: '2023-11-05 09:30',
      size: '87MB',
      status: '已发布',
      isCurrent: true,
      metrics: {
        mAP: '68.5%',
        precision: '72.1%',
        recall: '65.8%',
        f1: '68.8%',
      },
      trainParams: {
        batchSize: 16,
        epochs: 100,
        learningRate: '0.01',
        optimizer: 'SGD',
      },
      changelog: '增加了3种新的目标检测类型，优化了小目标检测性能',
    },
    {
      id: '2',
      version: '1.0.0',
      description: '初始版本',
      creator: '赵六',
      createTime: '2023-10-20 14:15',
      size: '85MB',
      status: '存档',
      isCurrent: false,
      metrics: {
        mAP: '65.2%',
        precision: '70.5%',
        recall: '61.2%',
        f1: '65.5%',
      },
      trainParams: {
        batchSize: 16,
        epochs: 100,
        learningRate: '0.01',
        optimizer: 'SGD',
      },
      changelog: '初始训练版本，基于YOLOv5s架构',
    },
  ],
  'LSTM-Sentiment': [
    {
      id: '1',
      version: '0.9.1',
      description: '优化模型结构',
      creator: '张三',
      createTime: '2023-11-15 11:15',
      size: '45MB',
      status: '测试中',
      isCurrent: true,
      metrics: {
        accuracy: '85.3%',
        precision: '86.1%',
        recall: '84.5%',
        f1: '85.3%',
      },
      trainParams: {
        batchSize: 32,
        epochs: 10,
        learningRate: '0.001',
        optimizer: 'Adam',
      },
      changelog: '优化LSTM层结构，添加了注意力机制',
    },
  ],
});

// 当前显示的版本列表
const currentVersionList = ref(versionList[selectedModel.value]);

// 表格列配置
const columns = [
  {
    title: '版本号',
    dataIndex: 'version',
    key: 'version',
    customRender: ({ text, record }) => (
      <Space>
        {text}
        {record.isCurrent && <Tag color="green">当前版本</Tag>}
      </Space>
    ),
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
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
    title: '大小',
    dataIndex: 'size',
    key: 'size',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      const color = 
        text === '已发布' ? 'green' : 
        text === '训练中' ? 'blue' :
        text === '测试中' ? 'orange' : 
        text === '存档' ? 'default' : 'default';
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }) => (
      <Space size="middle">
        <a onClick={() => showVersionDetail(record)}>详情</a>
        {!record.isCurrent && record.status === '存档' && <a>设为当前版本</a>}
        {record.status === '已发布' && <a>下载</a>}
        {record.status !== '训练中' && <a>比较</a>}
      </Space>
    ),
  },
];

// 版本详情弹窗
const detailVisible = ref(false);
const selectedVersion = ref(null);

// 新版本弹窗
const newVersionVisible = ref(false);
const formRef = ref();
const formState = ref({
  version: '',
  description: '',
  baseVersion: '',
});

// 处理模型选择变更
const handleModelChange = (value) => {
  selectedModel.value = value;
  currentVersionList.value = versionList[value] || [];
};

// 显示版本详情
const showVersionDetail = (version) => {
  selectedVersion.value = version;
  detailVisible.value = true;
};

// 显示新版本弹窗
const showNewVersionModal = () => {
  // 获取当前版本作为基准版本
  const currentVersion = currentVersionList.value.find(v => v.isCurrent);
  formState.value = {
    version: incrementVersion(currentVersion ? currentVersion.version : '0.0.0'),
    description: '',
    baseVersion: currentVersion ? currentVersion.version : '',
  };
  newVersionVisible.value = true;
};

// 版本号递增
const incrementVersion = (version) => {
  const parts = version.split('.').map(Number);
  parts[2] = (parts[2] || 0) + 1;
  return parts.join('.');
};

// 取消弹窗
const handleCancel = (type) => {
  if (type === 'detail') {
    detailVisible.value = false;
  } else if (type === 'new') {
    newVersionVisible.value = false;
  }
};

// 处理新建版本
const handleNewVersionOk = () => {
  formRef.value
    .validate()
    .then(() => {
      // 添加新版本
      const newVersion = {
        id: String(currentVersionList.value.length + 1),
        version: formState.value.version,
        description: formState.value.description,
        creator: '当前用户',
        createTime: new Date().toLocaleString(),
        size: '-',
        status: '创建中',
        isCurrent: false,
        metrics: {},
        trainParams: {},
        changelog: formState.value.description,
      };
      
      // 更新数据
      versionList[selectedModel.value].unshift(newVersion);
      currentVersionList.value = [...versionList[selectedModel.value]];
      
      newVersionVisible.value = false;
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};

const activeKey = ref('1');
</script>

<template>
  <div class="model-version-container">
    <Card title="模型版本管理">
      <!-- 模型选择和按钮 -->
      <div class="model-select-bar">
        <Space>
          <span>选择模型:</span>
          <Select
            v-model:value="selectedModel"
            style="width: 200px"
            @change="handleModelChange"
          >
            <Select.Option
              v-for="option in modelOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
          <Button type="primary" @click="showNewVersionModal">创建新版本</Button>
        </Space>
      </div>

      <Divider style="margin: 16px 0" />

      <!-- 版本列表 -->
      <Table
        :columns="columns"
        :dataSource="currentVersionList"
        rowKey="id"
        :pagination="false"
      />
    </Card>

    <!-- 版本详情弹窗 -->
    <Modal
      v-model:visible="detailVisible"
      title="版本详情"
      width="800px"
      :footer="null"
      @cancel="() => handleCancel('detail')"
    >
      <div v-if="selectedVersion">
        <div class="version-header">
          <h2>{{ selectedModel }} - {{ selectedVersion.version }}</h2>
          <Space>
            <Tag color="blue">{{ selectedVersion.status }}</Tag>
            {selectedVersion.isCurrent && <Tag color="green">当前版本</Tag>}
          </Space>
        </div>

        <Divider />

        <Tabs v-model:activeKey="activeKey">
          <Tabs.TabPane key="1" tab="基本信息">
            <Descriptions bordered column={2}>
              <Descriptions.Item label="版本号">{{ selectedVersion.version }}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag :color="
                  selectedVersion.status === '已发布' ? 'green' : 
                  selectedVersion.status === '训练中' ? 'blue' :
                  selectedVersion.status === '测试中' ? 'orange' : 
                  'default'
                ">{{ selectedVersion.status }}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="描述">{{ selectedVersion.description }}</Descriptions.Item>
              <Descriptions.Item label="创建者">{{ selectedVersion.creator }}</Descriptions.Item>
              <Descriptions.Item label="创建时间">{{ selectedVersion.createTime }}</Descriptions.Item>
              <Descriptions.Item label="模型大小">{{ selectedVersion.size }}</Descriptions.Item>
            </Descriptions>

            <Divider>变更记录</Divider>
            <div class="changelog">
              <p>{{ selectedVersion.changelog || '无变更记录' }}</p>
            </div>
          </Tabs.TabPane>

          <Tabs.TabPane key="2" tab="性能指标">
            <div v-if="Object.keys(selectedVersion.metrics || {}).length > 0">
              <Descriptions bordered>
                <Descriptions.Item v-for="(value, key) in selectedVersion.metrics" :key="key" :label="key">
                  {{ value }}
                </Descriptions.Item>
              </Descriptions>
            </div>
            <div v-else class="empty-metrics">
              <p>暂无性能指标数据</p>
            </div>
          </Tabs.TabPane>

          <Tabs.TabPane key="3" tab="训练参数">
            <div v-if="Object.keys(selectedVersion.trainParams || {}).length > 0">
              <Descriptions bordered>
                <Descriptions.Item v-for="(value, key) in selectedVersion.trainParams" :key="key" :label="key">
                  {{ value }}
                </Descriptions.Item>
              </Descriptions>
            </div>
            <div v-else class="empty-metrics">
              <p>暂无训练参数数据</p>
            </div>
          </Tabs.TabPane>

          <Tabs.TabPane key="4" tab="版本历史">
            <Timeline mode="left">
              <Timeline.Item v-for="version in versionList[selectedModel]" :key="version.id" :color="version.id === selectedVersion.id ? 'green' : 'blue'">
                <p>
                  <strong>{{ version.version }}</strong>
                  <span style="margin-left: 8px;">{{ version.createTime }}</span>
                  {version.isCurrent && <Tag color="green" style="margin-left: 8px;">当前版本</Tag>}
                </p>
                <p>{{ version.description }}</p>
              </Timeline.Item>
            </Timeline>
          </Tabs.TabPane>
        </Tabs>

        <div style="margin-top: 24px; text-align: right;">
          <Space>
            {!selectedVersion.isCurrent && selectedVersion.status === '已发布' && <Button type="primary">设为当前版本</Button>}
            {selectedVersion.status === '已发布' && <Button>下载模型</Button>}
            <Button @click="() => handleCancel('detail')">关闭</Button>
          </Space>
        </div>
      </div>
    </Modal>

    <!-- 新建版本弹窗 -->
    <Modal
      v-model:visible="newVersionVisible"
      title="创建新版本"
      @ok="handleNewVersionOk"
      @cancel="() => handleCancel('new')"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="version"
          label="版本号"
          :rules="[{ required: true, message: '请输入版本号' }]"
        >
          <Input v-model:value="formState.version" placeholder="请输入版本号，如1.0.0" />
        </Form.Item>
        <Form.Item
          name="baseVersion"
          label="基于版本"
        >
          <Input v-model:value="formState.baseVersion" disabled />
        </Form.Item>
        <Form.Item
          name="description"
          label="版本描述"
          :rules="[{ required: true, message: '请输入版本描述' }]"
        >
          <Input.TextArea
            v-model:value="formState.description"
            placeholder="请输入版本描述"
          />
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.model-version-container {
  padding: 0;
}
.model-select-bar {
  margin-bottom: 16px;
}
.version-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.changelog, .empty-metrics {
  background-color: #fafafa;
  padding: 16px;
  border: 1px solid #f0f0f0;
  border-radius: 2px;
}
.empty-metrics {
  text-align: center;
  color: #999;
  padding: 32px;
}
</style>
