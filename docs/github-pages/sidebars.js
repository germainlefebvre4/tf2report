/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  tutorialSidebar: [
    {
      type: 'category',
      label: 'Introduction',
      items: [
        'intro',
        'getting-started',
      ],
    },
    {
      type: 'category',
      label: 'User Guide',
      items: [
        'installation',
        'usage',
        'configuration',
        'output-formats',
        'filtering',
      ],
    },
    {
      type: 'category',
      label: 'Examples & Tutorials',
      items: [
        'examples/basic-usage',
        'examples/filtering',
        'examples/ci-cd-integration',
      ],
    },
    {
      type: 'category',
      label: 'Development',
      items: [
        'development/getting-started',
        'development/project-structure',
        'development/building',
        'development/testing',
        'development/api-reference',
        'development/contributing',
      ],
    },
    {
      type: 'category',
      label: 'Reference',
      items: [
        'reference/cli-options',
        'reference/configuration-schema',
        'reference/exit-codes',
      ],
    },
  ],
};

module.exports = sidebars;
