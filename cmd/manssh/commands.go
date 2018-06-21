package main

import "github.com/urfave/cli"

func commands() []cli.Command {
	return []cli.Command{
		{
			Name:    "add",
			Usage:   "Add a new ssh alias record",
			Action:  addCmd,
			Aliases: []string{"a"},
			Flags: []cli.Flag{
				cli.GenericFlag{Name: "config, c", Value: &kvFlag{}},
				cli.StringFlag{Name: "identityfile, i"},
				cli.StringFlag{Name: "addpath, ap", EnvVar: "MANSSH_ADD_PATH"},
				cli.BoolFlag{Name: "path, p", Usage: "display the file path of the alias", EnvVar: "MANSSH_SHOW_PATH"},
			},
		},
		{
			Name:    "list",
			Usage:   "List all or query ssh alias records",
			Action:  listCmd,
			Aliases: []string{"l"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "ignorecase, ic", Usage: "ignore case while searching"},
				cli.BoolFlag{Name: "path, p", Usage: "dispay the file path of the alias", EnvVar: "MANSSH_SHOW_PATH"},
			},
		},
		{
			Name:    "update",
			Usage:   "Update the specified ssh alias",
			Action:  updateCmd,
			Aliases: []string{"u"},
			Flags: []cli.Flag{
				cli.GenericFlag{Name: "config, c", Value: &kvFlag{}},
				cli.StringFlag{Name: "rename, r"},
				cli.StringFlag{Name: "identityfile, i"},
				cli.BoolFlag{Name: "path, p", Usage: "dispay the file path of the alias", EnvVar: "MANSSH_SHOW_PATH"},
			},
		},
		{
			Name:    "delete",
			Usage:   "Delete one or more ssh aliases",
			Action:  deleteCmd,
			Aliases: []string{"d"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "path, p", Usage: "dispay the file path of the alias", EnvVar: "MANSSH_SHOW_PATH"},
			},
		},
		{
			Name:    "backup",
			Usage:   "Backup SSH config files",
			Action:  backupCmd,
			Aliases: []string{"b"},
		},
	}
}
