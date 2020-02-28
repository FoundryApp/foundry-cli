package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/urfave/cli/v2"
)

// Build the program with `build -o deploy`
// Then use it as follows:
// $ ./deploy --user create_with_ws
// $ ./deploy --ws create
// $ ./deploy --ws delete

type DeployObject int
const (
	User DeployObject = iota
	Workspace
	BackendModule
)

const (
	USR_CREATE_WITH_WS = "create_with_ws"
)

const (
	WS_CREATE = "create"
	WS_DELETE = "delete"
	WS_UPDATE_SLUG = "update_slug"
)

const (
	BE_CREATE = "create"
	BE_DEPLOY = "deploy_function"
)

func main() {
	app := &cli.App{
		Name: "deploy",
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name: "user",
			},
			&cli.StringFlag{
				Name: "ws", // workspace
			},
			&cli.StringFlag{
				Name: "bem", // backend module
			},
		},
		Action: func(c *cli.Context) error {
			switch {
			case c.String("user") == USR_CREATE_WITH_WS:
				return deploy(User, USR_CREATE_WITH_WS)
			case c.String("ws") == WS_CREATE:
				return deploy(Workspace, WS_CREATE)
			case c.String("ws") == WS_DELETE:
				return deploy(Workspace, WS_DELETE)
			case c.String("ws") == WS_UPDATE_SLUG:
				return deploy(Workspace, WS_UPDATE_SLUG)
			case c.String("bem") == BE_CREATE:
				return deploy(BackendModule, BE_CREATE)
			case c.String("bem") == BE_DEPLOY:
				return deploy(BackendModule, BE_DEPLOY)
			default:
				return errors.New("Unsupported function")
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}


type errCmd struct {
	Dir 	string
	cmds 	[]*exec.Cmd
	Err 	error
}
func (ec *errCmd) addCmd(name string, arg ...string) {
	if ec.Err != nil {
		return
	}

	c := exec.Command(name, arg...)
	c.Dir = ec.Dir
	ec.cmds = append(ec.cmds, c)
}
func (ec *errCmd) run() {
	if ec.Err != nil {
		return
	}

	for _, cmd := range ec.cmds {
		fmt.Println(cmd)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ec.Err = err
			break
		}

		// go _ := func () {
		// 	for {
		// 		b := make([]byte, 1024, 1024)
		// 		n, err := pipe.Read(b)
		// 	}
		// }

		// TODO: Do periodically https://stackoverflow.com/questions/24282709/how-would-you-get-the-output-of-a-command-as-its-running-when-exec-ing-in-go
		// buf := new(bytes.Buffer)
    // buf.ReadFrom(stdout)
		fmt.Printf("%v\n", stdout)

		if err := cmd.Start(); err != nil {
			ec.Err = err
			break
		}

		if err := cmd.Wait(); err != nil {
			ec.Err = err
			break
		}

	}
}
func deploy(do DeployObject, fn string) error {
	// TODO: Not portable
	const root = "/Users/valenta.and.thomas/Developer/go-back"

	// gcloud functions deploy create-user-with-workspace --entry-point CreateWithWorkspace --runtime go113 --trigger-http
	var p string
	var deployedName string
	var entry string
	var trigger string

	switch do {
	case User:
		p = path.Join(root, "user")
		switch fn {
		case USR_CREATE_WITH_WS:
			deployedName = "create-user-with-ws"
			entry = "CreateWithWorkspace"
			trigger = "--trigger-http"
		}

	case Workspace:
		p = path.Join(root, "workspace")
		switch fn {
		case WS_CREATE:
			deployedName = "create-ws"
			entry = "Create"
			trigger = "--trigger-http"
		case WS_DELETE:
			deployedName = "delete-ws"
			entry = "Delete"
			trigger = "--trigger-http"
		case WS_UPDATE_SLUG:
			deployedName = "update-ws-slug"
			entry = "UpdateSlug"
			trigger = "--trigger-http"
		}

	case BackendModule:
		p = path.Join(root, "bemodule")
		switch fn {
		case BE_CREATE:
			deployedName = "create-bemodule"
			entry = "Create"
			trigger = "--trigger-http"
		case BE_DEPLOY:
			deployedName = "deploy-bemodule"
			entry = "Deploy"
			trigger = "--trigger-http"
		}
	}

	ec := &errCmd{Dir: p}
	ec.addCmd("go", "build")
	ec.addCmd("go", "mod", "vendor")
	ec.addCmd("gcloud", "functions", "deploy", deployedName, "--entry-point", entry, "--runtime", "go111", trigger)
	ec.addCmd("echo", "==== DONE ====")
	ec.run()
	return ec.Err
}
