package main

import (
	"diLesson/pkg/vault"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	var addr, mountPath, file, out string
	flag.StringVar(&addr, "address", "http://127.0.0.1:8200", "vault address")
	flag.StringVar(&mountPath, "mount", "terminals", "engine name")
	flag.StringVar(&file, "file", "terminals.json", "file with terminals credentials")
	flag.StringVar(&out, "out", "", "file to write VAULT_ADDRESS and VAULT_TOKEN")

	flag.Parse()

	fmt.Println("args: ", addr, mountPath, file, out)

	v, err := vault.NewVault(addr, mountPath)
	if err != nil {
		log.Fatal("NewVault error: ", err)
	}

	envs, err := v.Initialize()
	if err != nil {
		log.Fatal("Vault.Initialize error: ", err)
	}

	f := os.Stdout
	if out != "" {
		f, err = os.Create(out)
		defer f.Close()

		if err != nil {
			log.Fatalf("open file %s error: %s\n", out, err)
		}
	}

	for k, v := range envs {
		_, err := fmt.Fprintf(f, "%s=%s\n", k, v)
		if err != nil {
			log.Fatal("fmt.Fprintf error: ", err)
		}
	}

	f2, err := os.Open(file)
	if err != nil {
		log.Fatalf("open file %s error: %v\n", file, err)
	}

	uuids, err := v.Populate(f2)
	if err != nil {
		log.Fatal("Vault.Populate error : ", err)
	}

	fmt.Println("UUIDS:")
	for alias, uuid := range uuids {
		_, err := fmt.Printf("%s: %s\n", alias, uuid)
		if err != nil {
			log.Fatalf("open file %s error: %s\n", out, err)
		}
	}
}
