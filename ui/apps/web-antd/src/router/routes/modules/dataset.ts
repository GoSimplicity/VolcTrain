import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'ic:baseline-dataset',
      order: 3,
      title: '数据集管理',
    },
    name: 'Dataset',
    path: '/dataset',
    children: [
      {
        name: 'DatasetList',
        path: 'list',
        component: () => import('#/views/dataset/DatasetList.vue'),
        meta: {
          title: '数据集列表',
        },
      },
      {
        name: 'DatasetUpload',
        path: 'upload',
        component: () => import('#/views/dataset/DatasetUpload.vue'),
        meta: {
          title: '数据集上传',
        },
      },
      {
        name: 'DatasetProcess',
        path: 'process',
        component: () => import('#/views/dataset/DatasetProcess.vue'),
        meta: {
          title: '数据集处理',
        },
      },
    ],
  },
];

export default routes;
