package main

import "github.com/urfave/cli"

func commands() []cli.Command {
	return []cli.Command{
		{
			Name:    "add",
			Usage:   "Add a new SSH alias record",
			Action:  add,
			Aliases: []string{"a"},
			Flags: []cli.Flag{
				cli.GenericFlag{Name: "config, c", Value: &kvFlag{}},
				cli.StringFlag{Name: "identityfile, i"},
				cli.StringFlag{Name: "addpath, ap", EnvVar: "MANSSH_ADD_PATH"},
				cli.BoolFlag{Name: "path, p", Usage: "dispay the file path of the alias", EnvVar: "MANSSH_SHOW_PATH"},
			},
		},
		{
			Name:    "list",
			Usage:   "List or query SSH alias records",
			Action:  list,
			Aliases: []string{"l"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "ignorecase, ic", Usage: "ignore case while searching"},
				cli.BoolFlag{Name: "path, p", Usage: "dispay the file path of the alias", EnvVar: "MANSSH_SHOW_PATH"},
			},
		},
		{
			Name:    "update",
			Usage:   "Update SSH record by specifying alias name",
			Action:  update,
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
			Usage:   "Delete SSH records by specifying alias names",
			Action:  delete,
			Aliases: []string{"d"},
		},
	}
}
