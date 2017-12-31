## Participate in contributions

Contributors [list](https://github.com/yeasy/blockchain_guide/graphs/contributors).

Blockchain technology itself is still rapid development, the ecological environment is also booming.

The open source source of the book is hosted on Github, welcome to participate in the maintenance: [github.com/yeasy/blockchain_guide](https://github.com/yeasy/blockchain_guide).

First of all, on GitHub `fork` to your own repositories, such as `some_user/blockchain_guide`, and then `clone` to the local, and set user information.

```sh
$ git clone git@github.com:docker_user/blockchain_guide.git
$ cd blockchain_guide
$ git config user.name "yourname"
$ git config user.email "your email"
```

Submitted after the update, and pushed to their own repository.

```sh
$ #do some change on the content
$ git commit -am "Fix issue #1: change helo to hello"
$ git push
```

Finally, submit your pull request on GitHub.

In addition, it is recommended to periodically update the content of your own repository using the contents of the project's repository.

```sh
$ git remote add upstream https://github.com/yeasy/blockchain_guide
$ git fetch upstream
$ git checkout master
$ git rebase upstream/master
$ git push -f origin master
```
