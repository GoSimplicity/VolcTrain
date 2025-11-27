<template>
  <Form
    ref="formRef"
    :model="formData"
    :rules="formRules"
    layout="vertical"
    class="create-queue-form"
  >
    <Form.Item label="队列名称" name="name">
      <Input 
        v-model:value="formData.name"
        placeholder="请输入队列名称"
      />
    </Form.Item>

    <Form.Item label="队列描述" name="description">
      <Input.TextArea 
        v-model:value="formData.description"
        :rows="3"
        placeholder="请输入队列描述"
      />
    </Form.Item>

    <Row :gutter="16">
      <Col :span="12">
        <Form.Item label="优先级" name="priority">
          <InputNumber 
            v-model:value="formData.priority"
            :min="1"
            :max="10"
            style="width: 100%"
            placeholder="1-10"
          />
        </Form.Item>
      </Col>
      <Col :span="12">
        <Form.Item label="最大运行任务数" name="maxRunningJobs">
          <InputNumber 
            v-model:value="formData.maxRunningJobs"
            :min="1"
            :max="1000"
            style="width: 100%"
            placeholder="最大并发任务数"
          />
        </Form.Item>
      </Col>
    </Row>

    <!-- 资源配额 -->
    <div class="form-section">
      <h3 class="section-title">资源配额</h3>
      
      <Row :gutter="16">
        <Col :span="8">
          <Form.Item label="CPU" name="cpu">
            <Input 
              v-model:value="formData.resourceQuota.cpu"
              placeholder="如: 100"
              suffix="核"
            />
          </Form.Item>
        </Col>
        <Col :span="8">
          <Form.Item label="内存" name="memory">
            <Input 
              v-model:value="formData.resourceQuota.memory"
              placeholder="如: 200Gi"
            />
          </Form.Item>
        </Col>
        <Col :span="8">
          <Form.Item label="GPU数量" name="gpu">
            <InputNumber 
              v-model:value="formData.resourceQuota.gpu"
              :min="0"
              :max="1000"
              style="width: 100%"
              placeholder="GPU数量"
            />
          </Form.Item>
        </Col>
      </Row>
    </div>

    <!-- 调度策略 -->
    <div class="form-section">
      <h3 class="section-title">调度策略</h3>
      
      <div class="policies-list">
        <div 
          v-for="(policy, index) in formData.policies" 
          :key="index" 
          class="policy-item"
        >
          <Row :gutter="16">
            <Col :span="8">
              <Select 
                v-model:value="policy.type"
                placeholder="策略类型"
                style="width: 100%"
              >
                <Select.Option value="fairshare">公平共享</Select.Option>
                <Select.Option value="priority">优先级</Select.Option>
                <Select.Option value="drf">主导资源公平</Select.Option>
                <Select.Option value="gang">Gang调度</Select.Option>
              </Select>
            </Col>
            <Col :span="12">
              <Input 
                v-model:value="policy.config.description"
                placeholder="策略配置描述"
              />
            </Col>
            <Col :span="4">
              <Button 
                type="text" 
                danger
                @click="removePolicy(index)"
              >
                <TrashIcon />
              </Button>
            </Col>
          </Row>
        </div>
        
        <Button type="dashed" @click="addPolicy" block>
          <PlusIcon />
          添加调度策略
        </Button>
      </div>
    </div>
  </Form>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import {
  Form,
  Input,
  InputNumber,
  Select,
  Button,
  Row,
  Col,
} from 'ant-design-vue';
import type { Rule } from 'ant-design-vue/es/form';

// 图标
import { Trash as TrashIcon, Plus as PlusIcon } from 'lucide-vue-next';

// 类型定义
import type { TrainingApi } from '#/types/api';

// 表单引用
const formRef = ref();

// 表单数据
const formData = ref<Partial<TrainingApi.TrainingQueue>>({
  name: '',
  description: '',
  priority: 1,
  maxRunningJobs: 10,
  resourceQuota: {
    cpu: '50',
    memory: '100Gi',
    gpu: 10,
  },
  policies: [],
});

// 表单验证规则
const formRules: Record<string, Rule[]> = {
  name: [
    { required: true, message: '请输入队列名称' },
    { min: 3, max: 30, message: '队列名称长度应在 3-30 个字符' },
    { pattern: /^[a-zA-Z0-9-_]+$/, message: '队列名称只能包含字母、数字、横线和下划线' },
  ],
  priority: [
    { required: true, message: '请设置队列优先级' },
    { type: 'number', min: 1, max: 10, message: '优先级应在 1-10 之间' },
  ],
  maxRunningJobs: [
    { required: true, message: '请设置最大运行任务数' },
    { type: 'number', min: 1, message: '最大运行任务数至少为 1' },
  ],
  cpu: [
    { required: true, message: '请设置CPU配额' },
  ],
  memory: [
    { required: true, message: '请设置内存配额' },
  ],
  gpu: [
    { required: true, message: '请设置GPU配额' },
    { type: 'number', min: 0, message: 'GPU数量不能为负数' },
  ],
};

// 添加调度策略
const addPolicy = () => {
  formData.value.policies = formData.value.policies || [];
  formData.value.policies.push({
    type: 'fairshare',
    config: {
      description: '',
    },
  });
};

// 删除调度策略
const removePolicy = (index: number) => {
  if (formData.value.policies) {
    formData.value.policies.splice(index, 1);
  }
};

// 验证表单
const validate = async () => {
  const values = await formRef.value.validate();
  return values;
};

// 重置表单
const resetFields = () => {
  formRef.value?.resetFields();
  formData.value.policies = [];
};

// 暴露方法给父组件
defineExpose({
  validate,
  resetFields,
});
</script>

<style lang="scss" scoped>
.create-queue-form {
  .form-section {
    margin-bottom: 24px;
    padding: 16px;
    background: #fafafa;
    border-radius: 8px;

    .section-title {
      font-size: 14px;
      font-weight: 600;
      margin-bottom: 16px;
      color: #333;
      border-bottom: 1px solid #e8e8e8;
      padding-bottom: 8px;
    }
  }

  .policies-list {
    .policy-item {
      margin-bottom: 12px;
      padding: 12px;
      background: white;
      border-radius: 6px;
      border: 1px solid #e8e8e8;
    }
  }
}
</style>