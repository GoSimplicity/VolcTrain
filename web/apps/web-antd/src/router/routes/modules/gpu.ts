import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:cpu',
      order: 4,
      title: 'GPU管理',
    },
    name: 'GPU',
    path: '/gpu',
    children: [
      {
        name: 'GPUList',
        path: 'list',
        component: () => import('#/views/gpu/GPUList.vue'),
        meta: {
          icon: 'lucide:list',
          title: 'GPU列表',
        },
      },
      {
        name: 'GPUMonitor',
        path: 'monitor',
        component: () => import('#/views/gpu/GPUMonitor.vue'),
        meta: {
          icon: 'lucide:activity',
          title: 'GPU监控',
        },
      },
      {
        name: 'GPUSchedule',
        path: 'schedule',
        component: () => import('#/views/gpu/GPUSchedule.vue'),
        meta: {
          icon: 'lucide:calendar',
          title: 'GPU调度',
        },
      },
    ],
  },
];

export default routes;
