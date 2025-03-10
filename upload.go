package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
)

// PluginConfig declare a plugin file hosted in `object storage`
// host key is the storage key for `updatePlugins.xml`
// plugins are filepath for each plugin jar dist file
type PluginConfig struct {
	HostKey   string   `json:"host"`
	Plugins   []string `json:"plugins"`
	Overwrite bool     `json:"overwrite"`
}

type Uploader struct {
	config   PluginConfig
	provider Provider
}

func (u *Uploader) Run() error {
	repoXML, err := u.Build()
	if err != nil {
		return err
	}

	path := filepath.Join(u.config.HostKey, "updatePlugins.xml")
	_, err = u.UploadXML(repoXML, path)
	if err != nil {
		return err
	}

	return nil
}

// Build the `updatePlugins.xml`
func (u *Uploader) Build() (string, error) {
	plugins := new(Plugins)

	for _, plugin := range u.config.Plugins {
		pluginDTO, err := u.Parse(plugin)
		if err != nil {
			fmt.Println(err)
			continue
		}
		path, err := u.UploadPlugin(pluginDTO)
		if err != nil {
			fmt.Println(err)
			continue
		}
		println(path)

		plugins.Plugins = append(plugins.Plugins, pluginDTO)
	}

	bs, err := xml.MarshalIndent(plugins, "", "  ")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	return fmt.Sprintf("%s%s\n", xml.Header, bs), nil
}

func (u *Uploader) buildPluginXml(plugin *PluginDTO) string {
	output, err := xml.MarshalIndent(plugin, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	return string(output)
}

// UploadXML store the `updatePlugins.xml` to remote
func (u *Uploader) UploadXML(content, path string) (string, error) {
	return "", nil

}

// UploadPlugin store the plugin jar file to remote
func (u *Uploader) UploadPlugin(file *PluginDTO) (string, error) {
	return "", nil
}

func (u *Uploader) GenerateID() int {
	return 0
}

// Parse the plugin jar file to get its metadata like description, deps and so on.
func (u *Uploader) Parse(file string) (*PluginDTO, error) {

	pluginXML, err := u.parseJar(file, nil)
	if err != nil {
		return nil, err
	}

	plugin := &Plugin{}
	r := &PluginRelease{}

	// populate info in xml
	plugin.ID = u.GenerateID()
	plugin.Name = pluginXML.Name
	plugin.XMLID = pluginXML.ID
	r.Version = pluginXML.Version

	// Create Release DTO
	release := ReleaseDTO{}
	release.ID = r.ID
	release.File = r.File
	release.Since = r.Since
	release.Until = r.Until
	release.Version = r.Version

	// Create pluginDTO
	pluginDto := PluginDTO{}
	pluginDto.ID = plugin.ID
	pluginDto.Name = plugin.Name
	pluginDto.XMLID = plugin.XMLID
	pluginDto.Release = release
	pluginDto.Version = r.Version
	pluginDto.URL = r.File

	return &pluginDto, nil
}

// parseJar is a helper function to recursively parse jar files
func (u *Uploader) parseJar(file string, parentReader *zip.Reader) (*PluginXML, error) {
	var r *zip.Reader

	if parentReader == nil {
		// Open the jar file
		zipFile, err := zip.OpenReader(file)
		if err != nil {
			return nil, err
		}
		defer zipFile.Close()
		r = &zipFile.Reader
	} else {
		r = parentReader
	}

	// Find and read the plugin.xml file or nested jar files
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "plugin.xml") {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			// Parse the plugin.xml file
			var pluginXML PluginXML
			if err := xml.NewDecoder(rc).Decode(&pluginXML); err != nil {
				return nil, err
			}

			return &pluginXML, nil
		} else if strings.HasSuffix(f.Name, ".jar") {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			// Read the nested jar file into memory
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, rc); err != nil {
				return nil, err
			}

			nestedReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
			if err != nil {
				return nil, err
			}

			// Recursively parse the nested jar file
			nestedPluginXML, err := u.parseJar("", nestedReader)
			if err == nil {
				return nestedPluginXML, nil
			}
		}
	}

	return nil, fmt.Errorf("plugin.xml not found in %s", file)
}
