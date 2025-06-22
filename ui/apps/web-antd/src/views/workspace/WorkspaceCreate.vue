<script lang="ts" setup>
import { ref } from 'vue';
import {
  Card,
  Form,
  Input,
  Button,
  Select,
  Divider,
  message,
  Steps,
  Upload,
} from 'ant-design-vue';
import { UploadOutlined } from '@ant-design/icons-vue';

defineOptions({ name: 'WorkspaceCreate' });

const formRef = ref();
const currentStep = ref(0);
const formState = ref({
  name: '',
  description: '',
  type: 'personal',
  template: 'empty',
});

const fileList = ref([]);
const uploading = ref(false);

const workspaceTypes = [
  { label: '个人工作空间', value: 'personal' },
  { label: '团队工作空间', value: 'team' },
  { label: '项目工作空间', value: 'project' },
];

const templates = [
  { label: '空白工作空间', value: 'empty' },
  { label: '机器学习工作空间', value: 'ml' },
  { label: '深度学习工作空间', value: 'dl' },
  { label: '数据分析工作空间', value: 'data' },
];

const steps = [
  {
    title: '基本信息',
    content: 'first-content',
  },
  {
    title: '配置选项',
    content: 'second-content',
  },
  {
    title: '确认创建',
    content: 'third-content',
  },
];

const nextStep = () => {
  if (currentStep.value === 0) {
    formRef.value
      .validateFields(['name', 'description'])
      .then(() => {
        currentStep.value++;
      })
      .catch((error) => {
        console.log('验证失败:', error);
      });
  } else {
    currentStep.value++;
  }
};

const prevStep = () => {
  currentStep.value--;
};

const handleFinish = () => {
  message.success('工作空间创建成功！');
  // 这里可以添加跳转到工作空间列表页的逻辑
};
</script>

<template>
  <div class="workspace-create-container">
    <Card title="创建工作空间" :bordered="false">
      <Steps :current="currentStep" style="margin-bottom: 30px">
        <Steps.Step v-for="(item, index) in steps" :key="index" :title="item.title" />
      </Steps>

      <div class="steps-content">
        <div v-if="currentStep === 0" class="step-1">
          <Form ref="formRef" :model="formState" layout="vertical">
            <Form.Item
              name="name"
              label="工作空间名称"
              :rules="[{ required: true, message: '请输入工作空间名称' }]"
            >
              <Input v-model:value="formState.name" placeholder="请输入工作空间名称" />
            </Form.Item>
            <Form.Item
              name="description"
              label="工作空间描述"
              :rules="[{ required: true, message: '请输入工作空间描述' }]"
            >
              <Input.TextArea
                v-model:value="formState.description"
                placeholder="请输入工作空间描述"
                :rows="4"
              />
            </Form.Item>
          </Form>
        </div>

        <div v-else-if="currentStep === 1" class="step-2">
          <Form :model="formState" layout="vertical">
            <Form.Item name="type" label="工作空间类型">
              <Select
                v-model:value="formState.type"
                style="width: 100%"
                placeholder="请选择工作空间类型"
              >
                <Select.Option
                  v-for="type in workspaceTypes"
                  :key="type.value"
                  :value="type.value"
                >
                  {{ type.label }}
                </Select.Option>
              </Select>
            </Form.Item>
            <Form.Item name="template" label="工作空间模板">
              <Select
                v-model:value="formState.template"
                style="width: 100%"
                placeholder="请选择工作空间模板"
              >
                <Select.Option
                  v-for="template in templates"
                  :key="template.value"
                  :value="template.value"
                >
                  {{ template.label }}
                </Select.Option>
              </Select>
            </Form.Item>
            <Form.Item name="files" label="上传初始文件（可选）">
              <Upload v-model:fileList="fileList" :multiple="true">
                <Button><UploadOutlined />上传文件</Button>
              </Upload>
            </Form.Item>
          </Form>
        </div>

        <div v-else class="step-3">
          <h3>确认工作空间信息</h3>
          <Divider />
          <p><strong>名称：</strong> {{ formState.name }}</p>
          <p><strong>描述：</strong> {{ formState.description }}</p>
          <p>
            <strong>类型：</strong>
            {{ workspaceTypes.find((t) => t.value === formState.type)?.label }}
          </p>
          <p>
            <strong>模板：</strong>
            {{ templates.find((t) => t.value === formState.template)?.label }}
          </p>
          <p>
            <strong>上传文件：</strong>
            {{ fileList.length > 0 ? `${fileList.length}个文件` : '无' }}
          </p>
          <Divider />
          <p>点击"完成"按钮创建工作空间</p>
        </div>
      </div>

      <div class="steps-action">
        <Button v-if="currentStep > 0" style="margin-right: 8px" @click="prevStep">
          上一步
        </Button>
        <Button
          v-if="currentStep < steps.length - 1"
          type="primary"
          @click="nextStep"
        >
          下一步
        </Button>
        <Button
          v-if="currentStep === steps.length - 1"
          type="primary"
          @click="handleFinish"
        >
          完成
        </Button>
      </div>
    </Card>
  </div>
</template>

<style scoped>
.workspace-create-container {
  padding: 24px;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
}
.steps-content {
  margin-top: 16px;
  padding: 20px;
  background-color: #fafafa;
  border: 1px dashed #e9e9e9;
  border-radius: 2px;
  min-height: 200px;
  margin-bottom: 16px;
}
.steps-action {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}
.step-1, .step-2, .step-3 {
  max-width: 800px;
  margin: 0 auto;
}
</style>
