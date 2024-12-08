# observer

計測機器。
ansibleを使ってデプロイできる。ansibleの設定に関しては[ansibleのREADME](./ansible/README.md)にある。

それぞれのサービスがdocker composeで動いている。

## monitor (Grafana + Prometheus + Mimir + Loki)

- Grafana
  - メトリクスの可視化
  - ログの可視化
  - https://grafana.com/oss/grafana/
- Prometheus
  - メトリクスの集計
  - https://prometheus.io/
- Mimir
  - メトリクスの収集・蓄積
  - https://grafana.com/oss/mimir/
  - バックエンドとしてMinIOを、ロードバランサとしてNginxのコンテナを使っている。
- Loki
  - ログの収集・蓄積
  - https://grafana.com/oss/loki/

これらが1つのcompose.ymlで記述されている。

## pprotein

[kaz/pprotein](https://github.com/kaz/pprotein)をフォークして改造した[reactive-futakotamagawa/pprotein](https://github.com/reactive-futakotamagawa/pprotein)を使っている。具体的には、pt-query-digestが使えるようになっている。

GitHub Container Registry(ghcr.io)に Docker Image を配置して使っている。

## adminer

MySQLをブラウザから操作する。

<!-- ![](observer.svg) -->

## ポートフォワーディングについて

pprotein や adminer を使うためには、問題サーバーのデフォルトで空いている以外のポートを使う必要がある。問題サーバーのポートを空けると整合性チェックで落ちる可能性があるので、ポートフォワーディング用のコンテナを問題サーバーの数だけdocker composeに含めて、そこを通して操作を行うようになっている。

ポートフォワーディング用のコンテナは GitHub Container Registry に配置している。
https://ghcr.io/reactive-futakotamagawa/isu-isu-h/tunnel
