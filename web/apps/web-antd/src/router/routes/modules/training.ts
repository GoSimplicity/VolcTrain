import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:cpu',
      order: 1,
      title: '训练任务管理',
    },
    name: 'Training',
    path: '/training',
    children: [
      {
        name: 'TrainingDashboard',
        path: 'dashboard',
        component: () => import('#/views/training/TrainingDashboard.vue'),
        meta: {
          icon: 'carbon:dashboard',
          title: '训练概览',
        },
      },
      {
        name: 'TrainingQueue',
        path: 'queue',
        component: () => import('#/views/training/TrainingQueue.vue'),
        meta: {
          icon: 'carbon:task',
          title: '任务队列',
        },
      },
      {
        name: 'TrainingHistory',
        path: 'history',
        component: () => import('#/views/training/TrainingHistory.vue'),
        meta: {
          icon: 'carbon:task',
          title: '历史记录',
        },
      },
      {
        name: 'ResourceAllocation',
        path: 'allocation',
        component: () => import('#/views/training/ResourceAllocation.vue'),
        meta: {
          icon: 'carbon:cloud-service-management',
          title: '资源分配',
        },
      },
      {
        name: 'TrainingTemplate',
        path: 'template',
        component: () => import('#/views/training/TrainingTemplate.vue'),
        meta: {
          icon: 'carbon:template',
          title: '训练模板',
        },
      },
      {
        name: 'NotebookServices',
        path: 'notebooks',
        component: () => import('#/views/training/NotebookServices.vue'),
        meta: {
          icon: 'carbon:tools',
          title: 'Notebook服务',
        },
      },
      {
        name: 'ExperimentTracking',
        path: 'experiments',
        component: () => import('#/views/training/ExperimentTracking.vue'),
        meta: {
          icon: 'carbon:tools',
          title: '实验跟踪',
        },
      },
    ],
  },
];

export default routes;
