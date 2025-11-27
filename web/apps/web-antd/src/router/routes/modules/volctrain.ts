/**
 * VolcTrain 火山引擎训练平台路由配置
 * 重新设计的现代化路由结构
 */
import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  // 仪表板
  {
    name: 'Dashboard',
    path: '/dashboard',
    component: () => import('#/views/dashboard/index.vue'),
    meta: {
      affixTab: true,
      icon: 'lucide:monitor-dashboard',
      title: $t('page.dashboard.title'),
      order: 1,
    },
  },

  // 训练管理
  {
    name: 'Training',
    path: '/training',
    meta: {
      icon: 'lucide:brain-circuit',
      title: $t('page.training.title'),
      order: 2,
    },
    children: [
      {
        name: 'TrainingJobs',
        path: 'jobs',
        component: () => import('#/views/training/jobs/index.vue'),
        meta: {
          icon: 'lucide:play-circle',
          title: $t('page.training.jobs'),
        },
      },
      {
        name: 'TrainingQueues',
        path: 'queues',
        component: () => import('#/views/training/queues/index.vue'),
        meta: {
          icon: 'lucide:list-ordered',
          title: $t('page.training.queues'),
        },
      },
      {
        name: 'TrainingExperiments',
        path: 'experiments',
        component: () => import('#/views/training/experiments/index.vue'),
        meta: {
          icon: 'lucide:flask-conical',
          title: $t('page.training.experiments'),
        },
      },
      {
        name: 'TrainingTemplates',
        path: 'templates',
        component: () => import('#/views/training/templates/index.vue'),
        meta: {
          icon: 'lucide:file-text',
          title: $t('page.training.templates'),
        },
      },
    ],
  },

  // GPU 资源管理
  {
    name: 'GPU',
    path: '/gpu',
    meta: {
      icon: 'lucide:cpu',
      title: $t('page.gpu.title'),
      order: 3,
    },
    children: [
      {
        name: 'GPUClusters',
        path: 'clusters',
        component: () => import('#/views/gpu/clusters/index.vue'),
        meta: {
          icon: 'lucide:server',
          title: $t('page.gpu.clusters'),
        },
      },
      {
        name: 'GPUNodes',
        path: 'nodes',
        component: () => import('#/views/gpu/nodes/index.vue'),
        meta: {
          icon: 'lucide:hard-drive',
          title: $t('page.gpu.nodes'),
        },
      },
      {
        name: 'GPUDevices',
        path: 'devices',
        component: () => import('#/views/gpu/devices/index.vue'),
        meta: {
          icon: 'lucide:microchip',
          title: $t('page.gpu.devices'),
        },
      },
      {
        name: 'GPUUsage',
        path: 'usage',
        component: () => import('#/views/gpu/usage/index.vue'),
        meta: {
          icon: 'lucide:activity',
          title: $t('page.gpu.usage'),
        },
      },
    ],
  },

  // 数据管理
  {
    name: 'Data',
    path: '/data',
    meta: {
      icon: 'lucide:database',
      title: $t('page.data.title'),
      order: 4,
    },
    children: [
      {
        name: 'Datasets',
        path: 'datasets',
        component: () => import('#/views/data/datasets/index.vue'),
        meta: {
          icon: 'lucide:folder-open',
          title: $t('page.data.datasets'),
        },
      },
      {
        name: 'DatasetDetail',
        path: 'datasets/:id',
        component: () => import('#/views/data/datasets/detail.vue'),
        meta: {
          hideInMenu: true,
          title: $t('page.data.datasetDetail'),
        },
      },
      {
        name: 'DataUpload',
        path: 'upload',
        component: () => import('#/views/data/upload/index.vue'),
        meta: {
          icon: 'lucide:cloud-upload',
          title: $t('page.data.upload'),
        },
      },
    ],
  },

  // 模型管理
  {
    name: 'Models',
    path: '/models',
    meta: {
      icon: 'lucide:box',
      title: $t('page.models.title'),
      order: 5,
    },
    children: [
      {
        name: 'ModelRepository',
        path: 'repository',
        component: () => import('#/views/models/repository/index.vue'),
        meta: {
          icon: 'lucide:package',
          title: $t('page.models.repository'),
        },
      },
      {
        name: 'ModelDetail',
        path: 'repository/:id',
        component: () => import('#/views/models/repository/detail.vue'),
        meta: {
          hideInMenu: true,
          title: $t('page.models.detail'),
        },
      },
      {
        name: 'ModelDeployment',
        path: 'deployment',
        component: () => import('#/views/models/deployment/index.vue'),
        meta: {
          icon: 'lucide:rocket',
          title: $t('page.models.deployment'),
        },
      },
    ],
  },

  // 监控中心
  {
    name: 'Monitoring',
    path: '/monitoring',
    meta: {
      icon: 'lucide:bar-chart-3',
      title: $t('page.monitoring.title'),
      order: 6,
    },
    children: [
      {
        name: 'SystemMonitor',
        path: 'system',
        component: () => import('#/views/monitoring/system/index.vue'),
        meta: {
          icon: 'lucide:monitor',
          title: $t('page.monitoring.system'),
        },
      },
      {
        name: 'AlertManagement',
        path: 'alerts',
        component: () => import('#/views/monitoring/alerts/index.vue'),
        meta: {
          icon: 'lucide:bell',
          title: $t('page.monitoring.alerts'),
        },
      },
      {
        name: 'PerformanceAnalysis',
        path: 'performance',
        component: () => import('#/views/monitoring/performance/index.vue'),
        meta: {
          icon: 'lucide:trending-up',
          title: $t('page.monitoring.performance'),
        },
      },
    ],
  },

  // 工作空间
  {
    name: 'Workspace',
    path: '/workspace',
    meta: {
      icon: 'lucide:users',
      title: $t('page.workspace.title'),
      order: 7,
    },
    children: [
      {
        name: 'WorkspaceOverview',
        path: 'overview',
        component: () => import('#/views/workspace/overview/index.vue'),
        meta: {
          icon: 'lucide:layout-dashboard',
          title: $t('page.workspace.overview'),
        },
      },
      {
        name: 'WorkspaceMembers',
        path: 'members',
        component: () => import('#/views/workspace/members/index.vue'),
        meta: {
          icon: 'lucide:user-plus',
          title: $t('page.workspace.members'),
        },
      },
      {
        name: 'WorkspaceSettings',
        path: 'settings',
        component: () => import('#/views/workspace/settings/index.vue'),
        meta: {
          icon: 'lucide:settings',
          title: $t('page.workspace.settings'),
        },
      },
    ],
  },

  // 系统管理（管理员权限）
  {
    name: 'System',
    path: '/system',
    meta: {
      icon: 'lucide:shield-check',
      title: $t('page.system.title'),
      authority: ['admin'],
      order: 8,
    },
    children: [
      {
        name: 'UserManagement',
        path: 'users',
        component: () => import('#/views/system/users/index.vue'),
        meta: {
          icon: 'lucide:users',
          title: $t('page.system.users'),
          authority: ['admin'],
        },
      },
      {
        name: 'RoleManagement',
        path: 'roles',
        component: () => import('#/views/system/roles/index.vue'),
        meta: {
          icon: 'lucide:shield',
          title: $t('page.system.roles'),
          authority: ['admin'],
        },
      },
      {
        name: 'SystemSettings',
        path: 'settings',
        component: () => import('#/views/system/settings/index.vue'),
        meta: {
          icon: 'lucide:cog',
          title: $t('page.system.settings'),
          authority: ['admin'],
        },
      },
      {
        name: 'SystemLogs',
        path: 'logs',
        component: () => import('#/views/system/logs/index.vue'),
        meta: {
          icon: 'lucide:file-text',
          title: $t('page.system.logs'),
          authority: ['admin'],
        },
      },
    ],
  },

  // 个人中心
  {
    name: 'Profile',
    path: '/profile',
    meta: {
      icon: 'lucide:user',
      title: $t('page.profile.title'),
      order: 9,
    },
    children: [
      {
        name: 'ProfileSettings',
        path: 'settings',
        component: () => import('#/views/profile/settings/index.vue'),
        meta: {
          icon: 'lucide:user-cog',
          title: $t('page.profile.settings'),
        },
      },
      {
        name: 'ProfileSecurity',
        path: 'security',
        component: () => import('#/views/profile/security/index.vue'),
        meta: {
          icon: 'lucide:key',
          title: $t('page.profile.security'),
        },
      },
    ],
  },
];

export default routes;