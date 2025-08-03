<template>
  <Drawer
    v-model:open="visible"
    title="版本管理"
    width="900"
    placement="right"
    class="model-version-drawer"
  >
    <div v-if="model" class="drawer-content">
      <!-- 版本管理头部 -->
      <div class="version-header">
        <div class="model-info">
          <h3>{{ model.name }}</h3>
          <Tag color="blue">当前版本: {{ model.version }}</Tag>
        </div>
        
        <div class="header-actions">
          <Space>
            <Button @click="refreshVersions" :loading="loading">
              <ReloadOutlined />
              刷新
            </Button>
            <Button type="primary" @click="showCreateVersionModal">
              <PlusOutlined />
              创建新版本
            </Button>
          </Space>
        </div>
      </div>

      <Divider style="margin: 16px 0" />

      <!-- 版本统计 -->
      <Row :gutter="16" class="version-stats">
        <Col :span="6">
          <Card>
            <Statistic
              title="总版本数"
              :value="versionList.length"
              :value-style="{ color: '#3f8600' }"
              prefix="📦"
            />
          </Card>
        </Col>
        <Col :span="6">
          <Card>
            <Statistic
              title="活跃版本"
              :value="activeVersionCount"
              :value-style="{ color: '#52c41a' }"
              prefix="✅"
            />
          </Card>
        </Col>
        <Col :span="6">
          <Card>
            <Statistic
              title="总下载量"
              :value="totalDownloads"
              :value-style="{ color: '#1890ff' }"
              prefix="📥"
            />
          </Card>
        </Col>
        <Col :span="6">
          <Card>
            <Statistic
              title="总大小"
              :value="formatFileSize(totalSize)"
              :value-style="{ color: '#722ed1' }"
              prefix="💾"
            />
          </Card>
        </Col>
      </Row>

      <!-- 版本列表 -->
      <Card title="版本历史" class="version-list-card">
        <Table
          :columns="versionColumns"
          :data-source="versionList"
          :loading="loading"
          row-key="version"
          :pagination="{ pageSize: 10, size: 'small' }"
        >
          <!-- 版本信息 -->
          <template #versionInfo="{ record }">
            <div class="version-info">
              <div class="version-header">
                <span class="version-number">{{ record.version }}</span>
                <Tag v-if="record.version === model.version" color="green" size="small">
                  当前版本
                </Tag>
                <Tag v-if="record.isLatest" color="blue" size="small">
                  最新版本
                </Tag>
              </div>
              <div class="version-desc">{{ record.description || '无描述' }}</div>
              <div class="version-changelog" v-if="record.changeLog">
                <Text type="secondary" style="font-size: 12px">
                  {{ record.changeLog }}
                </Text>
              </div>
            </div>
          </template>

          <!-- 性能指标 -->
          <template #metrics="{ record }">
            <div class="metrics-info">
              <div class="metric-item" v-if="record.accuracy">
                <span class="metric-label">准确率:</span>
                <Progress
                  :percent="record.accuracy * 100"
                  size="small"
                  :stroke-color="record.accuracy >= 0.9 ? '#52c41a' : record.accuracy >= 0.8 ? '#faad14' : '#ff4d4f'"
                  style="width: 80px; margin-left: 8px"
                />
                <span class="metric-value">{{ (record.accuracy * 100).toFixed(1) }}%</span>
              </div>
              <div class="metric-item">
                <span class="metric-label">大小:</span>
                <span class="metric-value">{{ formatFileSize(record.size) }}</span>
              </div>
            </div>
          </template>

          <!-- 统计信息 -->
          <template #stats="{ record }">
            <div class="stats-info">
              <div class="stat-item">
                <DownloadOutlined style="margin-right: 4px" />
                {{ record.downloadCount || 0 }}
              </div>
              <div class="stat-item">
                <EyeOutlined style="margin-right: 4px" />
                {{ record.useCount || 0 }}
              </div>
            </div>
          </template>

          <!-- 创建信息 -->
          <template #creator="{ record }">
            <div class="creator-info">
              <Avatar size="small">{{ record.creatorName?.[0] }}</Avatar>
              <span style="margin-left: 8px">{{ record.creatorName }}</span>
            </div>
          </template>

          <!-- 创建时间 -->
          <template #createTime="{ record }">
            <div class="time-info">
              <div>{{ formatDateTime(record.createTime, 'MM-DD') }}</div>
              <div class="time-detail">{{ formatDateTime(record.createTime, 'HH:mm') }}</div>
            </div>
          </template>

          <!-- 操作 -->
          <template #action="{ record }">
            <Space size="small">
              <Tooltip title="下载版本">
                <Button type="text" size="small" @click="downloadVersion(record)">
                  <DownloadOutlined />
                </Button>
              </Tooltip>
              <Tooltip title="查看详情">
                <Button type="text" size="small" @click="viewVersionDetail(record)">
                  <EyeOutlined />
                </Button>
              </Tooltip>
              <Tooltip title="比较版本">
                <Button type="text" size="small" @click="compareVersion(record)">
                  <DiffOutlined />
                </Button>
              </Tooltip>
              <Dropdown>
                <Button type="text" size="small">
                  <MoreOutlined />
                </Button>
                <template #overlay>
                  <Menu>
                    <Menu.Item 
                      @click="setAsActive(record)"
                      :disabled="record.version === model.version"
                    >
                      <CheckOutlined />
                      设为当前版本
                    </Menu.Item>
                    <Menu.Item @click="editVersion(record)">
                      <EditOutlined />
                      编辑信息
                    </Menu.Item>
                    <Menu.Item 
                      @click="deleteVersion(record)" 
                      danger
                      :disabled="record.version === model.version"
                    >
                      <DeleteOutlined />
                      删除版本
                    </Menu.Item>
                  </Menu>
                </template>
              </Dropdown>
            </Space>
          </template>
        </Table>
      </Card>
    </div>

    <!-- 创建新版本模态框 -->
    <Modal
      v-model:open="createVersionModalVisible"
      title="创建新版本"
      width="600px"
      @ok="handleCreateVersionSubmit"
      @cancel="handleCreateVersionCancel"
      :confirm-loading="createVersionLoading"
    >
      <Form
        ref="createVersionFormRef"
        :model="createVersionForm"
        :rules="createVersionFormRules"
        layout="vertical"
      >
        <Form.Item label="版本号" name="version">
          <Input v-model:value="createVersionForm.version" placeholder="例如: v1.1.0" />
        </Form.Item>

        <Form.Item label="版本描述" name="description">
          <Input v-model:value="createVersionForm.description" placeholder="简要描述这个版本的特点" />
        </Form.Item>

        <Form.Item label="更新日志" name="changeLog">
          <Input.TextArea
            v-model:value="createVersionForm.changeLog"
            placeholder="详细描述本版本的更新内容、修复的问题等"
            :rows="4"
          />
        </Form.Item>

        <Form.Item label="模型文件" required>
          <Upload
            :before-upload="beforeUpload"
            :file-list="createVersionForm.file ? [{ uid: '1', name: createVersionForm.file.name, status: 'done' }] : []"
            accept=".h5,.pkl,.pth,.pt,.onnx,.pb,.tflite,.tar.gz,.zip"
          >
            <Button>
              <UploadOutlined />
              选择文件
            </Button>
            <div style="margin-top: 8px; color: #999; font-size: 12px">
              支持格式: .h5, .pkl, .pth, .pt, .onnx, .pb, .tflite, .tar.gz, .zip
            </div>
          </Upload>
        </Form.Item>

        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="准确率" name="accuracy">
              <InputNumber
                v-model:value="createVersionForm.accuracy"
                :min="0"
                :max="1"
                :step="0.001"
                :precision="3"
                placeholder="0.000"
                style="width: 100%"
              />
            </Form.Item>
          </Col>
        </Row>

        <Form.Item label="超参数配置" name="hyperParameters">
          <Input.TextArea
            v-model:value="createVersionForm.hyperParametersText"
            placeholder='JSON格式，例如: {"learning_rate": 0.001, "batch_size": 32}'
            :rows="4"
          />
        </Form.Item>
      </Form>
    </Modal>

    <!-- 版本详情模态框 -->
    <Modal
      v-model:open="versionDetailModalVisible"
      title="版本详情"
      width="800px"
      :footer="null"
    >
      <div v-if="selectedVersion" class="version-detail">
        <Descriptions :column="2" bordered>
          <Descriptions.Item label="版本号">
            {{ selectedVersion.version }}
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag v-if="selectedVersion.version === model.version" color="green">
              当前版本
            </Tag>
            <Tag v-else color="default">
              历史版本
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="文件大小">
            {{ formatFileSize(selectedVersion.size) }}
          </Descriptions.Item>
          <Descriptions.Item label="准确率" v-if="selectedVersion.accuracy">
            <Progress
              :percent="selectedVersion.accuracy * 100"
              size="small"
              :stroke-color="selectedVersion.accuracy >= 0.9 ? '#52c41a' : selectedVersion.accuracy >= 0.8 ? '#faad14' : '#ff4d4f'"
              style="width: 120px"
            />
            {{ (selectedVersion.accuracy * 100).toFixed(1) }}%
          </Descriptions.Item>
          <Descriptions.Item label="下载次数">
            {{ selectedVersion.downloadCount || 0 }}
          </Descriptions.Item>
          <Descriptions.Item label="使用次数">
            {{ selectedVersion.useCount || 0 }}
          </Descriptions.Item>
          <Descriptions.Item label="创建者">
            <div class="creator-info">
              <Avatar size="small">{{ selectedVersion.creatorName?.[0] }}</Avatar>
              <span style="margin-left: 8px">{{ selectedVersion.creatorName }}</span>
            </div>
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">
            {{ formatDateTime(selectedVersion.createTime) }}
          </Descriptions.Item>
          <Descriptions.Item label="描述" :span="2">
            {{ selectedVersion.description || '无描述' }}
          </Descriptions.Item>
          <Descriptions.Item label="更新日志" :span="2">
            <div class="changelog-content">
              {{ selectedVersion.changeLog || '无更新日志' }}
            </div>
          </Descriptions.Item>
        </Descriptions>

        <!-- 超参数信息 -->
        <Card title="超参数配置" style="margin-top: 16px" v-if="selectedVersion.hyperParameters">
          <pre class="hyperparams-display">{{ JSON.stringify(selectedVersion.hyperParameters, null, 2) }}</pre>
        </Card>
      </div>
    </Modal>
  </Drawer>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, onMounted } from 'vue';
