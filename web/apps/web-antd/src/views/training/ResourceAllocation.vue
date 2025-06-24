<template>
  <div class="resource-allocation-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <DatabaseOutlined class="title-icon" />
            <span class="title-text">资源分配管理</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text">管理和分配 GPU 训练资源</span>
          </p>
        </div>
        <div class="action-section">
          <a-button
            type="primary"
            size="large"
            @click="showCreateAllocationModal"
            class="create-btn"
          >
            <PlusOutlined />
            新建分配
          </a-button>
        </div>
      </div>
    </div>

    <!-- 资源概览卡片 -->
    <div class="overview-section">
      <a-row :gutter="16">
        <a-col :xs="24" :sm="12" :md="6" :lg="6">
          <a-card class="overview-card glass-card" :bordered="false">
            <div class="overview-item">
              <div class="overview-icon gpu-icon">
                <ThunderboltOutlined />
              </div>
              <div class="overview-content">
                <div class="overview-value">{{ totalGPUs }}</div>
                <div class="overview-label">总 GPU 数</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6" :lg="6">
          <a-card class="overview-card glass-card" :bordered="false">
            <div class="overview-item">
              <div class="overview-icon allocated-icon">
                <CheckCircleOutlined />
              </div>
              <div class="overview-content">
                <div class="overview-value">{{ allocatedGPUs }}</div>
                <div class="overview-label">已分配</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6" :lg="6">
          <a-card class="overview-card glass-card" :bordered="false">
            <div class="overview-item">
              <div class="overview-icon available-icon">
                <ClockCircleOutlined />
              </div>
              <div class="overview-content">
                <div class="overview-value">{{ availableGPUs }}</div>
                <div class="overview-label">可用</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :xs="24" :sm="12" :md="6" :lg="6">
          <a-card class="overview-card glass-card" :bordered="false">
            <div class="overview-item">
              <div class="overview-icon utilization-icon">
                <FireOutlined />
              </div>
              <div class="overview-content">
                <div class="overview-value">{{ utilizationRate }}%</div>
                <div class="overview-label">利用率</div>
              </div>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 筛选器 -->
    <div class="filter-section">
      <a-card class="filter-card glass-card" :bordered="false">
        <a-row :gutter="16" align="middle">
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterStatus"
              placeholder="选择状态"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部状态</a-select-option>
              <a-select-option value="allocated">已分配</a-select-option>
              <a-select-option value="pending">等待中</a-select-option>
              <a-select-option value="running">运行中</a-select-option>
              <a-select-option value="completed">已完成</a-select-option>
              <a-select-option value="failed">失败</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterNodeType"
              placeholder="选择节点类型"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部类型</a-select-option>
              <a-select-option value="gpu-v100">V100</a-select-option>
              <a-select-option value="gpu-a100">A100</a-select-option>
              <a-select-option value="gpu-t4">T4</a-select-option>
              <a-select-option value="gpu-rtx3090">RTX 3090</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="16" :md="8" :lg="8">
            <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索任务名称或用户"
              allow-clear
              @search="handleSearch"
              @change="handleSearchChange"
              class="search-input"
            />
          </a-col>
          <a-col :xs="24" :sm="8" :md="4" :lg="4" class="refresh-btn-col">
            <a-button
              @click="refreshData"
              :loading="loading"
              class="refresh-btn"
            >
              <ReloadOutlined />
              刷新
            </a-button>
          </a-col>
        </a-row>
      </a-card>
    </div>

    <!-- 资源分配表格 -->
    <div class="table-section">
      <a-card class="table-card glass-card" :bordered="false">
        <a-table
          :columns="columns"
          :data-source="filteredAllocations"
          :loading="loading"
          :pagination="paginationConfig"
          row-key="id"
          size="middle"
          :scroll="{ x: 'max-content' }"
          class="sci-fi-table"
        >
          <!-- 状态列 -->
          <template #status="{ record }">
            <div class="status-wrapper">
              <a-tag :color="getStatusColor(record.status)" class="status-tag">
                <component
                  :is="getStatusIcon(record.status)"
                  class="status-icon"
                />
                {{ getStatusText(record.status) }}
              </a-tag>
              <div
                class="status-indicator"
                :class="`indicator-${record.status}`"
              ></div>
            </div>
          </template>

          <!-- 资源配置列 -->
          <template #resources="{ record }">
            <div class="resources-info">
              <div class="resource-item">
                <ThunderboltOutlined class="resource-icon" />
                <span class="resource-label">GPU:</span>
                <span class="resource-value">{{ record.resources.gpu }}</span>
              </div>
              <div class="resource-item">
                <DatabaseOutlined class="resource-icon" />
                <span class="resource-label">CPU:</span>
                <span class="resource-value">{{ record.resources.cpu }}核</span>
              </div>
              <div class="resource-item">
                <BugOutlined class="resource-icon" />
                <span class="resource-label">内存:</span>
                <span class="resource-value"
                  >{{ record.resources.memory }}GB</span
                >
              </div>
            </div>
          </template>

          <!-- 使用时长列 -->
          <template #duration="{ record }">
            <div class="duration-info">
              <div class="duration-item">
                <span class="duration-label">已用:</span>
                <span class="duration-value">{{
                  formatDuration(record.usedTime)
                }}</span>
              </div>
              <div class="duration-item" v-if="record.status === 'running'">
                <span class="duration-label">预计:</span>
                <span class="duration-value">{{
                  formatDuration(record.estimatedTime)
                }}</span>
              </div>
            </div>
          </template>

          <!-- 创建时间列 -->
          <template #createTime="{ record }">
            <a-tooltip :title="record.createTime">
              <span class="time-text">{{
                formatRelativeTime(record.createTime)
              }}</span>
            </a-tooltip>
          </template>

          <!-- 操作列 -->
          <template #action="{ record }">
            <a-space class="action-buttons">
              <a-button
                type="link"
                size="small"
                @click="viewDetails(record)"
                class="action-btn"
              >
                <EyeOutlined />
                详情
              </a-button>
              <a-dropdown>
                <a-button type="link" size="small" class="action-btn">
                  <MoreOutlined />
                </a-button>
                <template #overlay>
                  <a-menu
                    @click="(item: any) => handleMenuAction(item.key, record)"
                    class="action-menu"
                  >
                    <a-menu-item
                      key="modify"
                      :disabled="
                        !['allocated', 'pending'].includes(record.status)
                      "
                    >
                      <EditOutlined />
                      修改配置
                    </a-menu-item>
                    <a-menu-item
                      key="extend"
                      :disabled="record.status !== 'running'"
                    >
                      <ClockCircleOutlined />
                      延长时间
                    </a-menu-item>
                    <a-menu-item
                      key="terminate"
                      :disabled="
                        !['allocated', 'running', 'pending'].includes(
                          record.status,
                        )
                      "
                    >
                      <StopOutlined />
                      终止任务
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item key="delete" class="danger-item">
                      <DeleteOutlined />
                      删除记录
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </a-space>
          </template>
        </a-table>
      </a-card>
    </div>

    <!-- 新建资源分配模态框 -->
    <a-modal
      v-model:open="createModalVisible"
      title="新建资源分配"
      width="900px"
      :confirm-loading="createLoading"
      @ok="handleCreateSubmit"
      @cancel="handleCreateCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="createFormRef"
        :model="createForm"
        :rules="createFormRules"
        layout="vertical"
        class="create-form"
      >
        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="任务名称" name="taskName">
              <a-input
                v-model:value="createForm.taskName"
                placeholder="请输入任务名称"
                class="form-input"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="用户" name="user">
              <a-select
                v-model:value="createForm.user"
                placeholder="选择用户"
                class="form-select"
                show-search
                :filter-option="filterUserOption"
              >
                <a-select-option value="admin">admin</a-select-option>
                <a-select-option value="researcher">researcher</a-select-option>
                <a-select-option value="developer">developer</a-select-option>
                <a-select-option value="student">student</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="队列" name="queue">
              <a-select
                v-model:value="createForm.queue"
                placeholder="选择队列"
                class="form-select"
              >
                <a-select-option value="high-priority"
                  >高优先级</a-select-option
                >
                <a-select-option value="normal">普通</a-select-option>
                <a-select-option value="low-priority">低优先级</a-select-option>
                <a-select-option value="preemptible">可抢占</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="节点类型" name="nodeType">
              <a-select
                v-model:value="createForm.nodeType"
                placeholder="选择节点类型"
                class="form-select"
                @change="handleNodeTypeChange"
              >
                <a-select-option value="gpu-v100">V100 (32GB)</a-select-option>
                <a-select-option value="gpu-a100">A100 (40GB)</a-select-option>
                <a-select-option value="gpu-t4">T4 (16GB)</a-select-option>
                <a-select-option value="gpu-rtx3090"
                  >RTX 3090 (24GB)</a-select-option
                >
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider class="form-divider">资源配置</a-divider>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="8">
            <a-form-item label="GPU 卡数" name="gpu">
              <a-input-number
                v-model:value="createForm.gpu"
                :min="1"
                :max="availableGPUsByType"
                style="width: 100%"
                addon-after="卡"
                class="form-input-number"
              />
              <div class="resource-hint">
                最大可用: {{ availableGPUsByType }} 卡
              </div>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="8">
            <a-form-item label="CPU 核数" name="cpu">
              <a-input-number
                v-model:value="createForm.cpu"
                :min="1"
                :max="64"
                style="width: 100%"
                addon-after="核"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="8">
            <a-form-item label="内存" name="memory">
              <a-input-number
                v-model:value="createForm.memory"
                :min="1"
                :max="256"
                style="width: 100%"
                addon-after="GB"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="预计运行时间" name="estimatedTime">
              <a-input-number
                v-model:value="createForm.estimatedTime"
                :min="1"
                :max="72"
                style="width: 100%"
                addon-after="小时"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="优先级" name="priority">
              <a-select
                v-model:value="createForm.priority"
                placeholder="选择优先级"
                class="form-select"
              >
                <a-select-option value="urgent">紧急 (3)</a-select-option>
                <a-select-option value="high">高 (2)</a-select-option>
                <a-select-option value="normal">普通 (1)</a-select-option>
                <a-select-option value="low">低 (0)</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="任务描述" name="description">
          <a-textarea
            v-model:value="createForm.description"
            placeholder="请输入任务描述"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>

        <!-- 资源预估 -->
        <a-card class="resource-preview-card" title="资源预估" size="small">
          <div class="resource-preview">
            <div class="preview-item">
              <span class="preview-label">预计费用:</span>
              <span class="preview-value">¥{{ estimatedCost }}</span>
            </div>
            <div class="preview-item">
              <span class="preview-label">预计完成时间:</span>
              <span class="preview-value">{{ estimatedCompletionTime }}</span>
            </div>
            <div class="preview-item">
              <span class="preview-label">队列等待时间:</span>
              <span class="preview-value">{{ estimatedWaitTime }}</span>
            </div>
          </div>
        </a-card>
      </a-form>
    </a-modal>

    <!-- 详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="资源分配详情"
      width="1000px"
      :footer="null"
      class="sci-fi-modal detail-modal"
    >
      <div v-if="selectedAllocation" class="detail-content">
        <a-tabs v-model:activeKey="detailActiveTab" class="detail-tabs">
          <a-tab-pane key="basic" tab="基本信息">
            <a-descriptions
              :column="{ xs: 1, sm: 2 }"
              bordered
              class="detail-descriptions"
            >
              <a-descriptions-item label="任务名称">
                {{ selectedAllocation.taskName }}
              </a-descriptions-item>
              <a-descriptions-item label="状态">
                <a-tag
                  :color="getStatusColor(selectedAllocation.status)"
                  class="status-tag"
                >
                  <component :is="getStatusIcon(selectedAllocation.status)" />
                  {{ getStatusText(selectedAllocation.status) }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="用户">
                {{ selectedAllocation.user }}
              </a-descriptions-item>
              <a-descriptions-item label="队列">
                {{ selectedAllocation.queue }}
              </a-descriptions-item>
              <a-descriptions-item label="节点类型">
                {{ selectedAllocation.nodeType }}
              </a-descriptions-item>
              <a-descriptions-item label="创建时间">
                {{ selectedAllocation.createTime }}
              </a-descriptions-item>
              <a-descriptions-item label="GPU">
                {{ selectedAllocation.resources.gpu }} 卡
              </a-descriptions-item>
              <a-descriptions-item label="CPU">
                {{ selectedAllocation.resources.cpu }} 核
              </a-descriptions-item>
              <a-descriptions-item label="内存">
                {{ selectedAllocation.resources.memory }} GB
              </a-descriptions-item>
              <a-descriptions-item label="优先级">
                {{ getPriorityText(selectedAllocation.priority) }}
              </a-descriptions-item>
              <a-descriptions-item label="已用时间">
                {{ formatDuration(selectedAllocation.usedTime) }}
              </a-descriptions-item>
              <a-descriptions-item label="预计时间">
                {{ formatDuration(selectedAllocation.estimatedTime) }}
              </a-descriptions-item>
              <a-descriptions-item label="描述" :span="2">
                {{ selectedAllocation.description || '暂无描述' }}
              </a-descriptions-item>
            </a-descriptions>
          </a-tab-pane>

          <a-tab-pane key="monitoring" tab="资源监控">
            <div class="monitoring-content">
              <a-row :gutter="16">
                <a-col :xs="24" :sm="12">
                  <a-card
                    title="GPU 使用率"
                    size="small"
                    class="monitoring-card"
                  >
                    <div class="metric-item">
                      <a-progress
                        :percent="75"
                        status="active"
                        stroke-color="#52c41a"
                      />
                      <div class="metric-label">GPU 0: 75%</div>
                    </div>
                    <div class="metric-item">
                      <a-progress
                        :percent="82"
                        status="active"
                        stroke-color="#52c41a"
                      />
                      <div class="metric-label">GPU 1: 82%</div>
                    </div>
                  </a-card>
                </a-col>
                <a-col :xs="24" :sm="12">
                  <a-card
                    title="内存使用率"
                    size="small"
                    class="monitoring-card"
                  >
                    <div class="metric-item">
                      <a-progress
                        :percent="68"
                        status="active"
                        stroke-color="#1890ff"
                      />
                      <div class="metric-label">CPU 内存: 68%</div>
                    </div>
                    <div class="metric-item">
                      <a-progress
                        :percent="45"
                        status="active"
                        stroke-color="#722ed1"
                      />
                      <div class="metric-label">GPU 内存: 45%</div>
                    </div>
                  </a-card>
                </a-col>
              </a-row>
            </div>
          </a-tab-pane>

          <a-tab-pane key="logs" tab="运行日志">
            <div class="log-container">
              <div class="log-header">
                <span class="log-title">实时日志</span>
                <a-button
                  size="small"
                  @click="refreshLogs"
                  class="log-refresh-btn"
                >
                  <ReloadOutlined />
                  刷新
                </a-button>
              </div>
              <div class="log-content">
                <pre
                  v-for="(log, index) in logs"
                  :key="index"
                  class="log-line"
                  >{{ log }}</pre
                >
              </div>
            </div>
          </a-tab-pane>
        </a-tabs>
      </div>
    </a-modal>

    <!-- 修改配置模态框 -->
    <a-modal
      v-model:open="modifyModalVisible"
      title="修改资源配置"
      width="700px"
      :confirm-loading="modifyLoading"
      @ok="handleModifySubmit"
      @cancel="handleModifyCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="modifyFormRef"
        :model="modifyForm"
        :rules="modifyFormRules"
        layout="vertical"
      >
        <a-alert
          message="提示"
          description="修改配置可能会导致任务重新排队，请谨慎操作。"
          type="warning"
          style="margin-bottom: 16px"
        />

        <a-row :gutter="16">
          <a-col :xs="24" :sm="8">
            <a-form-item label="GPU 卡数" name="gpu">
              <a-input-number
                v-model:value="modifyForm.gpu"
                :min="1"
                :max="8"
                style="width: 100%"
                addon-after="卡"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="8">
            <a-form-item label="CPU 核数" name="cpu">
              <a-input-number
                v-model:value="modifyForm.cpu"
                :min="1"
                :max="64"
                style="width: 100%"
                addon-after="核"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="8">
            <a-form-item label="内存" name="memory">
              <a-input-number
                v-model:value="modifyForm.memory"
                :min="1"
                :max="256"
                style="width: 100%"
                addon-after="GB"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="优先级" name="priority">
          <a-select
            v-model:value="modifyForm.priority"
            placeholder="选择优先级"
            class="form-select"
          >
            <a-select-option value="urgent">紧急 (3)</a-select-option>
            <a-select-option value="high">高 (2)</a-select-option>
            <a-select-option value="normal">普通 (1)</a-select-option>
            <a-select-option value="low">低 (0)</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 延长时间模态框 -->
    <a-modal
      v-model:open="extendModalVisible"
      title="延长运行时间"
      width="500px"
      :confirm-loading="extendLoading"
      @ok="handleExtendSubmit"
      @cancel="handleExtendCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="extendFormRef"
        :model="extendForm"
        :rules="extendFormRules"
        layout="vertical"
      >
        <a-form-item label="延长时间" name="extendTime">
          <a-input-number
            v-model:value="extendForm.extendTime"
            :min="1"
            :max="24"
            style="width: 100%"
            addon-after="小时"
            class="form-input-number"
          />
        </a-form-item>

        <a-form-item label="延长原因" name="reason">
          <a-textarea
            v-model:value="extendForm.reason"
            placeholder="请输入延长原因"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>

        <a-alert
          :message="`延长费用: ¥${extendCost}`"
          type="info"
          style="margin-top: 16px"
        />
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { message, Modal } from 'ant-design-vue';
import type { FormInstance, TableColumnsType } from 'ant-design-vue';
import {
  DatabaseOutlined,
  PlusOutlined,
  ReloadOutlined,
  ThunderboltOutlined,
  BugOutlined,
  EyeOutlined,
  MoreOutlined,
  EditOutlined,
  DeleteOutlined,
  StopOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  CloseCircleOutlined,
  PlayCircleOutlined,
  FireOutlined,
} from '@ant-design/icons-vue';

// ===== 类型定义 =====
interface ResourceConfig {
  gpu: number;
  cpu: number;
  memory: number;
}

interface ResourceAllocation {
  id: string;
  taskName: string;
  user: string;
  queue: string;
  nodeType: string;
  status:
    | 'allocated'
    | 'pending'
    | 'running'
    | 'completed'
    | 'failed'
    | 'terminated';
  priority: 'urgent' | 'high' | 'normal' | 'low';
  resources: ResourceConfig;
  usedTime: number; // 分钟
  estimatedTime: number; // 分钟
  createTime: string;
  description?: string;
}

interface CreateAllocationForm {
  taskName: string;
  user: string;
  queue: string;
  nodeType: string;
  gpu: number;
  cpu: number;
  memory: number;
  estimatedTime: number;
  priority: string;
  description: string;
}

interface ModifyForm {
  gpu: number;
  cpu: number;
  memory: number;
  priority: string;
}

interface ExtendForm {
  extendTime: number;
  reason: string;
}

// ===== 响应式数据 =====
const loading = ref<boolean>(false);
const createModalVisible = ref<boolean>(false);
const detailModalVisible = ref<boolean>(false);
const modifyModalVisible = ref<boolean>(false);
const extendModalVisible = ref<boolean>(false);
const createLoading = ref<boolean>(false);
const modifyLoading = ref<boolean>(false);
const extendLoading = ref<boolean>(false);

const filterStatus = ref<string>('');
const filterNodeType = ref<string>('');
const searchKeyword = ref<string>('');
const detailActiveTab = ref<string>('basic');

const selectedAllocation = ref<ResourceAllocation | null>(null);

// ===== 表单引用 =====
const createFormRef = ref<FormInstance>();
const modifyFormRef = ref<FormInstance>();
const extendFormRef = ref<FormInstance>();

// ===== 表单数据 =====
const createForm = reactive<CreateAllocationForm>({
  taskName: '',
  user: '',
  queue: 'normal',
  nodeType: '',
  gpu: 1,
  cpu: 4,
  memory: 16,
  estimatedTime: 2,
  priority: 'normal',
  description: '',
});

const modifyForm = reactive<ModifyForm>({
  gpu: 1,
  cpu: 4,
  memory: 16,
  priority: 'normal',
});

const extendForm = reactive<ExtendForm>({
  extendTime: 1,
  reason: '',
});

// ===== 配置数据 =====
const STATUS_CONFIG = {
  allocated: { color: 'processing', text: '已分配', icon: CheckCircleOutlined },
  pending: { color: 'warning', text: '等待中', icon: ClockCircleOutlined },
  running: { color: 'success', text: '运行中', icon: PlayCircleOutlined },
  completed: { color: 'success', text: '已完成', icon: CheckCircleOutlined },
  failed: { color: 'error', text: '失败', icon: CloseCircleOutlined },
  terminated: { color: 'default', text: '已终止', icon: StopOutlined },
} as const;

const PRIORITY_CONFIG = {
  urgent: '紧急 (3)',
  high: '高 (2)',
  normal: '普通 (1)',
  low: '低 (0)',
} as const;

const NODE_TYPE_CONFIG = {
  'gpu-v100': { available: 8, gpuMemory: 32 },
  'gpu-a100': { available: 4, gpuMemory: 40 },
  'gpu-t4': { available: 16, gpuMemory: 16 },
  'gpu-rtx3090': { available: 12, gpuMemory: 24 },
} as const;

// ===== 模拟数据 =====
const allocations = ref<ResourceAllocation[]>([
  {
    id: 'alloc-001',
    taskName: 'bert-training-v2',
    user: 'admin',
    queue: 'high-priority',
    nodeType: 'gpu-a100',
    status: 'running',
    priority: 'high',
    resources: { gpu: 2, cpu: 16, memory: 64 },
    usedTime: 245, // 4小时5分钟
    estimatedTime: 480, // 8小时
    createTime: '2024-06-23 08:30:00',
    description: 'BERT 模型大规模训练',
  },
  {
    id: 'alloc-002',
    taskName: 'image-classification',
    user: 'researcher',
    queue: 'normal',
    nodeType: 'gpu-v100',
    status: 'pending',
    priority: 'normal',
    resources: { gpu: 1, cpu: 8, memory: 32 },
    usedTime: 0,
    estimatedTime: 180, // 3小时
    createTime: '2024-06-23 09:15:00',
    description: '图像分类模型训练',
  },
  {
    id: 'alloc-003',
    taskName: 'nlp-experiment',
    user: 'developer',
    queue: 'normal',
    nodeType: 'gpu-t4',
    status: 'allocated',
    priority: 'normal',
    resources: { gpu: 4, cpu: 32, memory: 128 },
    usedTime: 0,
    estimatedTime: 360, // 6小时
    createTime: '2024-06-23 10:00:00',
    description: 'NLP 模型实验',
  },
  {
    id: 'alloc-004',
    taskName: 'data-preprocessing',
    user: 'student',
    queue: 'low-priority',
    nodeType: 'gpu-rtx3090',
    status: 'completed',
    priority: 'low',
    resources: { gpu: 1, cpu: 4, memory: 16 },
    usedTime: 135, // 2小时15分钟
    estimatedTime: 120, // 2小时
    createTime: '2024-06-22 14:30:00',
    description: '数据预处理任务',
  },
  {
    id: 'alloc-005',
    taskName: 'model-evaluation',
    user: 'researcher',
    queue: 'normal',
    nodeType: 'gpu-v100',
    status: 'failed',
    priority: 'normal',
    resources: { gpu: 2, cpu: 16, memory: 64 },
    usedTime: 45,
    estimatedTime: 240,
    createTime: '2024-06-23 07:45:00',
    description: '模型性能评估',
  },
]);

const logs = ref<string[]>([
  '2024-06-23 10:30:15 INFO: Task started successfully',
  '2024-06-23 10:30:16 INFO: GPU allocation complete',
  '2024-06-23 10:30:17 INFO: Loading training data...',
  '2024-06-23 10:30:20 INFO: Model training started',
  '2024-06-23 10:35:15 INFO: Epoch 1/100 completed, loss: 0.8324',
  '2024-06-23 10:40:15 INFO: Epoch 2/100 completed, loss: 0.7891',
]);

// ===== 表单验证规则 =====
const createFormRules = {
  taskName: [
    { required: true, message: '请输入任务名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
  ],
  user: [{ required: true, message: '请选择用户', trigger: 'change' }],
  queue: [{ required: true, message: '请选择队列', trigger: 'change' }],
  nodeType: [{ required: true, message: '请选择节点类型', trigger: 'change' }],
  gpu: [{ required: true, message: '请输入 GPU 卡数', trigger: 'blur' }],
  cpu: [{ required: true, message: '请输入 CPU 核数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
  estimatedTime: [
    { required: true, message: '请输入预计运行时间', trigger: 'blur' },
  ],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }],
};

const modifyFormRules = {
  gpu: [{ required: true, message: '请输入 GPU 卡数', trigger: 'blur' }],
  cpu: [{ required: true, message: '请输入 CPU 核数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }],
};

const extendFormRules = {
  extendTime: [{ required: true, message: '请输入延长时间', trigger: 'blur' }],
  reason: [{ required: true, message: '请输入延长原因', trigger: 'blur' }],
};

// ===== 表格列配置 =====
const columns: TableColumnsType<ResourceAllocation> = [
  {
    title: '任务名称',
    dataIndex: 'taskName',
    key: 'taskName',
    width: 180,
    ellipsis: true,
  },
  {
    title: '用户',
    dataIndex: 'user',
    key: 'user',
    width: 100,
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
    slots: { customRender: 'status' },
  },
  {
    title: '队列',
    dataIndex: 'queue',
    key: 'queue',
    width: 120,
  },
  {
    title: '节点类型',
    dataIndex: 'nodeType',
    key: 'nodeType',
    width: 120,
  },
  {
    title: '资源配置',
    key: 'resources',
    width: 180,
    slots: { customRender: 'resources' },
  },
  {
    title: '使用时长',
    key: 'duration',
    width: 150,
    slots: { customRender: 'duration' },
  },
  {
    title: '创建时间',
    key: 'createTime',
    width: 150,
    slots: { customRender: 'createTime' },
  },
  {
    title: '操作',
    key: 'action',
    width: 120,
    fixed: 'right',
    slots: { customRender: 'action' },
  },
];

// ===== 分页配置 =====
const paginationConfig = {
  total: computed(() => filteredAllocations.value.length),
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// ===== 计算属性 =====
const totalGPUs = computed(() => {
  return Object.values(NODE_TYPE_CONFIG).reduce(
    (sum, config) => sum + config.available,
    0,
  );
});

const allocatedGPUs = computed(() => {
  return allocations.value
    .filter((item) => ['allocated', 'running'].includes(item.status))
    .reduce((sum, item) => sum + item.resources.gpu, 0);
});

const availableGPUs = computed(() => {
  return totalGPUs.value - allocatedGPUs.value;
});

const utilizationRate = computed(() => {
  return totalGPUs.value > 0
    ? Math.round((allocatedGPUs.value / totalGPUs.value) * 100)
    : 0;
});

const availableGPUsByType = computed(() => {
  if (!createForm.nodeType) return 0;
  const typeConfig =
    NODE_TYPE_CONFIG[createForm.nodeType as keyof typeof NODE_TYPE_CONFIG];
  if (!typeConfig) return 0;

  const usedByType = allocations.value
    .filter(
      (item) =>
        item.nodeType === createForm.nodeType &&
        ['allocated', 'running'].includes(item.status),
    )
    .reduce((sum, item) => sum + item.resources.gpu, 0);

  return Math.max(0, typeConfig.available - usedByType);
});

const estimatedCost = computed(() => {
  const basePrice = 10; // 每GPU每小时基础价格
  const gpuCount = createForm.gpu || 1;
  const hours = createForm.estimatedTime || 1;
  return (basePrice * gpuCount * hours).toFixed(2);
});

const estimatedCompletionTime = computed(() => {
  const now = new Date();
  const estimatedMinutes = (createForm.estimatedTime || 1) * 60;
  const completionTime = new Date(now.getTime() + estimatedMinutes * 60000);
  return completionTime.toLocaleString();
});

const estimatedWaitTime = computed(() => {
  // 模拟队列等待时间计算
  const queueMultiplier = {
    'high-priority': 0.5,
    normal: 1,
    'low-priority': 2,
    preemptible: 0.8,
  };
  const baseWaitTime = 15; // 基础等待时间（分钟）
  const multiplier =
    queueMultiplier[createForm.queue as keyof typeof queueMultiplier] || 1;
  const waitMinutes = Math.round(baseWaitTime * multiplier);
  return `约 ${waitMinutes} 分钟`;
});

const extendCost = computed(() => {
  const basePrice = 10; // 每GPU每小时基础价格
  const gpuCount = selectedAllocation.value?.resources.gpu || 1;
  const hours = extendForm.extendTime || 1;
  return (basePrice * gpuCount * hours).toFixed(2);
});

const filteredAllocations = computed(() => {
  let result = allocations.value;

  if (filterStatus.value) {
    result = result.filter((item) => item.status === filterStatus.value);
  }

  if (filterNodeType.value) {
    result = result.filter((item) => item.nodeType === filterNodeType.value);
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter(
      (item) =>
        item.taskName.toLowerCase().includes(keyword) ||
        item.user.toLowerCase().includes(keyword),
    );
  }

  return result;
});

// ===== 工具函数 =====
const getStatusColor = (status: string): string => {
  return (
    STATUS_CONFIG[status as keyof typeof STATUS_CONFIG]?.color || 'default'
  );
};

const getStatusIcon = (status: string) => {
  return (
    STATUS_CONFIG[status as keyof typeof STATUS_CONFIG]?.icon ||
    ClockCircleOutlined
  );
};

const getStatusText = (status: string): string => {
  return STATUS_CONFIG[status as keyof typeof STATUS_CONFIG]?.text || status;
};

const getPriorityText = (priority: string): string => {
  return PRIORITY_CONFIG[priority as keyof typeof PRIORITY_CONFIG] || priority;
};

const formatDuration = (minutes: number): string => {
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  if (hours > 0) {
    return `${hours}h ${mins}m`;
  }
  return `${mins}m`;
};

const formatRelativeTime = (time: string): string => {
  const now = new Date();
  const target = new Date(time);
  const diffMs = now.getTime() - target.getTime();
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));

  if (diffHours < 1) {
    const diffMinutes = Math.floor(diffMs / (1000 * 60));
    return `${diffMinutes} 分钟前`;
  } else if (diffHours < 24) {
    return `${diffHours} 小时前`;
  } else {
    const diffDays = Math.floor(diffHours / 24);
    return `${diffDays} 天前`;
  }
};

const filterUserOption = (input: string, option: any): boolean => {
  return option.value.toLowerCase().includes(input.toLowerCase());
};

// ===== 事件处理函数 =====
const showCreateAllocationModal = (): void => {
  createModalVisible.value = true;
};

const handleNodeTypeChange = (): void => {
  // 节点类型改变时重置GPU数量
  createForm.gpu = 1;
};

const handleCreateSubmit = async (): Promise<void> => {
  try {
    await createFormRef.value?.validate();

    if (createForm.gpu > availableGPUsByType.value) {
      message.error(`选择的GPU数量超过可用数量 (${availableGPUsByType.value})`);
      return;
    }

    createLoading.value = true;

    // 模拟 API 调用
    await new Promise((resolve) => setTimeout(resolve, 2000));

    const newAllocation: ResourceAllocation = {
      id: `alloc-${Date.now()}`,
      taskName: createForm.taskName,
      user: createForm.user,
      queue: createForm.queue,
      nodeType: createForm.nodeType,
      status: 'pending',
      priority: createForm.priority as ResourceAllocation['priority'],
      resources: {
        gpu: createForm.gpu,
        cpu: createForm.cpu,
        memory: createForm.memory,
      },
      usedTime: 0,
      estimatedTime: createForm.estimatedTime * 60, // 转换为分钟
      createTime: new Date().toLocaleString(),
      description: createForm.description,
    };

    allocations.value.unshift(newAllocation);
    createModalVisible.value = false;
    message.success('资源分配创建成功');

    // 重置表单
    createFormRef.value?.resetFields();
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    createLoading.value = false;
  }
};

const handleCreateCancel = (): void => {
  createModalVisible.value = false;
  createFormRef.value?.resetFields();
};

const viewDetails = (record: ResourceAllocation): void => {
  selectedAllocation.value = record;
  detailModalVisible.value = true;
  detailActiveTab.value = 'basic';
};

const handleMenuAction = (key: string, record: ResourceAllocation): void => {
  const actions = {
    modify: () => handleModify(record),
    extend: () => handleExtend(record),
    terminate: () => handleTerminate(record),
    delete: () => handleDelete(record),
  };

  const action = actions[key as keyof typeof actions];
  if (action) {
    action();
  }
};

const handleModify = (record: ResourceAllocation): void => {
  selectedAllocation.value = record;
  modifyForm.gpu = record.resources.gpu;
  modifyForm.cpu = record.resources.cpu;
  modifyForm.memory = record.resources.memory;
  modifyForm.priority = record.priority;
  modifyModalVisible.value = true;
};

const handleModifySubmit = async (): Promise<void> => {
  try {
    await modifyFormRef.value?.validate();
    modifyLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 1500));

    if (selectedAllocation.value) {
      const index = allocations.value.findIndex(
        (item) => item.id === selectedAllocation.value!.id,
      );
      if (index !== -1) {
        allocations.value[index]!.resources = {
          gpu: modifyForm.gpu,
          cpu: modifyForm.cpu,
          memory: modifyForm.memory,
        };
        allocations.value[index]!.priority =
          modifyForm.priority as ResourceAllocation['priority'];
      }
    }

    modifyModalVisible.value = false;
    message.success('资源配置修改成功');
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    modifyLoading.value = false;
  }
};

const handleModifyCancel = (): void => {
  modifyModalVisible.value = false;
  modifyFormRef.value?.resetFields();
};

const handleExtend = (record: ResourceAllocation): void => {
  selectedAllocation.value = record;
  extendModalVisible.value = true;
};

const handleExtendSubmit = async (): Promise<void> => {
  try {
    await extendFormRef.value?.validate();
    extendLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 1500));

    if (selectedAllocation.value) {
      const index = allocations.value.findIndex(
        (item) => item.id === selectedAllocation.value!.id,
      );
      if (index !== -1) {
        allocations.value[index]!.estimatedTime += extendForm.extendTime * 60;
      }
    }

    extendModalVisible.value = false;
    message.success('运行时间延长成功');

    // 重置表单
    extendFormRef.value?.resetFields();
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    extendLoading.value = false;
  }
};

