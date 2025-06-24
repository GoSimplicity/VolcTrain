import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:folder',
      order: 0,
      title: '工作空间管理',
    },
    name: 'Workspace',
    path: '/workspace',
    children: [
      {
        name: 'WorkspaceList',
        path: 'list',
        component: () => import('#/views/workspace/WorkspaceList.vue'),
        meta: {
          icon: 'lucide:list',
          title: '工作空间列表',
        },
      },
      {
        name: 'WorkspaceCreate',
        path: 'create',
        component: () => import('#/views/workspace/WorkspaceCreate.vue'),
        meta: {
          icon: 'lucide:plus-circle',
          title: '创建工作空间',
        },
      },
      {
        name: 'WorkspaceProject',
        path: 'project',
        component: () => import('#/views/workspace/WorkspaceProject.vue'),
        meta: {
          icon: 'lucide:folder-kanban',
          title: '项目管理',
        },
      },
    ],
  },
];

export default routes;
