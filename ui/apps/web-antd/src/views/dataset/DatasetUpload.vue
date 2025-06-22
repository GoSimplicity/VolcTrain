<script lang="ts" setup>
import { ref, reactive } from 'vue';
import {
  Card,
  Steps,
  Button,
  Upload,
  Form,
  Input,
  Select,
  Divider,
  Progress,
  Table,
  Tag,
  Space,
  message,
} from 'ant-design-vue';
import { UploadOutlined, InboxOutlined, CheckCircleOutlined, LoadingOutlined } from '@ant-design/icons-vue';

defineOptions({ name: 'DatasetUpload' });

// 步骤状态
const currentStep = ref(0);
const uploadStatus = ref('waiting'); // waiting, uploading, success, error

// 表单数据
const formRef = ref();
const formState = ref({
  name: '',
  type: '',
  description: '',
  shared: false,
  labelRequired: false,
});

// 上传文件列表
const fileList = ref([]);
const uploadProgress = ref(0);

// 上传文件信息表格数据
const uploadedFiles = ref([]);
const uploadingFile = ref(null);
const uploadTotal = ref(0);
const uploadedSize = ref(0);

// 类型选项
const typeOptions = [
  { label: '图像分类', value: '图像分类' },
  { label: '文本', value: '文本' },
  { label: '音频', value: '音频' },
  { label: '视频', value: '视频' },
  { label: '表格', value: '表格' },
];

// 步骤配置
const steps = [
  {
    title: '填写基本信息',
    status: 'process',
  },
  {
    title: '上传文件',
    status: 'wait',
  },
  {
    title: '处理设置',
    status: 'wait',
  },
  {
    title: '完成',
    status: 'wait',
  },
];

// 文件列表表格列配置
const fileColumns = [
  {
    title: '文件名',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '大小',
    dataIndex: 'size',
    key: 'size',
    customRender: ({ text }) => {
      // 将字节转换为更易读的格式
      if (text < 1024) {
        return `${text}B`;
      } else if (text < 1024 * 1024) {
        return `${(text / 1024).toFixed(2)}KB`;
      } else if (text < 1024 * 1024 * 1024) {
        return `${(text / 1024 / 1024).toFixed(2)}MB`;
      } else {
        return `${(text / 1024 / 1024 / 1024).toFixed(2)}GB`;
      }
    },
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      const color = 
        text === '已上传' ? 'green' : 
        text === '上传中' ? 'blue' : 
        text === '等待上传' ? 'orange' : 
        text === '上传失败' ? 'red' : 'default';
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '进度',
    dataIndex: 'progress',
    key: 'progress',
    customRender: ({ text, record }) => {
      if (record.status === '已上传') {
        return <Progress percent={100} size="small" status="success" />;
      } else if (record.status === '上传中') {
        return <Progress percent={text} size="small" status="active" />;
      } else if (record.status === '上传失败') {
        return <Progress percent={text} size="small" status="exception" />;
      } else {
        return <Progress percent={0} size="small" />;
      }
    },
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }) => (
      <Space size="middle">
        {record.status === '上传失败' && <a>重试</a>}
        {record.status !== '上传中' && <a>删除</a>}
      </Space>
    ),
  },
];

// 处理步骤切换
const nextStep = () => {
  if (currentStep.value === 0) {
    // 验证表单
    formRef.value.validate().then(() => {
      currentStep.value++;
      steps[currentStep.value].status = 'process';
    }).catch(error => {
      console.error('表单验证失败:', error);
    });
  } else if (currentStep.value === 1) {
    // 检查是否有文件上传
    if (fileList.value.length === 0) {
      message.error('请至少上传一个文件');
      return;
    }
    // 模拟开始上传
    uploadStatus.value = 'uploading';
    startUpload();
  } else if (currentStep.value === 2) {
    // 完成处理设置，进入完成阶段
    currentStep.value++;
    steps[currentStep.value].status = 'process';
  }
};

const prevStep = () => {
  if (currentStep.value > 0) {
    steps[currentStep.value].status = 'wait';
    currentStep.value--;
  }
};

