<template>
  <div class="annotation-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <TagsOutlined class="title-icon" />
            <span class="title-text">数据集标注</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text">管理和标注您的数据集</span>
          </p>
        </div>
        <div class="action-section">
          <a-space size="large">
            <a-button
              type="primary"
              size="large"
              @click="showCreateModal"
              class="create-btn"
            >
              <PlusOutlined />
              创建标注任务
            </a-button>
            <a-button size="large" @click="showUploadModal" class="upload-btn">
              <UploadOutlined />
              上传数据集
            </a-button>
          </a-space>
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
              <a-select-option value="pending">待标注</a-select-option>
              <a-select-option value="annotating">标注中</a-select-option>
              <a-select-option value="reviewing">审核中</a-select-option>
              <a-select-option value="completed">已完成</a-select-option>
              <a-select-option value="failed">失败</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterType"
              placeholder="选择类型"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部类型</a-select-option>
              <a-select-option value="image">图像分类</a-select-option>
              <a-select-option value="object-detection"
                >目标检测</a-select-option
              >
              <a-select-option value="text">文本标注</a-select-option>
              <a-select-option value="audio">音频标注</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="16" :md="8" :lg="8">
            <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索标注任务名称或创建者"
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
          :data-source="filteredAnnotations"
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

          <!-- 类型列 -->
          <template #type="{ record }">
            <a-tag :color="getTypeColor(record.type)" class="type-tag">
              <component :is="getTypeIcon(record.type)" class="type-icon" />
              {{ getTypeText(record.type) }}
            </a-tag>
          </template>

          <!-- 进度列 -->
          <template #progress="{ record }">
            <div class="progress-wrapper">
              <a-progress
                :percent="record.progress"
                size="small"
                :show-info="false"
                class="progress-bar"
              />
              <span class="progress-text">{{ record.progress }}%</span>
              <div class="progress-detail">
                {{ record.annotatedCount }}/{{ record.totalCount }}
              </div>
            </div>
          </template>

          <!-- 标注者列 -->
          <template #annotators="{ record }">
            <a-avatar-group :max-count="3" :size="24" class="annotator-group">
              <a-avatar
                v-for="annotator in record.annotators"
                :key="annotator.id"
                :title="annotator.name"
                class="annotator-avatar"
              >
                {{ annotator.name.charAt(0) }}
              </a-avatar>
            </a-avatar-group>
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
                @click="startAnnotation(record)"
                :disabled="record.status === 'completed'"
                class="action-btn"
              >
                <EditOutlined />
                标注
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
                      key="assign"
                      :disabled="record.status === 'completed'"
                    >
                      <UserAddOutlined />
                      分配标注者
                    </a-menu-item>
                    <a-menu-item key="export" :disabled="record.progress === 0">
                      <ExportOutlined />
                      导出标注
                    </a-menu-item>
                    <a-menu-item key="clone">
                      <CopyOutlined />
                      克隆任务
                    </a-menu-item>
                    <a-menu-item
                      key="pause"
                      :disabled="record.status !== 'annotating'"
                    >
                      <PauseCircleOutlined />
                      暂停
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

    <!-- 创建标注任务模态框 -->
    <a-modal
      v-model:open="createModalVisible"
      title="创建标注任务"
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
            <a-form-item label="标注类型" name="type">
              <a-select
                v-model:value="createForm.type"
                placeholder="选择标注类型"
                class="form-select"
                @change="handleTypeChange"
              >
                <a-select-option value="image">图像分类</a-select-option>
                <a-select-option value="object-detection"
                  >目标检测</a-select-option
                >
                <a-select-option value="text">文本标注</a-select-option>
                <a-select-option value="audio">音频标注</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="数据集" name="dataset">
          <a-select
            v-model:value="createForm.dataset"
            placeholder="选择数据集"
            class="form-select"
            show-search
            :filter-option="filterDatasetOption"
          >
            <a-select-option
              v-for="dataset in availableDatasets"
              :key="dataset.id"
              :value="dataset.id"
            >
              <div class="dataset-option">
                <span class="dataset-name">{{ dataset.name }}</span>
                <span class="dataset-info">{{ dataset.fileCount }} 文件</span>
              </div>
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="标注说明" name="description">
          <a-textarea
            v-model:value="createForm.description"
            placeholder="请输入标注说明和要求"
            :rows="4"
            class="form-textarea"
          />
        </a-form-item>

        <!-- 标签配置 -->
        <a-divider class="form-divider">标签配置</a-divider>

        <a-form-item label="标签列表" name="labels">
          <div class="labels-config">
            <div
              v-for="(label, index) in createForm.labels"
              :key="index"
              class="label-item"
            >
              <a-input
                v-model:value="label.name"
                placeholder="标签名称"
                class="label-input"
              />
              <a-input
                v-model:value="label.color"
                type="color"
                class="color-picker"
              />
              <a-button
                type="text"
                danger
                @click="removeLabel(index)"
                class="remove-label-btn"
              >
                <DeleteOutlined />
              </a-button>
            </div>
            <a-button
              type="dashed"
              @click="addLabel"
              class="add-label-btn"
              block
            >
              <PlusOutlined />
              添加标签
            </a-button>
          </div>
        </a-form-item>

        <!-- 分配标注者 -->
        <a-divider class="form-divider">分配标注者</a-divider>

        <a-form-item label="标注者" name="annotators">
          <a-select
            v-model:value="createForm.annotators"
            placeholder="选择标注者"
            mode="multiple"
            class="form-select"
            show-search
            :filter-option="filterAnnotatorOption"
          >
            <a-select-option
              v-for="user in availableUsers"
              :key="user.id"
              :value="user.id"
            >
              <div class="annotator-option">
                <a-avatar :size="20" class="annotator-option-avatar">
                  {{ user.name.charAt(0) }}
                </a-avatar>
                <span class="annotator-name">{{ user.name }}</span>
                <span class="annotator-role">{{ user.role }}</span>
              </div>
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="截止日期" name="deadline">
              <a-date-picker
                v-model:value="createForm.deadline"
                placeholder="选择截止日期"
                style="width: 100%"
                class="form-date-picker"
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
                <a-select-option value="low">低</a-select-option>
                <a-select-option value="medium">中</a-select-option>
                <a-select-option value="high">高</a-select-option>
                <a-select-option value="urgent">紧急</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <!-- 质量控制 -->
        <a-divider class="form-divider">质量控制</a-divider>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="审核方式" name="reviewMode">
              <a-radio-group v-model:value="createForm.reviewMode">
                <a-radio value="none">无需审核</a-radio>
                <a-radio value="manual">人工审核</a-radio>
                <a-radio value="auto">自动审核</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="一致性要求" name="consistency">
              <a-input-number
                v-model:value="createForm.consistency"
                :min="0"
                :max="100"
                :step="5"
                addon-after="%"
                style="width: 100%"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <!-- 上传数据集模态框 -->
    <a-modal
      v-model:open="uploadModalVisible"
      title="上传数据集"
      width="700px"
      :confirm-loading="uploadLoading"
      @ok="handleUploadSubmit"
      @cancel="handleUploadCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="uploadFormRef"
        :model="uploadForm"
        :rules="uploadFormRules"
        layout="vertical"
        class="upload-form"
      >
        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="数据集名称" name="name">
              <a-input
                v-model:value="uploadForm.name"
                placeholder="请输入数据集名称"
                class="form-input"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="数据类型" name="dataType">
              <a-select
                v-model:value="uploadForm.dataType"
                placeholder="选择数据类型"
                class="form-select"
              >
                <a-select-option value="image">图像</a-select-option>
                <a-select-option value="text">文本</a-select-option>
                <a-select-option value="audio">音频</a-select-option>
                <a-select-option value="video">视频</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="描述" name="description">
          <a-textarea
            v-model:value="uploadForm.description"
            placeholder="请输入数据集描述"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>

        <a-form-item label="上传文件" name="files">
          <a-upload-dragger
            v-model:file-list="uploadForm.fileList"
            name="files"
            multiple
            :before-upload="beforeUpload"
            @change="handleUploadChange"
            class="upload-dragger"
          >
            <p class="ant-upload-drag-icon">
              <InboxOutlined />
            </p>
            <p class="ant-upload-text">点击或拖拽文件到此区域上传</p>
            <p class="ant-upload-hint">
              支持单个或批量上传。支持 jpg、png、txt、mp3、mp4 等格式
            </p>
          </a-upload-dragger>
        </a-form-item>

        <a-form-item label="标签配置">
          <a-radio-group v-model:value="uploadForm.labelConfig">
            <a-radio value="manual">手动创建标签</a-radio>
            <a-radio value="import">从文件导入标签</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item
          v-if="uploadForm.labelConfig === 'import'"
          label="标签文件"
          name="labelFile"
        >
          <a-upload
            v-model:file-list="uploadForm.labelFileList"
            name="labelFile"
            :before-upload="beforeLabelUpload"
            @change="handleLabelUploadChange"
            :max-count="1"
          >
            <a-button>
              <UploadOutlined />
              选择标签文件
            </a-button>
          </a-upload>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="标注任务详情"
      width="1200px"
      :footer="null"
      class="sci-fi-modal detail-modal"
    >
      <div v-if="selectedAnnotation" class="detail-content">
        <a-row :gutter="24">
          <a-col :xs="24" :lg="12">
            <a-descriptions
              :column="1"
              bordered
              class="detail-descriptions"
              title="基本信息"
            >
              <a-descriptions-item label="任务名称">
                {{ selectedAnnotation.name }}
              </a-descriptions-item>
              <a-descriptions-item label="状态">
                <a-tag
                  :color="getStatusColor(selectedAnnotation.status)"
                  class="status-tag"
                >
                  <component :is="getStatusIcon(selectedAnnotation.status)" />
                  {{ getStatusText(selectedAnnotation.status) }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="类型">
                <a-tag
                  :color="getTypeColor(selectedAnnotation.type)"
                  class="type-tag"
                >
                  <component :is="getTypeIcon(selectedAnnotation.type)" />
                  {{ getTypeText(selectedAnnotation.type) }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="创建者">
                {{ selectedAnnotation.creator }}
              </a-descriptions-item>
              <a-descriptions-item label="创建时间">
                {{ selectedAnnotation.createTime }}
              </a-descriptions-item>
              <a-descriptions-item label="截止日期">
                {{ selectedAnnotation.deadline || '无' }}
              </a-descriptions-item>
              <a-descriptions-item label="优先级">
                <a-tag :color="getPriorityColor(selectedAnnotation.priority)">
                  {{ getPriorityText(selectedAnnotation.priority) }}
                </a-tag>
              </a-descriptions-item>
            </a-descriptions>
          </a-col>
          <a-col :xs="24" :lg="12">
            <div class="progress-section">
              <h4>标注进度</h4>
              <div class="progress-stats">
                <a-progress
                  :percent="selectedAnnotation.progress"
                  type="circle"
                  :width="120"
                  class="main-progress"
                />
                <div class="stats-grid">
                  <div class="stat-item">
                    <div class="stat-value">
                      {{ selectedAnnotation.totalCount }}
                    </div>
                    <div class="stat-label">总数据量</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">
                      {{ selectedAnnotation.annotatedCount }}
                    </div>
                    <div class="stat-label">已标注</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">
                      {{
                        selectedAnnotation.totalCount -
                        selectedAnnotation.annotatedCount
                      }}
                    </div>
                    <div class="stat-label">待标注</div>
                  </div>
                </div>
              </div>
            </div>

            <div class="annotators-section">
              <h4>标注团队</h4>
              <div class="annotators-list">
                <div
                  v-for="annotator in selectedAnnotation.annotators"
                  :key="annotator.id"
                  class="annotator-card"
                >
                  <a-avatar class="annotator-avatar">
                    {{ annotator.name.charAt(0) }}
                  </a-avatar>
                  <div class="annotator-info">
                    <div class="annotator-name">{{ annotator.name }}</div>
                    <div class="annotator-role">{{ annotator.role }}</div>
                  </div>
                  <div class="annotator-progress">
                    <a-progress
                      :percent="annotator.progress || 0"
                      size="small"
                      :show-info="false"
                    />
                    <span class="progress-text"
                      >{{ annotator.progress || 0 }}%</span
                    >
                  </div>
                </div>
              </div>
            </div>
          </a-col>
        </a-row>

        <a-divider />

        <div class="labels-section">
          <h4>标签配置</h4>
          <div class="labels-grid">
            <div
              v-for="label in selectedAnnotation.labels"
              :key="label.id"
              class="label-card"
            >
              <div
                class="label-color"
                :style="{ backgroundColor: label.color }"
              ></div>
              <span class="label-name">{{ label.name }}</span>
              <span class="label-count">{{ label.count || 0 }} 个</span>
            </div>
          </div>
        </div>

        <a-divider />

        <div class="description-section">
          <h4>任务描述</h4>
          <div class="description-content">
            {{ selectedAnnotation.description || '暂无描述' }}
          </div>
        </div>
      </div>
    </a-modal>

    <!-- 分配标注者模态框 -->
    <a-modal
      v-model:open="assignModalVisible"
      title="分配标注者"
      width="600px"
      :confirm-loading="assignLoading"
      @ok="handleAssignSubmit"
      @cancel="handleAssignCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="assignFormRef"
        :model="assignForm"
        :rules="assignFormRules"
        layout="vertical"
      >
        <a-form-item label="选择标注者" name="annotators">
          <a-select
            v-model:value="assignForm.annotators"
            placeholder="选择标注者"
            mode="multiple"
            class="form-select"
            show-search
            :filter-option="filterAnnotatorOption"
          >
            <a-select-option
              v-for="user in availableUsers"
              :key="user.id"
              :value="user.id"
            >
              <div class="annotator-option">
                <a-avatar :size="20" class="annotator-option-avatar">
                  {{ user.name.charAt(0) }}
                </a-avatar>
                <span class="annotator-name">{{ user.name }}</span>
                <span class="annotator-role">{{ user.role }}</span>
              </div>
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="分配说明" name="note">
          <a-textarea
            v-model:value="assignForm.note"
            placeholder="请输入分配说明"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 标注界面模态框 -->
    <a-modal
      v-model:open="annotationModalVisible"
      title="数据标注"
      width="90%"
      :mask-closable="false"
      :footer="null"
      class="sci-fi-modal annotation-modal"
    >
      <div v-if="currentAnnotationTask" class="annotation-workspace">
        <div class="annotation-header">
          <div class="task-info">
            <h3>{{ currentAnnotationTask.name }}</h3>
            <a-tag :color="getTypeColor(currentAnnotationTask.type)">
              {{ getTypeText(currentAnnotationTask.type) }}
            </a-tag>
          </div>
          <div class="annotation-progress">
            <span
              >进度: {{ currentAnnotationIndex + 1 }} /
              {{ annotationSamples.length }}</span
            >
            <a-progress
              :percent="
                Math.round(
                  ((currentAnnotationIndex + 1) / annotationSamples.length) *
                    100,
                )
              "
              :show-info="false"
              class="header-progress"
            />
          </div>
        </div>

        <div class="annotation-content">
          <div class="annotation-main">
            <!-- 图像标注界面 -->
            <div
              v-if="currentAnnotationTask.type === 'image'"
              class="image-annotation"
            >
              <div class="image-container">
                <img
                  :src="currentSample?.url"
                  :alt="currentSample?.name"
                  class="annotation-image"
                />
                <!-- 这里可以添加标注工具，如边界框、多边形等 -->
              </div>
            </div>

            <!-- 文本标注界面 -->
            <div
              v-else-if="currentAnnotationTask.type === 'text'"
              class="text-annotation"
            >
              <div class="text-container">
                <div class="text-content">
                  {{ currentSample?.content }}
                </div>
              </div>
            </div>
          </div>

          <div class="annotation-sidebar">
            <div class="labels-panel">
              <h4>标签</h4>
              <div class="labels-list">
                <div
                  v-for="label in currentAnnotationTask.labels"
                  :key="label.id"
                  class="label-option"
                  :class="{ active: selectedLabels.includes(label.id) }"
                  @click="toggleLabel(label.id)"
                >
                  <div
                    class="label-color-dot"
                    :style="{ backgroundColor: label.color }"
                  ></div>
                  <span class="label-name">{{ label.name }}</span>
                </div>
              </div>
            </div>

            <div class="annotation-actions">
              <a-space direction="vertical" style="width: 100%">
                <a-button
                  type="primary"
                  @click="saveAnnotation"
                  :disabled="selectedLabels.length === 0"
                  block
                >
                  保存标注
                </a-button>
                <a-button @click="skipAnnotation" block> 跳过 </a-button>
                <a-button
                  @click="previousSample"
                  :disabled="currentAnnotationIndex === 0"
                  block
                >
                  上一个
                </a-button>
                <a-button
                  @click="nextSample"
                  :disabled="
                    currentAnnotationIndex === annotationSamples.length - 1
                  "
                  block
                >
                  下一个
                </a-button>
              </a-space>
            </div>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import type {
  FormInstance,
  TableColumnsType,
  UploadFile,
} from 'ant-design-vue';
import type { Dayjs } from 'dayjs';
import {
  TagsOutlined,
  PlusOutlined,
  UploadOutlined,
  ReloadOutlined,
  EditOutlined,
  EyeOutlined,
  MoreOutlined,
  UserAddOutlined,
  ExportOutlined,
  CopyOutlined,
  PauseCircleOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  CloseCircleOutlined,
  FileImageOutlined,
  FileTextOutlined,
  SoundOutlined,
  InboxOutlined,
} from '@ant-design/icons-vue';

// ===== 类型定义 =====
interface Label {
  id: string;
  name: string;
  color: string;
  count?: number;
}

interface Annotator {
  id: string;
  name: string;
  role: string;
  progress?: number;
}

interface AnnotationTask {
  id: string;
  name: string;
  type: 'image' | 'object-detection' | 'text' | 'audio';
  status: 'pending' | 'annotating' | 'reviewing' | 'completed' | 'failed';
  creator: string;
  createTime: string;
  deadline?: string;
  priority: 'low' | 'medium' | 'high' | 'urgent';
  progress: number;
  totalCount: number;
  annotatedCount: number;
  annotators: Annotator[];
  labels: Label[];
  description?: string;
  dataset: string;
  reviewMode: 'none' | 'manual' | 'auto';
  consistency: number;
}

interface Dataset {
  id: string;
  name: string;
  fileCount: number;
  dataType: string;
}

interface User {
  id: string;
  name: string;
  role: string;
}

interface CreateForm {
  name: string;
  type: string;
  dataset: string;
  description: string;
  labels: Label[];
  annotators: string[];
  deadline?: Dayjs;
  priority: string;
  reviewMode: string;
  consistency: number;
}

interface UploadForm {
  name: string;
  dataType: string;
  description: string;
  fileList: UploadFile[];
  labelConfig: string;
  labelFileList: UploadFile[];
}

interface AssignForm {
  annotators: string[];
  note: string;
}

interface AnnotationSample {
  id: string;
  name: string;
  url?: string;
  content?: string;
  labels: string[];
}

// ===== 响应式数据 =====
const loading = ref<boolean>(false);
const createModalVisible = ref<boolean>(false);
const uploadModalVisible = ref<boolean>(false);
const detailModalVisible = ref<boolean>(false);
const assignModalVisible = ref<boolean>(false);
const annotationModalVisible = ref<boolean>(false);
const createLoading = ref<boolean>(false);
const uploadLoading = ref<boolean>(false);
const assignLoading = ref<boolean>(false);

const filterStatus = ref<string>('');
const filterType = ref<string>('');
const searchKeyword = ref<string>('');

const selectedAnnotation = ref<AnnotationTask | null>(null);
const currentAnnotationTask = ref<AnnotationTask | null>(null);
const currentAnnotationIndex = ref<number>(0);
const selectedLabels = ref<string[]>([]);

// ===== 表单引用 =====
const createFormRef = ref<FormInstance>();
const uploadFormRef = ref<FormInstance>();
const assignFormRef = ref<FormInstance>();

// ===== 表单数据 =====
const createForm = reactive<CreateForm>({
  name: '',
  type: 'image',
  dataset: '',
  description: '',
  labels: [
    { id: '1', name: '正样本', color: '#52c41a' },
    { id: '2', name: '负样本', color: '#ff4d4f' },
  ],
  annotators: [],
  deadline: undefined,
  priority: 'medium',
  reviewMode: 'manual',
  consistency: 85,
});

const uploadForm = reactive<UploadForm>({
  name: '',
  dataType: 'image',
  description: '',
  fileList: [],
  labelConfig: 'manual',
  labelFileList: [],
});

const assignForm = reactive<AssignForm>({
  annotators: [],
  note: '',
});

// ===== 配置数据 =====
const STATUS_CONFIG = {
  pending: { color: 'default', text: '待标注', icon: ClockCircleOutlined },
  annotating: { color: 'processing', text: '标注中', icon: EditOutlined },
  reviewing: { color: 'warning', text: '审核中', icon: EyeOutlined },
  completed: { color: 'success', text: '已完成', icon: CheckCircleOutlined },
  failed: { color: 'error', text: '失败', icon: CloseCircleOutlined },
} as const;

const TYPE_CONFIG = {
  image: { color: 'blue', text: '图像分类', icon: FileImageOutlined },
  'object-detection': {
    color: 'purple',
    text: '目标检测',
    icon: FileImageOutlined,
  },
  text: { color: 'green', text: '文本标注', icon: FileTextOutlined },
  audio: { color: 'orange', text: '音频标注', icon: SoundOutlined },
} as const;

const PRIORITY_CONFIG = {
  low: { color: 'default', text: '低' },
  medium: { color: 'processing', text: '中' },
  high: { color: 'warning', text: '高' },
  urgent: { color: 'error', text: '紧急' },
} as const;

// ===== 模拟数据 =====
const annotations = ref<AnnotationTask[]>([
  {
    id: 'ann-001',
    name: '图像分类-动物识别',
    type: 'image',
    status: 'annotating',
    creator: 'admin',
    createTime: '2024-06-23 09:00:00',
    deadline: '2024-07-15 18:00:00',
    priority: 'high',
    progress: 65,
    totalCount: 1000,
    annotatedCount: 650,
    annotators: [
      { id: '1', name: '张三', role: '标注员', progress: 70 },
      { id: '2', name: '李四', role: '标注员', progress: 60 },
    ],
    labels: [
      { id: '1', name: '猫', color: '#ff4d4f', count: 250 },
      { id: '2', name: '狗', color: '#52c41a', count: 200 },
      { id: '3', name: '鸟', color: '#1890ff', count: 200 },
    ],
    description: '对动物图像进行分类标注，区分猫、狗、鸟三个类别',
    dataset: 'dataset-001',
    reviewMode: 'manual',
    consistency: 85,
  },
  {
    id: 'ann-002',
    name: '目标检测-车辆识别',
    type: 'object-detection',
    status: 'reviewing',
    creator: 'manager',
    createTime: '2024-06-22 14:30:00',
    deadline: '2024-07-20 18:00:00',
    priority: 'medium',
    progress: 90,
    totalCount: 800,
    annotatedCount: 720,
    annotators: [
      { id: '3', name: '王五', role: '高级标注员', progress: 95 },
      { id: '4', name: '赵六', role: '标注员', progress: 85 },
    ],
    labels: [
      { id: '4', name: '汽车', color: '#ff4d4f', count: 400 },
      { id: '5', name: '卡车', color: '#52c41a', count: 200 },
      { id: '6', name: '摩托车', color: '#1890ff', count: 120 },
    ],
    description: '在道路图像中标注各种车辆的边界框',
    dataset: 'dataset-002',
    reviewMode: 'auto',
    consistency: 90,
  },
  {
    id: 'ann-003',
    name: '文本情感分析',
    type: 'text',
    status: 'completed',
    creator: 'researcher',
    createTime: '2024-06-20 10:15:00',
    deadline: '2024-07-10 18:00:00',
    priority: 'low',
    progress: 100,
    totalCount: 2000,
    annotatedCount: 2000,
    annotators: [
      { id: '5', name: '孙七', role: '标注员', progress: 100 },
      { id: '6', name: '周八', role: '标注员', progress: 100 },
    ],
    labels: [
      { id: '7', name: '正面', color: '#52c41a', count: 800 },
      { id: '8', name: '负面', color: '#ff4d4f', count: 600 },
      { id: '9', name: '中性', color: '#faad14', count: 600 },
    ],
    description: '对用户评论进行情感极性标注',
    dataset: 'dataset-003',
    reviewMode: 'manual',
    consistency: 92,
  },
  {
    id: 'ann-004',
    name: '音频分类-语音识别',
    type: 'audio',
    status: 'pending',
    creator: 'developer',
    createTime: '2024-06-23 16:20:00',
    priority: 'urgent',
    progress: 0,
    totalCount: 500,
    annotatedCount: 0,
    annotators: [],
    labels: [
      { id: '10', name: '男性', color: '#1890ff', count: 0 },
      { id: '11', name: '女性', color: '#eb2f96', count: 0 },
      { id: '12', name: '儿童', color: '#faad14', count: 0 },
    ],
    description: '对语音数据进行性别和年龄分类',
    dataset: 'dataset-004',
    reviewMode: 'none',
    consistency: 80,
  },
]);

const availableDatasets = ref<Dataset[]>([
  {
    id: 'dataset-001',
    name: '动物图像数据集',
    fileCount: 1000,
    dataType: 'image',
  },
  {
    id: 'dataset-002',
    name: '道路车辆数据集',
    fileCount: 800,
    dataType: 'image',
  },
  {
    id: 'dataset-003',
    name: '用户评论数据集',
    fileCount: 2000,
    dataType: 'text',
  },
  {
    id: 'dataset-004',
    name: '语音样本数据集',
    fileCount: 500,
    dataType: 'audio',
  },
]);

const availableUsers = ref<User[]>([
  { id: '1', name: '张三', role: '标注员' },
  { id: '2', name: '李四', role: '标注员' },
  { id: '3', name: '王五', role: '高级标注员' },
  { id: '4', name: '赵六', role: '标注员' },
  { id: '5', name: '孙七', role: '标注员' },
  { id: '6', name: '周八', role: '标注员' },
  { id: '7', name: '吴九', role: '审核员' },
  { id: '8', name: '郑十', role: '管理员' },
]);

const annotationSamples = ref<AnnotationSample[]>([
  {
    id: 'sample-1',
    name: 'cat_001.jpg',
    url: 'https://via.placeholder.com/400x300/1890ff/ffffff?text=Cat+Image',
    labels: [],
  },
  {
    id: 'sample-2',
    name: 'dog_001.jpg',
    url: 'https://via.placeholder.com/400x300/52c41a/ffffff?text=Dog+Image',
    labels: [],
  },
  {
    id: 'sample-3',
    name: 'bird_001.jpg',
    url: 'https://via.placeholder.com/400x300/faad14/ffffff?text=Bird+Image',
    labels: [],
  },
]);

// ===== 表单验证规则 =====
const createFormRules = {
  name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度在 2 到 50 个字符', trigger: 'blur' },
  ],
  type: [{ required: true, message: '请选择标注类型', trigger: 'change' }],
  dataset: [{ required: true, message: '请选择数据集', trigger: 'change' }],
  description: [{ required: true, message: '请输入标注说明', trigger: 'blur' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }],
  reviewMode: [
    { required: true, message: '请选择审核方式', trigger: 'change' },
  ],
  consistency: [
    { required: true, message: '请输入一致性要求', trigger: 'blur' },
  ],
};

const uploadFormRules = {
  name: [
    { required: true, message: '请输入数据集名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度在 2 到 50 个字符', trigger: 'blur' },
  ],
  dataType: [{ required: true, message: '请选择数据类型', trigger: 'change' }],
  description: [
    { required: true, message: '请输入数据集描述', trigger: 'blur' },
  ],
};

const assignFormRules = {
  annotators: [{ required: true, message: '请选择标注者', trigger: 'change' }],
};

// ===== 表格列配置 =====
const columns: TableColumnsType<AnnotationTask> = [
  {
    title: '任务名称',
    dataIndex: 'name',
    key: 'name',
    width: 200,
    ellipsis: true,
  },
  {
    title: '类型',
    key: 'type',
    width: 120,
    slots: { customRender: 'type' },
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
    slots: { customRender: 'status' },
  },
  {
    title: '进度',
    key: 'progress',
    width: 150,
    slots: { customRender: 'progress' },
  },
  {
    title: '标注者',
    key: 'annotators',
    width: 150,
    slots: { customRender: 'annotators' },
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
    width: 200,
    fixed: 'right',
    slots: { customRender: 'action' },
  },
];

// ===== 分页配置 =====
const paginationConfig = {
  total: computed(() => filteredAnnotations.value.length),
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// ===== 计算属性 =====
const filteredAnnotations = computed(() => {
  let result = annotations.value;

  if (filterStatus.value) {
    result = result.filter((item) => item.status === filterStatus.value);
  }

  if (filterType.value) {
    result = result.filter((item) => item.type === filterType.value);
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

const currentSample = computed(() => {
  if (
    currentAnnotationIndex.value >= 0 &&
    currentAnnotationIndex.value < annotationSamples.value.length
  ) {
    return annotationSamples.value[currentAnnotationIndex.value];
  }
  return null;
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
    TYPE_CONFIG[type as keyof typeof TYPE_CONFIG]?.icon || FileImageOutlined
  );
};

const getTypeText = (type: string): string => {
  return TYPE_CONFIG[type as keyof typeof TYPE_CONFIG]?.text || type;
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
};

const showUploadModal = (): void => {
  uploadModalVisible.value = true;
};

const handleCreateSubmit = async (): Promise<void> => {
  try {
    await createFormRef.value?.validate();
    createLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 2000));

    const newAnnotation: AnnotationTask = {
      id: `ann-${Date.now()}`,
      name: createForm.name,
      type: createForm.type as any,
      status: 'pending',
      creator: 'current-user',
      createTime: new Date().toLocaleString(),
      deadline: createForm.deadline?.format('YYYY-MM-DD HH:mm:ss'),
      priority: createForm.priority as any,
      progress: 0,
      totalCount: 100,
      annotatedCount: 0,
      annotators: availableUsers.value.filter((user) =>
        createForm.annotators.includes(user.id),
      ),
      labels: [...createForm.labels],
      description: createForm.description,
      dataset: createForm.dataset,
      reviewMode: createForm.reviewMode as any,
      consistency: createForm.consistency,
    };

    annotations.value.unshift(newAnnotation);
    createModalVisible.value = false;
    message.success('标注任务创建成功');

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

const handleUploadSubmit = async (): Promise<void> => {
  try {
    await uploadFormRef.value?.validate();
    uploadLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 3000));

    const newDataset: Dataset = {
      id: `dataset-${Date.now()}`,
      name: uploadForm.name,
      fileCount: uploadForm.fileList.length,
      dataType: uploadForm.dataType,
    };

    availableDatasets.value.unshift(newDataset);
    uploadModalVisible.value = false;
    message.success('数据集上传成功');

    uploadFormRef.value?.resetFields();
    uploadForm.fileList = [];
    uploadForm.labelFileList = [];
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    uploadLoading.value = false;
  }
};

const handleUploadCancel = (): void => {
  uploadModalVisible.value = false;
  uploadFormRef.value?.resetFields();
  uploadForm.fileList = [];
  uploadForm.labelFileList = [];
};

const startAnnotation = (record: AnnotationTask): void => {
  currentAnnotationTask.value = record;
  currentAnnotationIndex.value = 0;
  selectedLabels.value = [];
  annotationModalVisible.value = true;
};

const viewDetails = (record: AnnotationTask): void => {
  selectedAnnotation.value = record;
  detailModalVisible.value = true;
};

const handleMenuAction = (key: string, record: AnnotationTask): void => {
  const actions = {
    assign: () => handleAssign(record),
    export: () => handleExport(record),
    clone: () => handleClone(record),
    pause: () => handlePause(record),
    delete: () => handleDelete(record),
  };

  const action = actions[key as keyof typeof actions];
  if (action) {
    action();
  }
};

const handleAssign = (record: AnnotationTask): void => {
  selectedAnnotation.value = record;
  assignForm.annotators = record.annotators.map((a) => a.id);
  assignModalVisible.value = true;
};

const handleExport = async (_: AnnotationTask): Promise<void> => {
  message.loading('导出中...', 0);
  try {
    await new Promise((resolve) => setTimeout(resolve, 2000));
    message.destroy();
    message.success('标注数据导出成功');
  } catch (error) {
    message.destroy();
    message.error('导出失败');
  }
};

const handleClone = async (record: AnnotationTask): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1500));

    const clonedAnnotation: AnnotationTask = {
      ...record,
      id: `ann-${Date.now()}`,
      name: `${record.name}-副本`,
      status: 'pending',
      createTime: new Date().toLocaleString(),
      progress: 0,
      annotatedCount: 0,
    };

    annotations.value.unshift(clonedAnnotation);
    message.success('任务克隆成功');
  } catch (error) {
    message.error('克隆失败');
  } finally {
    loading.value = false;
  }
};

