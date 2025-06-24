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
} from 'ant-design-vue';

defineOptions({ name: 'SystemUser' });

const dataSource = ref([
  {
    id: '1',
    username: 'admin',
    name: '管理员',
    email: 'admin@example.com',
    phone: '13800138000',
    role: '超级管理员',
    status: '启用',
    createTime: '2023-10-01 08:00:00',
    lastLoginTime: '2023-11-15 09:30:00',
  },
  {
    id: '2',
    username: 'user01',
    name: '张三',
    email: 'zhangsan@example.com',
    phone: '13800138001',
    role: '普通用户',
    status: '启用',
    createTime: '2023-10-05 10:15:00',
    lastLoginTime: '2023-11-14 16:20:00',
  },
  {
    id: '3',
    username: 'user02',
    name: '李四',
    email: 'lisi@example.com',
    phone: '13800138002',
    role: '开发者',
    status: '禁用',
    createTime: '2023-10-10 14:30:00',
    lastLoginTime: '2023-11-01 11:45:00',
  },
]);

const columns = [
  {
    title: '用户ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '用户名',
    dataIndex: 'username',
    key: 'username',
  },
  {
    title: '姓名',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '邮箱',
    dataIndex: 'email',
    key: 'email',
  },
  {
    title: '手机号',
    dataIndex: 'phone',
    key: 'phone',
  },
  {
    title: '角色',
    dataIndex: 'role',
    key: 'role',
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
    title: '最后登录时间',
    dataIndex: 'lastLoginTime',
    key: 'lastLoginTime',
  },
  {
    title: '操作',
    key: 'action',
    customRender: () => {
      return h(Space, { size: 'middle' }, [
        h('a', {}, '编辑'),
        h('a', {}, '重置密码'),
        h('a', {}, '禁用')
      ]);
    },
  },
];

const visible = ref(false);
const formRef = ref();
const formState = ref({
  username: '',
  name: '',
  email: '',
  phone: '',
  role: '普通用户',
  password: '',
  confirmPassword: '',
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
        username: formState.value.username,
        name: formState.value.name,
        email: formState.value.email,
        phone: formState.value.phone,
        role: formState.value.role,
        status: '启用',
        createTime: new Date().toLocaleString(),
        lastLoginTime: '-',
      });
      visible.value = false;
      formState.value = {
        username: '',
        name: '',
        email: '',
        phone: '',
        role: '普通用户',
        password: '',
        confirmPassword: '',
      };
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};
</script>

<template>
  <div class="system-user-container">
    <Card title="用户管理" :bordered="false">
      <template #extra>
        <Button type="primary" @click="showModal">新增用户</Button>
      </template>
      <Table :columns="columns" :dataSource="dataSource" rowKey="id" />
    </Card>

    <Modal
      v-model:visible="visible"
      title="新增用户"
      @ok="handleOk"
      @cancel="handleCancel"
      :maskClosable="false"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="username"
          label="用户名"
          :rules="[{ required: true, message: '请输入用户名' }]"
        >
          <Input
            v-model:value="formState.username"
            placeholder="请输入用户名"
          />
        </Form.Item>
        <Form.Item
          name="name"
          label="姓名"
          :rules="[{ required: true, message: '请输入姓名' }]"
        >
          <Input v-model:value="formState.name" placeholder="请输入姓名" />
        </Form.Item>
        <Form.Item
          name="email"
          label="邮箱"
          :rules="[
            { required: true, message: '请输入邮箱' },
            { type: 'email', message: '请输入有效的邮箱地址' },
          ]"
        >
          <Input v-model:value="formState.email" placeholder="请输入邮箱" />
        </Form.Item>
        <Form.Item
          name="phone"
          label="手机号"
          :rules="[{ required: true, message: '请输入手机号' }]"
        >
          <Input v-model:value="formState.phone" placeholder="请输入手机号" />
        </Form.Item>
        <Form.Item name="role" label="角色">
          <Select v-model:value="formState.role">
            <Select.Option value="超级管理员">超级管理员</Select.Option>
            <Select.Option value="开发者">开发者</Select.Option>
            <Select.Option value="普通用户">普通用户</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="password"
          label="密码"
          :rules="[{ required: true, message: '请输入密码' }]"
        >
          <Input.Password
            v-model:value="formState.password"
            placeholder="请输入密码"
          />
        </Form.Item>
        <Form.Item
          name="confirmPassword"
          label="确认密码"
          :rules="[
            { required: true, message: '请确认密码' },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('password') === value) {
                  return Promise.resolve();
                }
                return Promise.reject('两次输入的密码不匹配');
              },
            }),
          ]"
        >
          <Input.Password
            v-model:value="formState.confirmPassword"
            placeholder="请确认密码"
          />
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.system-user-container {
  padding: 0;
}
</style>
