package main

import (
	"errors"
	"github.com/codegangsta/cli"
	"github.com/hiroara/tienv/target"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "tienv"
	app.Usage = "tiny command to convert settings to switch between environment of Titanium project."
	app.Author = "Hiroki Arai"
	app.Email = "hiroara62@gmail.com"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:      "convert",
			ShortName: "c",
			Usage:     "convert setting file",
			Flags:     flags,
			Subcommands: []cli.Command{
				{
					Name:      "config",
					ShortName: "c",
					Usage:     "convert app/config.json",
					Action:    convertConfig,
				},
				{
					Name:      "tiapp",
					ShortName: "t",
					Usage:     "convert tiapp.xml",
					Action:    convertTiapp,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "config, c",
							Value: "tiapp-convert.json",
							Usage: "configuration file path",
						},
					},
				},
			},
		},
		{
			Name:      "restore",
			ShortName: "r",
			Usage:     "restore setting file",
			Flags:     flags,
			Action:    restoreAll,
			Subcommands: []cli.Command{
				{
					Name:      "config",
					ShortName: "c",
					Usage:     "restore app/config.json",
					Action:    restoreConfig,
				},
				{
					Name:      "tiapp",
					ShortName: "t",
					Usage:     "restore tiapp.xml",
					Action:    restoreTiapp,
				},
			},
		},
	}

	app.Run(os.Args)
}

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "directory, d",
		Value: ".",
		Usage: "application directory",
	},
	cli.StringFlag{
		Name:  "backup, b",
		Value: "bak",
		Usage: "extension of backup file",
	},
}

func convertTiapp(c *cli.Context) {
	conf, err := target.GetTiappWithRestore(c.GlobalString("directory"), c.GlobalString("backup"))
	handleError(err)

	handleError(backupTarget(conf, c))

	defer conf.Free()

	handleError(conf.ReplaceWithConf(c.String("config")))
	handleError(conf.AppendWithConf(c.String("config")))

	target.Write(conf, []byte(conf.Document.String()))

	path, _ := conf.GetFilePath()
	println("Convered config file: " + path)
}

func convertConfig(c *cli.Context) {
	if len(c.Args()) != 2 {
		handleError(errors.New("Usage: tienv <global options> config convert <target env> <as env>"))
	}

	targetEnv := c.Args().First()
	asEnv := c.Args().Get(1)

	conf := target.GetConfigWithRestore(c.GlobalString("directory"), c.GlobalString("backup"))

	handleError(backupTarget(conf, c))

	encoded, err := conf.ConvertEnv(targetEnv, asEnv)
	handleError(err)

	_, err = target.Write(conf, encoded)
	handleError(err)

	path, _ := conf.GetFilePath()
	println("Convered config file: " + path)
}

func restoreTiapp(c *cli.Context) {
	conf, err := target.GetTiapp(c.GlobalString("directory"))
	defer conf.Free()
	handleError(err)
	handleError(restoreTarget(conf, c))
}

func restoreConfig(c *cli.Context) {
	handleError(restoreTarget(target.GetConfig(c.GlobalString("directory")), c))
}

func restoreAll(c *cli.Context) {
	conf, _ := target.GetTiapp(c.GlobalString("directory"))
	defer conf.Free()
	restoreTarget(conf, c)
	restoreTarget(target.GetConfig(c.GlobalString("directory")), c)
}

func restoreTarget(t target.Target, c *cli.Context) error {
	backupPath, err := target.Restore(t, c.GlobalString("backup"))
	if err != nil {
		return err
	}
	println("Restored from backup file: " + backupPath)
	return nil
}

type convertTaget interface {
	GetFrom(directory string) convertTaget
	Backup(ext string) (backupPath string, err error)
}

func backupTarget(t target.Target, c *cli.Context) (err error) {
	backupPath, err := target.Backup(t, c.GlobalString("backup"))
	if err == nil {
		println("Backup file created: " + backupPath)
	}
	return
}

func handleError(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
