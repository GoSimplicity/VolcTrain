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
  Divider,
  Tooltip,
} from 'ant-design-vue';
import type { SelectValue } from 'ant-design-vue/es/select';

defineOptions({ name: 'DatasetList' });

interface DatasetItem {
  id: string;
  name: string;
  type: string;
  description: string;
  size: string;
  files: number;
  format: string;
  creator: string;
  createTime: string;
  updateTime: string;
  status: string;
  shared: boolean;
}

// 数据集列表数据
const dataSource = ref<DatasetItem[]>([
  {
    id: '1',
    name: 'COCO数据集',
    type: '图像分类',
    description: '常用目标检测数据集',
    size: '18.3GB',
    files: 123287,
    format: 'JPG, JSON',
    creator: '张三',
    createTime: '2023-10-05',
    updateTime: '2023-11-10',
    status: '已处理',
    shared: true,
  },
  {
    id: '2',
    name: '中文语料库',
    type: '文本',
    description: '大规模中文自然语言处理语料库',
    size: '5.2GB',
    files: 8521,
    format: 'TXT, CSV',
    creator: '李四',
    createTime: '2023-09-15',
    updateTime: '2023-11-12',
    status: '已处理',
    shared: true,
  },
  {
    id: '3',
    name: 'ImageNet子集',
    type: '图像分类',
    description: 'ImageNet数据集的一个子集',
    size: '45.8GB',
    files: 50000,
    format: 'JPG, PNG',
    creator: '王五',
    createTime: '2023-11-01',
    updateTime: '2023-11-14',
    status: '处理中',
    shared: false,
  },
  {
    id: '4',
    name: '语音识别数据集',
    type: '音频',
    description: '中英文混合语音识别数据集',
    size: '12.6GB',
    files: 25000,
    format: 'WAV, MP3',
    creator: '赵六',
    createTime: '2023-10-20',
    updateTime: '2023-11-05',
    status: '已处理',
    shared: false,
  },
  {
    id: '5',
    name: '电商评论数据',
    type: '文本',
    description: '电商产品评论情感分析数据集',
    size: '1.8GB',
    files: 15000,
    format: 'CSV, JSON',
    creator: '张三',
    createTime: '2023-11-08',
    updateTime: '2023-11-15',
    status: '未处理',
    shared: true,
  },
]);

// 表格列配置
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '数据集名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
    customRender: ({ text }: { text: string }) => {
      const color =
        text === '图像分类'
          ? 'blue'
          : text === '文本'
            ? 'green'
            : text === '音频'
              ? 'purple'
              : 'default';
      return h(Tag, { color }, () => text);
    },
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
    customRender: ({ text }: { text: string }) => {
      return h(Tooltip, { placement: 'topLeft', title: text }, () =>
        h(
          'div',
          {
            style:
              'width: 150px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;',
          },
          text,
        ),
      );
    },
  },
  {
    title: '大小',
    dataIndex: 'size',
    key: 'size',
  },
  {
    title: '文件数',
    dataIndex: 'files',
    key: 'files',
  },
  {
    title: '格式',
    dataIndex: 'format',
    key: 'format',
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
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
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }: { text: string }) => {
      const color =
        text === '已处理'
          ? 'green'
          : text === '处理中'
            ? 'blue'
            : text === '未处理'
              ? 'orange'
              : 'default';
      return h(Tag, { color }, () => text);
    },
  },
  {
    title: '共享',
    dataIndex: 'shared',
    key: 'shared',
    customRender: ({ text }: { text: boolean }) => {
      return text
        ? h(Tag, { color: 'green' }, () => '是')
        : h(Tag, {}, () => '否');
    },
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }: { record: DatasetItem }) =>
      h(Space, { size: 'middle' }, [
        h('a', { onClick: () => console.log('查看', record) }, '查看'),
        h('a', { onClick: () => console.log('下载', record) }, '下载'),
        record.status === '未处理' &&
          h('a', { onClick: () => console.log('处理', record) }, '处理'),
        h('a', { onClick: () => console.log('删除', record) }, '删除'),
      ]),
  },
];

