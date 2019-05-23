package postmortem

import (
	"bytes"
	"github.com/google/sredocs/parser"
	"log"
)

const (
	shortLinkStr     = `Short Link:`
	lastUpdatedStr   = `Last Updated:`
	teamNameStr      = `Team Name:`
	collaboratorsStr = `Collaborators:`
	statusStr        = `Status:`

	severityStr    = `Minor, Medium or High Severity:`
	impactStr      = `Impact:`
	descriptionStr = `Incident Description:`

	ttdStr                 = `Time to detect in minutes:`
	tteStr                 = `Time to initiate response in minutes:`
	ttmStr                 = `Time to mitigate in minutes:`
	SLOLinkStr             = `Link to impacted SLO:`
	impactedProductsStr    = `Impacted products:`
	firstServiceStr        = `First known impacted service:`
	blastRadiusStr         = `Known services in the blast radius:`
	noteworthyCustomersStr = `Noteworthy customers impacted:`
	triggerStr             = `Deploy, Cloud or Other Trigger:`

	// TODO(stratus): Parse AIs table.

	backgroundStr = `Background`
	wellStr       = `Things that went well`
	luckyStr      = `Where we got lucky`
	improvedStr   = `Things that could be improved`
	timelineStr   = `Timeline`

	footerStr = `DO NOT REMOVE THIS AND THE CONTENT BELOW`
)

var (
	// must be kept sorted by how it's expected to appear in the doc.
	Fields = []string{shortLinkStr, lastUpdatedStr, teamNameStr, collaboratorsStr, statusStr,
		severityStr, impactStr, descriptionStr, ttdStr, tteStr, ttmStr, SLOLinkStr,
		impactedProductsStr, firstServiceStr, blastRadiusStr, noteworthyCustomersStr,
		triggerStr, backgroundStr, wellStr, luckyStr, improvedStr, timelineStr, footerStr}
)

func Parse(fields []string, b []byte) (*bytes.Buffer, error) {
	log.Println("postmortem")
	p := &parser.DefaultParser{}
	csv, err := p.Parse(fields, b)
	if err != nil {
		return nil, err
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