import {
  Drawer,
  Space,
  Button,
  Divider,
  Row,
  Col,
  Card,
  Statistic,
  Table,
  Tag,
  Progress,
  Avatar,
  Tooltip,
  Dropdown,
  Menu,
  Modal,
  Form,
  Input,
  InputNumber,
  Upload,
  Descriptions,
  Typography,
  message,
} from 'ant-design-vue';
import {
  ReloadOutlined,
  PlusOutlined,
  DownloadOutlined,
  EyeOutlined,
  DiffOutlined,
  MoreOutlined,
  CheckOutlined,
  EditOutlined,
  DeleteOutlined,
  UploadOutlined,
} from '@ant-design/icons-vue';
import type { Model } from '#/api/types';
import { formatDateTime, formatFileSize } from '#/utils/date';

const { Text } = Typography;

const props = defineProps<{
  visible: boolean;
  model: Model | null;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
}>();

// 响应式数据
const loading = ref(false);
const createVersionModalVisible = ref(false);
const versionDetailModalVisible = ref(false);
const createVersionLoading = ref(false);
const versionList = ref<any[]>([]);
const selectedVersion = ref<any>(null);

// 创建版本表单
interface CreateVersionForm {
  version: string;
  description: string;
  changeLog: string;
  accuracy: number | null;
  hyperParametersText: string;
  file: File | null;
}

