import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:activity',
      order: 5,
      title: '监控管理',
    },
    name: 'Monitor',
    path: '/monitor',
    children: [
      {
        name: 'ClusterMonitor',
        path: 'cluster',
        component: () => import('#/views/monitor/ClusterMonitor.vue'),
        meta: {
          icon: 'lucide:network',
          title: '集群监控',
        },
      },
      {
        name: 'TaskMonitor',
        path: 'task',
        component: () => import('#/views/monitor/TaskMonitor.vue'),
        meta: {
          icon: 'lucide:check-circle',
          title: '任务监控',
        },
      },
      {
        name: 'AlertManagement',
        path: 'alert',
        component: () => import('#/views/monitor/AlertManagement.vue'),
        meta: {
          icon: 'lucide:bell',
          title: '告警管理',
        },
      },
    ],
  },
];

export default routes;
