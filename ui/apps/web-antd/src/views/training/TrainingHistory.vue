<script lang="ts" setup>
import { ref, h } from 'vue';
import { Table, Card, Button, Space, Tag, DatePicker } from 'ant-design-vue';

defineOptions({ name: 'TrainingHistory' });

const dataSource = ref([
  {
    id: '1',
    name: 'BERT微调任务',
    dataset: 'NLP语料集',
    model: 'BERT-base',
    gpu: 'NVIDIA A100',
    status: '已完成',
    startTime: '2023-11-01 10:30:00',
    endTime: '2023-11-01 12:30:00',
    duration: '2小时',
    user: '张三',
    result: '成功',
  },
  {
    id: '2',
    name: '图像分类训练',
    dataset: '图像数据集',
    model: 'ResNet50',
    gpu: 'NVIDIA V100',
    status: '已完成',
    startTime: '2023-11-02 09:15:00',
    endTime: '2023-11-02 14:15:00',
    duration: '5小时',
    user: '李四',
    result: '成功',
  },
  {
    id: '3',
    name: 'GPT-2微调',
    dataset: '文本数据集',
    model: 'GPT-2',
    gpu: 'NVIDIA A100 x4',
    status: '已完成',
    startTime: '2023-11-03 08:00:00',
    endTime: '2023-11-03 20:00:00',
    duration: '12小时',
    user: '王五',
    result: '失败',
  },
]);

const columns = [
  {
    title: '任务ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '任务名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '数据集',
    dataIndex: 'dataset',
    key: 'dataset',
  },
  {
    title: '模型',
    dataIndex: 'model',
    key: 'model',
  },
  {
    title: 'GPU资源',
    dataIndex: 'gpu',
    key: 'gpu',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
  },
  {
    title: '开始时间',
    dataIndex: 'startTime',
    key: 'startTime',
    sorter: (a, b) => new Date(a.startTime) - new Date(b.startTime),
  },
  {
    title: '结束时间',
    dataIndex: 'endTime',
    key: 'endTime',
  },
  {
    title: '运行时长',
    dataIndex: 'duration',
    key: 'duration',
  },
  {
    title: '执行用户',
    dataIndex: 'user',
    key: 'user',
  },
  {
    title: '结果',
    dataIndex: 'result',
    key: 'result',
    customRender: ({ text }) => {
      const color = text === '成功' ? 'green' : 'red';
      return h(Tag, { color }, () => text);
    },
  },
  {
    title: '操作',
    key: 'action',
    customRender: () => {
      return h(Space, { size: 'middle' }, [
        h('a', {}, '查看详情'),
        h('a', {}, '查看日志'),
        h('a', {}, '下载模型')
      ]);
    },
  },
];

// 日期筛选
const { RangePicker } = DatePicker;
const dateRange = ref([]);
const handleDateRangeChange = (value) => {
  dateRange.value = value;
  // 这里可以添加基于日期筛选的逻辑
};
</script>

<template>
  <div class="training-history-container">
    <Card title="历史任务记录" :bordered="false">
      <template #extra>
        <Space>
          <RangePicker
            v-model:value="dateRange"
            @change="handleDateRangeChange"
          />
          <Button type="primary">导出记录</Button>
        </Space>
      </template>
      <Table
        :columns="columns"
        :dataSource="dataSource"
        rowKey="id"
        :pagination="{ pageSize: 10 }"
      />
    </Card>
  </div>
</template>

<style scoped>
.training-history-container {
  padding: 0;
}
</style>
