package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {

	app := cli.NewApp()
	app.Name = "kicad plugin"
	app.Usage = "kicad plugin"
	app.Action = run
	app.Version = fmt.Sprintf("0.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "client.code",
			Usage:  "enterprise client code",
			EnvVar: "PLUGIN_CLIENT_CODE",
		},
		cli.StringFlag{
			Name:   "client.name",
			Usage:  "client name",
			EnvVar: "PLUGIN_CLIENT_NAME",
		},
		cli.StringSliceFlag{
			Name:   "project.codes",
			Usage:  "enterprise project Codes",
			EnvVar: "PLUGIN_PROJECTS_CODES",
		},
		cli.StringSliceFlag{
			Name:   "projects.names",
			Usage:  "project names",
			EnvVar: "PLUGIN_PROJECTS_NAMES",
		},
		cli.BoolFlag{
			Name:   "options.schematic",
			Usage:  "generate schematic",
			EnvVar: "PLUGIN_SCHEMATIC",
		},
		cli.BoolFlag{
			Name:   "options.bom",
			Usage:  "generate bom",
			EnvVar: "PLUGIN_BOM",
		},
		cli.StringFlag{
			Name:   "options.gerber",
			Usage:  "generate gerber files",
			EnvVar: "PLUGIN_GERBER",
		},
		cli.StringSliceFlag{
			Name:   "deps.libs",
			Usage:  "Download dependencies",
			EnvVar: "PLUGIN_LIBRARY",
		},
		cli.StringSliceFlag{
			Name:   "deps.pretty",
			Usage:  "Download footprints",
			EnvVar: "PLUGIN_PRETTY",
		},
		cli.StringSliceFlag{
			Name:   "deps.3d",
			Usage:  "Download 3d Models",
			EnvVar: "PLUGIN_3D",
		},
		cli.StringSliceFlag{
			Name:   "deps.template",
			Usage:  "Download templates",
			EnvVar: "PLUGIN_TEMPLATE",
		},
		cli.StringFlag{
			Name:   "deps.basedir",
			Usage:  "Base directory for dependencies",
			EnvVar: "PLUGIN_DEPS_DIR",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {

	plugin := Plugin{
		Client: Client{
			Code: c.String("client.code"),
			Name: c.String("client.name"),
		},
		Projects: Projects{
			Codes: c.StringSlice("projects.codes"),
			Names: c.StringSlice("projects.names"),
		},
		Options: Options{
			Sch:    c.Bool("options.schematic"),
			Bom:    c.Bool("options.bom"),
			GrbGen: c.IsSet("options.gerber"),
		},
		Dependencies: Dependencies{
			Libraries:  c.StringSlice("deps.libs"),
			Footprints: c.StringSlice("deps.pretty"),
			Modules3d:  c.StringSlice("deps.3d"),
			Basedir:    c.String("deps.basedir"),
			Templates:  c.StringSlice("deps.template"),
		},
	}

	if plugin.Options.GrbGen {
		err := json.Unmarshal([]byte(c.String("options.gerber")), &plugin.Options.Grb)
		if err != nil {
			return err
		}
	}

	return plugin.Exec()
}
