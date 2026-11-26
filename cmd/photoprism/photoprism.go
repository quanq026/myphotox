/*
MyPhotoX - Lightweight Self-Hosted Photos App
Licensed under Version 3 of the GNU Affero General Public License (AGPL-3.0).
*/
package main

import (
	"fmt"
	"os"

	// Embeds the IANA time zone database, which Go uses only for zones the system does not provide.
	_ "time/tzdata"

	"github.com/urfave/cli/v2"

	"github.com/photoprism/photoprism/internal/commands"
	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/pkg/clean"
)

var version = "development"

const appName = "MyPhotoX"
const appAbout = "MyPhotoX"
const appEdition = "lite"
const appDescription = "MyPhotoX is a lightweight, self-hosted photos and videos app."
const appCopyright = "MyPhotoX"

// Metadata contains build specific information.
var Metadata = map[string]any{
	"Name":        appName,
	"About":       appAbout,
	"Edition":     appEdition,
	"Description": appDescription,
	"Version":     version,
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			event.LogPanic(r)
			os.Exit(1)
		}
	}()

	app := cli.NewApp()
	app.Usage = appAbout
	app.Description = appDescription
	app.Version = version
	app.Copyright = appCopyright
	app.EnableBashCompletion = true
	app.Flags = config.Flags.Cli()
	app.Commands = commands.PhotoPrism
	app.Metadata = Metadata
	app.ExitErrHandler = commands.HandleExitErr

	// HandleExitErr prints and exits with the code of a cli.Exit(...) error before
	// app.Run returns, so other errors, including wrapped ones, reach this point.
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, clean.DisplayLines(fmt.Sprint(err)))
		os.Exit(commands.ExitCode(err)) //nolint:gocritic // the deferred recover only handles panics
	}
}
