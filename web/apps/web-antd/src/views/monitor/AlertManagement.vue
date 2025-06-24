<script lang="ts" setup>
import { ref, h } from 'vue';
import {
  Card,
  Table,
  Tag,
  Space,
  Button,
  Tabs,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  Switch,
  Popconfirm,
  Checkbox,
} from 'ant-design-vue';

defineOptions({ name: 'AlertManagement' });

interface AlertRule {
  id: string;
  name: string;
  description: string;
  type: string;
  metric: string;
  threshold: string;
  duration: string;
  severity: string;
  notifyChannel: string;
  status: string;
  createTime: string;
}

interface AlertHistory {
  id: string;
  ruleName: string;
  target: string;
  description: string;
  severity: string;
  time: string;
  duration: string;
  status: string;
  handler: string;
  handleTime: string;
}

interface FormState {
  name: string;
  description: string;
  type: string;
  metric: string;
  operator: string;
  threshold: number;
  duration: number;
  severity: string;
  notifyMail: boolean;
  notifySms: boolean;
  notifyPhone: boolean;
  status: boolean;
}

interface Option {
  label: string;
  value: string;
}

type CheckedType = boolean | (() => boolean);

// 告警规则数据
const alertRules = ref<AlertRule[]>([
  {
    id: '1',
    name: 'CPU使用率过高',
    description: '集群CPU使用率超过阈值告警',
    type: '集群资源',
    metric: 'CPU使用率',
    threshold: '> 90%',
    duration: '5分钟',
    severity: '严重',
    notifyChannel: '邮件、短信',
    status: '启用',
    createTime: '2023-10-15 09:30',
  },
  {
    id: '2',
    name: 'GPU显存不足',
    description: '训练任务GPU显存不足告警',
    type: '任务资源',
    metric: 'GPU显存使用率',
    threshold: '> 95%',
    duration: '3分钟',
    severity: '严重',
    notifyChannel: '邮件',
    status: '启用',
    createTime: '2023-10-16 10:15',
  },
  {
    id: '3',
    name: '节点离线告警',
    description: '计算节点离线告警',
    type: '节点状态',
    metric: '节点状态',
    threshold: '离线',
    duration: '即时',
    severity: '紧急',
    notifyChannel: '邮件、短信、电话',
    status: '启用',
    createTime: '2023-10-17 14:20',
  },
  {
    id: '4',
    name: '训练任务失败',
    description: '训练任务执行失败告警',
    type: '任务状态',
    metric: '任务状态',
    threshold: '失败',
    duration: '即时',
    severity: '中等',
    notifyChannel: '邮件',
    status: '禁用',
    createTime: '2023-10-18 16:45',
  },
]);

// 告警历史数据
const alertHistory = ref<AlertHistory[]>([
  {
    id: '1001',
    ruleName: 'CPU使用率过高',
    target: 'node-01',
    description: '节点CPU使用率达到95%',
    severity: '严重',
    time: '2023-11-15 10:23:45',
    duration: '15分钟',
    status: '已处理',
    handler: '张三',
    handleTime: '2023-11-15 10:38:12',
  },
  {
    id: '1002',
    ruleName: 'GPU显存不足',
    target: '任务ID:2',
    description: '训练任务GPU显存使用率达到98%',
    severity: '严重',
    time: '2023-11-15 08:12:30',
    duration: '10分钟',
    status: '已处理',
    handler: '李四',
    handleTime: '2023-11-15 08:25:45',
  },
  {
    id: '1003',
    ruleName: '节点离线告警',
    target: 'node-04',
    description: '计算节点离线',
    severity: '紧急',
    time: '2023-11-14 23:45:12',
    duration: '3小时25分钟',
    status: '未处理',
    handler: '-',
    handleTime: '-',
  },
  {
    id: '1004',
    ruleName: '训练任务失败',
    target: '任务ID:5',
    description: '训练任务执行失败',
    severity: '中等',
    time: '2023-11-14 19:30:45',
    duration: '-',
    status: '未处理',
    handler: '-',
    handleTime: '-',
  },
]);

// 告警规则表格列
const ruleColumns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '规则名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
  },
  {
    title: '监控指标',
    dataIndex: 'metric',
    key: 'metric',
  },
  {
    title: '阈值',
    dataIndex: 'threshold',
    key: 'threshold',
  },
  {
    title: '持续时间',
    dataIndex: 'duration',
    key: 'duration',
  },
  {
    title: '严重程度',
    dataIndex: 'severity',
    key: 'severity',
    customRender: ({ text }: { text: string }) => {
      const color =
        text === '紧急'
          ? 'red'
          : text === '严重'
            ? 'orange'
            : text === '中等'
              ? 'blue'
              : 'green';
      return h(Tag, { color }, () => text);
    },
  },
  {
    title: '通知渠道',
    dataIndex: 'notifyChannel',
    key: 'notifyChannel',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }: { text: string }) => {
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
    title: '操作',
    key: 'action',
    customRender: ({ record }: { record: AlertRule }) =>
      h(Space, { size: 'middle' }, [
        h('a', { onClick: () => handleEditRule(record) }, '编辑'),
        h(Switch, {
          checked: record.status === '启用',
          onChange: (checked: CheckedType) =>
            handleToggleStatus(record, checked as boolean),
        }),
        h(Popconfirm, {
          title: '确定要删除这条规则吗?',
          onConfirm: () => handleDeleteRule(record),
          okText: '是',
          cancelText: '否',
          children: [h('a', {}, '删除')],
        }),
      ]),
  },
];

