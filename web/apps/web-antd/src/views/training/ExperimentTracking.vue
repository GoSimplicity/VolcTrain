<template>
  <div class="experiment-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <ExperimentOutlined class="title-icon" />
            <span class="title-text">实验跟踪</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text"> 管理和监控您的机器学习实验 </span>
          </p>
        </div>
        <div class="action-section">
          <a-button
            v-if="selectedExperimentIds.length >= 2"
            type="primary"
            size="large"
            @click="handleCompareSelected"
            class="compare-btn"
            style="margin-right: 12px"
          >
            <SwapOutlined />
            对比选中 ({{ selectedExperimentIds.length }})
          </a-button>
          <a-button
            type="primary"
            size="large"
            @click="showCreateModal"
            class="create-btn"
          >
            <PlusOutlined />
            创建实验
          </a-button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-section">
      <a-row :gutter="16">
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <a-statistic
              title="总实验数"
              :value="experiments.length"
              :value-style="{ color: '#1890ff', fontSize: '24px' }"
            >
              <template #prefix>
                <ExperimentOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <a-statistic
              title="运行中"
              :value="runningCount"
              :value-style="{ color: '#52c41a', fontSize: '24px' }"
            >
              <template #prefix>
                <PlayCircleOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <a-statistic
              title="已完成"
              :value="completedCount"
              :value-style="{ color: '#52c41a', fontSize: '24px' }"
            >
              <template #prefix>
                <CheckCircleOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <a-statistic
              title="失败"
              :value="failedCount"
              :value-style="{ color: '#ff4d4f', fontSize: '24px' }"
            >
              <template #prefix>
                <CloseCircleOutlined />
              </template>
            </a-statistic>
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
              <a-select-option value="completed">已完成</a-select-option>
              <a-select-option value="failed">失败</a-select-option>
              <a-select-option value="stopped">已停止</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterProject"
              placeholder="选择项目"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部项目</a-select-option>
              <a-select-option value="image-classification">
                图像分类
              </a-select-option>
              <a-select-option value="nlp-sentiment">情感分析</a-select-option>
              <a-select-option value="recommendation">推荐系统</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="16" :md="8" :lg="8">
            <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索实验名称或创建者"
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
          :data-source="filteredExperiments"
          :loading="loading"
          :pagination="paginationConfig"
          :row-selection="rowSelection"
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

          <!-- 指标列 -->
          <template #metrics="{ record }">
            <div class="metrics-info">
              <div
                v-for="(value, key) in record.metrics"
                :key="key"
                class="metric-item"
              >
                <span class="metric-label">{{ key }}:</span>
                <span class="metric-value">{{ formatMetric(value) }}</span>
              </div>
            </div>
          </template>

          <!-- 持续时间列 -->
          <template #duration="{ record }">
            <span class="duration-text">
              {{ formatDuration(record.duration) }}
            </span>
          </template>

          <!-- 创建时间列 -->
          <template #createTime="{ record }">
            <a-tooltip :title="record.createTime">
              <span class="time-text">
                {{ formatRelativeTime(record.createTime) }}
              </span>
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
              <a-button
                type="link"
                size="small"
                @click="viewMetrics(record)"
                class="action-btn"
              >
                <LineChartOutlined />
                指标
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
                      key="stop"
                      :disabled="record.status !== 'running'"
                    >
                      <StopOutlined />
                      停止
                    </a-menu-item>
                    <a-menu-item
                      key="restart"
                      :disabled="record.status === 'running'"
                    >
                      <ReloadOutlined />
                      重新运行
                    </a-menu-item>
                    <a-menu-item key="clone">
                      <CopyOutlined />
                      克隆
                    </a-menu-item>
                    <a-menu-item key="compare">
                      <SwapOutlined />
                      对比
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

    <!-- 创建实验模态框 -->
    <a-modal
      v-model:open="createModalVisible"
      title="创建实验"
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
            <a-form-item label="实验名称" name="name">
              <a-input
                v-model:value="createForm.name"
                placeholder="请输入实验名称"
                class="form-input"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="项目" name="project">
              <a-select
                v-model:value="createForm.project"
                placeholder="选择项目"
                class="form-select"
              >
                <a-select-option value="image-classification">
                  图像分类
                </a-select-option>
                <a-select-option value="nlp-sentiment">
                  情感分析
                </a-select-option>
                <a-select-option value="recommendation">
                  推荐系统
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="算法框架" name="framework">
          <a-select
            v-model:value="createForm.framework"
            placeholder="选择框架"
            class="form-select"
          >
            <a-select-option value="tensorflow">TensorFlow</a-select-option>
            <a-select-option value="pytorch">PyTorch</a-select-option>
            <a-select-option value="sklearn">Scikit-Learn</a-select-option>
            <a-select-option value="xgboost">XGBoost</a-select-option>
          </a-select>
        </a-form-item>

        <a-divider class="form-divider">训练配置</a-divider>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="8">
            <a-form-item label="学习率" name="learningRate">
              <a-input-number
                v-model:value="createForm.learningRate"
                :min="0.0001"
                :max="1"
                :step="0.001"
                style="width: 100%"
                placeholder="0.001"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="8">
            <a-form-item label="批次大小" name="batchSize">
              <a-input-number
                v-model:value="createForm.batchSize"
                :min="1"
                :max="1024"
                style="width: 100%"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="8">
            <a-form-item label="训练轮数" name="epochs">
              <a-input-number
                v-model:value="createForm.epochs"
                :min="1"
                :max="1000"
                style="width: 100%"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider class="form-divider">资源配置</a-divider>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="8">
            <a-form-item label="CPU 核数" name="cpu">
              <a-input-number
                v-model:value="createForm.cpu"
                :min="1"
                :max="32"
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
                :max="128"
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

        <a-form-item label="描述" name="description">
          <a-textarea
            v-model:value="createForm.description"
            placeholder="请输入实验描述"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>

        <a-form-item label="代码库地址" name="codeRepository">
          <a-input
            v-model:value="createForm.codeRepository"
            placeholder="https://github.com/username/repo.git"
            class="form-input"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 实验详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="实验详情"
      width="1000px"
      :footer="null"
      class="sci-fi-modal detail-modal"
    >
      <div v-if="selectedExperiment" class="detail-content">
        <a-tabs v-model:activeKey="activeTabKey" class="detail-tabs">
          <a-tab-pane key="overview" tab="概览">
            <a-descriptions
              :column="{ xs: 1, sm: 2 }"
              bordered
              class="detail-descriptions"
            >
              <a-descriptions-item label="实验名称">
                {{ selectedExperiment.name }}
              </a-descriptions-item>
              <a-descriptions-item label="状态">
                <a-tag
                  :color="getStatusColor(selectedExperiment.status)"
                  class="status-tag"
                >
                  <component :is="getStatusIcon(selectedExperiment.status)" />
                  {{ getStatusText(selectedExperiment.status) }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="项目">
                {{ selectedExperiment.project }}
              </a-descriptions-item>
              <a-descriptions-item label="框架">
                {{ selectedExperiment.framework }}
              </a-descriptions-item>
              <a-descriptions-item label="创建者">
                {{ selectedExperiment.creator }}
              </a-descriptions-item>
              <a-descriptions-item label="持续时间">
                {{ formatDuration(selectedExperiment.duration) }}
              </a-descriptions-item>
              <a-descriptions-item label="学习率">
                {{ selectedExperiment.hyperParams.learningRate }}
              </a-descriptions-item>
              <a-descriptions-item label="批次大小">
                {{ selectedExperiment.hyperParams.batchSize }}
              </a-descriptions-item>
              <a-descriptions-item label="训练轮数">
                {{ selectedExperiment.hyperParams.epochs }}
              </a-descriptions-item>
              <a-descriptions-item label="创建时间">
                {{ selectedExperiment.createTime }}
              </a-descriptions-item>
              <a-descriptions-item label="代码库" :span="2">
                <a
                  v-if="selectedExperiment.codeRepository"
                  :href="selectedExperiment.codeRepository"
                  target="_blank"
                  class="repo-link"
                >
                  {{ selectedExperiment.codeRepository }}
                  <ExportOutlined />
                </a>
                <span v-else>未配置</span>
              </a-descriptions-item>
              <a-descriptions-item label="描述" :span="2">
                {{ selectedExperiment.description || '暂无描述' }}
              </a-descriptions-item>
            </a-descriptions>
          </a-tab-pane>

          <a-tab-pane key="metrics" tab="训练指标">
            <div class="metrics-section">
              <a-row :gutter="16">
                <a-col
                  v-for="(value, key) in selectedExperiment.metrics"
                  :key="key"
                  :xs="12"
                  :sm="8"
                  :md="6"
                >
                  <a-card class="metric-card" size="small">
                    <a-statistic
                      :title="key"
                      :value="formatMetric(value)"
                      :precision="4"
                    />
                  </a-card>
                </a-col>
              </a-row>
            </div>
          </a-tab-pane>

          <a-tab-pane key="logs" tab="训练日志">
            <div class="log-container">
              <div class="log-header">
                <span class="log-title">训练日志</span>
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

          <a-tab-pane key="resources" tab="资源使用">
            <div class="resources-section">
              <a-row :gutter="16">
                <a-col :xs="24" :sm="8">
                  <a-card class="resource-card" size="small" title="CPU 使用率">
                    <div class="resource-progress">
                      <a-progress
                        :percent="75"
                        :stroke-color="{ from: '#108ee9', to: '#87d068' }"
                      />
                      <span class="resource-text">75%</span>
                    </div>
                  </a-card>
                </a-col>
                <a-col :xs="24" :sm="8">
                  <a-card class="resource-card" size="small" title="内存使用率">
                    <div class="resource-progress">
                      <a-progress
                        :percent="60"
                        :stroke-color="{ from: '#108ee9', to: '#87d068' }"
                      />
                      <span class="resource-text">60%</span>
                    </div>
                  </a-card>
                </a-col>
                <a-col :xs="24" :sm="8">
                  <a-card class="resource-card" size="small" title="GPU 使用率">
                    <div class="resource-progress">
                      <a-progress
                        :percent="90"
                        :stroke-color="{ from: '#108ee9', to: '#87d068' }"
                      />
                      <span class="resource-text">90%</span>
                    </div>
                  </a-card>
                </a-col>
              </a-row>
            </div>
          </a-tab-pane>
        </a-tabs>
      </div>
    </a-modal>

    <!-- 指标对比模态框 -->
    <a-modal
      v-model:open="metricsModalVisible"
      title="训练指标图表"
      width="1200px"
      :footer="null"
      class="sci-fi-modal metrics-modal"
    >
      <div v-if="selectedExperiment" class="metrics-charts">
        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-card title="Loss 曲线" size="small" class="chart-card">
              <div class="chart-placeholder">
                <div class="chart-info">
                  <LineChartOutlined class="chart-icon" />
                  <span>Loss 训练曲线图表</span>
                </div>
              </div>
            </a-card>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-card title="Accuracy 曲线" size="small" class="chart-card">
              <div class="chart-placeholder">
                <div class="chart-info">
                  <LineChartOutlined class="chart-icon" />
                  <span>Accuracy 训练曲线图表</span>
                </div>
              </div>
            </a-card>
          </a-col>
        </a-row>
        <a-row :gutter="16" style="margin-top: 16px">
          <a-col :xs="24" :sm="12">
            <a-card title="学习率变化" size="small" class="chart-card">
              <div class="chart-placeholder">
                <div class="chart-info">
                  <LineChartOutlined class="chart-icon" />
                  <span>学习率变化曲线图表</span>
                </div>
              </div>
            </a-card>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-card title="GPU 内存使用" size="small" class="chart-card">
              <div class="chart-placeholder">
                <div class="chart-info">
                  <AreaChartOutlined class="chart-icon" />
                  <span>GPU 内存使用图表</span>
                </div>
              </div>
            </a-card>
          </a-col>
        </a-row>
      </div>
    </a-modal>

    <!-- 克隆实验模态框 -->
    <a-modal
      v-model:open="cloneModalVisible"
      title="克隆实验"
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
        <a-form-item label="新实验名称" name="name">
          <a-input
            v-model:value="cloneForm.name"
            placeholder="请输入新实验名称"
            class="form-input"
          />
        </a-form-item>
        <a-form-item label="项目" name="project">
          <a-select
            v-model:value="cloneForm.project"
            placeholder="选择项目"
            class="form-select"
          >
            <a-select-option value="image-classification">
              图像分类
            </a-select-option>
            <a-select-option value="nlp-sentiment">情感分析</a-select-option>
            <a-select-option value="recommendation">推荐系统</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述" name="description">
          <a-textarea
            v-model:value="cloneForm.description"
            placeholder="请输入描述"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 实验对比模态框 -->
    <a-modal
      v-model:open="compareModalVisible"
      title="实验对比分析"
      width="1400px"
      :footer="null"
      class="sci-fi-modal compare-modal"
    >
      <div class="compare-content">
        <!-- 对比头部 -->
        <div class="compare-header">
          <div class="compare-title">
            <SwapOutlined class="compare-icon" />
            <span>正在对比 {{ compareExperiments.length }} 个实验</span>
          </div>
          <div class="compare-actions">
            <a-button @click="exportCompareReport" class="export-btn">
              <ExportOutlined />
              导出报告
            </a-button>
            <a-button @click="addMoreExperiments" class="add-more-btn">
              <PlusOutlined />
              添加更多
            </a-button>
          </div>
        </div>

        <a-tabs v-model:activeKey="compareTabKey" class="compare-tabs">
          <!-- 基础信息对比 -->
          <a-tab-pane key="basic" tab="基础信息">
            <div class="compare-basic">
              <a-table
                :columns="basicCompareColumns"
                :data-source="basicCompareData"
                :pagination="false"
                size="small"
                class="compare-table"
              >
                <template #name="{ record }">
                  <div class="experiment-name-cell">
                    <a-tag
                      :color="getExperimentColor(record.experimentIndex)"
                      class="experiment-tag"
                    >
                      实验 {{ record.experimentIndex + 1 }}
                    </a-tag>
                    <span class="experiment-name">{{ record.name }}</span>
                  </div>
                </template>
                <template #status="{ record }">
                  <a-tag
                    :color="getStatusColor(record.status)"
                    class="status-tag"
                  >
                    <component :is="getStatusIcon(record.status)" />
                    {{ getStatusText(record.status) }}
                  </a-tag>
                </template>
              </a-table>
            </div>
          </a-tab-pane>

          <!-- 训练指标对比 -->
          <a-tab-pane key="metrics" tab="训练指标">
            <div class="compare-metrics">
              <!-- 指标对比表格 -->
              <a-card
                title="指标对比"
                class="metrics-compare-card"
                size="small"
              >
                <a-table
                  :columns="metricsCompareColumns"
                  :data-source="metricsCompareData"
                  :pagination="false"
                  size="small"
                  class="metrics-compare-table"
                >
                  <template #metric="{ record }">
                    <strong>{{ record.metric }}</strong>
                  </template>
                  <template #value="{ record, index }">
                    <div class="metric-value-cell">
                      <span
                        class="metric-value"
                        :class="{
                          'best-value': record.bestIndex === index - 1,
                          'worst-value': record.worstIndex === index - 1,
                        }"
                      >
                        {{ record.values[index - 1] }}
                      </span>
                      <div
                        v-if="record.bestIndex === index - 1"
                        class="best-indicator"
                      >
                        <CrownOutlined />
                      </div>
                    </div>
                  </template>
                </a-table>
              </a-card>

              <!-- 指标对比图表 -->
              <a-row :gutter="16" style="margin-top: 16px">
                <a-col :xs="24" :sm="12">
                  <a-card title="准确率对比" size="small" class="chart-card">
                    <div class="metric-chart">
                      <div class="chart-bars">
                        <div
                          v-for="(exp, index) in compareExperiments"
                          :key="exp.id"
                          class="chart-bar"
                        >
                          <div class="bar-container">
                            <div
                              class="bar"
                              :style="{
                                height: `${exp.metrics.accuracy * 100}%`,
                                backgroundColor: getExperimentColor(index),
                              }"
                            ></div>
                          </div>
                          <div class="bar-label">
                            {{ formatMetric(exp.metrics.accuracy) }}
                          </div>
                          <div class="bar-name">实验{{ index + 1 }}</div>
                        </div>
                      </div>
                    </div>
                  </a-card>
                </a-col>
                <a-col :xs="24" :sm="12">
                  <a-card title="损失对比" size="small" class="chart-card">
                    <div class="metric-chart">
                      <div class="chart-bars">
                        <div
                          v-for="(exp, index) in compareExperiments"
                          :key="exp.id"
                          class="chart-bar"
                        >
                          <div class="bar-container">
                            <div
                              class="bar"
                              :style="{
                                height: `${(1 - exp.metrics.loss) * 100}%`,
                                backgroundColor: getExperimentColor(index),
                              }"
                            ></div>
                          </div>
                          <div class="bar-label">
                            {{ formatMetric(exp.metrics.loss) }}
                          </div>
                          <div class="bar-name">实验{{ index + 1 }}</div>
                        </div>
                      </div>
                    </div>
                  </a-card>
                </a-col>
              </a-row>
            </div>
          </a-tab-pane>

          <!-- 超参数对比 -->
          <a-tab-pane key="hyperparams" tab="超参数">
            <div class="compare-hyperparams">
              <a-table
                :columns="hyperparamsCompareColumns"
                :data-source="hyperparamsCompareData"
                :pagination="false"
                size="small"
                class="hyperparams-compare-table"
              >
                <template #param="{ record }">
                  <strong>{{ record.param }}</strong>
                </template>
                <template #value="{ record, index }">
                  <span
                    class="hyperparam-value"
                    :class="{
                      'different-value':
                        !record.allSame &&
                        record.values[index - 1] !== record.values[0],
                    }"
                  >
                    {{ record.values[index - 1] }}
                  </span>
                </template>
              </a-table>

              <!-- 超参数差异提示 -->
              <a-alert
                v-if="hyperparamsDifferences.length > 0"
                type="info"
                show-icon
                class="hyperparams-alert"
              >
                <template #message>
                  <span>发现以下超参数存在差异：</span>
                  <a-tag
                    v-for="diff in hyperparamsDifferences"
                    :key="diff"
                    color="orange"
                    style="margin-left: 8px"
                  >
                    {{ diff }}
                  </a-tag>
                </template>
              </a-alert>
            </div>
          </a-tab-pane>

          <!-- 性能分析 -->
          <a-tab-pane key="performance" tab="性能分析">
            <div class="compare-performance">
              <!-- 性能排名 -->
              <a-card
                title="综合性能排名"
                class="performance-ranking-card"
                size="small"
              >
                <div class="ranking-list">
                  <div
                    v-for="(ranking, index) in performanceRanking"
                    :key="ranking.id"
                    class="ranking-item"
                    :class="`rank-${index + 1}`"
                  >
                    <div class="rank-number">
                      <span class="rank-text">{{ index + 1 }}</span>
                      <component :is="getRankIcon(index)" class="rank-icon" />
                    </div>
                    <div class="experiment-info">
                      <div class="experiment-name">{{ ranking.name }}</div>
                      <div class="experiment-score">
                        综合得分: {{ ranking.score.toFixed(2) }}
                      </div>
                    </div>
                    <div class="experiment-highlights">
                      <a-tag
                        v-for="highlight in ranking.highlights"
                        :key="highlight"
                        :color="getHighlightColor(highlight)"
                        size="small"
                      >
                        {{ highlight }}
                      </a-tag>
                    </div>
                  </div>
                </div>
              </a-card>

              <!-- 性能维度对比 -->
              <a-row :gutter="16" style="margin-top: 16px">
                <a-col :xs="24" :sm="8">
                  <a-card
                    title="准确性表现"
                    size="small"
                    class="performance-card"
                  >
                    <div class="performance-metric">
                      <div
                        v-for="(exp, index) in sortedByAccuracy"
                        :key="exp.id"
                        class="performance-item"
                      >
                        <div class="item-rank">{{ index + 1 }}</div>
                        <div class="item-name">
                          实验{{ getExperimentIndex(exp.id) + 1 }}
                        </div>
                        <div class="item-value">
                          {{ formatMetric(exp.metrics.accuracy) }}
                        </div>
                      </div>
                    </div>
                  </a-card>
                </a-col>
                <a-col :xs="24" :sm="8">
                  <a-card
                    title="训练效率"
                    size="small"
                    class="performance-card"
                  >
                    <div class="performance-metric">
                      <div
                        v-for="(exp, index) in sortedByEfficiency"
                        :key="exp.id"
                        class="performance-item"
                      >
                        <div class="item-rank">{{ index + 1 }}</div>
                        <div class="item-name">
                          实验{{ getExperimentIndex(exp.id) + 1 }}
                        </div>
                        <div class="item-value">
                          {{ formatDuration(exp.duration) }}
                        </div>
                      </div>
                    </div>
                  </a-card>
                </a-col>
                <a-col :xs="24" :sm="8">
                  <a-card title="稳定性" size="small" class="performance-card">
                    <div class="performance-metric">
                      <div
                        v-for="(exp, index) in sortedByStability"
                        :key="exp.id"
                        class="performance-item"
                      >
                        <div class="item-rank">{{ index + 1 }}</div>
                        <div class="item-name">
                          实验{{ getExperimentIndex(exp.id) + 1 }}
                        </div>
                        <div class="item-value">
                          {{ formatMetric(exp.metrics.f1_score) }}
                        </div>
                      </div>
                    </div>
                  </a-card>
                </a-col>
              </a-row>
            </div>
          </a-tab-pane>

          <!-- 对比报告 -->
          <a-tab-pane key="report" tab="对比报告">
            <div class="compare-report">
              <a-card title="实验对比总结" class="report-card">
                <div class="report-content">
                  <div class="report-section">
                    <h4>📊 实验概览</h4>
                    <p>
                      本次对比分析了 {{ compareExperiments.length }} 个实验，
                      涵盖了 {{ uniqueProjects.length }} 个项目类型， 使用了
                      {{ uniqueFrameworks.length }} 种不同的框架。
                    </p>
                  </div>

                  <div class="report-section">
                    <h4>🏆 最佳表现</h4>
                    <p>
                      <strong>{{ bestExperiment?.name }}</strong>
                      在综合指标上表现最佳， 准确率达到
                      {{
                        formatMetric(bestExperiment?.metrics?.accuracy || 0)
                      }}， 是本次对比中的推荐方案。
                    </p>
                  </div>

                  <div class="report-section">
                    <h4>⚡ 效率分析</h4>
                    <p>
                      训练时间最短的是
                      <strong>{{ fastestExperiment?.name }}</strong> ({{
                        formatDuration(fastestExperiment?.duration || 0)
                      }})， 而
                      <strong>{{ slowestExperiment?.name }}</strong> 耗时最长
                      ({{ formatDuration(slowestExperiment?.duration || 0) }})。
                    </p>
                  </div>

                  <div class="report-section">
                    <h4>📈 关键发现</h4>
                    <ul class="findings-list">
                      <li v-for="finding in keyFindings" :key="finding">
                        {{ finding }}
                      </li>
                    </ul>
                  </div>

                  <div class="report-section">
                    <h4>💡 优化建议</h4>
                    <ul class="suggestions-list">
                      <li
                        v-for="suggestion in optimizationSuggestions"
                        :key="suggestion"
                      >
                        {{ suggestion }}
                      </li>
                    </ul>
                  </div>
                </div>
              </a-card>
            </div>
          </a-tab-pane>
        </a-tabs>
      </div>
    </a-modal>

    <!-- 选择实验对比模态框 -->
    <a-modal
      v-model:open="selectCompareModalVisible"
      title="选择要对比的实验"
      width="800px"
      @ok="handleSelectCompareSubmit"
      @cancel="handleSelectCompareCancel"
      class="sci-fi-modal"
    >
      <div class="select-compare-content">
        <a-alert
          message="请选择2-5个实验进行对比分析"
          type="info"
          show-icon
          style="margin-bottom: 16px"
        />

        <a-table
          :columns="selectCompareColumns"
          :data-source="experiments"
          :row-selection="selectCompareRowSelection"
          :pagination="{ pageSize: 5 }"
          size="small"
          row-key="id"
        >
          <template #status="{ record }">
            <a-tag :color="getStatusColor(record.status)" size="small">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>
          <template #metrics="{ record }">
            <span>{{ formatMetric(record.metrics.accuracy) }}</span>
          </template>
        </a-table>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import type { FormInstance, TableColumnsType } from 'ant-design-vue';
