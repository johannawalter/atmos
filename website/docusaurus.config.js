// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion
// https://ricard.dev/how-to-set-docs-as-homepage-for-docusaurus
// https://docusaurus.io/docs/api/themes/configuration#theme
// https://docusaurus.io/docs/markdown-features/code-blocks#line-highlighting
// https://github.com/FormidableLabs/prism-react-renderer/tree/master/src/themes

const lightCodeTheme = require('prism-react-renderer/themes/vsDark');
const darkCodeTheme = require('prism-react-renderer/themes/nightOwl');

const BASE_URL = '';

/** @type {import('@docusaurus/types').Config} */
const config = {
    title: 'atmos',
    tagline: 'Universal tool for DevOps and Cloud Automation',
    url: 'https://atmos.tools',
    baseUrl: `${BASE_URL}/`,
    onBrokenLinks: 'throw',
    onBrokenMarkdownLinks: 'warn',
    favicon: 'img/atmos-logo.png',

    // GitHub pages deployment config.
    // If you aren't using GitHub pages, you don't need these.
    organizationName: 'cloudposse',
    projectName: 'atmos',

    // Even if you don't use internalization, you can use this field to set useful
    // metadata like html lang. For example, if your site is Chinese, you may want
    // to replace "en" with "zh-Hans".
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
                    routeBasePath: '/',
                    sidebarPath: require.resolve('./sidebars.js'),
                    editUrl: ({versionDocsDirPath, docPath, locale}) => {
                        return `https://github.com/cloudposse/atmos/edit/master/website/${versionDocsDirPath}/${docPath}`;
                    },
                    exclude: ['README.md'],
                },
                blog: {
                    showReadingTime: true,
                    editUrl: ({versionDocsDirPath, docPath, locale}) => {
                        return `https://github.com/cloudposse/atmos/edit/master/website/${versionDocsDirPath}/${docPath}`;
                    },
                    exclude: ['README.md'],
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
            docs: {
                sidebar: {
                    hideable: true,
                },
            },
            navbar: {
                title: 'atmos',
                logo: {
                    alt: 'Atmos Logo',
                    src: '/img/atmos-logo.svg',
                    srcDark: '/img/atmos-logo.svg',
                    href: '/',
                    target: '_self',
                    height: 36
                },
                items: [
                    {
                        type: 'doc',
                        docId: 'introduction',
                        position: 'left',
                        label: 'Docs',
                    },
                    {
                        to: '/cli/configuration',
                        position: 'left',
                        label: 'CLI'
                    },
                    {
                        to: 'https://github.com/cloudposse/community/discussions',
                        label: 'GitHub Discussions',
                        position: 'right'
                    },
                    {
                        to: 'https://slack.cloudposse.com/',
                        label: 'Slack Community',
                        position: 'right'
                    },
                    {
                        to: 'https://cloudposse.com/services/',
                        label: 'Get Help',
                        position: 'right',
                        className: 'button button--primary navbar-cta-button'
                    },
                    {
                        href: 'https://github.com/cloudposse/atmos',
                        position: 'right',
                        className: 'header-github-link',
                        'aria-label': 'GitHub repository',
                    }
                ],
            },
            prism: {
                theme: lightCodeTheme,
                darkTheme: darkCodeTheme,
                // https://prismjs.com/#supported-languages
                additionalLanguages: ['hcl', 'bash']
            },
            algolia: {
                appId: process.env.ALGOLIA_APP_ID || '32YOERUX83',
                apiKey: process.env.ALGOLIA_SEARCH_API_KEY || '557985309adf0e4df9dcf3cb29c61928', // this is SEARCH ONLY API key and is not sensitive information
                indexName: process.env.ALGOLIA_INDEX_NAME || 'docs.cloudposse.com',
                contextualSearch: false
            },
        }),

    markdown: {
        mermaid: true,
    },

    themes: ['@docusaurus/theme-mermaid']
};

module.exports = config;
