package ewn

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"github.com/go-gomail/gomail"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Notify send notifications over all chanels
func Notify(msg *Message, cfg *viper.Viper) {
	if err := sendEmail(msg, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Email sending failed with error: %s\n", err)
	}
	if err := sendGelf(msg, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "GELF sending failed with error: %s\n", err)
	}
	return
}

func isFailed(msg *Message) bool {
	if msg.GeneralError != nil {
		return true
	}
	for _, ret := range msg.Retries {
		for _, val := range msg.Args.ValidExitCode {
			if val == ret.ExitCode {
				return false
			}
		}
	}
	return true
}

func sendEmail(msg *Message, cfg *viper.Viper) error {
	if cfg.GetBool("email.enabled") && isFailed(msg) {
		recipients := map[string][]string{"To": cfg.GetStringSlice("email.recipients")}
		var messageFull string
		var sender *gomail.Dialer
		conn, err2 := connStrToStruct(cfg.GetString("email.host"))
		if err2 != nil {
			return err2
		}
		messageHeader := fmt.Sprintf("Command: %s\n"+
			"Host: %s\n"+
			"Comment: %s\n"+
			"Args: %+v\n\n",
			msg.Args.Command,
			msg.Host,
			msg.Args.Comment,
			msg.Args)

		eMessage := gomail.NewMessage()
		if cfg.GetString("email.from") == "" {
			eMessage.SetHeader("From", fmt.Sprintf("ewn@%s", msg.Host))
		} else {
			eMessage.SetHeader("From", cfg.GetString("email.from"))
		}
		eMessage.SetHeaders(recipients)
		eMessage.SetHeader("Subject", fmt.Sprintf("ewn@%s FAILED: %s", msg.Host, msg.Args.Command))
		if msg.GeneralError != nil {
			messageFull = messageHeader + fmt.Sprintf("General Error: %s", msg.GeneralError)
		} else {
			messageFull = messageHeader
			for _, v := range msg.Retries {
				messageFull += fmt.Sprintf("Retry number: %d\n"+
					"Start time: %s\n"+
					"End time: %s\n"+
					"Duration: %s\n"+
					"Exit code: %d\n"+
					"Output: \n%s\n\n",
					v.Retry,
					v.StartTime,
					v.EndTime,
					v.Duration,
					v.ExitCode,
					v.Output)
			}
		}
		eMessage.SetBody("text/plain", messageFull)
		if cfg.GetBool("email.secure") {
			sender = gomail.NewDialer(conn.host, conn.port, cfg.GetString("email.user"), cfg.GetString("email.pass"))
		} else {
			sender = gomail.NewPlainDialer(conn.host, conn.port, cfg.GetString("email.user"), cfg.GetString("email.pass"))
		}
		err := sender.DialAndSend(eMessage)
		if err != nil {
			return err
		}
	}
	return nil
}

type smtpConn struct {
	host string
	port int
}

func connStrToStruct(s string) (conn smtpConn, err error) {
	sSlice := strings.Split(s, ":")
	conn.host = sSlice[0]
	if len(sSlice) > 1 {
		port, err := strconv.Atoi(sSlice[1])
		if err != nil {
			return conn, err
		}
		conn.port = port
	} else {
		conn.port = 25
	}
	return conn, nil
}

func sendGelf(msg *Message, cfg *viper.Viper) error {
	if cfg.GetBool("graylog.enabled") {
		if msg.GeneralError != nil {
			out, err := stripOutput(fmt.Sprintf("General Error: %s", msg.GeneralError), 32765, "\n<Output truncated>")
			if err != nil {
				return err
			}
			gelfMessage := gelf.Message{Version: "1.1",
				Host:     msg.Host,
				Short:    msg.Args.Command,
				TimeUnix: float64(time.Now().Unix()),
				Level:    1,
				Extra: map[string]interface{}{
					"command":     msg.Args.Command,
					"start_date":  time.Now().Unix(),
					"finish_date": time.Now().Unix(),
					"duration":    0,
					"exitcode":    1,
					"comment":     msg.Args.Comment,
					"retry":       "0/0",
					"tag":         cfg.GetString("graylog.tag"),
					"output":      out,
					"failed":      1,
				}}
			gelfWriter, err := gelf.NewUDPWriter(fmt.Sprintf("%s:%d", cfg.GetString("graylog.host"), cfg.GetInt("graylog.port")))
			if err != nil {
				return err
			}
			defer gelfWriter.Close()
			err = gelfWriter.WriteMessage(&gelfMessage)
			if err != nil {
				return err
			}
			return nil
		}
		for _, retry := range msg.Retries {
			out, err := stripOutput(retry.Output, 32765, "\n<Output truncated>")
			if err != nil {
				return err
			}
			gelfMessage := gelf.Message{Version: "1.1",
				Host:     msg.Host,
				Short:    msg.Args.Command,
				TimeUnix: float64(time.Now().Unix()),
				Level:    1,
				Extra: map[string]interface{}{
					"command":     msg.Args.Command,
					"start_date":  retry.StartTime.Unix(),
					"finish_date": retry.EndTime.Unix(),
					"duration":    retry.Duration / time.Second,
					"exitcode":    retry.ExitCode,
					"comment":     msg.Args.Comment,
					"retry":       fmt.Sprintf("%d/%d", retry.Retry, msg.Args.Retry),
					"tag":         cfg.GetString("graylog.tag"),
					"output":      out,
				}}

			if isFailed(msg) {
				gelfMessage.Extra["failed"] = 1
			} else {
				gelfMessage.Extra["failed"] = 0
			}

			gelfWriter, err := gelf.NewUDPWriter(fmt.Sprintf("%s:%d", cfg.GetString("graylog.host"), cfg.GetInt("graylog.port")))
			if err != nil {
				return err
			}
			defer gelfWriter.Close()
			err = gelfWriter.WriteMessage(&gelfMessage)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func stripOutput(output string, maxLenBytes int, text string) (string, error) {
	if len(output) <= maxLenBytes {
		return output, nil
	} else if maxLenBytes == len(text) {
		return text, nil
	} else if maxLenBytes < len(text) {
		return "", fmt.Errorf("max string length less than text")
	}

	maxPayloadBytes := maxLenBytes - len(text)
	var out string
	i, char := 0, 0
	for {
		runeValue, size := utf8.DecodeRuneInString(output[i:])
		if i + size > maxPayloadBytes {
			break
		}
		out += string(runeValue)
		i += size
		char++
	}
	out += text
	return out, nil
}