import {
  ExperimentOutlined,
  PlusOutlined,
  ReloadOutlined,
  EyeOutlined,
  MoreOutlined,
  PlayCircleOutlined,
  StopOutlined,
  CopyOutlined,
  DeleteOutlined,
  ExportOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  CloseCircleOutlined,
  LineChartOutlined,
  AreaChartOutlined,
  SwapOutlined,
  CrownOutlined,
  TrophyOutlined,
  StarOutlined,
} from '@ant-design/icons-vue';

// ===== 类型定义 =====
interface ExperimentMetrics {
  accuracy: number;
  loss: number;
  f1_score: number;
  precision: number;
  recall: number;
}

interface HyperParams {
  learningRate: number;
  batchSize: number;
  epochs: number;
}

interface ExperimentItem {
  id: string;
  name: string;
  project: string;
  status: 'running' | 'completed' | 'failed' | 'stopped';
  framework: string;
  creator: string;
  createTime: string;
  duration: number; // 秒
  hyperParams: HyperParams;
  metrics: ExperimentMetrics;
  description?: string;
  codeRepository?: string;
}

interface CreateForm {
  name: string;
  project: string;
  framework: string;
  learningRate: number;
  batchSize: number;
  epochs: number;
  cpu: number;
  memory: number;
  gpu: number;
  description: string;
  codeRepository: string;
}

