package awx

import (
	"fmt"

	"get.porter.sh/porter/pkg/exec/builder"
	yaml "gopkg.in/yaml.v2"
)

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the awx mixin in porter.yaml
// mixins:
// - awx:
//	  clientVersion: "v0.0.0"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`
}

// This is an example. Replace the following with whatever steps are needed to
// install required components into
const dockerfileLines = `RUN apt-get install curl --yes`

// I need the awx tooling to invoke a recipe on the tower
// awx-cli
// https://github.com/awx/awx/blob/devel/INSTALL.md#installing-the-awx-cli
// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build() error {

	// Create new Builder.
	var input BuildInput

	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	suppliedClientVersion := input.Config.ClientVersion

	if suppliedClientVersion != "" {
		m.ClientVersion = suppliedClientVersion
	}

	fmt.Fprintf(m.Out, dockerfileLines)
	fmt.Fprintf(m.Out, "\nRUN curl -L https://files.pythonhosted.org/packages/4b/fe/bfe6efb7f5a57c24ed6ba9b670ce3c9d15dd8a5b148cb87f8fafe559bef6/awxkit-19.2.2-py3-none-any.whl -o awxkit-19.2.2-py3-none-any.whl")
	fmt.Fprintf(m.Out, "\nRUN /usr/local/bin/python -m pip install --upgrade pip")
	fmt.Fprintf(m.Out, "\nRUN ls")
	fmt.Fprintf(m.Out, "\nRUN pip3 install ./awxkit-19.2.2-py3-none-any.whl")
	// Example of pulling and defining a client version for your mixin
	// fmt.Fprintf(m.Out, "\nRUN curl https://get.helm.sh/helm-%s-linux-amd64.tar.gz --output helm3.tar.gz", m.ClientVersion)

	return nil
}
