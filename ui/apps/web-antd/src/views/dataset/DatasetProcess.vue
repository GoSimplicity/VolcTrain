<script lang="ts" setup>
import { ref } from 'vue';
import {
  Card,
  Tabs,
  Table,
  Tag,
  Space,
  Button,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  Progress,
  Alert,
  Divider,
  Descriptions,
} from 'ant-design-vue';

defineOptions({ name: 'DatasetProcess' });

// 数据集处理任务列表
const processTasks = ref([
  {
    id: '1',
    name: 'COCO数据集预处理',
    dataset: 'COCO数据集',
    type: '图像处理',
    processType: '数据清洗',
    status: '已完成',
    progress: 100,
    startTime: '2023-11-10 08:30',
    endTime: '2023-11-10 10:45',
    creator: '张三',
  },
  {
    id: '2',
    name: '中文语料库分词',
    dataset: '中文语料库',
    type: '文本处理',
    processType: '分词',
    status: '处理中',
    progress: 65,
    startTime: '2023-11-15 14:20',
    endTime: '-',
    creator: '李四',
  },
  {
    id: '3',
    name: 'ImageNet数据增强',
    dataset: 'ImageNet子集',
    type: '图像处理',
    processType: '数据增强',
    status: '等待中',
    progress: 0,
    startTime: '-',
    endTime: '-',
    creator: '王五',
  },
  {
    id: '4',
    name: '语音识别数据标准化',
    dataset: '语音识别数据集',
    type: '音频处理',
    processType: '标准化',
    status: '失败',
    progress: 45,
    startTime: '2023-11-12 09:15',
    endTime: '2023-11-12 09:48',
    creator: '赵六',
  },
]);

// 处理模板列表
const processTemplates = ref([
  {
    id: '1',
    name: '图像分类数据预处理',
    type: '图像处理',
    description: '针对图像分类任务的数据预处理模板，包含图像缩放、裁剪、数据增强等步骤',
    steps: ['图像缩放', '图像裁剪', '数据增强'],
    creator: '系统',
    createTime: '2023-10-01',
    usage: 28,
  },
  {
    id: '2',
    name: '文本分类预处理',
    type: '文本处理',
    description: '针对文本分类任务的数据预处理模板，包含分词、去停用词、向量化等步骤',
    steps: ['分词', '去停用词', '向量化'],
    creator: '系统',
    createTime: '2023-10-01',
    usage: 15,
  },
  {
    id: '3',
    name: '语音数据标准化',
    type: '音频处理',
    description: '语音数据标准化处理模板，包含音频格式统一、采样率调整、音量归一化等步骤',
    steps: ['格式转换', '采样率调整', '音量归一化'],
    creator: '张三',
    createTime: '2023-10-15',
    usage: 6,
  },
  {
    id: '4',
    name: '自定义图像预处理',
    type: '图像处理',
    description: '自定义图像预处理模板，包含多种图像处理操作',
    steps: ['图像缩放', '旋转', '颜色调整', '数据增强'],
    creator: '李四',
    createTime: '2023-11-01',
    usage: 3,
  },
]);

// 任务列表表格列配置
const taskColumns = [
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
    title: '数据类型',
    dataIndex: 'type',
    key: 'type',
  },
  {
    title: '处理类型',
    dataIndex: 'processType',
    key: 'processType',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      const color = 
        text === '已完成' ? 'green' : 
        text === '处理中' ? 'blue' : 
        text === '等待中' ? 'orange' : 
        text === '失败' ? 'red' : 'default';
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '进度',
    dataIndex: 'progress',
    key: 'progress',
    customRender: ({ text, record }) => {
      if (record.status === '已完成') {
        return <Progress percent={100} size="small" status="success" />;
      } else if (record.status === '处理中') {
        return <Progress percent={text} size="small" status="active" />;
      } else if (record.status === '失败') {
        return <Progress percent={text} size="small" status="exception" />;
      } else {
        return <Progress percent={0} size="small" />;
      }
    },
  },
  {
    title: '开始时间',
    dataIndex: 'startTime',
    key: 'startTime',
  },
  {
    title: '结束时间',
    dataIndex: 'endTime',
    key: 'endTime',
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }) => (
      <Space size="middle">
        <a onClick={() => showTaskDetail(record)}>详情</a>
        {record.status === '失败' && <a>重试</a>}
        {record.status !== '处理中' && <a>删除</a>}
        {record.status === '已完成' && <a>导出</a>}
      </Space>
    ),
  },
];

