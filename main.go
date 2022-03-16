package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.finema.co/finema/etda/key-repository-api/consts"
	"gitlab.finema.co/finema/etda/key-repository-api/helpers"
	"gitlab.finema.co/finema/etda/key-repository-api/home"
	core "ssi-gitlab.teda.th/ssi/core"
)

func main() {
	env := core.NewEnv()

	mysql, err := core.NewDatabase(env.Config()).Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "MySQL: %v", err)
		os.Exit(1)
	}

	hsm, err := helpers.NewHSMSession(env.Int(consts.ENVHSMSlot), env.String(consts.ENVHSMPin))
	if err != nil {
		fmt.Fprintf(os.Stderr, "HSM: %v", err)
		os.Exit(1)
	}

	contextOptions := &core.ContextOptions{
		DB:  mysql,
		ENV: env,
		DATA: map[string]interface{}{
			consts.ContextKeyHSMSession: hsm,
		},
	}
	sqlDB, err := mysql.DB()
	sqlDB.SetMaxIdleConns(20000)
	sqlDB.SetConnMaxIdleTime(time.Hour)
	go helpers.KeepHSMAlive(contextOptions, env.Int(consts.ENVHSMSlot), env.String(consts.ENVHSMPin))

	e := core.NewHTTPServer(&core.HTTPContextOptions{
		ContextOptions: contextOptions,
	})

	home.NewHomeHTTPHandler(e)

	core.StartHTTPServer(e, env)
}
