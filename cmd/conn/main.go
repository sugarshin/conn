package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"

	conn "github.com/sugarshin/conn"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	endpoint := flag.String("endpoint", os.Getenv("CONN_ENDPOINT"), "confluence reset api endpoint")
	username := flag.String("username", os.Getenv("CONN_USERNAME"), "confluence username")
	password := flag.String("password", os.Getenv("CONN_PASSWORD"), "confluence password or token")
	if *endpoint == "" || *username == "" || *password == "" {
		return errors.New("endpoint, username, password are required")
	}

	subpageCmd := flag.NewFlagSet("subpage", flag.ExitOnError)
	subpageCreate := subpageCmd.Bool("create", false, "create")
	parentPageID := subpageCmd.String("parentPageID", "", "parent page id")
	content := &conn.Content{}
	json := jsonValue{content}
	subpageCmd.Var(&json, "content", "content json unmarshal")
	if len(os.Args) < 2 {
		return errors.New("currently, expected 'subpage' subcommands")
	}
	client, err := conn.New(*endpoint, *username, *password)
	if err != nil {
		return err
	}
	switch os.Args[1] {
	case "subpage":
		subpageCmd.Parse(os.Args[2:])
		if *subpageCreate == true {
			if *parentPageID == "" {
				return errors.New("parentPageID is required")
			} else if json.String() != "" {
				return errors.New("content is required")
			}

			if _, err := client.CreateSubPageContent(*parentPageID, content); err != nil {
				return err
			}
		}
	default:
		return errors.New("currently, expected 'subpage' subcommands")
	}
	return nil
}

type jsonValue struct {
	Content *conn.Content
}

func (v jsonValue) Set(s string) error {
	if err := json.Unmarshal([]byte(s), v.Content); err != nil {
		return err
	}
	return nil
}

func (v jsonValue) String() string {
	if v.Content != nil {
		b, _ := json.Marshal(v.Content)
		return string(b)
	}
	return ""
}
