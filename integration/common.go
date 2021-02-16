// +build linux

/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package integration

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	cri "k8s.io/cri-api/pkg/apis"
)

// ImageList holds public image references
type ImageList struct {
	Alpine          string
	BusyBox         string
	Pause           string
	VolumeCopyUp    string
	VolumeOwnership string
}

var (
	imageService cri.ImageManagerService
	imageMap     map[int]string
	imageList    ImageList
	pauseImage   string // This is the same with default sandbox image
)

func initImages(imageListFile string) {
	imageList = ImageList{
		Alpine:          "docker.io/library/alpine:latest",
		BusyBox:         "docker.io/library/busybox:latest",
		Pause:           "k8s.gcr.io/pause:3.5",
		VolumeCopyUp:    "gcr.io/k8s-cri-containerd/volume-copy-up:2.0",
		VolumeOwnership: "gcr.io/k8s-cri-containerd/volume-ownership:2.0",
	}

	if imageListFile != "" {
		fileContent, err := ioutil.ReadFile(imageListFile)
		if err != nil {
			panic(fmt.Errorf("Error reading '%v' file contents: %v", imageList, err))
		}

		err = toml.Unmarshal(fileContent, &imageList)
		if err != nil {
			panic(fmt.Errorf("Error unmarshalling '%v' TOML file: %v", imageList, err))
		}
	}

	logrus.Infof("Using the following image list: %+v", imageList)

	imageMap = initImageMap(imageList)
	pauseImage = GetImage(Pause)
}

const (
	// None is to be used for unset/default images
	None = iota
	// Alpine image
	Alpine
	// BusyBox image
	BusyBox
	// Pause image
	Pause
	// VolumeCopyUp image
	VolumeCopyUp
	// VolumeOwnership image
	VolumeOwnership
)

func initImageMap(imageList ImageList) map[int]string {
	images := map[int]string{}
	images[Alpine] = imageList.Alpine
	images[BusyBox] = imageList.BusyBox
	images[Pause] = imageList.Pause
	images[VolumeCopyUp] = imageList.VolumeCopyUp
	images[VolumeOwnership] = imageList.VolumeOwnership
	return images
}

// GetImage returns the fully qualified URI to an image (including version)
func GetImage(image int) string {
	return imageMap[image]
}
