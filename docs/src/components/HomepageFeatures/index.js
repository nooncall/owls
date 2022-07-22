import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: '易于使用',
    Svg: require('@site/static/img/undraw_docusaurus_mountain.svg').default,
    description: (
      <>
        方便的用户接入系统、完善的审批流程、丰富的规则审批等，带给您高效、愉悦的使用体验。
      </>
    ),
  },
  {
    title: '带给您稳定',
    Svg: require('@site/static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
        通过规范化人与数据的交互接口，在线、离线的访问请求过滤，经验沉淀出的优化建议反馈等，使您的软件服务更加稳定。
      </>
    ),
  },
  {
    title: '解放创造力',
    Svg: require('@site/static/img/undraw_docusaurus_tree.svg').default,
    description: (
      <>
        尽可能的把时间和精力从无尽的故障、oncall中收回来，去做一些自己喜欢的、有创造力的、有意思的事情！
      </>
    ),
  },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
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
