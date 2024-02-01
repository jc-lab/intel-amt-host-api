/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package main

import (
	"github.com/jc-lab/intel-amt-host-api/internal/amt"
	"github.com/jc-lab/intel-amt-host-api/internal/flags"
	"github.com/jc-lab/intel-amt-host-api/internal/local"
	"github.com/jc-lab/intel-amt-host-api/internal/rps"
	"github.com/jc-lab/intel-amt-host-api/pkg/utils"
	"os"

	log "github.com/sirupsen/logrus"
)

const AccessErrMsg = "Failed to execute due to access issues. " +
	"Please ensure that Intel ME is present, " +
	"the MEI driver is installed, " +
	"and the runtime has administrator or root privileges."

func checkAccess() (utils.ReturnCode, error) {
	amtCommand := amt.NewAMTCommand()
	rc, err := amtCommand.Initialize()
	if rc != utils.Success || err != nil {
		return utils.AmtNotDetected, err
	}
	return utils.Success, nil
}

func runRPC(args []string) utils.ReturnCode {
	flags, rc := parseCommandLine(args)
	if rc != utils.Success {
		return rc
	}
	if flags.Local {
		rc = local.ExecuteCommand(flags)
	} else {
		rc = rps.ExecuteCommand(flags)
	}
	return rc
}

func parseCommandLine(args []string) (*flags.Flags, utils.ReturnCode) {
	//process flags
	flags := flags.NewFlags(args)
	rc := flags.ParseFlags()

	if flags.Verbose {
		log.SetLevel(log.TraceLevel)
	} else {
		lvl, err := log.ParseLevel(flags.LogLevel)
		if err != nil {
			log.Warn(err)
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(lvl)
		}
	}

	if flags.JsonOutput {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	}
	return flags, rc
}

func main() {
	rc, err := checkAccess()
	if rc != utils.Success {
		if err != nil {
			log.Error(err.Error())
		}
		log.Error(AccessErrMsg)
		os.Exit(int(rc))
	}
	rc = runRPC(os.Args)
	if err != nil {
		log.Error(err.Error())
	}
	os.Exit(int(rc))
}
