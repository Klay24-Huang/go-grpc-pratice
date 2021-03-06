# Golang 微服務範例

## 概要

這是一個練習專案，使用 Golang 的 gin 框架，實作一個簡短的 Restful Api server
使用後端微服務概念，每個服務之間利用 gRPC 溝通，並利用雲端空間 oKeteto 的 kubernetes 服務呈現作品

## 架構

![架構](https://i.imgur.com/wzTl4l0.png "架構")

## 在本地端運行

clone 並進入專案底下後

##### `1`. 使用 docker-compose

```
$docker-compose up
```

##### `2`. 使用 kubernetes

```
$ cd kubernetes
$ kubectl apply -f user-service.yaml
$ kubectl apply -f product-service.yaml
$ kubectl apply -f order-service.yaml
$ kubectl apply -f server-service.yaml
```

## 使用 swagger 測試

太長時間沒被訪問，此線上服務可能會進入沉睡狀態(會回覆 503 or 504)，請稍後再試試看
[swagger url](https://server-klay24-huang.cloud.okteto.net/swagger/index.html "link")

如果是用 docker compose 或 kubernetes 運行  
可以從 http://localhost:1231/swagger/index.html  
進行測試

## 使用 postman 測試

[下載 OpenApi json](https://drive.google.com/file/d/1CEKd1Pbq7R-WD376Ifvx2Ldgs754nIn5/view?usp=sharing "link")  
匯入至 postman 或其他測試工具後，即可使用線上 server 進行測試  
如需使用 local 端測試，請更改環境變數中 host 變數為 localhost:1231
