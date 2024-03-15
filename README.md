# Simple Recommend API

此專案滿足以下測驗需求 
1. **註冊**
    - **帳號要求：** 使用者必須註冊一個以 email 為帳號的用戶。
    - **密碼強度：** 密碼長度必須介於 6 到 16 個字符之間，並且至少包含一個大寫字母、一個小寫字母和一個特殊符號。

2. **驗證 Email**
    - 使用者註冊後，系統會發送一封包含驗證連結的 email 到他們的郵箱。
    - 使用者需要點擊該連結或輸入接收到的驗證碼以完成 email 驗證過程。

3. **登入**
    - 使用者可以使用其 email 和密碼進行登入。

4. **獲取推薦商品**
    - **接口：** `GET /recommendation`
    - **訪問限制：** 使用者必須登入後才能訪問此接口。
    - **功能描述：** 此接口回應用戶推薦的商品資料。

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- [Go](https://golang.org/dl/) (Version 1.22 or later)
- [Docker](https://www.docker.com/get-started) (運作開發環境需要的 MYSQL & Redis)
- [Migrate](https://github.com/golang-migrate/migrate) (用於 dev 環境快速建立 DB table)
- [Sqlc](https://github.com/sqlc-dev/sqlc) (A SQL Compiler, optional, 如果要動到 CRUD 的部分才要安裝)

### Running Locally

To run the SimpleRecommend API locally for debugging purposes, follow these steps:

1. **準備開發環境**
   確保您的開發環境中的 3306 (MySQL) 和 6379 (Redis) 端口未被佔用，因為我們將使用 Docker 來建立 MySQL 和 Redis 實例。請按照以下步驟使用我們在專案根目錄提供的 `docker-compose.yaml` 文件：
   ```sh
   docker-compose -p simplerecommend up -d
   ```

2. Download the Go modules required by the project:
    ```sh
    go mod download
    ```
   
3. 基於您已經安裝 [Migrate](https://github.com/golang-migrate/migrate) 和 啟動 MySQL&Redis
   ```sh
   migrate -path scripts/db/migration -database "mysql://user:password@tcp(localhost:3306)/simplerecommend?multiStatements=true" -verbose up
   ```
   
4. .env
   ```sh
   cp env.dev .env
   ```
   注意！！ EMAIL_SENDER_PASSWORD 不提供，如果要驗證寄信，請換上自己的 [AppPassword](https://myaccount.google.com/apppasswords)
   
5. 將開發用的 Postman collection 匯入 (節省檢查時間), 檔案位於
   ```
   ./scripts/postman/simplerecommend.postman_collection.json
   ```

6. Run the application:
    ```sh
    go run ./cmd/main.go
    ```
## 功能驗證(happy case 展示)
### happy case 系統 log
![Happy Log](./resultImage/createAccountFlow/happylog.png)
### 註冊成功 >> 寄信 >> 使用 code 驗證通過 >> login >> 存取 /recommendation
1. 剛註冊成功的時候屬於尚未驗證狀態
![Flow](./resultImage/createAccountFlow/account_table.png)
2. 收到信件(由於放上 github 會抽掉 EMAIL_SENDER_PASSWORD 如果要驗這部分要抽換自己的 EMAIL_SENDER_PASSWORD)
![Email](./resultImage/createAccountFlow/email.png)
3. 用 postman 模擬使用者點擊了連結, 可以參考 happy case 系統 log
4. 使用者帳號被刷上 is_valid true & verified_at
![Verify](./resultImage/createAccountFlow/verified_account.png)
5. 使用者 login 成功, 可以參考 happy case 系統 log
6. 使用者 access /recommendation, log 中可以看到一開始花了 2 秒(我開發時覺得等3秒覺得有點久就改成2秒), 之後靠 redis 變成微秒
