<template>
  <Modal
    v-model:open="visible"
    title="成员管理"
    width="800px"
    :footer="null"
  >
    <div v-if="workspace">
      <!-- 工作空间信息 -->
      <Card size="small" style="margin-bottom: 16px">
        <div style="display: flex; justify-content: space-between; align-items: center">
          <div>
            <h4 style="margin: 0">{{ workspace.name }}</h4>
            <p style="margin: 4px 0 0 0; color: #666">
              总成员: {{ memberList.length }} 人 / 最大成员: {{ workspace.config.maxMembers }} 人
            </p>
          </div>
          <Button type="primary" @click="showInviteModal">
            <UserAddOutlined />
            邀请成员
          </Button>
        </div>
      </Card>

      <!-- 成员搜索 -->
      <div style="margin-bottom: 16px">
        <Input
          v-model:value="searchKeyword"
          placeholder="搜索成员姓名或用户名"
          style="width: 300px"
          @change="handleSearch"
        >
          <template #prefix>
            <SearchOutlined />
          </template>
        </Input>
      </div>

      <!-- 成员列表 -->
      <Table
        :columns="columns"
        :data-source="filteredMembers"
        :loading="loading"
        :pagination="false"
        row-key="id"
        size="small"
      >
        <!-- 用户信息 -->
        <template #user="{ record }">
          <div style="display: flex; align-items: center">
            <Avatar size="small" :src="record.avatar">
              {{ record.realName?.[0] || record.username?.[0] }}
            </Avatar>
            <div style="margin-left: 8px">
              <div>{{ record.realName }}</div>
              <div style="font-size: 12px; color: #999">@{{ record.username }}</div>
            </div>
          </div>
        </template>

        <!-- 权限 -->
        <template #permission="{ record }">
          <Select
            :value="record.permission"
            style="width: 120px"
            size="small"
            @change="(value) => handlePermissionChange(record, value)"
            :disabled="record.permission === 'owner'"
          >
            <Select.Option value="owner" disabled>所有者</Select.Option>
            <Select.Option value="admin">管理员</Select.Option>
            <Select.Option value="member">成员</Select.Option>
            <Select.Option value="viewer">查看者</Select.Option>
          </Select>
        </template>

        <!-- 状态 -->
        <template #status="{ record }">
          <Tag :color="record.status === 'active' ? 'success' : 'default'">
            {{ record.status === 'active' ? '活跃' : '非活跃' }}
          </Tag>
        </template>

        <!-- 加入时间 -->
        <template #joinTime="{ record }">
          <span>{{ formatDateTime(record.joinTime, 'YYYY-MM-DD') }}</span>
        </template>

        <!-- 最后活跃 -->
        <template #lastActive="{ record }">
          <span>{{ formatDateTime(record.lastActiveTime, 'MM-DD HH:mm') || '从未活跃' }}</span>
        </template>

        <!-- 操作 -->
        <template #action="{ record }">
          <Space size="small">
            <Button 
              type="link" 
              size="small" 
              @click="viewMemberDetail(record)"
            >
              详情
            </Button>
            <Button 
              type="link" 
              size="small" 
              danger
              @click="removeMember(record)"
              :disabled="record.permission === 'owner'"
            >
              移除
            </Button>
          </Space>
        </template>
      </Table>
    </div>

    <!-- 邀请成员模态框 -->
    <Modal
      v-model:open="inviteModalVisible"
      title="邀请成员"
      @ok="handleInviteSubmit"
      @cancel="inviteModalVisible = false"
      :confirm-loading="inviteLoading"
    >
      <Form
        ref="inviteFormRef"
        :model="inviteForm"
        :rules="inviteRules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
      >
        <Form.Item label="邀请用户" name="userIds">
          <Select
            v-model:value="inviteForm.userIds"
            mode="multiple"
            placeholder="选择要邀请的用户"
            style="width: 100%"
            :filter-option="false"
            :not-found-content="userSearchLoading ? '搜索中...' : '暂无数据'"
            @search="handleUserSearch"
            @change="handleUserSelect"
          >
            <Select.Option
              v-for="user in userOptions"
              :key="user.id"
              :value="user.id"
            >
              <div style="display: flex; align-items: center">
                <Avatar size="small" :src="user.avatar">
                  {{ user.realName?.[0] }}
                </Avatar>
                <span style="margin-left: 8px">{{ user.realName }} (@{{ user.username }})</span>
              </div>
            </Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label="初始权限" name="permission">
          <Select v-model:value="inviteForm.permission">
            <Select.Option value="admin">管理员</Select.Option>
            <Select.Option value="member">成员</Select.Option>
            <Select.Option value="viewer">查看者</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label="邀请消息" name="message">
          <Input.TextArea
            v-model:value="inviteForm.message"
            :rows="3"
            placeholder="可选的邀请消息"
            :maxlength="200"
            show-count
          />
        </Form.Item>
      </Form>
    </Modal>
  </Modal>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, watch } from 'vue';
import {
  Modal,
  Card,
  Button,
  Input,
  Table,
  Avatar,
  Tag,
  Select,
  Space,
  Form,
  message,
} from 'ant-design-vue';
import { UserAddOutlined, SearchOutlined } from '@ant-design/icons-vue';
import type { 
  Workspace, 
  WorkspaceMember, 
  InviteMemberRequest,
  UpdateMemberPermissionRequest,
  WorkspacePermission
} from '#/api/types';
import { 
  getWorkspaceMembers, 
  inviteMembers, 
  updateMemberPermission,
  removeMember as removeMemberApi 
} from '#/api';
import { formatDateTime } from '#/utils/date';