const handleExtendCancel = (): void => {
  extendModalVisible.value = false;
  extendFormRef.value?.resetFields();
};

const handleTerminate = (record: ResourceAllocation): void => {
  Modal.confirm({
    title: '确认终止',
    content: `确定要终止任务 "${record.taskName}" 吗？此操作不可恢复。`,
    okText: '确认',
    cancelText: '取消',
    type: 'warning',
    onOk: async () => {
      const index = allocations.value.findIndex(
        (item) => item.id === record.id,
      );
      if (index !== -1) {
        allocations.value[index]!.status = 'terminated';
        message.success('任务终止成功');
      }
    },
  });
};

const handleDelete = (record: ResourceAllocation): void => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除任务记录 "${record.taskName}" 吗？此操作不可恢复。`,
    okText: '确认',
    cancelText: '取消',
    type: 'warning',
    onOk: () => {
      const index = allocations.value.findIndex(
        (item) => item.id === record.id,
      );
      if (index !== -1) {
        allocations.value.splice(index, 1);
        message.success('记录删除成功');
      }
    },
  });
};

const refreshData = async (): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    message.success('数据刷新成功');
  } catch (error) {
    message.error('刷新失败');
  } finally {
    loading.value = false;
  }
};

const refreshLogs = async (): Promise<void> => {
  const newLogs = [
    ...logs.value,
    `${new Date().toLocaleString()} INFO: Log refreshed`,
  ];
  logs.value = newLogs.slice(-50);
};

const handleFilterChange = (): void => {
  // 筛选变化时的处理逻辑
};

const handleSearch = (): void => {
  // 搜索处理逻辑
};

const handleSearchChange = (): void => {
  // 搜索输入变化时的处理逻辑
};

// ===== 生命周期 =====
onMounted(() => {
  refreshData();
});
</script>

<style scoped>
/* ===== 基础样式 ===== */
.resource-allocation-container {
  padding: 24px;
  min-height: 100vh;
}

/* ===== 卡片样式 ===== */
.glass-card {
  border-radius: 8px !important;
}

/* ===== 页面头部 ===== */
.page-header {
  margin-bottom: 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
}

.title-section {
  flex: 1;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  font-size: 32px;
  color: #1890ff;
}

.page-description {
  font-size: 16px;
  margin: 0;
}

/* ===== 资源概览样式 ===== */
.overview-section {
  margin-bottom: 24px;
}

.overview-card {
  border-radius: 8px !important;
  transition: all 0.3s ease;
}

.overview-card:hover {
  transform: translateY(-2px);
}

.overview-item {
  display: flex;
  align-items: center;
  gap: 16px;
}

.overview-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.gpu-icon {
  background: linear-gradient(135deg, #1890ff, #096dd9);
}

.allocated-icon {
  background: linear-gradient(135deg, #52c41a, #389e0d);
}

.available-icon {
  background: linear-gradient(135deg, #faad14, #d48806);
}

.utilization-icon {
  background: linear-gradient(135deg, #722ed1, #531dab);
}

.overview-content {
  flex: 1;
}

.overview-value {
  font-size: 24px;
  font-weight: 700;
  line-height: 1;
  margin-bottom: 4px;
}

.overview-label {
  font-size: 14px;
  opacity: 0.8;
}

/* ===== 按钮样式 ===== */
.create-btn {
  border: none !important;
  height: 40px !important;
  padding: 0 24px !important;
  font-size: 14px !important;
  font-weight: 500 !important;
  border-radius: 6px !important;
  transition: all 0.3s ease !important;
}

.create-btn:hover {
  background: #1890ff !important;
  transform: translateY(-1px);
}

/* ===== 筛选器样式 ===== */
.filter-section {
  margin-bottom: 24px;
}

.filter-card {
  border-radius: 8px !important;
}

.filter-select,
.search-input,
.refresh-btn {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

.refresh-btn:hover {
  color: #1890ff !important;
  border-color: #1890ff !important;
}

/* ===== 表格样式 ===== */
.table-section {
  margin-bottom: 24px;
}

.table-card {
  border-radius: 8px !important;
}

.sci-fi-table :deep(.ant-table-thead > tr > th) {
  font-weight: 600 !important;
}

/* ===== 状态标签 ===== */
.status-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border-radius: 6px !important;
  font-weight: 500 !important;
  padding: 4px 8px !important;
}

.status-icon {
  font-size: 12px;
}

.status-indicator {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.indicator-allocated {
  background: #1890ff;
}

.indicator-pending {
  background: #faad14;
}

.indicator-running {
  background: #52c41a;
}

.indicator-completed {
  background: #52c41a;
}

.indicator-failed {
  background: #ff4d4f;
}

.indicator-terminated {
  background: #8c8c8c;
}

/* ===== 资源信息 ===== */
.resources-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.resource-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.resource-item:hover {
  color: #1890ff;
}

.resource-icon {
  font-size: 12px;
  color: #1890ff;
}

.resource-label {
  font-weight: 500;
}

.resource-value {
  font-weight: 600;
}

/* ===== 使用时长信息 ===== */
.duration-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.duration-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
}

.duration-label {
  font-weight: 500;
  opacity: 0.8;
}

.duration-value {
  font-weight: 600;
}

/* ===== 时间显示 ===== */
.time-text {
  font-size: 12px;
}

/* ===== 操作按钮 ===== */
.action-buttons {
  display: flex;
  gap: 4px;
}

.action-btn {
  border: none !important;
  background: transparent !important;
  border-radius: 4px !important;
  padding: 4px 8px !important;
  height: auto !important;
  font-size: 12px !important;
  transition: all 0.3s ease !important;
}

.action-btn:hover {
  color: #1890ff !important;
}

.action-btn:disabled {
  color: #bfbfbf !important;
  background: transparent !important;
}

.action-menu {
  border-radius: 8px !important;
}

.action-menu :deep(.ant-menu-item) {
  border-radius: 4px !important;
  margin: 2px !important;
  transition: all 0.3s ease;
}

.action-menu :deep(.ant-menu-item:hover) {
  color: #1890ff !important;
}

.danger-item {
  color: #ff4d4f !important;
}

.danger-item:hover {
  color: #ff4d4f !important;
}

/* ===== 模态框样式 ===== */
.sci-fi-modal :deep(.ant-modal-content) {
  border-radius: 8px !important;
}

.sci-fi-modal :deep(.ant-modal-header) {
  border-radius: 8px 8px 0 0 !important;
}

.sci-fi-modal :deep(.ant-modal-title) {
  font-weight: 600 !important;
  font-size: 16px !important;
}

/* ===== 表单样式 ===== */
.create-form :deep(.ant-form-item-label > label) {
  font-weight: 500 !important;
}

.form-input,
.form-select,
.form-textarea,
.form-input-number {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

.form-divider {
  font-weight: 500 !important;
}

/* ===== 资源提示 ===== */
.resource-hint {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 4px;
}

/* ===== 资源预估卡片 ===== */
.resource-preview-card {
  margin-top: 16px;
  border-radius: 6px !important;
}

.resource-preview {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.preview-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
}

.preview-label {
  font-weight: 500;
}

.preview-value {
  font-weight: 600;
  color: #1890ff;
}

/* ===== 详情标签页 ===== */
.detail-tabs :deep(.ant-tabs-tab) {
  border-radius: 6px 6px 0 0 !important;
}

/* ===== 监控内容 ===== */
.monitoring-content {
  padding: 16px 0;
}

.monitoring-card {
  border-radius: 6px !important;
}

.metric-item {
  margin-bottom: 12px;
}

.metric-item:last-child {
  margin-bottom: 0;
}

.metric-label {
  font-size: 12px;
  margin-top: 4px;
  color: #8c8c8c;
}

/* ===== 日志容器 ===== */
.log-container {
  margin-top: 16px;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.log-title {
  font-weight: 600;
  font-size: 14px;
}

.log-refresh-btn {
  border-radius: 4px !important;
  transition: all 0.3s ease;
}

.log-refresh-btn:hover {
  color: #1890ff !important;
  border-color: #1890ff !important;
}

.log-content {
  border-radius: 6px !important;
  padding: 12px !important;
  max-height: 300px;
  overflow-y: auto;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace !important;
}

.log-line {
  margin: 0;
  font-size: 12px;
  line-height: 1.4;
}

/* ===== 响应式设计 ===== */
@media (max-width: 768px) {
  .resource-allocation-container {
    padding: 16px;
  }

  .header-content {
    flex-direction: column;
    gap: 16px;
  }

  .page-title {
    font-size: 24px;
  }

  .title-icon {
    font-size: 28px;
  }

  .action-section {
    align-self: stretch;
  }

  .create-btn,
  .refresh-btn {
    width: 100% !important;
    justify-content: center !important;
  }

  .refresh-btn-col {
    margin-top: 12px;
  }

  .overview-item {
    gap: 12px;
  }

  .overview-icon {
    width: 40px;
    height: 40px;
    font-size: 20px;
  }

  .overview-value {
    font-size: 20px;
  }

  .sci-fi-modal :deep(.ant-modal) {
    margin: 16px !important;
    max-width: calc(100vw - 32px) !important;
  }
}

@media (max-width: 480px) {
  .page-title {
    font-size: 20px;
    flex-direction: column;
    gap: 8px;
    text-align: center;
  }

  .title-icon {
    font-size: 24px;
  }

  .overview-item {
    flex-direction: column;
    text-align: center;
    gap: 8px;
  }

  .overview-icon {
    width: 36px;
    height: 36px;
    font-size: 18px;
  }

  .overview-value {
    font-size: 18px;
  }

  .resources-info {
    gap: 2px;
  }

  .resource-item {
    font-size: 11px;
    padding: 2px 4px;
  }

  .action-buttons {
    flex-direction: column;
    gap: 2px;
  }

  .action-btn {
    font-size: 11px !important;
    padding: 3px 6px !important;
  }
}

/* ===== 滚动条样式 ===== */
::-webkit-scrollbar {
  width: 4px;
}

::-webkit-scrollbar-track,
::-webkit-scrollbar-thumb {
  border-radius: 2px;
}
</style>
