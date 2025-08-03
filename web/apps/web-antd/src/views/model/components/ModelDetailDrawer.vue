<template>
  <Drawer
    v-model:open="visible"
    title="模型详情"
    width="800"
    placement="right"
    class="model-detail-drawer"
  >
    <div v-if="model" class="drawer-content">
      <!-- 模型基本信息 -->
      <div class="model-header">
        <div class="model-title">
          <h3>{{ model.name }}</h3>
          <div class="title-tags">
            <Tag v-if="model.isPublic" color="blue">公开</Tag>
            <Tag :color="getModelStatusColor(model.status)">
              {{ getModelStatusText(model.status) }}
            </Tag>
            <Tag :color="getModelTypeColor(model.type)">
              {{ getModelTypeText(model.type) }}
            </Tag>
          </div>
        </div>
        
        <div class="model-actions">
          <Space>
            <Button @click="downloadModel" :loading="downloadLoading">
              <DownloadOutlined />
              下载模型
            </Button>
            <Button type="primary" @click="deployModel" :disabled="model.status !== ModelStatus.AVAILABLE">
              <RocketOutlined />
              部署模型
            </Button>
          </Space>
        </div>
      </div>

      <Divider style="margin: 16px 0" />

      <!-- 模型信息标签页 -->
      <Tabs v-model:activeKey="activeTab">
        <Tabs.TabPane key="overview" tab="概览">
          <div class="overview-content">
            <!-- 基础信息 -->
            <Card title="基础信息" class="info-card">
              <Descriptions :column="2" bordered>
                <Descriptions.Item label="模型名称">
                  {{ model.name }}
                </Descriptions.Item>
                <Descriptions.Item label="版本">
                  {{ model.version }}
                </Descriptions.Item>
                <Descriptions.Item label="类型">
                  <Tag :color="getModelTypeColor(model.type)">
                    {{ getModelTypeText(model.type) }}
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="框架">
                  {{ model.framework }}
                </Descriptions.Item>
                <Descriptions.Item label="大小">
                  {{ formatFileSize(model.size) }}
                </Descriptions.Item>
                <Descriptions.Item label="准确率">
                  <Progress
                    :percent="model.accuracy * 100"
                    size="small"
                    :stroke-color="model.accuracy >= 0.9 ? '#52c41a' : model.accuracy >= 0.8 ? '#faad14' : '#ff4d4f'"
                    style="width: 120px"
                  />
                  {{ (model.accuracy * 100).toFixed(1) }}%
                </Descriptions.Item>
                <Descriptions.Item label="状态">
                  <Badge 
                    :status="getModelStatusColor(model.status) as any" 
                    :text="getModelStatusText(model.status)"
                  />
                </Descriptions.Item>
                <Descriptions.Item label="可见性">
                  <Tag :color="model.isPublic ? 'blue' : 'default'">
                    {{ model.isPublic ? '公开' : '私有' }}
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="下载次数" :span="2">
                  <Statistic
                    :value="model.downloadCount"
                    suffix="次"
                    :value-style="{ fontSize: '14px' }"
                  />
                </Descriptions.Item>
              </Descriptions>
            </Card>

            <!-- 创建信息 -->
            <Card title="创建信息" class="info-card">
              <Descriptions :column="2" bordered>
                <Descriptions.Item label="创建者">
                  <div class="creator-info">
                    <Avatar size="small">{{ model.creatorName?.[0] }}</Avatar>
                    <span style="margin-left: 8px">{{ model.creatorName }}</span>
                  </div>
                </Descriptions.Item>
                <Descriptions.Item label="工作空间">
                  {{ model.workspaceName }}
                </Descriptions.Item>
                <Descriptions.Item label="项目" v-if="model.projectName">
                  {{ model.projectName }}
                </Descriptions.Item>
                <Descriptions.Item label="训练任务" v-if="model.trainedJobId">
                  <Button type="link" size="small" @click="viewTrainingJob">
                    {{ model.trainedJobId }}
                  </Button>
                </Descriptions.Item>
                <Descriptions.Item label="训练数据集" v-if="model.trainingDataset">
                  {{ model.trainingDataset }}
                </Descriptions.Item>
                <Descriptions.Item label="创建时间">
                  {{ formatDateTime(model.createTime) }}
                </Descriptions.Item>
                <Descriptions.Item label="更新时间">
                  {{ formatDateTime(model.updateTime) }}
                </Descriptions.Item>
              </Descriptions>
            </Card>

            <!-- 模型描述 -->
            <Card title="描述" class="info-card" v-if="model.description">
              <div class="description-content">
                {{ model.description }}
              </div>
            </Card>

            <!-- 标签 -->
            <Card title="标签" class="info-card" v-if="model.tags?.length">
              <div class="tags-content">
                <Tag v-for="tag in model.tags" :key="tag" style="margin-bottom: 8px">
                  {{ tag }}
                </Tag>
              </div>
            </Card>
          </div>
        </Tabs.TabPane>

        <Tabs.TabPane key="hyperparameters" tab="超参数">
          <Card title="超参数配置" class="info-card">
            <div v-if="model.hyperParameters && Object.keys(model.hyperParameters).length > 0">
              <Table
                :columns="hyperParamsColumns"
                :data-source="hyperParamsData"
                :pagination="false"
                size="small"
              >
                <template #value="{ record }">
                  <code class="param-value">{{ formatParamValue(record.value) }}</code>
                </template>
              </Table>
            </div>
            <Empty v-else description="暂无超参数信息" />
          </Card>
        </Tabs.TabPane>

        <Tabs.TabPane key="metrics" tab="性能指标">
          <Card title="模型性能" class="info-card">
            <Row :gutter="16">
              <Col :span="8">
                <Statistic
                  title="准确率"
                  :value="model.accuracy * 100"
                  suffix="%"
                  :precision="2"
                  :value-style="{ color: '#52c41a' }"
                />
                <Progress
                  :percent="model.accuracy * 100"
                  :stroke-color="model.accuracy >= 0.9 ? '#52c41a' : model.accuracy >= 0.8 ? '#faad14' : '#ff4d4f'"
                  style="margin-top: 8px"
                />
              </Col>
              <Col :span="8">
                <Statistic
                  title="模型大小"
                  :value="formatFileSize(model.size)"
                  :value-style="{ color: '#1890ff' }"
                />
              </Col>
              <Col :span="8">
                <Statistic
                  title="使用次数"
                  :value="model.useCount"
                  suffix="次"
                  :value-style="{ color: '#722ed1' }"
                />
              </Col>
            </Row>

            <!-- 性能图表占位 -->
            <div class="metrics-charts" style="margin-top: 24px">
              <Alert
                message="性能指标图表"
                description="这里可以展示模型的详细性能指标图表，如ROC曲线、混淆矩阵等。"
                type="info"
                show-icon
              />
            </div>
          </Card>
        </Tabs.TabPane>

        <Tabs.TabPane key="usage" tab="使用记录">
          <Card title="使用历史" class="info-card">
            <div class="usage-stats">
              <Row :gutter="16" style="margin-bottom: 16px">
                <Col :span="8">
                  <Statistic
                    title="总下载量"
                    :value="model.downloadCount"
                    :value-style="{ color: '#52c41a' }"
                    prefix="📥"
                  />
                </Col>
                <Col :span="8">
                  <Statistic
                    title="总使用次数"
                    :value="model.useCount"
                    :value-style="{ color: '#1890ff' }"
                    prefix="🚀"
                  />
                </Col>
                <Col :span="8">
                  <Statistic
                    title="活跃用户"
                    value="15"
                    :value-style="{ color: '#722ed1' }"
                    prefix="👥"
                  />
                </Col>
              </Row>
            </div>

            <!-- 使用记录表格占位 -->
            <Alert
              message="使用记录详情"
              description="这里可以展示模型的详细使用记录，包括使用时间、使用者、使用场景等信息。"
              type="info"
              show-icon
              style="margin-top: 16px"
            />
          </Card>
        </Tabs.TabPane>

        <Tabs.TabPane key="files" tab="文件信息">
          <Card title="模型文件" class="info-card">
            <Descriptions :column="1" bordered>
              <Descriptions.Item label="文件路径">
                <code>{{ model.filePath }}</code>
              </Descriptions.Item>
              <Descriptions.Item label="下载地址" v-if="model.downloadUrl">
                <Button type="link" size="small" @click="copyDownloadUrl">
                  <CopyOutlined />
                  复制下载链接
                </Button>
              </Descriptions.Item>
              <Descriptions.Item label="文件大小">
                {{ formatFileSize(model.size) }}
              </Descriptions.Item>
              <Descriptions.Item label="文件格式">
                {{ getFileExtension(model.filePath) }}
              </Descriptions.Item>
            </Descriptions>

            <Divider />

            <!-- 文件操作 -->
            <div class="file-actions">
              <Space>
                <Button @click="downloadModel" :loading="downloadLoading">
                  <DownloadOutlined />
                  下载模型
                </Button>
                <Button @click="previewModel">
                  <EyeOutlined />
                  预览结构
                </Button>
                <Button @click="validateModel">
                  <SafetyCertificateOutlined />
                  验证模型
                </Button>
              </Space>
            </div>
          </Card>
        </Tabs.TabPane>
      </Tabs>
    </div>
  </Drawer>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue';