// 自定义上传处理
const handleChange = info => {
  fileList.value = [...info.fileList];
  
  // 将文件添加到上传队列
  const newFiles = info.fileList
    .filter(file => !uploadedFiles.value.some(f => f.uid === file.uid))
    .map(file => ({
      uid: file.uid,
      name: file.name,
      size: file.size,
      type: file.type,
      status: '等待上传',
      progress: 0,
    }));
  
  uploadedFiles.value = [...uploadedFiles.value, ...newFiles];
  
  // 计算总上传大小
  uploadTotal.value = uploadedFiles.value.reduce((total, file) => total + file.size, 0);
};

// 模拟文件上传过程
const startUpload = () => {
  // 找到第一个等待上传的文件
  const fileToUpload = uploadedFiles.value.find(file => file.status === '等待上传');
  
  if (!fileToUpload) {
    // 所有文件已上传完成
    uploadStatus.value = 'success';
    currentStep.value++;
    steps[currentStep.value].status = 'process';
    return;
  }
  
  // 更新状态为上传中
  const index = uploadedFiles.value.findIndex(file => file.uid === fileToUpload.uid);
  uploadedFiles.value[index].status = '上传中';
  uploadingFile.value = uploadedFiles.value[index];
  
  // 模拟上传进度
  let progress = 0;
  const timer = setInterval(() => {
    progress += Math.floor(Math.random() * 10) + 1;
    
    if (progress >= 100) {
      progress = 100;
      clearInterval(timer);
      
      // 更新当前文件状态
      uploadedFiles.value[index].status = '已上传';
      uploadedFiles.value[index].progress = 100;
      
      // 更新已上传大小
      uploadedSize.value += uploadedFiles.value[index].size;
      
      // 计算总进度
      uploadProgress.value = Math.floor((uploadedSize.value / uploadTotal.value) * 100);
      
      // 继续上传下一个文件
      setTimeout(() => {
        startUpload();
      }, 500);
    } else {
      // 更新当前文件进度
      uploadedFiles.value[index].progress = progress;
      
      // 计算总进度（已上传文件 + 当前上传进度的一部分）
      const currentProgress = (uploadedSize.value + uploadedFiles.value[index].size * (progress / 100)) / uploadTotal.value;
      uploadProgress.value = Math.floor(currentProgress * 100);
    }
  }, 200);
};

// 完成所有步骤后跳转到数据集列表
const handleFinish = () => {
  message.success('数据集上传成功！');
  // 这里通常会跳转到数据集列表页
};
</script>

