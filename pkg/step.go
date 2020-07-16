package pkg

import (
	"os"
	"os/exec"
	"strings"

	"github.com/jenkins-x/jx/v2/pkg/util"

	"github.com/jenkins-x/jx/v2/pkg/cmd/clients"
	"github.com/jenkins-x/jx/v2/pkg/cmd/opts"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	githubappLabel             = "jenkins.io/githubapp-owner="
	ownerAnnotation            = "jenkins.io/githubapp-owner"
	nonGithubAppSecretSelector = "jenkins.io/kind=git,jenkins.io/service-kind=github"
	githubPassword             = "password"
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

	selector := githubappLabel + o.organisation
	listOpts := metav1.ListOptions{
		LabelSelector: selector,
	}
	secretInterface := client.CoreV1().Secrets(ns)
	secrets, err := secretInterface.List(listOpts)
	if err != nil && !apierrors.IsNotFound(err) {
		return "", errors.Wrapf(err, "failed to get secrets for selector: %s", selector)
	}

	for _, s := range secrets.Items {
		token := s.Data[githubPassword]
		if len(token) > 0 {
			return string(token), nil
		}
	}

	// lets try find a non-github app secret
	listOpts = metav1.ListOptions{
		LabelSelector: nonGithubAppSecretSelector,
	}

	secrets, err = secretInterface.List(listOpts)
	if err != nil && !apierrors.IsNotFound(err) {
		return "", errors.Wrapf(err, "failed to get secrets for selector: %s", nonGithubAppSecretSelector)
	}
	for _, s := range secrets.Items {
		token := s.Data[githubPassword]
		if len(token) > 0 {
			return string(token), nil
		}
	}
	return "", errors.Errorf("could not find a secret for selector %s or %s", selector, nonGithubAppSecretSelector)
}

func (o *options) run() (string, error) {

	e := exec.Command(o.Runner.CurrentName(), o.Runner.CurrentArgs()...)
	e.Stdout = o.Out
	e.Stderr = o.Err
	os.Setenv("PATH", util.PathWithBinary())

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

	err := e.Run()
	if err != nil {
		if err != nil {
			errors.Wrapf(err, "failed to run command")
		}
	}

	return text, err
}
