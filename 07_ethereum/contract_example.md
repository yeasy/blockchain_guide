## 智能合约案例：投票

本节将介绍一个用 Solidity 语言编写的智能合约案例。代码来源于 [Solidity 官方文档](https://solidity.readthedocs.io/en/latest/index.html) 中的示例。

该智能合约实现了一个自动化的、透明的投票应用。投票发起人可以发起投票，将投票权赋予投票人；投票人可以自己投票，或将自己的票委托给其他投票人；任何人都可以公开查询投票的结果。

### 智能合约代码

实现上述功能的合约代码如下所示，并不复杂，语法跟 JavaScript 十分类似。

```js
pragma solidity ^0.4.11;

contract Ballot {
    struct Voter {
        uint weight;
        bool voted;
        address delegate;
        uint vote;
    }

    struct Proposal {
        bytes32 name;
        uint voteCount;
    }

    address public chairperson;
    mapping(address => Voter) public voters;
    Proposal[] public proposals;

    // Create a new ballot to choose one of `proposalNames`
    function Ballot(bytes32[] proposalNames) {
        chairperson = msg.sender;
        voters[chairperson].weight = 1;

        for (uint i = 0; i < proposalNames.length; i++) {
            proposals.push(Proposal({
                name: proposalNames[i],
                voteCount: 0
            }));
        }
    }

    // Give `voter` the right to vote on this ballot.
    // May only be called by `chairperson`.
    function giveRightToVote(address voter) {
        require((msg.sender == chairperson) && !voters[voter].voted);
        voters[voter].weight = 1;
    }

    // Delegate your vote to the voter `to`.
    function delegate(address to) {
        Voter sender = voters[msg.sender];
        require(!sender.voted);
        require(to != msg.sender);

        while (voters[to].delegate != address(0)) {
            to = voters[to].delegate;

            // We found a loop in the delegation, not allowed.
            require(to != msg.sender);
        }

        sender.voted = true;
        sender.delegate = to;
        Voter delegate = voters[to];
        if (delegate.voted) {
            proposals[delegate.vote].voteCount += sender.weight;
        } else {
            delegate.weight += sender.weight;
        }
    }

    // Give your vote (including votes delegated to you)
    // to proposal `proposals[proposal].name`.
    function vote(uint proposal) {
        Voter sender = voters[msg.sender];
        require(!sender.voted);
        sender.voted = true;
        sender.vote = proposal;

        proposals[proposal].voteCount += sender.weight;
    }

    // @dev Computes the winning proposal taking all
    // previous votes into account.
    function winningProposal() constant
            returns (uint winningProposal)
    {
        uint winningVoteCount = 0;
        for (uint p = 0; p < proposals.length; p++) {
            if (proposals[p].voteCount > winningVoteCount) {
                winningVoteCount = proposals[p].voteCount;
                winningProposal = p;
            }
        }
    }

    // Calls winningProposal() function to get the index
    // of the winner contained in the proposals array and then
    // returns the name of the winner
    function winnerName() constant
            returns (bytes32 winnerName)
    {
        winnerName = proposals[winningProposal()].name;
    }
}
```

### 代码解析

#### 指定版本

在第一行，`pragma` 关键字指定了和该合约兼容的编译器版本。

```
pragma solidity ^0.4.11;
```

该合约指定，不兼容比 `0.4.11` 更旧的编译器版本，且 `^` 符号表示也不兼容从 `0.5.0` 起的新编译器版本。即兼容版本范围是 `0.4.11 <= version < 0.5.0`。该语法与 npm 的版本描述语法一致。

#### 结构体类型

Solidity 中的合约（contract）类似面向对象编程语言中的类。每个合约可以包含状态变量、函数、事件、结构体类型和枚举类型等。一个合约也可以继承另一个合约。

在本例命名为 `Ballot` 的合约中，声明了 2 个结构体类型：`Voter` 和 `Proposal`。

* `struct Voter`：投票人，其属性包括 `uint weight`（该投票人的权重）、`bool voted`（是否已投票）、`address delegate`（如果该投票人将投票委托给他人，则记录受委托人的账户地址）和 `uint vote`（投票做出的选择，即相应提案的索引号）。
* `struct Proposal`：提案，其属性包括 `bytes32 name`（名称）和 `uint voteCount`（已获得的票数）。

需要注意，`address` 类型记录了一个以太坊账户的地址。`address` 可看作一个数值类型，但也包括一些与以太币相关的方法，如查询余额 `<address>.balance`、向该地址转账 `<address>.transfer(uint256 amount)` 等。

#### 状态变量

合约中的状态变量会长期保存在区块链中。通过调用合约中的函数，这些状态变量可以被读取和改写。

本例中定义了 3 个状态变量：`chairperson`、`voters`、`proposals`。

* `address public chairperson`：投票发起人，类型为 `address`。
* `mapping(address => Voter) public voters`：所有投票人，类型为 `address` 到 `Voter` 的映射。
* `Proposal[] public proposals`：所有提案，类型为动态大小的 `Proposal` 数组。

3 个状态变量都使用了 `public` 关键字，使得变量可以被外部访问（即通过消息调用）。事实上，编译器会自动为 `public` 的变量创建同名的 getter 函数，供外部直接读取。

状态变量还可设置为 `internal` 或 `private`。`internal` 的状态变量只能被该合约和继承该合约的子合约访问，`private` 的状态变量只能被该合约访问。状态变量默认为 `internal`。

将上述关键状态信息设置为 `public` 能够增加投票的公平性和透明性。

#### 函数

合约中的函数用于处理业务逻辑。函数的可见性默认为 `public`，即可以从内部或外部调用，是合约的对外接口。函数可见性也可设置为 `external`、`internal` 和 `private`。

本例实现了 6 个 `public` 函数，可看作 6 个对外接口，功能分别如下。

##### 创建投票

函数 `function Ballot(bytes32[] proposalNames)` 用于创建一个新的投票。

所有提案的名称通过参数 `bytes32[] proposalNames` 传入，逐个记录到状态变量 `proposals` 中。同时用 `msg.sender` 获取当前调用消息的发送者的地址，记录为投票发起人 `chairperson`，该发起人投票权重设为 1。

##### 赋予投票权

函数 `function giveRightToVote(address voter)` 实现给投票人赋予投票权。

该函数给 `address voter` 赋予投票权，即将 `voter` 的投票权重设为 1，存入 `voters` 状态变量。

这个函数只有投票发起人 `chairperson` 可以调用。这里用到了 `require((msg.sender == chairperson) && !voters[voter].voted)` 函数。如果 `require` 中表达式结果为 `false`，这次调用会中止，且回滚所有状态和以太币余额的改变到调用前。但已消耗的 Gas 不会返还。

##### 委托投票权

函数 `function delegate(address to)` 把投票委托给其他投票人。

其中，用 `voters[msg.sender]` 获取委托人，即此次调用的发起人。用 `require` 确保发起人没有投过票，且不是委托给自己。由于被委托人也可能已将投票委托出去，所以接下来，用 `while` 循环查找最终的投票代表。找到后，如果投票代表已投票，则将委托人的权重加到所投的提案上；如果投票代表还未投票，则将委托人的权重加到代表的权重上。

该函数使用了 `while` 循环，这里合约编写者需要十分谨慎，防止调用者消耗过多 Gas，甚至出现死循环。

##### 进行投票

函数 `function vote(uint proposal)` 实现投票过程。

其中，用 `voters[msg.sender]` 获取投票人，即此次调用的发起人。接下来检查是否是重复投票，如果不是，进行投票后相关状态变量的更新。

##### 查询获胜提案

函数 `function winningProposal() constant returns (uint winningProposal)` 将返回获胜提案的索引号。

这里，`returns (uint winningProposal)` 指定了函数的返回值类型，`constant` 表示该函数不会改变合约状态变量的值。

函数通过遍历所有提案进行记票，得到获胜提案。

##### 查询获胜者名称

函数 `function winnerName() constant returns (bytes32 winnerName)` 实现返回获胜者的名称。

这里采用内部调用 `winningProposal()` 函数的方式获得获胜提案。如果需要采用外部调用，则需要写为 `this.winningProposal()`。
