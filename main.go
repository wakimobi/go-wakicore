package main

import (
	"time"

	"github.com/idprm/go-pass-tsel/src/cmd"
	"github.com/idprm/go-pass-tsel/src/config"
)

func main() {
	/**
	 * LOAD CONFIG
	 */
	cfg, err := config.LoadSecret("secret.yaml")
	if err != nil {
		panic(err)
	}

	/**
	 * SET TIMEZONE
	 */
	loc, _ := time.LoadLocation(cfg.App.TimeZone)
	time.Local = loc

	cmd.Execute()

}