// 表单相关
const visible = ref(false);
const formRef = ref();
const formState = ref({
  name: '',
  type: '',
  description: '',
  shared: false as boolean,
});

const typeOptions = [
  { label: '图像分类', value: '图像分类' },
  { label: '文本', value: '文本' },
  { label: '音频', value: '音频' },
  { label: '视频', value: '视频' },
  { label: '表格', value: '表格' },
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
      // 模拟添加一条新数据
      dataSource.value.push({
        id: String(dataSource.value.length + 1),
        name: formState.value.name,
        type: formState.value.type,
        description: formState.value.description,
        size: '0KB',
        files: 0,
        format: '-',
        creator: '当前用户',
        createTime: new Date().toLocaleDateString(),
        updateTime: new Date().toLocaleDateString(),
        status: '未处理',
        shared: formState.value.shared,
      });

      visible.value = false;
      formState.value = {
        name: '',
        type: '',
        description: '',
        shared: false,
      };
    })
    .catch((error: any) => {
      console.error('验证失败:', error);
    });
};

// 搜索表单
const searchForm = ref({
  name: '',
  type: '',
  status: '',
  creator: '',
});

const handleSearch = () => {
  // 实际场景中，这里应该根据搜索条件过滤数据
  console.log('搜索条件:', searchForm.value);
};

const resetSearch = () => {
  searchForm.value = {
    name: '',
    type: '',
    status: '',
    creator: '',
  };
};
</script>

<template>
  <div class="dataset-list-container">
    <Card title="数据集管理">
      <template #extra>
        <Space>
          <Button @click="showModal" type="primary">新建数据集</Button>
          <Button>导入数据集</Button>
        </Space>
      </template>

      <!-- 搜索区域 -->
      <div class="search-area">
        <Form layout="inline" :model="searchForm">
          <Form.Item name="name" label="数据集名称">
            <Input
              v-model:value="searchForm.name"
              placeholder="请输入数据集名称"
            />
          </Form.Item>
          <Form.Item name="type" label="数据类型">
            <Select
              v-model:value="searchForm.type"
              placeholder="请选择类型"
              style="width: 120px"
              allowClear
            >
              <Select.Option value="图像分类">图像分类</Select.Option>
              <Select.Option value="文本">文本</Select.Option>
              <Select.Option value="音频">音频</Select.Option>
              <Select.Option value="视频">视频</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select
              v-model:value="searchForm.status"
              placeholder="请选择状态"
              style="width: 120px"
              allowClear
            >
              <Select.Option value="已处理">已处理</Select.Option>
              <Select.Option value="处理中">处理中</Select.Option>
              <Select.Option value="未处理">未处理</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" @click="handleSearch">搜索</Button>
            <Button style="margin-left: 8px" @click="resetSearch">重置</Button>
          </Form.Item>
        </Form>
      </div>

      <Divider style="margin: 12px 0" />

      <Table :columns="columns" :dataSource="dataSource" rowKey="id" />
    </Card>

    <Modal
      v-model:visible="visible"
      title="新建数据集"
      @ok="handleOk"
      @cancel="handleCancel"
      :maskClosable="false"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="name"
          label="数据集名称"
          :rules="[{ required: true, message: '请输入数据集名称' }]"
        >
          <Input
            v-model:value="formState.name"
            placeholder="请输入数据集名称"
          />
        </Form.Item>
        <Form.Item
          name="type"
          label="数据类型"
          :rules="[{ required: true, message: '请选择数据类型' }]"
        >
          <Select v-model:value="formState.type" placeholder="请选择数据类型">
            <Select.Option
              v-for="option in typeOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item name="description" label="描述">
          <Input.TextArea
            v-model:value="formState.description"
            placeholder="请输入数据集描述"
          />
        </Form.Item>
        <Form.Item name="shared" label="是否共享">
          <Select v-model:value="formState.shared">
            <Select.Option :value="true">是</Select.Option>
            <Select.Option :value="false">否</Select.Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.dataset-list-container {
  padding: 0;
}
.search-area {
  margin-bottom: 16px;
}
</style>
