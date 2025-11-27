<template>
  <Form
    ref="formRef"
    :model="formData"
    :rules="formRules"
    layout="vertical"
    class="create-experiment-form"
  >
    <Form.Item label="实验名称" name="name">
      <Input 
        v-model:value="formData.name"
        placeholder="请输入实验名称"
      />
    </Form.Item>

    <Form.Item label="实验描述" name="description">
      <Input.TextArea 
        v-model:value="formData.description"
        :rows="3"
        placeholder="请输入实验描述"
      />
    </Form.Item>

    <Form.Item label="实验类型" name="type">
      <Select v-model:value="formData.type" placeholder="选择实验类型">
        <Select.Option value="image_classification">图像分类</Select.Option>
        <Select.Option value="object_detection">目标检测</Select.Option>
        <Select.Option value="nlp">自然语言处理</Select.Option>
        <Select.Option value="recommendation">推荐系统</Select.Option>
        <Select.Option value="time_series">时序预测</Select.Option>
        <Select.Option value="other">其他</Select.Option>
      </Select>
    </Form.Item>

    <Form.Item label="标签" name="tags">
      <Select 
        v-model:value="formData.tags"
        mode="tags"
        placeholder="输入标签"
        :token-separators="[',']"
      >
        <Select.Option value="深度学习">深度学习</Select.Option>
        <Select.Option value="计算机视觉">计算机视觉</Select.Option>
        <Select.Option value="自然语言处理">自然语言处理</Select.Option>
        <Select.Option value="强化学习">强化学习</Select.Option>
      </Select>
    </Form.Item>

    <Form.Item label="目标指标" name="targetMetric">
      <Input 
        v-model:value="formData.targetMetric"
        placeholder="如: accuracy, f1_score, loss"
      />
    </Form.Item>

    <Form.Item label="是否公开" name="isPublic">
      <Switch v-model:checked="formData.isPublic" />
      <span style="margin-left: 8px;">公开实验可被其他用户查看</span>
    </Form.Item>
  </Form>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import {
  Form,
  Input,
  Select,
  Switch,
} from 'ant-design-vue';
import type { Rule } from 'ant-design-vue/es/form';

// 表单引用
const formRef = ref();

// 表单数据
const formData = ref({
  name: '',
  description: '',
  type: undefined,
  tags: [],
  targetMetric: '',
  isPublic: false,
});

// 表单验证规则
const formRules: Record<string, Rule[]> = {
  name: [
    { required: true, message: '请输入实验名称' },
    { min: 3, max: 50, message: '实验名称长度应在 3-50 个字符' },
  ],
  type: [
    { required: true, message: '请选择实验类型' },
  ],
  targetMetric: [
    { required: true, message: '请输入目标指标' },
  ],
};

// 验证表单
const validate = async () => {
  const values = await formRef.value.validate();
  return values;
};

// 重置表单
const resetFields = () => {
  formRef.value?.resetFields();
};

// 暴露方法给父组件
defineExpose({
  validate,
  resetFields,
});
</script>

<style lang="scss" scoped>
.create-experiment-form {
  padding: 16px 0;
}
</style>