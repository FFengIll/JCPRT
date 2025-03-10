package main

import (
	"encoding/xml"
)

type ReleaseDTO struct {
	ID      int    `json:"id" xml:"-"`
	File    string `json:"file" xml:"-"`
	Since   string `json:"since" xml:"since-build,attr"`
	Until   string `json:"until" xml:"until-build,attr"`
	Version string `json:"version" xml:"-"`
}

type PluginDTO struct {
	ID      int        `json:"id" xml:"-"`
	Name    string     `json:"name" xml:"-"`
	XMLID   string     `json:"xmlId" xml:"id,attr"`
	Version string     `json:"version" xml:"version,attr"`
	Release ReleaseDTO `json:"releases" xml:"idea-version"`
	URL     string     `json:"url" xml:"url,attr"`
}

type Plugins struct {
	XMLName xml.Name     `xml:"plugins"`
	Plugins []*PluginDTO `xml:"plugin"`
}

// PluginXML represents the structure of plugin.xml
type PluginXML struct {
	ID          string `xml:"id"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Version     string `xml:"version"`
	IdeaVersion struct {
		Since string `xml:"since-build,attr"`
		Until string `xml:"until-build,attr"`
	} `xml:"idea-version"`
	Depends []string `xml:"depends"`
}