// 告警历史表格列
const historyColumns = [
  {
    title: '告警ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '告警规则',
    dataIndex: 'ruleName',
    key: 'ruleName',
  },
  {
    title: '告警对象',
    dataIndex: 'target',
    key: 'target',
  },
  {
    title: '告警描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '严重程度',
    dataIndex: 'severity',
    key: 'severity',
    customRender: ({ text }: { text: string }) => {
      const color =
        text === '紧急'
          ? 'red'
          : text === '严重'
            ? 'orange'
            : text === '中等'
              ? 'blue'
              : 'green';
      return h(Tag, { color }, () => text);
    },
  },
  {
    title: '告警时间',
    dataIndex: 'time',
    key: 'time',
  },
  {
    title: '持续时间',
    dataIndex: 'duration',
    key: 'duration',
  },
  {
    title: '处理状态',
    dataIndex: 'status',
    key: 'status',
    customRender: ({ text }: { text: string }) => {
      const color = text === '已处理' ? 'green' : 'orange';
      return h(Tag, { color }, () => text);
    },
  },
  {
    title: '处理人',
    dataIndex: 'handler',
    key: 'handler',
  },
  {
    title: '处理时间',
    dataIndex: 'handleTime',
    key: 'handleTime',
  },
  {
    title: '操作',
    key: 'action',
    customRender: ({ record }: { record: AlertHistory }) =>
      h(
        Space,
        { size: 'middle' },
        [
          h('a', {}, '详情'),
          record.status === '未处理'
            ? h('a', { onClick: () => handleAlert(record) }, '处理')
            : null,
        ].filter(Boolean),
      ),
  },
];

const activeKey = ref('1');
const ruleModalVisible = ref(false);
const editingRule = ref<AlertRule | null>(null);
const formRef = ref();
const formState = ref<FormState>({
  name: '',
  description: '',
  type: '',
  metric: '',
  operator: '>',
  threshold: 90,
  duration: 5,
  severity: '中等',
  notifyMail: true,
  notifySms: false,
  notifyPhone: false,
  status: true,
});

const severityOptions: Option[] = [
  { label: '紧急', value: '紧急' },
  { label: '严重', value: '严重' },
  { label: '中等', value: '中等' },
  { label: '一般', value: '一般' },
];

const typeOptions: Option[] = [
  { label: '集群资源', value: '集群资源' },
  { label: '任务资源', value: '任务资源' },
  { label: '节点状态', value: '节点状态' },
  { label: '任务状态', value: '任务状态' },
];

const metricOptions: Option[] = [
  { label: 'CPU使用率', value: 'CPU使用率' },
  { label: '内存使用率', value: '内存使用率' },
  { label: 'GPU使用率', value: 'GPU使用率' },
  { label: 'GPU显存使用率', value: 'GPU显存使用率' },
  { label: '节点状态', value: '节点状态' },
  { label: '任务状态', value: '任务状态' },
];

const handleEditRule = (rule: AlertRule) => {
  editingRule.value = rule;
  formState.value = {
    name: rule.name,
    description: rule.description,
    type: rule.type,
    metric: rule.metric,
    operator: rule.threshold.split(' ')[0],
    threshold: parseInt(rule.threshold.split(' ')[1]) || 0,
    duration: parseInt(rule.duration) || 5,
    severity: rule.severity,
    notifyMail: rule.notifyChannel.includes('邮件'),
    notifySms: rule.notifyChannel.includes('短信'),
    notifyPhone: rule.notifyChannel.includes('电话'),
    status: rule.status === '启用',
  };
  ruleModalVisible.value = true;
};

const showAddRuleModal = () => {
  editingRule.value = null;
  formState.value = {
    name: '',
    description: '',
    type: '',
    metric: '',
    operator: '>',
    threshold: 90,
    duration: 5,
    severity: '中等',
    notifyMail: true,
    notifySms: false,
    notifyPhone: false,
    status: true,
  };
  ruleModalVisible.value = true;
};

const handleCancel = () => {
  ruleModalVisible.value = false;
};

