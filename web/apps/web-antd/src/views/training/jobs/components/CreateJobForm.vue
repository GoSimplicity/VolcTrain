<template>
  <Form
    ref="formRef"
    :model="formData"
    :rules="formRules"
    layout="vertical"
    class="create-job-form"
  >
    <Row :gutter="16">
      <!-- 基本信息 -->
      <Col :span="24">
        <div class="form-section">
          <h3 class="section-title">基本信息</h3>
          
          <Row :gutter="16">
            <Col :span="12">
              <Form.Item label="任务名称" name="name">
                <Input 
                  v-model:value="formData.name"
                  placeholder="请输入任务名称"
                />
              </Form.Item>
            </Col>
            <Col :span="12">
              <Form.Item label="优先级" name="priority">
                <Select v-model:value="formData.priority" placeholder="选择优先级">
                  <Select.Option value="low">低</Select.Option>
                  <Select.Option value="normal">正常</Select.Option>
                  <Select.Option value="high">高</Select.Option>
                  <Select.Option value="urgent">紧急</Select.Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Form.Item label="任务描述" name="description">
            <Input.TextArea 
              v-model:value="formData.description"
              :rows="3"
              placeholder="请输入任务描述"
            />
          </Form.Item>

          <Row :gutter="16">
            <Col :span="12">
              <Form.Item label="训练框架" name="framework">
                <Select v-model:value="formData.framework" placeholder="选择训练框架">
                  <Select.Option value="tensorflow">TensorFlow</Select.Option>
                  <Select.Option value="pytorch">PyTorch</Select.Option>
                  <Select.Option value="paddlepaddle">PaddlePaddle</Select.Option>
                  <Select.Option value="mindspore">MindSpore</Select.Option>
                  <Select.Option value="mpi">MPI</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col :span="12">
              <Form.Item label="训练队列" name="queueId">
                <Select v-model:value="formData.queueId" placeholder="选择训练队列">
                  <Select.Option 
                    v-for="queue in queues"
                    :key="queue.id"
                    :value="queue.id"
                  >
                    {{ queue.name }}
                  </Select.Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
        </div>
      </Col>

      <!-- 资源配置 -->
      <Col :span="24">
        <div class="form-section">
          <h3 class="section-title">资源配置</h3>
          
          <Row :gutter="16">
            <Col :span="8">
              <Form.Item label="CPU" name="cpu">
                <Input 
                  v-model:value="formData.resourceConfig.cpu"
                  placeholder="如: 4"
                  suffix="核"
                />
              </Form.Item>
            </Col>
            <Col :span="8">
              <Form.Item label="内存" name="memory">
                <Input 
                  v-model:value="formData.resourceConfig.memory"
                  placeholder="如: 8Gi"
                />
              </Form.Item>
            </Col>
            <Col :span="8">
              <Form.Item label="副本数" name="replicas">
                <InputNumber 
                  v-model:value="formData.resourceConfig.replicas"
                  :min="1"
                  :max="100"
                  style="width: 100%"
                />
              </Form.Item>
            </Col>
          </Row>

          <Row :gutter="16">
            <Col :span="12">
              <Form.Item label="GPU 类型" name="gpuType">
                <Select v-model:value="formData.resourceConfig.gpu.type" placeholder="选择GPU类型">
                  <Select.Option value="V100">NVIDIA V100</Select.Option>
                  <Select.Option value="A100">NVIDIA A100</Select.Option>
                  <Select.Option value="T4">NVIDIA T4</Select.Option>
                  <Select.Option value="RTX3080">NVIDIA RTX 3080</Select.Option>
                  <Select.Option value="RTX3090">NVIDIA RTX 3090</Select.Option>
                  <Select.Option value="H100">NVIDIA H100</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col :span="12">
              <Form.Item label="GPU 数量" name="gpuCount">
                <InputNumber 
                  v-model:value="formData.resourceConfig.gpu.count"
                  :min="0"
                  :max="8"
                  style="width: 100%"
                />
              </Form.Item>
            </Col>
          </Row>
        </div>
      </Col>

      <!-- 运行配置 -->
      <Col :span="24">
        <div class="form-section">
          <h3 class="section-title">运行配置</h3>
          
          <Form.Item label="镜像地址" name="image">
            <Input 
              v-model:value="formData.image"
              placeholder="如: tensorflow/tensorflow:2.11.0-gpu"
            />
          </Form.Item>

          <Form.Item label="启动命令" name="command">
            <Input.TextArea 
              v-model:value="commandText"
              :rows="3"
              placeholder="如: python train.py --epochs 100 --batch-size 32"
            />
          </Form.Item>

          <Form.Item label="工作目录" name="workingDir">
            <Input 
              v-model:value="formData.workingDir"
              placeholder="如: /workspace"
            />
          </Form.Item>

          <Form.Item label="环境变量">
            <div class="env-vars">
              <div 
                v-for="(env, index) in formData.env" 
                :key="index" 
                class="env-var-item"
              >
                <Input 
                  v-model:value="env.key"
                  placeholder="变量名"
                  style="width: 40%; margin-right: 8px;"
                />
                <Input 
                  v-model:value="env.value"
                  placeholder="变量值"
                  style="width: 40%; margin-right: 8px;"
                />
                <Button 
                  type="text" 
                  danger
                  @click="removeEnvVar(index)"
                >
                  <TrashIcon />
                </Button>
              </div>
              <Button type="dashed" @click="addEnvVar" block>
                <PlusIcon />
                添加环境变量
              </Button>
            </div>
          </Form.Item>
        </div>
      </Col>

      <!-- 数据配置 -->
      <Col :span="24">
        <div class="form-section">
          <h3 class="section-title">数据配置</h3>
          
          <Form.Item label="数据集">
            <Select 
              v-model:value="formData.datasetIds"
              mode="multiple"
              placeholder="选择数据集"
            >
              <Select.Option 
                v-for="dataset in availableDatasets"
                :key="dataset.id"
                :value="dataset.id"
              >
                {{ dataset.name }}
              </Select.Option>
            </Select>
          </Form.Item>

          <Form.Item label="输出路径" name="outputPath">
            <Input 
              v-model:value="formData.outputPath"
              placeholder="如: /data/output"
            />
          </Form.Item>
        </div>
      </Col>

      <!-- Volcano 高级配置 -->
      <Col :span="24">
        <div class="form-section">
          <h3 class="section-title">
            高级配置
            <Switch 
              v-model:checked="showAdvanced" 
              size="small"
              style="margin-left: 12px;"
            />
          </h3>
          
          <div v-if="showAdvanced">
            <Row :gutter="16">
              <Col :span="12">
                <Form.Item label="最小可用实例" name="minAvailable">
                  <InputNumber 
                    v-model:value="formData.volcanoSpec.minAvailable"
                    :min="1"
                    style="width: 100%"
                  />
                </Form.Item>
              </Col>
              <Col :span="12">
                <Form.Item label="最大重试次数" name="maxRetry">
                  <InputNumber 
                    v-model:value="formData.volcanoSpec.maxRetry"
                    :min="0"
                    :max="10"
                    style="width: 100%"
                  />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item label="调度器名称" name="schedulerName">
              <Input 
                v-model:value="formData.volcanoSpec.schedulerName"
                placeholder="volcano"
              />
            </Form.Item>

            <Form.Item label="优先级类名" name="priorityClassName">
              <Input 
                v-model:value="formData.volcanoSpec.priorityClassName"
                placeholder="high-priority"
              />
            </Form.Item>
          </div>
        </div>
      </Col>
    </Row>
  </Form>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue';
