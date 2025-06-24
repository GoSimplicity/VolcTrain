<template>
  <div class="job-queue-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <ThunderboltOutlined class="title-icon" />
            <span class="title-text">训练任务队列</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text">
              管理和监控基于 Volcano 的 AI 训练任务
            </span>
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
      <a-row :gutter="16">
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <div class="stat-content">
              <div class="stat-icon running-icon">
                <PlayCircleOutlined />
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ runningCount }}</div>
                <div class="stat-label">运行中</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <div class="stat-content">
              <div class="stat-icon pending-icon">
                <ClockCircleOutlined />
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ pendingCount }}</div>
                <div class="stat-label">排队中</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <div class="stat-content">
              <div class="stat-icon completed-icon">
                <CheckCircleOutlined />
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ completedCount }}</div>
                <div class="stat-label">已完成</div>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <div class="stat-content">
              <div class="stat-icon failed-icon">
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
              <a-select-option value="pending">排队中</a-select-option>
              <a-select-option value="completed">已完成</a-select-option>
              <a-select-option value="failed">失败</a-select-option>
              <a-select-option value="cancelled">已取消</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterQueue"
              placeholder="选择队列"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部队列</a-select-option>
              <a-select-option value="default">default</a-select-option>
              <a-select-option value="high-priority"
                >high-priority</a-select-option
              >
              <a-select-option value="gpu-queue">gpu-queue</a-select-option>
              <a-select-option value="cpu-queue">cpu-queue</a-select-option>
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
          :data-source="filteredJobs"
          :loading="loading"
          :pagination="paginationConfig"
          row-key="id"
          size="middle"
          :scroll="{ x: 'max-content' }"
          class="sci-fi-table"
          :row-class-name="getRowClassName"
        >
          <!-- 任务名称列 -->
          <template #name="{ record }">
            <div class="job-name-wrapper">
              <div class="job-name">{{ record.name }}</div>
              <div class="job-id">ID: {{ record.id }}</div>
            </div>
          </template>

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

          <!-- 优先级列 -->
          <template #priority="{ record }">
            <a-tag
              :color="getPriorityColor(record.priority)"
              class="priority-tag"
            >
              {{ getPriorityText(record.priority) }}
            </a-tag>
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
                class="job-progress"
              />
              <span class="progress-text">{{ record.progress }}%</span>
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

          <!-- 运行时间列 -->
          <template #duration="{ record }">
            <span class="duration-text">{{
              formatDuration(record.duration)
            }}</span>
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
              <a-button
                type="link"
                size="small"
                @click="viewLogs(record)"
                class="action-btn"
              >
                <FileTextOutlined />
                日志
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
                        !['running', 'pending', 'paused'].includes(
                          record.status,
                        )
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
            <a-form-item label="队列" name="queue">
              <a-select
                v-model:value="createForm.queue"
                placeholder="选择队列"
                class="form-select"
              >
                <a-select-option value="default">default</a-select-option>
                <a-select-option value="high-priority"
                  >high-priority</a-select-option
                >
                <a-select-option value="gpu-queue">gpu-queue</a-select-option>
                <a-select-option value="cpu-queue">cpu-queue</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="优先级" name="priority">
              <a-select
                v-model:value="createForm.priority"
                placeholder="选择优先级"
                class="form-select"
              >
                <a-select-option value="low">低</a-select-option>
                <a-select-option value="medium">中</a-select-option>
                <a-select-option value="high">高</a-select-option>
                <a-select-option value="urgent">紧急</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="训练镜像" name="image">
          <a-select
            v-model:value="createForm.image"
            placeholder="选择训练镜像"
            class="form-select"
          >
            <a-select-option
              value="pytorch/pytorch:1.12.0-cuda11.3-cudnn8-runtime"
            >
              PyTorch 1.12.0
            </a-select-option>
            <a-select-option value="tensorflow/tensorflow:2.9.0-gpu">
              TensorFlow 2.9.0
            </a-select-option>
            <a-select-option value="nvcr.io/nvidia/pytorch:22.05-py3">
              NVIDIA PyTorch 22.05
            </a-select-option>
            <a-select-option value="custom">自定义镜像</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item
          v-if="createForm.image === 'custom'"
          label="自定义镜像地址"
          name="customImage"
        >
          <a-input
            v-model:value="createForm.customImage"
            placeholder="请输入镜像地址"
            class="form-input"
          />
        </a-form-item>

        <a-divider class="form-divider">资源配置</a-divider>

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
                :min="2"
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

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="副本数" name="replicas">
              <a-input-number
                v-model:value="createForm.replicas"
                :min="1"
                :max="10"
                style="width: 100%"
                addon-after="个"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="最大重试次数" name="maxRetries">
              <a-input-number
                v-model:value="createForm.maxRetries"
                :min="0"
                :max="10"
                style="width: 100%"
                addon-after="次"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider class="form-divider">训练配置</a-divider>

        <a-form-item label="训练脚本" name="script">
          <a-textarea
            v-model:value="createForm.script"
            placeholder="请输入训练脚本或命令"
            :rows="4"
            class="form-textarea"
          />
        </a-form-item>

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
                class="env-key"
              />
              <a-input
                v-model:value="env.value"
                placeholder="变量值"
                class="env-value"
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

        <a-form-item label="数据集挂载路径" name="dataPath">
          <a-input
            v-model:value="createForm.dataPath"
            placeholder="例如: /data/dataset"
            class="form-input"
          />
        </a-form-item>

        <a-form-item label="模型输出路径" name="outputPath">
          <a-input
            v-model:value="createForm.outputPath"
            placeholder="例如: /output/models"
            class="form-input"
          />
        </a-form-item>

        <a-form-item label="描述" name="description">
          <a-textarea
            v-model:value="createForm.description"
            placeholder="请输入任务描述"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 任务详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="任务详情"
      width="1000px"
      :footer="null"
      class="sci-fi-modal detail-modal"
    >
      <div v-if="selectedJob" class="detail-content">
        <a-tabs default-active-key="overview" class="detail-tabs">
          <a-tab-pane key="overview" tab="概览">
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
              <a-descriptions-item label="队列">
                {{ selectedJob.queue }}
              </a-descriptions-item>
              <a-descriptions-item label="优先级">
                <a-tag
                  :color="getPriorityColor(selectedJob.priority)"
                  class="priority-tag"
                >
                  {{ getPriorityText(selectedJob.priority) }}
                </a-tag>
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
              <a-descriptions-item label="CPU">
                {{ selectedJob.resources.cpu }} 核
              </a-descriptions-item>
              <a-descriptions-item label="内存">
                {{ selectedJob.resources.memory }} GB
              </a-descriptions-item>
              <a-descriptions-item label="GPU" v-if="selectedJob.resources.gpu">
                {{ selectedJob.resources.gpu }} 卡
              </a-descriptions-item>
              <a-descriptions-item label="副本数">
                {{ selectedJob.replicas }}
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
              <a-descriptions-item label="描述" :span="2">
                {{ selectedJob.description || '暂无描述' }}
              </a-descriptions-item>
            </a-descriptions>
          </a-tab-pane>

          <a-tab-pane key="pods" tab="Pod 状态">
            <a-table
              :columns="podColumns"
              :data-source="mockPods"
              :pagination="false"
              size="small"
              class="pod-table"
            >
              <template #podStatus="{ record }">
                <a-tag
                  :color="getPodStatusColor(record.status)"
                  class="status-tag"
                >
                  {{ record.status }}
                </a-tag>
              </template>
              <template #resources="{ record }">
                <div class="pod-resources">
                  <div>CPU: {{ record.cpu }}</div>
                  <div>内存: {{ record.memory }}</div>
                  <div v-if="record.gpu">GPU: {{ record.gpu }}</div>
                </div>
              </template>
            </a-table>
          </a-tab-pane>

          <a-tab-pane key="events" tab="事件">
            <a-timeline class="event-timeline">
              <a-timeline-item
                v-for="event in mockEvents"
                :key="event.id"
                :color="getEventColor(event.type)"
              >
                <div class="event-item">
                  <div class="event-header">
                    <span class="event-type">{{ event.type }}</span>
                    <span class="event-time">{{ event.time }}</span>
                  </div>
                  <div class="event-message">{{ event.message }}</div>
                </div>
              </a-timeline-item>
            </a-timeline>
          </a-tab-pane>
        </a-tabs>
      </div>
    </a-modal>

    <!-- 日志查看模态框 -->
    <a-modal
      v-model:open="logModalVisible"
      title="任务日志"
      width="1000px"
      :footer="null"
      class="sci-fi-modal log-modal"
    >
      <div class="log-container">
        <div class="log-header">
          <a-space>
            <a-select
              v-model:value="selectedPod"
              placeholder="选择 Pod"
              style="width: 200px"
              class="pod-select"
            >
              <a-select-option
                v-for="pod in mockPods"
                :key="pod.name"
                :value="pod.name"
              >
                {{ pod.name }}
              </a-select-option>
            </a-select>
            <a-button @click="refreshLogs" class="log-refresh-btn">
              <ReloadOutlined />
              刷新
            </a-button>
            <a-button @click="downloadLogs" class="log-download-btn">
              <DownloadOutlined />
              下载
            </a-button>
          </a-space>
        </div>
        <div class="log-content">
          <pre
            v-for="(log, index) in logs"
            :key="index"
            class="log-line"
            :class="getLogLineClass(log)"
            >{{ log }}</pre
          >
        </div>
      </div>
    </a-modal>

    <!-- 克隆任务模态框 -->
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
        <a-form-item label="队列" name="queue">
          <a-select
            v-model:value="cloneForm.queue"
            placeholder="选择队列"
            class="form-select"
          >
            <a-select-option value="default">default</a-select-option>
            <a-select-option value="high-priority"
              >high-priority</a-select-option
            >
            <a-select-option value="gpu-queue">gpu-queue</a-select-option>
            <a-select-option value="cpu-queue">cpu-queue</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="优先级" name="priority">
          <a-select
            v-model:value="cloneForm.priority"
            placeholder="选择优先级"
            class="form-select"
          >
            <a-select-option value="low">低</a-select-option>
            <a-select-option value="medium">中</a-select-option>
            <a-select-option value="high">高</a-select-option>
            <a-select-option value="urgent">紧急</a-select-option>
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
  ThunderboltOutlined,
  PlusOutlined,
  ReloadOutlined,
  DatabaseOutlined,
  BugOutlined,
  EyeOutlined,
  MoreOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  StopOutlined,
  CopyOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  CloseCircleOutlined,
  FileTextOutlined,
  DownloadOutlined,
} from '@ant-design/icons-vue';