const props = defineProps<{
  visible: boolean;
  workspace: Workspace | null;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
}>();

const loading = ref(false);
const memberList = ref<WorkspaceMember[]>([]);
const searchKeyword = ref('');

// 邀请成员相关
const inviteModalVisible = ref(false);
const inviteLoading = ref(false);
const inviteFormRef = ref();
const userSearchLoading = ref(false);
const userOptions = ref<any[]>([]);

const inviteForm = reactive<InviteMemberRequest>({
  workspaceId: '',
  userIds: [],
  permission: 'member' as WorkspacePermission,
  message: '',
});

const inviteRules = {
  userIds: [
    { required: true, message: '请选择要邀请的用户', trigger: 'change' },
  ],
  permission: [
    { required: true, message: '请选择初始权限', trigger: 'change' },
  ],
};

// 表格列定义
const columns = [
  {
    title: '用户',
    key: 'user',
    slots: { customRender: 'user' },
    width: 200,
  },
  {
    title: '邮箱',
    dataIndex: 'email',
    key: 'email',
    width: 180,
  },
  {
    title: '权限',
    key: 'permission',
    slots: { customRender: 'permission' },
    width: 130,
  },
  {
    title: '状态',
    key: 'status',
    slots: { customRender: 'status' },
    width: 80,
  },
  {
    title: '加入时间',
    key: 'joinTime',
    slots: { customRender: 'joinTime' },
    width: 100,
  },
  {
    title: '最后活跃',
    key: 'lastActive',
    slots: { customRender: 'lastActive' },
    width: 100,
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 100,
  },
];

// 过滤后的成员列表
const filteredMembers = computed(() => {
  if (!searchKeyword.value) {
    return memberList.value;
  }
  const keyword = searchKeyword.value.toLowerCase();
  return memberList.value.filter(member => 
    member.realName.toLowerCase().includes(keyword) ||
    member.username.toLowerCase().includes(keyword)
  );
});

// 监听workspace变化
watch([() => props.visible, () => props.workspace], ([visible, workspace]) => {
  if (visible && workspace) {
    loadMembers(workspace.id);
    inviteForm.workspaceId = workspace.id;
  }
});

// 加载成员列表
const loadMembers = async (workspaceId: string) => {
  try {
    loading.value = true;
    const response = await getWorkspaceMembers(workspaceId);
    memberList.value = response.data;
  } catch (error) {
    message.error('加载成员列表失败');
  } finally {
    loading.value = false;
  }
};

// 搜索处理
const handleSearch = () => {
  // 搜索逻辑已在computed中处理
};

// 权限变更
const handlePermissionChange = async (member: WorkspaceMember, permission: WorkspacePermission) => {
  try {
    const request: UpdateMemberPermissionRequest = {
      workspaceId: props.workspace!.id,
      userId: member.userId,
      permission,
    };
    await updateMemberPermission(request);
    message.success('权限更新成功');
    
    // 更新本地数据
    const index = memberList.value.findIndex(m => m.id === member.id);
    if (index !== -1) {
      memberList.value[index].permission = permission;
    }
  } catch (error) {
    message.error('权限更新失败');
  }
};

// 移除成员
const removeMember = async (member: WorkspaceMember) => {
  try {
    await removeMemberApi(props.workspace!.id, member.userId);
    message.success('成员移除成功');
    
    // 从本地列表中移除
    const index = memberList.value.findIndex(m => m.id === member.id);
    if (index !== -1) {
      memberList.value.splice(index, 1);
    }
  } catch (error) {
    message.error('移除成员失败');
  }
};

// 查看成员详情
const viewMemberDetail = (member: WorkspaceMember) => {
  console.log('查看成员详情:', member);
  // 实现成员详情查看逻辑
};

// 显示邀请模态框
const showInviteModal = () => {
  inviteModalVisible.value = true;
  // 重置表单
  inviteForm.userIds = [];
  inviteForm.permission = 'member' as WorkspacePermission;
  inviteForm.message = '';
  userOptions.value = [];
};

// 用户搜索
const handleUserSearch = (keyword: string) => {
  if (!keyword) {
    userOptions.value = [];
    return;
  }
  
  // 模拟用户搜索
  userSearchLoading.value = true;
  setTimeout(() => {
    userOptions.value = [
      { id: '1', username: 'user1', realName: '用户一', avatar: '' },
      { id: '2', username: 'user2', realName: '用户二', avatar: '' },
    ].filter(user => 
      user.realName.includes(keyword) || user.username.includes(keyword)
    );
    userSearchLoading.value = false;
  }, 500);
};

// 用户选择
const handleUserSelect = (userIds: string[]) => {
  console.log('选择的用户:', userIds);
};

// 提交邀请
const handleInviteSubmit = async () => {
  try {
    await inviteFormRef.value.validateFields();
    inviteLoading.value = true;

    await inviteMembers(inviteForm);
    message.success('邀请发送成功');
    inviteModalVisible.value = false;
    
    // 重新加载成员列表
    if (props.workspace) {
      loadMembers(props.workspace.id);
    }
  } catch (error: any) {
    if (error.errorFields) {
      message.error('请完善邀请信息');
    } else {
      message.error('邀请失败：' + (error.message || '未知错误'));
    }
  } finally {
    inviteLoading.value = false;
  }
};
</script>

<style scoped>
.ant-table-small {
  font-size: 13px;
}
</style>