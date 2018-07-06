## PKI 体系

在非对称加密中，公钥可以通过证书机制来进行保护，但证书的生成、分发、撤销等过程并没有在 X.509 规范中进行定义。

实际上，如何安全地管理和分发证书可以遵循 PKI（Public Key Infrastructure）体系来完成。

PKI 体系核心解决的是证书生命周期相关的认证和管理问题，在现代密码学应用领域处于十分基础和重要的地位。

需要注意，PKI 是建立在公私钥基础上实现安全可靠传递消息和身份确认的一个通用框架，并不代表某个特定的密码学技术和流程。实现了 PKI 规范的平台可以安全可靠地管理网络中用户的密钥和证书。目前包括多个实现和规范，知名的有 RSA 公司的 PKCS（Public Key Cryptography Standards）标准和 X.509 相关规范等。

### PKI 基本组件

一般情况下，PKI 至少包括如下核心组件：

* CA（Certification Authority）：负责证书的颁发和吊销（Revoke），接收来自 RA 的请求，是最核心的部分；
* RA（Registration Authority）：对用户身份进行验证，校验数据合法性，负责登记，审核过了就发给 CA；
* 证书数据库：存放证书，多采用 X.500 系列标准格式。可以配合LDAP 目录服务管理用户信息。

其中，CA 是最核心的组件，主要完成对证书信息的维护。

常见的操作流程为，用户通过 RA 登记申请证书，提供身份和认证信息等；CA 审核后完成证书的制造，颁发给用户。用户如果需要撤销证书则需要再次向 CA 发出申请。

### 证书的签发

CA 对用户签发证书实际上是对某个用户公钥，使用 CA 的私钥对其进行签名。这样任何人都可以用 CA 的公钥对该证书进行合法性验证。验证成功则认可该证书中所提供的用户公钥内容，实现用户公钥的安全分发。

用户证书的签发可以有两种方式。一般可以由 CA 直接来生成证书（内含公钥）和对应的私钥发给用户；也可以由用户自己生成公钥和私钥，然后由 CA 来对公钥内容进行签名。

后者情况下，用户一般会首先自行生成一个私钥和证书申请文件（Certificate Signing Request，即 csr 文件），该文件中包括了用户对应的公钥和一些基本信息，如通用名（common name，即 cn）、组织信息、地理位置等。CA 只需要对证书请求文件进行签名，生成证书文件，颁发给用户即可。整个过程中，用户可以保持私钥信息的私密性，不会被其他方获知（包括 CA 方）。

生成证书申请文件的过程并不复杂，用户可以很容易地使用开源软件 openssl 来生成 csr 文件和对应的私钥文件。

例如，安装 openssl 后可以执行如下命令来生成私钥和对应的证书请求文件：

```bash
$ openssl req -new -keyout private.key -out for_request.csr
Generating a 1024 bit RSA private key
...........................++++++
............................................++++++
writing new private key to 'private.key'
Enter PEM pass phrase:
Verifying - Enter PEM pass phrase:
-----
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:CN
State or Province Name (full name) [Some-State]:Beijing
Locality Name (eg, city) []:Beijing
Organization Name (eg, company) [Internet Widgits Pty Ltd]:Blockchain
Organizational Unit Name (eg, section) []:Dev
Common Name (e.g. server FQDN or YOUR name) []:example.com
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
```

生成过程中需要输入地理位置、组织、通用名等信息。生成的私钥和 csr 文件默认以 PEM 格式存储，内容为 base64 编码。

如生成的 csr 文件内容可能为：

