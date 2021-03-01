// Copyright [2020] [thinkgos] thinkgo@aliyun.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

// Discard is an logger on which all Write calls succeed
// without doing anything.
type Discard struct{}

var _ Logger = (*Discard)(nil)

// NewDiscard a discard logger on which always succeed without doing anything
func NewDiscard() Discard { return Discard{} }

// Debugf implement Logger interface.
func (sf Discard) Debugf(string, ...interface{}) {}

// Infof implement Logger interface.
func (sf Discard) Infof(string, ...interface{}) {}

// Errorf implement Logger interface.
func (sf Discard) Errorf(string, ...interface{}) {}

// Warnf implement Logger interface.
func (sf Discard) Warnf(string, ...interface{}) {}

// DPanicf implement Logger interface.
func (sf Discard) DPanicf(string, ...interface{}) {}

// Fatalf implement Logger interface.
func (sf Discard) Fatalf(string, ...interface{}) {}
