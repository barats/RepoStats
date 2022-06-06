# RepoStats 

开源代码仓库的 star、fork、commit、pull request、issue 等相关数据，是分析和了解代码仓库的客观依据，这些数据在一定程度上反应了开源项目的受欢迎程度、活跃度、影响力等。

RepoStats 致力于解决的痛点问题是：  
1. 开源代码仓库的数据抓取、存储、分析及统计
1. 开源代码仓库的相关数据可视化展示
1. 做到全网、全平台平台打通，并支持分隔、组合展示


## 功能说明
1. 支持平台：当前版本的 RepoStats 仅支持 Gitee 平台相关数据获取   
2. 管理后台：支持界面化的 Gitee Oauth 配置、Grafana Token 获取配置
3. 管理后台：支持添加单个仓库、支持批量添加个人帐号及组织帐号下的公开仓库
4. 管理后台：支持禁用、启用 Gitee 数据抓取(启动抓取除外)
5. 管理后台：支持 Commit 列表显示及查询、Issue 列表显示及查询、Pull Request 列表显示及查询
6. Grafana 标签：每个面板均有附带仓库拥有者标签、仓库名称、平台名称等信息支持查询过滤

<p align="center">
<a target="_blank" href="https://www.repostats.cn">https://www.repostats.cn</a> <br/>
<a target="_blank" href="https://github.com/barats/RepoStats/stargazers"><img src="https://img.shields.io/github/stars/barats/RepoStats"/></a>
<a target="_blank" href="https://github.com/barats/RepoStats/network/members"><img src="https://img.shields.io/github/forks/barats/RepoStats"/></a>
<a target="_blank" href="https://github.com/barats/RepoStats/issues"><img src="https://img.shields.io/github/issues/barats/RepoStats"/></a>  
<a target="_blank" href='https://gitee.com/barat/repostats/stargazers'><img src='https://gitee.com/barat/repostats/badge/star.svg?theme=dark' /></a>
<a target="_blank" href='https://gitee.com/barat/repostats/members'><img src='https://gitee.com/barat/repostats/badge/fork.svg?theme=dark' /></a>
<a target="_blank" href='https://www.oschina.net/comment/project/64186'><img src='https://www.oschina.net/comment/badge/project/64186'/></a>
</p>

## 安装及配置说明

