// Copyright 2022 PingCAP, Inc.
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

package core_test

import (
	"testing"

	"github.com/pingcap/tidb/planner/core"
	"github.com/pingcap/tidb/testkit"
	"github.com/pingcap/tidb/testkit/testdata"
)

func TestIndexMergeJSONMemberOf(t *testing.T) {
	store := testkit.CreateMockStore(t)
	tk := testkit.NewTestKit(t, store)
	tk.MustExec("use test")
	tk.MustExec(`create table t(
a int, j0 json, j1 json,
index j0_0((cast(j0->'$.path0' as signed array))),
index j0_1((cast(j0->'$.path1' as signed array))),
index j1((cast(j1 as signed array))))`)

	var input []string
	var output []struct {
		SQL  string
		Plan []string
	}
	planSuiteData := core.GetIndexMergeSuiteData()
	planSuiteData.LoadTestCases(t, &input, &output)

	for i, query := range input {
		testdata.OnRecord(func() {
			output[i].SQL = query
		})
		result := tk.MustQuery("explain format = 'brief' " + query)
		testdata.OnRecord(func() {
			output[i].Plan = testdata.ConvertRowsToStrings(result.Rows())
		})
		result.Check(testkit.Rows(output[i].Plan...))
	}
}
