<template>
  <div class="training-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <RocketOutlined class="title-icon" />
            <span class="title-text">训练任务管理</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text">管理和监控您的AI模型训练任务</span>
          </p>
        </div>
        <div class="action-section">
          <a-button
            type="primary"
            size="large"
            @click="showCreateModal"
            class="create-btn"
          >
            <PlusOutlined />
            创建训练任务
          </a-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-section">
      <a-row :gutter="24">
        <a-col :xs="12" :sm="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <div class="stat-content">
              <div class="stat-icon running">
                <PlayCircleOutlined />
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ runningCount }}</div>
                <div class="stat-label">运行中</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <div class="stat-content">
              <div class="stat-icon pending">
                <ClockCircleOutlined />
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ pendingCount }}</div>
                <div class="stat-label">等待中</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <div class="stat-content">
              <div class="stat-icon completed">
                <CheckCircleOutlined />
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ completedCount }}</div>
                <div class="stat-label">已完成</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <div class="stat-content">
              <div class="stat-icon failed">
                <CloseCircleOutlined />
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ failedCount }}</div>
                <div class="stat-label">失败</div>
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
              <a-select-option value="running">运行中</a-select-option>
              <a-select-option value="pending">等待中</a-select-option>
              <a-select-option value="completed">已完成</a-select-option>
              <a-select-option value="failed">失败</a-select-option>
              <a-select-option value="cancelled">已取消</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterFramework"
              placeholder="选择框架"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部框架</a-select-option>
              <a-select-option value="tensorflow">TensorFlow</a-select-option>
              <a-select-option value="pytorch">PyTorch</a-select-option>
              <a-select-option value="mindspore">MindSpore</a-select-option>
              <a-select-option value="paddlepaddle"
                >PaddlePaddle</a-select-option
              >
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="16" :md="8" :lg="8">
            <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索任务名称或创建者"
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

    <!-- 数据表格 -->
    <div class="table-section">
      <a-card class="table-card glass-card" :bordered="false">
        <a-table
          :columns="columns"
          :data-source="filteredTrainingJobs"
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

          <!-- 训练框架列 -->
          <template #framework="{ record }">
            <div class="framework-wrapper">
              <component
                :is="getFrameworkIcon(record.framework)"
                class="framework-icon"
              />
              <span class="framework-text">{{
                getFrameworkText(record.framework)
              }}</span>
            </div>
          </template>

          <!-- 资源配置列 -->
          <template #resources="{ record }">
            <div class="resources-info">
              <div class="resource-item">
                <DatabaseOutlined class="resource-icon" />
                <span class="resource-label">CPU:</span>
                <span class="resource-value">{{ record.resources.cpu }}</span>
              </div>
              <div class="resource-item">
                <ThunderboltOutlined class="resource-icon" />
                <span class="resource-label">内存:</span>
                <span class="resource-value"
                  >{{ record.resources.memory }}GB</span
                >
              </div>
              <div class="resource-item" v-if="record.resources.gpu">
                <BugOutlined class="resource-icon" />
                <span class="resource-label">GPU:</span>
                <span class="resource-value">{{ record.resources.gpu }}</span>
              </div>
            </div>
          </template>

          <!-- 进度列 -->
          <template #progress="{ record }">
            <div class="progress-wrapper">
              <a-progress
                :percent="record.progress"
                :status="getProgressStatus(record.status)"
                size="small"
                :show-info="false"
                class="progress-bar"
              />
              <span class="progress-text">{{ record.progress }}%</span>
            </div>
          </template>

          <!-- 运行时间列 -->
          <template #duration="{ record }">
            <div class="duration-wrapper">
              <ClockCircleOutlined class="duration-icon" />
              <span class="duration-text">{{
                formatDuration(record.duration)
              }}</span>
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
                @click="viewLogs(record)"
                class="action-btn"
              >
                <FileTextOutlined />
                日志
              </a-button>
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
                      key="pause"
                      :disabled="record.status !== 'running'"
                    >
                      <PauseCircleOutlined />
                      暂停
                    </a-menu-item>
                    <a-menu-item
                      key="resume"
                      :disabled="record.status !== 'paused'"
                    >
                      <PlayCircleOutlined />
                      恢复
                    </a-menu-item>
                    <a-menu-item
                      key="stop"
                      :disabled="
                        !['running', 'pending'].includes(record.status)
                      "
                    >
                      <StopOutlined />
                      停止
                    </a-menu-item>
                    <a-menu-item key="clone">
                      <CopyOutlined />
                      克隆
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item key="delete" class="danger-item">
                      <DeleteOutlined />
                      删除
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </a-space>
          </template>
        </a-table>
      </a-card>
    </div>

    <!-- 创建训练任务模态框 -->
    <a-modal
      v-model:open="createModalVisible"
      title="创建训练任务"
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
        <a-tabs v-model:activeKey="activeTab" type="card" class="create-tabs">
          <a-tab-pane key="basic" tab="基础配置">
            <a-row :gutter="16">
              <a-col :xs="24" :sm="12">
                <a-form-item label="任务名称" name="name">
                  <a-input
                    v-model:value="createForm.name"
                    placeholder="请输入任务名称"
                    class="form-input"
                  />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :sm="12">
                <a-form-item label="命名空间" name="namespace">
                  <a-select
                    v-model:value="createForm.namespace"
                    placeholder="选择命名空间"
                    class="form-select"
                  >
                    <a-select-option value="default">default</a-select-option>
                    <a-select-option value="ai-training"
                      >ai-training</a-select-option
                    >
                    <a-select-option value="research">research</a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>

            <a-row :gutter="16">
              <a-col :xs="24" :sm="12">
                <a-form-item label="训练框架" name="framework">
                  <a-select
                    v-model:value="createForm.framework"
                    placeholder="选择训练框架"
                    class="form-select"
                  >
                    <a-select-option value="tensorflow"
                      >TensorFlow</a-select-option
                    >
                    <a-select-option value="pytorch">PyTorch</a-select-option>
                    <a-select-option value="mindspore"
                      >MindSpore</a-select-option
                    >
                    <a-select-option value="paddlepaddle"
                      >PaddlePaddle</a-select-option
                    >
                  </a-select>
                </a-form-item>
              </a-col>
              <a-col :xs="24" :sm="12">
                <a-form-item label="任务类型" name="jobType">
                  <a-select
                    v-model:value="createForm.jobType"
                    placeholder="选择任务类型"
                    class="form-select"
                  >
                    <a-select-option value="single">单机训练</a-select-option>
                    <a-select-option value="distributed"
                      >分布式训练</a-select-option
                    >
                    <a-select-option value="horovod"
                      >Horovod训练</a-select-option
                    >
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>

            <a-form-item label="镜像" name="image">
              <a-input
                v-model:value="createForm.image"
                placeholder="请输入镜像地址"
                class="form-input"
              />
            </a-form-item>

            <a-form-item label="训练脚本" name="script">
              <a-textarea
                v-model:value="createForm.script"
                placeholder="请输入训练脚本路径或命令"
                :rows="3"
                class="form-textarea"
              />
            </a-form-item>
          </a-tab-pane>

          <a-tab-pane key="resources" tab="资源配置">
            <a-form-item label="工作节点数量" name="workers">
              <a-input-number
                v-model:value="createForm.workers"
                :min="1"
                :max="10"
                style="width: 100%"
                addon-after="个"
                class="form-input-number"
              />
            </a-form-item>

            <a-divider class="form-divider">单节点资源配置</a-divider>

            <a-row :gutter="16">
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
                    :min="4"
                    :max="512"
                    style="width: 100%"
                    addon-after="GB"
                    class="form-input-number"
                  />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :sm="8">
                <a-form-item label="GPU 卡数" name="gpu">
                  <a-input-number
                    v-model:value="createForm.gpu"
                    :min="0"
                    :max="8"
                    style="width: 100%"
                    addon-after="卡"
                    class="form-input-number"
                  />
                </a-form-item>
              </a-col>
            </a-row>

            <a-form-item label="存储卷" name="storage">
              <a-input-number
                v-model:value="createForm.storage"
                :min="50"
                :max="2000"
                style="width: 100%"
                addon-after="GB"
                class="form-input-number"
              />
            </a-form-item>
          </a-tab-pane>

          <a-tab-pane key="advanced" tab="高级配置">
            <a-form-item label="环境变量" name="envVars">
              <div class="env-vars-container">
                <div
                  v-for="(env, index) in createForm.envVars"
                  :key="index"
                  class="env-var-item"
                >
                  <a-input
                    v-model:value="env.key"
                    placeholder="变量名"
                    class="env-key-input"
                  />
                  <a-input
                    v-model:value="env.value"
                    placeholder="变量值"
                    class="env-value-input"
                  />
                  <a-button
                    type="text"
                    danger
                    @click="removeEnvVar(index)"
                    class="env-remove-btn"
                  >
                    <DeleteOutlined />
                  </a-button>
                </div>
                <a-button type="dashed" @click="addEnvVar" class="add-env-btn">
                  <PlusOutlined />
                  添加环境变量
                </a-button>
              </div>
            </a-form-item>

            <a-form-item label="数据集路径" name="dataPath">
              <a-input
                v-model:value="createForm.dataPath"
                placeholder="请输入数据集路径"
                class="form-input"
              />
            </a-form-item>

            <a-form-item label="输出路径" name="outputPath">
              <a-input
                v-model:value="createForm.outputPath"
                placeholder="请输入模型输出路径"
                class="form-input"
              />
            </a-form-item>

            <a-form-item label="描述" name="description">
              <a-textarea
                v-model:value="createForm.description"
                placeholder="请输入任务描述"
                :rows="4"
                class="form-textarea"
              />
            </a-form-item>
          </a-tab-pane>
        </a-tabs>
      </a-form>
    </a-modal>

    <!-- 详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="训练任务详情"
      width="1000px"
      :footer="null"
      class="sci-fi-modal detail-modal"
    >
      <div v-if="selectedJob" class="detail-content">
        <a-tabs type="card" class="detail-tabs">
          <a-tab-pane key="info" tab="基本信息">
            <a-descriptions
              :column="{ xs: 1, sm: 2 }"
              bordered
              class="detail-descriptions"
            >
              <a-descriptions-item label="任务名称">
                {{ selectedJob.name }}
              </a-descriptions-item>
              <a-descriptions-item label="状态">
                <a-tag
                  :color="getStatusColor(selectedJob.status)"
                  class="status-tag"
                >
                  <component :is="getStatusIcon(selectedJob.status)" />
                  {{ getStatusText(selectedJob.status) }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="训练框架">
                <div class="framework-wrapper">
                  <component
                    :is="getFrameworkIcon(selectedJob.framework)"
                    class="framework-icon"
                  />
                  <span>{{ getFrameworkText(selectedJob.framework) }}</span>
                </div>
              </a-descriptions-item>
              <a-descriptions-item label="任务类型">
                {{ getJobTypeText(selectedJob.jobType) }}
              </a-descriptions-item>
              <a-descriptions-item label="命名空间">
                {{ selectedJob.namespace }}
              </a-descriptions-item>
              <a-descriptions-item label="创建者">
                {{ selectedJob.creator }}
              </a-descriptions-item>
              <a-descriptions-item label="镜像">
                {{ selectedJob.image }}
              </a-descriptions-item>
              <a-descriptions-item label="创建时间">
                {{ selectedJob.createTime }}
              </a-descriptions-item>
              <a-descriptions-item label="工作节点">
                {{ selectedJob.workers }} 个
              </a-descriptions-item>
              <a-descriptions-item label="运行时间">
                {{ formatDuration(selectedJob.duration) }}
              </a-descriptions-item>
              <a-descriptions-item label="进度">
                <a-progress
                  :percent="selectedJob.progress"
                  :status="getProgressStatus(selectedJob.status)"
                  size="small"
                />
              </a-descriptions-item>
              <a-descriptions-item label="数据集路径" :span="2">
                {{ selectedJob.dataPath || '未设置' }}
              </a-descriptions-item>
              <a-descriptions-item label="输出路径" :span="2">
                {{ selectedJob.outputPath || '未设置' }}
              </a-descriptions-item>
              <a-descriptions-item label="描述" :span="2">
                {{ selectedJob.description || '暂无描述' }}
              </a-descriptions-item>
            </a-descriptions>
          </a-tab-pane>

          <a-tab-pane key="resources" tab="资源使用">
            <div class="resource-details">
              <a-row :gutter="24">
                <a-col :xs="24" :sm="8">
                  <a-card class="resource-card" :bordered="false">
                    <div class="resource-header">
                      <DatabaseOutlined class="resource-card-icon cpu" />
                      <span class="resource-title">CPU 使用率</span>
                    </div>
                    <div class="resource-value">
                      {{ selectedJob.resources.cpu }} 核
                    </div>
                    <a-progress :percent="75" size="small" :show-info="false" />
                  </a-card>
                </a-col>
                <a-col :xs="24" :sm="8">
                  <a-card class="resource-card" :bordered="false">
                    <div class="resource-header">
                      <ThunderboltOutlined class="resource-card-icon memory" />
                      <span class="resource-title">内存使用</span>
                    </div>
                    <div class="resource-value">
                      {{ selectedJob.resources.memory }} GB
                    </div>
                    <a-progress :percent="68" size="small" :show-info="false" />
                  </a-card>
                </a-col>
                <a-col :xs="24" :sm="8" v-if="selectedJob.resources.gpu">
                  <a-card class="resource-card" :bordered="false">
                    <div class="resource-header">
                      <BugOutlined class="resource-card-icon gpu" />
                      <span class="resource-title">GPU 使用率</span>
                    </div>
                    <div class="resource-value">
                      {{ selectedJob.resources.gpu }} 卡
                    </div>
                    <a-progress :percent="92" size="small" :show-info="false" />
                  </a-card>
                </a-col>
              </a-row>
            </div>
          </a-tab-pane>

          <a-tab-pane key="logs" tab="运行日志">
            <div class="log-container">
              <div class="log-header">
                <span class="log-title">训练日志</span>
                <a-space>
                  <a-select
                    v-model:value="selectedPod"
                    placeholder="选择Pod"
                    style="width: 200px"
                    class="pod-select"
                  >
                    <a-select-option value="worker-0">worker-0</a-select-option>
                    <a-select-option value="worker-1">worker-1</a-select-option>
                    <a-select-option value="parameter-server"
                      >parameter-server</a-select-option
                    >
                  </a-select>
                  <a-button
                    size="small"
                    @click="refreshLogs"
                    class="log-refresh-btn"
                  >
                    <ReloadOutlined />
                    刷新
                  </a-button>
                  <a-button
                    size="small"
                    @click="downloadLogs"
                    class="log-download-btn"
                  >
                    <DownloadOutlined />
                    下载
                  </a-button>
                </a-space>
              </div>
              <div class="log-content">
                <pre
                  v-for="(log, index) in trainingLogs"
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

    <!-- 克隆模态框 -->
    <a-modal
      v-model:open="cloneModalVisible"
      title="克隆训练任务"
      width="600px"
      :confirm-loading="cloneLoading"
      @ok="handleCloneSubmit"
      @cancel="handleCloneCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="cloneFormRef"
        :model="cloneForm"
        :rules="cloneFormRules"
        layout="vertical"
      >
        <a-form-item label="新任务名称" name="name">
          <a-input
            v-model:value="cloneForm.name"
            placeholder="请输入新任务名称"
            class="form-input"
          />
        </a-form-item>
        <a-form-item label="命名空间" name="namespace">
          <a-select
            v-model:value="cloneForm.namespace"
            placeholder="选择命名空间"
            class="form-select"
          >
            <a-select-option value="default">default</a-select-option>
            <a-select-option value="ai-training">ai-training</a-select-option>
            <a-select-option value="research">research</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import type { FormInstance, TableColumnsType } from 'ant-design-vue';
import {
  RocketOutlined,
  PlusOutlined,
  ReloadOutlined,
  DatabaseOutlined,
  ThunderboltOutlined,
  BugOutlined,
  EyeOutlined,
  MoreOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  CopyOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  StopOutlined,
  CloseCircleOutlined,
  FileTextOutlined,
  DownloadOutlined,
  ApiOutlined,
  CodeOutlined,
  CloudOutlined,
  RobotOutlined,
} from '@ant-design/icons-vue';

// ===== 类型定义 =====
interface TrainingJobResources {
  cpu: number;
  memory: number;
  gpu?: number;
}

interface EnvVar {
  key: string;
  value: string;
}

interface TrainingJob {
  id: string;
  name: string;
  namespace: string;
  status:
    | 'running'
    | 'pending'
    | 'completed'
    | 'failed'
    | 'cancelled'
    | 'paused';
  framework: 'tensorflow' | 'pytorch' | 'mindspore' | 'paddlepaddle';
  jobType: 'single' | 'distributed' | 'horovod';
  creator: string;
  image: string;
  createTime: string;
  resources: TrainingJobResources;
  workers: number;
  progress: number;
  duration: number; // 运行时间（分钟）
  script: string;
  dataPath?: string;
  outputPath?: string;
  description?: string;
}

interface CreateForm {
  name: string;
  namespace: string;
  framework: string;
  jobType: string;
  image: string;
  script: string;
  workers: number;
  cpu: number;
  memory: number;
  gpu: number;
  storage: number;
  envVars: EnvVar[];
  dataPath: string;
  outputPath: string;
  description: string;
}

interface CloneForm {
  name: string;
  namespace: string;
}

// ===== 响应式数据 =====
const loading = ref<boolean>(false);
const createModalVisible = ref<boolean>(false);
const detailModalVisible = ref<boolean>(false);
const cloneModalVisible = ref<boolean>(false);
const createLoading = ref<boolean>(false);
const cloneLoading = ref<boolean>(false);

const filterStatus = ref<string>('');
const filterFramework = ref<string>('');
const searchKeyword = ref<string>('');
const activeTab = ref<string>('basic');
const selectedPod = ref<string>('worker-0');

const selectedJob = ref<TrainingJob | null>(null);

const trainingLogs = ref<string[]>([
  '2024-06-23 10:30:15 INFO: Starting training job...',
  '2024-06-23 10:30:16 INFO: Loading dataset from /data/train',
  '2024-06-23 10:30:17 INFO: Model architecture initialized',
  '2024-06-23 10:30:18 INFO: Starting epoch 1/100',
  '2024-06-23 10:30:20 INFO: Epoch 1 - loss: 0.8567, accuracy: 0.7234',
  '2024-06-23 10:31:15 INFO: Epoch 2 - loss: 0.7234, accuracy: 0.7856',
  '2024-06-23 10:32:10 INFO: Epoch 3 - loss: 0.6789, accuracy: 0.8123',
]);

// ===== 表单引用 =====
const createFormRef = ref<FormInstance>();
const cloneFormRef = ref<FormInstance>();

// ===== 表单数据 =====
const createForm = reactive<CreateForm>({
  name: '',
  namespace: 'default',
  framework: 'tensorflow',
  jobType: 'single',
  image: 'tensorflow/tensorflow:2.13.0-gpu',
  script: 'python train.py',
  workers: 1,
  cpu: 4,
  memory: 8,
  gpu: 1,
  storage: 100,
  envVars: [{ key: '', value: '' }],
  dataPath: '/data/train',
  outputPath: '/data/output',
  description: '',
});

const cloneForm = reactive<CloneForm>({
  name: '',
  namespace: 'default',
});

// ===== 配置数据 =====
const STATUS_CONFIG = {
  running: { color: 'processing', text: '运行中', icon: PlayCircleOutlined },
  pending: { color: 'default', text: '等待中', icon: ClockCircleOutlined },
  completed: { color: 'success', text: '已完成', icon: CheckCircleOutlined },
  failed: { color: 'error', text: '失败', icon: CloseCircleOutlined },
  cancelled: { color: 'warning', text: '已取消', icon: StopOutlined },
  paused: { color: 'processing', text: '已暂停', icon: PauseCircleOutlined },
} as const;

const FRAMEWORK_CONFIG = {
  tensorflow: { text: 'TensorFlow', icon: ApiOutlined },
  pytorch: { text: 'PyTorch', icon: RobotOutlined },
  mindspore: { text: 'MindSpore', icon: CloudOutlined },
  paddlepaddle: { text: 'PaddlePaddle', icon: CodeOutlined },
} as const;

const JOB_TYPE_CONFIG = {
  single: '单机训练',
  distributed: '分布式训练',
  horovod: 'Horovod训练',
} as const;

// ===== 模拟数据 =====
const trainingJobs = ref<TrainingJob[]>([
  {
    id: 'job-001',
    name: 'resnet50-imagenet',
    namespace: 'ai-training',
    status: 'running',
    framework: 'tensorflow',
    jobType: 'distributed',
    creator: 'admin',
    image: 'tensorflow/tensorflow:2.13.0-gpu',
    createTime: '2024-06-23 09:00:00',
    resources: { cpu: 8, memory: 16, gpu: 2 },
    workers: 4,
    progress: 65,
    duration: 180,
    script: 'python train_resnet.py --epochs 100 --batch-size 64',
    dataPath: '/data/imagenet',
    outputPath: '/data/models/resnet50',
    description: 'ResNet-50 在 ImageNet 数据集上的训练',
  },
  {
    id: 'job-002',
    name: 'bert-fine-tuning',
    namespace: 'research',
    status: 'completed',
    framework: 'pytorch',
    jobType: 'single',
    creator: 'researcher',
    image: 'pytorch/pytorch:2.0.1-cuda11.7-cudnn8-devel',
    createTime: '2024-06-23 08:30:00',
    resources: { cpu: 4, memory: 8, gpu: 1 },
    workers: 1,
    progress: 100,
    duration: 120,
    script: 'python fine_tune_bert.py --model bert-base-uncased',
    dataPath: '/data/nlp/sentiment',
    outputPath: '/data/models/bert-sentiment',
    description: 'BERT 情感分析模型微调',
  },
  {
    id: 'job-003',
    name: 'yolo-object-detection',
    namespace: 'ai-training',
    status: 'pending',
    framework: 'pytorch',
    jobType: 'single',
    creator: 'developer',
    image: 'ultralytics/yolov8:latest',
    createTime: '2024-06-23 10:15:00',
    resources: { cpu: 6, memory: 12, gpu: 1 },
    workers: 1,
    progress: 0,
    duration: 0,
    script:
      'python train.py --data coco.yaml --epochs 100 --weights yolov8n.pt',
    dataPath: '/data/coco',
    outputPath: '/data/models/yolo',
    description: 'YOLOv8 目标检测模型训练',
  },
  {
    id: 'job-004',
    name: 'gan-face-generation',
    namespace: 'research',
    status: 'failed',
    framework: 'tensorflow',
    jobType: 'single',
    creator: 'researcher',
    image: 'tensorflow/tensorflow:2.13.0-gpu',
    createTime: '2024-06-23 07:45:00',
    resources: { cpu: 8, memory: 16, gpu: 2 },
    workers: 1,
    progress: 23,
    duration: 45,
    script: 'python train_gan.py --dataset faces --epochs 500',
    dataPath: '/data/faces',
    outputPath: '/data/models/gan',
    description: 'GAN 人脸生成模型训练',
  },
  {
    id: 'job-005',
    name: 'transformer-translation',
    namespace: 'ai-training',
    status: 'running',
    framework: 'pytorch',
    jobType: 'distributed',
    creator: 'admin',
    image: 'pytorch/pytorch:2.0.1-cuda11.7-cudnn8-devel',
    createTime: '2024-06-23 09:30:00',
    resources: { cpu: 16, memory: 32, gpu: 4 },
    workers: 2,
    progress: 35,
    duration: 90,
    script: 'python train_transformer.py --src-lang en --tgt-lang zh',
    dataPath: '/data/translation',
    outputPath: '/data/models/transformer',
    description: 'Transformer 英中翻译模型训练',
  },
]);

// ===== 表单验证规则 =====
const createFormRules = {
  name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
    {
      pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
      message: '名称只能包含小写字母、数字和连字符',
      trigger: 'blur',
    },
  ],
  namespace: [{ required: true, message: '请选择命名空间', trigger: 'change' }],
  framework: [{ required: true, message: '请选择训练框架', trigger: 'change' }],
  jobType: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  image: [{ required: true, message: '请输入镜像地址', trigger: 'blur' }],
  script: [{ required: true, message: '请输入训练脚本', trigger: 'blur' }],
  workers: [{ required: true, message: '请输入工作节点数量', trigger: 'blur' }],
  cpu: [{ required: true, message: '请输入 CPU 核数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
  storage: [{ required: true, message: '请输入存储大小', trigger: 'blur' }],
};

const cloneFormRules = {
  name: [
    { required: true, message: '请输入新任务名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
    {
      pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
      message: '名称只能包含小写字母、数字和连字符',
      trigger: 'blur',
    },
  ],
  namespace: [{ required: true, message: '请选择命名空间', trigger: 'change' }],
};

// ===== 表格列配置 =====
const columns: TableColumnsType<TrainingJob> = [
  {
    title: '任务名称',
    dataIndex: 'name',
    key: 'name',
    width: 200,
    ellipsis: true,
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
    slots: { customRender: 'status' },
  },
  {
    title: '训练框架',
    key: 'framework',
    width: 120,
    slots: { customRender: 'framework' },
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
    width: 100,
  },
  {
    title: '资源配置',
    key: 'resources',
    width: 200,
    slots: { customRender: 'resources' },
  },
  {
    title: '进度',
    key: 'progress',
    width: 120,
    slots: { customRender: 'progress' },
  },
  {
    title: '运行时间',
    key: 'duration',
    width: 120,
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
    width: 180,
    fixed: 'right',
    slots: { customRender: 'action' },
  },
];

// ===== 分页配置 =====
const paginationConfig = {
  total: computed(() => filteredTrainingJobs.value.length),
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// ===== 计算属性 =====
const filteredTrainingJobs = computed(() => {
  let result = trainingJobs.value;

  if (filterStatus.value) {
    result = result.filter((item) => item.status === filterStatus.value);
  }

  if (filterFramework.value) {
    result = result.filter((item) => item.framework === filterFramework.value);
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter(
      (item) =>
        item.name.toLowerCase().includes(keyword) ||
        item.creator.toLowerCase().includes(keyword),
    );
  }

  return result;
});

const runningCount = computed(
  () => trainingJobs.value.filter((job) => job.status === 'running').length,
);

const pendingCount = computed(
  () => trainingJobs.value.filter((job) => job.status === 'pending').length,
);

const completedCount = computed(
  () => trainingJobs.value.filter((job) => job.status === 'completed').length,
);

const failedCount = computed(
  () => trainingJobs.value.filter((job) => job.status === 'failed').length,
);

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

const getFrameworkIcon = (framework: string) => {
  return (
    FRAMEWORK_CONFIG[framework as keyof typeof FRAMEWORK_CONFIG]?.icon ||
    ApiOutlined
  );
};

const getFrameworkText = (framework: string): string => {
  return (
    FRAMEWORK_CONFIG[framework as keyof typeof FRAMEWORK_CONFIG]?.text ||
    framework
  );
};

const getJobTypeText = (jobType: string): string => {
  return JOB_TYPE_CONFIG[jobType as keyof typeof JOB_TYPE_CONFIG] || jobType;
};

const getProgressStatus = (
  status: string,
): 'success' | 'exception' | 'normal' => {
  if (status === 'completed') return 'success';
  if (status === 'failed') return 'exception';
  return 'normal';
};

const formatDuration = (minutes: number): string => {
  if (minutes === 0) return '0分钟';

  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;

  if (hours === 0) {
    return `${remainingMinutes}分钟`;
  } else if (remainingMinutes === 0) {
    return `${hours}小时`;
  } else {
    return `${hours}小时${remainingMinutes}分钟`;
  }
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

// ===== 事件处理函数 =====
const showCreateModal = (): void => {
  createModalVisible.value = true;
  activeTab.value = 'basic';
};

const handleCreateSubmit = async (): Promise<void> => {
  try {
    await createFormRef.value?.validate();
    createLoading.value = true;

    // 模拟 API 调用
    await new Promise((resolve) => setTimeout(resolve, 2000));

    const newJob: TrainingJob = {
      id: `job-${Date.now()}`,
      name: createForm.name,
      namespace: createForm.namespace,
      status: 'pending',
      framework: createForm.framework as any,
      jobType: createForm.jobType as any,
      creator: 'current-user',
      image: createForm.image,
      createTime: new Date().toLocaleString(),
      resources: {
        cpu: createForm.cpu,
        memory: createForm.memory,
        ...(createForm.gpu > 0 && { gpu: createForm.gpu }),
      },
      workers: createForm.workers,
      progress: 0,
      duration: 0,
      script: createForm.script,
      dataPath: createForm.dataPath,
      outputPath: createForm.outputPath,
      description: createForm.description,
    };

    trainingJobs.value.unshift(newJob);
    createModalVisible.value = false;
    message.success('训练任务创建成功');

    // 重置表单
    createFormRef.value?.resetFields();
    createForm.envVars = [{ key: '', value: '' }];
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    createLoading.value = false;
  }
};

const handleCreateCancel = (): void => {
  createModalVisible.value = false;
  createFormRef.value?.resetFields();
  createForm.envVars = [{ key: '', value: '' }];
};

const addEnvVar = (): void => {
  createForm.envVars.push({ key: '', value: '' });
};

const removeEnvVar = (index: number): void => {
  createForm.envVars.splice(index, 1);
};

const viewLogs = (record: TrainingJob): void => {
  selectedJob.value = record;
  detailModalVisible.value = true;
  // 在详情模态框中切换到日志tab
  setTimeout(() => {
    const tabElement = document.querySelector(
      '[data-node-key="logs"]',
    ) as HTMLElement;
    if (tabElement) {
      tabElement.click();
    }
  }, 100);
};

const viewDetails = (record: TrainingJob): void => {
  selectedJob.value = record;
  detailModalVisible.value = true;
};

const handleMenuAction = (key: string, record: TrainingJob): void => {
  const actions = {
    pause: () => handlePause(record),
    resume: () => handleResume(record),
    stop: () => handleStop(record),
    clone: () => handleClone(record),
    delete: () => handleDelete(record),
  };

  const action = actions[key as keyof typeof actions];
  if (action) {
    action();
  }
};

const handlePause = async (record: TrainingJob): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const index = trainingJobs.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      trainingJobs.value[index]!.status = 'paused';
    }
    message.success('训练任务暂停成功');
  } catch (error) {
    message.error('暂停失败');
  } finally {
    loading.value = false;
  }
};

const handleResume = async (record: TrainingJob): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const index = trainingJobs.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      trainingJobs.value[index]!.status = 'running';
    }
    message.success('训练任务恢复成功');
  } catch (error) {
    message.error('恢复失败');
  } finally {
    loading.value = false;
  }
};