import {
  Form,
  Input,
  Select,
  InputNumber,
  Switch,
  Button,
  Row,
  Col,
} from 'ant-design-vue';
import type { Rule } from 'ant-design-vue/es/form';

// 图标
import { Trash as TrashIcon, Plus as PlusIcon } from 'lucide-vue-next';

// 类型定义
import type { TrainingApi } from '#/types/api';

interface Props {
  queues: TrainingApi.TrainingQueue[];
}

const props = defineProps<Props>();

// 表单引用
const formRef = ref();

// 显示高级配置
const showAdvanced = ref(false);

// 可用数据集（模拟数据）
const availableDatasets = ref([
  { id: '1', name: 'ImageNet-2012' },
  { id: '2', name: 'COCO-2017' },
  { id: '3', name: 'CIFAR-10' },
]);

// 表单数据
const formData = ref<Partial<TrainingApi.CreateTrainingJobRequest>>({
  name: '',
  description: '',
  priority: 'normal',
  framework: undefined,
  queueId: '',
  resourceConfig: {
    cpu: '4',
    memory: '8Gi',
    replicas: 1,
    gpu: {
      type: 'V100',
      count: 1,
    },
  },
  image: '',
  command: [],
  workingDir: '/workspace',
  env: {},
  datasetIds: [],
  outputPath: '/data/output',
  volcanoSpec: {
    minAvailable: 1,
    maxRetry: 3,
    schedulerName: 'volcano',
    priorityClassName: '',
  },
});

