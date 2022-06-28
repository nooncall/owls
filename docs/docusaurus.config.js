// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Owls',
  tagline: '稳定、高效、创造',
  url: 'https://your-docusaurus-test-site.com',
  baseUrl: '/docs/',
  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'nooncall', // Usually your GitHub org/user name.
  projectName: 'owls', // Usually your repo name.

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  plugins: [
    [
      'content-docs',
      /** @type {import('@docusaurus/plugin-content-docs').Options} */
      ({
        id: 'user_guide',
        path: 'user_guide',
        routeBasePath: 'user_guide',
        editUrl: ({locale, versionDocsDirPath, docPath}) => {
          if (locale !== 'en') {
            return `https://crowdin.com/project/docusaurus-v2/${locale}`;
          }
          return `https://github.com/facebook/docusaurus/edit/main/website/${versionDocsDirPath}/${docPath}`;
        },
        //remarkPlugins: [npm2yarn],
        editCurrentVersion: true,
        // sidebarPath: require.resolve('./sidebarsCommunity.js'),
        showLastUpdateAuthor: true,
        showLastUpdateTime: true,
      }),
    ],
    [
      'content-docs',
      /** @type {import('@docusaurus/plugin-content-docs').Options} */
      ({
        id: 'db',
        path: 'db',
        routeBasePath: 'db',
        editUrl: ({locale, versionDocsDirPath, docPath}) => {
          if (locale !== 'en') {
            return `https://crowdin.com/project/docusaurus-v2/${locale}`;
          }
          return `https://github.com/facebook/docusaurus/edit/main/website/${versionDocsDirPath}/${docPath}`;
        },
        editCurrentVersion: true,
        showLastUpdateAuthor: true,
        showLastUpdateTime: true,
      }),
    ],
    [
      'content-docs',
      /** @type {import('@docusaurus/plugin-content-docs').Options} */
      ({
        id: 'roadmap',
        path: 'roadmap',
        routeBasePath: 'roadmap',
        editUrl: ({locale, versionDocsDirPath, docPath}) => {
          if (locale !== 'en') {
            return `https://crowdin.com/project/docusaurus-v2/${locale}`;
          }
          return `https://github.com/facebook/docusaurus/edit/main/website/${versionDocsDirPath}/${docPath}`;
        },
        editCurrentVersion: true,
        showLastUpdateAuthor: true,
        showLastUpdateTime: true,
      }),
    ],
  ],

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/facebook/docusaurus/tree/main/packages/create-docusaurus/templates/shared/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/facebook/docusaurus/tree/main/packages/create-docusaurus/templates/shared/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: 'Owls',
        logo: {
          alt: 'Owls',
          src: 'img/docusaurus.png',
        },
        items: [
          
          {to: '/user_guide/intro', label: '用户手册', position: 'left', activeBaseRegex: `/user_guide/`,},
          {to: '/db/intro', label: 'DB', position: 'left', activeBaseRegex: `/db/`,},
          {to: '/roadmap/intro', label: '研发文档', position: 'left', activeBaseRegex: `/roadmap/`,},

          // {
          //   type: 'doc',
          //   docId: 'intro',
          //   position: 'left',
          //   label: '文档',
          // },
          {to: '/blog', label: 'Blog', position: 'left'},
          {
            href: 'https://github.com/facebook/docusaurus',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Tutorial',
                to: '/docs/intro',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'Stack Overflow',
                href: 'https://stackoverflow.com/questions/tagged/docusaurus',
              },
              {
                label: 'Discord',
                href: 'https://discordapp.com/invite/docusaurus',
              },
              {
                label: 'Twitter',
                href: 'https://twitter.com/docusaurus',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'Blog',
                to: '/blog',
              },
              {
                label: 'GitHub',
                href: 'https://github.com/facebook/docusaurus',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} My Project, Inc. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;
