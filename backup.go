package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/SBanczyk/backup/backend/fs"
)

func main() {
	initFsCommand := flag.NewFlagSet("init fs", flag.ExitOnError)
	initS3Command := flag.NewFlagSet("initS3", flag.ExitOnError)
	backupCommand := flag.NewFlagSet("backup", flag.ExitOnError)
	shadowCommand := flag.NewFlagSet("shadow", flag.ExitOnError)
	statusCommand := flag.NewFlagSet("status", flag.ExitOnError)
	unstageCommand := flag.NewFlagSet("status", flag.ExitOnError)
	destroyCommand := flag.NewFlagSet("status", flag.ExitOnError)
	getCommand := flag.NewFlagSet("status", flag.ExitOnError)
	pullCommand := flag.NewFlagSet("status", flag.ExitOnError)
	pushCommand := flag.NewFlagSet("status", flag.ExitOnError)
	targetDirPtr := initFsCommand.String("target-dir", "", "Target-dir")
	bucketPtr := initS3Command.String("bucket-name", "", "Bucket-name")
	apiKeyPtr := initS3Command.String("api-key", "", "api-key")

	switch os.Args[1] {
	case "init":
		switch os.Args[2] {
		case "fs":
			initFsCommand.Parse(os.Args[3:])
		case "s3":
			initS3Command.Parse(os.Args[3:])
		default:
			fmt.Printf("Wrong argument: %v\n", os.Args[2])
			os.Exit(1)
		}
	case "backup":
		backupCommand.Parse(os.Args[2:])
	case "shadow":
		shadowCommand.Parse(os.Args[2:])
	case "status":
		statusCommand.Parse(os.Args[2:])
	case "unstage":
		unstageCommand.Parse(os.Args[2:])
	case "destroy":
		destroyCommand.Parse(os.Args[2:])
	case "get":
		getCommand.Parse(os.Args[2:])
	case "pull":
		pullCommand.Parse(os.Args[2:])
	case "push":
		pushCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("Wrong argument: %v\n", os.Args[1])
		os.Exit(1)
	}

	if initFsCommand.Parsed() {
		if *targetDirPtr == "" {
			initFsCommand.PrintDefaults()
			os.Exit(1)
		}
		fmt.Printf("init fs: %s\n", *targetDirPtr)
		backend, err := fs.Init(*targetDirPtr)
		if err != nil {
			log.Fatal(err)
		}
		err = backend.UploadFile("downloaded", "uploaded")
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("Removed: %v", *targetDirPtr)
	}

	if initS3Command.Parsed() {
		if *bucketPtr == "" || *apiKeyPtr == "" {
			initS3Command.PrintDefaults()
			os.Exit(1)
		}
		fmt.Printf("init s3: %s %s\n", *bucketPtr, *apiKeyPtr)
	}

	if backupCommand.Parsed() {
		if len(os.Args[2:]) != 0 {
			fmt.Printf("backup: %s\n", os.Args[2:])
		} else {
			fmt.Printf("No arguments\n")
		}

	}

	if shadowCommand.Parsed() {
		if len(os.Args[2:]) != 0 {
			fmt.Printf("shadow: %s\n", os.Args[2:])
		} else {
			fmt.Printf("No arguments\n")
		}

	}

	if statusCommand.Parsed() {
		if len(os.Args[2:]) == 0 {
			fmt.Printf("status\n")
		} else {
			fmt.Printf("Too many arguments\n")
		}

	}

	if unstageCommand.Parsed() {
		if len(os.Args[2:]) != 0 {
			fmt.Printf("unstage: %v\n", os.Args[2:])
		} else {
			fmt.Printf("No arguments\n")
		}

	}

	if destroyCommand.Parsed() {
		if len(os.Args[2:]) != 0 {
			fmt.Printf("destroy: %v\n", os.Args[2:])
		} else {
			fmt.Printf("No arguments\n")
		}

	}

	if getCommand.Parsed() {
		if len(os.Args[2:]) == 0 {
			fmt.Printf("get\n")
		} else {
			fmt.Printf("Too many arguments\n")
		}

	}

	if pullCommand.Parsed() {
		if len(os.Args[2:]) == 0 {
			fmt.Printf("pull\n")
		} else {
			fmt.Printf("Too many arguments\n")
		}

	}

	if pushCommand.Parsed() {
		if len(os.Args[2:]) == 0 {
			fmt.Printf("push\n")
		} else {
			fmt.Printf("Too many arguments\n")
		}

	}

	flag.Parse()
}
