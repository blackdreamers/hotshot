package new

import (
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xlab/treeprint"

	"github.com/blackdreamers/hotshot/cmd"
	tmpl "github.com/blackdreamers/hotshot/template"
)

const rootPath = "github.com/blackdreamers/"

func init() {
	cmd.Register(&cli.Command{
		Name:        "new",
		Usage:       "Create a service template",
		Description: `'hotshot new' scaffolds a new service skeleton. Example: 'hotshot new helloworld && cd helloworld'`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "type",
				Aliases:     []string{"t"},
				Usage:       "project type",
				DefaultText: "service",
			},
		},
		Action: Run,
	})
}

type config struct {
	// foo
	Alias string
	// hotshot new example
	Command string
	// api, service
	Type string
	// github.com/blackdreamers/foo
	Dir string
	// $GOPATH/src/github.com/blackdreamers/foo
	GoDir string
	// $GOPATH
	GoPath string
	// UseGoPath
	UseGoPath bool
	// Files
	Files []file
	// Comments
	Comments []string
}

type file struct {
	Path string
	Tmpl string
}

func write(c config, file, tmpl string) error {
	fn := template.FuncMap{
		"title": func(s string) string {
			return strings.ReplaceAll(strings.Title(s), "-", "")
		},
		"dehyphen": func(s string) string {
			return strings.ReplaceAll(s, "-", "")
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.New("f").Funcs(fn).Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(f, c)
}

func create(c config) error {
	// check if dir exists
	if _, err := os.Stat(c.Alias); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", c.Alias)
	}

	// just wait
	<-time.After(time.Millisecond * 250)

	fmt.Printf("Creating service %s\n\n", c.Alias)

	t := treeprint.New()

	// write the files
	for _, file := range c.Files {
		f := filepath.Join(c.Alias, file.Path)
		dir := filepath.Dir(f)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		addFileToTree(t, file.Path)
		if err := write(c, f, file.Tmpl); err != nil {
			return err
		}
	}

	// print tree
	fmt.Println(t.String())

	for _, comment := range c.Comments {
		fmt.Println(comment)
	}

	// just wait
	<-time.After(time.Millisecond * 250)

	return nil
}

func addFileToTree(
	root treeprint.Tree,
	file string,
) {
	split := strings.Split(file, "/")
	curr := root
	for i := 0; i < len(split)-1; i++ {
		n := curr.FindByValue(split[i])
		if n != nil {
			curr = n
		} else {
			curr = curr.AddBranch(split[i])
		}
	}
	if curr.FindByValue(split[len(split)-1]) == nil {
		curr.AddNode(split[len(split)-1])
	}
}

func protoComments(goDir, alias string) []string {
	return []string{
		"\ndownload protoc zip packages (protoc-$VERSION-$PLATFORM.zip) and install:\n",
		"visit https://github.com/protocolbuffers/protobuf/releases",
		"\ndownload protobuf for micro:\n",
		"go get -u github.com/golang/protobuf/proto",
		"go get -u github.com/golang/protobuf/protoc-gen-go",
		"go get github.com/gogo/protobuf/protoc-gen-gofast",
		"go get github.com/micro/go-micro/cmd/protoc-gen-micro/v2",
		"\ncompile the proto file " + alias + ".proto:\n",
		"cd " + alias,
		"make proto\n",
	}
}

func Run(ctx *cli.Context) error {
	atype := ctx.String("type")
	alias := ctx.Args().First()
	if len(alias) == 0 {
		fmt.Println("specify service name")
		return nil
	}

	// set the command
	command := "hotshot new"
	if len(atype) > 0 {
		command += " --type=" + atype
	}

	// default type
	if len(atype) == 0 {
		atype = "srv"
	}

	// check if the path is absolute, we don't want this
	// we want to a relative path so we can install in GOPATH
	if path.IsAbs(alias) {
		fmt.Println("require relative path as service will be installed in GOPATH")
		return nil
	}

	var goPath string
	var goDir string

	goPath = build.Default.GOPATH

	// don't know GOPATH, runaway....
	if len(goPath) == 0 {
		fmt.Println("unknown GOPATH")
		return nil
	}

	// attempt to split path if not windows
	if runtime.GOOS == "windows" {
		goPath = strings.Split(goPath, ";")[0]
	} else {
		goPath = strings.Split(goPath, ":")[0]
	}
	goDir = filepath.Join(goPath, "src", path.Clean(alias))

	c := config{
		Command:   command,
		Alias:     alias,
		Type:      atype,
		Comments:  protoComments(goDir, alias),
		Dir:       rootPath + alias,
		GoDir:     goDir,
		GoPath:    goPath,
		UseGoPath: false,
	}

	switch atype {
	case "srv":
	case "service":
		c.Files = append(c.Files, []file{
			{"main.go", tmpl.MainSRV},
			{"subscriber/" + alias + ".go", tmpl.SubscriberSRV},
		}...)
	case "api":
		c.Files = append(c.Files, []file{
			{"main.go", tmpl.MainAPI},
		}...)
	default:
	}

	c.Files = append(c.Files, []file{
		{"handler/" + alias + ".go", tmpl.HandlerSRV},
		{"proto/" + alias + "/" + alias + ".proto", tmpl.ProtoSRV},
		{"plugin.go", tmpl.Plugin},
		{"generate.go", tmpl.GenerateFile},
		{"Dockerfile", tmpl.DockerSRV},
		{"Makefile", tmpl.Makefile},
		{"README.md", tmpl.Readme},
		{".gitignore", tmpl.GitIgnore},
	}...)

	// set gomodule
	if os.Getenv("GO111MODULE") != "off" {
		c.Files = append(c.Files, file{"go.mod", tmpl.Module})
	}

	// create the files
	return create(c)
}