// 模板列表表格列配置
const templateColumns = [
  {
    title: '模板ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '模板名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '适用类型',
    dataIndex: 'type',
    key: 'type',
    customRender: ({ text }) => {
      const color = 
        text === '图像处理' ? 'blue' : 
        text === '文本处理' ? 'green' : 
        text === '音频处理' ? 'purple' : 'default';
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '处理步骤',
    dataIndex: 'steps',
    key: 'steps',
    customRender: ({ text }) => (
      <>
        {text.map((step, index) => (
          <Tag key={index} style={{ marginBottom: '5px' }}>{step}</Tag>
        ))}
      </>
    ),
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
    title: '使用次数',
    dataIndex: 'usage',
    key: 'usage',
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }) => (
      <Space size="middle">
        <a>使用</a>
        <a>编辑</a>
        <a>导出</a>
      </Space>
    ),
  },
];

// 新建处理任务表单
const newTaskVisible = ref(false);
const taskFormRef = ref();
const taskFormState = ref({
  name: '',
  dataset: '',
  template: '',
  description: '',
});

// 新建模板表单
const newTemplateVisible = ref(false);
const templateFormRef = ref();
const templateFormState = ref({
  name: '',
  type: '',
  description: '',
  steps: [],
});

// 任务详情弹窗
const taskDetailVisible = ref(false);
const selectedTask = ref(null);

// 数据集选项
const datasetOptions = [
  { label: 'COCO数据集', value: 'COCO数据集' },
  { label: '中文语料库', value: '中文语料库' },
  { label: 'ImageNet子集', value: 'ImageNet子集' },
  { label: '语音识别数据集', value: '语音识别数据集' },
  { label: '电商评论数据', value: '电商评论数据' },
];

// 模板选项
const templateOptions = processTemplates.value.map(template => ({
  label: template.name,
  value: template.name,
}));

// 类型选项
const typeOptions = [
  { label: '图像处理', value: '图像处理' },
  { label: '文本处理', value: '文本处理' },
  { label: '音频处理', value: '音频处理' },
  { label: '视频处理', value: '视频处理' },
];

// 处理步骤选项
const stepsOptions = [
  { label: '图像缩放', value: '图像缩放' },
  { label: '图像裁剪', value: '图像裁剪' },
  { label: '数据增强', value: '数据增强' },
  { label: '分词', value: '分词' },
  { label: '去停用词', value: '去停用词' },
  { label: '向量化', value: '向量化' },
  { label: '格式转换', value: '格式转换' },
  { label: '采样率调整', value: '采样率调整' },
  { label: '音量归一化', value: '音量归一化' },
  { label: '旋转', value: '旋转' },
  { label: '颜色调整', value: '颜色调整' },
];

// 显示任务详情
const showTaskDetail = (task) => {
  selectedTask.value = task;
  taskDetailVisible.value = true;
};

// 显示新建处理任务弹窗
const showNewTask = () => {
  newTaskVisible.value = true;
};

// 显示新建处理模板弹窗
const showNewTemplate = () => {
  newTemplateVisible.value = true;
};

