module Printer

go 1.18

require (
	Printer/db v0.0.0-00010101000000-000000000000 // indirect
	Lib/models v0.0.0-00010101000000-000000000000 // indirect
	Lib/filemanager v0.0.0-00010101000000-000000000000
	github.com/chromedp/cdproto v0.0.0-20230722233645-dbf72f61037f // indirect
	github.com/chromedp/chromedp v0.9.1 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.12.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

require github.com/fabrizio2210/photobook v0.0.0-00010101000000-000000000000

replace (
	Printer/db => ../lib/db
	Lib/models => ../lib/models
	Lib/filemanager => ../lib/filemanager
	github.com/fabrizio2210/photobook => ../lib/github.com/fabrizio2210/photobook
)
