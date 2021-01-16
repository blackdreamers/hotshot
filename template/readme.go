package template

var (
	Readme = `# {{title .Alias}}

This is the {{title .Alias}} {{.Type}}

Generated with
` + "```" +
		`
{{.Command}} {{.Alias}}
` + "```" + `

## Usage

Generate the proto code
` + "```" +
		`
make proto
` + "```" + `

Build the binary
` + "```" +
		`
make build
` + "```" + `

Run the service
` + "```" +
		`
./{{.Alias}}-{{.Type}}
` + "```" + `

Build a docker image
` + "```" +
		`
make docker
` + "```"
)
