## Python 客户端
前面应用案例，都是直接通过 HTTP API 来跟 hyperledger 进行交互，操作比较麻烦。

还可以直接通过 [hyperledger-py](https://github.com/yeasy/hyperledger-py) 客户端来进行更方便的操作。

### 安装

```sh
$ pip install hyperledger --upgrade
```

或直接源码安装

```sh
$ git clone https://github.com/yeasy/hyperledger-py.git
$ cd hyperledger-py
$ pip install -r requirements.txt
$ python setup.py install
```

### 使用

```py
>>> from hyperledger.client import Client
>>> c = Client(base_url="http://127.0.0.1:7050")
>>> c.peer_list()
{u'peers': [{u'type': 1, u'ID': {u'name': u'vp1'}, u'address': u'172.17.0.2:30303'}, {u'type': 1, u'ID': {u'name': u'vp2'}, u'address': u'172.17.0.3:30303'}]}
```

更多使用方法，可以参考 [API 文档](https://github.com/yeasy/hyperledger-py/blob/master/docs/api.md)。



