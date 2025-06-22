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
  Select,
  Tabs,
} from 'ant-design-vue';

defineOptions({ name: 'WorkspaceProject' });

const dataSource = ref([
  {
    id: '1',
    name: 'BERT微调项目',
    workspace: 'AI训练项目1',
    type: '模型训练',
    description: 'BERT模型微调的项目',
    owner: '张三',
    members: 3,
    status: '进行中',
    createTime: '2023-10-05',
    lastActive: '2023-11-15',
  },
  {
    id: '2',
    name: '图像分类研究',
    workspace: '图像识别研究',
    type: '实验研究',
    description: '基于ResNet的图像分类模型研究',
    owner: '李四',
    members: 5,
    status: '进行中',
    createTime: '2023-09-20',
    lastActive: '2023-11-14',
  },
  {
    id: '3',
    name: '推荐算法开发',
    workspace: '推荐系统开发',
    type: '开发实现',
    description: '电商推荐算法实现与优化',
    owner: '王五',
    members: 2,
    status: '已完成',
    createTime: '2023-08-25',
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
    title: '项目名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '所属工作空间',
    dataIndex: 'workspace',
    key: 'workspace',
  },
  {
    title: '项目类型',
    dataIndex: 'type',
    key: 'type',
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
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      const color =
        text === '进行中' ? 'blue' : text === '已完成' ? 'green' : 'default';
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
        h('a', {}, '详情'),
        h('a', {}, '编辑'),
        h('a', {}, '删除')
      ]);
    },
  },
];

const visible = ref(false);
const formRef = ref();
const formState = ref({
  name: '',
  workspace: '',
  type: '',
  description: '',
});

const workspaceOptions = [
  { label: 'AI训练项目1', value: 'AI训练项目1' },
  { label: '图像识别研究', value: '图像识别研究' },
  { label: '推荐系统开发', value: '推荐系统开发' },
];

const typeOptions = [
  { label: '模型训练', value: '模型训练' },
  { label: '实验研究', value: '实验研究' },
  { label: '开发实现', value: '开发实现' },
  { label: '数据处理', value: '数据处理' },
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
        workspace: formState.value.workspace,
        type: formState.value.type,
        description: formState.value.description,
        owner: '当前用户',
        members: 1,
        status: '进行中',
        createTime: new Date().toLocaleDateString(),
        lastActive: new Date().toLocaleDateString(),
      });
      visible.value = false;
      formState.value = {
        name: '',
        workspace: '',
        type: '',
        description: '',
      };
    })
    .catch((error) => {
      console.error('验证失败:', error);
    });
};

const activeKey = ref('1');
</script>

<template>
  <div class="workspace-project-container">
    <Card title="项目管理" :bordered="false">
      <template #extra>
        <Button type="primary" @click="showModal">创建项目</Button>
      </template>
      
      <Tabs v-model:activeKey="activeKey">
        <Tabs.TabPane key="1" tab="全部项目">
          <Table :columns="columns" :dataSource="dataSource" rowKey="id" />
        </Tabs.TabPane>
        <Tabs.TabPane key="2" tab="我的项目">
          <Table 
            :columns="columns" 
            :dataSource="dataSource.filter(item => item.owner === '当前用户')" 
            rowKey="id" 
          />
        </Tabs.TabPane>
        <Tabs.TabPane key="3" tab="参与的项目">
          <Table 
            :columns="columns" 
            :dataSource="dataSource.filter(item => item.owner !== '当前用户')" 
            rowKey="id" 
          />
        </Tabs.TabPane>
      </Tabs>
    </Card>

    <Modal
      v-model:visible="visible"
      title="创建项目"
      @ok="handleOk"
      @cancel="handleCancel"
      :maskClosable="false"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="name"
          label="项目名称"
          :rules="[{ required: true, message: '请输入项目名称' }]"
        >
          <Input v-model:value="formState.name" placeholder="请输入项目名称" />
        </Form.Item>
        <Form.Item
          name="workspace"
          label="所属工作空间"
          :rules="[{ required: true, message: '请选择所属工作空间' }]"
        >
          <Select
            v-model:value="formState.workspace"
            placeholder="请选择所属工作空间"
          >
            <Select.Option
              v-for="option in workspaceOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="type"
          label="项目类型"
          :rules="[{ required: true, message: '请选择项目类型' }]"
        >
          <Select v-model:value="formState.type" placeholder="请选择项目类型">
            <Select.Option
              v-for="option in typeOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item name="description" label="项目描述">
          <Input.TextArea
            v-model:value="formState.description"
            placeholder="请输入项目描述"
          />
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.workspace-project-container {
  padding: 0;
}
</style>
