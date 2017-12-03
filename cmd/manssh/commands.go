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
			},
		},
		{
			Name:    "list",
			Usage:   "List or query SSH alias records",
			Action:  list,
			Aliases: []string{"l"},
		},
		{
			Name:    "update",
			Usage:   "Update SSH record by specified alias name",
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
			Usage:   "Delete SSH records by specified alias names",
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
