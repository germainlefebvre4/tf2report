import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Easy to Use',
    description: (
      <>
        tf2report is designed from the ground up to be easily installed and
        used to analyze Terraform plans and generate reports quickly.
      </>
    ),
  },
  {
    title: 'Multiple Output Formats',
    description: (
      <>
        Generate reports in Markdown, plain text, or JSON format. Perfect for
        documentation, terminal output, or programmatic consumption.
      </>
    ),
  },
  {
    title: 'Powerful Filtering',
    description: (
      <>
        Filter changes by resource type and action. Focus on what matters most
        for code reviews, security audits, or compliance checks.
      </>
    ),
  },
  {
    title: 'CI/CD Ready',
    description: (
      <>
        Integrate seamlessly into GitHub Actions, GitLab CI, or any CI/CD pipeline.
        Automatically generate reports for every infrastructure change.
      </>
    ),
  },
  {
    title: 'Configuration Driven',
    description: (
      <>
        Use YAML configuration files for repeatable setups. Define filters,
        output formats, and defaults once and reuse across environments.
      </>
    ),
  },
  {
    title: 'Open Source',
    description: (
      <>
        Built with Go, fully open source. Contribute, extend, or customize
        tf2report to fit your infrastructure workflow needs.
      </>
    ),
  },
];

function Feature({title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
