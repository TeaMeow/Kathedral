package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fmt"

	"github.com/codegangsta/cli"
	api "gopkg.in/telegram-BOT-api.v4"
)

var (
	DIR, TOKEN, PORT, ADDR string
	BOT                    *api.BotAPI
)

// fileServer 呼叫內建的檔案伺服器來提供瀏覽器在網頁上瀏覽檔案。
func fileServer() {
	http.Handle("/", http.FileServer(http.Dir(fmt.Sprintln("%s/files", DIR))))
	http.ListenAndServe(":"+PORT, nil)
}

// getImage 會取得遠端圖片並放置於本機位置。
func getImage(url, id string) error {
	// 建立空白檔案。
	img, _ := os.Create(fmt.Sprintf("%s/files/%s.jpg", DIR, id))
	defer img.Close()

	// 從遠端取得圖片。
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	// 將圖片複製於本機。
	io.Copy(img, resp.Body)
	return nil
}

// send 傳送基本的訊息到指定的頻道。
func send(u api.Update, m string) (api.Message, error) {
	msg := api.NewMessage(u.Message.Chat.ID, m)
	msg.ParseMode = "markdown"
	msg.ReplyToMessageID = u.Message.MessageID
	return BOT.Send(msg)
}

// bot 會啟動 Kathedral 並持續監聽 Telegram 訊息事件。
func bot(c *cli.Context) {
	// 配置全域變數。
	DIR, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	TOKEN = c.String("token")
	ADDR = c.String("addr")
	PORT = c.String("port")

	// 建立檔案資料夾如果不存在的話。
	newpath := filepath.Join(DIR, "files")
	os.MkdirAll(newpath, os.ModePerm)

	// 啟動 Golang 原生檔案伺服器。
	go fileServer()

	// 以 Token 建立新的 BOT API。
	BOT, _ = api.NewBotAPI(TOKEN)
	BOT.Debug = true

	log.Printf("連線已經成功驗證於 %s", BOT.Self.UserName)

	u := api.NewUpdate(0)
	u.Timeout = 60

	// 捕捉每個訊息的更新。
	updates, _ := BOT.GetUpdatesChan(u)
	for update := range updates {
		// 如果訊息是空的則略過。
		if update.Message == nil {
			continue
		}

		// 是否為 Mention（使用者以 `@` 標記 BOT）型態。
		var isMentioned bool
		if update.Message.Entities != nil {
			isMentioned = true
		}

		go func(u api.Update) {

			// 捕捉到了底層異常則回復並繼續執行，避免 BOT 因為 panic() 而死亡。
			defer func() {
				if err := recover(); err != nil {
					send(u, fmt.Sprintln(err))
				}
			}()

			// 決定圖片的來源：來自「使用者標記的訊息」或是「使用者自己傳送的訊息」。
			var photos *[]api.PhotoSize
			if isMentioned {
				photos = u.Message.ReplyToMessage.Photo
			} else {
				photos = u.Message.Photo
			}
			// 如果圖片是空的則略過。
			if photos == nil {
				return
			}

			// 發送「Fetching...」訊息讓使用者知道 BOT 接收到了資料。
			msg, _ := send(u, "_Fetching..._")
			// 尺寸文字清單，用以排列文字。
			list := ""
			// 取得開始時間。
			start := time.Now()

			// 按鈕陣列，稍後會將每個圖片尺寸作為按鈕然後推至此陣列。
			var buttons []api.InlineKeyboardButton
			// 讀取每個尺寸的圖片。
			for i, v := range *photos {
				// 此尺寸的圖片檔案編號。
				id := v.FileID

				// 取得圖片的存取網址。
				conf := api.FileConfig{FileID: id}
				file, _ := BOT.GetFile(conf)
				url := file.Link(TOKEN)

				// 從 Telegram 伺服器獲取該圖片並存至本機。
				getImage(url, id)

				// 尺寸格式：寬度x高度
				size := strconv.Itoa(v.Width) + "x" + strconv.Itoa(v.Height)
				// 在清單中加入這個圖片尺寸。
				list += fmt.Sprintf("\n[%d] - %9s | %6d Bytes", i, size, v.FileSize)
				// 建立一個按鈕給這個圖片尺寸。
				buttons = append(buttons, api.NewInlineKeyboardButtonURL(size, fmt.Sprintf("http://%s:%s/%s.jpg", ADDR, PORT, id)))
			}

			// 建立鍵盤按鍵列。
			kbd := api.NewInlineKeyboardMarkup(
				api.NewInlineKeyboardRow(buttons...),
			)

			// 稍後將會輸出的訊息。
			formattedMsg := "```"
			formattedMsg += list
			formattedMsg += "```"
			formattedMsg += "--\n"
			formattedMsg += "Took: %s\n"
			formattedMsg += "_The file will be served at_ [%s](http://%s:%s)"
			formattedMsg = fmt.Sprintf(formattedMsg, time.Since(start), ADDR, ADDR, PORT)

			// 建立編輯訊息的設定建構體。
			edit := api.EditMessageTextConfig{
				BaseEdit: api.BaseEdit{
					ChatID:    u.Message.Chat.ID,
					MessageID: msg.MessageID,
				},
				Text: formattedMsg,
			}
			edit.ParseMode = "markdown"
			edit.BaseEdit.ReplyMarkup = &kbd

			// 編輯原先的訊息。
			BOT.Send(edit)
		}(update)
	}
}
