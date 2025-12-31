# tf2report Documentation Site

This directory contains the comprehensive Docusaurus documentation for tf2report.

## Quick Start

### Install Dependencies

```bash
cd docs/github-pages
npm install
```

### Start Development Server

```bash
npm start
```

The site opens at `http://localhost:3000`.

### Build for Production

```bash
npm run build
```

Output in `build/` directory.

### Serve Production Build

```bash
npm run serve
```

## Documentation Structure

```
docs/github-pages/
├── docs/                       # Documentation content
│   ├── intro.md               # Introduction
│   ├── getting-started.md     # Quick start guide
│   ├── installation.md        # Installation instructions
│   ├── usage.md               # Usage guide
│   ├── configuration.md       # Configuration reference
│   ├── output-formats.md      # Output formats guide
│   ├── filtering.md           # Filtering guide
│   ├── examples/              # Examples and tutorials
│   │   ├── basic-usage.md
│   │   ├── filtering.md
│   │   └── ci-cd-integration.md
│   ├── development/           # Developer documentation
│   │   ├── getting-started.md
│   │   ├── project-structure.md
│   │   ├── building.md
│   │   ├── testing.md
│   │   ├── api-reference.md
│   │   └── contributing.md
│   └── reference/             # Reference documentation
│       ├── cli-options.md
│       ├── configuration-schema.md
│       └── exit-codes.md
├── src/                        # React components
│   ├── components/
│   │   └── HomepageFeatures/
│   ├── css/
│   │   └── custom.css
│   └── pages/
│       └── index.js           # Homepage
├── static/                     # Static assets
│   └── img/
├── docusaurus.config.js       # Docusaurus configuration
├── sidebars.js                # Sidebar structure
├── package.json               # Dependencies
└── README.md                  # This file
```

## Documentation Categories

### User Documentation

- **Introduction** - Overview and key features
- **Getting Started** - 5-minute quick start
- **Installation** - Installation methods
- **Usage** - Comprehensive usage guide
- **Configuration** - Configuration file reference
- **Output Formats** - Format details
- **Filtering** - Filtering techniques

### Examples & Tutorials

- **Basic Usage** - Simple examples
- **Filtering** - Advanced filtering
- **CI/CD Integration** - Pipeline integration
- **Advanced Scenarios** - Complex use cases

### Developer Documentation

- **Getting Started** - Development setup
- **Project Structure** - Codebase organization
- **Building** - Build instructions
- **Testing** - Testing guide
- **API Reference** - Package API documentation
- **Contributing** - Contribution guidelines

### Reference

- **CLI Options** - Command-line reference
- **Configuration Schema** - YAML schema
- **Exit Codes** - Exit code reference

## Customization

### Update Configuration

Edit `docusaurus.config.js`:

```javascript
const config = {
  title: 'tf2report',
  tagline: 'Summarize Terraform plan changes into human-readable reports',
  url: 'https://germainlefebvre4.github.io',
  baseUrl: '/tf2report/',
  // ... more config
};
```

### Update Sidebar

Edit `sidebars.js` to change navigation structure.

### Add New Page

1. Create `.md` file in `docs/`
2. Add front matter:
   ```yaml
   ---
   sidebar_position: 1
   ---
   ```
3. Page automatically appears in sidebar

### Update Homepage

Edit `src/pages/index.js` and `src/components/HomepageFeatures/index.js`.

## Deployment

### GitHub Pages

The site is configured to deploy to GitHub Pages at:
`https://germainlefebvre4.github.io/tf2report/`

### Deploy Command

```bash
npm run build
```

Then commit and push the `build/` directory or use GitHub Actions for automated deployment.

### GitHub Actions Deployment

Create `.github/workflows/deploy-docs.yml`:

```yaml
name: Deploy Docs

on:
  push:
    branches: [main]
    paths:
      - 'docs/github-pages/**'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 18
      - name: Install dependencies
        run: cd docs/github-pages && npm install
      - name: Build
        run: cd docs/github-pages && npm run build
      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./docs/github-pages/build
```

## Contributing to Documentation

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally with `npm start`
5. Build with `npm run build`
6. Submit a pull request

## Documentation Best Practices

1. **Clear and Concise** - Use simple language
2. **Code Examples** - Include working examples
3. **Screenshots** - Add visuals where helpful
4. **Links** - Cross-reference related content
5. **Up-to-date** - Keep in sync with code changes

## Search

The site includes Algolia search (configured in `docusaurus.config.js`). Update search credentials when deploying.

## Support

- Report documentation issues on GitHub
- Suggest improvements via pull requests
- Ask questions in discussions

## License

Documentation is part of the tf2report project and follows the same license.
