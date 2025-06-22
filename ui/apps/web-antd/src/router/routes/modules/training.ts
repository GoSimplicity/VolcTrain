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
        name: 'TrainingQueue',
        path: 'queue',
        component: () => import('#/views/training/TrainingQueue.vue'),
        meta: {
          icon: 'carbon:align-box-bottom-center',
          title: '训练队列管理',
        },
      },
      {
        name: 'TrainingHistory',
        path: 'history',
        component: () => import('#/views/training/TrainingHistory.vue'),
        meta: {
          icon: 'carbon:document',
          title: '历史任务记录',
        },
      },
      {
        name: 'TrainingTemplate',
        path: 'template',
        component: () => import('#/views/training/TrainingTemplate.vue'),
        meta: {
          icon: 'carbon:template',
          title: '训练模板管理',
        },
      },
      {
        name: 'TrainingConfig',
        path: 'config',
        component: () => import('#/views/training/TrainingConfig.vue'),
        meta: {
          icon: 'carbon:settings',
          title: '训练配置管理',
        },
      },
    ],
  },
];

export default routes;