### 1. 使用 Docker 环境
启动 Docker 镜像前请注意查看 `docker/vars.env` 文件，并根据自己的实际情况调整需要的参数(eg：本地端口号等)。启动 `docker/start_docker_repostats.sh` 可以通过 Docker 环境将所有依赖安装并启动。该命令将：
1. 拉取 [repostats](https://hub.docker.com/r/baratsemet/repostats) 镜像(本地构建可查看 `docker/repostats.Dockerfile` 文件) 
1. 通过 `docker/pull_build.yml` 其他描述内容构建 `Grafana` 和 `PostgreSQL` 镜像及服务，并对其运行状态做判断，再启动其他必要服务 (本地构建镜像请查阅 `local_build.yml`) 
1. 构建名为 `network_repostats` 的虚拟网络供上述服务使用
1. 开启本机 `9103` 端口应对 RepoStats 工具、启动 `13000` 端口应对 Grafana 工具，启动 `15432` 应对 PostgreSQL 数据库

### 2. 通过 `Makefile` 构建

构建 linux 平台对应的可执行文件：

```shell
make build-linux
```

压缩 linux 平台对应的可执行文件(压缩可执行文件需要 [upx](https://github.com/upx/upx) 支持)：

```shell
make compress-linux
```

### 3. 使用各系统分发版本

通过 Release 下载对应平台的分发版并启动运行，启动之前请确保 `repostats.ini` 配置文件中各项内容的正确性

```ini
[repostats]
debug = false
admin_port = 9103


[postgres]
host = localhost
port = 15432
user = postgres
password = DePmoG_123
database = repostats
max_open_conn = 20
max_idle_conn = 5

[grafana]
host = localhost
port = 13000
user = admin
password = admin
```

数据库说明

1. 请在 `PostgreSQL` 数据库中创建名称 `repostats` 的数据库
2. 请分别 **顺序执行** `sql/db.sql`、`sql/gitee.sql`、`sql/roles.sql` 创建必要的表及视图等

启动参数说明

```shell
repostats [-c config_file]
```

## 使用前置条件

RepoStats 启动成功之后，请使用帐号密码登录 Admin 管理后台。管理后台默认的账户名密码信息如下：  

```
repostats
-2aDzm=0(ln_9^1
```

在开始爬取 Gitee 数据并向 Grafana 推送相关统计结果之前，需要对其进行一定的配置：
1. 在 Admin 管理界面中的 `Gitee 配置` 页面根据提示配置必要的 `Oauth 参数` 从而保证能够正常获取 Gitee 相关数据 
1. 在 Admin 管理界面中的 `Grafana 配置` 页面中根据提示配置必要的 `Grafana 参数` 从而确保能够与 Grafana 进行通信
1. 上一步骤中，与 Grafana 通信成功之后，需要在 Grafana 操作页面中对数据源进行一次 `Test & Save` 才能保证数据源正常 (这个问题暂时似乎是没办法处理的，后续会想办法处理)
1. 在 Admin 管理界面中的 `代码仓库列表` 页面中，根据界面提示添加想要关注的代码仓库

## 截图分享

所有仓库总视图
![所有仓库总视图](https://oscimg.oschina.net/oscnet/up-1d0f56655abc5a92846614e9862620e55b4.jpg)

指定某个仓库的视图  
![指定某个仓库的视图](https://oscimg.oschina.net/oscnet/up-6fac497c4428602cc7a44a363c7b674165a.jpg)

Admin 后端管理界面 
![Admin 后端管理界面](https://oscimg.oschina.net/oscnet/up-101d6ca0c57de648c7fa20ec7b3f863fcd6.jpg)

## 统计指标说明

RepoStats 当前版本支持3大类共计21项统计数据可视化结果展示，这些统计数据**不能表示一个开源项目的好与坏**，仅从数据层面对开源代码仓库进行一定的展示。这些数据指标分类以下三类：

### 1. 统计汇总

`统计汇总` 分类中展示的数据，与时间无关，它们代表的是所有项目(Gitee Overview)后者是某个指定的项目的汇总数据结果，其中包括：

- 仓库统计  
当前抓取的仓库总数量、总 Star 人数、总 Fork 人数、总 Watch 人数

- 基本信息  
当前仓库的 Star 人数、Fork 人数、Watch 人数

- Commit 统计  
Commit 总数、Commit Author 总数(去重)、Commit Committer 总数(去重)

- Issue 统计  
Issue 综述、Issue 总人数、打开状态的 Issue 总数、已关闭状态的 Issue 总数、已拒绝状态的 Issue 总数、处理中的 Issue 总数 

- Issue 状态图  
已关闭、已拒绝、打开、处理中 状态的 Issue 占比示意图

- Pull Request 统计  
Pull Request 总数、Pull Request 人数、打开状态的 Pull Request 总数、已合并状态的 Pull Request 总数、已关闭状态的 Pull Request 总数

- Pull Request 状态图  
已合并、打开、已关闭 的 Pull Reqeust 占比示意图 

- Issue 处理时间分析  
所有 Issue 从 `created_at` 到 `finished_at` 的最小耗时、平均耗时、最大耗时，单位：小时

- Pull Request 合并时间分析  
所有 `可合并的` Pull Request 从 `created_at` 到 `merged_at` 的最小耗时、平均耗时、最大耗时，单位：小时

### 2. 动态趋势

`动态趋势` 分类中展示的数据，是 `某个时间段内` 数据量的动态变化过程，可以通过 Grafana 面板右上角的时间选项查看指定时间范围内的变化趋势，其中包括： 

- Star 趋势图  
指定时间范围内，关注仓库的总人数变化趋势 

- Commit 趋势图  
指定时间范围内，Commit 提交次数的变化趋势

- Issue 趋势图  
指定时间范围内，新增 Issue 数的变化趋势

- Pull Request 趋势图  
指定时间范围内，新增 Pull Request 数的变化趋势

- Pull Request 合并时间分析  
指定时间范围内，`可合并的` Pull Request 从 `created_at` 到 `merged_at` 的最小耗时、平均耗时、最大耗时，单位：小时

- Issue 处理时间分析  
指定时间范围内，新增的 Issue 从 `created_at` 到 `finished_at` 的最小耗时、平均耗时、最大耗时，单位：小时

### 3. 数据列表

`数据列表` 分类中展示的数据，与时间无关，它们代表的是所有项目(Gitee Overview)后者是某个指定的项目的汇总数据结果，其中包括：

- 仓库列表  
所有仓库的明细列表

- Commit 列表  
Commit 明细列表

- Issue 列表  
Issue 明细列表

- Pull Request 列表    
Pull Request 明细列表

- Commit Author 排行  
Commit Auhtor 次数排行

- Commit Committer 排行 
Commit Committer 次数排行

## Contributor License Agreement

在 **第一次提交 Pull Request 时** ，请您在 Pull Request 内容中明确写明「本人自愿接受并签署 [《RepoStats Contributor License Agreement》](CLA.md)」，并在 Pull Request 信息中附带该协议链接信息。

## Inspired By 
1. [CNCF DevStats](https://devstats.cncf.io/)
1. [cncf/devstatscode](https://github.com/cncf/devstatscode)

## Give Thanks To

由衷感谢以下开源软件、框架等（包括但不限于）

1. [Grafana](https://grafana.com/)
1. [gin-gonic/gin](https://github.com/gin-gonic/gin) 
1. [FomanticUI](https://fomantic-ui.com/)
1. [dchest/captcha](https://github.com/dchest/captcha) 
1. [Masterminds/sprig](https://github.com/Masterminds/sprig)
1. [jmoiron/sqlx](https://github.com/jmoiron/sqlx)
1. [go-ini/ini](https://github.com/go-ini/ini)

## RepoStats News
1. 2022-06-01 [RepoStats 代码仓库数据可视化工具，路线图发布](https://www.oschina.net/news/198120/repostats-roadmap)
1. 2022-05-11 [RepoStats v1.1 正式发布，开源代码仓库数据可视化工具](https://www.oschina.net/news/195251/repostats-1-1-released)
1. 2022-04-27 [安装并使用 RepoStats 代码仓库数据可视化工具](https://mp.weixin.qq.com/s/St3OItSpgcxl_wuuIGnuIA)
1. 2022-04-26 [RepoStats v1.0 发布，开源代码仓库统计数据可视化](https://www.oschina.net/news/193100/repostats-1-0-released)