const createVersionForm = reactive<CreateVersionForm>({
  version: '',
  description: '',
  changeLog: '',
  accuracy: null,
  hyperParametersText: '',
  file: null,
});

const createVersionFormRef = ref();

// 模拟版本数据
const mockVersions = [
  {
    version: 'v1.2.0',
    description: '改进准确率和性能',
    changeLog: '- 优化模型架构\n- 提升准确率到95%\n- 减少模型大小',
    size: 412000000,
    accuracy: 0.95,
    downloadCount: 1258,
    useCount: 89,
    isLatest: true,
    creatorId: 'user-001',
    creatorName: '张三',
    createTime: '2024-01-20 14:20:00',
    hyperParameters: {
      'learning_rate': 0.00005,
      'batch_size': 32,
      'epochs': 10,
    },
  },
  {
    version: 'v1.1.0',
    description: '修复训练稳定性问题',
    changeLog: '- 修复梯度爆炸问题\n- 改进数据预处理\n- 增加正则化',
    size: 420000000,
    accuracy: 0.92,
    downloadCount: 856,
    useCount: 67,
    isLatest: false,
    creatorId: 'user-001',
    creatorName: '张三',
    createTime: '2024-01-18 10:15:00',
    hyperParameters: {
      'learning_rate': 0.0001,
      'batch_size': 32,
      'epochs': 8,
    },
  },
  {
    version: 'v1.0.0',
    description: '初始版本',
    changeLog: '- 基础模型实现\n- 支持中文语料\n- 基本功能完整',
    size: 408000000,
    accuracy: 0.89,
    downloadCount: 234,
    useCount: 45,
    isLatest: false,
    creatorId: 'user-001',
    creatorName: '张三',
    createTime: '2024-01-15 10:30:00',
    hyperParameters: {
      'learning_rate': 0.0002,
      'batch_size': 16,
      'epochs': 5,
    },
  },
];

