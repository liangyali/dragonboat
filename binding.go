// Copyright 2017-2019 Lei Ni (nilei81@gmail.com)
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

// +build dragonboat_language_binding

package dragonboat

import (
	"time"

	"github.com/lni/dragonboat/client"
)

//
// Public methods in this file are used by language bindings, they are
// considered as a part of the internal interface that can change any
// time. Applications are not suppose to directly call any of them.
//

func (rr *RequestResult) GetCode() RequestResultCode {
	return rr.code
}

// ProposeCH is similar to the Propose method with an extra ICompleteHandler
// specified. On proposal's completion, the ICompleteHandler will be invoked
// by the system. The ICompleteHandler instance should not be used as a
// general callback function, it should only be used to notify the completion
// of the proposal. ProposeWithCH is mainly used by language bindings to
// implement async proposals, Go applications are expected to use the Propose
// method.
func (nh *NodeHost) ProposeCH(s *client.Session,
	data []byte, handler ICompleteHandler,
	timeout time.Duration) (*RequestState, error) {
	return nh.propose(s, data, handler, timeout)
}

// ReadIndexCH is similar to the ReadIndex method with an extra
// ICompleteHandler specified. On completion of the ReadIndex operation, the
// ICompleteHandler will be invoked by the system. The ICompleteHandler should
// not be used as a general callback function, it should only be used to notify
// the completion of the ReadIndex operation.
// ReadIndexCH is mainly used by language bindings to implement async
// ReadIndex operations, Go applications are expected to use the ReadIndex
// method.
func (nh *NodeHost) ReadIndexCH(clusterID uint64,
	handler ICompleteHandler,
	timeout time.Duration) (*RequestState, error) {
	rs, _, err := nh.readIndex(clusterID, handler, timeout)
	return rs, err
}

// ProposeSessionCH is similar to the ProposeSession method but with an extra
// ICompleteHandler specified as input parameter. The ICompleteHandler should
// not be used as a general callback function, it should only be used to notify
// the completion of the propose session operation.
func (nh *NodeHost) ProposeSessionCH(s *client.Session,
	handler ICompleteHandler, timeout time.Duration) (*RequestState, error) {
	v, ok := nh.getCluster(s.ClusterID)
	if !ok {
		return nil, ErrClusterNotFound
	}
	req, err := v.proposeSession(s, handler, timeout)
	nh.execEngine.setNodeReady(s.ClusterID)
	return req, err
}
