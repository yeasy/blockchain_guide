## 数字证书

对于非对称加密算法和数字签名来说，很重要的步骤就是公钥的分发。理论上任何人都可以获取到公开的公钥。然而这个公钥文件有没有可能是伪造的呢？传输过程中有没有可能被篡改呢？一旦公钥自身出了问题，则整个建立在其上的的安全性将不复成立。

数字证书机制正是为了解决这个问题，它就像日常生活中的证书一样，可以确保所记录信息的合法性。比如证明某个公钥是某个实体（个人或组织）拥有，并且确保任何篡改都能被检测出来，从而实现对用户公钥的安全分发。

根据所保护公钥的用途，数字证书可以分为加密数字证书（Encryption Certificate）和签名验证数字证书（Signature Certificate）。前者往往用于保护用于加密用途的公钥；后者则保护用于签名用途的公钥。两种类型的公钥也可以同时放在同一证书中。

一般情况下，证书需要由证书认证机构（Certification Authority，CA）来进行签发和背书。权威的商业证书认证机构包括 DigiCert、GlobalSign、VeriSign 等。用户也可以自行搭建本地 CA 系统，在私有网络中进行使用。

### X.509 证书规范
一般的，一个数字证书内容可能包括证书域（证书的版本、序列号、签名算法类型、签发者信息、有效期、被签发主体、**签发的公开密钥**）、CA 对证书的签名算法和签名值等。

目前使用最广泛的标准为 ITU 和 ISO 联合制定的 X.509 的 v3 版本规范（RFC 5280），其中定义了如下证书信息域：

* 版本号（Version Number）：规范的版本号，目前为版本 3，值为 0x2；
* 序列号（Serial Number）：由 CA 维护的为它所颁发的每个证书分配的唯一的序列号，用来追踪和撤销证书。只要拥有签发者信息和序列号，就可以唯一标识一个证书。最大不能超过 20 个字节；
* 签名算法（Signature Algorithm）：数字签名所采用的算法，如 sha256WithRSAEncryption 或 ecdsa-with-SHA256；
* 颁发者（Issuer）：颁发证书单位的信息，如 “C=CN, ST=Beijing, L=Beijing, O=org.example.com, CN=ca.org.example.com”；
* 有效期（Validity）：证书的有效期限，包括起止时间（如 Not Before 2018-08-08-00-00UTC，Not After 2028-08-08-00-00UTC）；
* 被签发主体（Subject）：证书拥有者的标识信息（Distinguished Name），如 “C=CN, ST=Beijing, L=Beijing, CN=personA.org.example.com”；
* 主体的公钥信息（Subject Public Key Info）：所保护的公钥相关的信息；
    * 公钥算法（Public Key Algorithm）：公钥采用的算法；
    * 主体公钥（Subject Public Key）：公钥的内容；
* 颁发者唯一号（Issuer Unique Identifier，可选）：代表颁发者的唯一信息，仅 2、3 版本支持，可选；
* 主体唯一号（Subject Unique Identifier，可选）：代表拥有证书实体的唯一信息，仅 2、3 版本支持，可选；
* 扩展（Extensions，可选）：可选的一些扩展。可能包括：
    * Subject Key Identifier：实体的密钥标识符，区分实体的多对密钥；
    * Basic Constraints：一般指明该证书是否属于某个 CA；
    * Authority Key Identifier：颁发这个证书的颁发者的公钥标识符；
    * Authority Information Access：颁发相关的服务地址，如颁发者证书获取地址和吊销证书列表信息查询地址；
    * CRL Distribution Points：证书注销列表的发布地址；
    * Key Usage: 表明证书的用途或功能信息，如 Digital Signature、Key CertSign；
    * Subject Alternative Name：证书身份实体的别名，如该证书可以同样代表 *.org.example.com，org.example.com，*.example.com，example.com 身份等。

此外，证书的颁发者还需要对证书内容利用自己的私钥进行签名，以防止他人篡改证书内容。

### 证书格式

X.509 规范中一般推荐使用 PEM（Privacy Enhanced Mail）格式来存储证书相关的文件。证书文件的文件名后缀一般为 `.crt` 或 `.cer`，对应私钥文件的文件名后缀一般为 `.key`，证书请求文件的文件名后缀为 `.csr`。有时候也统一用 `.pem` 作为文件名后缀。

PEM 格式采用文本方式进行存储，一般包括首尾标记和内容块，内容块采用 base64 编码。

例如，一个示例证书文件的 PEM 格式如下所示。

