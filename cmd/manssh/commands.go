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
				cli.StringFlag{Name: "path, p", EnvVar: "MANSSH_DEFAULT_ADD_PATH"},
			},
		},
		{
			Name:    "list",
			Usage:   "List or query SSH alias records",
			Action:  list,
			Aliases: []string{"l"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "ignorecase, ic", Usage: "ignore case while searching"},
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
			},
		},
		{
			Name:    "delete",
			Usage:   "Delete SSH records by specifying alias names",
			Action:  delete,
			Aliases: []string{"d"},
		},
		{
			Name:    "backup",
			Usage:   "Backup SSH alias config records",
			Action:  backup,
			Aliases: []string{"b"},
		},
	}
}
