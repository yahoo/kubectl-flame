package utils

import (
	"bufio"
	"encoding/base64"
	"github.com/VerizonMedia/kubectl-flame/api"
	"io/ioutil"
	"os"
)

func PublishFlameGraph(flameFile string) error {
	file, err := os.Open(flameFile)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	fgData := api.FlameGraphData{EncodedFile: encoded}

	return api.PublishEvent(api.FlameGraph, fgData)
}