// 环境变量数组形式
const envVars = ref<Array<{ key: string; value: string }>>([]);

// 命令文本
const commandText = ref('');

// 表单验证规则
const formRules: Record<string, Rule[]> = {
  name: [
    { required: true, message: '请输入任务名称' },
    { min: 3, max: 50, message: '任务名称长度应在 3-50 个字符' },
  ],
  framework: [
    { required: true, message: '请选择训练框架' },
  ],
  queueId: [
    { required: true, message: '请选择训练队列' },
  ],
  image: [
    { required: true, message: '请输入镜像地址' },
  ],
};

// 监听环境变量变化
watch(envVars, (newVars) => {
  const envObj: Record<string, string> = {};
  newVars.forEach(({ key, value }) => {
    if (key && value) {
      envObj[key] = value;
    }
  });
  if (formData.value.env) {
    formData.value.env = envObj;
  }
}, { deep: true });

// 监听命令文本变化
watch(commandText, (newCommand) => {
  if (formData.value.command) {
    formData.value.command = newCommand ? newCommand.split(/\s+/) : [];
  }
});

// 添加环境变量
const addEnvVar = () => {
  envVars.value.push({ key: '', value: '' });
};

// 删除环境变量
const removeEnvVar = (index: number) => {
  envVars.value.splice(index, 1);
};

// 验证表单
const validate = async () => {
  const values = await formRef.value.validate();
  
  // 处理环境变量
  const env: Record<string, string> = {};
  envVars.value.forEach(({ key, value }) => {
    if (key && value) {
      env[key] = value;
    }
  });

  return {
    ...values,
    env,
    command: commandText.value ? commandText.value.split(/\s+/) : [],
  };
};

// 重置表单
const resetFields = () => {
  formRef.value?.resetFields();
  envVars.value = [];
  commandText.value = '';
  showAdvanced.value = false;
};

// 暴露方法给父组件
defineExpose({
  validate,
  resetFields,
});
</script>

<style lang="scss" scoped>
.create-job-form {
  .form-section {
    margin-bottom: 32px;
    padding: 20px;
    background: #fafafa;
    border-radius: 8px;

    .section-title {
      font-size: 16px;
      font-weight: 600;
      margin-bottom: 16px;
      color: #333;
      border-bottom: 1px solid #e8e8e8;
      padding-bottom: 8px;
      display: flex;
      align-items: center;
    }
  }

  .env-vars {
    .env-var-item {
      display: flex;
      align-items: center;
      margin-bottom: 8px;
    }
  }
}
</style>