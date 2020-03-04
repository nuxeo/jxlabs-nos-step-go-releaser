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
wraps goreleaser so we can get the git token from a kubernetes secret 
`

	createExample = `
step-go-releaser --organisation=jenkins-x-labs --revision=1b59ffc --branch=master --build-date=20200303-22:14:54 --go-version=1.12.17 --root-package=github.com/jenkins-x-labs/gsm-controller --version=0.0.17
`
)

// NewCmdHelloWorld creates a command object for the "hello world" command
func NewCmdGoReleaser() *cobra.Command {
	o := &options{}

	cmd := &cobra.Command{
		Use:     "step-go-releaser",
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
