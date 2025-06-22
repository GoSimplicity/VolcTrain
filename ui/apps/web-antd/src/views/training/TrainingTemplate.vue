<script lang="ts" setup>
import { ref, h } from 'vue';
import {
  Table,
  Card,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
} from 'ant-design-vue';

defineOptions({ name: 'TrainingTemplate' });

const dataSource = ref([
  {
    id: '1',
    name: 'BERT预训练模板',
    description: 'BERT预训练基础模板',
    framework: 'PyTorch',
    gpuType: 'NVIDIA A100',
    gpuCount: 4,
    parameters: 'batch_size=32, learning_rate=5e-5, epochs=3',
    createTime: '2023-10-15 10:30:00',
    creator: '张三',
  },
  {
    id: '2',
    name: '图像分类训练模板',
    description: '用于图像分类的训练模板',
    framework: 'TensorFlow',
    gpuType: 'NVIDIA V100',
    gpuCount: 2,
    parameters: 'batch_size=64, learning_rate=0.001, epochs=10',
    createTime: '2023-10-18 14:20:00',
    creator: '李四',
  },
  {
    id: '3',
    name: 'NLP序列标注模板',
    description: '用于NER等序列标注任务的训练模板',
    framework: 'PyTorch',
    gpuType: 'NVIDIA T4',
    gpuCount: 1,
    parameters: 'batch_size=16, learning_rate=3e-5, epochs=5',
    createTime: '2023-10-20 09:00:00',
    creator: '王五',
  },
]);

const columns = [
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
    title: '描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '框架',
    dataIndex: 'framework',
    key: 'framework',
  },
  {
    title: 'GPU类型',
    dataIndex: 'gpuType',
    key: 'gpuType',
  },
  {
    title: 'GPU数量',
    dataIndex: 'gpuCount',
    key: 'gpuCount',
  },
  {
    title: '训练参数',
    dataIndex: 'parameters',
    key: 'parameters',
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    key: 'createTime',
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
  },
  {
    title: '操作',
    key: 'action',
    customRender: () => {
      return h(Space, { size: 'middle' }, [
        h('a', {}, '编辑'),
        h('a', {}, '复制'),
        h('a', {}, '删除')
      ]);
    },
  },
];

const visible = ref(false);
const formRef = ref();
const formState = ref({
  name: '',
  description: '',
  framework: 'PyTorch',
  gpuType: 'NVIDIA A100',
  gpuCount: 1,
  parameters: '',
});

const showModal = () => {
  visible.value = true;
};

const handleCancel = () => {
  visible.value = false;
};

const handleOk = () => {
  formRef.value
    .validate()
    .then(() => {
      dataSource.value.push({
        id: String(dataSource.value.length + 1),
        name: formState.value.name,
        description: formState.value.description,
        framework: formState.value.framework,
        gpuType: formState.value.gpuType,
        gpuCount: formState.value.gpuCount,
        parameters: formState.value.parameters,
        createTime: new Date().toLocaleString(),
        creator: '当前用户',
      });
      visible.value = false;
      formState.value = {
        name: '',
        description: '',
        framework: 'PyTorch',
        gpuType: 'NVIDIA A100',
        gpuCount: 1,
        parameters: '',
      };
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};
</script>

<template>
  <div class="training-template-container">
    <Card title="训练模板管理" :bordered="false">
      <template #extra>
        <Button type="primary" @click="showModal">新增模板</Button>
      </template>
      <Table
        :columns="columns"
        :dataSource="dataSource"
        rowKey="id"
        :pagination="{ pageSize: 10 }"
      />
    </Card>

    <Modal
      v-model:visible="visible"
      title="新增训练模板"
      @ok="handleOk"
      @cancel="handleCancel"
      :maskClosable="false"
      width="700px"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="name"
          label="模板名称"
          :rules="[{ required: true, message: '请输入模板名称' }]"
        >
          <Input v-model:value="formState.name" placeholder="请输入模板名称" />
        </Form.Item>
        <Form.Item name="description" label="描述">
          <Input.TextArea
            v-model:value="formState.description"
            placeholder="请输入模板描述"
          />
        </Form.Item>
        <Form.Item
          name="framework"
          label="框架"
          :rules="[{ required: true, message: '请选择框架' }]"
        >
          <Select v-model:value="formState.framework" placeholder="请选择框架">
            <Select.Option value="PyTorch">PyTorch</Select.Option>
            <Select.Option value="TensorFlow">TensorFlow</Select.Option>
            <Select.Option value="JAX">JAX</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="gpuType"
          label="GPU类型"
          :rules="[{ required: true, message: '请选择GPU类型' }]"
        >
          <Select v-model:value="formState.gpuType" placeholder="请选择GPU类型">
            <Select.Option value="NVIDIA A100">NVIDIA A100</Select.Option>
            <Select.Option value="NVIDIA V100">NVIDIA V100</Select.Option>
            <Select.Option value="NVIDIA T4">NVIDIA T4</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="gpuCount"
          label="GPU数量"
          :rules="[{ required: true, message: '请输入GPU数量' }]"
        >
          <InputNumber v-model:value="formState.gpuCount" :min="1" :max="16" />
        </Form.Item>
        <Form.Item
          name="parameters"
          label="训练参数"
          :rules="[{ required: true, message: '请输入训练参数' }]"
        >
          <Input.TextArea
            v-model:value="formState.parameters"
            placeholder="例如: batch_size=32, learning_rate=5e-5, epochs=3"
          />
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.training-template-container {
  padding: 0;
}
</style>
