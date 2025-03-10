package main

import (
	"fmt"
	"testing"
)

func TestUploader_parseJar(t *testing.T) {
	u := &Uploader{}
	pluginXML, err := u.Parse("/Users/yz/Project/tt-project/synapse_cov/plugin/synapse_cov_intellij/build/distributions/synapse_cov_intellij-1.2.1-SNAPSHOT.zip")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", pluginXML)
}
