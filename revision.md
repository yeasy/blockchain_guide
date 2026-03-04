## Version History

* 2.2.4: 2026-03-03
  * Update to new version;
  * Add latest progress;
  * Support BFT consensus algorithm;
  * Add more details.

* 2.0.0: 2026-01-12
  * Major update: Comprehensive revision of content to adapt to the technical landscape of 2026;
  * Add: Chapter 14 frontier trends (AI, DePIN, RWA);
  * Update: Bitcoin (Taproot, Ordinals, ETF) and Ethereum (PoS, Layer 2, Account Abstraction);
  * Update: Cryptography (Post-Quantum, ZK-SNARKs/STARKs) and Consensus Algorithms (PoS, HotStuff);
  * Optimize: Remove obsolete EOS chapter, integrate Hyperledger project status.

* 1.8.6: 2026-03-04
  * Fix terminology "miner" -> "validator" on line 19 of `07_ethereum/design.md` (context after Merge);
  * Add full Contract API example to `11_app_dev/chaincode.md` (recommended over shim API);
  * Add Fabric Gateway chapter and Go language code examples to `10_fabric_op/sdk.md`;
  * Update REST API deprecation notice in `11_app_dev/chaincode_example01.md` (v0.6 only);
  * Update ordering node section in `13_fabric_design/design.md`: add Fabric v3.0 BFT ordering service description;
  * Update BFT ordering performance characteristics in `13_fabric_design/performance.md`;
  * Fix various external links: Burrow archive, gRPC/Go official website HTTPS upgrade, Remix HTTPS, Cello Git repo migration, Kafka document link correction;
  * Fix deprecated `grpc.WithInsecure()` -> `insecure.NewCredentials()` modern pattern in `appendix/grpc.md`;
  * Correct Vyper positioning description in `07_ethereum/concept.md` (security alternative language rather than mainstream), fix double periods;
  * Update ZKML limitations description in `15_frontiers/ai_web3.md`, reflecting 2025 technical progress;
  * Correct modular blockchain date format in `02_overview/classify.md`.

* 1.8.5: 2026-03-03
  * Comprehensive technical review by experts, fix 8 P0 critical errors and 8 P1 important issues;
  * Fix mining reward data (12.5 -> 3.125 BTC) in `06_bitcoin/mining.md`, add halving history and hash rate data;
  * Fix PoS consensus context in `07_ethereum/concept.md`, remove outdated "mining" expressions and ETH price references;
  * Fix duplicate rows in comparison table and UTXO spelling in `07_ethereum/design.md`;
  * Update Solidity ^0.8.0 syntax and Geth --http command in `07_ethereum/smart_contract.md`;
  * Fix Prometheus "push" -> "pull" method in `10_fabric_op/operation.md`;
  * Fix Caliper/Grid URLs, Fabric 3.0 official release status in `08_hyperledger/project.md`;
  * Fix Org2 anchor node variable name error in `09_fabric_deploy/start_local.md`;
  * Fix import path error in `appendix/golang/tools.md`, update golint -> golangci-lint;
  * Update protobuf import path and protoc command migration in `appendix/grpc.md`;
  * Add Tendermint BFT chapter and six consensus algorithms comparison table in `04_distributed_system/bft.md`;
  * Add new DeFi attacks (flash loans, MEV, cross-chain bridges, governance attacks) to `05_crypto/smart_contract_vulns.md`;
  * Add MEV chapter in `12_web3/defi.md` and governance attack case in `12_web3/dao.md`;
  * Add Web3 challenges and reality chapter in `12_web3/intro.md`;
  * Add ZKML limitations analysis in `15_frontiers/ai_web3.md`;
  * Add sustainability risk analysis in `15_frontiers/depin.md`;
  * Add core regulatory framework checklist (Reg D, MiFID II, MiCA, etc.) in `15_frontiers/rwa.md`;
  * Add modular blockchain classification in `02_overview/classify.md`;
  * Update 12 modern terms in `appendix/terms.md`;
  * Update Docker version recommendations (18.03 -> 24.0+), Kafka deprecation notes, Layer 2 performance context;
  * Update USDT/USDC market cap data, add time labels;
  * Improve Cello image version warning, Azure service shutdown notice.

