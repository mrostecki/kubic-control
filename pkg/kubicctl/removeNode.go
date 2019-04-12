// Copyright 2019 Thorsten Kukuk
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubicctl

import (
	"context"
	"time"
	"fmt"

        log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pb "github.com/thkukuk/kubic-control/api"
)

func RemoveNodeCmd() *cobra.Command {
        var subCmd = &cobra.Command {
                Use:   "remove <node>",
                Short: "Remove node form cluster",
                Run: removeNode,
		Args: cobra.ExactArgs(1),
	}

	return subCmd
}

func removeNode(cmd *cobra.Command, args []string) {
	// Set up a connection to the server.

	nodes := args[0]

	conn, err := CreateConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	c := pb.NewKubeadmClient(conn)

	// var deadlineMin = flag.Int("deadline_min", 10, "Default deadline in minutes.")
	// clientDeadline := time.Now().Add(time.Duration(*deadlineMin) * time.Minute)
	// ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.RemoveNode(ctx, &pb.RemoveNodeRequest{NodeNames: nodes})
	if err != nil {
		log.Errorf("could not initialize: %v", err)
		return
	}
	if r.Success {
		fmt.Printf("Node %s removed\n", nodes)
	} else {
		log.Errorf("Removing node %s failed: %s", nodes, r.Message)
	}
}