const handlePause = async (record: AnnotationTask): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const index = annotations.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      annotations.value[index]!.status = 'pending';
    }
    message.success('任务已暂停');
  } catch (error) {
    message.error('暂停失败');
  } finally {
    loading.value = false;
  }
};

const handleDelete = (record: AnnotationTask): void => {
  const deleteConfirm = () => {
    const index = annotations.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      annotations.value.splice(index, 1);
      message.success('标注任务删除成功');
    }
  };

  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除标注任务 "${record.name}" 吗？此操作不可恢复。`,
      okText: '确认',
      cancelText: '取消',
      type: 'warning',
      onOk: deleteConfirm,
    });
  });
};

const handleAssignSubmit = async (): Promise<void> => {
  try {
    await assignFormRef.value?.validate();
    assignLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 1500));

    if (selectedAnnotation.value) {
      const index = annotations.value.findIndex(
        (item) => item.id === selectedAnnotation.value!.id,
      );
      if (index !== -1) {
        annotations.value[index]!.annotators = availableUsers.value.filter(
          (user) => assignForm.annotators.includes(user.id),
        );
      }
      message.success('标注者分配成功');
      assignModalVisible.value = false;
    }
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    assignLoading.value = false;
  }
};

const handleAssignCancel = (): void => {
  assignModalVisible.value = false;
  assignFormRef.value?.resetFields();
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

const handleTypeChange = (): void => {
  // 类型变化时重置标签
  createForm.labels = [
    { id: '1', name: '正样本', color: '#52c41a' },
    { id: '2', name: '负样本', color: '#ff4d4f' },
  ];
};

const addLabel = (): void => {
  const newLabel: Label = {
    id: Date.now().toString(),
    name: '',
    color: '#1890ff',
  };
  createForm.labels.push(newLabel);
};

const removeLabel = (index: number): void => {
  createForm.labels.splice(index, 1);
};

const filterDatasetOption = (input: string, option: any): boolean => {
  return option.children.props.children[0].children
    .toLowerCase()
    .includes(input.toLowerCase());
};

const filterAnnotatorOption = (input: string, option: any): boolean => {
  const annotatorName = option.children.props.children[1].children;
  return annotatorName.toLowerCase().includes(input.toLowerCase());
};

const beforeUpload = (): boolean => {
  return false; // 阻止自动上传
};

const beforeLabelUpload = (): boolean => {
  return false; // 阻止自动上传
};

const handleUploadChange = (info: any): void => {
  uploadForm.fileList = info.fileList;
};

const handleLabelUploadChange = (info: any): void => {
  uploadForm.labelFileList = info.fileList;
};

const toggleLabel = (labelId: string): void => {
  const index = selectedLabels.value.indexOf(labelId);
  if (index > -1) {
    selectedLabels.value.splice(index, 1);
  } else {
    selectedLabels.value.push(labelId);
  }
};

const saveAnnotation = async (): Promise<void> => {
  if (currentSample.value && selectedLabels.value.length > 0) {
    currentSample.value.labels = [...selectedLabels.value];

    // 模拟保存
    await new Promise((resolve) => setTimeout(resolve, 500));
    message.success('标注已保存');

    // 自动跳到下一个
    nextSample();
  }
};

const skipAnnotation = (): void => {
  nextSample();
};

const previousSample = (): void => {
  if (currentAnnotationIndex.value > 0) {
    currentAnnotationIndex.value--;
    selectedLabels.value = currentSample.value?.labels || [];
  }
};

const nextSample = (): void => {
  if (currentAnnotationIndex.value < annotationSamples.value.length - 1) {
    currentAnnotationIndex.value++;
    selectedLabels.value = currentSample.value?.labels || [];
  } else {
    message.success('所有样本标注完成！');
    annotationModalVisible.value = false;
  }
};

// ===== 生命周期 =====
onMounted(() => {
  refreshData();
});
</script>

<style scoped>
/* ===== 基础样式 ===== */
.annotation-container {
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
.create-btn,
.upload-btn {
  height: 40px !important;
  padding: 0 24px !important;
  font-size: 14px !important;
  font-weight: 500 !important;
  border-radius: 6px !important;
  transition: all 0.3s ease !important;
}

.create-btn {
  background: #1890ff !important;
  border: none !important;
}

.upload-btn {
  background: transparent !important;
  border: 1px solid #d9d9d9 !important;
}

.create-btn:hover,
.upload-btn:hover {
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

.indicator-pending {
  background: #8c8c8c;
}

.indicator-annotating {
  background: #1890ff;
}

.indicator-reviewing {
  background: #faad14;
}

.indicator-completed {
  background: #52c41a;
}

.indicator-failed {
  background: #ff4d4f;
}

/* ===== 进度样式 ===== */
.progress-wrapper {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.progress-bar {
  width: 100px;
}

.progress-text {
  font-size: 12px;
  font-weight: 600;
}

.progress-detail {
  font-size: 11px;
  color: #8c8c8c;
}

/* ===== 标注者样式 ===== */
.annotator-group {
  display: flex;
  gap: 4px;
}

.annotator-avatar {
  background: #1890ff !important;
  font-size: 12px !important;
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
.create-form :deep(.ant-form-item-label > label),
.upload-form :deep(.ant-form-item-label > label) {
  font-weight: 500 !important;
}

.form-input,
.form-select,
.form-textarea,
.form-input-number,
.form-date-picker {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

.form-divider {
  font-weight: 500 !important;
}

/* ===== 标签配置样式 ===== */
.labels-config {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.label-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.label-input {
  flex: 1;
}

.color-picker {
  width: 60px !important;
  height: 32px !important;
  padding: 0 !important;
  border-radius: 6px !important;
}

.remove-label-btn {
  color: #ff4d4f !important;
}

.add-label-btn {
  border-radius: 6px !important;
  border-style: dashed !important;
  border-color: #d9d9d9 !important;
  color: #8c8c8c !important;
  transition: all 0.3s ease;
}

.add-label-btn:hover {
  border-color: #1890ff !important;
  color: #1890ff !important;
}

/* ===== 选项样式 ===== */
.dataset-option,
.annotator-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dataset-name,
.annotator-name {
  font-weight: 500;
}

.dataset-info,
.annotator-role {
  font-size: 12px;
  color: #8c8c8c;
  margin-left: auto;
}

.annotator-option-avatar {
  background: #1890ff !important;
}

/* ===== 上传样式 ===== */
.upload-dragger {
  border-radius: 8px !important;
}

.upload-dragger :deep(.ant-upload-btn) {
  border-radius: 8px !important;
}

/* ===== 详情页样式 ===== */
.detail-content {
  max-height: 700px;
  overflow-y: auto;
}

.detail-descriptions {
  border-radius: 6px !important;
}

.progress-section,
.annotators-section,
.labels-section,
.description-section {
  margin-bottom: 24px;
}

.progress-section h4,
.annotators-section h4,
.labels-section h4,
.description-section h4 {
  font-weight: 600;
  margin-bottom: 16px;
  color: #262626;
}

.progress-stats {
  display: flex;
  align-items: center;
  gap: 24px;
}

.main-progress {
  flex-shrink: 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  flex: 1;
}

.stat-item {
  text-align: center;
  padding: 16px;
  border-radius: 8px;
  background: #fafafa;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1890ff;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #8c8c8c;
}

.annotators-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.annotator-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  background: #fafafa;
}

.annotator-info {
  flex: 1;
}

.annotator-name {
  font-weight: 500;
  margin-bottom: 2px;
}

.annotator-role {
  font-size: 12px;
  color: #8c8c8c;
}

.annotator-progress {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
  min-width: 80px;
}

.labels-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 12px;
}

.label-card {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  background: #fafafa;
  border: 1px solid #f0f0f0;
}

.label-color {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.label-name {
  font-weight: 500;
  flex: 1;
}

.label-count {
  font-size: 12px;
  color: #8c8c8c;
}

.description-content {
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
  line-height: 1.6;
}

/* ===== 标注界面样式 ===== */
.annotation-modal :deep(.ant-modal-body) {
  padding: 0 !important;
}

.annotation-workspace {
  height: 80vh;
  display: flex;
  flex-direction: column;
}

.annotation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.task-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.task-info h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.annotation-progress {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-progress {
  width: 200px;
}

.annotation-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.annotation-main {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: #fafafa;
}

.image-annotation {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-container {
  max-width: 100%;
  max-height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.annotation-image {
  max-width: 100%;
  max-height: 100%;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.text-annotation {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.text-container {
  width: 100%;
  max-width: 800px;
  padding: 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.text-content {
  font-size: 16px;
  line-height: 1.8;
  color: #262626;
}

.annotation-sidebar {
  width: 300px;
  border-left: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
}

.labels-panel {
  flex: 1;
  padding: 24px;
}

.labels-panel h4 {
  font-weight: 600;
  margin-bottom: 16px;
}

.labels-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.label-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid #f0f0f0;
  cursor: pointer;
  transition: all 0.3s ease;
}

.label-option:hover {
  border-color: #1890ff;
  background: #f6ffed;
}

.label-option.active {
  border-color: #1890ff;
  background: #e6f7ff;
}

.label-color-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.annotation-actions {
  padding: 24px;
  border-top: 1px solid #f0f0f0;
}

/* ===== 响应式设计 ===== */
@media (max-width: 768px) {
  .annotation-container {
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
  .upload-btn,
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

  .annotation-content {
    flex-direction: column;
  }

  .annotation-sidebar {
    width: 100%;
    border-left: none;
    border-top: 1px solid #f0f0f0;
  }

  .stats-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }

  .progress-stats {
    flex-direction: column;
    gap: 16px;
  }

  .labels-grid {
    grid-template-columns: 1fr;
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