// 表格列定义
const versionColumns = [
  {
    title: '版本信息',
    key: 'versionInfo',
    slots: { customRender: 'versionInfo' },
    width: 200,
  },
  {
    title: '性能指标',
    key: 'metrics',
    slots: { customRender: 'metrics' },
    width: 180,
  },
  {
    title: '统计',
    key: 'stats',
    slots: { customRender: 'stats' },
    width: 100,
  },
  {
    title: '创建者',
    key: 'creator',
    slots: { customRender: 'creator' },
    width: 120,
  },
  {
    title: '创建时间',
    key: 'createTime',
    slots: { customRender: 'createTime' },
    width: 120,
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 150,
    fixed: 'right' as const,
  },
];

// 计算属性
const activeVersionCount = computed(() => {
  return versionList.value.length; // 简化处理，所有版本都算活跃
});

const totalDownloads = computed(() => {
  return versionList.value.reduce((sum, version) => sum + (version.downloadCount || 0), 0);
});

const totalSize = computed(() => {
  return versionList.value.reduce((sum, version) => sum + version.size, 0);
});

// 数据加载
const loadVersions = async () => {
  if (!props.model) return;
  
  try {
    loading.value = true;
    // const response = await getModelVersions(props.model.id);
    // versionList.value = response.data.items;
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    versionList.value = mockVersions;
  } catch (error) {
    message.error('加载版本列表失败');
  } finally {
    loading.value = false;
  }
};

const refreshVersions = () => {
  loadVersions();
};

// 事件处理
const showCreateVersionModal = () => {
  createVersionModalVisible.value = true;
  resetCreateVersionForm();
};

const resetCreateVersionForm = () => {
  Object.assign(createVersionForm, {
    version: '',
    description: '',
    changeLog: '',
    accuracy: null,
    hyperParametersText: '',
    file: null,
  });
};

const beforeUpload = (file: File) => {
  createVersionForm.file = file;
  return false; // 阻止自动上传
};

const handleCreateVersionSubmit = async () => {
  try {
    await createVersionFormRef.value?.validate();
    
    if (!createVersionForm.file) {
      message.error('请选择模型文件');
      return;
    }
    
    createVersionLoading.value = true;
    
    // 解析超参数
    let hyperParameters = null;
    if (createVersionForm.hyperParametersText.trim()) {
      try {
        hyperParameters = JSON.parse(createVersionForm.hyperParametersText);
      } catch (error) {
        message.error('超参数格式错误，请使用有效的JSON格式');
        return;
      }
    }
    
    // const request = {
    //   version: createVersionForm.version,
    //   description: createVersionForm.description,
    //   changeLog: createVersionForm.changeLog,
    //   file: createVersionForm.file,
    //   accuracy: createVersionForm.accuracy,
    //   hyperParameters,
    // };
    
    // const response = await createModelVersion(props.model!.id, request);
    
    // 模拟创建成功
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    message.success('新版本创建成功');
    createVersionModalVisible.value = false;
    loadVersions();
  } catch (error) {
    message.error('创建失败');
  } finally {
    createVersionLoading.value = false;
  }
};

const handleCreateVersionCancel = () => {
  createVersionModalVisible.value = false;
};

const downloadVersion = async (version: any) => {
  try {
    // await downloadModel(props.model!.id, version.version);
    message.success(`版本 ${version.version} 下载中...`);
  } catch (error) {
    message.error('下载失败');
  }
};

