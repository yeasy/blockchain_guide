## 开发必备工具

工欲善其事，必先利其器。开源社区提供了大量易用的开发协作工具。掌握好这些工具，对于高效的开发来说十分重要。

### Linux Foundation ID

超级账本项目受到 Linux 基金会的支持，采用 Linux Foundation ID（LF ID）作为社区唯一的 ID。

个人申请 ID 是完全免费的。可以到 https://identity.linuxfoundation.org/ 进行注册。

用户使用该 ID 即可访问到包括 Jira、Gerrit、RocketChat 等社区的开发工具。

### Jira - 任务和进度管理

![Jira 任务管理](_images/jira.png)

Jira 是 Atlassian 公司开发的一套任务管理和事项跟踪的追踪平台，提供 Web 操作界面，使用十分方面。

社区采用 jira.hyperledger.org 作为所有项目开发计划和任务追踪的入口，使用 LF ID 即可登录。

登录之后，可以通过最上面的 Project 菜单来查看某个项目相关的事项，还可以通过 Create 按钮来快速创建事项（常见的包括 task、bug、improvement 等）。

用户打开事项后可以通过 assign 按钮分配给自己来领取该事项。

一般情况下，事项分为 TODO（待处理）、In Process（处理中）、In Review（补丁已提交、待审查）、Done（事项已完成）等多个状态，由事项所有者来进行维护。

### Gerrit - 代码仓库和 Review 管理

![Gerrit 代码仓库管理](_images/gerrit.png)

Gerrit 是一个负责代码协同的开源项目，很多企业和团队都使用它负责代码仓库管理和代码的审阅工作。Gerrit 使用十分方便，提供了基于 Web 的操作界面。


社区采用 gerrit.hyperledger.org 作为官方的代码仓库，并实时同步代码到 github.com/hyperledger 作为只读的镜像。

用户使用自己的 LF ID 登录之后，可以查看所有项目信息，也可以查看自己提交的补丁等信息。每个补丁的页面上会自动追踪修改历史，审阅人可以通过页面进行审阅操作，赞同提交则可以加分，发现问题则注明问题并进行减分。

### RocketChat - 在线沟通

![RocketChat 在线沟通](_images/rocket_chat.png)

除了邮件列表外，社区也为开发者们提供了在线沟通的渠道—— RocketChat。

RocketChat 是一款功能十分强大的在线沟通软件，支持多媒体消息、附件、提醒、搜索等功能，虽然是开源软件，但在体验上可以跟商业软件 Slack 媲美。支持包括网页、桌面端、移动端等多种客户端。

社区采用 chat.hyperledger.org 作为服务器。最简单的，用户直接使用自己的 LF ID 登录该网站，即可访问。之后可以自行添加感兴趣项目的频道。

用户也可以下载 RocketChat 客户端，添加 chat.hyperledger.org 作为服务器即可访问社区内的频道，跟广大开发者进行在线交流。

一般地，每个项目都有一个同名的频道作为主频道，例如 `#Fabric`，`#Cello` 等。同时各个工作组也往往有自己的频道，例如大中华区技术工作组的频道为 `#twg-china`。



