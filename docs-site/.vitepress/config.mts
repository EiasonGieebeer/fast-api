import { defineConfig } from 'vitepress'

export default defineConfig({
  lang: 'zh-CN',
  title: 'Fast API 文档',
  description: 'Fast API 接入与使用文档',
  cleanUrls: true,
  lastUpdated: true,
  head: [
    [
      'link',
      {
        rel: 'icon',
        href: 'https://www.jetbrains.top/fast-api/logo.png',
      },
    ],
  ],
  themeConfig: {
    logo: 'https://www.jetbrains.top/fast-api/logo.png',
    siteTitle: 'Fast API 文档',
    nav: [
      { text: '快速开始', link: '/guide/getting-started' },
      { text: 'API 调用', link: '/guide/api' },
      { text: '控制台', link: 'https://www.jetbrains.top/dashboard' },
      { text: '模型广场', link: 'https://www.jetbrains.top/pricing' },
    ],
    sidebar: [
      {
        text: '开始使用',
        items: [
          { text: '文档首页', link: '/' },
          { text: '快速开始', link: '/guide/getting-started' },
          { text: 'API 调用', link: '/guide/api' },
          { text: '客户端配置', link: '/guide/clients' },
        ],
      },
      {
        text: '账户与帮助',
        items: [
          { text: '额度与计费', link: '/guide/billing' },
          { text: '常见问题', link: '/guide/faq' },
          { text: '服务说明', link: '/legal/terms' },
        ],
      },
    ],
    outline: {
      level: [2, 3],
      label: '本页目录',
    },
    search: {
      provider: 'local',
    },
    socialLinks: [
      {
        icon: 'github',
        link: 'https://github.com/EiasonGieebeer/fast-api',
      },
    ],
    editLink: {
      pattern:
        'https://github.com/EiasonGieebeer/fast-api/edit/my-custom/docs-site/:path',
      text: '在 GitHub 上编辑此页',
    },
    lastUpdated: {
      text: '最后更新',
      formatOptions: {
        dateStyle: 'medium',
        timeStyle: 'short',
      },
    },
    docFooter: {
      prev: '上一篇',
      next: '下一篇',
    },
    returnToTopLabel: '返回顶部',
    sidebarMenuLabel: '菜单',
    darkModeSwitchLabel: '主题',
    lightModeSwitchTitle: '切换到浅色模式',
    darkModeSwitchTitle: '切换到深色模式',
    footer: {
      message: '请妥善保管 API 密钥，并遵守适用法律法规。',
      copyright: '© 2026 Fast API',
    },
  },
})
