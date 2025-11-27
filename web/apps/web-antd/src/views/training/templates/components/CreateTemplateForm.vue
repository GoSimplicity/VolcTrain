<template>
  <Form ref="formRef" :model="formData" :rules="formRules" layout="vertical">
    <Form.Item label="模板名称" name="name">
      <Input v-model:value="formData.name" placeholder="请输入模板名称" />
    </Form.Item>
    
    <Form.Item label="模板描述" name="description">
      <Input.TextArea v-model:value="formData.description" :rows="3" placeholder="请输入模板描述" />
    </Form.Item>
    
    <Row :gutter="16">
      <Col :span="12">
        <Form.Item label="训练框架" name="framework">
          <Select v-model:value="formData.framework" placeholder="选择训练框架">
            <Select.Option value="tensorflow">TensorFlow</Select.Option>
            <Select.Option value="pytorch">PyTorch</Select.Option>
          </Select>
        </Form.Item>
      </Col>
      <Col :span="12">
        <Form.Item label="任务类型" name="type">
          <Select v-model:value="formData.type" placeholder="选择任务类型">
            <Select.Option value="classification">图像分类</Select.Option>
            <Select.Option value="detection">目标检测</Select.Option>
          </Select>
        </Form.Item>
      </Col>
    </Row>
  </Form>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { Form, Input, Select, Row, Col } from 'ant-design-vue';

const formRef = ref();
const formData = ref({ name: '', description: '', framework: undefined, type: undefined });
const formRules = {
  name: [{ required: true, message: '请输入模板名称' }],
  framework: [{ required: true, message: '请选择训练框架' }],
};

const validate = () => formRef.value.validate();
const resetFields = () => formRef.value?.resetFields();

defineExpose({ validate, resetFields });
</script>