interface CloneForm {
  name: string;
  project: string;
  description: string;
}

interface PerformanceRanking {
  id: string;
  name: string;
  score: number;
  highlights: string[];
}

// ===== 响应式数据 =====
const loading = ref<boolean>(false);
const createModalVisible = ref<boolean>(false);
const detailModalVisible = ref<boolean>(false);
const metricsModalVisible = ref<boolean>(false);
const cloneModalVisible = ref<boolean>(false);
const compareModalVisible = ref<boolean>(false);
const selectCompareModalVisible = ref<boolean>(false);
const createLoading = ref<boolean>(false);
const cloneLoading = ref<boolean>(false);

const filterStatus = ref<string>('');
const filterProject = ref<string>('');
const searchKeyword = ref<string>('');
const activeTabKey = ref<string>('overview');
const compareTabKey = ref<string>('basic');

const selectedExperiment = ref<ExperimentItem | null>(null);
const selectedExperimentIds = ref<string[]>([]);
const compareExperiments = ref<ExperimentItem[]>([]);
const selectCompareIds = ref<string[]>([]);

const logs = ref<string[]>([
  '2024-06-23 10:30:15 INFO: Experiment started',
  '2024-06-23 10:30:16 INFO: Loading dataset...',
  '2024-06-23 10:30:18 INFO: Dataset loaded successfully, 10000 samples',
  '2024-06-23 10:30:20 INFO: Model initialized',
  '2024-06-23 10:30:25 INFO: Training started...',
  '2024-06-23 10:31:00 INFO: Epoch 1/100 - loss: 0.6524 - accuracy: 0.7234',
  '2024-06-23 10:31:35 INFO: Epoch 2/100 - loss: 0.5892 - accuracy: 0.7856',
  '2024-06-23 10:32:10 INFO: Epoch 3/100 - loss: 0.5234 - accuracy: 0.8123',
]);

