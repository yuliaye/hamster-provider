package corehttp

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/hamster-shared/hamster-provider/core/context"
	"github.com/hamster-shared/hamster-provider/log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func StartApi(ctx *context.CoreContext) error {
	r := NewMyServer(ctx)

	// router
	v1 := r.Group("/api/v1")
	{

		// basic configuration
		config := v1.Group("/config")
		{
			config.GET("/settting", getConfig)
			config.POST("/settting", setConfig)
			config.POST("/boot", setBootState)
			config.GET("/boot", getBootState)
		}
		chain := v1.Group("/chain")
		{
			chain.GET("/resource", getChainResource)
			chain.GET("/expiration-time", getCalculateInstanceOverdue)
			chain.GET("/account-info", getAccountInfo)
			chain.GET("/staking-info", getStakingInfo)
			chain.POST("/pledge", stakingAmount)
			chain.POST("/withdraw-amount", withdrawAmount)
			chain.POST("/price", changeUnitPrice)
			chain.GET("/reward", queryReward)
			chain.POST("/reward", payoutReward)
		}
		// container routing
		container := v1.Group("/container")
		{
			container.GET("/start", startContainer)
			container.GET("/delete", deleteContainer)
		}

		p2p := v1.Group("/p2p")
		// p2p
		{
			p2p.POST("/listen", listenP2p)
			p2p.POST("/forward", forwardP2p)
			p2p.GET("/ls", lsP2p)
			p2p.POST("/close", closeP2p)
			p2p.POST("/check", checkP2p)
		}
		vm := v1.Group("/vm")
		{
			vm.POST("/create", createVm)
		}
		resource := v1.Group("/resource")
		{
			resource.POST("/modify-price", modifyPrice)
			resource.POST("/add-duration", addDuration)
			resource.POST("/receive-income", receiveIncome)
			resource.POST("/rent-again", rentAgain)
			resource.POST("/delete-resource", deleteResource)
			resource.GET("/receive-income-judge", receiveIncomeJudge)
		}

		thegraph := v1.Group("/thegraph")
		{
			thegraph.POST("/deploy", deployTheGraph)
			thegraph.GET("/ws", execHandler)
			thegraph.GET("/wslog", logHandler)
			thegraph.GET("/status", deployStatus)
		}
	}
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	path, _ := os.Getwd()
	fmt.Println("static path: ", filepath.Join(path, "frontend/dist"))
	r.Use(static.Serve("/", static.LocalFile(filepath.Join(path, "frontend/dist"), true)))
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	port := ctx.GetConfig().ApiPort

	err := OpenWeb(port)

	if err != nil {
		log.GetLogger().Warn("cannot open Explore, http://127.0.0.1:10771, error is :", err.Error())
	}

	err = r.Run(fmt.Sprintf("0.0.0.0:%d", port))

	return err
}

var commands = map[string]string{
	"windows": "start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

func OpenWeb(port int) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/c", fmt.Sprintf("start http://127.0.0.1:%d", port))
	} else {
		cmd = exec.Command(run, fmt.Sprintf("http://127.0.0.1:%d", port))
	}
	return cmd.Start()
}
