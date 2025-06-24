<template>
  <div class="dataset-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <BarChartOutlined class="title-icon" />
            <span class="title-text">数据集分析</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text">数据集质量检测与统计分析</span>
          </p>
        </div>
        <div class="action-section">
          <a-button
            type="primary"
            size="large"
            @click="showAnalysisModal"
            class="create-btn"
          >
            <PlayCircleOutlined />
            开始分析
          </a-button>
        </div>
      </div>
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
              <a-select-option value="completed">已完成</a-select-option>
              <a-select-option value="running">分析中</a-select-option>
              <a-select-option value="pending">等待中</a-select-option>
              <a-select-option value="failed">失败</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterType"
              placeholder="选择数据类型"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部类型</a-select-option>
              <a-select-option value="image">图像</a-select-option>
              <a-select-option value="text">文本</a-select-option>
              <a-select-option value="tabular">表格</a-select-option>
              <a-select-option value="audio">音频</a-select-option>
              <a-select-option value="video">视频</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="16" :md="8" :lg="8">
            <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索数据集名称"
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
          :data-source="filteredDatasets"
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

          <!-- 数据类型列 -->
          <template #dataType="{ record }">
            <a-tag :color="getTypeColor(record.dataType)" class="type-tag">
              <component :is="getTypeIcon(record.dataType)" class="type-icon" />
              {{ getTypeText(record.dataType) }}
            </a-tag>
          </template>

          <!-- 数据集规模列 -->
          <template #scale="{ record }">
            <div class="scale-info">
              <div class="scale-item">
                <FileOutlined class="scale-icon" />
                <span class="scale-value">{{
                  formatFileCount(record.fileCount)
                }}</span>
              </div>
              <div class="scale-item">
                <DatabaseOutlined class="scale-icon" />
                <span class="scale-value">{{
                  formatFileSize(record.totalSize)
                }}</span>
              </div>
            </div>
          </template>

          <!-- 质量评分列 -->
          <template #qualityScore="{ record }">
            <div class="quality-score">
              <a-progress
                :percent="record.qualityScore"
                :status="getQualityStatus(record.qualityScore)"
                :stroke-color="getQualityColor(record.qualityScore)"
                size="small"
                :show-info="false"
                class="quality-progress"
              />
              <span class="score-text">{{ record.qualityScore }}%</span>
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
                @click="viewReport(record)"
                :disabled="record.status !== 'completed'"
                class="action-btn"
              >
                <EyeOutlined />
                报告
              </a-button>
              <a-button
                type="link"
                size="small"
                @click="downloadReport(record)"
                :disabled="record.status !== 'completed'"
                class="action-btn"
              >
                <DownloadOutlined />
                下载
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
                      key="rerun"
                      :disabled="record.status === 'running'"
                    >
                      <PlayCircleOutlined />
                      重新分析
                    </a-menu-item>
                    <a-menu-item key="export">
                      <ExportOutlined />
                      导出数据
                    </a-menu-item>
                    <a-menu-item key="share">
                      <ShareAltOutlined />
                      分享
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

    <!-- 分析配置模态框 -->
    <a-modal
      v-model:open="analysisModalVisible"
      title="数据集分析配置"
      width="800px"
      :confirm-loading="analysisLoading"
      @ok="handleAnalysisSubmit"
      @cancel="handleAnalysisCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="analysisFormRef"
        :model="analysisForm"
        :rules="analysisFormRules"
        layout="vertical"
        class="analysis-form"
      >
        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="数据集名称" name="name">
              <a-input
                v-model:value="analysisForm.name"
                placeholder="请输入数据集名称"
                class="form-input"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="数据类型" name="dataType">
              <a-select
                v-model:value="analysisForm.dataType"
                placeholder="选择数据类型"
                class="form-select"
              >
                <a-select-option value="image">图像</a-select-option>
                <a-select-option value="text">文本</a-select-option>
                <a-select-option value="tabular">表格</a-select-option>
                <a-select-option value="audio">音频</a-select-option>
                <a-select-option value="video">视频</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="数据源路径" name="sourcePath">
          <a-input
            v-model:value="analysisForm.sourcePath"
            placeholder="请输入数据源路径，如: /data/datasets/imagenet"
            class="form-input"
          />
        </a-form-item>

        <a-divider class="form-divider">分析配置</a-divider>

        <a-form-item label="分析项目" name="analysisTypes">
          <a-checkbox-group
            v-model:value="analysisForm.analysisTypes"
            class="analysis-options"
          >
            <a-row :gutter="[16, 16]">
              <a-col :xs="24" :sm="12" :md="8">
                <a-checkbox value="basic">基础统计</a-checkbox>
              </a-col>
              <a-col :xs="24" :sm="12" :md="8">
                <a-checkbox value="quality">数据质量</a-checkbox>
              </a-col>
              <a-col :xs="24" :sm="12" :md="8">
                <a-checkbox value="distribution">分布分析</a-checkbox>
              </a-col>
              <a-col :xs="24" :sm="12" :md="8">
                <a-checkbox value="duplicates">重复检测</a-checkbox>
              </a-col>
              <a-col :xs="24" :sm="12" :md="8">
                <a-checkbox value="anomaly">异常检测</a-checkbox>
              </a-col>
              <a-col :xs="24" :sm="12" :md="8">
                <a-checkbox value="format">格式验证</a-checkbox>
              </a-col>
            </a-row>
          </a-checkbox-group>
        </a-form-item>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="采样比例" name="sampleRatio">
              <a-slider
                v-model:value="analysisForm.sampleRatio"
                :min="10"
                :max="100"
                :step="10"
                :marks="{ 10: '10%', 50: '50%', 100: '100%' }"
                class="form-slider"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="并行度" name="parallelism">
              <a-input-number
                v-model:value="analysisForm.parallelism"
                :min="1"
                :max="16"
                style="width: 100%"
                addon-after="并发"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="备注" name="description">
          <a-textarea
            v-model:value="analysisForm.description"
            placeholder="请输入分析任务备注"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 分析报告模态框 -->
    <a-modal
      v-model:open="reportModalVisible"
      title="数据集分析报告"
      width="1200px"
      :footer="null"
      class="sci-fi-modal report-modal"
    >
      <div v-if="selectedDataset" class="report-content">
        <!-- 报告概览 -->
        <div class="report-overview">
          <a-row :gutter="16">
            <a-col :xs="24" :sm="6">
              <a-statistic
                title="总文件数"
                :value="selectedDataset.fileCount"
                :value-style="{ color: '#1890ff' }"
                class="stat-card"
              >
                <template #suffix>
                  <FileOutlined />
                </template>
              </a-statistic>
            </a-col>
            <a-col :xs="24" :sm="6">
              <a-statistic
                title="总大小"
                :value="selectedDataset.totalSize"
                suffix="GB"
                :value-style="{ color: '#52c41a' }"
                class="stat-card"
              >
                <template #suffix>
                  <DatabaseOutlined />
                </template>
              </a-statistic>
            </a-col>
            <a-col :xs="24" :sm="6">
              <a-statistic
                title="质量评分"
                :value="selectedDataset.qualityScore"
                suffix="%"
                :value-style="{
                  color: getQualityColor(selectedDataset.qualityScore),
                }"
                class="stat-card"
              >
                <template #suffix>
                  <TrophyOutlined />
                </template>
              </a-statistic>
            </a-col>
            <a-col :xs="24" :sm="6">
              <a-statistic
                title="异常文件"
                :value="reportData.anomalyCount"
                :value-style="{ color: '#ff4d4f' }"
                class="stat-card"
              >
                <template #suffix>
                  <ExclamationCircleOutlined />
                </template>
              </a-statistic>
            </a-col>
          </a-row>
        </div>

        <a-divider class="report-divider">详细分析</a-divider>

        <!-- 分析报告标签页 -->
        <a-tabs v-model:activeKey="activeReportTab" class="report-tabs">
          <a-tab-pane key="basic" tab="基础统计">
            <div class="basic-stats">
              <a-row :gutter="16">
                <a-col :xs="24" :md="12">
                  <a-card title="文件类型分布" class="chart-card">
                    <div class="chart-container">
                      <canvas
                        ref="fileTypeChartRef"
                        class="chart-canvas"
                      ></canvas>
                    </div>
                  </a-card>
                </a-col>
                <a-col :xs="24" :md="12">
                  <a-card title="文件大小分布" class="chart-card">
                    <div class="chart-container">
                      <canvas
                        ref="fileSizeChartRef"
                        class="chart-canvas"
                      ></canvas>
                    </div>
                  </a-card>
                </a-col>
              </a-row>
            </div>
          </a-tab-pane>

          <a-tab-pane key="quality" tab="质量分析">
            <div class="quality-analysis">
              <a-row :gutter="16">
                <a-col :xs="24" :md="8">
                  <a-card title="质量维度评分" class="quality-card">
                    <div class="quality-metrics">
                      <div class="metric-item">
                        <span class="metric-label">完整性</span>
                        <a-progress
                          :percent="reportData.qualityMetrics.completeness"
                          size="small"
                          :stroke-color="
                            getQualityColor(
                              reportData.qualityMetrics.completeness,
                            )
                          "
                        />
                      </div>
                      <div class="metric-item">
                        <span class="metric-label">一致性</span>
                        <a-progress
                          :percent="reportData.qualityMetrics.consistency"
                          size="small"
                          :stroke-color="
                            getQualityColor(
                              reportData.qualityMetrics.consistency,
                            )
                          "
                        />
                      </div>
                      <div class="metric-item">
                        <span class="metric-label">准确性</span>
                        <a-progress
                          :percent="reportData.qualityMetrics.accuracy"
                          size="small"
                          :stroke-color="
                            getQualityColor(reportData.qualityMetrics.accuracy)
                          "
                        />
                      </div>
                      <div class="metric-item">
                        <span class="metric-label">有效性</span>
                        <a-progress
                          :percent="reportData.qualityMetrics.validity"
                          size="small"
                          :stroke-color="
                            getQualityColor(reportData.qualityMetrics.validity)
                          "
                        />
                      </div>
                    </div>
                  </a-card>
                </a-col>
                <a-col :xs="24" :md="16">
                  <a-card title="质量问题分布" class="chart-card">
                    <div class="chart-container">
                      <canvas
                        ref="qualityIssuesChartRef"
                        class="chart-canvas"
                      ></canvas>
                    </div>
                  </a-card>
                </a-col>
              </a-row>
            </div>
          </a-tab-pane>

          <a-tab-pane key="duplicates" tab="重复检测">
            <div class="duplicates-analysis">
              <a-alert
                :message="`发现 ${reportData.duplicateGroups.length} 组重复文件，共 ${reportData.totalDuplicates} 个文件`"
                type="warning"
                show-icon
                class="duplicate-alert"
              />
              <a-table
                :columns="duplicateColumns"
                :data-source="reportData.duplicateGroups"
                :pagination="{ pageSize: 5 }"
                size="small"
                class="duplicate-table"
              >
                <template #files="{ record }">
                  <a-tag
                    v-for="file in record.files.slice(0, 3)"
                    :key="file"
                    class="file-tag"
                  >
                    {{ file.split('/').pop() }}
                  </a-tag>
                  <a-tag v-if="record.files.length > 3" class="more-tag">
                    +{{ record.files.length - 3 }}
                  </a-tag>
                </template>
              </a-table>
            </div>
          </a-tab-pane>

          <a-tab-pane key="anomaly" tab="异常检测">
            <div class="anomaly-analysis">
              <a-table
                :columns="anomalyColumns"
                :data-source="reportData.anomalies"
                :pagination="{ pageSize: 10 }"
                size="small"
                class="anomaly-table"
              >
                <template #type="{ record }">
                  <a-tag :color="getAnomalyTypeColor(record.type)">
                    {{ getAnomalyTypeText(record.type) }}
                  </a-tag>
                </template>
                <template #severity="{ record }">
                  <a-tag :color="getSeverityColor(record.severity)">
                    {{ getSeverityText(record.severity) }}
                  </a-tag>
                </template>
              </a-table>
            </div>
          </a-tab-pane>
        </a-tabs>

        <!-- 报告操作栏 -->
        <div class="report-actions">
          <a-space>
            <a-button
              @click="downloadFullReport(selectedDataset)"
              class="action-btn"
            >
              <DownloadOutlined />
              下载完整报告
            </a-button>
            <a-button
              @click="exportReportData(selectedDataset)"
              class="action-btn"
            >
              <ExportOutlined />
              导出原始数据
            </a-button>
            <a-button @click="shareReport(selectedDataset)" class="action-btn">
              <ShareAltOutlined />
              分享报告
            </a-button>
          </a-space>
        </div>
      </div>
    </a-modal>

    <!-- 导出进度模态框 -->
    <a-modal
      v-model:open="exportModalVisible"
      title="导出进度"
      width="500px"
      :footer="null"
      class="sci-fi-modal"
    >
      <div class="export-progress">
        <a-progress
          :percent="exportProgress"
          :status="exportProgress === 100 ? 'success' : 'active'"
          :stroke-color="exportProgress === 100 ? '#52c41a' : '#1890ff'"
        />
        <p class="progress-text">
          {{ exportProgress === 100 ? '导出完成' : '正在导出数据...' }}
        </p>
        <div v-if="exportProgress === 100" class="export-complete">
          <a-button
            type="primary"
            @click="downloadExportedFile"
            class="download-btn"
          >
            <DownloadOutlined />
            下载文件
          </a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue';
