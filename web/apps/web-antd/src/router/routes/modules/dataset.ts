import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'ic:baseline-dataset',
      order: 3,
      title: '数据集处理',
    },
    name: 'Dataset',
    path: '/dataset',
    children: [
      {
        name: 'DatasetManagement',
        path: 'management',
        component: () => import('#/views/dataset/DatasetManagement.vue'),
        meta: {
          icon: 'carbon:data-view',
          title: '数据集管理',
        },
      },
      {
        name: 'DatasetAnnotation',
        path: 'annotation',
        component: () => import('#/views/dataset/DatasetAnnotation.vue'),
        meta: {
          icon: 'carbon:tag',
          title: '数据标注',
        },
      },
      {
        name: 'DatasetAnalytics',
        path: 'analytics',
        component: () => import('#/views/dataset/DatasetAnalytics.vue'),
        meta: {
          icon: 'carbon:chart-line',
          title: '数据分析',
        },
      },
      {
        name: 'DatasetVersion',
        path: 'version',
        component: () => import('#/views/dataset/DatasetVersion.vue'),
        meta: {
          icon: 'carbon:version',
          title: '版本管理',
        },
      },
    ],
  },
];

export default routes;
