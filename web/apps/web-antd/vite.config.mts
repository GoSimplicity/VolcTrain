import { defineConfig } from '@vben/vite-config';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      server: {
        proxy: {
          '/api': {
            changeOrigin: true,
            // 代理到后端API服务器
            target: 'http://localhost:8888',
            ws: true,
          },
        },
      },
    },
  };
});
