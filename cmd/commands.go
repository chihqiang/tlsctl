package cmd

import "github.com/urfave/cli/v3"

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
