<script lang="ts" setup>
import { ref } from 'vue';
import {
  Card,
  Table,
  Tag,
  Space,
  Button,
  Modal,
  Form,
  Input,
  DatePicker,
  TimePicker,
  Select,
  InputNumber,
  Badge,
  Tabs,
} from 'ant-design-vue';

defineOptions({ name: 'GPUSchedule' });

// GPU调度任务数据
const dataSource = ref([
  {
    id: '1',
    name: 'BERT微调训练',
    user: '张三',
    gpuType: 'NVIDIA A100',
    gpuCount: 2,
    priority: '高',
    status: '运行中',
    startTime: '2023-11-15 08:30',
    endTime: '2023-11-15 14:30',
    duration: '6小时',
  },
  {
    id: '2',
    name: 'ResNet图像分类训练',
    user: '李四',
    gpuType: 'NVIDIA A100',
    gpuCount: 4,
    priority: '中',
    status: '等待中',
    startTime: '2023-11-15 14:30',
    endTime: '2023-11-16 02:30',
    duration: '12小时',
  },
  {
    id: '3',
    name: 'GPT-2模型训练',
    user: '王五',
    gpuType: 'NVIDIA V100',
    gpuCount: 8,
    priority: '中',
    status: '等待中',
    startTime: '2023-11-16 08:00',
    endTime: '2023-11-18 08:00',
    duration: '48小时',
  },
  {
    id: '4',
    name: '推荐算法训练',
    user: '赵六',
    gpuType: 'NVIDIA V100',
    gpuCount: 4,
    priority: '低',
    status: '已预约',
    startTime: '2023-11-16 14:00',
    endTime: '2023-11-17 14:00',
    duration: '24小时',
  },
  {
    id: '5',
    name: '语音识别模型训练',
    user: '孙七',
    gpuType: 'NVIDIA A100',
    gpuCount: 2,
    priority: '高',
    status: '已完成',
    startTime: '2023-11-14 10:00',
    endTime: '2023-11-15 04:00',
    duration: '18小时',
  },
]);

// 可用资源数据
const resourcesData = ref([
  {
    id: '1',
    type: 'NVIDIA A100',
    total: 16,
    used: 6,
    available: 10,
    maintenance: 0,
    utilization: 37.5,
  },
  {
    id: '2',
    type: 'NVIDIA V100',
    total: 24,
    used: 8,
    available: 14,
    maintenance: 2,
    utilization: 33.3,
  },
  {
    id: '3',
    type: 'NVIDIA T4',
    total: 32,
    used: 12,
    available: 18,
    maintenance: 2,
    utilization: 37.5,
  },
]);

// 表格列配置
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
    title: '用户',
    dataIndex: 'user',
    key: 'user',
  },
  {
    title: 'GPU型号',
    dataIndex: 'gpuType',
    key: 'gpuType',
  },
  {
    title: 'GPU数量',
    dataIndex: 'gpuCount',
    key: 'gpuCount',
  },
  {
    title: '优先级',
    dataIndex: 'priority',
    key: 'priority',
    customRender: ({ text }) => {
      const color = 
        text === '高' ? 'red' : 
        text === '中' ? 'orange' : 
        text === '低' ? 'blue' : 'default';
      return <Tag color={color}>{text}</Tag>;
    },
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }) => {
      let color = 'default';
      let status = 'default';
      
      switch (text) {
        case '运行中':
          color = 'green';
          status = 'processing';
          break;
        case '等待中':
          color = 'orange';
          status = 'warning';
          break;
        case '已预约':
          color = 'blue';
          status = 'default';
          break;
        case '已完成':
          color = 'default';
          status = 'success';
          break;
      }
      
      return <Badge status={status} text={text} />;
    },
  },
  {
    title: '开始时间',
    dataIndex: 'startTime',
    key: 'startTime',
  },
  {
    title: '结束时间',
    dataIndex: 'endTime',
    key: 'endTime',
  },
  {
    title: '持续时间',
    dataIndex: 'duration',
    key: 'duration',
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }) => {
      return (
        <Space size="middle">
          {record.status === '等待中' && <a onClick={() => handleStart(record)}>启动</a>}
          {record.status === '运行中' && <a onClick={() => handleStop(record)}>停止</a>}
          {(record.status === '等待中' || record.status === '已预约') && <a onClick={() => handleCancel(record)}>取消</a>}
          {record.status !== '已完成' && <a onClick={() => handleModify(record)}>修改</a>}
          <a onClick={() => handleViewDetails(record)}>详情</a>
        </Space>
      );
    },
  },
]);

// 资源表格列配置
const resourceColumns = [
  {
    title: 'GPU型号',
    dataIndex: 'type',
    key: 'type',
  },
  {
    title: '总数量',
    dataIndex: 'total',
    key: 'total',
  },
  {
    title: '使用中',
    dataIndex: 'used',
    key: 'used',
  },
  {
    title: '可用数量',
    dataIndex: 'available',
    key: 'available',
  },
  {
    title: '维护中',
    dataIndex: 'maintenance',
    key: 'maintenance',
  },
  {
    title: '使用率',
    dataIndex: 'utilization',
    key: 'utilization',
    customRender: ({ text }) => {
      return `${text}%`;
    },
  },
  {
    title: '操作',
    key: 'action',
    customRender: () => (
      <Space size="middle">
        <Button type="primary" size="small">预约</Button>
      </Space>
    ),
  },
];

const visible = ref(false);
const formRef = ref();
const formState = ref({
  name: '',
  gpuType: '',
  gpuCount: 1,
  startDate: null,
  startTime: null,
  duration: 1,
  priority: '中',
});