* 1.8.4: 2026-01-13
  * Comprehensive editorial review for publication readiness;
  * Fix Chinese monetary history TBD marker in `06_bitcoin/currency.md`;
  * Fix Kafka migration TBD marker in `10_fabric_op/upgrade.md`, noting Kafka is deprecated;
  * Verify all chapter content relevance (all updated to 2025).

* 1.8.3: 2026-01-12
  * Comprehensive editorial review for publication readiness;
  * Fix structural issues: create missing `10_fabric_op/node.md`, `12_web3/summary.md`;
  * Fix SUMMARY.md links for golang appendix and add Web3 summary;
  * Update `07_ethereum/design.md`: PoW→PoS consensus, Rollup-centric scalability;
  * Rewrite `07_ethereum/install.md` for modern Geth installation;
  * Update Go/Docker versions in `08_hyperledger/contribute.md`;
  * Remove TODO markers from `03_scenario/finance.md` and `03_scenario/others.md`;
  * Update outdated scenarios (2026 context): TradeLens (discontinued), JPM Coin (Onyx), IoT (DePIN);
  * Verify all 110+ image references (relative paths in `_images/`) and 120+ external links.

* 1.8.2: 2026-01-11
  * Add new chapters: Layer 2, CBDC (with global overview), Web3/DeFi/DAO, Zero-Knowledge Proofs;
  * Rewrite NFT chapter with 2022-2025 market updates (Ordinals, NFTFi, SBT);
  * Rewrite Fabric chaincode lifecycle for v2.x (Package-Install-Approve-Commit);
  * Update Ethereum tools (Hardhat, Foundry, MetaMask);
  * Update Go language guide (Go Modules, GOPROXY);
  * Update Fabric install guide (install-fabric.sh, LevelDB vs CouchDB);
  * Fix outdated data: Bitcoin block reward, SHA-3 status, DeFi TVL, pizza BTC value;
  * Add version notice to Fabric architecture design chapter;
  * Update appendix resource links.

* 1.8.0: 2026-01-10
  * Update Bitcoin history (ETF, halving) and stats;
  * Update Ethereum history (The Merge, Dencun upgrade) and features;
  * Update Hyperledger Fabric version (v2.5 LTS, v3.0) and project status;
  * Update BaaS platform info (IBM Blockchain Platform, Azure).

* 1.7.0: 2025-12-28
  * Update outdated project status;
  * Fix typos and links.

* 1.6.0: 2021-12-01
  * Fix expressions;
  * Fix typos.

* 1.5.0: 2021-01-21
  * Add operation chapter;
  * Fix typos and polish expression.

* 1.4.0: 2020-06-18
  * Refine deployment fabric with v2.0 version;
  * Update hyperledger community and projects;
  * Add operation guide and best practices.

* 1.3.0: 2019-12-31
  * Add more crypto techniques;
  * Update go and related tools;
  * Update bitcoin project.

* 1.2.0: 2018-12-31
  * Add common Golang tools and tips;
  * Update cryptography related knowledge, add bloom filters etc;
  * Update content of Hyperledger projects;
  * Update distributed system chapter.

* 1.1.0: 2018-04-24
  * Update group signature;
  * Update evolution of blockchain and distributed ledgers;
  * Update latest progress of Bitcoin and Ethereum.

* 1.0.0: 2017-12-31
  * Update BaaS design;
  * Update appendix section;
  * Correct some expressions.

* 0.9.0: 2017-08-24
  * Correct wording;
  * Add content for fabric 1.0;
  * "Blockchain Principles, Design and Application" officially published.

* 0.8.0: 2017-03-07
  * Improve application scenarios etc;
  * Improve distributed system technologies;
  * Improve cryptography technologies;
  * Update Hyperledger usage according to latest code.

* 0.7.0: 2016-09-10
  * Improve consensus technologies etc;
  * Correct wording.

* 0.6.0: 2016-08-05
  * Modify wording;
  * Add more smart contracts;
  * Add more business scenarios.

* 0.5.0: 2016-07-10
  * Add content for Hyperledger project;
  * Add content for Ethereum project;
  * Add Lightning Network introduction and key technology analysis;
  * Add Blockchain as a Service (BaaS);
  * Add Bitcoin project.

* 0.4.0: 2016-06-02
  * Add application scenario analysis.

* 0.3.0: 2016-05-12
  * Add digital currency issue analysis.

* 0.2.0: 2016-04-07
  * Add Hyperledger project introduction.

* 0.1.0: 2016-01-17
  * Add blockchain introduction.
