package cmd

import "github.com/chihqiang/cli"

func CreateCommands() []*cli.Command {
	return []*cli.Command{
		createCommand(),
		localhostCommand(),
		listCommand(),
		clearCommand(),

		deployCommand(),
		scheduledRunCommand(),
		listScheduledCommand(),
		removeScheduledCommand(),

		deployHelpCommand(),
		dnsHelpCommand(),
	}
}