```bash
$ cat for_request.csr                                                                                                                                       
-----BEGIN CERTIFICATE REQUEST-----
MIIBrzCCARgCAQAwbzELMAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxEDAO
BgNVBAcTB0JlaWppbmcxEzARBgNVBAoTCkJsb2NrY2hhaW4xDDAKBgNVBAsTA0Rl
djEZMBcGA1UEAxMQeWVhc3kuZ2l0aHViLmNvbTCBnzANBgkqhkiG9w0BAQEFAAOB
jQAwgYkCgYEA8fzVl7MJpFOuKRH+BWqJY0RPTQK4LB7fEgQFTIotO264ZlVJVbk8
Yfl42F7dh/8SgHqmGjPGZgDb3hhIJLoxSOI0vJweU9v6HiOVrFWE7BZEvhvEtP5k
lXXEzOewLvhLMNQpG0kBwdIh2EcwmlZKcTSITJmdulEvoZXr/DHXnyUCAwEAAaAA
MA0GCSqGSIb3DQEBBQUAA4GBAOtQDyJmfP64anQtRuEZPZji/7G2+y3LbqWLQIcj
IpZbexWJvORlyg+iEbIGno3Jcia7lKLih26lr04W/7DHn19J6Kb/CeXrjDHhKGLO
I7s4LuE+2YFSemzBVr4t/g24w9ZB4vKjN9X9i5hc6c6uQ45rNlQ8UK5nAByQ/TWD
OxyG
-----END CERTIFICATE REQUEST-----
```

openssl 工具提供了查看 PEM 格式文件明文的功能，如使用如下命令可以查看生成的 csr 文件的明文：

```bash
$ openssl req -in for_request.csr -noout -text
Certificate Request:
    Data:
        Version: 0 (0x0)
        Subject: C=CN, ST=Beijing, L=Beijing, O=Blockchain, OU=Dev, CN=yeasy.github.com
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
            RSA Public Key: (1024 bit)
                Modulus (1024 bit):
                    00:f1:fc:d5:97:b3:09:a4:53:ae:29:11:fe:05:6a:
                    89:63:44:4f:4d:02:b8:2c:1e:df:12:04:05:4c:8a:
                    2d:3b:6e:b8:66:55:49:55:b9:3c:61:f9:78:d8:5e:
                    dd:87:ff:12:80:7a:a6:1a:33:c6:66:00:db:de:18:
                    48:24:ba:31:48:e2:34:bc:9c:1e:53:db:fa:1e:23:
                    95:ac:55:84:ec:16:44:be:1b:c4:b4:fe:64:95:75:
                    c4:cc:e7:b0:2e:f8:4b:30:d4:29:1b:49:01:c1:d2:
                    21:d8:47:30:9a:56:4a:71:34:88:4c:99:9d:ba:51:
                    2f:a1:95:eb:fc:31:d7:9f:25
                Exponent: 65537 (0x10001)
        Attributes:
            a0:00
    Signature Algorithm: sha1WithRSAEncryption
        eb:50:0f:22:66:7c:fe:b8:6a:74:2d:46:e1:19:3d:98:e2:ff:
        b1:b6:fb:2d:cb:6e:a5:8b:40:87:23:22:96:5b:7b:15:89:bc:
        e4:65:ca:0f:a2:11:b2:06:9e:8d:c9:72:26:bb:94:a2:e2:87:
        6e:a5:af:4e:16:ff:b0:c7:9f:5f:49:e8:a6:ff:09:e5:eb:8c:
        31:e1:28:62:ce:23:bb:38:2e:e1:3e:d9:81:52:7a:6c:c1:56:
        be:2d:fe:0d:b8:c3:d6:41:e2:f2:a3:37:d5:fd:8b:98:5c:e9:
        ce:ae:43:8e:6b:36:54:3c:50:ae:67:00:1c:90:fd:35:83:3b:
        1c:86
```

需要注意，用户自行生成私钥情况下，私钥文件一旦丢失，CA 方由于不持有私钥信息，无法进行恢复，意味着通过该证书中公钥加密的内容将无法被解密。

### 证书的撤销

证书超出有效期后会作废，用户也可以主动向 CA 申请撤销某证书文件。

由于 CA 无法强制收回已经颁发出去的数字证书，因此为了实现证书的作废，往往还需要维护一个撤销证书列表（Certificate Revocation List，CRL），用于记录已经撤销的证书序号。

因此，通常情况下，当第三方对某个证书进行验证时，需要首先检查该证书是否在撤销列表中。如果存在，则该证书无法通过验证。如果不在，则继续进行后续的证书验证过程。
