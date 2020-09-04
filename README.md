# StartKit

See [omo-msa-startkit](https://github.com/xtech-cloud/omo-msa-startkit)

# Protoc

See [omo-msp-account](https://github.com/xtech-cloud/omo-msp-account)

# Docker

See [omo-docker-account](https://github.com/xtech-cloud/omo-docker-account)

# 消息订阅

- 地址
  omo.msa.account.notification

- 消息
  | Action | Head | Body|
  |:--|:--|:--|
  |/signup||uuid|
  |/signin|accessToken|uuid|
  |/signout|accessToken||
  |/reset/password|accessToken||
  |/profile/update|accessToken||
  |/profile/query|accessToken||
