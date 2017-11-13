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
			},
		},
		{
			Name:    "list",
			Usage:   "Listing or query SSH alias records",
			Action:  list,
			Aliases: []string{"l"},
		},
		{
			Name:    "update",
			Usage:   "Update SSH record by alias name",
			Action:  update,
			Aliases: []string{"u"},
			Flags: []cli.Flag{
				cli.GenericFlag{Name: "config, c", Value: &kvFlag{}},
				cli.StringFlag{Name: "rename, r"},
			},
		},
		{
			Name:    "delete",
			Usage:   "Delete SSH records by alias name",
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
