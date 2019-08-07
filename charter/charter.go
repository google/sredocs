// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package charter

import (
	"bytes"
	"github.com/google/sredocs/parser"
	"log"
)

const (
	lastUpdatedStr          = `Last Updated:`
	companyNameStr          = `Company Name:`
	teamNameStr             = `Team Name:`
	collaboratorsStr        = `Collaborators:`
	approversStr            = `Approvers:`
	statusStr               = `Status:`
	whoAreWeStr             = `Who Are We`
	servicesSupportedStr    = `Services Supported`
	howDoWeInvestOurTimeStr = `How Do We Invest Our Time`
	teamValuesStr           = `Team Values`
	footerStr               = `DO NOT REMOVE THIS AND THE CONTENT BELOW`
)

var (
	// must be kept sorted by how it's expected to appear in the doc.
	Fields = []string{lastUpdatedStr, companyNameStr, teamNameStr,
		collaboratorsStr, approversStr, statusStr, whoAreWeStr, servicesSupportedStr,
		howDoWeInvestOurTimeStr, teamValuesStr, footerStr}
)

func Parse(fields []string, b []byte) (*bytes.Buffer, error) {
	log.Println("charter")
	p := &parser.DefaultParser{}
	csv, err := p.Parse(fields, b)
	if err != nil {
		return csv, err
	}
	return csv, nil
}

func Save(b *bytes.Buffer, filename string) error {
	p := &parser.DefaultParser{}
	err := p.Save(b, filename)
	if err != nil {
		return err
	}
	return nil
}