// ===== 类型定义 =====
interface JobResources {
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
  queue: string;
  status:
    | 'running'
    | 'pending'
    | 'completed'
    | 'failed'
    | 'cancelled'
    | 'paused';
  priority: 'low' | 'medium' | 'high' | 'urgent';
  creator: string;
  image: string;
  createTime: string;
  duration: number; // 秒
  progress: number;
  resources: JobResources;
  replicas: number;
  maxRetries: number;
  script: string;
  dataPath?: string;
  outputPath?: string;
  description?: string;
}

interface CreateForm {
  name: string;
  namespace: string;
  queue: string;
  priority: string;
  image: string;
  customImage: string;
  cpu: number;
  memory: number;
  gpu: number;
  replicas: number;
  maxRetries: number;
  script: string;
  envVars: EnvVar[];
  dataPath: string;
  outputPath: string;
  description: string;
}

interface CloneForm {
  name: string;
  queue: string;
  priority: string;
}

interface PodInfo {
  name: string;
  status: string;
  cpu: string;
  memory: string;
  gpu?: string;
  node: string;
  startTime: string;
}

interface EventInfo {
  id: string;
  type: string;
  time: string;
  message: string;
}

// ===== 响应式数据 =====
const loading = ref<boolean>(false);
const createModalVisible = ref<boolean>(false);
const detailModalVisible = ref<boolean>(false);
const logModalVisible = ref<boolean>(false);
const cloneModalVisible = ref<boolean>(false);
const createLoading = ref<boolean>(false);
const cloneLoading = ref<boolean>(false);

