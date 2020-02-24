package pkg

import (
	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/spf13/cobra"
)

type options struct {
	Cmd    *cobra.Command
	Args   []string
	Runner util.Commander
	*opts.CommonOptions

	organisation string
	revision     string
	branch       string
	buildDate    string
	goVersion    string
	version      string
	rootPackage  string
}

const (
	version      = "version"
	organisation = "organisation"
	revision     = "revision"
	branch       = "branch"
	buildDate    = "build-date"
	goVersion    = "go-version"
	rootPackage  = "root-package"
)

var (
	createLong = `
This is a hello world quickstart for writing CLI's in Go.  This quickstart will setup automatic CI and release
pipelines using Jenkins X, upon release you will get cross platform binaries uploaded as a GitHub release.
`

	createExample = `
# print hello to the terminal
goreleaser --org foo
`
)

// NewCmdHelloWorld creates a command object for the "hello world" command
func NewCmdHelloWorld() *cobra.Command {
	o := &options{}

	cmd := &cobra.Command{
		Use:     "goreleaser",
		Short:   "wraps the go releaser tool getting required secrets needed to upload artifacts to GitHub",
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Args = args
			return o.Run()
		},
	}

	cmd.Flags().StringVarP(&o.organisation, organisation, "", "", "the git organisation")
	cmd.Flags().StringVarP(&o.revision, revision, "", "", "git revision")
	cmd.Flags().StringVarP(&o.branch, branch, "", "", "git branch")
	cmd.Flags().StringVarP(&o.buildDate, buildDate, "", "", "build date")
	cmd.Flags().StringVarP(&o.version, version, "", "", "version")
	cmd.Flags().StringVarP(&o.goVersion, goVersion, "", "", "go version")
	cmd.Flags().StringVarP(&o.rootPackage, rootPackage, "", "", "root package")
	o.Cmd = cmd

	return cmd
}