const handleSaveRule = () => {
  formRef.value.validate().then(() => {
    // 构建通知渠道字符串
    let channels: string[] = [];
    if (formState.value.notifyMail) channels.push('邮件');
    if (formState.value.notifySms) channels.push('短信');
    if (formState.value.notifyPhone) channels.push('电话');

    const rule: AlertRule = {
      id: editingRule.value
        ? editingRule.value.id
        : String(alertRules.value.length + 1),
      name: formState.value.name,
      description: formState.value.description,
      type: formState.value.type,
      metric: formState.value.metric,
      threshold: `${formState.value.operator} ${formState.value.threshold}%`,
      duration: `${formState.value.duration}分钟`,
      severity: formState.value.severity,
      notifyChannel: channels.join('、'),
      status: formState.value.status ? '启用' : '禁用',
      createTime: new Date().toLocaleString(),
    };

    if (editingRule.value) {
      // 更新现有规则
      const index = alertRules.value.findIndex(
        (r) => r.id === editingRule.value!.id,
      );
      alertRules.value[index] = { ...alertRules.value[index], ...rule };
    } else {
      // 添加新规则
      alertRules.value.push(rule);
    }

    ruleModalVisible.value = false;
  });
};

const handleToggleStatus = (record: AlertRule, checked: boolean) => {
  const index = alertRules.value.findIndex((r) => r.id === record.id);
  alertRules.value[index].status = checked ? '启用' : '禁用';
};

const handleDeleteRule = (record: AlertRule) => {
  const index = alertRules.value.findIndex((r) => r.id === record.id);
  alertRules.value.splice(index, 1);
};

const handleAlert = (record: AlertHistory) => {
  const index = alertHistory.value.findIndex((a) => a.id === record.id);
  alertHistory.value[index].status = '已处理';
  alertHistory.value[index].handler = '当前用户';
  alertHistory.value[index].handleTime = new Date().toLocaleString();
};
</script>

<template>
  <div class="alert-management-container">
    <Tabs v-model:activeKey="activeKey">
      <Tabs.TabPane key="1" tab="告警规则">
        <Card title="告警规则管理">
          <template #extra>
            <Button type="primary" @click="showAddRuleModal">
              新增告警规则
            </Button>
          </template>
          <Table :columns="ruleColumns" :dataSource="alertRules" rowKey="id" />
        </Card>
      </Tabs.TabPane>
      <Tabs.TabPane key="2" tab="告警历史">
        <Card title="告警历史记录">
          <Table
            :columns="historyColumns"
            :dataSource="alertHistory"
            rowKey="id"
          />
        </Card>
      </Tabs.TabPane>
    </Tabs>

    <!-- 告警规则表单模态框 -->
    <Modal
      v-model:visible="ruleModalVisible"
      :title="editingRule ? '编辑告警规则' : '新增告警规则'"
      @ok="handleSaveRule"
      @cancel="handleCancel"
      width="700px"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <Form.Item
          name="name"
          label="规则名称"
          :rules="[{ required: true, message: '请输入规则名称' }]"
        >
          <Input v-model:value="formState.name" placeholder="请输入规则名称" />
        </Form.Item>
        <Form.Item name="description" label="规则描述">
          <Input.TextArea
            v-model:value="formState.description"
            placeholder="请输入规则描述"
          />
        </Form.Item>
        <Form.Item
          name="type"
          label="规则类型"
          :rules="[{ required: true, message: '请选择规则类型' }]"
        >
          <Select v-model:value="formState.type" placeholder="请选择规则类型">
            <Select.Option
              v-for="option in typeOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          name="metric"
          label="监控指标"
          :rules="[{ required: true, message: '请选择监控指标' }]"
        >
          <Select v-model:value="formState.metric" placeholder="请选择监控指标">
            <Select.Option
              v-for="option in metricOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item label="告警条件" required>
          <Input.Group compact>
            <Select v-model:value="formState.operator" style="width: 15%">
              <Select.Option value=">">></Select.Option>
              <Select.Option value="<"><</Select.Option>
              <Select.Option value="==">==</Select.Option>
              <Select.Option value=">=">>=</Select.Option>
              <Select.Option value="<="><=</Select.Option>
            </Select>
            <InputNumber
              v-model:value="formState.threshold"
              style="width: 35%"
              :min="0"
              :max="100"
              addonAfter="%"
            />
            <span style="padding: 0 8px; line-height: 32px">持续</span>
            <InputNumber
              v-model:value="formState.duration"
              style="width: 25%"
              :min="0"
              addonAfter="分钟"
            />
          </Input.Group>
        </Form.Item>
        <Form.Item
          name="severity"
          label="严重程度"
          :rules="[{ required: true, message: '请选择严重程度' }]"
        >
          <Select
            v-model:value="formState.severity"
            placeholder="请选择严重程度"
          >
            <Select.Option
              v-for="option in severityOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </Select.Option>
          </Select>
        </Form.Item>
        <Form.Item label="通知渠道">
          <Space>
            <Checkbox v-model:checked="formState.notifyMail">邮件</Checkbox>
            <Checkbox v-model:checked="formState.notifySms">短信</Checkbox>
            <Checkbox v-model:checked="formState.notifyPhone">电话</Checkbox>
          </Space>
        </Form.Item>
        <Form.Item name="status" label="状态">
          <Switch v-model:checked="formState.status" />
          <span style="margin-left: 8px">
            {{ formState.status ? '启用' : '禁用' }}
          </span>
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.alert-management-container {
  padding: 0;
}
</style>