const filterStatus = ref<string>('');
const filterQueue = ref<string>('');
const searchKeyword = ref<string>('');
const selectedPod = ref<string>('');

const selectedJob = ref<TrainingJob | null>(null);

// ===== 表单引用 =====
const createFormRef = ref<FormInstance>();
const cloneFormRef = ref<FormInstance>();

// ===== 表单数据 =====
const createForm = reactive<CreateForm>({
  name: '',
  namespace: 'default',
  queue: 'default',
  priority: 'medium',
  image: 'pytorch/pytorch:1.12.0-cuda11.3-cudnn8-runtime',
  customImage: '',
  cpu: 4,
  memory: 8,
  gpu: 1,
  replicas: 1,
  maxRetries: 3,
  script: '',
  envVars: [],
  dataPath: '',
  outputPath: '',
  description: '',
});

const cloneForm = reactive<CloneForm>({
  name: '',
  queue: 'default',
  priority: 'medium',
});

// ===== 配置数据 =====
const STATUS_CONFIG = {
  running: { color: 'success', text: '运行中', icon: PlayCircleOutlined },
  pending: { color: 'processing', text: '排队中', icon: ClockCircleOutlined },
  completed: { color: 'success', text: '已完成', icon: CheckCircleOutlined },
  failed: { color: 'error', text: '失败', icon: CloseCircleOutlined },
  cancelled: { color: 'default', text: '已取消', icon: StopOutlined },
  paused: { color: 'warning', text: '已暂停', icon: PauseCircleOutlined },
} as const;