```
-----BEGIN CERTIFICATE-----
MIICMzCCAdmgAwIBAgIQIhMiRzqkCljq3ZXnsl6EijAKBggqhkjOPQQDAjBmMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEUMBIGA1UEChMLZXhhbXBsZS5jb20xFDASBgNVBAMTC2V4YW1wbGUu
Y29tMB4XDTE3MDQyNTAzMzAzN1oXDTI3MDQyMzAzMzAzN1owZjELMAkGA1UEBhMC
VVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28x
FDASBgNVBAoTC2V4YW1wbGUuY29tMRQwEgYDVQQDEwtleGFtcGxlLmNvbTBZMBMG
ByqGSM49AgEGCCqGSM49AwEHA0IABCkIHZ3mJCEPbIbUdh/Kz3zWW1C9wxnZOwfy
yrhr6aHwWREW3ZpMWKUcbsYup5kbouBc2dvMFUgoPBoaFYJ9D0SjaTBnMA4GA1Ud
DwEB/wQEAwIBpjAZBgNVHSUEEjAQBgRVHSUABggrBgEFBQcDATAPBgNVHRMBAf8E
BTADAQH/MCkGA1UdDgQiBCBIA/DmemwTGibbGe8uWjt5hnlE63SUsXuNKO9iGEhV
qDAKBggqhkjOPQQDAgNIADBFAiEAyoMO2BAQ3c9gBJOk1oSyXP70XRk4dTwXMF7q
R72ijLECIFKLANpgWFoMoo3W91uzJeUmnbJJt8Jlr00ByjurfAvv
-----END CERTIFICATE-----
```

可以通过 openssl 工具来查看其内容。

```bash
# openssl x509 -in example.com-cert.pem -noout -text
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            22:13:22:47:3a:a4:0a:58:ea:dd:95:e7:b2:5e:84:8a
    Signature Algorithm: ecdsa-with-SHA256
        Issuer: C=US, ST=California, L=San Francisco, O=example.com, CN=example.com
        Validity
            Not Before: Apr 25 03:30:37 2017 GMT
            Not After : Apr 23 03:30:37 2027 GMT
        Subject: C=US, ST=California, L=San Francisco, O=example.com, CN=example.com
        Subject Public Key Info:
            Public Key Algorithm: id-ecPublicKey
                Public-Key: (256 bit)
                pub:
                    04:29:08:1d:9d:e6:24:21:0f:6c:86:d4:76:1f:ca:
                    cf:7c:d6:5b:50:bd:c3:19:d9:3b:07:f2:ca:b8:6b:
                    e9:a1:f0:59:11:16:dd:9a:4c:58:a5:1c:6e:c6:2e:
                    a7:99:1b:a2:e0:5c:d9:db:cc:15:48:28:3c:1a:1a:
                    15:82:7d:0f:44
                ASN1 OID: prime256v1
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature, Key Encipherment, Certificate Sign, CRL Sign
            X509v3 Extended Key Usage:
                Any Extended Key Usage, TLS Web Server Authentication
            X509v3 Basic Constraints: critical
                CA:TRUE
            X509v3 Subject Key Identifier:
                48:03:F0:E6:7A:6C:13:1A:26:DB:19:EF:2E:5A:3B:79:86:79:44:EB:74:94:B1:7B:8D:28:EF:62:18:48:55:A8
    Signature Algorithm: ecdsa-with-SHA256
         30:45:02:21:00:ca:83:0e:d8:10:10:dd:cf:60:04:93:a4:d6:
         84:b2:5c:fe:f4:5d:19:38:75:3c:17:30:5e:ea:47:bd:a2:8c:
         b1:02:20:52:8b:00:da:60:58:5a:0c:a2:8d:d6:f7:5b:b3:25:
         e5:26:9d:b2:49:b7:c2:65:af:4d:01:ca:3b:ab:7c:0b:ef
```

此外，还有 DER（Distinguished Encoding Rules）格式，是采用二进制对证书进行保存，可以与 PEM 格式互相转换。

### 证书信任链

证书中记录了大量信息，其中最重要的包括 `签发的公开密钥` 和 `CA 数字签名` 两个信息。因此，只要使用 CA 的公钥再次对这个证书进行签名比对，就能证明所记录的公钥是否合法。

读者可能会想到，怎么证明用来验证对实体证书进行签名的 CA 公钥自身是否合法呢？毕竟在获取 CA 公钥的过程中，它也可能被篡改掉。

实际上，CA 的公钥是否合法，一方面可以通过更上层的 CA 颁发的证书来进行认证；另一方面某些根 CA（Root CA）可以通过预先分发证书来实现信任基础。例如，主流操作系统和浏览器里面，往往会提前预置一些权威 CA 的证书（通过自身的私钥签名，系统承认这些是合法的证书）。之后所有基于这些 CA 认证过的中间层 CA（Intermediate CA）和后继 CA 都会被验证合法。这样就从预先信任的根证书，经过中间层证书，到最底下的实体证书，构成一条完整的证书信任链。

某些时候用户在使用浏览器访问某些网站时，可能会被提示是否信任对方的证书。这说明该网站证书无法被当前系统中的证书信任链进行验证，需要进行额外检查。另外，当信任链上任一证书不可靠时，则依赖它的所有后继证书都将失去保障。

可见，证书作为公钥信任的基础，对其生命周期进行安全管理十分关键。后面章节将介绍的 PKI 体系提供了一套完整的证书管理的框架，包括生成、颁发、撤销过程等。
