# observer/ansible

observer をデプロイするための ansible。isu-isu-h の機能としての ansible は [../../ansible](../../ansible/README.md) を見てね。

## 使い方

1. `vault.txt`に ansible vault のパスワードを書く。
2. [`hosts`](./hosts)にデプロイしたいサーバーのIPアドレスやURLを書く。
3. [group_vars/all/vars.yml](./group_vars/all/vars.yml)を編集する。
4. `ansible-vault edit group_vars/all/vault.yml`を実行し、変数を編集する。
5. 以下を実行する。(デプロイ)

```sh
ansible-playbook main.yml
```

あんまりちゃんとデバッグしてないから、動かない環境ありそう。apt使ってるから少なくともUbuntuが必要？