const PRIORITY_CONFIG = {
  low: { color: 'default', text: '低' },
  medium: { color: 'processing', text: '中' },
  high: { color: 'warning', text: '高' },
  urgent: { color: 'error', text: '紧急' },
} as const;

// ===== 模拟数据 =====
const jobs = ref<TrainingJob[]>([
  {
    id: 'job-001',
    name: 'bert-training-task',
    namespace: 'ai-training',
    queue: 'gpu-queue',
    status: 'running',
    priority: 'high',
    creator: 'admin',
    image: 'pytorch/pytorch:1.12.0-cuda11.3-cudnn8-runtime',
    createTime: '2024-06-23 09:30:00',
    duration: 3600,
    progress: 65,
    resources: { cpu: 8, memory: 32, gpu: 4 },
    replicas: 2,
    maxRetries: 3,
    script: 'python train.py --model bert --epochs 100',
    dataPath: '/data/nlp-dataset',
    outputPath: '/output/bert-model',
    description: 'BERT 模型训练任务',
  },
  {
    id: 'job-002',
    name: 'resnet-image-classification',
    namespace: 'research',
    queue: 'default',
    status: 'pending',
    priority: 'medium',
    creator: 'researcher',
    image: 'tensorflow/tensorflow:2.9.0-gpu',
    createTime: '2024-06-23 10:15:00',
    duration: 0,
    progress: 0,
    resources: { cpu: 4, memory: 16, gpu: 2 },
    replicas: 1,
    maxRetries: 5,
    script: 'python train_resnet.py --dataset imagenet',
    dataPath: '/data/imagenet',
    outputPath: '/output/resnet',
    description: 'ResNet 图像分类模型训练',
  },
  {
    id: 'job-003',
    name: 'llama-fine-tuning',
    namespace: 'ai-training',
    queue: 'high-priority',
    status: 'completed',
    priority: 'urgent',
    creator: 'ml-engineer',
    image: 'nvcr.io/nvidia/pytorch:22.05-py3',
    createTime: '2024-06-22 14:20:00',
    duration: 7200,
    progress: 100,
    resources: { cpu: 16, memory: 64, gpu: 8 },
    replicas: 4,
    maxRetries: 2,
    script: 'python finetune_llama.py --model llama-7b',
    dataPath: '/data/text-corpus',
    outputPath: '/output/llama-ft',
    description: 'LLaMA 大模型微调任务',
  },
  {
    id: 'job-004',
    name: 'yolo-object-detection',
    namespace: 'default',
    queue: 'cpu-queue',
    status: 'failed',
    priority: 'low',
    creator: 'developer',
    image: 'pytorch/pytorch:1.12.0-cuda11.3-cudnn8-runtime',
    createTime: '2024-06-23 08:45:00',
    duration: 1800,
    progress: 25,
    resources: { cpu: 8, memory: 16, gpu: 2 },
    replicas: 1,
    maxRetries: 3,
    script: 'python train_yolo.py --config yolo_config.yaml',
    dataPath: '/data/coco-dataset',
    outputPath: '/output/yolo',
    description: 'YOLO 目标检测模型训练',
  },
  {
    id: 'job-005',
    name: 'transformer-translation',
    namespace: 'research',
    queue: 'gpu-queue',
    status: 'paused',
    priority: 'medium',
    creator: 'researcher',
    image: 'tensorflow/tensorflow:2.9.0-gpu',
    createTime: '2024-06-23 07:30:00',
    duration: 2700,
    progress: 45,
    resources: { cpu: 6, memory: 24, gpu: 3 },
    replicas: 2,
    maxRetries: 4,
    script: 'python train_transformer.py --task translation',
    dataPath: '/data/translation-corpus',
    outputPath: '/output/transformer',
    description: 'Transformer 机器翻译模型训练',
  },
]);

