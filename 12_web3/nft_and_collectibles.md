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
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

interface IERC165 {
    function supportsInterface(bytes4 interfaceId) external view returns (bool);
}

// ERC-721 核心接口。name/symbol/tokenURI 属于 IERC721Metadata 扩展。
interface IERC721 is IERC165 {
    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
    event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId);
    event ApprovalForAll(address indexed owner, address indexed operator, bool approved);

    function balanceOf(address owner) external view returns (uint256 balance);
    function ownerOf(uint256 tokenId) external view returns (address owner);
    function approve(address to, uint256 tokenId) external;
    function getApproved(uint256 tokenId) external view returns (address operator);
    function setApprovalForAll(address operator, bool approved) external;
    function isApprovedForAll(address owner, address operator) external view returns (bool);
    function transferFrom(address from, address to, uint256 tokenId) external;
    function safeTransferFrom(address from, address to, uint256 tokenId) external;
    function safeTransferFrom(address from, address to, uint256 tokenId, bytes calldata data) external;
}

interface IERC721Metadata is IERC721 {
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function tokenURI(uint256 tokenId) external view returns (string memory);
}
```

**安全最小实现示例**：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract SimpleCollectible is ERC721, Ownable {
    uint256 private nextTokenId = 1;
    string private baseTokenURI;

    constructor(string memory baseURI_) ERC721("Simple Collectible", "SCOL") Ownable(msg.sender) {
        baseTokenURI = baseURI_;
    }

    function safeMint(address to) external onlyOwner returns (uint256 tokenId) {
        tokenId = nextTokenId++;
        _safeMint(to, tokenId);
    }

    function _baseURI() internal view override returns (string memory) {
        return baseTokenURI;
    }
}
```

#### 2.2 ERC-1155 标准

支持同时发行可交换和不可交换代币，更灵活且 Gas 高效。

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

// ERC-1155 核心接口节选；实际开发应直接使用 OpenZeppelin 的 IERC1155/ERC1155。
interface IERC1155 {
    event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value);
    event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values);
    event ApprovalForAll(address indexed account, address indexed operator, bool approved);

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

#### 3.2 使用链外存储保存元数据

```javascript
// 伪代码：使用当前存储服务的官方 SDK 或 HTTP API。
// NFT.Storage Classic 的免费上传/API/SDK 已停止接受新上传，不应作为新项目模板。
async function uploadMetadata(storage, imageFile) {
  const imageCid = await storage.uploadFile(imageFile);

  const metadata = {
    name: "My NFT",
    description: "A beautiful digital art",
    image: `ipfs://${imageCid}`,
    attributes: [
      { trait_type: "Color", value: "Blue" },
      { trait_type: "Rarity", value: "Rare" }
    ]
  };

  const metadataCid = await storage.uploadJson(metadata);
  return `ipfs://${metadataCid}`;
}
```

关键点是把合约中的 `tokenURI` 指向内容寻址的元数据 URI，并为生产项目准备多个存储/固定副本，避免单一网关或单一服务影响可用性。

### 4. NFT 应用场景

#### 4.1 数字艺术品与收藏品

OpenSea、Blur 等平台使开发者和艺术家无需编写代码即可创建和交易 NFT。

**特点**：
- 高度个性化和稀有性
- 社区驱动的价值发现
- 版税机制（艺术家可在每次转售中获得分成）

**代码示例 - 声明版税的 NFT**：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {ERC721Royalty} from "@openzeppelin/contracts/token/ERC721/extensions/ERC721Royalty.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract RoyaltyCollectible is ERC721Royalty, Ownable {
    uint256 private nextTokenId = 1;

    constructor(address royaltyReceiver)
        ERC721("Royalty Collectible", "RCL")
        Ownable(msg.sender)
    {
        // 500 = 5%，OpenZeppelin ERC2981 默认分母为 10000。
        _setDefaultRoyalty(royaltyReceiver, 500);
    }

    function safeMint(address to) external onlyOwner returns (uint256 tokenId) {
        tokenId = nextTokenId++;
        _safeMint(to, tokenId);
    }
}
```

ERC-2981 只负责让市场通过 `royaltyInfo(tokenId, salePrice)` 查询收款地址和金额；NFT 合约不应暴露自定义销售回调或在合约内直接转账。是否支付、如何支付应由支持该标准的市场在交易结算中处理。

#### 4.2 游戏资产

链上游戏（GameFi）使用 NFT 表示游戏内物品，玩家真正拥有这些资产。

