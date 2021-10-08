package scraper

import (
	"archive/zip"
	"encoding/xml"
	"strings"
)

type OfficeCoreProperty struct {
	XMLName			xml.Name	`xml:"coreProperties"`
	Creator			string		`xml:"creator"`
	LastModifiedBy	string		`xml:"lastModifiedBy"`
}

type OfficeAppProperty struct {
	XMLName			xml.Name	`xml:"Properties"`
	Application		string		`xml:"Application"`
	Company			string		`xml:"Company"`
	AppVersion		string		`xml:"AppVersion"`
}

var OfficeVersions = map[string]string{
	"16":	"2016",
	"15":	"2013",
	"14":	"2010",
	"12":	"2007",
	"11":	"2003",
}

// GetMajorVersion computes the human readable office version from the version code
func (a *OfficeAppProperty) GetMajorVersion() string {
	tokens := strings.Split(a.AppVersion, ".")
	if len(tokens) < 2 {
		return "Unknown"
	}

	v, ok := OfficeVersions[tokens[0]]
	if !ok {
		return "Unknown"
	}
	return v
}

// NewProperties parses a zip document and
// returns the attributes specified in the OfficeProperty structs
func NewProperties(r * zip.Reader) (*OfficeAppProperty, *OfficeCoreProperty, error) {
	var appProperty OfficeAppProperty
	var coreProperty OfficeCoreProperty
	for _, f := range r.File{
		switch f.Name {
		case "docProp/core.xml":
			if err := process(f, &coreProperty); err != nil {
				return nil, nil, err
			}
		case "docProp/app.xml":
			if err := process(f, &appProperty); err!= nil {
				return nil, nil, err
			}
		default:
			continue
		}
	}
	return &appProperty, &coreProperty, nil
}


func process(f *zip.File, prop interface{}) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}

	defer rc.Close()

	if err := xml.NewDecoder(rc).Decode(&prop); err != nil {
		return err
	}
	return nil
}

