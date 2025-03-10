package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"testing"
)

func TestPlugins(t *testing.T) {
	plugins := Plugins{
		Plugins: []*PluginDTO{
			{
				ID:      1,
				Name:    "CSV Editor",
				XMLID:   "intellij.plugins.id",
				Version: "test",
				URL:     "test.jar",
				Release: ReleaseDTO{
					File:    "CSVEditor-3.4.0-242.zip",
					Since:   "242.0",
					Until:   "242.*",
					Version: "3.4.0-242",
				},
			},
		},
	}

	output, err := xml.MarshalIndent(plugins, "", "  ")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	fmt.Printf("%s%s\n", xml.Header, output)
}
