
## 使用 Hyperledger Fabric SDK Node 进行测试

[Hyperledger Fabric Client SDK](https://github.com/hyperledger/fabric-sdk-node) 能够非常简单的使用API和 Hyperledger Fabric Blockchain 网络进行交互。其`v1.1`及其以上的版本添加了一个重要的功能[Conection-Profile](https://fabric-sdk-node.github.io/tutorial-network-config.html)来保存整个network中必要的配置信息，方便client读取和配置。
该Demo基于`Connection-Profile`测试了整个网络的如下功能：
* Fabric CA 相关
  * Enroll用户
  * Register用户
* Channel 相关
  * 创建Channel
  * 将指定Peer join Channel
  * 查询Channel相关信息
  * 动态更新Channel配置信息
* Chaincode 相关
  * Install Chaincode
  * Instantiate Chaincode
  * Invoke Chaincode
  * Query Chaincode
  * 查询Chaincode相关信息

### 主要依赖

* Node v8.9.0 或更高 （注意目前v9.0+还不支持）
* npm v5.5.1 或更高
* gulp命令。 必须要进行全局安装 `npm install -g gulp`
* docker运行环境
* docker compose工具

主要fabric环境可参考[Fabric 1.0](https://github.com/yeasy/blockchain_guide/blob/master/fabric/1.0.md)。

### 下载 Demo 工程

```sh
$ git clone https://github.com/Sunnykaby/Hyperledger-fabric-node-sdk-demo
```


进入 `Hyperledger-fabric-node-sdk-demo` 目录，查看各文件夹和文件，功能如下。

文件/文件夹 | 功能 
-- | --
artifacts-local | 本地准备好构建fabric网络的基础材料
artifacts-remote | 使用官方fabric-sample动态构建网络
extra | 一些拓展性的材料
node |  基于Fabric SDK Node的demo核心代码 
src | 测试用chaincode
Init.sh | 构建Demo的初始化脚本

### 构建Demo

该项目提供两种Demo构建方式：
* 利用本地已经准备好的相关网络资源，启动fabric network。
* 利用官方fabric-sample项目，动态启动fabric network。

当然，你也可以使用自己已经创建好的fabric network和其相关的connection-profile来测试Demo。

```sh
##进入项目根目录

##使用本地资源构建Demo
./Init.sh local

##使用官方资源构建Demo
./Init.sh remote
```

执行之后，会在根目录中生成一个`demo`文件夹，其就是Demo程序的入口。

清理Demo资源，使用`./Init.sh clean`

### 启动Fabric网络

首先，我们需要准备一个fabric网络来进行测试。
进入到`demo`文件夹。

#### 本地资源构建网络

进入资源目录，利用脚本启动网络即可。
```sh
cd artifacts
##启动网络
./net.sh up
##关闭网络
./net.sh down
```
用该脚本启动网络中包含：1个orderer， 2个organisation， 4个peer（每个组织有2个peer）和两个ca（每个组织一个）。

#### 官方资源构建网络

在demo目录，利用脚本启动网络即可。
```sh
##启动网络,并配置本地资源
./net.sh init
##关闭网络并清理资源
./net.sh clean
```
用该脚本启动网络中包含：1个orderer， 2个organisation， 4个peer（每个组织有2个peer）和两个ca（每个组织一个）。

与本地资源启动不同，该方案主要有以下步骤：
* 将官方[fabric-sample项目](https://github.com/hyperledger/fabric-samples)clone到本地
* 利用`fabric-sample/first-network/bynf.sh up`启动fabric脚本
* 将一些资源文件连接到指定位置，方便node程序使用
* 通过资源文件构建connection-profile（替换密钥等）
* 创建一个新的channel的binary

详细信息可以直接查看`net.sh`脚本。

>`clean`命令会将所有相关的docker 容器和remote的动态资源全部删除。还原到最初的demo文件状态。

#### 资源清单

无论是remote还是local模式，最终资源和网络准备完成之后，核心资源列表如下：
```
demo/artifacts/  
├── channel-artifacts                
│   ├── channel2.tx    
│   ├── channel.tx  
│   ├── genesis.block  
│   ├── Org1MSPanchors.tx  
│   └── Org2MSPanchors.tx  
├── connection-profile              
│   ├── network.yaml  
│   ├── org1.yaml  
│   ├── org2.yaml  
├── crypto-config  
│   ├── ordererOrganizations  
│   │   └── example.com  
│   └── peerOrganizations  
│       ├── org1.example.com  
│       └── org2.example.com  
```

### 运行Demo

网络和相关资源准备成功之后，进入`demo/node`目录。
其主要结构为：
```
├── app                             //核心应用接口
│   ├── api-handler.js              //接口定义文件
│   ├── *.js                        //应用实现模块
│   ├── tools                       //通用工具类
│   │   ├── ca-tools.js
│   │   ├── config-tool.js
│   │   └── helper.js
├── app-test.js                     //Demo程序启动文件
├── package.json
└── readme.md
```

使用命令`node app-test.js`即可进行一个完整workflow的测试，包括最开始我们提到的所有功能。
同时可以使用`node app-test.js -m ca|createChannel|joinChannel|install|instantiate|invoke|query|queryChaincodeInfo|queryChannelInfo`来运行单个功能。

程序使用的均为默认参数，其定义在`app-test.js`文件中。可以按照需求修改对应的参数，再运行程序即可。

### 持续更新

如果在使用途中发现任何问题，或者有任何需求可以在该项目的issue中提出改进方案或者建议。
Github地址：[Hyperledger-fabric-node-sdk-demo](https://github.com/Sunnykaby/Hyperledger-fabric-node-sdk-demo)
