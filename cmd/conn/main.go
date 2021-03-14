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
	token := flag.String("token", os.Getenv("CONN_TOKEN"), "confluence token or password")
	if *endpoint == "" || *username == "" || *token == "" {
		return errors.New("endpoint, username, token are required")
	}

	childpageCmd := flag.NewFlagSet("childpage", flag.ExitOnError)
	childpageCreate := childpageCmd.Bool("create", false, "create")
	parentPageID := childpageCmd.String("parentPageID", "", "parent page id")
	content := &conn.Content{}
	json := jsonValue{content}
	childpageCmd.Var(&json, "content", "content json unmarshal")
	if len(os.Args) < 2 {
		return errors.New("currently, expected 'childpage' subcommand")
	}
	client, err := conn.New(*endpoint, *username, *token)
	if err != nil {
		return err
	}
	switch os.Args[1] {
	case "childpage":
		childpageCmd.Parse(os.Args[2:])
		if *childpageCreate == true {
			if *parentPageID == "" {
				return errors.New("parentPageID is required")
			} else if json.String() != "" {
				return errors.New("content is required")
			}

			if _, err := client.CreateChildPageContent(*parentPageID, content); err != nil {
				return err
			}
		}
	default:
		return errors.New("currently, expected 'childpage' subcommand")
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
