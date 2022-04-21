/**
 * 网站配置文件
 */

const config = {
  appName: 'Owl',
  // todo, set a better image addr
  appLogo: 'https://img0.baidu.com/it/u=2822765666,2555722031&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=501',
  showViteLogo: true
}

export const viteLogo = (env) => {
  if (config.showViteLogo) {
    const chalk = require('chalk')
    console.log(
      chalk.green(
        `> 欢迎使用Owls，开源地址：https://github.com/qingfeng777/owls`
      )
    )
    console.log(
      chalk.green(
        `> 当前版本:V0.0.1`
      )
    )
    console.log(
      chalk.green(
        `> 加群方式:微信：xxxx QQ群：xxxx`
      )
    )
    console.log(
      chalk.green(
        `> 默认自动化文档地址:http://127.0.0.1:${env.VITE_SERVER_PORT}/swagger/index.html`
      )
    )
    console.log(
      chalk.green(
        `> 默认前端文件运行地址:http://127.0.0.1:${env.VITE_CLI_PORT}`
      )
    )
    console.log(
      chalk.green(
        `> 如果项目让您获得了收益，那就帮忙宣传一下吧！`
      )
    )
    console.log('\n')
  }
}

export default config