<template>
  <div class="dataset-upload-container">
    <Card title="数据集上传">
      <!-- 步骤条 -->
      <Steps :current="currentStep" style="margin-bottom: 30px">
        <Steps.Step 
          v-for="(step, index) in steps" 
          :key="index" 
          :title="step.title" 
          :status="index === currentStep ? 'process' : index < currentStep ? 'finish' : 'wait'"
        />
      </Steps>

      <!-- 步骤内容 -->
      <div class="steps-content">
        <!-- 步骤1：基本信息 -->
        <div v-if="currentStep === 0" class="step-form">
          <Form ref="formRef" :model="formState" layout="vertical">
            <Form.Item
              name="name"
              label="数据集名称"
              :rules="[{ required: true, message: '请输入数据集名称' }]"
            >
              <Input v-model:value="formState.name" placeholder="请输入数据集名称" />
            </Form.Item>
            <Form.Item
              name="type"
              label="数据类型"
              :rules="[{ required: true, message: '请选择数据类型' }]"
            >
              <Select v-model:value="formState.type" placeholder="请选择数据类型">
                <Select.Option
                  v-for="option in typeOptions"
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
                placeholder="请输入数据集描述"
              />
            </Form.Item>
            <Form.Item name="shared" label="是否共享">
              <Select v-model:value="formState.shared">
                <Select.Option :value="true">是</Select.Option>
                <Select.Option :value="false">否</Select.Option>
              </Select>
            </Form.Item>
            <Form.Item name="labelRequired" label="是否需要标注">
              <Select v-model:value="formState.labelRequired">
                <Select.Option :value="true">是</Select.Option>
                <Select.Option :value="false">否</Select.Option>
              </Select>
            </Form.Item>
          </Form>
        </div>

        <!-- 步骤2：文件上传 -->
        <div v-else-if="currentStep === 1" class="step-upload">
          <div class="upload-area" v-if="uploadStatus === 'waiting'">
            <Upload.Dragger
              v-model:fileList="fileList"
              multiple
              :beforeUpload="() => false"
              @change="handleChange"
            >
              <p class="ant-upload-drag-icon">
                <InboxOutlined />
              </p>
              <p class="ant-upload-text">点击或拖拽文件到此区域上传</p>
              <p class="ant-upload-hint">
                支持单个或批量上传。严禁上传公司内部资料或其他违禁文件。
              </p>
            </Upload.Dragger>
          </div>

          <div class="upload-progress" v-if="uploadStatus === 'uploading'">
            <Progress 
              :percent="uploadProgress" 
              :status="uploadProgress < 100 ? 'active' : 'success'" 
              :showInfo="true"
            />
            <p>正在上传 {{ uploadedFiles.length }} 个文件，
               总大小 {{ Math.round(uploadTotal / 1024 / 1024 * 100) / 100 }}MB</p>
          </div>

          <Divider />

          <div class="file-list">
            <Table 
              :columns="fileColumns" 
              :dataSource="uploadedFiles"
              rowKey="uid"
              :pagination="false"
            />
          </div>
        </div>

        <!-- 步骤3：处理设置 -->
        <div v-else-if="currentStep === 2" class="step-process">
          <div class="process-settings">
            <h3>数据处理设置</h3>
            <Form layout="vertical">
              <Form.Item label="是否进行数据清洗">
                <Select defaultValue="true" style="width: 200px">
                  <Select.Option value="true">是</Select.Option>
                  <Select.Option value="false">否</Select.Option>
                </Select>
              </Form.Item>
              <Form.Item label="是否进行数据增强">
                <Select defaultValue="false" style="width: 200px">
                  <Select.Option value="true">是</Select.Option>
                  <Select.Option value="false">否</Select.Option>
                </Select>
              </Form.Item>
              <Form.Item label="是否切分训练测试集">
                <Select defaultValue="true" style="width: 200px">
                  <Select.Option value="true">是</Select.Option>
                  <Select.Option value="false">否</Select.Option>
                </Select>
              </Form.Item>
              <Form.Item label="训练集比例">
                <Input style="width: 200px" defaultValue="0.8" addonAfter="%" />
              </Form.Item>
            </Form>
          </div>
        </div>

        <!-- 步骤4：完成 -->
        <div v-else class="step-finish">
          <div class="finish-content">
            <div class="success-icon">
              <CheckCircleOutlined style="font-size: 72px; color: #52c41a;" />
            </div>
            <h2 style="margin-top: 24px; margin-bottom: 16px;">数据集上传成功！</h2>
            <p>数据集 <b>{{ formState.name }}</b> 已成功上传，您可以在数据集列表中查看和管理。</p>
            <div style="margin-top: 36px;">
              <Button type="primary" @click="handleFinish">返回数据集列表</Button>
            </div>
          </div>
        </div>
      </div>

      <!-- 步骤按钮 -->
      <div class="steps-action">
        <Button
          v-if="currentStep > 0 && currentStep < 3"
          style="margin-right: 8px"
          @click="prevStep"
        >
          上一步
        </Button>
        <Button
          v-if="currentStep < 3"
          type="primary"
          @click="nextStep"
          :loading="uploadStatus === 'uploading'"
        >
          {{ currentStep < 1 ? '下一步' : currentStep === 1 ? '开始上传' : '完成设置' }}
        </Button>
      </div>
    </Card>
  </div>
</template>

<style scoped>
.dataset-upload-container {
  padding: 0;
}
.steps-content {
  margin-top: 16px;
  padding: 20px;
  background-color: #fafafa;
  border: 1px dashed #e9e9e9;
  border-radius: 2px;
  min-height: 400px;
  margin-bottom: 16px;
}
.steps-action {
  margin-top: 24px;
}
.step-form {
  max-width: 600px;
}
.file-list {
  margin-top: 16px;
}
.upload-progress {
  margin: 20px 0;
}
.success-icon {
  margin-bottom: 24px;
}
.finish-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
}
</style>
