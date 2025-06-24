import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'mingcute:cube-3d-line',
      order: 4,
      title: '模型管理',
    },
    name: 'Model',
    path: '/model',
    children: [
      {
        name: 'ModelList',
        path: 'list',
        component: () => import('#/views/model/ModelList.vue'),
        meta: {
          title: '模型列表',
        },
      },
      {
        name: 'ModelVersion',
        path: 'version',
        component: () => import('#/views/model/ModelVersion.vue'),
        meta: {
          title: '版本管理',
        },
      },
      {
        name: 'ModelRegistry',
        path: 'registry',
        component: () => import('#/views/model/ModelRegistry.vue'),
        meta: {
          title: '模型仓库',
        },
      },
    ],
  },
];

export default routes;
