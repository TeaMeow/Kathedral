# Kathedral

卡西卓是一個 Telegram 機器人，能夠將接收、被標記到的圖片儲存至本地端並且供他人在瀏覽器上觀賞（這將會啟動 Golang 內建檔案伺服器）。

## 螢幕截圖

## 相依性
無。

## 使用方法

```
名稱：
   Kathedral - 啟動 Kathrdral 常駐程式。
使用方式：
   kathedral [全域選項] 指令 [指令選項] [參數...]

版本：
   1.0.0

指令：
     help, h  顯示 Kathedral 的說明與指令清單。

全域選項：
   --token 值     Telegram 的機器人 Token。
   --addr  值     檔案伺服器的網址，將會被用在連結按鈕上。（預設：“example.com”）
   --port  值     Kathedral 將會在此埠口部署 Golang 原生檔案伺服器以暴露取得到的圖片。（預設：“:8888”）
   --help, -h     顯示說明。
   --version, -v  顯示版本號碼。
```

## 編譯

此程式以 Golang 撰寫，請安裝 Go 以重新編譯。更多平台請參考[此列表](https://golang.org/doc/install/source#environment)。

```bash
# For macOS
GOOS=darwin GOARCH=386 go build -o kathedral.macos

# For Windows
GOOS=windows GOARCH=386 go build -o kathedral.exe

# For Linux
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o kathedral.linux
```
