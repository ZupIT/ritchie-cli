/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metric

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	pb "github.com/ZupIT/ritchie-cli/internal/proto"
)

func TestSend(t *testing.T) {
	type in struct {
		client  pb.ProcessorClient
		dataset APIData
	}

	tests := []struct {
		name string
		in   in
	}{
		{
			name: "success",
			in: in{
				client: grpcProcessMock{},
				dataset: APIData{
					Id:         "metric-id",
					UserId:     "user-id",
					Timestamp:  time.Now(),
					Os:         runtime.GOOS,
					RitVersion: "2.0.0",
					Data:       nil,
				},
			},
		},
		{
			name: "ignore error",
			in: in{
				client: grpcProcessMock{err: errors.New("error to send metric")},
				dataset: APIData{
					Id:         "metric-id",
					UserId:     "user-id",
					Timestamp:  time.Now(),
					Os:         runtime.GOOS,
					RitVersion: "2.0.0",
					Data:       nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rpcSender := NewRpcSender(tt.in.client)
			rpcSender.Send(tt.in.dataset)
		})
	}

}

type grpcProcessMock struct {
	err error
}

func (g grpcProcessMock) Process(ctx context.Context, in *pb.DatasetRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return &empty.Empty{}, g.err
}
