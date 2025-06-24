<script lang="ts" setup>
import { ref, h } from 'vue';
import {
  Table,
  Card,
  Button,
  Space,
  Tag,
  Modal,
  Form,
  Input,
  Select,
  Radio,
} from 'ant-design-vue';

defineOptions({ name: 'SystemApi' });

const dataSource = ref([
  {
    id: '1',
    name: '获取用户列表',
    path: '/api/users',
    method: 'GET',
    category: '用户管理',
    description: '分页获取用户列表数据',
    status: '启用',
    createTime: '2023-10-01 08:00:00',
    updateTime: '2023-11-15 09:30:00',
  },
  {
    id: '2',
    name: '创建用户',
    path: '/api/users',
    method: 'POST',
    category: '用户管理',
    description: '创建新用户',
    status: '启用',
    createTime: '2023-10-01 08:30:00',
    updateTime: '2023-11-15 10:00:00',
  },
  {
    id: '3',
    name: '获取角色列表',
    path: '/api/roles',
    method: 'GET',
    category: '角色管理',
    description: '分页获取角色列表数据',
    status: '启用',
    createTime: '2023-10-02 09:00:00',
    updateTime: '2023-11-14 11:20:00',
  },
  {
    id: '4',
    name: '创建角色',
    path: '/api/roles',
    method: 'POST',
    category: '角色管理',
    description: '创建新角色',
    status: '禁用',
    createTime: '2023-10-02 09:30:00',
    updateTime: '2023-11-10 15:40:00',
  },
]);

const columns = [
  {
    title: 'API ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: 'API名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '请求路径',
    dataIndex: 'path',
    key: 'path',
  },
  {
    title: '请求方法',
    dataIndex: 'method',
    key: 'method',
    customRender: ({ text }) => {
      const colorMap = {
        GET: 'green',
        POST: 'blue',
        PUT: 'orange',
        DELETE: 'red',
      };
      return h(Tag, { color: colorMap[text] || 'default' }, () => text);
    },
  },
  {
    title: '所属分类',
    dataIndex: 'category',
    key: 'category',
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      const color = text === '启用' ? 'green' : 'red';
      return h(Tag, { color }, () => text);
    },
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    key: 'createTime',
  },
  {
    title: '更新时间',
    dataIndex: 'updateTime',
    key: 'updateTime',
  },
  {
    title: '操作',
    key: 'action',
    customRender: () => {
      return h(Space, { size: 'middle' }, [
        h('a', {}, '编辑'),
        h('a', {}, '测试'),
        h('a', {}, '删除')
      ]);
    },
  },
];

const visible = ref(false);
const formRef = ref();
const formState = ref({
  name: '',
  path: '',
  method: 'GET',
  category: '',
  description: '',
  status: 1,
});

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
];

const categoryOptions = [
  { label: '用户管理', value: '用户管理' },
  { label: '角色管理', value: '角色管理' },
  { label: '权限管理', value: '权限管理' },
  { label: '系统配置', value: '系统配置' },
];

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
        path: formState.value.path,
        method: formState.value.method,
        category: formState.value.category,
        description: formState.value.description,
        status: formState.value.status === 1 ? '启用' : '禁用',
        createTime: new Date().toLocaleString(),
        updateTime: new Date().toLocaleString(),
      });
      visible.value = false;
      formState.value = {
        name: '',
        path: '',
        method: 'GET',
        category: '',
        description: '',
        status: 1,
      };
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};
</script>

<template>
  <div class="system-api-container">
    <Card title="API管理" :bordered="false">
      <template #extra>
        <Button type="primary" @click="showModal">新增API</Button>
      </template>
      <Table :columns="columns" :dataSource="dataSource" rowKey="id" />
    </Card>

    <Modal
      v-model:visible="visible"
      title="新增API"
      @ok="handleOk"
      @cancel="handleCancel"
      :maskClosable="false"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="name"
          label="API名称"
          :rules="[{ required: true, message: '请输入API名称' }]"
        >
          <Input v-model:value="formState.name" placeholder="请输入API名称" />
        </Form.Item>
        <Form.Item
          name="path"
          label="请求路径"
          :rules="[{ required: true, message: '请输入请求路径' }]"
        >
          <Input
            v-model:value="formState.path"
            placeholder="请输入请求路径，如/api/users"
          />
        </Form.Item>
        <Form.Item
          name="method"
          label="请求方法"
          :rules="[{ required: true, message: '请选择请求方法' }]"
        >
          <Select v-model:value="formState.method" placeholder="请选择请求方法">
            <Select.Option
              v-for="option in methodOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="category"
          label="所属分类"
          :rules="[{ required: true, message: '请选择所属分类' }]"
        >
          <Select
            v-model:value="formState.category"
            placeholder="请选择所属分类"
          >
            <Select.Option
              v-for="option in categoryOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item name="description" label="API描述">
          <Input.TextArea
            v-model:value="formState.description"
            placeholder="请输入API描述"
          />
        </Form.Item>
        <Form.Item name="status" label="状态">
          <Radio.Group v-model:value="formState.status">
            <Radio :value="1">启用</Radio>
            <Radio :value="0">禁用</Radio>
          </Radio.Group>
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.system-api-container {
  padding: 0;
}
</style>