const gpuTypeOptions = [
  { label: 'NVIDIA A100', value: 'NVIDIA A100' },
  { label: 'NVIDIA V100', value: 'NVIDIA V100' },
  { label: 'NVIDIA T4', value: 'NVIDIA T4' },
];

const priorityOptions = [
  { label: '高', value: '高' },
  { label: '中', value: '中' },
  { label: '低', value: '低' },
];

const activeKey = ref('1');

// 处理函数
const showModal = () => {
  visible.value = true;
  formState.value = {
    name: '',
    gpuType: '',
    gpuCount: 1,
    startDate: null,
    startTime: null,
    duration: 1,
    priority: '中',
  };
};

const handleCancel = () => {
  visible.value = false;
};

const handleOk = () => {
  formRef.value.validate().then(() => {
    // 构建开始和结束时间字符串
    const startDateStr = formState.value.startDate.format('YYYY-MM-DD');
    const startTimeStr = formState.value.startTime.format('HH:mm');
    const startTimeFullStr = `${startDateStr} ${startTimeStr}`;
    
    // 计算结束时间
    const durationHours = formState.value.duration;
    const endTimeFullStr = '计算得到的结束时间'; // 实际应用中应该计算真实的结束时间
    
    dataSource.value.push({
      id: String(dataSource.value.length + 1),
      name: formState.value.name,
      user: '当前用户', // 实际应用中应该是登录用户
      gpuType: formState.value.gpuType,
      gpuCount: formState.value.gpuCount,
      priority: formState.value.priority,
      status: '已预约',
      startTime: startTimeFullStr,
      endTime: endTimeFullStr,
      duration: `${durationHours}小时`,
    });
    
    visible.value = false;
  }).catch(error => {
    console.log('表单验证失败:', error);
  });
};

const handleStart = (record) => {
  const index = dataSource.value.findIndex(item => item.id === record.id);
  dataSource.value[index].status = '运行中';
};

const handleStop = (record) => {
  const index = dataSource.value.findIndex(item => item.id === record.id);
  dataSource.value[index].status = '已完成';
};

const handleCancel = (record) => {
  const index = dataSource.value.findIndex(item => item.id === record.id);
  dataSource.value.splice(index, 1);
};

const handleModify = (record) => {
  // 实际应用中应该打开修改表单
  console.log('修改任务:', record);
};

const handleViewDetails = (record) => {
  // 实际应用中应该打开详情页面
  console.log('查看详情:', record);
};
</script>

<template>
  <div class="gpu-schedule-container">
    <Tabs v-model:activeKey="activeKey">
      <Tabs.TabPane key="1" tab="调度任务">
        <Card title="GPU调度任务">
          <template #extra>
            <Button type="primary" @click="showModal">预约GPU资源</Button>
          </template>
          <Table :columns="columns" :dataSource="dataSource" rowKey="id" />
        </Card>
      </Tabs.TabPane>
      <Tabs.TabPane key="2" tab="可用资源">
        <Card title="GPU可用资源">
          <Table :columns="resourceColumns" :dataSource="resourcesData" rowKey="id" />
        </Card>
      </Tabs.TabPane>
      <Tabs.TabPane key="3" tab="调度日历">
        <Card title="GPU资源调度日历">
          <div class="calendar-placeholder">
            <p>这里应该是一个资源调度日历，显示GPU资源的预约和使用情况</p>
            <p>可以按照时间轴展示各个GPU资源的分配情况</p>
          </div>
        </Card>
      </Tabs.TabPane>
    </Tabs>

    <!-- 预约GPU资源表单 -->
    <Modal
      v-model:visible="visible"
      title="预约GPU资源"
      @ok="handleOk"
      @cancel="handleCancel"
      :maskClosable="false"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="name"
          label="任务名称"
          :rules="[{ required: true, message: '请输入任务名称' }]"
        >
          <Input v-model:value="formState.name" placeholder="请输入任务名称" />
        </Form.Item>
        <Form.Item
          name="gpuType"
          label="GPU型号"
          :rules="[{ required: true, message: '请选择GPU型号' }]"
        >
          <Select v-model:value="formState.gpuType" placeholder="请选择GPU型号">
            <Select.Option
              v-for="option in gpuTypeOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="gpuCount"
          label="GPU数量"
          :rules="[{ required: true, message: '请输入GPU数量' }]"
        >
          <InputNumber
            v-model:value="formState.gpuCount"
            :min="1"
            style="width: 100%"
          />
        </Form.Item>
        <Form.Item
          name="startDate"
          label="开始日期"
          :rules="[{ required: true, message: '请选择开始日期' }]"
        >
          <DatePicker
            v-model:value="formState.startDate"
            style="width: 100%"
            placeholder="选择日期"
          />
        </Form.Item>
        <Form.Item
          name="startTime"
          label="开始时间"
          :rules="[{ required: true, message: '请选择开始时间' }]"
        >
          <TimePicker
            v-model:value="formState.startTime"
            style="width: 100%"
            placeholder="选择时间"
            format="HH:mm"
          />
        </Form.Item>
        <Form.Item
          name="duration"
          label="持续时间(小时)"
          :rules="[{ required: true, message: '请输入持续时间' }]"
        >
          <InputNumber
            v-model:value="formState.duration"
            :min="1"
            style="width: 100%"
          />
        </Form.Item>
        <Form.Item name="priority" label="优先级">
          <Select v-model:value="formState.priority">
            <Select.Option
              v-for="option in priorityOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.gpu-schedule-container {
  padding: 0;
}
.calendar-placeholder {
  height: 500px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #fafafa;
  border: 1px dashed #d9d9d9;
  border-radius: 2px;
}
</style>
