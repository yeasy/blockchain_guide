###安装部署

如果你是首次接触ethereum，推荐使用下面的步骤安装部署。

####安装Go
    curl -O https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz
    mkdir -p ~/go; echo "export GOPATH=$HOME/go" >> ~/.bashrc
    echo "export PATH=$PATH:$HOME/go/bin:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc


####安装ethereum
    sudo apt-get install software-properties-common
    sudo add-apt-repository -y ppa:ethereum/ethereum
    sudo add-apt-repository -y ppa:ethereum/ethereum-dev
    sudo apt-get update
    sudo apt-get install ethereum


####安装solc编译器

    sudo add-apt-repository ppa:ethereum/ethereum-qt
    sudo add-apt-repository ppa:ethereum/ethereum
    sudo apt-get update
    sudo apt-get install cpp-ethereum
    
安装后可以使用geth命令创建ethereum账户
```
    geth account new
```
