package main

//#cgo LDFLAGS: -L . -lntgcalls -Wl,-rpath=./
import "C"
import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	tele "gopkg.in/telebot.v3"
	ntgc "santiirepair.dev/kreacher/ntgcalls"

	tg "github.com/amarnathcjd/gogram/telegram"
	redis "github.com/redis/go-redis/v9"
)

func init() {
	
	var err error

	_, err = exec.LookPath("yt-dlp")
	if err != nil {
		panic("yt-dlp isn't installed")
	}

	_, err = exec.LookPath("speedtest")
	if err != nil {
		panic("speedtest isn't installed")
	}

	ntgcalls = ntgc.NTgCalls()
	// defer ntgcalls.Free()

	bot, err = tele.NewBot(tele.Settings{
		Token:  BotConfig().BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, ctx tele.Context) {
			Error(err.Error())
		},
	})

	if err != nil {
		panic(err)
	}

	defer bot.Close()
	cy.Printf("✔️ Bot client connected to %s", bot.URL)

	ubot, err = tg.NewClient(tg.ClientConfig{
		AppID:    int32(BotConfig().APIID),
		AppHash:  BotConfig().APIHash,
		Session:  ".mtproto",
		LogLevel: tg.LogDisable,
	})

	if err != nil {
		panic(err)
	}

	ubot.Start()
	cy.Println("\n✔️ Initialized new MTProto client, waiting to connect...")

	rdc = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", BotConfig().RedisHost, BotConfig().RedisPort),
		Password: BotConfig().RedisPassword,
		DB:       0, // use default DB
		Protocol: 3, // specify 2 for RESP 2 or 3 for RESP 3
	})

	defer rdc.Close()

	cy.Print("✔️ Redis client connected, waiting for requests...")

	/*
		dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			BotConfig().PostgresUser,
			BotConfig().PostgresPassword,
			BotConfig().PostgresHost,
			BotConfig().PostgresPort,
			BotConfig().PostgresDB,
		)

		idbc, err := pgx.Connect(ctx, dbUrl)

		if err != nil {
			panic(err)
		}

		defer idbc.Close(ctx)
		cy.Println("\n✔️ PGX client connected to PostgreSQL database, waiting for requests...")
	*/

	// Set client instances to late vars in vars.go file.
}

func main() {

	var wg sync.WaitGroup

	if err := bot.SetCommands([]tele.Command{
		{Text: "config", Description: "Set the bot's configuration"},
		{Text: "help", Description: "How to use this"},
		{Text: "ping", Description: "Check the server's latency"},
		{Text: "play_book", Description: "Play a pdf or epub file as an audio book"},
		{Text: "play_song", Description: "Play audio in the voice chat"},
		{Text: "play_video", Description: "Play video in the voice chat"},
		{Text: "streaming", Description: "Any movie or series"},
	}); err != nil {
		panic(err)
	}

	ntgcalls.OnStreamEnd(func(chatId int64, streamType ntgc.StreamType) {
		fmt.Println(chatId)
	})

	ntgcalls.OnConnectionChange(func(chatId int64, state ntgc.ConnectionState) {
		switch state {
		case ntgc.Connecting:
			fmt.Println("Connecting with chatId:", chatId)
		case ntgc.Connected:
			fmt.Println("Connected with chatId:", chatId)
		case ntgc.Failed:
			fmt.Println("Failed with chatId:", chatId)
		case ntgc.Timeout:
			fmt.Println("Timeout with chatId:", chatId)
		case ntgc.Closed:
			fmt.Println("Closed with chatId:", chatId)
		}
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go commands()

	go func() {
		<-sigChan
		fmt.Println("\nReceived SIGINT, exiting...")

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bot.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ubot.Idle()
	}()

	cy.Printf("\n\nBot @%s started, receiving updates...\n", bot.Me.Username)

	wg.Wait()

}
