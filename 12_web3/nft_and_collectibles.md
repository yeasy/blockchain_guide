## NFT 与数字收藏品

### 概述

NFT（非同质化代币，Non-Fungible Token）是指具有唯一性和不可互换性的数字资产。与 ERC-20 代币（每个代币完全相同）不同，每个 NFT 都有独特的标识和属性，代表真实世界或数字世界中的唯一物品的所有权。

### 1. 与 ERC-20 的本质区别

| 特性 | ERC-20 代币 | NFT (ERC-721) |
|------|-----------|--------------|
| 同质性 | 完全同质，可互换 | 非同质，每个唯一 |
| 转账 | 任意数量转账 | 一次一个，不可分割 |
| 识别 | 通过数量识别 | 通过 Token ID 识别 |
| 应用场景 | 货币、支付 | 艺术品、收藏、游戏资产 |
| 智能合约标准 | ERC-20 | ERC-721, ERC-1155 |

### 2. NFT 的技术标准

#### 2.1 ERC-721 标准

最基础的 NFT 标准，每个代币通过唯一的 tokenId 识别。

```solidity
pragma solidity ^0.8.0;

interface IERC721 {
    // 返回账户拥有的 NFT 数量
    function balanceOf(address owner) external view returns (uint256 balance);

    // 返回 NFT 的所有者
    function ownerOf(uint256 tokenId) external view returns (address owner);

    // 转账 NFT
    function transferFrom(address from, address to, uint256 tokenId) external;

    // 授权第三方进行操作
    function approve(address to, uint256 tokenId) external;

    // 返回被授权的地址
    function getApproved(uint256 tokenId) external view returns (address operator);

    // 批量授权
    function setApprovalForAll(address operator, bool approved) external;

    // 检查是否被授权
    function isApprovedForAll(address owner, address operator) external view returns (bool);
}

// 事件
event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId);
event ApprovalForAll(address indexed owner, address indexed operator, bool approved);
```

**完整实现示例**：

```solidity
pragma solidity ^0.8.0;

contract SimpleNFT {
    string public name = "Simple NFT";
    string public symbol = "SNFT";

    // Token 元数据
    mapping(uint256 => string) private tokenURI;
    mapping(uint256 => address) private tokenOwner;
    mapping(address => uint256) private balances;
    mapping(uint256 => address) private tokenApprovals;
    mapping(address => mapping(address => bool)) private operatorApprovals;

    uint256 private tokenIdCounter;

    // 事件
    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
    event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId);
    event Mint(address indexed to, uint256 indexed tokenId, string uri);

    // 铸造 NFT
    function mint(address to, string memory uri) public returns (uint256) {
        require(to != address(0), "Invalid address");

        uint256 tokenId = tokenIdCounter++;
        tokenOwner[tokenId] = to;
        balances[to]++;

        tokenURI[tokenId] = uri;

        emit Mint(to, tokenId, uri);
        emit Transfer(address(0), to, tokenId);

        return tokenId;
    }

    // 销毁 NFT
    function burn(uint256 tokenId) public {
        require(msg.sender == tokenOwner[tokenId], "Not owner");

        address owner = tokenOwner[tokenId];
        balances[owner]--;
        delete tokenOwner[tokenId];
        delete tokenURI[tokenId];

        emit Transfer(owner, address(0), tokenId);
    }

    // 转账（带授权检查）
    function transferFrom(address from, address to, uint256 tokenId) public {
        require(from == tokenOwner[tokenId], "Not from owner");
        require(to != address(0), "Invalid address");
        require(
            msg.sender == from || msg.sender == tokenApprovals[tokenId] || operatorApprovals[from][msg.sender],
            "Not approved"
        );

        // 转移所有权
        tokenOwner[tokenId] = to;
        balances[from]--;
        balances[to]++;

        // 清除单次授权
        if (tokenApprovals[tokenId] != address(0)) {
            delete tokenApprovals[tokenId];
        }

        emit Transfer(from, to, tokenId);
    }

    // 获取 Token 元数据 URI
    function getTokenURI(uint256 tokenId) public view returns (string memory) {
        require(tokenOwner[tokenId] != address(0), "Token does not exist");
        return tokenURI[tokenId];
    }

    // 获取余额
    function balanceOf(address owner) public view returns (uint256) {
        return balances[owner];
    }

    // 获取所有者
    function ownerOf(uint256 tokenId) public view returns (address) {
        return tokenOwner[tokenId];
    }

    // 授权
    function approve(address to, uint256 tokenId) public {
        address owner = tokenOwner[tokenId];
        require(msg.sender == owner || operatorApprovals[owner][msg.sender], "Not authorized");

        tokenApprovals[tokenId] = to;
        emit Approval(owner, to, tokenId);
    }

    // 批量授权
    function setApprovalForAll(address operator, bool approved) public {
        operatorApprovals[msg.sender][operator] = approved;
        emit ApprovalForAll(msg.sender, operator, approved);
    }
}
```

