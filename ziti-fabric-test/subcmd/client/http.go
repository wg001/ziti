/*
	Copyright 2019 Netfoundry, Inc.

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

package client

import (
	"github.com/netfoundry/ziti-cmd/ziti-fabric-test/subcmd"
	"github.com/netfoundry/ziti-fabric/fabric/xgress_transport"
	"github.com/netfoundry/ziti-foundation/identity/dotziti"
	"github.com/netfoundry/ziti-foundation/identity/identity"
	"github.com/netfoundry/ziti-foundation/transport"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/michaelquigley/pfxlog"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net"
	"net/http"
)

func init() {
	httpCmd.Flags().StringVarP(&httpCmdIdentity, "identityName", "i", "default", "dotzeet identity name")
	httpCmd.Flags().StringVarP(&httpCmdIngress, "ingressEndpoint", "e", "tls:127.0.0.1:7001", "ingress endpoint address")
	httpCmd.Flags().StringVar(&httpCmdHost, "host", "", "optional host header")
	httpCmd.Flags().BoolVar(&httpCmdInsecure, "insecure", false, "Disable SSL security checks")
	subcmd.Root.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   "http <http[s]://service/path>",
	Short: "Simple HTTP client",
	Args:  cobra.ExactArgs(1),
	Run:   doHttp,
}
var httpCmdIdentity string
var httpCmdIngress string
var httpCmdHost string
var httpCmdInsecure bool

func doHttp(cmd *cobra.Command, args []string) {
	if _, id, err := dotziti.LoadIdentity(httpCmdIdentity); err == nil {
		if ingressAddr, err := transport.ParseAddress(httpCmdIngress); err == nil {
			tr := &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
					host, _, err := net.SplitHostPort(addr)
					if err != nil {
						return nil, err
					}
					serviceId := &identity.TokenId{Token: host}
					if peer, err := xgress_transport.ClientDial(ingressAddr, id, serviceId); err == nil {
						pfxlog.Logger().Debug("connected")
						return peer.Conn(), nil
					} else {
						return nil, err
					}
				},
			}
			if httpCmdInsecure {
				pfxlog.Logger().Warn("disabled SSL security checks")
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
			c := &http.Client{Transport: tr}
			request, err := http.NewRequest("GET", args[0], nil)
			if err != nil {
				panic(err)
			}
			if httpCmdHost != "" {
				pfxlog.Logger().Infof("set host header to [%s]", httpCmdHost)
				request.Host = httpCmdHost
			}
			if response, err := c.Do(request); err == nil {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					panic(err)
				}
				if err := response.Body.Close(); err != nil {
					panic(err)
				}
				fmt.Println(string(body))
			} else {
				panic(err)
			}

		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}
