// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package master

import (
	"fmt"

	"github.com/pingcap/dm/dm/ctl/common"
	"github.com/pingcap/dm/dm/pb"
	"github.com/pingcap/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// NewStartTaskCmd creates a StartTask command
func NewStartTaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-task [-w worker ...] <config_file>",
		Short: "start a task with config file",
		Run:   startTaskFunc,
	}
	return cmd
}

// startTaskFunc does start task request
func startTaskFunc(cmd *cobra.Command, _ []string) {
	if len(cmd.Flags().Args()) != 1 {
		fmt.Println(cmd.Usage())
		return
	}
	content, err := common.GetFileContent(cmd.Flags().Arg(0))
	if err != nil {
		common.PrintLines("get file content error:\n%v", errors.ErrorStack(err))
		return
	}

	workers, err := common.GetWorkerArgs(cmd)
	if err != nil {
		common.PrintLines("%s", errors.ErrorStack(err))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start task
	cli := common.MasterClient()
	resp, err := cli.StartTask(ctx, &pb.StartTaskRequest{
		Task:    string(content),
		Workers: workers,
	})
	if err != nil {
		common.PrintLines("can not start task:\n%v", errors.ErrorStack(err))
		return
	}

	common.PrettyPrintResponse(resp)
}
