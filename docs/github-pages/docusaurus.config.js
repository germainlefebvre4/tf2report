// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer').themes.github;
const darkCodeTheme = require('prism-react-renderer').themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'tf2report',
  tagline: 'Summarize Terraform plan changes into human-readable reports',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://tf2report.germainlefebvre.fr',
  // Set the /<baseUrl>/ pathname under which your site is served
  baseUrl: '/',

  // GitHub pages deployment config.
  organizationName: 'germainlefebvre4',
  projectName: 'tf2report',

  onBrokenLinks: 'throw',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl: 'https://github.com/germainlefebvre4/tf2report/tree/main/docs/github-pages/',
        },
        blog: false,
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      image: 'img/tf2report-social-card.png',
      navbar: {
        title: 'tf2report',
        logo: {
          alt: 'tf2report Logo',
          src: 'img/logo.svg',
        },
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'tutorialSidebar',
            position: 'left',
            label: 'Documentation',
          },
          {
            href: 'https://github.com/germainlefebvre4/tf2report',
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
                label: 'Getting Started',
                to: '/docs/getting-started',
              },
              {
                label: 'Usage Guide',
                to: '/docs/usage',
              },
              {
                label: 'Configuration',
                to: '/docs/configuration',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/germainlefebvre4/tf2report',
              },
              {
                label: 'Issues',
                href: 'https://github.com/germainlefebvre4/tf2report/issues',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'Development',
                to: '/docs/development/getting-started',
              },
              {
                label: 'API Reference',
                to: '/docs/development/api-reference',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Germain LEFEBVRE. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
        additionalLanguages: ['bash', 'json', 'yaml', 'hcl', 'go'],
      },
      algolia: {
        // The application ID provided by Algolia
        appId: 'YOUR_APP_ID',
        // Public API key: it is safe to commit it
        apiKey: 'YOUR_SEARCH_API_KEY',
        indexName: 'tf2report',
        // Optional: see doc section below
        contextualSearch: true,
        // Optional: Algolia search parameters
        searchParameters: {},
        // Optional: path for search page that enabled by default (`false` to disable it)
        searchPagePath: 'search',
      },
    }),

  markdown: {
    hooks: {
      onBrokenMarkdownLinks: 'warn',
    },
  },
};

module.exports = config;
