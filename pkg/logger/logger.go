// package logger

// import (
// 	"fmt"
// 	"log"
// 	"log/slog"
// 	"math/rand"
// 	"net"
// 	"os"
// 	"time"
// )

// type Message struct {
// 	Msg string `json:"message"`
// }

// // logrus
// // zap
// func NewLogger() *slog.Logger {
// 	opts := slog.HandlerOptions{
// 		Level: slog.LevelDebug,
// 	}

// 	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
// 	if err != nil {
// 		log.Fatal("Log file ochishda xatolik", err)
// 	}

// 	logger := slog.New(slog.NewTextHandler(file, &opts))
// 	return logger
// }

// var logLevels = []string{"INFO", "WARNING", "ERROR"}

// func SendLog(conn net.Conn) {
// 	for {
// 		logMessage := fmt.Sprintf("%s - Server 1 - %s log message", time.Now().Format(time.RFC3339), logLevels[rand.Intn(len(logLevels))])

// 		conn.Write([]byte(logMessage + "\n"))
// 		fmt.Println("Sent:", logMessage)
// 		time.Sleep(2 * time.Second)
// 	}
// }

package logs

import (
	"google_docs_user/config"
	"net"

	"github.com/sirupsen/logrus"
)

type TCPHook struct {
	Address string
	conn    net.Conn
}

func (hook *TCPHook) Fire(entry *logrus.Entry) error {
	if hook.conn == nil {
		var err error
		hook.conn, err = net.Dial("tcp", hook.Address)
		if err != nil {
			return err
		}
	}

	logMessage, err := entry.String()
	if err != nil {
		return err
	}

	_, err = hook.conn.Write([]byte(logMessage))
	return err
}

func (hook *TCPHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func InitLogger(cfg config.Config) *logrus.Logger {
	logger := logrus.New()

	// JSON formatida log yozish uchun sozlash
	logger.SetFormatter(&logrus.JSONFormatter{})

	// TCP Hook ni ulash
	tcpHook := &TCPHook{Address: "localhost:9090"}
	logger.AddHook(tcpHook)

	return logger
}
