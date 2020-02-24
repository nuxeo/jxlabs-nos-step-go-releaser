package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jenkins-x/jx/pkg/log"

	"github.com/jenkins-x/jx/pkg/util"

	"github.com/jenkins-x/jx/pkg/cmd/clients"
	"github.com/jenkins-x/jx/pkg/cmd/opts"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"
)

const (
	githubappLabel = "jenkins.io/githubapp-owner="
	githubPassword = "password"
)

// Run implements the command
func (o *options) Run() error {
	if o.organisation == "" {
		return util.MissingArgument(organisation)
	}
	if o.revision == "" {
		return util.MissingArgument(revision)
	}
	if o.branch == "" {
		return util.MissingArgument(branch)
	}
	if o.version == "" {
		return util.MissingArgument(version)
	}
	if o.buildDate == "" {
		return util.MissingArgument(buildDate)
	}
	if o.goVersion == "" {
		return util.MissingArgument(goVersion)
	}
	if o.rootPackage == "" {
		return util.MissingArgument(rootPackage)
	}

	f := clients.NewFactory()
	o.CommonOptions = opts.NewCommonOptionsWithTerm(f, os.Stdin, os.Stdout, os.Stderr)

	o.Runner = &util.Command{
		Out: o.Out,
		Err: o.Err,
	}

	return o.goReleaser()

}

func (o *options) goReleaser() error {

	token, err := o.getToken()
	if err != nil {
		return errors.Wrapf(err, "failed to get github token for organisation %s", o.organisation)
	}
	o.Runner.SetEnvVariable("GITHUB_TOKEN", token)
	o.Runner.SetEnvVariable("REV", o.revision)
	o.Runner.SetEnvVariable("BRANCH", o.branch)
	o.Runner.SetEnvVariable("VERSION", o.version)
	o.Runner.SetEnvVariable("BUILDDATE", o.buildDate)
	o.Runner.SetEnvVariable("GOVERSION", o.goVersion)
	o.Runner.SetEnvVariable("ROOTPACKAGE", o.rootPackage)

	o.Runner.SetName("goreleaser")

	args := []string{"release", "--config=.goreleaser.yml", "--rm-dist", "--release-notes=./changelog.md", "--skip-validate"}
	o.Runner.SetArgs(args)
	o.Out = os.Stdout
	o.Err = os.Stderr
	_, err = o.run()

	return err
}

func (o *options) getToken() (string, error) {
	client, ns, err := o.CommonOptions.KubeClientAndDevNamespace()
	if err != nil {
		return "", errors.Wrap(err, "failed to get a kubernetes client")
	}

	listOpts := metav1.ListOptions{
		LabelSelector: githubappLabel + o.organisation,
	}
	gitHubAppSecrets, err := client.CoreV1().Secrets(ns).List(listOpts)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get secrets for organisation %s", o.organisation)
	}

	if len(gitHubAppSecrets.Items) != 1 {
		return "", fmt.Errorf("found %v github app secrets with label oragisation %s, should be 1", len(gitHubAppSecrets.Items), o.organisation)
	}

	secret := gitHubAppSecrets.Items[0]
	token := secret.Data[githubPassword]
	return string(token), nil
}

func (o *options) run() (string, error) {

	e := exec.Command(o.Runner.CurrentName(), o.Runner.CurrentArgs()...)
	e.Stdout = o.Out
	e.Stderr = o.Err
	os.Setenv("PATH", util.PathWithBinary())
	err := e.Run()
	if err != nil {
		log.Logger().Errorf("Error: Command failed  %s %s", o.Runner.CurrentName(), strings.Join(o.Runner.CurrentArgs(), " "))
	}

	if len(o.Runner.CurrentEnv()) > 0 {
		m := map[string]string{}
		environ := os.Environ()
		for _, kv := range environ {
			paths := strings.SplitN(kv, "=", 2)
			if len(paths) == 2 {
				m[paths[0]] = paths[1]
			}
		}
		for k, v := range o.Runner.CurrentEnv() {
			m[k] = v
		}
		envVars := []string{}
		for k, v := range m {
			envVars = append(envVars, k+"="+v)
		}
		e.Env = envVars
	}

	e.Stdout = o.Out
	e.Stderr = o.Err

	var text string

	err = e.Run()
	if err != nil {
		if err != nil {
			errors.Wrapf(err, "failed to run command")
		}
	}

	return text, err
}
