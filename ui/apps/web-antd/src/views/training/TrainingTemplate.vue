<template>
  <div class="training-template-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <CodeOutlined class="title-icon" />
            <span class="title-text">训练模板</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text">管理和配置AI训练任务模板</span>
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
            创建模板
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
              <a-select-option value="paddle">PaddlePaddle</a-select-option>
              <a-select-option value="mindspore">MindSpore</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterCategory"
              placeholder="选择类别"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部类别</a-select-option>
              <a-select-option value="cv">计算机视觉</a-select-option>
              <a-select-option value="nlp">自然语言处理</a-select-option>
              <a-select-option value="audio">语音识别</a-select-option>
              <a-select-option value="recommendation">推荐系统</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="16" :md="8" :lg="8">
            <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索模板名称或描述"
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
          :data-source="filteredTemplates"
          :loading="loading"
          :pagination="paginationConfig"
          row-key="id"
          size="middle"
          :scroll="{ x: 'max-content' }"
          class="sci-fi-table"
        >
          <!-- 框架列 -->
          <template #framework="{ record }">
            <a-tag
              :color="getFrameworkColor(record.framework)"
              class="framework-tag"
            >
              <component
                :is="getFrameworkIcon(record.framework)"
                class="framework-icon"
              />
              {{ getFrameworkText(record.framework) }}
            </a-tag>
          </template>

          <!-- 类别列 -->
          <template #category="{ record }">
            <a-tag
              :color="getCategoryColor(record.category)"
              class="category-tag"
            >
              {{ getCategoryText(record.category) }}
            </a-tag>
          </template>

          <!-- 版本列 -->
          <template #version="{ record }">
            <span class="version-text">v{{ record.version }}</span>
          </template>

          <!-- 资源配置列 -->
          <template #defaultResources="{ record }">
            <div class="resources-info">
              <div class="resource-item">
                <DatabaseOutlined class="resource-icon" />
                <span class="resource-label">CPU:</span>
                <span class="resource-value">{{
                  record.defaultResources.cpu
                }}</span>
              </div>
              <div class="resource-item">
                <ThunderboltOutlined class="resource-icon" />
                <span class="resource-label">内存:</span>
                <span class="resource-value"
                  >{{ record.defaultResources.memory }}GB</span
                >
              </div>
              <div class="resource-item" v-if="record.defaultResources.gpu">
                <BugOutlined class="resource-icon" />
                <span class="resource-label">GPU:</span>
                <span class="resource-value">{{
                  record.defaultResources.gpu
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
                @click="useTemplate(record)"
                class="action-btn"
              >
                <PlayCircleOutlined />
                使用
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
                    <a-menu-item key="edit">
                      <EditOutlined />
                      编辑
                    </a-menu-item>
                    <a-menu-item key="duplicate">
                      <CopyOutlined />
                      复制
                    </a-menu-item>
                    <a-menu-item key="export">
                      <ExportOutlined />
                      导出
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

    <!-- 创建模板模态框 -->
    <a-modal
      v-model:open="createModalVisible"
      title="创建训练模板"
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
            <a-form-item label="模板名称" name="name">
              <a-input
                v-model:value="createForm.name"
                placeholder="请输入模板名称"
                class="form-input"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="版本" name="version">
              <a-input
                v-model:value="createForm.version"
                placeholder="请输入版本号，如：1.0.0"
                class="form-input"
              />
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
                <a-select-option value="tensorflow">TensorFlow</a-select-option>
                <a-select-option value="pytorch">PyTorch</a-select-option>
                <a-select-option value="paddle">PaddlePaddle</a-select-option>
                <a-select-option value="mindspore">MindSpore</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="应用类别" name="category">
              <a-select
                v-model:value="createForm.category"
                placeholder="选择应用类别"
                class="form-select"
              >
                <a-select-option value="cv">计算机视觉</a-select-option>
                <a-select-option value="nlp">自然语言处理</a-select-option>
                <a-select-option value="audio">语音识别</a-select-option>
                <a-select-option value="recommendation"
                  >推荐系统</a-select-option
                >
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="基础镜像" name="image">
          <a-select
            v-model:value="createForm.image"
            placeholder="选择基础镜像"
            class="form-select"
          >
            <a-select-option value="tensorflow/tensorflow:2.13.0-gpu">
              TensorFlow 2.13.0 GPU
            </a-select-option>
            <a-select-option
              value="pytorch/pytorch:2.0.1-cuda11.7-cudnn8-runtime"
            >
              PyTorch 2.0.1 CUDA 11.7
            </a-select-option>
            <a-select-option
              value="paddlepaddle/paddle:2.5.1-gpu-cuda11.7-cudnn8.4"
            >
              PaddlePaddle 2.5.1 GPU
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

        <a-divider class="form-divider">默认资源配置</a-divider>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="8">
            <a-form-item label="CPU 核数" name="cpu">
              <a-input-number
                v-model:value="createForm.cpu"
                :min="1"
                :max="64"
                :step="1"
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
                :max="16"
                style="width: 100%"
                addon-after="卡"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider class="form-divider">启动配置</a-divider>

        <a-form-item label="启动命令" name="command">
          <a-textarea
            v-model:value="createForm.command"
            placeholder="请输入启动命令，支持多行"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>

        <a-form-item label="环境变量" name="envVars">
          <div class="env-vars-container">
            <div
              v-for="(envVar, index) in createForm.envVars"
              :key="index"
              class="env-var-item"
            >
              <a-input
                v-model:value="envVar.key"
                placeholder="变量名"
                class="env-var-key"
              />
              <a-input
                v-model:value="envVar.value"
                placeholder="变量值"
                class="env-var-value"
              />
              <a-button
                type="text"
                danger
                @click="removeEnvVar(index, 'create')"
                class="env-var-remove"
              >
                <DeleteOutlined />
              </a-button>
            </div>
            <a-button
              type="dashed"
              @click="addEnvVar('create')"
              style="width: 100%"
              class="add-env-var-btn"
            >
              <PlusOutlined />
              添加环境变量
            </a-button>
          </div>
        </a-form-item>

        <a-form-item label="描述" name="description">
          <a-textarea
            v-model:value="createForm.description"
            placeholder="请输入模板描述"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 编辑模板模态框 -->
    <a-modal
      v-model:open="editModalVisible"
      title="编辑训练模板"
      width="900px"
      :confirm-loading="editLoading"
      @ok="handleEditSubmit"
      @cancel="handleEditCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="editFormRef"
        :model="editForm"
        :rules="editFormRules"
        layout="vertical"
        class="edit-form"
      >
        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="模板名称" name="name">
              <a-input
                v-model:value="editForm.name"
                placeholder="请输入模板名称"
                class="form-input"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="版本" name="version">
              <a-input
                v-model:value="editForm.version"
                placeholder="请输入版本号，如：1.0.0"
                class="form-input"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="训练框架" name="framework">
              <a-select
                v-model:value="editForm.framework"
                placeholder="选择训练框架"
                class="form-select"
              >
                <a-select-option value="tensorflow">TensorFlow</a-select-option>
                <a-select-option value="pytorch">PyTorch</a-select-option>
                <a-select-option value="paddle">PaddlePaddle</a-select-option>
                <a-select-option value="mindspore">MindSpore</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="应用类别" name="category">
              <a-select
                v-model:value="editForm.category"
                placeholder="选择应用类别"
                class="form-select"
              >
                <a-select-option value="cv">计算机视觉</a-select-option>
                <a-select-option value="nlp">自然语言处理</a-select-option>
                <a-select-option value="audio">语音识别</a-select-option>
                <a-select-option value="recommendation"
                  >推荐系统</a-select-option
                >
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="基础镜像" name="image">
          <a-select
            v-model:value="editForm.image"
            placeholder="选择基础镜像"
            class="form-select"
          >
            <a-select-option value="tensorflow/tensorflow:2.13.0-gpu">
              TensorFlow 2.13.0 GPU
            </a-select-option>
            <a-select-option
              value="pytorch/pytorch:2.0.1-cuda11.7-cudnn8-runtime"
            >
              PyTorch 2.0.1 CUDA 11.7
            </a-select-option>
            <a-select-option
              value="paddlepaddle/paddle:2.5.1-gpu-cuda11.7-cudnn8.4"
            >
              PaddlePaddle 2.5.1 GPU
            </a-select-option>
            <a-select-option value="custom">自定义镜像</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item
          v-if="editForm.image === 'custom'"
          label="自定义镜像地址"
          name="customImage"
        >
          <a-input
            v-model:value="editForm.customImage"
            placeholder="请输入镜像地址"
            class="form-input"
          />
        </a-form-item>

        <a-divider class="form-divider">默认资源配置</a-divider>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="8">
            <a-form-item label="CPU 核数" name="cpu">
              <a-input-number
                v-model:value="editForm.cpu"
                :min="1"
                :max="64"
                :step="1"
                style="width: 100%"
                addon-after="核"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="8">
            <a-form-item label="内存" name="memory">
              <a-input-number
                v-model:value="editForm.memory"
                :min="1"
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
                v-model:value="editForm.gpu"
                :min="0"
                :max="16"
                style="width: 100%"
                addon-after="卡"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider class="form-divider">启动配置</a-divider>

        <a-form-item label="启动命令" name="command">
          <a-textarea
            v-model:value="editForm.command"
            placeholder="请输入启动命令，支持多行"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>

        <a-form-item label="环境变量" name="envVars">
          <div class="env-vars-container">
            <div
              v-for="(envVar, index) in editForm.envVars"
              :key="index"
              class="env-var-item"
            >
              <a-input
                v-model:value="envVar.key"
                placeholder="变量名"
                class="env-var-key"
              />
              <a-input
                v-model:value="envVar.value"
                placeholder="变量值"
                class="env-var-value"
              />
              <a-button
                type="text"
                danger
                @click="removeEnvVar(index, 'edit')"
                class="env-var-remove"
              >
                <DeleteOutlined />
              </a-button>
            </div>
            <a-button
              type="dashed"
              @click="addEnvVar('edit')"
              style="width: 100%"
              class="add-env-var-btn"
            >
              <PlusOutlined />
              添加环境变量
            </a-button>
          </div>
        </a-form-item>

        <a-form-item label="描述" name="description">
          <a-textarea
            v-model:value="editForm.description"
            placeholder="请输入模板描述"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="模板详情"
      width="1000px"
      :footer="null"
      class="sci-fi-modal detail-modal"
    >
      <div v-if="selectedTemplate" class="detail-content">
        <a-tabs default-active-key="basic" class="detail-tabs">
          <!-- 基本信息 -->
          <a-tab-pane key="basic" tab="基本信息">
            <a-descriptions
              :column="{ xs: 1, sm: 2 }"
              bordered
              class="detail-descriptions"
            >
              <a-descriptions-item label="模板名称">
                {{ selectedTemplate.name }}
              </a-descriptions-item>
              <a-descriptions-item label="版本">
                v{{ selectedTemplate.version }}
              </a-descriptions-item>
              <a-descriptions-item label="训练框架">
                <a-tag :color="getFrameworkColor(selectedTemplate.framework)">
                  <component
                    :is="getFrameworkIcon(selectedTemplate.framework)"
                  />
                  {{ getFrameworkText(selectedTemplate.framework) }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="应用类别">
                <a-tag :color="getCategoryColor(selectedTemplate.category)">
                  {{ getCategoryText(selectedTemplate.category) }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="创建者">
                {{ selectedTemplate.creator }}
              </a-descriptions-item>
              <a-descriptions-item label="创建时间">
                {{ selectedTemplate.createTime }}
              </a-descriptions-item>
              <a-descriptions-item label="基础镜像" :span="2">
                <code class="image-code">{{ selectedTemplate.image }}</code>
              </a-descriptions-item>
              <a-descriptions-item label="描述" :span="2">
                {{ selectedTemplate.description || '暂无描述' }}
              </a-descriptions-item>
            </a-descriptions>
          </a-tab-pane>

          <!-- 资源配置 -->
          <a-tab-pane key="resources" tab="资源配置">
            <div class="resource-config">
              <a-row :gutter="24">
                <a-col :xs="24" :sm="8">
                  <div class="resource-card">
                    <div class="resource-card-header">
                      <DatabaseOutlined class="resource-card-icon" />
                      <span class="resource-card-title">CPU</span>
                    </div>
                    <div class="resource-card-value">
                      {{ selectedTemplate.defaultResources.cpu }} 核
                    </div>
                  </div>
                </a-col>
                <a-col :xs="24" :sm="8">
                  <div class="resource-card">
                    <div class="resource-card-header">
                      <ThunderboltOutlined class="resource-card-icon" />
                      <span class="resource-card-title">内存</span>
                    </div>
                    <div class="resource-card-value">
                      {{ selectedTemplate.defaultResources.memory }} GB
                    </div>
                  </div>
                </a-col>
                <a-col :xs="24" :sm="8">
                  <div class="resource-card">
                    <div class="resource-card-header">
                      <BugOutlined class="resource-card-icon" />
                      <span class="resource-card-title">GPU</span>
                    </div>
                    <div class="resource-card-value">
                      {{ selectedTemplate.defaultResources.gpu || 0 }} 卡
                    </div>
                  </div>
                </a-col>
              </a-row>
            </div>
          </a-tab-pane>

          <!-- 启动配置 -->
          <a-tab-pane key="config" tab="启动配置">
            <div class="config-section">
              <h4 class="config-title">启动命令</h4>
              <div class="command-container">
                <pre class="command-content">{{
                  selectedTemplate.command
                }}</pre>
              </div>

              <h4 class="config-title">环境变量</h4>
              <div class="env-vars-display">
                <div
                  v-for="envVar in selectedTemplate.envVars"
                  :key="envVar.key"
                  class="env-var-display-item"
                >
                  <span class="env-var-key-display">{{ envVar.key }}</span>
                  <span class="env-var-separator">=</span>
                  <span class="env-var-value-display">{{ envVar.value }}</span>
                </div>
              </div>
            </div>
          </a-tab-pane>
        </a-tabs>
      </div>
    </a-modal>

    <!-- 使用模板模态框 -->
    <a-modal
      v-model:open="useTemplateModalVisible"
      title="使用模板创建训练任务"
      width="800px"
      :confirm-loading="useTemplateLoading"
      @ok="handleUseTemplateSubmit"
      @cancel="handleUseTemplateCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="useTemplateFormRef"
        :model="useTemplateForm"
        :rules="useTemplateFormRules"
        layout="vertical"
      >
        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="任务名称" name="jobName">
              <a-input
                v-model:value="useTemplateForm.jobName"
                placeholder="请输入任务名称"
                class="form-input"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="命名空间" name="namespace">
              <a-select
                v-model:value="useTemplateForm.namespace"
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

        <a-divider class="form-divider">资源配置（可调整）</a-divider>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="8">
            <a-form-item label="CPU 核数" name="cpu">
              <a-input-number
                v-model:value="useTemplateForm.cpu"
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
                v-model:value="useTemplateForm.memory"
                :min="1"
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
                v-model:value="useTemplateForm.gpu"
                :min="0"
                :max="16"
                style="width: 100%"
                addon-after="卡"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="启动参数（可选）" name="args">
          <a-textarea
            v-model:value="useTemplateForm.args"
            placeholder="请输入额外的启动参数"
            :rows="3"
            class="form-textarea"
          />
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
  CodeOutlined,
  PlusOutlined,
  ReloadOutlined,
  DatabaseOutlined,
  ThunderboltOutlined,
  BugOutlined,
  EyeOutlined,
  MoreOutlined,
  PlayCircleOutlined,
  EditOutlined,
  CopyOutlined,
  DeleteOutlined,
  ExportOutlined,
  AppstoreOutlined,
  ExperimentOutlined,
  CloudOutlined,
  RocketOutlined,
} from '@ant-design/icons-vue';

// ===== 类型定义 =====
interface TrainingResources {
  cpu: number;
  memory: number;
  gpu?: number;
}

interface EnvVar {
  key: string;
  value: string;
}

interface TrainingTemplate {
  id: string;
  name: string;
  version: string;
  framework: 'tensorflow' | 'pytorch' | 'paddle' | 'mindspore';
  category: 'cv' | 'nlp' | 'audio' | 'recommendation';
  image: string;
  defaultResources: TrainingResources;
  command: string;
  envVars: EnvVar[];
  creator: string;
  createTime: string;
  updateTime?: string;
  description?: string;
}

interface CreateForm {
  name: string;
  version: string;
  framework: string;
  category: string;
  image: string;
  customImage: string;
  cpu: number;
  memory: number;
  gpu: number;
  command: string;
  envVars: EnvVar[];
  description: string;
}

interface EditForm {
  id: string;
  name: string;
  version: string;
  framework: string;
  category: string;
  image: string;
  customImage: string;
  cpu: number;
  memory: number;
  gpu: number;
  command: string;
  envVars: EnvVar[];
  description: string;
}

interface UseTemplateForm {
  jobName: string;
  namespace: string;
  cpu: number;
  memory: number;
  gpu: number;
  args: string;
}

// ===== 响应式数据 =====
const loading = ref<boolean>(false);
const createModalVisible = ref<boolean>(false);
const editModalVisible = ref<boolean>(false);
const detailModalVisible = ref<boolean>(false);
const useTemplateModalVisible = ref<boolean>(false);
const createLoading = ref<boolean>(false);
const editLoading = ref<boolean>(false);
const useTemplateLoading = ref<boolean>(false);

const filterFramework = ref<string>('');
const filterCategory = ref<string>('');
const searchKeyword = ref<string>('');

const selectedTemplate = ref<TrainingTemplate | null>(null);

// ===== 表单引用 =====
const createFormRef = ref<FormInstance>();
const editFormRef = ref<FormInstance>();
const useTemplateFormRef = ref<FormInstance>();

// ===== 表单数据 =====
const createForm = reactive<CreateForm>({
  name: '',
  version: '1.0.0',
  framework: '',
  category: '',
  image: '',
  customImage: '',
  cpu: 4,
  memory: 8,
  gpu: 1,
  command: 'python train.py',
  envVars: [{ key: '', value: '' }],
  description: '',
});

const editForm = reactive<EditForm>({
  id: '',
  name: '',
  version: '',
  framework: '',
  category: '',
  image: '',
  customImage: '',
  cpu: 4,
  memory: 8,
  gpu: 1,
  command: '',
  envVars: [{ key: '', value: '' }],
  description: '',
});

const useTemplateForm = reactive<UseTemplateForm>({
  jobName: '',
  namespace: 'default',
  cpu: 4,
  memory: 8,
  gpu: 1,
  args: '',
});

// ===== 配置数据 =====
const FRAMEWORK_CONFIG = {
  tensorflow: { color: 'orange', text: 'TensorFlow', icon: ExperimentOutlined },
  pytorch: { color: 'red', text: 'PyTorch', icon: RocketOutlined },
  paddle: { color: 'blue', text: 'PaddlePaddle', icon: CloudOutlined },
  mindspore: { color: 'purple', text: 'MindSpore', icon: AppstoreOutlined },
} as const;

const CATEGORY_CONFIG = {
  cv: { color: 'cyan', text: '计算机视觉' },
  nlp: { color: 'green', text: '自然语言处理' },
  audio: { color: 'magenta', text: '语音识别' },
  recommendation: { color: 'gold', text: '推荐系统' },
} as const;

// ===== 模拟数据 =====
const templates = ref<TrainingTemplate[]>([
  {
    id: 'tpl-001',
    name: 'resnet50-imagenet',
    version: '1.2.0',
    framework: 'tensorflow',
    category: 'cv',
    image: 'tensorflow/tensorflow:2.13.0-gpu',
    defaultResources: { cpu: 8, memory: 16, gpu: 2 },
    command: 'python train.py --model resnet50 --dataset imagenet --epochs 100',
    envVars: [
      { key: 'CUDA_VISIBLE_DEVICES', value: '0,1' },
      { key: 'TF_FORCE_GPU_ALLOW_GROWTH', value: 'true' },
    ],
    creator: 'admin',
    createTime: '2024-06-20 14:30:00',
    description: 'ResNet50 在 ImageNet 数据集上的训练模板',
  },
  {
    id: 'tpl-002',
    name: 'bert-base-chinese',
    version: '2.1.0',
    framework: 'pytorch',
    category: 'nlp',
    image: 'pytorch/pytorch:2.0.1-cuda11.7-cudnn8-runtime',
    defaultResources: { cpu: 16, memory: 32, gpu: 4 },
    command: 'python train_bert.py --model bert-base-chinese --max_epochs 10',
    envVars: [
      { key: 'CUDA_VISIBLE_DEVICES', value: '0,1,2,3' },
      { key: 'NCCL_DEBUG', value: 'INFO' },
    ],
    creator: 'nlp-team',
    createTime: '2024-06-21 09:15:00',
    description: 'BERT 中文预训练模型训练模板',
  },
  {
    id: 'tpl-003',
    name: 'yolov8-detection',
    version: '1.0.0',
    framework: 'pytorch',
    category: 'cv',
    image: 'ultralytics/yolov8:latest',
    defaultResources: { cpu: 12, memory: 24, gpu: 2 },
    command: 'yolo train model=yolov8n.yaml data=coco.yaml epochs=300',
    envVars: [{ key: 'CUDA_VISIBLE_DEVICES', value: '0,1' }],
    creator: 'cv-team',
    createTime: '2024-06-22 16:45:00',
    description: 'YOLOv8 目标检测训练模板',
  },
  {
    id: 'tpl-004',
    name: 'deepfm-recommendation',
    version: '1.1.0',
    framework: 'paddle',
    category: 'recommendation',
    image: 'paddlepaddle/paddle:2.5.1-gpu-cuda11.7-cudnn8.4',
    defaultResources: { cpu: 8, memory: 16, gpu: 1 },
    command: 'python deepfm_train.py --config config.yaml',
    envVars: [{ key: 'FLAGS_gpu_memory_limit_mb', value: '8192' }],
    creator: 'rec-team',
    createTime: '2024-06-23 08:20:00',
    description: 'DeepFM 推荐系统训练模板',
  },
  {
    id: 'tpl-005',
    name: 'wav2vec2-speech',
    version: '1.0.0',
    framework: 'pytorch',
    category: 'audio',
    image: 'pytorch/pytorch:2.0.1-cuda11.7-cudnn8-runtime',
    defaultResources: { cpu: 16, memory: 32, gpu: 4 },
    command: 'python train_wav2vec2.py --config-name base',
    envVars: [
      { key: 'CUDA_VISIBLE_DEVICES', value: '0,1,2,3' },
      { key: 'HYDRA_FULL_ERROR', value: '1' },
    ],
    creator: 'audio-team',
    createTime: '2024-06-22 11:30:00',
    description: 'Wav2Vec2 语音识别训练模板',
  },
]);

// ===== 表单验证规则 =====
const createFormRules = {
  name: [
    { required: true, message: '请输入模板名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
  ],
  version: [
    { required: true, message: '请输入版本号', trigger: 'blur' },
    {
      pattern: /^\d+\.\d+\.\d+$/,
      message: '版本号格式：x.y.z',
      trigger: 'blur',
    },
  ],
  framework: [{ required: true, message: '请选择训练框架', trigger: 'change' }],
  category: [{ required: true, message: '请选择应用类别', trigger: 'change' }],
  image: [{ required: true, message: '请选择基础镜像', trigger: 'change' }],
  customImage: [
    { required: true, message: '请输入自定义镜像地址', trigger: 'blur' },
  ],
  command: [{ required: true, message: '请输入启动命令', trigger: 'blur' }],
  cpu: [{ required: true, message: '请输入 CPU 核数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
};

const editFormRules = {
  name: [
    { required: true, message: '请输入模板名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
  ],
  version: [
    { required: true, message: '请输入版本号', trigger: 'blur' },
    {
      pattern: /^\d+\.\d+\.\d+$/,
      message: '版本号格式：x.y.z',
      trigger: 'blur',
    },
  ],
  framework: [{ required: true, message: '请选择训练框架', trigger: 'change' }],
  category: [{ required: true, message: '请选择应用类别', trigger: 'change' }],
  image: [{ required: true, message: '请选择基础镜像', trigger: 'change' }],
  customImage: [
    { required: true, message: '请输入自定义镜像地址', trigger: 'blur' },
  ],
  command: [{ required: true, message: '请输入启动命令', trigger: 'blur' }],
  cpu: [{ required: true, message: '请输入 CPU 核数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
};

const useTemplateFormRules = {
  jobName: [
    { required: true, message: '请输入任务名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
  ],
  namespace: [{ required: true, message: '请选择命名空间', trigger: 'change' }],
  cpu: [{ required: true, message: '请输入 CPU 核数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
};

// ===== 表格列配置 =====
const columns: TableColumnsType<TrainingTemplate> = [
  {
    title: '模板名称',
    dataIndex: 'name',
    key: 'name',
    width: 200,
    ellipsis: true,
  },
  {
    title: '版本',
    key: 'version',
    width: 100,
    slots: { customRender: 'version' },
  },
  {
    title: '训练框架',
    key: 'framework',
    width: 130,
    slots: { customRender: 'framework' },
  },
  {
    title: '应用类别',
    key: 'category',
    width: 130,
    slots: { customRender: 'category' },
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
    width: 100,
  },
  {
    title: '默认资源',
    key: 'defaultResources',
    width: 200,
    slots: { customRender: 'defaultResources' },
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
  total: computed(() => filteredTemplates.value.length),
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// ===== 计算属性 =====
const filteredTemplates = computed(() => {
  let result = templates.value;

  if (filterFramework.value) {
    result = result.filter((item) => item.framework === filterFramework.value);
  }

  if (filterCategory.value) {
    result = result.filter((item) => item.category === filterCategory.value);
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter(
      (item) =>
        item.name.toLowerCase().includes(keyword) ||
        (item.description && item.description.toLowerCase().includes(keyword)),
    );
  }

  return result;
});

// ===== 工具函数 =====
const getFrameworkColor = (framework: string): string => {
  return (
    FRAMEWORK_CONFIG[framework as keyof typeof FRAMEWORK_CONFIG]?.color ||
    'default'
  );
};

const getFrameworkIcon = (framework: string) => {
  return (
    FRAMEWORK_CONFIG[framework as keyof typeof FRAMEWORK_CONFIG]?.icon ||
    ExperimentOutlined
  );
};

const getFrameworkText = (framework: string): string => {
  return (
    FRAMEWORK_CONFIG[framework as keyof typeof FRAMEWORK_CONFIG]?.text ||
    framework
  );
};

const getCategoryColor = (category: string): string => {
  return (
    CATEGORY_CONFIG[category as keyof typeof CATEGORY_CONFIG]?.color ||
    'default'
  );
};

const getCategoryText = (category: string): string => {
  return (
    CATEGORY_CONFIG[category as keyof typeof CATEGORY_CONFIG]?.text || category
  );
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

// ===== 重置表单函数 =====
const resetCreateForm = (): void => {
  Object.assign(createForm, {
    name: '',
    version: '1.0.0',
    framework: '',
    category: '',
    image: '',
    customImage: '',
    cpu: 4,
    memory: 8,
    gpu: 1,
    command: 'python train.py',
    envVars: [{ key: '', value: '' }],
    description: '',
  });
};

const resetEditForm = (): void => {
  Object.assign(editForm, {
    id: '',
    name: '',
    version: '',
    framework: '',
    category: '',
    image: '',
    customImage: '',
    cpu: 4,
    memory: 8,
    gpu: 1,
    command: '',
    envVars: [{ key: '', value: '' }],
    description: '',
  });
};

// ===== 创建模板相关函数 =====
const showCreateModal = (): void => {
  createModalVisible.value = true;
};

const handleCreateSubmit = async (): Promise<void> => {
  try {
    await createFormRef.value?.validate();
    createLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 2000));

    const newTemplate: TrainingTemplate = {
      id: `tpl-${Date.now()}`,
      name: createForm.name,
      version: createForm.version,
      framework: createForm.framework as TrainingTemplate['framework'],
      category: createForm.category as TrainingTemplate['category'],
      image:
        createForm.image === 'custom'
          ? createForm.customImage
          : createForm.image,
      defaultResources: {
        cpu: createForm.cpu,
        memory: createForm.memory,
        ...(createForm.gpu > 0 && { gpu: createForm.gpu }),
      },
      command: createForm.command,
      envVars: createForm.envVars.filter((env) => env.key && env.value),
      creator: 'current-user',
      createTime: new Date().toLocaleString(),
      description: createForm.description,
    };

    templates.value.unshift(newTemplate);
    createModalVisible.value = false;
    message.success('训练模板创建成功');

    createFormRef.value?.resetFields();
    resetCreateForm();
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    createLoading.value = false;
  }
};

const handleCreateCancel = (): void => {
  createModalVisible.value = false;
  createFormRef.value?.resetFields();
  resetCreateForm();
};

// ===== 编辑模板相关函数 =====
const handleEdit = (record: TrainingTemplate): void => {
  // 填充编辑表单数据
  Object.assign(editForm, {
    id: record.id,
    name: record.name,
    version: record.version,
    framework: record.framework,
    category: record.category,
    image: record.image,
    customImage:
      record.image.startsWith('tensorflow/') ||
      record.image.startsWith('pytorch/') ||
      record.image.startsWith('paddlepaddle/')
        ? ''
        : record.image,
    cpu: record.defaultResources.cpu,
    memory: record.defaultResources.memory,
    gpu: record.defaultResources.gpu || 0,
    command: record.command,
    envVars:
      record.envVars.length > 0
        ? [...record.envVars]
        : [{ key: '', value: '' }],
    description: record.description || '',
  });

  // 如果是自定义镜像，设置相应的选项
  if (
    !record.image.startsWith('tensorflow/') &&
    !record.image.startsWith('pytorch/') &&
    !record.image.startsWith('paddlepaddle/')
  ) {
    editForm.image = 'custom';
    editForm.customImage = record.image;
  }

  editModalVisible.value = true;
};

const handleEditSubmit = async (): Promise<void> => {
  try {
    await editFormRef.value?.validate();
    editLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 2000));

    // 查找并更新模板
    const index = templates.value.findIndex((item) => item.id === editForm.id);
    if (index !== -1 && editForm.id) {
      const updatedTemplate: TrainingTemplate = {
        id: editForm.id,
        name: editForm.name,
        version: editForm.version,
        framework: editForm.framework as TrainingTemplate['framework'],
        category: editForm.category as TrainingTemplate['category'],
        image:
          editForm.image === 'custom' ? editForm.customImage : editForm.image,
        defaultResources: {
          cpu: editForm.cpu,
          memory: editForm.memory,
          ...(editForm.gpu > 0 && { gpu: editForm.gpu }),
        },
        command: editForm.command,
        envVars: editForm.envVars.filter((env) => env.key && env.value),
        description: editForm.description,
        updateTime: new Date().toLocaleString(),
        creator: templates.value[index]?.creator || '',
        createTime: templates.value[index]?.createTime || '',
      };

      templates.value[index] = updatedTemplate;

      // 如果当前选中的模板是被编辑的模板，更新详情显示
      if (selectedTemplate.value?.id === editForm.id) {
        selectedTemplate.value = updatedTemplate;
      }
    }

    editModalVisible.value = false;
    message.success('训练模板更新成功');

    editFormRef.value?.resetFields();
    resetEditForm();
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    editLoading.value = false;
  }
};

const handleEditCancel = (): void => {
  editModalVisible.value = false;
  editFormRef.value?.resetFields();
  resetEditForm();
};

// ===== 环境变量操作函数 =====
const addEnvVar = (formType: 'create' | 'edit'): void => {
  if (formType === 'create') {
    createForm.envVars.push({ key: '', value: '' });
  } else {
    editForm.envVars.push({ key: '', value: '' });
  }
};

const removeEnvVar = (index: number, formType: 'create' | 'edit'): void => {
  const envVars = formType === 'create' ? createForm.envVars : editForm.envVars;
  if (envVars.length > 1) {
    envVars.splice(index, 1);
  }
};

// ===== 使用模板相关函数 =====
const useTemplate = (record: TrainingTemplate): void => {
  selectedTemplate.value = record;
  useTemplateForm.jobName = `${record.name}-job`;
  useTemplateForm.cpu = record.defaultResources.cpu;
  useTemplateForm.memory = record.defaultResources.memory;
  useTemplateForm.gpu = record.defaultResources.gpu || 0;
  useTemplateModalVisible.value = true;
};

const handleUseTemplateSubmit = async (): Promise<void> => {
  try {
    await useTemplateFormRef.value?.validate();
    useTemplateLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 1500));

    useTemplateModalVisible.value = false;
    message.success('训练任务创建成功，请前往任务管理查看');

    useTemplateFormRef.value?.resetFields();
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    useTemplateLoading.value = false;
  }
};

const handleUseTemplateCancel = (): void => {
  useTemplateModalVisible.value = false;
  useTemplateFormRef.value?.resetFields();
};

// ===== 查看详情函数 =====
const viewDetails = (record: TrainingTemplate): void => {
  selectedTemplate.value = record;
  detailModalVisible.value = true;
};

// ===== 菜单操作函数 =====
const handleMenuAction = (key: string, record: TrainingTemplate): void => {
  const actions = {
    edit: () => handleEdit(record),
    duplicate: () => handleDuplicate(record),
    export: () => handleExport(record),
    delete: () => handleDelete(record),
  };

  const action = actions[key as keyof typeof actions];
  if (action) {
    action();
  }
};

const handleDuplicate = async (record: TrainingTemplate): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));

    const duplicatedTemplate: TrainingTemplate = {
      ...record,
      id: `tpl-${Date.now()}`,
      name: `${record.name}-copy`,
      version: '1.0.0',
      createTime: new Date().toLocaleString(),
      updateTime: undefined,
    };

    templates.value.unshift(duplicatedTemplate);
    message.success('模板复制成功');
  } catch (error) {
    message.error('复制失败');
  } finally {
    loading.value = false;
  }
};

const handleExport = (record: TrainingTemplate): void => {
  const templateData = JSON.stringify(record, null, 2);
  const blob = new Blob([templateData], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${record.name}-template.json`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
  message.success('模板导出成功');
};

const handleDelete = (record: TrainingTemplate): void => {
  const deleteConfirm = () => {
    const index = templates.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      templates.value.splice(index, 1);
      message.success('模板删除成功');

      // 如果删除的是当前查看详情的模板，关闭详情弹窗
      if (selectedTemplate.value?.id === record.id) {
        detailModalVisible.value = false;
        selectedTemplate.value = null;
      }
    }
  };

  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除模板 "${record.name}" 吗？此操作不可恢复。`,
      okText: '确认',
      cancelText: '取消',
      type: 'warning',
      onOk: deleteConfirm,
    });
  });
};

// ===== 其他操作函数 =====
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
/* ===== 基础样式 ===== */
.training-template-container {
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

/* ===== 框架标签 ===== */
.framework-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border-radius: 6px !important;
  font-weight: 500 !important;
  padding: 4px 8px !important;
}

.framework-icon {
  font-size: 12px;
}

/* ===== 类别标签 ===== */
.category-tag {
  border-radius: 6px !important;
  font-weight: 500 !important;
  padding: 4px 8px !important;
}

/* ===== 版本显示 ===== */
.version-text {
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
  font-weight: 600;
  color: #1890ff;
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
.create-form :deep(.ant-form-item-label > label),
.edit-form :deep(.ant-form-item-label > label) {
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

/* ===== 环境变量配置 ===== */
.env-vars-container {
  border-radius: 6px;
  padding: 16px;
  gap: 12px;
  display: flex;
  flex-direction: column;
}

.env-var-item {
  display: flex;
  gap: 8px;
  align-items: center;
}

.env-var-key,
.env-var-value {
  flex: 1;
  border-radius: 4px !important;
}

.env-var-remove {
  border-radius: 4px !important;
  color: #ff4d4f !important;
}

.add-env-var-btn {
  border-radius: 6px !important;
  border-style: dashed !important;
  transition: all 0.3s ease;
}

.add-env-var-btn:hover {
  color: #1890ff !important;
  border-color: #1890ff !important;
}

/* ===== 详情页样式 ===== */
.detail-content {
  max-height: 600px;
  overflow-y: auto;
}

.detail-tabs :deep(.ant-tabs-tab) {
  font-weight: 500 !important;
}

.image-code {
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
  font-size: 12px;
}

/* ===== 资源配置卡片 ===== */
.resource-config {
  margin-bottom: 24px;
}

.resource-card {
  border-radius: 8px !important;
  padding: 20px;
  text-align: center;
  transition: all 0.3s ease;
}

.resource-card:hover {
  transform: translateY(-2px);
}

.resource-card-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 12px;
}

.resource-card-icon {
  font-size: 20px;
  color: #1890ff;
}

.resource-card-title {
  font-weight: 600;
  font-size: 14px;
}

.resource-card-value {
  font-size: 24px;
  font-weight: 700;
  color: #1890ff;
}

/* ===== 启动配置 ===== */
.config-section {
  margin-bottom: 24px;
}

.config-title {
  font-weight: 600;
  margin-bottom: 12px;
}

.command-container {
  border-radius: 6px !important;
  padding: 16px !important;
  margin-bottom: 24px;
}

.command-content {
  margin: 0;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  word-wrap: break-word;
  white-space: pre-wrap;
}

.env-vars-display {
  border-radius: 6px !important;
  padding: 16px !important;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.env-var-display-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
  font-size: 13px;
}

.env-var-key-display {
  font-weight: 600;
  color: #1890ff;
}

.env-var-separator {
  font-weight: 600;
}

.env-var-value-display {
  font-weight: 500;
}

/* ===== 响应式设计 ===== */
@media (max-width: 768px) {
  .training-template-container {
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

  .resource-card {
    margin-bottom: 16px;
  }

  .env-var-item {
    flex-direction: column;
    align-items: stretch;
  }

  .env-var-remove {
    align-self: flex-end;
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

  .resource-card-value {
    font-size: 20px;
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
