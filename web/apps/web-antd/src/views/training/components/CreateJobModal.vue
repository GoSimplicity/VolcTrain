<template>
  <Modal
    v-model:open="visible"
    title="创建训练任务"
    width="800px"
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
      <!-- 基本信息 -->
      <div class="form-section">
        <h4>基本信息</h4>
        
        <Form.Item label="任务名称" name="name">
          <Input
            v-model:value="formData.name"
            placeholder="请输入训练任务名称"
            :maxlength="50"
          />
        </Form.Item>

        <Form.Item label="描述" name="description">
          <Input.TextArea
            v-model:value="formData.description"
            :rows="2"
            placeholder="请输入任务描述"
            :maxlength="200"
            show-count
          />
        </Form.Item>

        <Form.Item label="工作空间" name="workspaceId">
          <Select
            v-model:value="formData.workspaceId"
            placeholder="选择工作空间"
            @change="handleWorkspaceChange"
          >
            <Select.Option
              v-for="workspace in workspaces"
              :key="workspace.id"
              :value="workspace.id"
            >
              {{ workspace.name }}
            </Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label="训练队列" name="queueId">
          <Select
            v-model:value="formData.queueId"
            placeholder="选择训练队列"
          >
            <Select.Option
              v-for="queue in availableQueues"
              :key="queue.id"
              :value="queue.id"
            >
              {{ queue.name }}
              <span style="color: #999; margin-left: 8px">
                ({{ queue.runningJobCount }}/{{ queue.maxRunningJobs }} 运行中)
              </span>
            </Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label="优先级" name="priority">
          <Select v-model:value="formData.priority">
            <Select.Option value="low">低</Select.Option>
            <Select.Option value="medium">中</Select.Option>
            <Select.Option value="high">高</Select.Option>
            <Select.Option value="urgent">紧急</Select.Option>
          </Select>
        </Form.Item>
      </div>

      <!-- 训练配置 -->
      <Divider />
      <div class="form-section">
        <h4>训练配置</h4>
        
        <Form.Item label="训练框架" name="framework">
          <Select v-model:value="formData.framework">
            <Select.Option value="pytorch">PyTorch</Select.Option>
            <Select.Option value="tensorflow">TensorFlow</Select.Option>
            <Select.Option value="keras">Keras</Select.Option>
            <Select.Option value="paddlepaddle">PaddlePaddle</Select.Option>
            <Select.Option value="mindspore">MindSpore</Select.Option>
            <Select.Option value="custom">自定义</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label="分布式类型" name="distributedType">
          <Select v-model:value="formData.distributedType">
            <Select.Option value="single">单机训练</Select.Option>
            <Select.Option value="data_parallel">数据并行</Select.Option>
            <Select.Option value="model_parallel">模型并行</Select.Option>
            <Select.Option value="pipeline_parallel">流水线并行</Select.Option>
            <Select.Option value="hybrid">混合并行</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label="容器镜像" name="image">
          <Input
            v-model:value="formData.image"
            placeholder="请输入容器镜像地址"
          />
        </Form.Item>

        <Form.Item label="启动命令" name="command">
          <Select
            v-model:value="formData.command"
            mode="tags"
            placeholder="请输入启动命令"
            style="width: 100%"
          >
            <Select.Option value="python">python</Select.Option>
            <Select.Option value="python3">python3</Select.Option>
            <Select.Option value="bash">bash</Select.Option>
            <Select.Option value="sh">sh</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label="启动参数" name="args">
          <Select
            v-model:value="formData.args"
            mode="tags"
            placeholder="请输入启动参数"
            style="width: 100%"
          />
        </Form.Item>
      </div>

      <!-- 资源配置 -->
      <Divider />
      <div class="form-section">
        <h4>资源配置</h4>
        
        <Form.Item label="副本数量" name="replicas">
          <InputNumber
            v-model:value="formData.replicas"
            :min="1"
            :max="16"
            style="width: 100%"
            addon-after="个"
          />
        </Form.Item>

        <Form.Item label="CPU需求" name="cpu">
          <InputNumber
            v-model:value="formData.resourceRequirements.cpu"
            :min="1"
            :max="32"
            style="width: 100%"
            addon-after="核"
          />
        </Form.Item>

        <Form.Item label="内存需求" name="memory">
          <InputNumber
            v-model:value="formData.resourceRequirements.memory"
            :min="1"
            :max="256"
            style="width: 100%"
            addon-after="GB"
          />
        </Form.Item>

        <Form.Item label="GPU需求" name="gpu">
          <InputNumber
            v-model:value="formData.resourceRequirements.gpu"
            :min="0"
            :max="8"
            style="width: 100%"
            addon-after="个"
          />
        </Form.Item>

        <Form.Item label="存储需求" name="storage">
          <InputNumber
            v-model:value="formData.resourceRequirements.storage"
            :min="1"
            :max="1000"
            style="width: 100%"
            addon-after="GB"
          />
        </Form.Item>
      </div>

      <!-- 数据配置 -->
      <Divider />
      <div class="form-section">
        <h4>数据配置</h4>
        
        <Form.Item label="输入数据路径" name="inputDataPath">
          <Input
            v-model:value="formData.inputDataPath"
            placeholder="请输入输入数据路径"
          />
        </Form.Item>

        <Form.Item label="输出数据路径" name="outputDataPath">
          <Input
            v-model:value="formData.outputDataPath"
            placeholder="请输入输出数据路径"
          />
        </Form.Item>
      </div>

      <!-- 高级配置 -->
      <Divider />
      <div class="form-section">
        <h4>高级配置</h4>
        
        <Form.Item label="最大重试次数" name="maxRetries">
          <InputNumber
            v-model:value="formData.maxRetries"
            :min="0"
            :max="10"
            style="width: 100%"
            addon-after="次"
          />
        </Form.Item>

        <Form.Item label="启用检查点">
          <Switch v-model:checked="formData.checkpointEnabled" />
        </Form.Item>

        <Form.Item 
          v-if="formData.checkpointEnabled"
          label="检查点间隔" 
          name="checkpointInterval"
        >
          <InputNumber
            v-model:value="formData.checkpointInterval"
            :min="1"
            :max="10000"
            style="width: 100%"
            addon-after="步"
          />
        </Form.Item>
      </div>
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
import type { 
  CreateTrainingJobRequest, 
  TrainingQueue,
  TrainingFramework,
  TrainingPriority,
  DistributedType 
} from '#/api/types';
import { createTrainingJob } from '#/api';

