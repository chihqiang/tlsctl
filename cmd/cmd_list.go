package cmd

import (
	"context"

	"github.com/chihqiang/tlsctl/pkg/stdout"
	"github.com/chihqiang/tlsctl/resource"
	"github.com/urfave/cli/v3"
)

func listCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "list",
		Usage:                  "List all certificates installed on this machine",
		Flags:                  []cli.Flag{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cStorage := setupResourceCache(cmd)
			resources, err := cStorage.GetAllDomainResources()
			if err != nil {
				return err
			}
			table := stdout.NewTablePrinter()
			table.Add([]string{"Domain", "Start Time", "End Time", "Save Path"})
			for _, res := range resources {
				certificate, err := cStorage.ParseResourceFindCertificate(res)
				if err != nil {
					continue
				}
				domain, _ := resource.SanitizedDomain(res.Domain)
				table.Add([]string{
					res.Domain,
					certificate.NotBefore.Format("2006-01-02 15:04:05"),
					certificate.NotAfter.Format("2006-01-02 15:04:05"),
					cStorage.GetSanitizedDomainSavePath(domain),
				})
			}
			return table.Print()
		},
	}
}
