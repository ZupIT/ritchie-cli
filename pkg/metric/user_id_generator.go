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
	"crypto/sha256"
	"encoding/hex"
	"hash"

	"github.com/denisbrodbeck/machineid"
)

var _ UserIdGenerator = UserIdManager{}

type UserIdManager struct {
	hash hash.Hash
}

func NewUserIdGenerator() UserIdManager {
	return UserIdManager{hash: sha256.New()}
}

func (us UserIdManager) Generate() (UserId, error) {
	id, err := machineid.ID()
	if err != nil {
		return "", err
	}

	us.hash.Reset()
	if _, err := us.hash.Write([]byte(id)); err != nil {
		return "", err
	}
	userId := hex.EncodeToString(us.hash.Sum(nil))

	return UserId(userId), nil
}
