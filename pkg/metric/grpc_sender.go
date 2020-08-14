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
	"encoding/json"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	pb "github.com/ZupIT/ritchie-cli/internal/proto"
)

var _ Sender = SendManagerRpc{}

type SendManagerRpc struct {
	rpcClient *grpc.ClientConn
}

//  Example of how to create a grpc client:
//  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	conn, _ := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
func NewRpcSender(rpcClient *grpc.ClientConn) SendManagerRpc {
	return SendManagerRpc{rpcClient: rpcClient}
}

func (sm SendManagerRpc) Send(dataset Dataset) {
	defer sm.rpcClient.Close()

	c := pb.NewProcessorClient(sm.rpcClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, _ = c.Process(ctx, convert(dataset))
}

func convert(dataset Dataset) *pb.DatasetRequest {
	timestamp, _ := ptypes.TimestampProto(dataset.Timestamp)
	data, _ := json.Marshal(dataset.Data)

	return &pb.DatasetRequest{
		MetricId:   dataset.Id.String(),
		UserId:     dataset.UserId.String(),
		Timestamp:  timestamp,
		So:         dataset.So,
		RitVersion: dataset.RitVersion,
		Data:       data,
	}
}
