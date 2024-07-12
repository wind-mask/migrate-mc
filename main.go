package main

import (
	"io"
	"log"
	"os"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/wind-mask/migrate-mc/api"
	"github.com/wind-mask/migrate-mc/modfile"
)

func main() {
	log.SetOutput(io.Discard)
	api.Logger.SetOutput(io.Discard)
	modfile.Logger.SetOutput(io.Discard)
	app := &cli.App{
		Name:    " migrate-mc",
		Usage:   "migration tool for minecraft ",
		Version: VERSION,
		Suggest: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"ver"},
				Usage:   "print verbose log",
				Action: func(c *cli.Context, b bool) error {
					if b {
						log.SetOutput(os.Stdout)
						api.Logger.SetOutput(os.Stdout)
						modfile.Logger.SetOutput(os.Stdout)
					}
					return nil
				},
			},
		},
		Commands: []*cli.Command{
			&migrate_command,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var migrate_command cli.Command = cli.Command{
	Name:    "migrate",
	Aliases: []string{"m"},
	Usage:   "migrate mods from one path to another",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "from",
			Required: true,
			Aliases:  []string{"f"},
			Usage:    "path to migrate mods from",
		},
		&cli.StringFlag{
			Name:     "to",
			Required: true,
			Aliases:  []string{"t"},
			Usage:    "path to migrate mods to",
		},
		&cli.StringFlag{
			Name:     "minecraft-version",
			Required: true,
			Aliases:  []string{"mcv"},
			Usage:    "to minecraft version",
		},
		&cli.StringFlag{
			Name:     "loader",
			Required: true,
			Aliases:  []string{"l"},
			Usage:    "to loader",
		},
	},
	Action: func(cCtx *cli.Context) error {
		log.Println("migrate command")
		mods, err := modfile.Extract_mods(cCtx.String("from"))
		if err != nil {
			return err
		}
		log.Println("mcv:", cCtx.String("minecraft-version"), "loader:", cCtx.String("loader"))
		mUS, err := mods.UpdateToPath(cCtx.String("loader"), cCtx.String("minecraft-version"), cCtx.String("to"))
		if err != nil {
			return err
		}
		mUS = lo.Filter(mUS, func(mus modfile.ModUpdateStatus, _ int) bool {
			return mus.Err != nil
		})
		println("Mods updated failed:", len(mUS))
		for _, mus := range mUS {
			println("Mod update failed:", mus.Err.Error())
		}
		return nil
	},
}
