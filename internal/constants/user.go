/*
 * Apache License
 * Version 2.0, January 2004
 * http://www.apache.org/licenses/
 *
 * TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION
 *
 * Copyright 2025 Bamboo
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

package constants

import "errors"

// Http Error Code
const (
	UserSignUpFailedErrorCode = 400001 + iota
	UserExistErrorCode
	UserNotExistErrorCode
)

var (
	// UserService
	ErrorUserExist         = errors.New("user already exists, check your username or mobile, or try to login")
	ErrorUserNotExist      = errors.New("user not exists")
	ErrorUserSignUpFail    = errors.New("user sign up fail")
	ErrorPasswordIncorrect = errors.New("user password incorrect")

	// TreeService
	// Node DAO
	ErrorTreeNodeNotExist = errors.New("tree node not exists")

	// ECS DAO
	ErrorResourceEcsExist = errors.New("resource ecs already exists")
)
