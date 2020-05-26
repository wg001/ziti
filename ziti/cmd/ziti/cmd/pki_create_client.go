/*
	Copyright NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"

	cmdutil "github.com/openziti/ziti/ziti/cmd/ziti/cmd/factory"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/openziti/ziti/ziti/cmd/ziti/internal/log"
	"github.com/openziti/ziti/ziti/pki/certificate"
	"github.com/openziti/ziti/ziti/pki/pki"
	"github.com/openziti/ziti/ziti/pki/store"
)

// PKICreateClientOptions the options for the create spring command
type PKICreateClientOptions struct {
	PKICreateOptions
}

// NewCmdPKICreateClient creates a command object for the "create" command
func NewCmdPKICreateClient(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &PKICreateClientOptions{
		PKICreateOptions: PKICreateOptions{
			PKIOptions: PKIOptions{
				CommonOptions: CommonOptions{
					Factory: f,
					Out:     out,
					Err:     errOut,
				},
			},
		},
	}

	cmd := &cobra.Command{
		Use:     "client",
		Short:   "Creates new Client certificate (signed by previously created Intermediate-chain)",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			cmdhelper.CheckErr(err)
		},
	}

	options.addPKICreateClientFlags(cmd)
	return cmd
}

func (o *PKICreateClientOptions) addPKICreateClientFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Flags.PKIRoot, "pki-root", "", "", "Directory in which PKI resides")
	cmd.Flags().StringVarP(&o.Flags.CAName, "ca-name", "", "intermediate", "Name of Intermediate CA (within PKI_ROOT) to use to sign the new Client certificate")
	cmd.Flags().StringVarP(&o.Flags.ClientFile, "client-file", "", "client", "Name of file (under chosen CA) in which to store new Client certificate and private key")
	cmd.Flags().StringVarP(&o.Flags.KeyFile, "key-file", "", "", "Name of file (under chosen CA) containing private key to use when generating Client certificate")
	cmd.Flags().StringVarP(&o.Flags.ClientName, "client-name", "", "NetFoundry, Inc. Client", "Common Name (CN) to use for new Client certificate")
	cmd.Flags().StringSliceVar(&o.Flags.Email, "email", []string{}, "Email addr(s) to add to Subject Alternate Name (SAN) for new Client certificate")
	cmd.Flags().IntVarP(&o.Flags.CAExpire, "expire-limit", "", 365, "Expiration limit in days")
	cmd.Flags().IntVarP(&o.Flags.CAMaxpath, "max-path-len", "", -1, "Intermediate maximum path length")
	cmd.Flags().IntVarP(&o.Flags.CAPrivateKeySize, "private-key-size", "", 2048, "Size of the private key")
}

// Run implements this command
func (o *PKICreateClientOptions) Run() error {

	pkiroot, err := o.ObtainPKIRoot()
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	o.Flags.PKI = &pki.ZitiPKI{Store: &store.Local{}}
	local := o.Flags.PKI.Store.(*store.Local)
	local.Root = pkiroot

	commonName := o.Flags.ClientName

	clientCertFile, err := o.ObtainClientCertFile()
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	filename := o.ObtainFileName(clientCertFile, commonName)
	template := o.ObtainPKIRequestTemplate(commonName)

	template.IsCA = false
	template.EmailAddresses = o.Flags.Email

	var signer *certificate.Bundle

	caname, err := o.ObtainCAName(pkiroot)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	keyFile, err := o.ObtainKeyFile(false)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	signer, err = o.Flags.PKI.GetCA(caname)
	if err != nil {
		return fmt.Errorf("Cannot locate signer: %v", err)
	}

	req := &pki.Request{
		Name:                filename,
		KeyName:             keyFile,
		Template:            template,
		IsClientCertificate: true,
		PrivateKeySize:      o.Flags.CAPrivateKeySize,
	}

	if err := o.Flags.PKI.Sign(signer, req); err != nil {
		return fmt.Errorf("Cannot Sign: %v", err)
	}

	log.Infoln("Success")

	return nil

}
