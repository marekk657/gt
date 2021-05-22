package main

import (
	"encoding/json"
	"fmt"
	"gt/services"
	"gt/services/container"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	argCommand = iota
	argCommandFiles
	argCommandContainerPath
	argCommandSignatureID
)

const (
	argCommandCreate          = "create"
	argCommandOpen            = "open"
	argCommandAddSignature    = "add-signature"
	argCommandRemoveSignature = "remove-signature"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("please specify command:", []string{argCommandCreate, argCommandOpen, argCommandRemoveSignature, argCommandAddSignature})
		os.Exit(-1)
	}

	settings, err := readsettings()
	if err != nil {
		fmt.Println("failed to load settings")
		os.Exit(-1)
	}

	ksiSigner, err := services.NewKSISigner(settings.Endpoint, settings.Username, settings.Password)
	if err != nil {
		fmt.Println("ksi error:", err)
		os.Exit(-1)
	}

	sigCreator := container.NewSignatureCreator(ksiSigner)
	archiveService := container.NewZipArchiveService()
	signer := container.NewSigner(sigCreator, archiveService)
	creator := container.NewCreator(sigCreator, archiveService)

	cmd := args[0]
	switch cmd {
	case argCommandCreate:
		filepaths := strings.Split(args[2], ",")

		if err := creator.Create(filepaths, args[1]); err != nil {
			fmt.Println("error:", err)
			os.Exit(-1)
		}
	case argCommandOpen:
		paths, err := archiveService.Extract(args[1])
		if err != nil {
			fmt.Println("error", err)
			os.Exit(-1)
		}
		fmt.Println("extacted container files: ", paths)
	case argCommandAddSignature:
		if err := signer.AddSignature(args[1]); err != nil {
			fmt.Println("error", err)
			os.Exit(-1)
		}
	case argCommandRemoveSignature:
		i, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("invalid int provided as index:", args[2])
			os.Exit(-1)
		}

		if err := signer.RemoveSignature(args[1], i); err != nil {
			fmt.Println("error", err)
			os.Exit(-1)
		}
	default:
		fmt.Println("unknown command")
		os.Exit(-1)
	}
}

type settings struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Endpoint string `json:"endpoint"`
}

func readsettings() (settings, error) {
	f, err := os.Open("settings.json")
	if err != nil {
		return settings{}, err
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return settings{}, err
	}

	var set settings
	if err := json.Unmarshal(bs, &set); err != nil {
		return settings{}, err
	}

	return set, nil
}