import {
  Drawer,
  Tag,
  Space,
  Button,
  Divider,
  Tabs,
  Card,
  Descriptions,
  Progress,
  Badge,
  Statistic,
  Avatar,
  Row,
  Col,
  Table,
  Empty,
  Alert,
  message,
} from 'ant-design-vue';
import {
  DownloadOutlined,
  RocketOutlined,
  CopyOutlined,
  EyeOutlined,
  SafetyCertificateOutlined,
} from '@ant-design/icons-vue';
import type { Model, ModelType, ModelStatus } from '#/api/types';
import { formatDateTime, formatFileSize } from '#/utils/date';

const props = defineProps<{
  visible: boolean;
  model: Model | null;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
  'deploy': [model: Model];
}>();

// 响应式数据
const activeTab = ref('overview');
const downloadLoading = ref(false);

// 计算属性
const hyperParamsData = computed(() => {
  if (!props.model?.hyperParameters) return [];
  
  return Object.entries(props.model.hyperParameters).map(([key, value]) => ({
    key,
    parameter: key,
    value,
  }));
});

// 超参数表格列定义
const hyperParamsColumns = [
  {
    title: '参数名',
    dataIndex: 'parameter',
    key: 'parameter',
    width: 200,
  },
  {
    title: '参数值',
    key: 'value',
    slots: { customRender: 'value' },
  },
];

