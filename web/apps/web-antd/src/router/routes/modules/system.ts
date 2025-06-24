import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:layout-dashboard',
      order: -1,
      title: '系统管理',
    },
    name: 'System',
    path: '/system',
    children: [
      {
        name: 'User',
        path: 'user',
        component: () => import('#/views/system/SystemUser.vue'),
        meta: {
          icon: 'carbon:workspace',
          title: '用户管理',
        },
      },
      {
        name: 'Role',
        path: 'role',
        component: () => import('#/views/system/SystemRole.vue'),
        meta: {
          icon: 'carbon:workspace',
          title: '角色管理',
        },
      },
      {
        name: 'API',
        path: 'api',
        component: () => import('#/views/system/SystemApi.vue'),
        meta: {
          icon: 'carbon:workspace',
          title: 'API管理',
        },
      },
    ],
  },
];

export default routes;