// ===== 表单引用 =====
const createFormRef = ref<FormInstance>();
const cloneFormRef = ref<FormInstance>();

// ===== 表单数据 =====
const createForm = reactive<CreateForm>({
  name: '',
  project: '',
  framework: 'tensorflow',
  learningRate: 0.001,
  batchSize: 32,
  epochs: 100,
  cpu: 4,
  memory: 8,
  gpu: 1,
  description: '',
  codeRepository: '',
});

const cloneForm = reactive<CloneForm>({
  name: '',
  project: '',
  description: '',
});

// ===== 配置数据 =====
const STATUS_CONFIG = {
  running: { color: 'processing', text: '运行中', icon: PlayCircleOutlined },
  completed: { color: 'success', text: '已完成', icon: CheckCircleOutlined },
  failed: { color: 'error', text: '失败', icon: CloseCircleOutlined },
  stopped: { color: 'default', text: '已停止', icon: StopOutlined },
} as const;

// 实验颜色配置
const EXPERIMENT_COLORS = [
  '#1890ff',
  '#52c41a',
  '#faad14',
  '#f5222d',
  '#722ed1',
  '#13c2c2',
  '#eb2f96',
  '#fa8c16',
  '#a0d911',
  '#2f54eb',
];

// ===== 模拟数据 =====
const experiments = ref<ExperimentItem[]>([
  {
    id: 'exp-001',
    name: 'resnet50-classification',
    project: 'image-classification',
    status: 'running',
    framework: 'tensorflow',
    creator: 'admin',
    createTime: '2024-06-23 09:30:00',
    duration: 3600,
    hyperParams: { learningRate: 0.001, batchSize: 32, epochs: 100 },
    metrics: {
      accuracy: 0.8956,
      loss: 0.2543,
      f1_score: 0.8834,
      precision: 0.8921,
      recall: 0.8748,
    },
    description: 'ResNet50 图像分类模型训练',
    codeRepository: 'https://github.com/username/resnet-classification.git',
  },
  {
    id: 'exp-002',
    name: 'bert-sentiment-analysis',
    project: 'nlp-sentiment',
    status: 'completed',
    framework: 'pytorch',
    creator: 'researcher',
    createTime: '2024-06-23 08:15:00',
    duration: 7200,
    hyperParams: { learningRate: 0.0002, batchSize: 16, epochs: 50 },
    metrics: {
      accuracy: 0.9234,
      loss: 0.1876,
      f1_score: 0.9156,
      precision: 0.9287,
      recall: 0.9028,
    },
    description: 'BERT 情感分析模型',
    codeRepository: 'https://github.com/username/bert-sentiment.git',
  },
  {
    id: 'exp-003',
    name: 'collaborative-filtering',
    project: 'recommendation',
    status: 'failed',
    framework: 'sklearn',
    creator: 'developer',
    createTime: '2024-06-23 10:00:00',
    duration: 1800,
    hyperParams: { learningRate: 0.01, batchSize: 64, epochs: 200 },
    metrics: {
      accuracy: 0.7234,
      loss: 0.4567,
      f1_score: 0.7012,
      precision: 0.7345,
      recall: 0.6892,
    },
    description: '协同过滤推荐算法',
  },
  {
    id: 'exp-004',
    name: 'lstm-time-series',
    project: 'time-series',
    status: 'completed',
    framework: 'tensorflow',
    creator: 'admin',
    createTime: '2024-06-22 16:30:00',
    duration: 5400,
    hyperParams: { learningRate: 0.005, batchSize: 128, epochs: 150 },
    metrics: {
      accuracy: 0.8567,
      loss: 0.3234,
      f1_score: 0.8423,
      precision: 0.8612,
      recall: 0.8238,
    },
    description: 'LSTM 时间序列预测模型',
    codeRepository: 'https://github.com/username/lstm-timeseries.git',
  },
  {
    id: 'exp-005',
    name: 'xgboost-regression',
    project: 'regression',
    status: 'stopped',
    framework: 'xgboost',
    creator: 'developer',
    createTime: '2024-06-23 09:45:00',
    duration: 2700,
    hyperParams: { learningRate: 0.1, batchSize: 256, epochs: 300 },
    metrics: {
      accuracy: 0.7845,
      loss: 0.3876,
      f1_score: 0.7623,
      precision: 0.7912,
      recall: 0.7342,
    },
    description: 'XGBoost 回归模型',
  },
]);

