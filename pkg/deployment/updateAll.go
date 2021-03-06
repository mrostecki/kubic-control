// Copyright 2019, 2020 Thorsten Kukuk
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deployment

import (
	"gopkg.in/ini.v1"
	log "github.com/sirupsen/logrus"
	"github.com/thkukuk/kubic-control/pkg/tools"
)

func UpdateAll(forced bool) (bool, string) {

	cfg, err := ini.Load("/var/lib/kubic-control/k8s-yaml.conf")
	if err != nil {
		return false, "Cannot load k8s-yaml.conf: " + err.Error()
        }

	keys := cfg.Section("").KeyStrings()
	for _, key := range keys {
		if forced {
			// force, so always update even if not changed
			success, message := UpdateFile(key)
			if success != true {
				return success, message
			}
		} else {
			value := cfg.Section("").Key(key).String()
			hash, _ := tools.Sha256sum_f(key)

			if hash != value {
				log.Infof("%s has changed, updating")
				success, message := UpdateFile(key)
				if success != true {
					return success, message
				}
			} else {
				log.Infof("%s has not changed, ignoring")
			}
		}
	}


	// Update kustomize installed services
	cfg, err = ini.Load("/var/lib/kubic-control/k8s-kustomize.conf")
	if err != nil {
		return false, "Cannot load k8s-kustomize.conf: " + err.Error()
        }

	keys = cfg.Section("").KeyStrings()
	for _, key := range keys {
		if forced {
			// force, so always update even if not changed
			success, message := UpdateKustomize(key)
			if success != true {
				return success, message
			}
		} else {
			retval, message := tools.ExecuteCmd("kustomize", "build",
				StateDir + "/kustomize/" + key + "/overlay")
			if retval != true {
				return retval, message
			}

			value := cfg.Section("").Key(key).String()
			hash, _ := tools.Sha256sum_f(message)

			if hash != value {
				log.Infof("%s has changed, updating")
				success, message := UpdateKustomize(key)
				if success != true {
					return success, message
				}
			} else {
				log.Infof("%s has not changed, ignoring")
			}
		}
	}

 	return true, ""
}