const handleStop = async (record: TrainingJob): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const index = trainingJobs.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      trainingJobs.value[index]!.status = 'cancelled';
    }
    message.success('训练任务停止成功');
  } catch (error) {
    message.error('停止失败');
  } finally {
    loading.value = false;
  }
};

const handleClone = (record: TrainingJob): void => {
  cloneForm.name = `${record.name}-copy`;
  cloneForm.namespace = record.namespace;
  selectedJob.value = record;
  cloneModalVisible.value = true;
};

const handleCloneSubmit = async (): Promise<void> => {
  try {
    await cloneFormRef.value?.validate();
    cloneLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 1500));

    if (selectedJob.value) {
      const clonedJob: TrainingJob = {
        ...selectedJob.value,
        id: `job-${Date.now()}`,
        name: cloneForm.name,
        namespace: cloneForm.namespace,
        status: 'pending',
        createTime: new Date().toLocaleString(),
        progress: 0,
        duration: 0,
      };

      trainingJobs.value.unshift(clonedJob);
      cloneModalVisible.value = false;
      message.success('训练任务克隆成功');
    }
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    cloneLoading.value = false;
  }
};

const handleCloneCancel = (): void => {
  cloneModalVisible.value = false;
  cloneFormRef.value?.resetFields();
};

const handleDelete = (record: TrainingJob): void => {
  const deleteConfirm = () => {
    const index = trainingJobs.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      trainingJobs.value.splice(index, 1);
      message.success('训练任务删除成功');
    }
  };

  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除训练任务 "${record.name}" 吗？此操作不可恢复。`,
      okText: '确认',
      cancelText: '取消',
      type: 'warning',
      onOk: deleteConfirm,
    });
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
    ...trainingLogs.value,
    `${new Date().toLocaleString()} INFO: Log refreshed for ${selectedPod.value}`,
  ];
  trainingLogs.value = newLogs.slice(-50);
};

const downloadLogs = (): void => {
  const logContent = trainingLogs.value.join('\n');
  const blob = new Blob([logContent], { type: 'text/plain' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `training-logs-${selectedPod.value}-${Date.now()}.txt`;
  a.click();
  URL.revokeObjectURL(url);
  message.success('日志下载成功');
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
.training-container {
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

/* ===== 统计卡片样式 ===== */
.stats-section {
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 8px !important;
  transition: all 0.3s ease;
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.stat-icon.running {
  background: linear-gradient(135deg, #52c41a, #389e0d);
}

.stat-icon.pending {
  background: linear-gradient(135deg, #1890ff, #096dd9);
}

.stat-icon.completed {
  background: linear-gradient(135deg, #13c2c2, #08979c);
}

.stat-icon.failed {
  background: linear-gradient(135deg, #ff4d4f, #cf1322);
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.65;
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

.indicator-running {
  background: #52c41a;
}

.indicator-pending {
  background: #1890ff;
}

.indicator-completed {
  background: #13c2c2;
}

.indicator-failed {
  background: #ff4d4f;
}

.indicator-cancelled {
  background: #faad14;
}

.indicator-paused {
  background: #722ed1;
}

/* ===== 框架显示 ===== */
.framework-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.framework-icon {
  font-size: 16px;
  color: #1890ff;
}

.framework-text {
  font-weight: 500;
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

/* ===== 进度条 ===== */
.progress-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.progress-bar {
  flex: 1;
}

.progress-text {
  font-size: 12px;
  font-weight: 500;
  min-width: 35px;
}

/* ===== 运行时间 ===== */
.duration-wrapper {
  display: flex;
  align-items: center;
  gap: 6px;
}

.duration-icon {
  font-size: 12px;
  color: #1890ff;
}

.duration-text {
  font-size: 12px;
  font-weight: 500;
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

.create-tabs :deep(.ant-tabs-tab) {
  border-radius: 6px 6px 0 0 !important;
}

/* ===== 环境变量样式 ===== */
.env-vars-container {
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  padding: 16px;
}

.env-var-item {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  align-items: center;
}

.env-key-input,
.env-value-input {
  flex: 1;
}

.env-remove-btn {
  border: none !important;
  background: transparent !important;
  color: #ff4d4f !important;
}

.add-env-btn {
  width: 100%;
  border-radius: 6px !important;
}

/* ===== 详情页样式 ===== */
.detail-content {
  max-height: 700px;
  overflow-y: auto;
}

.detail-tabs :deep(.ant-tabs-tab) {
  border-radius: 6px 6px 0 0 !important;
}

.detail-descriptions {
  margin-bottom: 16px;
}

/* ===== 资源详情卡片 ===== */
.resource-details {
  margin: 16px 0;
}

.resource-card {
  text-align: center;
  border-radius: 8px !important;
  transition: all 0.3s ease;
}

.resource-card:hover {
  transform: translateY(-2px);
}

.resource-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 12px;
}

.resource-card-icon {
  font-size: 20px;
}

.resource-card-icon.cpu {
  color: #52c41a;
}

.resource-card-icon.memory {
  color: #1890ff;
}

.resource-card-icon.gpu {
  color: #722ed1;
}

.resource-title {
  font-weight: 600;
  font-size: 14px;
}

.resource-value {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 8px;
}

/* ===== 日志容器 ===== */
.log-container {
  margin: 16px 0;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.log-title {
  font-weight: 600;
  font-size: 16px;
}

.pod-select,
.log-refresh-btn,
.log-download-btn {
  border-radius: 6px !important;
}

.log-refresh-btn:hover,
.log-download-btn:hover {
  color: #1890ff !important;
  border-color: #1890ff !important;
}

.log-content {
  border-radius: 6px !important;
  padding: 16px !important;
  max-height: 400px;
  overflow-y: auto;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace !important;
  border: 1px solid #d9d9d9;
}

.log-line {
  margin: 0;
  font-size: 12px;
  line-height: 1.5;
}

/* ===== 响应式设计 ===== */
@media (max-width: 768px) {
  .training-container {
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

  .sci-fi-modal :deep(.ant-modal) {
    margin: 16px !important;
    max-width: calc(100vw - 32px) !important;
  }

  .stat-content {
    gap: 12px;
  }

  .stat-icon {
    width: 40px;
    height: 40px;
    font-size: 20px;
  }

  .stat-number {
    font-size: 20px;
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

  .env-var-item {
    flex-direction: column;
    gap: 8px;
  }

  .env-key-input,
  .env-value-input {
    width: 100%;
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
