import { defineConfig, type DefaultTheme } from 'vitepress'

const logo = 'https://www.jetbrains.top/fast-api/logo.png'
const repository = 'https://github.com/EiasonGieebeer/fast-api'

const sharedTheme: Pick<DefaultTheme.Config, 'logo' | 'socialLinks'> = {
  logo,
  socialLinks: [{ icon: 'github', link: repository }],
}

const zhTheme: DefaultTheme.Config = {
  ...sharedTheme,
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
  outline: { level: [2, 3], label: '本页目录' },
  editLink: {
    pattern: `${repository}/edit/my-custom/docs-site/:path`,
    text: '在 GitHub 上编辑此页',
  },
  lastUpdated: {
    text: '最后更新',
    formatOptions: { dateStyle: 'medium', timeStyle: 'short' },
  },
  docFooter: { prev: '上一篇', next: '下一篇' },
  returnToTopLabel: '返回顶部',
  sidebarMenuLabel: '菜单',
  darkModeSwitchLabel: '主题',
  lightModeSwitchTitle: '切换到浅色模式',
  darkModeSwitchTitle: '切换到深色模式',
  langMenuLabel: '切换语言',
  skipToContentLabel: '跳到正文',
  footer: {
    message: '请妥善保管 API 密钥，并遵守适用法律法规。',
    copyright: '© 2026 Fast API',
  },
}

const enTheme: DefaultTheme.Config = {
  ...sharedTheme,
  siteTitle: 'Fast API Docs',
  nav: [
    { text: 'Quick Start', link: '/en/guide/getting-started' },
    { text: 'API Usage', link: '/en/guide/api' },
    { text: 'Console', link: 'https://www.jetbrains.top/dashboard' },
    { text: 'Models', link: 'https://www.jetbrains.top/pricing' },
  ],
  sidebar: [
    {
      text: 'Getting Started',
      items: [
        { text: 'Documentation', link: '/en/' },
        { text: 'Quick Start', link: '/en/guide/getting-started' },
        { text: 'API Usage', link: '/en/guide/api' },
        { text: 'Client Setup', link: '/en/guide/clients' },
      ],
    },
    {
      text: 'Account and Help',
      items: [
        { text: 'Billing and Quota', link: '/en/guide/billing' },
        { text: 'FAQ', link: '/en/guide/faq' },
        { text: 'Service Terms', link: '/en/legal/terms' },
      ],
    },
  ],
  outline: { level: [2, 3], label: 'On this page' },
  editLink: {
    pattern: `${repository}/edit/my-custom/docs-site/:path`,
    text: 'Edit this page on GitHub',
  },
  lastUpdated: {
    text: 'Last updated',
    formatOptions: { dateStyle: 'medium', timeStyle: 'short' },
  },
  docFooter: { prev: 'Previous', next: 'Next' },
  returnToTopLabel: 'Return to top',
  sidebarMenuLabel: 'Menu',
  darkModeSwitchLabel: 'Theme',
  lightModeSwitchTitle: 'Switch to light theme',
  darkModeSwitchTitle: 'Switch to dark theme',
  langMenuLabel: 'Change language',
  skipToContentLabel: 'Skip to content',
  footer: {
    message: 'Keep your API keys secure and comply with applicable laws.',
    copyright: '© 2026 Fast API',
  },
}

export default defineConfig({
  title: 'Fast API Docs',
  description: 'Fast API integration and usage documentation',
  cleanUrls: true,
  lastUpdated: true,
  head: [['link', { rel: 'icon', href: logo }]],
  themeConfig: {
    search: {
      provider: 'local',
      options: {
        locales: {
          root: {
            translations: {
              button: {
                buttonText: '搜索',
                buttonAriaLabel: '搜索文档',
              },
              modal: {
                displayDetails: '显示详细列表',
                resetButtonTitle: '清除搜索',
                backButtonTitle: '关闭搜索',
                noResultsText: '没有找到相关结果',
                footer: {
                  selectText: '选择',
                  selectKeyAriaLabel: '回车',
                  navigateText: '切换',
                  navigateUpKeyAriaLabel: '向上',
                  navigateDownKeyAriaLabel: '向下',
                  closeText: '关闭',
                  closeKeyAriaLabel: 'Esc',
                },
              },
            },
          },
          en: {
            translations: {
              button: {
                buttonText: 'Search',
                buttonAriaLabel: 'Search documentation',
              },
              modal: {
                displayDetails: 'Display detailed list',
                resetButtonTitle: 'Clear search',
                backButtonTitle: 'Close search',
                noResultsText: 'No results found',
                footer: {
                  selectText: 'Select',
                  selectKeyAriaLabel: 'Enter',
                  navigateText: 'Navigate',
                  navigateUpKeyAriaLabel: 'Up',
                  navigateDownKeyAriaLabel: 'Down',
                  closeText: 'Close',
                  closeKeyAriaLabel: 'Escape',
                },
              },
            },
          },
        },
      },
    },
  },
  locales: {
    root: {
      label: '简体中文',
      lang: 'zh-CN',
      title: 'Fast API 文档',
      description: 'Fast API 接入与使用文档',
      themeConfig: zhTheme,
    },
    en: {
      label: 'English',
      lang: 'en-US',
      link: '/en/',
      title: 'Fast API Docs',
      description: 'Fast API integration and usage documentation',
      themeConfig: enTheme,
    },
  },
})
