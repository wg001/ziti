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

package edge_controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/Jeffail/gabs"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/common"
	cmdutil "github.com/openziti/ziti/ziti/cmd/ziti/cmd/factory"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/spf13/cobra"
)

type createConfigTypeOptions struct {
	commonOptions
	schemaFile string
}

// newCreateConfigTypeCmd creates the 'edge controller create service-policy' command
func newCreateConfigTypeCmd(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &createConfigTypeOptions{
		commonOptions: commonOptions{
			CommonOptions: common.CommonOptions{Factory: f, Out: out, Err: errOut},
		},
	}

	cmd := &cobra.Command{
		Use:   "config-type <name> <JSON schema>",
		Short: "creates a config type managed by the Ziti Edge Controller",
		Long:  "creates a config type managed by the Ziti Edge Controller",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := runCreateConfigType(options)
			cmdhelper.CheckErr(err)
		},
		SuggestFor: []string{},
	}

	// allow interspersing positional args and flags
	cmd.Flags().SetInterspersed(true)
	cmd.Flags().StringVarP(&options.schemaFile, "schema-file", "f", "", "Read config type JSON schema from a file instead of the command line")
	cmd.Flags().BoolVarP(&options.OutputJSONResponse, "output-json", "j", false, "Output the full JSON response from the Ziti Edge Controller")

	return cmd
}

// runCreateConfigType create a new configType on the Ziti Edge Controller
func runCreateConfigType(o *createConfigTypeOptions) error {
	var schemaMap map[string]interface{}

	var schemaBytes []byte

	if len(o.Args) == 2 {
		schemaBytes = []byte(o.Args[1])
	}

	if o.Cmd.Flags().Changed("schema-file") {
		if len(o.Args) == 2 {
			return errors.New("schema specified both in file and on command line. please pick one")
		}
		var err error
		if schemaBytes, err = ioutil.ReadFile(o.schemaFile); err != nil {
			return fmt.Errorf("failed to read schema file %v: %w", o.schemaFile, err)
		}
	}

	if len(schemaBytes) > 0 {
		schemaMap = map[string]interface{}{}
		if err := json.Unmarshal(schemaBytes, &schemaMap); err != nil {
			fmt.Printf("Attempted to parse: %v\n", string(schemaBytes))
			fmt.Printf("Failing parsing JSON: %+v\n", err)
			return errors.Errorf("unable to parse data as json: %v", err)
		}
	}

	entityData := gabs.New()
	setJSONValue(entityData, o.Args[0], "name")
	if schemaMap != nil {
		setJSONValue(entityData, schemaMap, "schema")
	}

	result, err := createEntityOfType("config-types", entityData.String(), &o.commonOptions)

	if err != nil {
		panic(err)
	}

	configTypeId := result.S("data", "id").Data()

	if _, err := fmt.Fprintf(o.Out, "%v\n", configTypeId); err != nil {
		panic(err)
	}

	return err
}