// 工具方法
const getModelTypeText = (type: ModelType) => {
  const types = {
    [ModelType.CLASSIFICATION]: '分类',
    [ModelType.REGRESSION]: '回归',
    [ModelType.OBJECT_DETECTION]: '目标检测',
    [ModelType.SEMANTIC_SEGMENTATION]: '语义分割',
    [ModelType.NLP]: '自然语言处理',
    [ModelType.RECOMMENDATION]: '推荐系统',
    [ModelType.GENERATIVE]: '生成模型',
    [ModelType.CUSTOM]: '自定义',
  };
  return types[type] || type;
};

const getModelTypeColor = (type: ModelType) => {
  const colors = {
    [ModelType.CLASSIFICATION]: 'blue',
    [ModelType.REGRESSION]: 'green',
    [ModelType.OBJECT_DETECTION]: 'orange',
    [ModelType.SEMANTIC_SEGMENTATION]: 'purple',
    [ModelType.NLP]: 'cyan',
    [ModelType.RECOMMENDATION]: 'magenta',
    [ModelType.GENERATIVE]: 'red',
    [ModelType.CUSTOM]: 'default',
  };
  return colors[type] || 'default';
};

const getModelStatusText = (status: ModelStatus) => {
  const statuses = {
    [ModelStatus.TRAINING]: '训练中',
    [ModelStatus.AVAILABLE]: '可用',
    [ModelStatus.DEPRECATED]: '已弃用',
    [ModelStatus.DELETED]: '已删除',
  };
  return statuses[status] || status;
};