#### 2.2 ERC-1155 标准

支持同时发行可交换和不可交换代币，更灵活且 Gas 高效。

```solidity
pragma solidity ^0.8.0;

interface IERC1155 {
    // 查询余额
    function balanceOf(address account, uint256 id) external view returns (uint256);

    // 批量查询余额
    function balanceOfBatch(address[] calldata accounts, uint256[] calldata ids)
        external view returns (uint256[] memory);

    // 转账
    function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes calldata data) external;

    // 批量转账
    function safeBatchTransferFrom(
        address from,
        address to,
        uint256[] calldata ids,
        uint256[] calldata amounts,
        bytes calldata data
    ) external;

    // 批量授权
    function setApprovalForAll(address operator, bool approved) external;

    // 检查授权
    function isApprovedForAll(address account, address operator) external view returns (bool);
}

event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value);
event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values);
event ApprovalForAll(address indexed account, address indexed operator, bool approved);
```

**ERC-1155 的优势**：
- 可同时管理多种资产（纯 NFT、可交换代币、半可交换代币）
- 单次交易可处理多个 ID，Gas 费用更低
- 适合游戏中各类资源的管理（武器、皮肤、货币等）

### 3. NFT 元数据标准

NFT 的真实内容（图片、视频、文本）通常不存储在区块链上，而是存储在链外存储（如 IPFS），智能合约仅保存指向这些资源的 URI。

#### 3.1 标准元数据格式

ERC-721 使用 JSON 格式描述元数据：

```json
{
  "name": "Bored Ape #1",
  "description": "A rare digital ape from the Bored Ape Yacht Club collection",
  "image": "ipfs://QmRZeF3mAL4VY5N8w8Wz8H6k1m2Zx4c7v8b9d0e1f2g3h4i",
  "attributes": [
    {
      "trait_type": "Background",
      "value": "Blue"
    },
    {
      "trait_type": "Fur",
      "value": "Tan"
    },
    {
      "trait_type": "Eyes",
      "value": "Blue Beams"
    },
    {
      "trait_type": "Mouth",
      "value": "Bored"
    },
    {
      "trait_type": "Rarity Score",
      "value": 8.5
    }
  ]
}
```

#### 3.2 使用 IPFS 存储元数据

```javascript
// 使用 web3.storage 或 NFT.storage 上传到 IPFS
const NFTStorage = require('nft.storage');
const fs = require('fs');

async function uploadMetadata() {
    const client = new NFTStorage({ token: process.env.NFT_STORAGE_TOKEN });

    const metadata = {
        name: "My NFT",
        description: "A beautiful digital art",
        image: new File([fs.readFileSync('art.jpg')], 'art.jpg', { type: 'image/jpeg' }),
        attributes: [
            { trait_type: "Color", value: "Blue" },
            { trait_type: "Rarity", value: "Rare" }
        ]
    };

    // 上传并获取 IPFS CID
    const cid = await client.store(metadata);
    const uri = `ipfs://${cid}`;

    console.log(`元数据已上传到: ${uri}`);
    return uri;
}
```

### 4. NFT 应用场景

#### 4.1 数字艺术品与收藏品

OpenSea、Blur 等平台使开发者和艺术家无需编写代码即可创建和交易 NFT。

**特点**：
- 高度个性化和稀有性
- 社区驱动的价值发现
- 版税机制（艺术家可在每次转售中获得分成）

**代码示例 - 支持版税的 NFT**：

```solidity
pragma solidity ^0.8.0;

contract RoyaltyNFT {
    mapping(uint256 => address) public creators;
    mapping(uint256 => uint256) public royaltyPercent;  // 百分比，基数 10000

    event RoyaltyTransfer(address indexed recipient, uint256 amount);

    // 销售时调用，自动分配版税
    function handleSale(uint256 tokenId, uint256 salePrice) internal {
        address creator = creators[tokenId];
        if (creator != address(0)) {
            uint256 royalty = (salePrice * royaltyPercent[tokenId]) / 10000;
            payable(creator).transfer(royalty);
            emit RoyaltyTransfer(creator, royalty);
        }
    }

    // 市场合约会调用此函数
    function onSale(uint256 tokenId, uint256 price) public {
        handleSale(tokenId, price);
    }
}
```

#### 4.2 游戏资产

链上游戏（GameFi）使用 NFT 表示游戏内物品，玩家真正拥有这些资产。

**示例架构**：

```solidity
pragma solidity ^0.8.0;