const mockPods: PodInfo[] = [
  {
    name: 'bert-training-task-worker-0',
    status: 'Running',
    cpu: '4 cores',
    memory: '16 GB',
    gpu: '2 cards',
    node: 'gpu-node-01',
    startTime: '2024-06-23 09:31:00',
  },
  {
    name: 'bert-training-task-worker-1',
    status: 'Running',
    cpu: '4 cores',
    memory: '16 GB',
    gpu: '2 cards',
    node: 'gpu-node-02',
    startTime: '2024-06-23 09:31:00',
  },
];

const mockEvents: EventInfo[] = [
  {
    id: 'event-1',
    type: 'Normal',
    time: '2024-06-23 09:30:00',
    message: '任务已创建',
  },
  {
    id: 'event-2',
    type: 'Normal',
    time: '2024-06-23 09:30:15',
    message: '任务已加入队列',
  },
  {
    id: 'event-3',
    type: 'Normal',
    time: '2024-06-23 09:31:00',
    message: 'Pod 开始调度',
  },
  {
    id: 'event-4',
    type: 'Normal',
    time: '2024-06-23 09:31:30',
    message: '训练开始执行',
  },
];

const logs = ref<string[]>([
  '2024-06-23 09:31:00 INFO: 初始化训练环境...',
  '2024-06-23 09:31:05 INFO: 加载数据集...',
  '2024-06-23 09:31:10 INFO: 模型初始化完成',
  '2024-06-23 09:31:15 INFO: 开始训练 Epoch 1/100',
  '2024-06-23 09:32:00 INFO: Epoch 1 完成，损失: 0.85',
  '2024-06-23 09:32:05 INFO: 开始训练 Epoch 2/100',
  '2024-06-23 09:33:00 INFO: Epoch 2 完成，损失: 0.78',
  '2024-06-23 09:33:05 INFO: 开始训练 Epoch 3/100',
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
  queue: [{ required: true, message: '请选择队列', trigger: 'change' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }],
  image: [{ required: true, message: '请选择镜像', trigger: 'change' }],
  customImage: [
    { required: true, message: '请输入自定义镜像地址', trigger: 'blur' },
  ],
  cpu: [{ required: true, message: '请输入 CPU 核数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
  replicas: [{ required: true, message: '请输入副本数', trigger: 'blur' }],
  script: [{ required: true, message: '请输入训练脚本', trigger: 'blur' }],
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
  queue: [{ required: true, message: '请选择队列', trigger: 'change' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }],
};

// ===== 表格列配置 =====
const columns: TableColumnsType<TrainingJob> = [
  {
    title: '任务名称',
    key: 'name',
    width: 200,
    ellipsis: true,
    slots: { customRender: 'name' },
  },
  {
    title: '队列',
    dataIndex: 'queue',
    key: 'queue',
    width: 120,
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
    slots: { customRender: 'status' },
  },
  {
    title: '优先级',
    key: 'priority',
    width: 100,
    slots: { customRender: 'priority' },
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
    width: 180,
    slots: { customRender: 'resources' },
  },
  {
    title: '进度',
    key: 'progress',
    width: 120,
    slots: { customRender: 'progress' },
  },
  {
    title: '创建时间',
    key: 'createTime',
    width: 150,
    slots: { customRender: 'createTime' },
  },
  {
    title: '运行时间',
    key: 'duration',
    width: 120,
    slots: { customRender: 'duration' },
  },
  {
    title: '操作',
    key: 'action',
    width: 200,
    fixed: 'right',
    slots: { customRender: 'action' },
  },
];

const podColumns: TableColumnsType<PodInfo> = [
  {
    title: 'Pod 名称',
    dataIndex: 'name',
    key: 'name',
    width: 200,
  },
  {
    title: '状态',
    key: 'podStatus',
    width: 100,
    slots: { customRender: 'podStatus' },
  },
  {
    title: '资源',
    key: 'resources',
    width: 150,
    slots: { customRender: 'resources' },
  },
  {
    title: '节点',
    dataIndex: 'node',
    key: 'node',
    width: 120,
  },
  {
    title: '启动时间',
    dataIndex: 'startTime',
    key: 'startTime',
    width: 150,
  },
];

// ===== 分页配置 =====
const paginationConfig = {
  total: computed(() => filteredJobs.value.length),
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// ===== 计算属性 =====
const filteredJobs = computed(() => {
  let result = jobs.value;

  if (filterStatus.value) {
    result = result.filter((item) => item.status === filterStatus.value);
  }

  if (filterQueue.value) {
    result = result.filter((item) => item.queue === filterQueue.value);
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
  () => jobs.value.filter((job) => job.status === 'running').length,
);

const pendingCount = computed(
  () => jobs.value.filter((job) => job.status === 'pending').length,
);

const completedCount = computed(
  () => jobs.value.filter((job) => job.status === 'completed').length,
);

const failedCount = computed(
  () => jobs.value.filter((job) => job.status === 'failed').length,
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

const getPriorityColor = (priority: string): string => {
  return (
    PRIORITY_CONFIG[priority as keyof typeof PRIORITY_CONFIG]?.color ||
    'default'
  );
};

const getPriorityText = (priority: string): string => {
  return (
    PRIORITY_CONFIG[priority as keyof typeof PRIORITY_CONFIG]?.text || priority
  );
};

const getProgressStatus = (status: string): string => {
  if (status === 'failed') return 'exception';
  if (status === 'completed') return 'success';
  return 'normal';
};

const getPodStatusColor = (status: string): string => {
  const colorMap: Record<string, string> = {
    Running: 'success',
    Pending: 'processing',
    Failed: 'error',
    Succeeded: 'success',
  };
  return colorMap[status] || 'default';
};

const getEventColor = (type: string): string => {
  const colorMap: Record<string, string> = {
    Normal: 'blue',
    Warning: 'orange',
    Error: 'red',
  };
  return colorMap[type] || 'blue';
};

const getLogLineClass = (log: string): string => {
  if (log.includes('ERROR')) return 'log-error';
  if (log.includes('WARNING')) return 'log-warning';
  if (log.includes('INFO')) return 'log-info';
  return '';
};

const getRowClassName = (record: TrainingJob): string => {
  return `row-${record.status}`;
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

const formatDuration = (seconds: number): string => {
  if (seconds === 0) return '0秒';

  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  if (hours > 0) {
    return `${hours}小时${minutes}分钟`;
  } else if (minutes > 0) {
    return `${minutes}分钟${secs}秒`;
  } else {
    return `${secs}秒`;
  }
};

// ===== 环境变量管理 =====
const addEnvVar = (): void => {
  createForm.envVars.push({ key: '', value: '' });
};

const removeEnvVar = (index: number): void => {
  createForm.envVars.splice(index, 1);
};

// ===== 事件处理函数 =====
const showCreateModal = (): void => {
  createModalVisible.value = true;
};

const handleCreateSubmit = async (): Promise<void> => {
  try {
    await createFormRef.value?.validate();
    createLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 2000));

    const newJob: TrainingJob = {
      id: `job-${Date.now()}`,
      name: createForm.name,
      namespace: createForm.namespace,
      queue: createForm.queue,
      status: 'pending',
      priority: createForm.priority as 'low' | 'medium' | 'high' | 'urgent',
      creator: 'current-user',
      image:
        createForm.image === 'custom'
          ? createForm.customImage
          : createForm.image,
      createTime: new Date().toLocaleString(),
      duration: 0,
      progress: 0,
      resources: {
        cpu: createForm.cpu,
        memory: createForm.memory,
        ...(createForm.gpu > 0 && { gpu: createForm.gpu }),
      },
      replicas: createForm.replicas,
      maxRetries: createForm.maxRetries,
      script: createForm.script,
      dataPath: createForm.dataPath,
      outputPath: createForm.outputPath,
      description: createForm.description,
    };

    jobs.value.unshift(newJob);
    createModalVisible.value = false;
    message.success('训练任务创建成功');

    createFormRef.value?.resetFields();
    createForm.envVars = [];
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    createLoading.value = false;
  }
};

const handleCreateCancel = (): void => {
  createModalVisible.value = false;
  createFormRef.value?.resetFields();
  createForm.envVars = [];
};

const viewDetails = (record: TrainingJob): void => {
  selectedJob.value = record;
  detailModalVisible.value = true;
};

const viewLogs = (record: TrainingJob): void => {
  selectedJob.value = record;
  if (mockPods.length > 0) {
    selectedPod.value = mockPods[0]!.name;
  }
  logModalVisible.value = true;
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
    const index = jobs.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      jobs.value[index]!.status = 'paused';
    }
    message.success('任务暂停成功');
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
    const index = jobs.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      jobs.value[index]!.status = 'running';
    }
    message.success('任务恢复成功');
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
    const index = jobs.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      jobs.value[index]!.status = 'cancelled';
    }
    message.success('任务停止成功');
  } catch (error) {
    message.error('停止失败');
  } finally {
    loading.value = false;
  }
};

const handleClone = (record: TrainingJob): void => {
  cloneForm.name = `${record.name}-copy`;
  cloneForm.queue = record.queue;
  cloneForm.priority = record.priority;
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
        queue: cloneForm.queue,
        priority: cloneForm.priority as 'low' | 'medium' | 'high' | 'urgent',
        status: 'pending',
        createTime: new Date().toLocaleString(),
        duration: 0,
        progress: 0,
      };

      jobs.value.unshift(clonedJob);
      cloneModalVisible.value = false;
      message.success('任务克隆成功');
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
    const index = jobs.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      jobs.value.splice(index, 1);
      message.success('任务删除成功');
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
    ...logs.value,
    `${new Date().toLocaleString()} INFO: 日志已刷新`,
  ];
  logs.value = newLogs.slice(-100);
};

