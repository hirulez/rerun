package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type processManager struct {
	conf  *config
	oscmd *exec.Cmd
}

func (pm *processManager) formatBuildTime(duration time.Duration) string {
	return fmt.Sprintf("%.2f(s)", duration.Seconds())
}

func (pm *processManager) run() {
	logger.Debugf("building application %s...", pm.conf.Build)

	start := time.Now()

	os.Remove(pm.conf.Build)
	out, err := exec.Command("go", "build", "-o", pm.conf.Build).CombinedOutput()
	if err != nil {
		logger.Errorf("build failed! %s", err.Error())
		fmt.Printf("%s", out)
		return
	}

	// build success, display build time
	logger.Infof("build took %s", pm.formatBuildTime(time.Since(start)))

	if pm.conf.Test {
		testOut, testErr := exec.Command("go", "test").CombinedOutput()
		if testErr != nil {
			logger.Error("Tests failed!")
			fmt.Printf("==========\n%s==========\n", testOut)
		} else {
			logger.Info("Tests OK!")
		}
	}

	pm.oscmd = exec.Command(pm.conf.Build, pm.conf.Args...)
	pm.oscmd.Stdout = os.Stdout
	pm.oscmd.Stdin = os.Stdin
	pm.oscmd.Stderr = os.Stderr

	logger.Debugf("starting application with arguments: %v", pm.conf.Args)
	err = pm.oscmd.Start()
	if err != nil {
		logger.Errorf("error while starting application! %s", err.Error())
	}
}

func (pm *processManager) stop() {
	logger.Debug("stopping application")

	if pm.oscmd == nil {
		return
	}

	err := pm.oscmd.Process.Kill()
	if err != nil {
		logger.Errorf("error while stopping application! %s", err.Error())
	}
}
