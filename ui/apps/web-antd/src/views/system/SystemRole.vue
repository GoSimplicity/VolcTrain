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
  Tree,
  Checkbox,
} from 'ant-design-vue';

defineOptions({ name: 'SystemRole' });

const dataSource = ref([
  {
    id: '1',
    name: '超级管理员',
    code: 'SUPER_ADMIN',
    description: '拥有所有权限',
    userCount: 1,
    status: '启用',
    createTime: '2023-10-01 08:00:00',
    updateTime: '2023-11-15 09:30:00',
  },
  {
    id: '2',
    name: '开发者',
    code: 'DEVELOPER',
    description: '拥有开发相关权限',
    userCount: 5,
    status: '启用',
    createTime: '2023-10-02 10:00:00',
    updateTime: '2023-11-10 15:45:00',
  },
  {
    id: '3',
    name: '普通用户',
    code: 'NORMAL_USER',
    description: '基本操作权限',
    userCount: 10,
    status: '启用',
    createTime: '2023-10-03 09:00:00',
    updateTime: '2023-11-05 11:30:00',
  },
  {
    id: '4',
    name: '只读用户',
    code: 'READONLY_USER',
    description: '只有查看权限',
    userCount: 8,
    status: '禁用',
    createTime: '2023-10-04 14:00:00',
    updateTime: '2023-10-20 16:20:00',
  },
]);

const columns = [
  {
    title: '角色ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '角色名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '角色编码',
    dataIndex: 'code',
    key: 'code',
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '用户数',
    dataIndex: 'userCount',
    key: 'userCount',
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
        h('a', {}, '权限设置'),
        h('a', {}, '删除')
      ]);
    },
  },
];

const visible = ref(false);
const formRef = ref();
const formState = ref({
  name: '',
  code: '',
  description: '',
  status: true,
});

// 权限树
const permissionTreeData = ref([
  {
    title: '系统管理',
    key: 'system',
    children: [
      {
        title: '用户管理',
        key: 'system:user',
        children: [
          { title: '查询用户', key: 'system:user:list' },
          { title: '新增用户', key: 'system:user:add' },
          { title: '修改用户', key: 'system:user:edit' },
          { title: '删除用户', key: 'system:user:delete' },
        ],
      },
      {
        title: '角色管理',
        key: 'system:role',
        children: [
          { title: '查询角色', key: 'system:role:list' },
          { title: '新增角色', key: 'system:role:add' },
          { title: '修改角色', key: 'system:role:edit' },
          { title: '删除角色', key: 'system:role:delete' },
        ],
      },
      {
        title: 'API管理',
        key: 'system:api',
        children: [
          { title: '查询API', key: 'system:api:list' },
          { title: '新增API', key: 'system:api:add' },
          { title: '修改API', key: 'system:api:edit' },
          { title: '删除API', key: 'system:api:delete' },
        ],
      },
    ],
  },
]);

const checkedKeys = ref(['system:user:list', 'system:role:list']);

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
        code: formState.value.code,
        description: formState.value.description,
        userCount: 0,
        status: formState.value.status ? '启用' : '禁用',
        createTime: new Date().toLocaleString(),
        updateTime: new Date().toLocaleString(),
      });
      visible.value = false;
      formState.value = {
        name: '',
        code: '',
        description: '',
        status: true,
      };
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};

const onCheck = (checked) => {
  checkedKeys.value = checked;
};
</script>

<template>
  <div class="system-role-container">
    <Card title="角色管理" :bordered="false">
      <template #extra>
        <Button type="primary" @click="showModal">新增角色</Button>
      </template>
      <Table :columns="columns" :dataSource="dataSource" rowKey="id" />
    </Card>

    <Modal
      v-model:visible="visible"
      title="新增角色"
      @ok="handleOk"
      @cancel="handleCancel"
      :maskClosable="false"
      width="700px"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="name"
          label="角色名称"
          :rules="[{ required: true, message: '请输入角色名称' }]"
        >
          <Input v-model:value="formState.name" placeholder="请输入角色名称" />
        </Form.Item>
        <Form.Item
          name="code"
          label="角色编码"
          :rules="[{ required: true, message: '请输入角色编码' }]"
        >
          <Input v-model:value="formState.code" placeholder="请输入角色编码" />
        </Form.Item>
        <Form.Item name="description" label="角色描述">
          <Input.TextArea
            v-model:value="formState.description"
            placeholder="请输入角色描述"
          />
        </Form.Item>
        <Form.Item name="status" label="角色状态">
          <Checkbox v-model:checked="formState.status">启用</Checkbox>
        </Form.Item>
        <Form.Item name="permissions" label="角色权限">
          <div class="permission-tree">
            <Tree
              checkable
              :treeData="permissionTreeData"
              v-model:checkedKeys="checkedKeys"
              @check="onCheck"
            />
          </div>
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.system-role-container {
  padding: 0;
}
.permission-tree {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #f0f0f0;
  padding: 10px;
}
</style>
