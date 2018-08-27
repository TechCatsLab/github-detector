/*
 * Revision History:
 *     Initial: 2018/08/02        Li Zebang
 */

package main

import (
	"github.com/TechCatsLab/logging/logrus"

	"github.com/fengyfei/github-detector/cmd/github-detector/app"
)

func main() {
	command := app.NewDetector()

	if err := command.Execute(); err != nil {
		logrus.Fatalf("%v", err)
	}
}
