package charter

import (
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

func Parse(fields []string, b []byte) (string, error) {
	log.Println("charter")
	p := &parser.DefaultParser{}
	csv, err := p.Parse(fields, b)
	if err != nil {
		return csv, err
	}
	return csv, nil
}
