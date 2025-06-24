<script lang="ts" setup>
import { ref, h } from 'vue';
import {
  Table,
  Card,
  Button,
  Space,
  Tag,
  Modal,
  Input,
  Form,
  DatePicker,
} from 'ant-design-vue';

defineOptions({ name: 'WorkspaceList' });

const dataSource = ref([
  {
    id: '1',
    name: 'AI训练项目1',
    description: '用于NLP模型训练的工作空间',
    owner: '张三',
    members: 5,
    projects: 3,
    status: '活跃',
    createTime: '2023-10-01',
    lastActive: '2023-11-15',
  },
  {
    id: '2',
    name: '图像识别研究',
    description: '针对图像识别算法的研发工作空间',
    owner: '李四',
    members: 8,
    projects: 2,
    status: '活跃',
    createTime: '2023-09-15',
    lastActive: '2023-11-14',
  },
  {
    id: '3',
    name: '推荐系统开发',
    description: '电商推荐系统算法工作空间',
    owner: '王五',
    members: 4,
    projects: 1,
    status: '归档',
    createTime: '2023-08-20',
    lastActive: '2023-10-30',
  },
]);

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '工作空间名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '创建人',
    dataIndex: 'owner',
    key: 'owner',
  },
  {
    title: '成员数',
    dataIndex: 'members',
    key: 'members',
  },
  {
    title: '项目数',
    dataIndex: 'projects',
    key: 'projects',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      const color = text === '活跃' ? 'green' : text === '归档' ? 'orange' : 'default';
      return h(Tag, { color }, () => text);
    },
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    key: 'createTime',
  },
  {
    title: '最近活动',
    dataIndex: 'lastActive',
    key: 'lastActive',
  },
  {
    title: '操作',
    key: 'action',
    customRender: () => {
      return h(Space, { size: 'middle' }, [
        h('a', {}, '进入'),
        h('a', {}, '编辑'),
        h('a', {}, '归档')
      ]);
    },
  },
];

const visible = ref(false);
const formRef = ref();
const formState = ref({
  name: '',
  description: '',
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
        owner: '当前用户',
        members: 1,
        projects: 0,
        status: '活跃',
        createTime: new Date().toLocaleDateString(),
        lastActive: new Date().toLocaleDateString(),
      });
      visible.value = false;
      formState.value = {
        name: '',
        description: '',
      };
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};
</script>

<template>
  <div class="workspace-list-container">
    <Card title="工作空间列表" :bordered="false">
      <template #extra>
        <Button type="primary" @click="showModal">创建工作空间</Button>
      </template>
      <Table :columns="columns" :dataSource="dataSource" rowKey="id" />
    </Card>

    <Modal
      v-model:visible="visible"
      title="创建工作空间"
      @ok="handleOk"
      @cancel="handleCancel"
      :maskClosable="false"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="name"
          label="工作空间名称"
          :rules="[{ required: true, message: '请输入工作空间名称' }]"
        >
          <Input v-model:value="formState.name" placeholder="请输入工作空间名称" />
        </Form.Item>
        <Form.Item name="description" label="描述">
          <Input.TextArea
            v-model:value="formState.description"
            placeholder="请输入工作空间描述"
          />
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.workspace-list-container {
  padding: 0;
}
</style>