import { message } from 'ant-design-vue';
import type { FormInstance, TableColumnsType } from 'ant-design-vue';
import {
  BarChartOutlined,
  PlayCircleOutlined,
  ReloadOutlined,
  FileOutlined,
  DatabaseOutlined,
  EyeOutlined,
  DownloadOutlined,
  MoreOutlined,
  ExportOutlined,
  ShareAltOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  LoadingOutlined,
  CloseCircleOutlined,
  PictureOutlined,
  FileTextOutlined,
  TableOutlined,
  AudioOutlined,
  VideoCameraOutlined,
  TrophyOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons-vue';

// ===== 类型定义 =====
interface DatasetItem {
  id: string;
  name: string;
  dataType: 'image' | 'text' | 'tabular' | 'audio' | 'video';
  status: 'completed' | 'running' | 'pending' | 'failed';
  fileCount: number;
  totalSize: number;
  qualityScore: number;
  createTime: string;
  finishTime?: string;
  creator: string;
  sourcePath: string;
  description?: string;
}

interface AnalysisForm {
  name: string;
  dataType: string;
  sourcePath: string;
  analysisTypes: string[];
  sampleRatio: number;
  parallelism: number;
  description: string;
}

interface QualityMetrics {
  completeness: number;
  consistency: number;
  accuracy: number;
  validity: number;
}

interface DuplicateGroup {
  id: string;
  hash: string;
  count: number;
  totalSize: number;
  files: string[];
}

interface AnomalyItem {
  id: string;
  file: string;
  type: 'size' | 'format' | 'encoding' | 'metadata';
  severity: 'low' | 'medium' | 'high';
  description: string;
}

interface ReportData {
  anomalyCount: number;
  qualityMetrics: QualityMetrics;
  duplicateGroups: DuplicateGroup[];
  totalDuplicates: number;
  anomalies: AnomalyItem[];
}

// ===== 响应式数据 =====
const loading = ref<boolean>(false);
const analysisModalVisible = ref<boolean>(false);
const reportModalVisible = ref<boolean>(false);
const exportModalVisible = ref<boolean>(false);
const analysisLoading = ref<boolean>(false);
const exportProgress = ref<number>(0);

const filterStatus = ref<string>('');
const filterType = ref<string>('');
const searchKeyword = ref<string>('');

const selectedDataset = ref<DatasetItem | null>(null);
const activeReportTab = ref<string>('basic');

// 图表引用
const fileTypeChartRef = ref<HTMLCanvasElement | null>(null);
const fileSizeChartRef = ref<HTMLCanvasElement | null>(null);
const qualityIssuesChartRef = ref<HTMLCanvasElement | null>(null);

// ===== 表单引用 =====
const analysisFormRef = ref<FormInstance>();

// ===== 表单数据 =====
const analysisForm = reactive<AnalysisForm>({
  name: '',
  dataType: 'image',
  sourcePath: '',
  analysisTypes: ['basic', 'quality'],
  sampleRatio: 100,
  parallelism: 4,
  description: '',
});

// ===== 配置数据 =====
const STATUS_CONFIG = {
  completed: { color: 'success', text: '已完成', icon: CheckCircleOutlined },
  running: { color: 'processing', text: '分析中', icon: LoadingOutlined },
  pending: { color: 'default', text: '等待中', icon: ClockCircleOutlined },
  failed: { color: 'error', text: '失败', icon: CloseCircleOutlined },
} as const;

const TYPE_CONFIG = {
  image: { color: 'blue', text: '图像', icon: PictureOutlined },
  text: { color: 'green', text: '文本', icon: FileTextOutlined },
  tabular: { color: 'orange', text: '表格', icon: TableOutlined },
  audio: { color: 'purple', text: '音频', icon: AudioOutlined },
  video: { color: 'red', text: '视频', icon: VideoCameraOutlined },
} as const;

// ===== 模拟数据 =====
const datasets = ref<DatasetItem[]>([
  {
    id: 'ds-001',
    name: 'imagenet-subset',
    dataType: 'image',
    status: 'completed',
    fileCount: 50000,
    totalSize: 15.6,
    qualityScore: 85,
    createTime: '2024-06-23 09:30:00',
    finishTime: '2024-06-23 11:45:00',
    creator: 'admin',
    sourcePath: '/data/datasets/imagenet',
    description: 'ImageNet数据集子集，用于图像分类训练',
  },
  {
    id: 'ds-002',
    name: 'sentiment-analysis-data',
    dataType: 'text',
    status: 'running',
    fileCount: 100000,
    totalSize: 2.3,
    qualityScore: 0,
    createTime: '2024-06-23 10:15:00',
    creator: 'researcher',
    sourcePath: '/data/datasets/sentiment',
    description: '情感分析文本数据集',
  },
  {
    id: 'ds-003',
    name: 'sales-data-2024',
    dataType: 'tabular',
    status: 'completed',
    fileCount: 12,
    totalSize: 0.8,
    qualityScore: 92,
    createTime: '2024-06-22 16:30:00',
    finishTime: '2024-06-22 16:45:00',
    creator: 'analyst',
    sourcePath: '/data/datasets/sales',
    description: '2024年销售数据分析',
  },
  {
    id: 'ds-004',
    name: 'speech-recognition-corpus',
    dataType: 'audio',
    status: 'pending',
    fileCount: 25000,
    totalSize: 45.2,
    qualityScore: 0,
    createTime: '2024-06-23 11:00:00',
    creator: 'developer',
    sourcePath: '/data/datasets/speech',
    description: '语音识别语料库',
  },
  {
    id: 'ds-005',
    name: 'video-classification-dataset',
    dataType: 'video',
    status: 'failed',
    fileCount: 1500,
    totalSize: 125.8,
    qualityScore: 0,
    createTime: '2024-06-23 08:30:00',
    creator: 'researcher',
    sourcePath: '/data/datasets/videos',
    description: '视频分类数据集',
  },
]);

const reportData = reactive<ReportData>({
  anomalyCount: 245,
  qualityMetrics: {
    completeness: 85,
    consistency: 78,
    accuracy: 92,
    validity: 88,
  },
  duplicateGroups: [
    {
      id: 'dup-001',
      hash: 'md5:a1b2c3d4e5f6',
      count: 15,
      totalSize: 2.4,
      files: [
        '/data/images/cat_001.jpg',
        '/data/images/duplicate/cat_001_copy.jpg',
        '/data/images/backup/cat_001.jpg',
        '/data/images/cat_001_v2.jpg',
        '/data/images/cat_001_backup.jpg',
      ],
    },
    {
      id: 'dup-002',
      hash: 'md5:b2c3d4e5f6g7',
      count: 8,
      totalSize: 1.2,
      files: [
        '/data/images/dog_025.jpg',
        '/data/images/duplicate/dog_025_copy.jpg',
        '/data/images/dog_025_v2.jpg',
      ],
    },
  ],
  totalDuplicates: 23,
  anomalies: [
    {
      id: 'anom-001',
      file: '/data/images/corrupted_001.jpg',
      type: 'format',
      severity: 'high',
      description: '文件格式损坏，无法正常读取',
    },
    {
      id: 'anom-002',
      file: '/data/images/oversized_002.jpg',
      type: 'size',
      severity: 'medium',
      description: '文件大小异常，超过预期范围',
    },
    {
      id: 'anom-003',
      file: '/data/images/invalid_encoding.jpg',
      type: 'encoding',
      severity: 'low',
      description: '字符编码异常',
    },
  ],
});

// ===== 表单验证规则 =====
const analysisFormRules = {
  name: [
    { required: true, message: '请输入数据集名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
  ],
  dataType: [{ required: true, message: '请选择数据类型', trigger: 'change' }],
  sourcePath: [
    { required: true, message: '请输入数据源路径', trigger: 'blur' },
  ],
  analysisTypes: [
    { required: true, message: '请至少选择一项分析内容', trigger: 'change' },
  ],
};

// ===== 表格列配置 =====
const columns: TableColumnsType<DatasetItem> = [
  {
    title: '数据集名称',
    dataIndex: 'name',
    key: 'name',
    width: 200,
    ellipsis: true,
  },
  {
    title: '数据类型',
    key: 'dataType',
    width: 100,
    slots: { customRender: 'dataType' },
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
    slots: { customRender: 'status' },
  },
  {
    title: '数据规模',
    key: 'scale',
    width: 150,
    slots: { customRender: 'scale' },
  },
  {
    title: '质量评分',
    key: 'qualityScore',
    width: 120,
    slots: { customRender: 'qualityScore' },
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
    width: 100,
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

// 重复文件表格列
const duplicateColumns: TableColumnsType<DuplicateGroup> = [
  {
    title: 'Hash值',
    dataIndex: 'hash',
    key: 'hash',
    width: 150,
    ellipsis: true,
  },
  {
    title: '重复数量',
    dataIndex: 'count',
    key: 'count',
    width: 100,
  },
  {
    title: '总大小(MB)',
    dataIndex: 'totalSize',
    key: 'totalSize',
    width: 120,
  },
  {
    title: '文件列表',
    key: 'files',
    slots: { customRender: 'files' },
  },
];

// 异常文件表格列
const anomalyColumns: TableColumnsType<AnomalyItem> = [
  {
    title: '文件路径',
    dataIndex: 'file',
    key: 'file',
    width: 250,
    ellipsis: true,
  },
  {
    title: '异常类型',
    key: 'type',
    width: 100,
    slots: { customRender: 'type' },
  },
  {
    title: '严重程度',
    key: 'severity',
    width: 100,
    slots: { customRender: 'severity' },
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
    ellipsis: true,
  },
];

// ===== 分页配置 =====
const paginationConfig = {
  total: computed(() => filteredDatasets.value.length),
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// ===== 计算属性 =====
const filteredDatasets = computed(() => {
  let result = datasets.value;

  if (filterStatus.value) {
    result = result.filter((item) => item.status === filterStatus.value);
  }

  if (filterType.value) {
    result = result.filter((item) => item.dataType === filterType.value);
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter((item) => item.name.toLowerCase().includes(keyword));
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

const getTypeColor = (type: string): string => {
  return TYPE_CONFIG[type as keyof typeof TYPE_CONFIG]?.color || 'default';
};

const getTypeIcon = (type: string) => {
  return (
    TYPE_CONFIG[type as keyof typeof TYPE_CONFIG]?.icon || FileTextOutlined
  );
};

const getTypeText = (type: string): string => {
  return TYPE_CONFIG[type as keyof typeof TYPE_CONFIG]?.text || type;
};

const getQualityStatus = (score: number): string => {
  if (score >= 80) return 'success';
  if (score >= 60) return 'normal';
  return 'exception';
};

const getQualityColor = (score: number): string => {
  if (score >= 80) return '#52c41a';
  if (score >= 60) return '#faad14';
  return '#ff4d4f';
};

const getAnomalyTypeColor = (type: string): string => {
  const colors = {
    size: 'orange',
    format: 'red',
    encoding: 'blue',
    metadata: 'purple',
  };
  return colors[type as keyof typeof colors] || 'default';
};

const getAnomalyTypeText = (type: string): string => {
  const texts = {
    size: '大小异常',
    format: '格式错误',
    encoding: '编码异常',
    metadata: '元数据错误',
  };
  return texts[type as keyof typeof texts] || type;
};

const getSeverityColor = (severity: string): string => {
  const colors = {
    low: 'green',
    medium: 'orange',
    high: 'red',
  };
  return colors[severity as keyof typeof colors] || 'default';
};

const getSeverityText = (severity: string): string => {
  const texts = {
    low: '低',
    medium: '中',
    high: '高',
  };
  return texts[severity as keyof typeof texts] || severity;
};

const formatFileCount = (count: number): string => {
  if (count >= 1000000) {
    return `${(count / 1000000).toFixed(1)}M`;
  } else if (count >= 1000) {
    return `${(count / 1000).toFixed(1)}K`;
  }
  return count.toString();
};

const formatFileSize = (size: number): string => {
  if (size >= 1024) {
    return `${(size / 1024).toFixed(1)}TB`;
  } else if (size >= 1) {
    return `${size.toFixed(1)}GB`;
  } else {
    return `${(size * 1024).toFixed(0)}MB`;
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

// ===== 图表绘制函数 =====
const drawCharts = (): void => {
  nextTick(() => {
    drawFileTypeChart();
    drawFileSizeChart();
    drawQualityIssuesChart();
  });
};

const drawFileTypeChart = (): void => {
  const canvas = fileTypeChartRef.value;
  if (!canvas) return;

  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  // 简单的饼图绘制
  const data = [
    { label: 'JPG', value: 45, color: '#1890ff' },
    { label: 'PNG', value: 30, color: '#52c41a' },
    { label: 'GIF', value: 15, color: '#faad14' },
    { label: '其他', value: 10, color: '#f5222d' },
  ];

  canvas.width = 300;
  canvas.height = 200;

  const centerX = canvas.width / 2;
  const centerY = canvas.height / 2;
  const radius = 60;

  let currentAngle = 0;
  data.forEach((item) => {
    const sliceAngle = (item.value / 100) * 2 * Math.PI;

    ctx.beginPath();
    ctx.arc(centerX, centerY, radius, currentAngle, currentAngle + sliceAngle);
    ctx.lineTo(centerX, centerY);
    ctx.fillStyle = item.color;
    ctx.fill();

    currentAngle += sliceAngle;
  });
};

const drawFileSizeChart = (): void => {
  const canvas = fileSizeChartRef.value;
  if (!canvas) return;

  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  canvas.width = 300;
  canvas.height = 200;

  // 简单的柱状图
  const data = [10, 25, 35, 20, 8, 2];
  const labels = ['<1MB', '1-5MB', '5-10MB', '10-50MB', '50-100MB', '>100MB'];

  const barWidth = 40;
  const maxHeight = 150;
  const maxValue = Math.max(...data);

  data.forEach((value, index) => {
    const height = (value / maxValue) * maxHeight;
    const x = index * 50 + 10;
    const y = canvas.height - height - 20;

    ctx.fillStyle = '#1890ff';
    ctx.fillRect(x, y, barWidth, height);

    ctx.fillStyle = '#333';
    ctx.font = '10px Arial';
    ctx.fillText(labels[index] || '', x, canvas.height - 5);
  });
};

const drawQualityIssuesChart = (): void => {
  const canvas = qualityIssuesChartRef.value;
  if (!canvas) return;

  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  canvas.width = 400;
  canvas.height = 200;

  // 简单的水平条形图
  const data = [
    { label: '格式错误', value: 45, color: '#f5222d' },
    { label: '大小异常', value: 30, color: '#faad14' },
    { label: '编码问题', value: 15, color: '#1890ff' },
    { label: '元数据缺失', value: 10, color: '#722ed1' },
  ];

  const barHeight = 30;
  const maxWidth = 300;
  const maxValue = Math.max(...data.map((d) => d.value));

  data.forEach((item, index) => {
    const width = (item.value / maxValue) * maxWidth;
    const y = index * 40 + 20;

    ctx.fillStyle = item.color;
    ctx.fillRect(50, y, width, barHeight);

    ctx.fillStyle = '#333';
    ctx.font = '12px Arial';
    ctx.fillText(item.label, 10, y + 20);
    ctx.fillText(`${item.value}`, width + 60, y + 20);
  });
};

// ===== 事件处理函数 =====
const showAnalysisModal = (): void => {
  analysisModalVisible.value = true;
};

const handleAnalysisSubmit = async (): Promise<void> => {
  try {
    await analysisFormRef.value?.validate();
    analysisLoading.value = true;

    // 模拟分析任务创建
    await new Promise((resolve) => setTimeout(resolve, 2000));

    const newDataset: DatasetItem = {
      id: `ds-${Date.now()}`,
      name: analysisForm.name,
      dataType: analysisForm.dataType as DatasetItem['dataType'],
      status: 'pending',
      fileCount: 0,
      totalSize: 0,
      qualityScore: 0,
      createTime: new Date().toLocaleString(),
      creator: 'current-user',
      sourcePath: analysisForm.sourcePath,
      description: analysisForm.description,
    };

    datasets.value.unshift(newDataset);
    analysisModalVisible.value = false;
    message.success('分析任务创建成功');

    // 模拟状态变化
    setTimeout(() => {
      const index = datasets.value.findIndex(
        (item) => item.id === newDataset.id,
      );
      if (index !== -1) {
        datasets.value[index]!.status = 'running';
      }
    }, 3000);

    // 重置表单
    analysisFormRef.value?.resetFields();
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    analysisLoading.value = false;
  }
};

const handleAnalysisCancel = (): void => {
  analysisModalVisible.value = false;
  analysisFormRef.value?.resetFields();
};

const viewReport = (record: DatasetItem): void => {
  selectedDataset.value = record;
  reportModalVisible.value = true;
  activeReportTab.value = 'basic';
  drawCharts();
};

const downloadReport = async (record: DatasetItem): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));

    // 模拟下载
    const link = document.createElement('a');
    link.href =
      'data:text/plain;charset=utf-8,' +
      encodeURIComponent(
        `数据集分析报告 - ${record.name}\n\n` +
          `文件数量: ${record.fileCount}\n` +
          `总大小: ${record.totalSize}GB\n` +
          `质量评分: ${record.qualityScore}%\n`,
      );
    link.download = `${record.name}-report.txt`;
    link.click();

    message.success('报告下载成功');
  } catch (error) {
    message.error('下载失败');
  } finally {
    loading.value = false;
  }
};

const handleMenuAction = (key: string, record: DatasetItem): void => {
  const actions = {
    rerun: () => handleRerun(record),
    export: () => handleExport(record),
    share: () => handleShare(record),
    delete: () => handleDelete(record),
  };

  const action = actions[key as keyof typeof actions];
  if (action) {
    action();
  }
};

const handleRerun = async (record: DatasetItem): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const index = datasets.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      datasets.value[index]!.status = 'running';
      datasets.value[index]!.qualityScore = 0;
    }
    message.success('重新分析任务已启动');
  } catch (error) {
    message.error('操作失败');
  } finally {
    loading.value = false;
  }
};

const handleExport = (_: DatasetItem): void => {
  exportModalVisible.value = true;
  exportProgress.value = 0;

  // 模拟导出进度
  const timer = setInterval(() => {
    exportProgress.value += 10;
    if (exportProgress.value >= 100) {
      clearInterval(timer);
    }
  }, 500);
};

const handleShare = (record: DatasetItem): void => {
  // 模拟生成分享链接
  const shareUrl = `${window.location.origin}/dataset/report/${record.id}`;
  navigator.clipboard
    .writeText(shareUrl)
    .then(() => {
      message.success('分享链接已复制到剪贴板');
    })
    .catch(() => {
      message.error('复制失败，请手动复制');
    });
};

const handleDelete = (record: DatasetItem): void => {
  const deleteConfirm = () => {
    const index = datasets.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      datasets.value.splice(index, 1);
      message.success('数据集分析记录删除成功');
    }
  };

  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除数据集分析 "${record.name}" 吗？此操作不可恢复。`,
      okText: '确认',
      cancelText: '取消',
      type: 'warning',
      onOk: deleteConfirm,
    });
  });
};

const downloadFullReport = (_: DatasetItem): void => {
  message.success('完整报告下载中...');
};

const exportReportData = (record: DatasetItem): void => {
  handleExport(record);
};

const shareReport = (record: DatasetItem): void => {
  handleShare(record);
};

const downloadExportedFile = (): void => {
  const link = document.createElement('a');
  link.href =
    'data:application/json;charset=utf-8,' +
    encodeURIComponent(JSON.stringify(reportData, null, 2));
  link.download = 'dataset-analysis-data.json';
  link.click();

  exportModalVisible.value = false;
  exportProgress.value = 0;
  message.success('文件下载成功');
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
/* 继承基础样式并扩展 */
.dataset-container {
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

/* ===== 状态和类型标签 ===== */
.status-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-tag,
.type-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border-radius: 6px !important;
  font-weight: 500 !important;
  padding: 4px 8px !important;
}

.status-icon,
.type-icon {
  font-size: 12px;
}

.status-indicator {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.indicator-completed {
  background: #52c41a;
}

.indicator-running {
  background: #1890ff;
}

.indicator-pending {
  background: #8c8c8c;
}

.indicator-failed {
  background: #ff4d4f;
}

/* ===== 数据规模信息 ===== */
.scale-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.scale-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.scale-item:hover {
  color: #1890ff;
}

.scale-icon {
  font-size: 12px;
  color: #1890ff;
}

.scale-value {
  font-weight: 600;
}

/* ===== 质量评分 ===== */
.quality-score {
  display: flex;
  align-items: center;
  gap: 8px;
}

.quality-progress {
  flex: 1;
  max-width: 80px;
}

.score-text {
  font-size: 12px;
  font-weight: 600;
  min-width: 40px;
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
.analysis-form :deep(.ant-form-item-label > label) {
  font-weight: 500 !important;
}

.form-input,
.form-select,
.form-textarea,
.form-input-number,
.form-slider {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

.form-divider {
  font-weight: 500 !important;
}

.analysis-options {
  width: 100%;
}

/* ===== 报告模态框样式 ===== */
.report-modal :deep(.ant-modal-content) {
  max-height: 90vh;
  overflow: hidden;
}

.report-content {
  max-height: 70vh;
  overflow-y: auto;
}

/* ===== 报告概览 ===== */
.report-overview {
  margin-bottom: 24px;
}

.stat-card :deep(.ant-statistic-title) {
  font-size: 12px !important;
  font-weight: 500 !important;
}

.stat-card :deep(.ant-statistic-content) {
  font-size: 20px !important;
  font-weight: 600 !important;
}

/* ===== 报告标签页 ===== */
.report-divider {
  font-weight: 500 !important;
  margin: 24px 0 16px 0 !important;
}

.report-tabs :deep(.ant-tabs-tab) {
  font-weight: 500 !important;
}

/* ===== 图表容器 ===== */
.chart-card {
  border-radius: 8px !important;
  height: 280px;
}

.chart-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}

.chart-canvas {
  max-width: 100%;
  max-height: 100%;
}

/* ===== 质量分析 ===== */
.quality-card {
  border-radius: 8px !important;
  height: 280px;
}

.quality-metrics {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px 0;
}

.metric-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.metric-label {
  font-size: 12px;
  font-weight: 500;
  color: #666;
}

/* ===== 重复检测 ===== */
.duplicate-alert {
  margin-bottom: 16px;
  border-radius: 6px !important;
}

.duplicate-table {
  margin-top: 16px;
}

.file-tag {
  margin: 2px;
  font-size: 11px;
  border-radius: 4px !important;
}

.more-tag {
  background: #f0f0f0;
  color: #666;
  border: 1px solid #d9d9d9;
}

/* ===== 异常检测 ===== */
.anomaly-table {
  margin-top: 16px;
}

/* ===== 报告操作栏 ===== */
.report-actions {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: center;
}

/* ===== 导出进度 ===== */
.export-progress {
  text-align: center;
  padding: 24px;
}

.progress-text {
  margin: 16px 0;
  font-size: 14px;
  color: #666;
}

.export-complete {
  margin-top: 24px;
}

.download-btn {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

/* ===== 响应式设计 ===== */
@media (max-width: 768px) {
  .dataset-container {
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

  .report-modal :deep(.ant-modal) {
    margin: 8px !important;
    max-width: calc(100vw - 16px) !important;
  }

  .chart-card {
    height: auto;
  }

  .chart-container {
    height: 150px;
  }

  .quality-card {
    height: auto;
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

  .scale-info {
    gap: 2px;
  }

  .scale-item {
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

  .quality-score {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .quality-progress {
    max-width: 100%;
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