const props = defineProps<{
  visible: boolean;
  availableQueues: TrainingQueue[];
  defaultQueueId?: string;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
  success: [];
}>();

const formRef = ref();
const loading = ref(false);

// 模拟工作空间数据
const workspaces = ref([
  { id: 'ws1', name: '默认工作空间' },
  { id: 'ws2', name: '项目工作空间' },
]);

// 表单数据
const formData = reactive<CreateTrainingJobRequest>({
  name: '',
  description: '',
  workspaceId: 'ws1',
  queueId: '',
  framework: 'pytorch' as TrainingFramework,
  distributedType: 'single' as DistributedType,
  image: 'pytorch/pytorch:latest',
  command: ['python'],
  args: [],
  resourceRequirements: {
    cpu: 4,
    memory: 8,
    gpu: 1,
    storage: 50,
  },
  replicas: 1,
  inputDataPath: '',
  outputDataPath: '',
  maxRetries: 3,
  checkpointEnabled: true,
  checkpointInterval: 1000,
  priority: 'medium' as TrainingPriority,
});

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度应在2-50个字符之间', trigger: 'blur' },
  ],
  workspaceId: [
    { required: true, message: '请选择工作空间', trigger: 'change' },
  ],
  queueId: [
    { required: true, message: '请选择训练队列', trigger: 'change' },
  ],
  framework: [
    { required: true, message: '请选择训练框架', trigger: 'change' },
  ],
  distributedType: [
    { required: true, message: '请选择分布式类型', trigger: 'change' },
  ],
  image: [
    { required: true, message: '请输入容器镜像', trigger: 'blur' },
  ],
  command: [
    { required: true, message: '请输入启动命令', trigger: 'change' },
  ],
  cpu: [
    { required: true, message: '请输入CPU需求', trigger: 'blur' },
    { type: 'number', min: 1, max: 32, message: 'CPU需求应在1-32核之间', trigger: 'blur' },
  ],
  memory: [
    { required: true, message: '请输入内存需求', trigger: 'blur' },
    { type: 'number', min: 1, max: 256, message: '内存需求应在1-256GB之间', trigger: 'blur' },
  ],
  replicas: [
    { required: true, message: '请输入副本数量', trigger: 'blur' },
    { type: 'number', min: 1, max: 16, message: '副本数量应在1-16个之间', trigger: 'blur' },
  ],
  maxRetries: [
    { required: true, message: '请输入最大重试次数', trigger: 'blur' },
    { type: 'number', min: 0, max: 10, message: '重试次数应在0-10次之间', trigger: 'blur' },
  ],
  checkpointInterval: [
    { type: 'number', min: 1, max: 10000, message: '检查点间隔应在1-10000步之间', trigger: 'blur' },
  ],
};

// 监听visible变化，重置表单
watch(() => props.visible, (newVal) => {
  if (newVal) {
    resetForm();
    if (props.defaultQueueId) {
      formData.queueId = props.defaultQueueId;
    }
  }
});

// 重置表单
const resetForm = () => {
  formData.name = '';
  formData.description = '';
  formData.workspaceId = 'ws1';
  formData.queueId = props.defaultQueueId || '';
  formData.framework = 'pytorch' as TrainingFramework;
  formData.distributedType = 'single' as DistributedType;
  formData.image = 'pytorch/pytorch:latest';
  formData.command = ['python'];
  formData.args = [];
  formData.resourceRequirements = {
    cpu: 4,
    memory: 8,
    gpu: 1,
    storage: 50,
  };
  formData.replicas = 1;
  formData.inputDataPath = '';
  formData.outputDataPath = '';
  formData.maxRetries = 3;
  formData.checkpointEnabled = true;
  formData.checkpointInterval = 1000;
  formData.priority = 'medium' as TrainingPriority;
  
  if (formRef.value) {
    formRef.value.clearValidate();
  }
};

// 事件处理
const handleWorkspaceChange = (workspaceId: string) => {
  console.log('工作空间变更:', workspaceId);
  // 可以根据工作空间加载对应的队列
};

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validateFields();
    loading.value = true;

    await createTrainingJob(formData);
    message.success('训练任务创建成功');
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
.form-section {
  h4 {
    margin: 0 0 16px 0;
    color: #333;
    font-weight: 500;
  }
}

.ant-form-item {
  margin-bottom: 16px;
}
</style>