const downloadLogs = (): void => {
  const logContent = logs.value.join('\n');
  const blob = new Blob([logContent], { type: 'text/plain' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${selectedJob.value?.name || 'job'}-logs.txt`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
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
.job-queue-container {
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

/* ===== 统计卡片 ===== */
.stats-section {
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 8px !important;
  transition: all 0.3s ease;
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

.running-icon {
  background: linear-gradient(135deg, #52c41a, #73d13d);
}

.pending-icon {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
}

.completed-icon {
  background: linear-gradient(135deg, #52c41a, #73d13d);
}

.failed-icon {
  background: linear-gradient(135deg, #ff4d4f, #ff7875);
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 4px;
}

.stat-label {
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

/* ===== 表格行样式 ===== */
.sci-fi-table :deep(.row-running) {
  background: rgba(82, 196, 26, 0.05) !important;
}

.sci-fi-table :deep(.row-failed) {
  background: rgba(255, 77, 79, 0.05) !important;
}

.sci-fi-table :deep(.row-completed) {
  background: rgba(82, 196, 26, 0.05) !important;
}

/* ===== 任务名称 ===== */
.job-name-wrapper {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.job-name {
  font-weight: 500;
  font-size: 14px;
}

.job-id {
  font-size: 12px;
  opacity: 0.6;
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
  background: #52c41a;
}

.indicator-failed {
  background: #ff4d4f;
}

.indicator-cancelled {
  background: #8c8c8c;
}

.indicator-paused {
  background: #faad14;
}

/* ===== 优先级标签 ===== */
.priority-tag {
  border-radius: 6px !important;
  font-weight: 500 !important;
  padding: 4px 8px !important;
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

.job-progress {
  flex: 1;
}

.progress-text {
  font-size: 12px;
  font-weight: 500;
  min-width: 35px;
}

/* ===== 时间显示 ===== */
.time-text,
.duration-text {
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

/* ===== 环境变量 ===== */
.env-vars-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.env-var-item {
  display: flex;
  gap: 8px;
  align-items: center;
}

.env-key,
.env-value {
  flex: 1;
  border-radius: 6px !important;
}

.env-remove-btn {
  border: none !important;
  background: transparent !important;
  color: #ff4d4f !important;
  padding: 4px !important;
  height: auto !important;
}

.add-env-btn {
  border-radius: 6px !important;
  border-style: dashed !important;
  transition: all 0.3s ease;
}

.add-env-btn:hover {
  color: #1890ff !important;
  border-color: #1890ff !important;
}

/* ===== 详情页样式 ===== */
.detail-content {
  max-height: 600px;
  overflow-y: auto;
}

.detail-tabs :deep(.ant-tabs-tab) {
  font-weight: 500;
}

.detail-descriptions {
  margin-bottom: 24px;
}

/* ===== Pod 表格 ===== */
.pod-table :deep(.ant-table-thead > tr > th) {
  font-weight: 600 !important;
}

.pod-resources {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 12px;
}

/* ===== 事件时间轴 ===== */
.event-timeline {
  max-height: 400px;
  overflow-y: auto;
}

.event-item {
  padding: 8px 0;
}

.event-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.event-type {
  font-weight: 500;
  font-size: 14px;
}

.event-time {
  font-size: 12px;
  opacity: 0.6;
}

.event-message {
  font-size: 13px;
}

/* ===== 日志容器 ===== */
.log-container {
  display: flex;
  flex-direction: column;
  height: 500px;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.pod-select {
  border-radius: 6px !important;
}

.log-refresh-btn,
.log-download-btn {
  border-radius: 4px !important;
  transition: all 0.3s ease;
}

.log-refresh-btn:hover,
.log-download-btn:hover {
  color: #1890ff !important;
  border-color: #1890ff !important;
}

.log-content {
  border-radius: 6px !important;
  padding: 12px !important;
  flex: 1;
  overflow-y: auto;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace !important;
  font-size: 12px;
  line-height: 1.4;
}

.log-line {
  margin: 0;
  padding: 2px 0;
}

.log-error {
  color: #ff4d4f;
}

.log-warning {
  color: #faad14;
}

.log-info {
  color: #1890ff;
}

/* ===== 响应式设计 ===== */
@media (max-width: 768px) {
  .job-queue-container {
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

  .stat-card {
    margin-bottom: 16px;
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
    font-size: 24px;
  }

  .env-var-item {
    flex-direction: column;
    align-items: stretch;
  }

  .log-container {
    height: 400px;
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

  .stat-content {
    flex-direction: column;
    text-align: center;
    gap: 8px;
  }

  .stat-icon {
    align-self: center;
  }

  .log-container {
    height: 300px;
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