const viewVersionDetail = (version: any) => {
  selectedVersion.value = version;
  versionDetailModalVisible.value = true;
};

const compareVersion = (version: any) => {
  message.info('版本比较功能开发中');
};

const setAsActive = async (version: any) => {
  Modal.confirm({
    title: '确认切换版本',
    content: `确定要将版本 ${version.version} 设置为当前活跃版本吗？`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // await setActiveModelVersion(props.model!.id, version.version);
        message.success('版本切换成功');
        loadVersions();
      } catch (error) {
        message.error('版本切换失败');
      }
    },
  });
};

const editVersion = (version: any) => {
  message.info('编辑版本功能开发中');
};

const deleteVersion = async (version: any) => {
  Modal.confirm({
    title: '确认删除版本',
    content: `确定要删除版本 ${version.version} 吗？此操作不可恢复。`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // await deleteModelVersion(props.model!.id, version.version);
        message.success('版本删除成功');
        loadVersions();
      } catch (error) {
        message.error('版本删除失败');
      }
    },
  });
};

// 表单验证规则
const createVersionFormRules = {
  version: [
    { required: true, message: '请输入版本号', trigger: 'blur' },
    { pattern: /^v?\d+\.\d+\.\d+$/, message: '版本号格式不正确，例如: v1.0.0 或 1.0.0', trigger: 'blur' },
  ],
  changeLog: [
    { required: true, message: '请输入更新日志', trigger: 'blur' },
  ],
};

// 监听模型变化
import { watch } from 'vue';
watch(() => props.model, (newModel) => {
  if (newModel && props.visible) {
    loadVersions();
  }
});

// 监听visible变化
watch(() => props.visible, (newVal) => {
  if (newVal && props.model) {
    loadVersions();
  }
});

// 初始化
onMounted(() => {
  if (props.visible && props.model) {
    loadVersions();
  }
});
</script>

<style scoped lang="scss">
.model-version-drawer {
  :deep(.ant-drawer-body) {
    padding: 0;
    display: flex;
    flex-direction: column;
  }
}

.drawer-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 24px;
  overflow-y: auto;
}

.version-header {
  flex-shrink: 0;
  
  .model-info {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
    
    h3 {
      margin: 0;
      color: #1890ff;
    }
  }
  
  .header-actions {
    display: flex;
    justify-content: flex-end;
  }
}

.version-stats {
  margin-bottom: 24px;
  
  .ant-statistic {
    text-align: center;
  }
}

.version-list-card {
  flex: 1;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
  
  :deep(.ant-card-head-title) {
    font-weight: 600;
    color: #1890ff;
    font-size: 14px;
  }
}

.version-info {
  .version-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
    
    .version-number {
      font-weight: 600;
      color: #1890ff;
      font-family: 'Monaco', 'Consolas', monospace;
    }
  }
  
  .version-desc {
    color: #666;
    font-size: 12px;
    margin-bottom: 4px;
  }
  
  .version-changelog {
    font-size: 11px;
    max-width: 180px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.metrics-info {
  .metric-item {
    display: flex;
    align-items: center;
    margin-bottom: 4px;
    
    .metric-label {
      font-size: 12px;
      color: #666;
      min-width: 50px;
    }
    
    .metric-value {
      font-size: 12px;
      font-weight: 500;
      margin-left: 8px;
    }
  }
}

.stats-info {
  .stat-item {
    display: flex;
    align-items: center;
    font-size: 12px;
    margin-bottom: 4px;
    color: #666;
  }
}

.creator-info {
  display: flex;
  align-items: center;
  
  span {
    font-size: 12px;
  }
}

.time-info {
  .time-detail {
    font-size: 12px;
    color: #999;
    margin-top: 2px;
  }
}

.version-detail {
  .creator-info {
    display: flex;
    align-items: center;
  }
  
  .changelog-content {
    white-space: pre-wrap;
    color: #666;
    line-height: 1.5;
  }
}

.hyperparams-display {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 6px;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 12px;
  line-height: 1.4;
  overflow-x: auto;
}

// 响应式设计
@media (max-width: 768px) {
  .drawer-content {
    padding: 16px;
  }
  
  .version-header {
    .model-info {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;
    }
    
    .header-actions {
      justify-content: flex-start;
      margin-top: 12px;
    }
  }
  
  .version-stats {
    :deep(.ant-col) {
      margin-bottom: 12px;
    }
  }
}
</style>