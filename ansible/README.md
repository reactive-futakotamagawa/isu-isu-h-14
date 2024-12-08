# ISUCON用 ansible

## 必要

- Ansible
  - `init`を実行するためには、`community.general.git_config`、`community.general.github_deploy_key`、`prometheus.prometheus.node_exporter`のモジュールが必要。`ansible-galaxy collection install community.general`、`ansible-galaxy collection install prometheus.prometheus`でインストールできる。
- Python3系
  - Ansibleを動かすのに必要。

## 使い方

### 競技用のリポジトリを用意する

```sh
source ./bin/init.sh {ディレクトリ名} {サーバーの数}
```

で必要なディレクトリ構成ができる。

```txt
.
├── Makefile
├── s1
│   └── etc
│       ├── mysql
│       │   └── .gitkeep
│       ├── nginx
│       │   └── .gitkeep
│       └── systemd
│           └── system
│               └── .gitkeep
├── s2
│   └── etc
│       ├── mysql
│       │   └── .gitkeep
│       ├── nginx
│       │   └── .gitkeep
│       └── systemd
│           └── system
│               └── .gitkeep
└── s3
    └── etc
        ├── mysql
        │   └── .gitkeep
        ├── nginx
        │   └── .gitkeep
        └── systemd
            └── system
                └── .gitkeep
```

GitHubに上げておく。

### vaultのパスワードを配置する

`./vault.txt`を置き、パスワードを書く。

### 変数を設定する

通常

- [`group_vars/all/vars.yml`](./group_vars/all/vars.yml)
- [`hosts.yml`](./hosts.yml)

vault

- [`group_vars/all/vault.yml`](./group_vars/all/vault.yml)

```sh
ansible-vault edit group_vars/all/vault.yml
```

### 実行

#### 初期化

- ツールのインストール
- Gitの設定
- GitHubリポジトリのセットアップ(コミット・プッシュはしない)

```sh
ansible-playbook 0_init.yml
```

#### デプロイ

- git pull
- DB、nginx、アプリの再起動
- ログローテーション

デフォルトブランチ

```sh
ansible-playbook 1_deploy.yml
```

ブランチを指定

```sh
ansible-playbook 1_deploy.yml -e "deploy_branch={ブランチ名}"
```