const getModelStatusColor = (status: ModelStatus) => {
  const colors = {
    [ModelStatus.TRAINING]: 'processing',
    [ModelStatus.AVAILABLE]: 'success',
    [ModelStatus.DEPRECATED]: 'warning',
    [ModelStatus.DELETED]: 'error',
  };
  return colors[status] || 'default';
};

const formatParamValue = (value: any) => {
  if (typeof value === 'object') {
    return JSON.stringify(value, null, 2);
  }
  return String(value);
};

const getFileExtension = (filePath: string) => {
  const ext = filePath.split('.').pop();
  return ext ? `.${ext}` : '未知';
};

// 事件处理
const downloadModel = async () => {
  if (!props.model) return;
  
  try {
    downloadLoading.value = true;
    // 实际应该调用下载API
    message.success('模型下载中...');
  } catch (error) {
    message.error('下载失败');
  } finally {
    downloadLoading.value = false;
  }
};

const deployModel = () => {
  if (props.model) {
    emit('deploy', props.model);
  }
};

const viewTrainingJob = () => {
  message.info('跳转到训练任务详情');
};

const copyDownloadUrl = async () => {
  if (props.model?.downloadUrl) {
    try {
      await navigator.clipboard.writeText(props.model.downloadUrl);
      message.success('下载链接已复制到剪贴板');
    } catch (error) {
      message.error('复制失败');
    }
  }
};

const previewModel = () => {
  message.info('模型结构预览功能开发中');
};

const validateModel = () => {
  message.info('模型验证功能开发中');
};

// 监听visible变化，重置activeTab
watch(() => props.visible, (newVal) => {
  if (newVal) {
    activeTab.value = 'overview';
  }
});
</script>

<style scoped lang="scss">
.model-detail-drawer {
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

.model-header {
  flex-shrink: 0;
  
  .model-title {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 16px;
    
    h3 {
      margin: 0;
      color: #1890ff;
      flex: 1;
    }
    
    .title-tags {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }
  }
  
  .model-actions {
    display: flex;
    justify-content: flex-end;
  }
}

.overview-content {
  .info-card {
    margin-bottom: 16px;
    
    :deep(.ant-card-head) {
      background: #fafafa;
    }
    
    :deep(.ant-card-head-title) {
      font-weight: 600;
      color: #1890ff;
      font-size: 14px;
    }
  }
}

.creator-info {
  display: flex;
  align-items: center;
}

.description-content {
  color: #666;
  line-height: 1.6;
  white-space: pre-wrap;
}

.tags-content {
  .ant-tag {
    margin-right: 8px;
  }
}

.param-value {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 12px;
}

.metrics-charts {
  min-height: 200px;
}

.usage-stats {
  .ant-statistic {
    text-align: center;
  }
}

.file-actions {
  margin-top: 16px;
}

// 响应式设计
@media (max-width: 768px) {
  .drawer-content {
    padding: 16px;
  }
  
  .model-title {
    flex-direction: column !important;
    gap: 12px;
    
    .title-tags {
      justify-content: flex-start;
    }
  }
  
  .model-actions {
    justify-content: flex-start !important;
    
    :deep(.ant-space) {
      flex-wrap: wrap;
    }
  }
  
  .overview-content {
    :deep(.ant-descriptions) {
      .ant-descriptions-item {
        padding: 8px 12px;
      }
    }
  }
}
</style>