contract GameAsset is ERC1155 {
    enum AssetType { Weapon, Armor, Consumable }

    mapping(uint256 => AssetType) public assetTypes;
    mapping(uint256 => uint256) public assetStats;  // 攻击力、防御力等

    function mintWeapon(address to, uint256 damage) public onlyGameMaster {
        uint256 assetId = totalAssets++;

        assetTypes[assetId] = AssetType.Weapon;
        assetStats[assetId] = damage;

        _mint(to, assetId, 1, "");
    }

    function upgradeAsset(uint256 assetId, uint256 increaseStats) public {
        require(balanceOf(msg.sender, assetId) > 0, "Not owner");
        assetStats[assetId] += increaseStats;
    }
}
```

#### 4.3 身份与凭证

可验证的链上身份、学位证书、会员证等。

**优势**：
- 防伪性强（由区块链保证）
- 可验证性（任何人可查证真伪）
- 可转移性（持有者可出售或转赠）

#### 4.4 域名和地址

ENS（Ethereum Name Service）将复杂的钱包地址映射到易记的域名。

```solidity
pragma solidity ^0.8.0;

contract SimpleDomainNFT {
    mapping(string => address) public domainOwners;
    mapping(address => string) public addressToDomain;

    event DomainRegistered(string indexed domain, address indexed owner);

    function registerDomain(string memory domain, address to) public {
        require(bytes(domain).length > 0, "Invalid domain");
        require(domainOwners[domain] == address(0), "Domain taken");

        domainOwners[domain] = to;
        addressToDomain[to] = domain;

        emit DomainRegistered(domain, to);
    }

    function resolveDomain(string memory domain) public view returns (address) {
        return domainOwners[domain];
    }
}
```

### 5. NFT 市场合约

二级市场允许持有者之间交易 NFT。

**简单的 NFT 市场实现**：

```solidity
pragma solidity ^0.8.0;

interface IERC721 {
    function transferFrom(address from, address to, uint256 tokenId) external;
}

contract NFTMarketplace {
    struct Listing {
        address seller;
        uint256 price;
        bool active;
    }

    IERC721 public nftContract;
    mapping(uint256 => Listing) public listings;

    event Listed(uint256 indexed tokenId, address indexed seller, uint256 price);
    event Sold(uint256 indexed tokenId, address indexed seller, address indexed buyer, uint256 price);

    constructor(address _nftContract) {
        nftContract = IERC721(_nftContract);
    }

    // 卖家上架 NFT
    function list(uint256 tokenId, uint256 price) public {
        require(price > 0, "Invalid price");

        listings[tokenId] = Listing({
            seller: msg.sender,
            price: price,
            active: true
        });

        emit Listed(tokenId, msg.sender, price);
    }

    // 买家购买 NFT
    function buy(uint256 tokenId) public payable {
        Listing memory listing = listings[tokenId];
        require(listing.active, "Not listed");
        require(msg.value >= listing.price, "Insufficient payment");

        // 转移 NFT
        nftContract.transferFrom(listing.seller, msg.sender, tokenId);

        // 转移支付
        payable(listing.seller).transfer(msg.value);

        // 标记为已售
        listings[tokenId].active = false;

        emit Sold(tokenId, listing.seller, msg.sender, listing.price);
    }

    // 卖家撤销上架
    function unlist(uint256 tokenId) public {
        require(listings[tokenId].seller == msg.sender, "Not seller");
        listings[tokenId].active = false;
    }
}
```

### 6. NFT 的风险和挑战

**技术风险**：
- 智能合约漏洞可导致资产被盗
- 元数据链接失效（IPFS 节点关闭，中心化服务器关闭）
- 跨链桥接风险（跨链 NFT 时可能丢失）

**市场风险**：
- 价格波动剧烈，投机成分大
- 洗盘行为（人为抬高价格）
- 地板价格（底价）可能长期低迷

**法律风险**：
- 知识产权问题（NFT 不等于版权）
- 监管不确定性
- 税收处理复杂

### 7. 最佳实践

1. **使用经审计的 ERC-721 标准库**
   ```solidity
   import "@openzeppelin/contracts/token/ERC721/ERC721.sol";

   contract MyNFT is ERC721 {
       constructor() ERC721("MyNFT", "MNFT") {}
   }
   ```

2. **元数据持久化**：使用 Arweave（永久存储）而非 IPFS（节点可能关闭）

3. **版税支持**：实现 ERC-2981 标准以获得市场支持

4. **提供升级机制**：使用代理模式以应对安全发现

这样的深入理解 NFT 技术和应用，将为进一步的区块链应用开发打下坚实基础。
