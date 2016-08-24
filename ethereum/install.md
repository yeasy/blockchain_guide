## 安装部署

如果你是首次接触 ethereum，推荐使用下面的步骤安装部署。

### 安装 Go 环境
```sh
    curl -O https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz
    mkdir -p ~/go; echo "export GOPATH=$HOME/go" >> ~/.bashrc
    echo "export PATH=$PATH:$HOME/go/bin:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
```

### 安装 ethereum
```sh
    sudo apt-get install software-properties-common
    sudo add-apt-repository -y ppa:ethereum/ethereum
    sudo add-apt-repository -y ppa:ethereum/ethereum-dev
    sudo apt-get update
    sudo apt-get install ethereum
```

### 安装 solc 编译器

```sh
    sudo add-apt-repository ppa:ethereum/ethereum-qt
    sudo add-apt-repository ppa:ethereum/ethereum
    sudo apt-get update
    sudo apt-get install cpp-ethereum
```
    
安装后可以使用 geth 命令创建 ethereum 账户

```sh
    geth account new
```

### Solidity 语言支持

[browser-solidity](https://ethereum.github.io/browser-solidity/#version=0.1.1) 提供了在线的 Solidity 语言测试。

用户也可以从 [这里](https://github.com/ethereum/webthree-umbrella/releases) 下载到包括 Solidity 运行环境的安装包。
