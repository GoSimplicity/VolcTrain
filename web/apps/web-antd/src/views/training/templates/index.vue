<template>
  <div class="training-templates-container">
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <FileTextIcon class="title-icon" />
            <span class="title-text">任务模板</span>
          </h1>
          <p class="page-description">
            <span class="description-text">管理和使用训练任务模板，快速创建标准化训练任务</span>
          </p>
        </div>
        <div class="action-section">
          <Space>
            <Button @click="refreshTemplates">
              <RefreshCwIcon :class="{ 'animate-spin': loading }" />
              刷新
            </Button>
            <Button type="primary" @click="showCreateModal">
              <PlusCircleIcon />
              创建模板
            </Button>
          </Space>
        </div>
      </div>
    </div>

    <Card class="templates-grid-card">
      <Row :gutter="[16, 16]">
        <Col
          v-for="template in templatesList"
          :key="template.id"
          :span="8"
        >
          <Card class="template-card" hoverable>
            <template #cover>
              <div class="template-cover">
                <div class="framework-badge">
                  <Tag :color="getFrameworkColor(template.framework)">
                    {{ getFrameworkText(template.framework) }}
                  </Tag>
                </div>
              </div>
            </template>
            
            <div class="template-content">
              <h3 class="template-name">{{ template.name }}</h3>
              <p class="template-description">{{ template.description }}</p>
              
              <div class="template-stats">
                <Space>
                  <span class="stat-item">
                    <StarIcon class="stat-icon" />
                    {{ template.stars || 0 }}
                  </span>
                  <span class="stat-item">
                    <PlayIcon class="stat-icon" />
                    {{ template.usageCount || 0 }}
                  </span>
                </Space>
              </div>

              <div class="template-actions">
                <Space>
                  <Button size="small" @click="viewTemplate(template.id)">
                    <EyeIcon />
                    查看
                  </Button>
                  <Button size="small" type="primary" @click="useTemplate(template.id)">
                    <PlayIcon />
                    使用
                  </Button>
                  <Dropdown>
                    <Button size="small">
                      <MoreHorizontalIcon />
                    </Button>
                    <template #overlay>
                      <Menu @click="(e) => handleTemplateAction(template, e.key)">
                        <Menu.Item key="edit">
                          <EditIcon />
                          编辑
                        </Menu.Item>
                        <Menu.Item key="clone">
                          <CopyIcon />
                          克隆
                        </Menu.Item>
                        <Menu.Item key="delete" danger>
                          <TrashIcon />
                          删除
                        </Menu.Item>
                      </Menu>
                    </template>
                  </Dropdown>
                </Space>
              </div>

              <div class="template-meta">
                <Text type="secondary">
                  创建于 {{ formatDateTime(template.createdAt) }}
                </Text>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </Card>

    <!-- 创建模板弹窗 -->
    <Modal
      v-model:open="createModalVisible"
      title="创建模板"
      width="800px"
      :confirm-loading="createLoading"
      @ok="handleCreateTemplate"
      @cancel="cancelCreate"
    >
      <CreateTemplateForm ref="createTemplateFormRef" />
    </Modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import {
  Card, Button, Space, Row, Col, Tag, Modal, Menu, Dropdown, Typography, message,
} from 'ant-design-vue';
import { format } from 'date-fns';
import {
  FileText as FileTextIcon,
  RefreshCw as RefreshCwIcon,
  PlusCircle as PlusCircleIcon,
  Star as StarIcon,
  Play as PlayIcon,
  Eye as EyeIcon,
  MoreHorizontal as MoreHorizontalIcon,
  Edit as EditIcon,
  Copy as CopyIcon,
  Trash as TrashIcon,
} from 'lucide-vue-next';

import CreateTemplateForm from './components/CreateTemplateForm.vue';

const { Text } = Typography;

defineOptions({ name: 'TrainingTemplates' });

const loading = ref(false);
const createLoading = ref(false);
const createModalVisible = ref(false);
const createTemplateFormRef = ref();

// 模板列表数据
const templatesList = ref([
  {
    id: '1',
    name: 'ResNet-50 图像分类',
    description: '基于ResNet-50的图像分类模板，适用于ImageNet等数据集',
    framework: 'pytorch',
    stars: 25,
    usageCount: 156,
    createdAt: '2024-01-15T10:00:00Z',
  },
  {
    id: '2', 
    name: 'BERT 文本分类',
    description: 'BERT模型文本分类模板，支持多种NLP任务',
    framework: 'tensorflow',
    stars: 18,
    usageCount: 89,
    createdAt: '2024-01-16T09:00:00Z',
  },
]);

const getFrameworkColor = (framework: string) => ({
  tensorflow: 'orange',
  pytorch: 'red',
  paddlepaddle: 'blue',
  mindspore: 'green',
}[framework] || 'default');

const getFrameworkText = (framework: string) => ({
  tensorflow: 'TensorFlow',
  pytorch: 'PyTorch', 
  paddlepaddle: 'PaddlePaddle',
  mindspore: 'MindSpore',
}[framework] || framework);

const formatDateTime = (dateTime: string) => format(new Date(dateTime), 'yyyy-MM-dd');

const refreshTemplates = () => loading.value = !loading.value;
const showCreateModal = () => createModalVisible.value = true;
const cancelCreate = () => createModalVisible.value = false;
const viewTemplate = (id: string) => message.info(`查看模板: ${id}`);
const useTemplate = (id: string) => message.success(`使用模板: ${id}`);

const handleTemplateAction = (template: any, action: string) => {
  message.info(`${action}: ${template.name}`);
};

const handleCreateTemplate = async () => {
  createLoading.value = true;
  setTimeout(() => {
    createLoading.value = false;
    createModalVisible.value = false;
    message.success('模板创建成功');
  }, 1000);
};

onMounted(() => {
  // 初始化
});
</script>

<style lang="scss" scoped>
.training-templates-container {
  padding: 0;
  background: #f5f7fa;
  min-height: 100vh;
}

.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 24px 32px;
  margin-bottom: 24px;

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .title-section .page-title {
    display: flex;
    align-items: center;
    font-size: 24px;
    font-weight: 600;
    margin: 0;
    
    .title-icon {
      margin-right: 12px;
      font-size: 28px;
    }
  }

  .title-section .page-description {
    margin: 8px 0 0 40px;
    font-size: 14px;
    opacity: 0.9;
  }
}

.templates-grid-card {
  margin: 0 24px;
  border-radius: 8px;
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.template-card {
  height: 100%;

  .template-cover {
    height: 60px;
    background: linear-gradient(45deg, #f0f2f5, #e6f7ff);
    display: flex;
    align-items: flex-end;
    justify-content: flex-end;
    padding: 8px;

    .framework-badge {
      position: absolute;
      top: 8px;
      right: 8px;
    }
  }

  .template-content {
    padding: 16px;

    .template-name {
      margin-bottom: 8px;
      font-size: 16px;
      font-weight: 600;
    }

    .template-description {
      color: #666;
      margin-bottom: 12px;
      font-size: 13px;
      line-height: 1.4;
    }

    .template-stats {
      margin-bottom: 16px;

      .stat-item {
        display: inline-flex;
        align-items: center;
        color: #666;
        font-size: 12px;

        .stat-icon {
          width: 14px;
          height: 14px;
          margin-right: 4px;
        }
      }
    }

    .template-actions {
      margin-bottom: 12px;
    }

    .template-meta {
      border-top: 1px solid #f0f0f0;
      padding-top: 8px;
      font-size: 11px;
    }
  }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>