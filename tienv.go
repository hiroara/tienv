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
	app.Usage = "TODO"
	app.Author = "Hiroki Arai"
	app.Email = "hiroara62@gmail.com"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:      "config",
			ShortName: "c",
			Usage:     "handle app/config.json",
			Flags:     flags,
			Subcommands: []cli.Command{
				{
					Name:      "convert",
					ShortName: "c",
					Usage:     "convert environment settings",
					Action:    convertConfig,
				},
				{
					Name:      "restore",
					ShortName: "r",
					Usage:     "restore backup file",
					Action:    restoreConfig,
				},
			},
		},
		{
			Name:      "tiapp",
			ShortName: "t",
			Usage:     "convert env of tiapp.xml",
			Flags:     flags,
			Subcommands: []cli.Command{
				{
					Name:      "convert",
					ShortName: "c",
					Usage:     "convert environment settings",
					Action:    convertTiapp,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "config, c",
							Value: "tiapp-convert.json",
							Usage: "configuration file path",
						},
					},
				},
				{
					Name:      "restore",
					ShortName: "r",
					Usage:     "restore backup file",
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

	handleError(getConfWithBackup(conf, c))

	defer conf.Free()

	err = conf.Replace("//property[@name=\"ti.facebook.appid\"]", "1234567890")
	handleError(err)

	handleError(conf.ReplaceWithConf(c.String("config")))
	handleError(conf.AppendWithConf(c.String("config")))

	target.Write(conf, []byte(conf.Document.String()))

	path, _ := conf.GetFilePath()
	println("Convered config file: " + path)
}

func restoreTiapp(c *cli.Context) {
	conf, err := target.GetTiapp(c.GlobalString("directory"))
	defer conf.Free()
	handleError(err)
	restoreTarget(conf, c)
}

func convertConfig(c *cli.Context) {
	if len(c.Args()) != 2 {
		handleError(errors.New("Usage: tienv <global options> config convert <target env> <as env>"))
	}

	targetEnv := c.Args().First()
	asEnv := c.Args().Get(1)

	conf := target.GetConfigWithRestore(c.GlobalString("directory"), c.GlobalString("backup"))

	handleError(getConfWithBackup(conf, c))

	encoded, err := conf.ConvertEnv(targetEnv, asEnv)
	handleError(err)

	_, err = target.Write(conf, encoded)
	handleError(err)

	path, _ := conf.GetFilePath()
	println("Convered config file: " + path)
}

func restoreConfig(c *cli.Context) {
	restoreTarget(target.GetConfig(c.GlobalString("directory")), c)
}

func restoreTarget(t target.Target, c *cli.Context) {
	backupPath, err := target.Restore(t, c.GlobalString("backup"))
	handleError(err)
	println("Restored from backup file: " + backupPath)
}

type convertTaget interface {
	GetFrom(directory string) convertTaget
	Backup(ext string) (backupPath string, err error)
}

func getConfWithBackup(t target.Target, c *cli.Context) (err error) {
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
