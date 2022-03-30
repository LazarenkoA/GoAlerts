package notify

import (
	"bytes"
	"errors"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

type CLI struct {
	BaseNotify

	Comand string   `yaml:"comand"`
	Args   []string `yaml:"args"`
}

func (cli *CLI) Init(logger *logrus.Entry) *CLI {
	cli.logger = logger.WithField("notifyType", "CLI")
	return cli
}

func (cli *CLI) Notify(items []interface{}) {
	if len(items) == 0 {
		return
	}

	cli.logger.Info("оповещение CLI")

	for _, item := range items {
		for i, _ := range cli.Args {
			cli.Args[i] = cli.buildMessages(cli.Args[i], item)
		}
		cli.run(cli.Args)
	}
}

func (cli *CLI) getParamsFromMsg(msg string, out map[string]interface{}) {
	start := strings.Index(msg, "%") + 1
	end := strings.Index(msg[start:], "%") + start

	if start < 0 || end < 0 || start > end {
		return
	}

	out[msg[start:end]] = nil
	cli.getParamsFromMsg(msg[end+1:], out)
}

func (cli *CLI) run(args []string) {
	cmd := exec.Command(cli.Comand, args...)
	cmd.Stdout = new(bytes.Buffer)
	cmd.Stderr = new(bytes.Buffer)

	if err := cmd.Run(); err != nil {
		cli.logger.WithError(err).WithField("Args", args).Errorf("ошибка выполнения команды %q", cli.Comand)
		return
	}
	cli.logger.Debug("Stdout: ", cmd.Stdout.(*bytes.Buffer).String())
	cli.logger.Debug("Stderr: ", cmd.Stderr.(*bytes.Buffer).String())
}

func (cli *CLI) CheckParams() error {
	cli.Comand = strings.Trim(cli.Comand, " ")
	if len(cli.Comand) == 0 {
		return errors.New("команда не заполнена")
	}

	return nil
}
