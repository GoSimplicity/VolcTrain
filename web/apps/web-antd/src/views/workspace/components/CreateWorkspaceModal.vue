<template>
  <Modal
    v-model:open="visible"
    title="创建工作空间"
    width="600px"
    @ok="handleSubmit"
    @cancel="handleCancel"
    :confirm-loading="loading"
  >
    <Form
      ref="formRef"
      :model="formData"
      :rules="rules"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 18 }"
    >
      <Form.Item label="工作空间名称" name="name">
        <Input
          v-model:value="formData.name"
          placeholder="请输入工作空间名称"
          :maxlength="50"
        />
      </Form.Item>

      <Form.Item label="工作空间类型" name="type">
        <Select v-model:value="formData.type" placeholder="请选择工作空间类型">
          <Select.Option value="personal">个人工作空间</Select.Option>
          <Select.Option value="team">团队工作空间</Select.Option>
          <Select.Option value="project">项目工作空间</Select.Option>
          <Select.Option value="department">部门工作空间</Select.Option>
        </Select>
      </Form.Item>

      <Form.Item label="描述" name="description">
        <Input.TextArea
          v-model:value="formData.description"
          :rows="3"
          placeholder="请输入工作空间描述"
          :maxlength="200"
          show-count
        />
      </Form.Item>

      <Divider orientation="left">资源配额</Divider>

      <Form.Item label="CPU核心数" name="cpu">
        <InputNumber
          v-model:value="formData.resourceQuota.cpu"
          :min="1"
          :max="100"
          style="width: 100%"
          addon-after="核"
        />
      </Form.Item>

      <Form.Item label="内存大小" name="memory">
        <InputNumber
          v-model:value="formData.resourceQuota.memory"
          :min="1"
          :max="1000"
          style="width: 100%"
          addon-after="GB"
        />
      </Form.Item>

      <Form.Item label="GPU数量" name="gpu">
        <InputNumber
          v-model:value="formData.resourceQuota.gpu"
          :min="0"
          :max="16"
          style="width: 100%"
          addon-after="个"
        />
      </Form.Item>

      <Form.Item label="存储空间" name="storage">
        <InputNumber
          v-model:value="formData.resourceQuota.storage"
          :min="10"
          :max="10000"
          style="width: 100%"
          addon-after="GB"
        />
      </Form.Item>

      <Divider orientation="left">工作空间配置</Divider>

      <Form.Item label="最大成员数" name="maxMembers">
        <InputNumber
          v-model:value="formData.config.maxMembers"
          :min="1"
          :max="1000"
          style="width: 100%"
          addon-after="人"
        />
      </Form.Item>

      <Form.Item label="允许外部访问">
        <Switch v-model:checked="formData.config.allowExternalAccess" />
      </Form.Item>

      <Form.Item label="启用自动清理">
        <Switch v-model:checked="formData.config.enableAutoCleanup" />
      </Form.Item>

      <Form.Item
        v-if="formData.config.enableAutoCleanup"
        label="清理天数"
        name="cleanupDays"
      >
        <InputNumber
          v-model:value="formData.config.cleanupDays"
          :min="1"
          :max="365"
          style="width: 100%"
          addon-after="天"
        />
      </Form.Item>
    </Form>
  </Modal>
</template>

<script lang="ts" setup>
import { ref, reactive, watch } from 'vue';
import {
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  Switch,
  Divider,
  message,
} from 'ant-design-vue';
import type { CreateWorkspaceRequest, WorkspaceType } from '#/api/types';
import { createWorkspace } from '#/api';

const props = defineProps<{
  visible: boolean;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
  success: [];
}>();

const formRef = ref();
const loading = ref(false);

// 表单数据
const formData = reactive<CreateWorkspaceRequest>({
  name: '',
  description: '',
  type: 'personal' as WorkspaceType,
  resourceQuota: {
    cpu: 4,
    memory: 8,
    gpu: 1,
    storage: 100,
  },
  config: {
    allowExternalAccess: false,
    enableAutoCleanup: true,
    cleanupDays: 30,
    maxMembers: 10,
    defaultResourceQuota: {
      cpu: 2,
      memory: 4,
      gpu: 0,
      storage: 50,
    },
  },
});

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入工作空间名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度应在2-50个字符之间', trigger: 'blur' },
  ],
  type: [
    { required: true, message: '请选择工作空间类型', trigger: 'change' },
  ],
  cpu: [
    { required: true, message: '请输入CPU核心数', trigger: 'blur' },
    { type: 'number', min: 1, max: 100, message: 'CPU核心数应在1-100之间', trigger: 'blur' },
  ],
  memory: [
    { required: true, message: '请输入内存大小', trigger: 'blur' },
    { type: 'number', min: 1, max: 1000, message: '内存大小应在1-1000GB之间', trigger: 'blur' },
  ],
  maxMembers: [
    { required: true, message: '请输入最大成员数', trigger: 'blur' },
    { type: 'number', min: 1, max: 1000, message: '最大成员数应在1-1000之间', trigger: 'blur' },
  ],
  cleanupDays: [
    { type: 'number', min: 1, max: 365, message: '清理天数应在1-365之间', trigger: 'blur' },
  ],
};

// 监听visible变化，重置表单
watch(() => props.visible, (newVal) => {
  if (newVal) {
    resetForm();
  }
});

// 重置表单
const resetForm = () => {
  formData.name = '';
  formData.description = '';
  formData.type = 'personal' as WorkspaceType;
  formData.resourceQuota = {
    cpu: 4,
    memory: 8,
    gpu: 1,
    storage: 100,
  };
  formData.config = {
    allowExternalAccess: false,
    enableAutoCleanup: true,
    cleanupDays: 30,
    maxMembers: 10,
    defaultResourceQuota: {
      cpu: 2,
      memory: 4,
      gpu: 0,
      storage: 50,
    },
  };
  
  if (formRef.value) {
    formRef.value.clearValidate();
  }
};

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validateFields();
    loading.value = true;

    await createWorkspace(formData);
    message.success('工作空间创建成功');
    emit('success');
    emit('update:visible', false);
  } catch (error: any) {
    if (error.errorFields) {
      message.error('请完善表单信息');
    } else {
      message.error('创建失败：' + (error.message || '未知错误'));
    }
  } finally {
    loading.value = false;
  }
};

// 取消
const handleCancel = () => {
  emit('update:visible', false);
};
</script>

<style scoped>
.ant-form-item {
  margin-bottom: 16px;
}
</style>