// ===== 表单验证规则 =====
const createFormRules = {
  name: [
    { required: true, message: '请输入实验名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
    {
      pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
      message: '名称只能包含小写字母、数字和连字符',
      trigger: 'blur',
    },
  ],
  project: [{ required: true, message: '请选择项目', trigger: 'change' }],
  framework: [{ required: true, message: '请选择框架', trigger: 'change' }],
  learningRate: [{ required: true, message: '请输入学习率', trigger: 'blur' }],
  batchSize: [{ required: true, message: '请输入批次大小', trigger: 'blur' }],
  epochs: [{ required: true, message: '请输入训练轮数', trigger: 'blur' }],
  cpu: [{ required: true, message: '请输入 CPU 核数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
};

const cloneFormRules = {
  name: [
    { required: true, message: '请输入新实验名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
    {
      pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
      message: '名称只能包含小写字母、数字和连字符',
      trigger: 'blur',
    },
  ],
  project: [{ required: true, message: '请选择项目', trigger: 'change' }],
};

// ===== 表格列配置 =====
const columns: TableColumnsType<ExperimentItem> = [
  {
    title: '实验名称',
    dataIndex: 'name',
    key: 'name',
    width: 200,
    ellipsis: true,
  },
  {
    title: '项目',
    dataIndex: 'project',
    key: 'project',
    width: 120,
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
    slots: { customRender: 'status' },
  },
  {
    title: '框架',
    dataIndex: 'framework',
    key: 'framework',
    width: 100,
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
    width: 100,
  },
  {
    title: '训练指标',
    key: 'metrics',
    width: 200,
    slots: { customRender: 'metrics' },
  },
  {
    title: '持续时间',
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
    width: 200,
    fixed: 'right',
    slots: { customRender: 'action' },
  },
];

// 对比相关的表格列配置
const basicCompareColumns = [
  { title: '实验名称', key: 'name', slots: { customRender: 'name' } },
  { title: '项目', dataIndex: 'project', key: 'project' },
  { title: '状态', key: 'status', slots: { customRender: 'status' } },
  { title: '框架', dataIndex: 'framework', key: 'framework' },
  { title: '创建者', dataIndex: 'creator', key: 'creator' },
  { title: '创建时间', dataIndex: 'createTime', key: 'createTime' },
];

const selectCompareColumns = [
  { title: '实验名称', dataIndex: 'name', key: 'name' },
  { title: '项目', dataIndex: 'project', key: 'project' },
  { title: '状态', key: 'status', slots: { customRender: 'status' } },
  { title: '框架', dataIndex: 'framework', key: 'framework' },
  { title: '准确率', key: 'metrics', slots: { customRender: 'metrics' } },
];

// ===== 分页配置 =====
const paginationConfig = {
  total: computed(() => filteredExperiments.value.length),
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// ===== 行选择配置 =====
const rowSelection = {
  selectedRowKeys: selectedExperimentIds,
  onChange: (selectedRowKeys: string[]) => {
    selectedExperimentIds.value = selectedRowKeys;
  },
  getCheckboxProps: (record: ExperimentItem) => ({
    disabled: record.status === 'running',
  }),
};

const selectCompareRowSelection = {
  selectedRowKeys: selectCompareIds,
  onChange: (selectedRowKeys: string[]) => {
    selectCompareIds.value = selectedRowKeys;
  },
  getCheckboxProps: () => ({
    disabled: false,
  }),
};

// ===== 计算属性 =====
const filteredExperiments = computed(() => {
  let result = experiments.value;

  if (filterStatus.value) {
    result = result.filter((item) => item.status === filterStatus.value);
  }

  if (filterProject.value) {
    result = result.filter((item) => item.project === filterProject.value);
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
  () => experiments.value.filter((exp) => exp.status === 'running').length,
);

const completedCount = computed(
  () => experiments.value.filter((exp) => exp.status === 'completed').length,
);

const failedCount = computed(
  () => experiments.value.filter((exp) => exp.status === 'failed').length,
);

// 对比相关计算属性
const basicCompareData = computed(() => {
  return compareExperiments.value.map((exp, index) => ({
    ...exp,
    experimentIndex: index,
  }));
});

const metricsCompareColumns = computed(() => {
  const columns = [
    { title: '指标', key: 'metric', slots: { customRender: 'metric' } },
  ];
  compareExperiments.value.forEach((exp, index) => {
    columns.push({
      title: `实验${index + 1}`,
      key: `exp${index}`,
      slots: { customRender: 'value' },
    });
  });
  return columns;
});

const metricsCompareData = computed(() => {
  const metricKeys = ['accuracy', 'loss', 'f1_score', 'precision', 'recall'];
  return metricKeys.map((key) => {
    const values = compareExperiments.value.map((exp) =>
      formatMetric(exp.metrics[key as keyof ExperimentMetrics]),
    );
    const numValues = compareExperiments.value.map(
      (exp) => exp.metrics[key as keyof ExperimentMetrics],
    );

    let bestIndex = 0;
    let worstIndex = 0;

    if (key === 'loss') {
      // 对于 loss，值越小越好
      bestIndex = numValues.indexOf(Math.min(...numValues));
      worstIndex = numValues.indexOf(Math.max(...numValues));
    } else {
      // 对于其他指标，值越大越好
      bestIndex = numValues.indexOf(Math.max(...numValues));
      worstIndex = numValues.indexOf(Math.min(...numValues));
    }

    return {
      metric: key,
      values,
      bestIndex,
      worstIndex,
    };
  });
});

const hyperparamsCompareColumns = computed(() => {
  const columns = [
    { title: '超参数', key: 'param', slots: { customRender: 'param' } },
  ];
  compareExperiments.value.forEach((exp, index) => {
    columns.push({
      title: `实验${index + 1}`,
      key: `exp${index}`,
      slots: { customRender: 'value' },
    });
  });
  return columns;
});

const hyperparamsCompareData = computed(() => {
  const paramKeys = ['learningRate', 'batchSize', 'epochs'];
  return paramKeys.map((key) => {
    const values = compareExperiments.value.map(
      (exp) => exp.hyperParams[key as keyof HyperParams],
    );
    const allSame = values.every((val) => val === values[0]);

    return {
      param: key,
      values,
      allSame,
    };
  });
});

const hyperparamsDifferences = computed(() => {
  return hyperparamsCompareData.value
    .filter((item) => !item.allSame)
    .map((item) => item.param);
});

const performanceRanking = computed((): PerformanceRanking[] => {
  return compareExperiments.value
    .map((exp) => {
      // 计算综合得分
      const score =
        (exp.metrics.accuracy * 0.4 +
          (1 - exp.metrics.loss) * 0.3 +
          exp.metrics.f1_score * 0.3) *
        100;

      // 确定亮点
      const highlights: string[] = [];
      const maxAccuracy = Math.max(
        ...compareExperiments.value.map((e) => e.metrics.accuracy),
      );
      const minLoss = Math.min(
        ...compareExperiments.value.map((e) => e.metrics.loss),
      );
      const minDuration = Math.min(
        ...compareExperiments.value.map((e) => e.duration),
      );

      if (exp.metrics.accuracy === maxAccuracy) highlights.push('最高准确率');
      if (exp.metrics.loss === minLoss) highlights.push('最低损失');
      if (exp.duration === minDuration) highlights.push('最快训练');

      return {
        id: exp.id,
        name: exp.name,
        score,
        highlights,
      };
    })
    .sort((a, b) => b.score - a.score);
});

const sortedByAccuracy = computed(() => {
  return [...compareExperiments.value].sort(
    (a, b) => b.metrics.accuracy - a.metrics.accuracy,
  );
});

const sortedByEfficiency = computed(() => {
  return [...compareExperiments.value].sort((a, b) => a.duration - b.duration);
});

const sortedByStability = computed(() => {
  return [...compareExperiments.value].sort(
    (a, b) => b.metrics.f1_score - a.metrics.f1_score,
  );
});

const uniqueProjects = computed(() => {
  return [...new Set(compareExperiments.value.map((exp) => exp.project))];
});

const uniqueFrameworks = computed(() => {
  return [...new Set(compareExperiments.value.map((exp) => exp.framework))];
});

const bestExperiment = computed(() => {
  return (
    performanceRanking.value[0] &&
    compareExperiments.value.find(
      (exp) => exp.id === performanceRanking.value[0]?.id,
    )
  );
});

const fastestExperiment = computed(() => {
  return sortedByEfficiency.value[0] || null;
});

const slowestExperiment = computed(() => {
  return sortedByEfficiency.value[sortedByEfficiency.value.length - 1];
});

const keyFindings = computed(() => {
  const findings: string[] = [];

  if (compareExperiments.value.length >= 2) {
    const accuracyRange =
      Math.max(...compareExperiments.value.map((e) => e.metrics.accuracy)) -
      Math.min(...compareExperiments.value.map((e) => e.metrics.accuracy));

    if (accuracyRange > 0.1) {
      findings.push('不同实验间准确率差异较大，建议进一步分析超参数影响');
    }

    const durationRange =
      Math.max(...compareExperiments.value.map((e) => e.duration)) -
      Math.min(...compareExperiments.value.map((e) => e.duration));

    if (durationRange > 3600) {
      findings.push('训练时间差异显著，可考虑优化计算资源配置');
    }

    findings.push(`${uniqueFrameworks.value.join('、')} 框架在性能上各有优势`);
  }

  return findings;
});

const optimizationSuggestions = computed(() => {
  const suggestions: string[] = [];

  if (bestExperiment.value) {
    suggestions.push(
      `建议采用 ${bestExperiment.value.name} 的超参数配置作为基准`,
    );
  }

  if (hyperparamsDifferences.value.length > 0) {
    suggestions.push('可通过网格搜索进一步优化差异较大的超参数');
  }

  suggestions.push('建议建立自动化超参数调优流程');
  suggestions.push('考虑使用早停机制避免过拟合');

  return suggestions;
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

const formatMetric = (value: number): string => {
  return value.toFixed(4);
};

const formatDuration = (seconds: number): string => {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  if (hours > 0) {
    return `${hours}h ${minutes}m ${secs}s`;
  } else if (minutes > 0) {
    return `${minutes}m ${secs}s`;
  } else {
    return `${secs}s`;
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

const getExperimentColor = (index: number): string => {
  return EXPERIMENT_COLORS[index % EXPERIMENT_COLORS.length] || '';
};

const getExperimentIndex = (id: string): number => {
  return compareExperiments.value.findIndex((exp) => exp.id === id);
};

const getRankIcon = (index: number) => {
  switch (index) {
    case 0:
      return TrophyOutlined;
    case 1:
      return TrophyOutlined;
    case 2:
      return StarOutlined;
    default:
      return null;
  }
};

const getHighlightColor = (highlight: string): string => {
  const colorMap: Record<string, string> = {
    最高准确率: 'gold',
    最低损失: 'green',
    最快训练: 'blue',
  };
  return colorMap[highlight] || 'default';
};

// ===== 对比相关函数 =====
const handleCompareSelected = (): void => {
  if (selectedExperimentIds.value.length < 2) {
    message.warning('请至少选择2个实验进行对比');
    return;
  }

  if (selectedExperimentIds.value.length > 5) {
    message.warning('最多只能对比5个实验');
    return;
  }

  const selectedExps = experiments.value.filter((exp) =>
    selectedExperimentIds.value.includes(exp.id),
  );

  compareExperiments.value = selectedExps;
  compareTabKey.value = 'basic';
  compareModalVisible.value = true;
};

const handleCompare = (record: ExperimentItem): void => {
  selectCompareIds.value = [record.id];
  selectCompareModalVisible.value = true;
};

const handleSelectCompareSubmit = (): void => {
  if (selectCompareIds.value.length < 2) {
    message.warning('请至少选择2个实验进行对比');
    return;
  }

  if (selectCompareIds.value.length > 5) {
    message.warning('最多只能对比5个实验');
    return;
  }

  const selectedExps = experiments.value.filter((exp) =>
    selectCompareIds.value.includes(exp.id),
  );

  compareExperiments.value = selectedExps;
  selectCompareModalVisible.value = false;
  compareModalVisible.value = true;
  compareTabKey.value = 'basic';
};

const handleSelectCompareCancel = (): void => {
  selectCompareModalVisible.value = false;
  selectCompareIds.value = [];
};

const addMoreExperiments = (): void => {
  selectCompareIds.value = compareExperiments.value.map((exp) => exp.id);
  selectCompareModalVisible.value = true;
};

const exportCompareReport = (): void => {
  // 生成对比报告
  const report = {
    title: '实验对比报告',
    date: new Date().toLocaleString(),
    experiments: compareExperiments.value.map((exp, index) => ({
      index: index + 1,
      name: exp.name,
      project: exp.project,
      framework: exp.framework,
      metrics: exp.metrics,
      hyperParams: exp.hyperParams,
      duration: exp.duration,
    })),
    ranking: performanceRanking.value,
    findings: keyFindings.value,
    suggestions: optimizationSuggestions.value,
  };

  // 模拟下载
  const blob = new Blob([JSON.stringify(report, null, 2)], {
    type: 'application/json',
  });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `experiment-compare-report-${Date.now()}.json`;
  a.click();
  URL.revokeObjectURL(url);

  message.success('对比报告导出成功');
};

// ===== 原有事件处理函数 =====
const showCreateModal = (): void => {
  createModalVisible.value = true;
};

const handleCreateSubmit = async (): Promise<void> => {
  try {
    await createFormRef.value?.validate();
    createLoading.value = true;

    // 模拟 API 调用
    await new Promise((resolve) => setTimeout(resolve, 2000));

    const newExperiment: ExperimentItem = {
      id: `exp-${Date.now()}`,
      name: createForm.name,
      project: createForm.project,
      status: 'running',
      framework: createForm.framework,
      creator: 'current-user',
      createTime: new Date().toLocaleString(),
      duration: 0,
      hyperParams: {
        learningRate: createForm.learningRate,
        batchSize: createForm.batchSize,
        epochs: createForm.epochs,
      },
      metrics: {
        accuracy: 0,
        loss: 0,
        f1_score: 0,
        precision: 0,
        recall: 0,
      },
      description: createForm.description,
      codeRepository: createForm.codeRepository,
    };

    experiments.value.unshift(newExperiment);
    createModalVisible.value = false;
    message.success('实验创建成功');

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

const viewDetails = (record: ExperimentItem): void => {
  selectedExperiment.value = record;
  activeTabKey.value = 'overview';
  detailModalVisible.value = true;
};

const viewMetrics = (record: ExperimentItem): void => {
  selectedExperiment.value = record;
  metricsModalVisible.value = true;
};

const handleMenuAction = (key: string, record: ExperimentItem): void => {
  const actions = {
    stop: () => handleStop(record),
    restart: () => handleRestart(record),
    clone: () => handleClone(record),
    compare: () => handleCompare(record),
    delete: () => handleDelete(record),
  };

  const action = actions[key as keyof typeof actions];
  if (action) {
    action();
  }
};

const handleStop = async (record: ExperimentItem): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const index = experiments.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      experiments.value[index]!.status = 'stopped';
    }
    message.success('实验停止成功');
  } catch (error) {
    message.error('停止失败');
  } finally {
    loading.value = false;
  }
};

const handleRestart = async (record: ExperimentItem): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 2000));
    const index = experiments.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      experiments.value[index]!.status = 'running';
      experiments.value[index]!.duration = 0;
    }
    message.success('实验重新运行成功');
  } catch (error) {
    message.error('重新运行失败');
  } finally {
    loading.value = false;
  }
};

const handleClone = (record: ExperimentItem): void => {
  cloneForm.name = `${record.name}-copy`;
  cloneForm.project = record.project;
  cloneForm.description = record.description || '';
  selectedExperiment.value = record;
  cloneModalVisible.value = true;
};

const handleCloneSubmit = async (): Promise<void> => {
  try {
    await cloneFormRef.value?.validate();
    cloneLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 1500));

    if (selectedExperiment.value) {
      const clonedExperiment: ExperimentItem = {
        ...selectedExperiment.value,
        id: `exp-${Date.now()}`,
        name: cloneForm.name,
        project: cloneForm.project,
        status: 'running',
        createTime: new Date().toLocaleString(),
        duration: 0,
        description: cloneForm.description,
        metrics: {
          accuracy: 0,
          loss: 0,
          f1_score: 0,
          precision: 0,
          recall: 0,
        },
      };

      experiments.value.unshift(clonedExperiment);
      cloneModalVisible.value = false;
      message.success('实验克隆成功');
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

const handleDelete = (record: ExperimentItem): void => {
  const deleteConfirm = () => {
    const index = experiments.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      experiments.value.splice(index, 1);
      message.success('实验删除成功');
    }
  };

  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除实验 "${record.name}" 吗？此操作不可恢复。`,
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
.experiment-container {
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
  text-align: center;
  border-radius: 8px !important;
}

.stat-card :deep(.ant-statistic-title) {
  font-weight: 500 !important;
  margin-bottom: 8px !important;
}

/* ===== 按钮样式 ===== */
.create-btn,
.compare-btn {
  border: none !important;
  height: 40px !important;
  padding: 0 24px !important;
  font-size: 14px !important;
  font-weight: 500 !important;
  border-radius: 6px !important;
  transition: all 0.3s ease !important;
}

.create-btn:hover,
.compare-btn:hover {
  background: #1890ff !important;
  transform: translateY(-1px);
}

.compare-btn {
  background: #52c41a !important;
}

.compare-btn:hover {
  background: #389e0d !important;
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
  background: #1890ff;
}

.indicator-completed {
  background: #52c41a;
}

.indicator-failed {
  background: #ff4d4f;
}

.indicator-stopped {
  background: #8c8c8c;
}

/* ===== 指标信息 ===== */
.metrics-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.metric-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.metric-item:hover {
  color: #1890ff;
}

.metric-label {
  font-weight: 500;
}

.metric-value {
  font-weight: 600;
  color: #1890ff;
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

/* ===== 详情页样式 ===== */
.detail-content {
  max-height: 600px;
  overflow-y: auto;
}

.detail-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 16px !important;
}

.detail-tabs :deep(.ant-tabs-tab) {
  font-weight: 500 !important;
}

.repo-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  text-decoration: none;
  transition: all 0.3s ease;
}

.repo-link:hover {
  text-decoration: underline;
}

/* ===== 指标部分 ===== */
.metrics-section {
  margin-top: 16px;
}

.metric-card {
  text-align: center;
  border-radius: 6px !important;
}

.metric-card :deep(.ant-statistic-title) {
  font-weight: 500 !important;
}

/* ===== 资源部分 ===== */
.resources-section {
  margin-top: 16px;
}

.resource-card {
  border-radius: 6px !important;
}

.resource-progress {
  display: flex;
  align-items: center;
  gap: 12px;
}

.resource-text {
  font-weight: 600;
  font-size: 14px;
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

/* ===== 图表样式 ===== */
.metrics-charts {
  padding: 16px 0;
}

.chart-card {
  border-radius: 6px !important;
}

.chart-placeholder {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.chart-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: #8c8c8c;
}

.chart-icon {
  font-size: 32px;
}

/* ===== 对比样式 ===== */
.compare-modal :deep(.ant-modal-body) {
  padding: 16px !important;
}

.compare-content {
  min-height: 500px;
}

.compare-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px 16px;
  border-radius: 8px;
}

.compare-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
  font-weight: 600;
}

.compare-icon {
  font-size: 20px;
  color: #1890ff;
}

.compare-actions {
  display: flex;
  gap: 8px;
}

.export-btn,
.add-more-btn {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

.compare-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 20px !important;
}

.compare-table,
.metrics-compare-table,
.hyperparams-compare-table {
  border-radius: 6px !important;
}

.experiment-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.experiment-tag {
  border-radius: 4px !important;
  font-size: 11px !important;
  padding: 2px 6px !important;
}

.experiment-name {
  font-weight: 500;
}

.metric-value-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
}

.metric-value {
  font-weight: 600;
}

.best-value {
  color: #52c41a !important;
  font-weight: 700 !important;
}

.worst-value {
  color: #ff4d4f !important;
}

.best-indicator {
  color: #faad14;
  font-size: 12px;
}

.hyperparam-value {
  font-weight: 500;
}

.different-value {
  color: #fa8c16 !important;
  font-weight: 600 !important;
}

.hyperparams-alert {
  margin-top: 16px;
}

.metric-chart {
  height: 200px;
  padding: 16px;
}

.chart-bars {
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  height: 100%;
  gap: 16px;
}

.chart-bar {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  gap: 8px;
}

.bar-container {
  height: 120px;
  width: 100%;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.bar {
  width: 60%;
  min-height: 4px;
  border-radius: 4px 4px 0 0;
  transition: all 0.3s ease;
}

.bar-label {
  font-weight: 600;
  font-size: 12px;
}

.bar-name {
  font-size: 11px;
  color: #8c8c8c;
}

.performance-ranking-card {
  margin-bottom: 16px;
}

.ranking-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  transition: all 0.3s ease;
}

.ranking-item:hover {
  border-color: #1890ff;
}

.rank-1 {
  border-color: #faad14 !important;
}

.rank-2 {
  border-color: #bfbfbf !important;
}

.rank-3 {
  border-color: #d4b106 !important;
}

.rank-number {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 60px;
}

.rank-text {
  font-size: 18px;
  font-weight: 700;
}

.rank-icon {
  font-size: 16px;
}

.rank-1 .rank-text,
.rank-1 .rank-icon {
  color: #faad14;
}

.rank-2 .rank-text,
.rank-2 .rank-icon {
  color: #bfbfbf;
}

.rank-3 .rank-text,
.rank-3 .rank-icon {
  color: #d4b106;
}

.experiment-info {
  flex: 1;
}

.experiment-name {
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
}

.experiment-score {
  font-size: 12px;
  color: #8c8c8c;
}

.experiment-highlights {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.performance-card {
  border-radius: 6px !important;
}

.performance-metric {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.performance-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.performance-item:hover {
  color: #1890ff;
}

.item-rank {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 600;
}

.item-name {
  flex: 1;
  font-weight: 500;
  font-size: 12px;
}

.item-value {
  font-weight: 600;
  font-size: 12px;
}

.report-card {
  border-radius: 8px !important;
}

.report-content {
  line-height: 1.6;
}

.report-section {
  margin-bottom: 24px;
}

.report-section h4 {
  margin-bottom: 8px;
  color: #1890ff;
  font-weight: 600;
}

.findings-list,
.suggestions-list {
  margin: 8px 0;
  padding-left: 20px;
}

.findings-list li,
.suggestions-list li {
  margin-bottom: 4px;
}

.select-compare-content {
  max-height: 400px;
  overflow-y: auto;
}

/* ===== 响应式设计 ===== */
@media (max-width: 768px) {
  .experiment-container {
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
  .compare-btn,
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

  .stats-section .ant-col {
    margin-bottom: 16px;
  }

  .compare-header {
    flex-direction: column;
    gap: 12px;
  }

  .chart-bars {
    gap: 8px;
  }

  .ranking-item {
    flex-direction: column;
    text-align: center;
    gap: 8px;
  }

  .experiment-highlights {
    justify-content: center;
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

  .metrics-info {
    gap: 2px;
  }

  .metric-item {
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

  .chart-placeholder {
    height: 150px;
  }

  .chart-icon {
    font-size: 24px;
  }

  .bar {
    width: 80%;
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
