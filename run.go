package main

import (
	"flock/flock"
	"flock/util"
	"fmt"
	"os"
	"syscall"
)

func help() {
	util.LogInfo(
		`Usage:
	run [OPTION] [SCOPE:]SCRIPT

Options:
	-h, --help      show this help message run, then exit
	-V, --version   output version information, then exit
	-v, --view      view the content of script, then exit
	-u, --update    force to update the script before run
	-c, --clean     clean out all scripts cached in local
	-i INTERPRETER  run with interpreter(e.g., bash, python)
	-I, --init      create configuration and cache directory

Examples:
	run pt-summary
	run github:runscripts/scripts/pt-summary

Report bugs to <https://github.com/runscripts/run/issues>.`,
	)
	util.LogInfo("\n")
}
func initialize() {
	for _, arg := range os.Args {
		if arg == "-I" || arg == "--init" {
			if util.IsRunInstalled() {
				util.LogInfo("Run is already installed\n")
			} else {
				if os.Geteuid() != 0 {
					util.LogError("Root privilege is required\n")
					os.Exit(1)
				}
				// Create script cache directory.
				err := os.MkdirAll(util.DATA_DIR, 0777)
				if err != nil {
					util.ExitError(err)
				}
				// Download run.conf, VERSION and run.1 from master branch.
				err = util.Fetch(util.MASTER_URL+"run.conf", util.CONFIG_PATH)
				if err != nil {
					util.ExitError(err)
				}
				err = util.Fetch(util.MASTER_URL+"VERSION", util.DATA_DIR+"/VERSION")
				if err != nil {
					util.ExitError(err)
				}
				err = util.Fetch(util.MASTER_URL+"man/run.1", "/usr/share/man/man1/run.1.gz")
				if err != nil {
					util.ExitError(err)
				}
			}
			os.Exit(0)
		}
	}
}
func main() {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	util.SetConfigPath()
	util.SetDataDir()
	initialize()
	// If run is not installed.
	if !util.IsRunInstalled() {
		util.LogError("Run is not installed yet. You need to 'run --init' as root.\n")
		os.Exit(1)
	}
	// Show help message if no parameter given.
	if len(os.Args) == 1 {
		help()
		return
	}
	// Parse configuration and runtime options.
	config, err := util.NewConfig()
	if err != nil {
		util.ExitError(err)
	}
	options, err := util.NewOptions(config)
	if err != nil {
		util.ExitError(err)
	}
	// If print help message.
	if options.Help {
		help()
		return
	}
	// If output version information.
	if options.Version {
		version, err := os.ReadFile(util.DATA_DIR + "/VERSION")
		if err != nil {
			util.ExitError(err)
		}
		util.LogInfo("Run version %s\n", version)
		return
	}
	// If clean out scripts.
	if options.Clean {
		util.LogInfo("Do you want to clear out the script cache? [Y/n] ")
		var answer string
		fmt.Scanln(&answer)
		if answer == "Y" || answer == "y" {
			// rm -rf $DATA_DIR/* will remove VERSION. Use $DATA_DIR/*/ instead.
			util.Exec([]string{"sh", "-x", "-c", "rm -rf " + util.DATA_DIR + "/*/"})
		}
		return
	}
	// If not script given.
	if options.Fields == nil {
		util.LogError("The script to run is not specified\n")
		os.Exit(1)
	}
	// Ensure the cache directory has been created.
	cacheID := options.CacheID
	cacheDir := util.DATA_DIR + "/" + options.Scope + "/" + cacheID
	err = os.MkdirAll(cacheDir, 0777)
	if err != nil {
		util.ExitError(err)
	}

	// Lock the script.
	lockPath := cacheDir + ".lock"
	err = flock.Flock(lockPath)
	if err != nil {
		util.LogError("%s: %v\n", lockPath, err)
		os.Exit(1)
	}

	// Download the script.
	scriptPath := cacheDir + "/" + options.Script
	_, err = os.Stat(scriptPath)
	if os.IsNotExist(err) || options.Update {
		err = util.Fetch(options.URL, scriptPath)
		if err != nil {
			util.ExitError(err)
		}
	}

	// If view the script.
	if options.View {
		flock.Funlock(lockPath)
		util.Exec([]string{"cat", scriptPath})
	}

	// Run the script.
	flock.Funlock(lockPath)
	if options.Interpreter == "" {
		util.Exec(append([]string{scriptPath}, options.Args...))
	} else {
		util.Exec(append([]string{options.Interpreter, scriptPath}, options.Args...))
	}
}