**示例架构**：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC1155} from "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract GameAsset is ERC1155, Ownable {
    enum AssetType { Weapon, Armor, Consumable }

    uint256 private nextAssetId = 1;
    mapping(uint256 => AssetType) public assetTypes;
    mapping(uint256 => uint256) public assetStats;  // 攻击力、防御力等

    constructor(string memory baseURI) ERC1155(baseURI) Ownable(msg.sender) {}

    function mintWeapon(address to, uint256 damage) external onlyOwner returns (uint256 assetId) {
        assetId = nextAssetId++;
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
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract SimpleDomainNFT is ERC721, Ownable {
    uint256 private nextTokenId = 1;
    mapping(string => uint256) public domainToTokenId;
    mapping(uint256 => string) public tokenDomains;

    event DomainRegistered(string indexed domain, address indexed owner);

    constructor() ERC721("Simple Domain", "SDOM") Ownable(msg.sender) {}

    function registerDomain(string memory domain, address to) external onlyOwner returns (uint256 tokenId) {
        require(bytes(domain).length > 0, "Invalid domain");
        require(domainToTokenId[domain] == 0, "Domain taken");

        tokenId = nextTokenId++;
        domainToTokenId[domain] = tokenId;
        tokenDomains[tokenId] = domain;
        _safeMint(to, tokenId);

        emit DomainRegistered(domain, to);
    }

    function resolveDomain(string memory domain) public view returns (address) {
        uint256 tokenId = domainToTokenId[domain];
        require(tokenId != 0, "Domain not registered");
        return ownerOf(tokenId);
    }
}
```

### 5. NFT 市场合约

二级市场允许持有者之间交易 NFT。

**简单的 NFT 市场实现**：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract NFTMarketplace is ReentrancyGuard {
    struct Listing {
        address seller;
        uint256 price;
    }

    IERC721 public immutable nft;
    mapping(uint256 => Listing) public listings;
    mapping(address => uint256) public proceeds;

    event Listed(uint256 indexed tokenId, address indexed seller, uint256 price);
    event Sold(uint256 indexed tokenId, address indexed seller, address indexed buyer, uint256 price);
    event Unlisted(uint256 indexed tokenId);
    event ProceedsWithdrawn(address indexed seller, uint256 amount);

    constructor(IERC721 nft_) {
        nft = nft_;
    }

    // 卖家上架 NFT
    function list(uint256 tokenId, uint256 price) external {
        require(price > 0, "Invalid price");
        require(nft.ownerOf(tokenId) == msg.sender, "Not owner");
        require(
            nft.getApproved(tokenId) == address(this) ||
                nft.isApprovedForAll(msg.sender, address(this)),
            "Marketplace not approved"
        );

        listings[tokenId] = Listing({
            seller: msg.sender,
            price: price
        });

        emit Listed(tokenId, msg.sender, price);
    }

    // 买家购买 NFT：检查 -> 更新状态 -> 外部交互
    function buy(uint256 tokenId) external payable nonReentrant {
        Listing memory listing = listings[tokenId];
        require(listing.seller != address(0), "Not listed");
        require(msg.value == listing.price, "Wrong payment");

        delete listings[tokenId];
        proceeds[listing.seller] += msg.value;

        nft.safeTransferFrom(listing.seller, msg.sender, tokenId);

        emit Sold(tokenId, listing.seller, msg.sender, listing.price);
    }

    // 卖家撤销上架
    function unlist(uint256 tokenId) external {
        require(listings[tokenId].seller == msg.sender, "Not seller");
        delete listings[tokenId];
        emit Unlisted(tokenId);
    }

    // 提款模式避免购买流程中直接向卖家转账
    function withdrawProceeds() external nonReentrant {
        uint256 amount = proceeds[msg.sender];
        require(amount > 0, "No proceeds");

        proceeds[msg.sender] = 0;
        (bool ok, ) = payable(msg.sender).call{value: amount}("");
        require(ok, "Withdraw failed");

        emit ProceedsWithdrawn(msg.sender, amount);
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
   // SPDX-License-Identifier: MIT
   pragma solidity ^0.8.24;

   import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";

   contract MyNFT is ERC721 {
       constructor() ERC721("MyNFT", "MNFT") {}
   }
   ```

2. **元数据持久化**：使用内容寻址 URI，并准备多节点 pinning、Filecoin、Arweave 等备份，避免依赖单一网关或单一服务

3. **版税支持**：实现 ERC-2981 标准以获得市场支持

4. **提供升级机制**：使用代理模式以应对安全发现

这样的深入理解 NFT 技术和应用，将为进一步的区块链应用开发打下坚实基础。