// 处理新建任务
const handleNewTaskOk = () => {
  taskFormRef.value
    .validate()
    .then(() => {
      const template = processTemplates.value.find(t => t.name === taskFormState.value.template);
      
      // 添加新任务
      processTasks.value.push({
        id: String(processTasks.value.length + 1),
        name: taskFormState.value.name,
        dataset: taskFormState.value.dataset,
        type: template ? template.type : '未知',
        processType: template ? template.steps.join('、') : '自定义',
        status: '等待中',
        progress: 0,
        startTime: '-',
        endTime: '-',
        creator: '当前用户',
      });
      
      newTaskVisible.value = false;
      taskFormState.value = {
        name: '',
        dataset: '',
        template: '',
        description: '',
      };
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};

// 处理新建模板
const handleNewTemplateOk = () => {
  templateFormRef.value
    .validate()
    .then(() => {
      // 添加新模板
      processTemplates.value.push({
        id: String(processTemplates.value.length + 1),
        name: templateFormState.value.name,
        type: templateFormState.value.type,
        description: templateFormState.value.description,
        steps: templateFormState.value.steps,
        creator: '当前用户',
        createTime: new Date().toLocaleDateString(),
        usage: 0,
      });
      
      newTemplateVisible.value = false;
      templateFormState.value = {
        name: '',
        type: '',
        description: '',
        steps: [],
      };
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};

// 取消弹窗
const handleCancel = (type) => {
  if (type === 'task') {
    newTaskVisible.value = false;
  } else if (type === 'template') {
    newTemplateVisible.value = false;
  } else if (type === 'detail') {
    taskDetailVisible.value = false;
  }
};

// 标签页
const activeKey = ref('tasks');
</script>

<template>
  <div class="dataset-process-container">
    <Tabs v-model:activeKey="activeKey">
      <Tabs.TabPane key="tasks" tab="处理任务">
        <Card title="数据集处理任务">
          <template #extra>
            <Button type="primary" @click="showNewTask">新建处理任务</Button>
          </template>
          <Table :columns="taskColumns" :dataSource="processTasks" rowKey="id" />
        </Card>
      </Tabs.TabPane>
      <Tabs.TabPane key="templates" tab="处理模板">
        <Card title="数据处理模板">
          <template #extra>
            <Button type="primary" @click="showNewTemplate">新建处理模板</Button>
          </template>
          <Table :columns="templateColumns" :dataSource="processTemplates" rowKey="id" />
        </Card>
      </Tabs.TabPane>
    </Tabs>

    <!-- 新建处理任务弹窗 -->
    <Modal
      v-model:visible="newTaskVisible"
      title="新建处理任务"
      @ok="handleNewTaskOk"
      @cancel="() => handleCancel('task')"
      :maskClosable="false"
    >
      <Form ref="taskFormRef" :model="taskFormState" layout="vertical">
        <Form.Item
          name="name"
          label="任务名称"
          :rules="[{ required: true, message: '请输入任务名称' }]"
        >
          <Input v-model:value="taskFormState.name" placeholder="请输入任务名称" />
        </Form.Item>
        <Form.Item
          name="dataset"
          label="选择数据集"
          :rules="[{ required: true, message: '请选择数据集' }]"
        >
          <Select v-model:value="taskFormState.dataset" placeholder="请选择数据集">
            <Select.Option
              v-for="option in datasetOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="template"
          label="处理模板"
          :rules="[{ required: true, message: '请选择处理模板' }]"
        >
          <Select v-model:value="taskFormState.template" placeholder="请选择处理模板">
            <Select.Option
              v-for="option in templateOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item name="description" label="描述">
          <Input.TextArea
            v-model:value="taskFormState.description"
            placeholder="请输入任务描述"
          />
        </Form.Item>
      </Form>
    </Modal>

    <!-- 新建处理模板弹窗 -->
    <Modal
      v-model:visible="newTemplateVisible"
      title="新建处理模板"
      @ok="handleNewTemplateOk"
      @cancel="() => handleCancel('template')"
      :maskClosable="false"
      width="700px"
    >
      <Form ref="templateFormRef" :model="templateFormState" layout="vertical">
        <Form.Item
          name="name"
          label="模板名称"
          :rules="[{ required: true, message: '请输入模板名称' }]"
        >
          <Input v-model:value="templateFormState.name" placeholder="请输入模板名称" />
        </Form.Item>
        <Form.Item
          name="type"
          label="适用类型"
          :rules="[{ required: true, message: '请选择适用类型' }]"
        >
          <Select v-model:value="templateFormState.type" placeholder="请选择适用类型">
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
          name="steps"
          label="处理步骤"
          :rules="[{ required: true, message: '请选择处理步骤' }]"
        >
          <Select
            v-model:value="templateFormState.steps"
            placeholder="请选择处理步骤"
            mode="multiple"
          >
            <Select.Option
              v-for="option in stepsOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item name="description" label="描述">
          <Input.TextArea
            v-model:value="templateFormState.description"
            placeholder="请输入模板描述"
          />
        </Form.Item>
        <Alert 
          message="您可以根据需要自定义处理步骤顺序，系统会按照选择的顺序执行处理。" 
          type="info" 
          showIcon 
          style="margin-bottom: 16px" 
        />
      </Form>
    </Modal>

    <!-- 任务详情弹窗 -->
    <Modal
      v-model:visible="taskDetailVisible"
      title="处理任务详情"
      @cancel="() => handleCancel('detail')"
      :footer="null"
      width="800px"
    >
      <div v-if="selectedTask">
        <Descriptions title="基本信息" bordered>
          <Descriptions.Item label="任务ID">{{ selectedTask.id }}</Descriptions.Item>
          <Descriptions.Item label="任务名称">{{ selectedTask.name }}</Descriptions.Item>
          <Descriptions.Item label="处理状态">
            <Tag :color="
              selectedTask.status === '已完成' ? 'green' : 
              selectedTask.status === '处理中' ? 'blue' : 
              selectedTask.status === '等待中' ? 'orange' : 
              'red'
            ">
              {{ selectedTask.status }}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="数据集">{{ selectedTask.dataset }}</Descriptions.Item>
          <Descriptions.Item label="数据类型">{{ selectedTask.type }}</Descriptions.Item>
          <Descriptions.Item label="创建者">{{ selectedTask.creator }}</Descriptions.Item>
          <Descriptions.Item label="开始时间">{{ selectedTask.startTime }}</Descriptions.Item>
          <Descriptions.Item label="结束时间">{{ selectedTask.endTime }}</Descriptions.Item>
          <Descriptions.Item label="处理进度">
            <Progress 
              :percent="selectedTask.progress" 
              :status="
                selectedTask.status === '已完成' ? 'success' : 
                selectedTask.status === '处理中' ? 'active' : 
                selectedTask.status === '失败' ? 'exception' : 
                'normal'
              "
            />
          </Descriptions.Item>
        </Descriptions>

        <Divider />

        <h3>处理步骤</h3>
        <Descriptions bordered>
          <Descriptions.Item label="处理类型">{{ selectedTask.processType }}</Descriptions.Item>
        </Descriptions>

        <div v-if="selectedTask.status === '失败'" style="margin-top: 16px">
          <Alert
            message="处理失败原因"
            description="在处理第45个文件时发生错误：文件格式不兼容。请检查数据集中的文件格式是否符合要求，或尝试使用其他处理模板。"
            type="error"
            showIcon
          />
        </div>

        <div v-if="selectedTask.status === '已完成'" style="margin-top: 16px">
          <Alert
            message="处理结果"
            description="处理完成，共处理123287个文件，生成处理后数据集。您可以在数据集列表中查看处理结果或下载导出。"
            type="success"
            showIcon
          />
        </div>

        <div style="margin-top: 24px; text-align: right;">
          <Space>
            {selectedTask.status === '已完成' && <Button type="primary">查看结果</Button>}
            {selectedTask.status === '失败' && <Button type="primary">重新处理</Button>}
            <Button @click="() => handleCancel('detail')">关闭</Button>
          </Space>
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.dataset-process-container {
  padding: 0;
}
</style>
