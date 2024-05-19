// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
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

package creator

import (
	"github.com/altinity/clickhouse-operator/pkg/model/chi/tags/annotator"
	"github.com/altinity/clickhouse-operator/pkg/model/chi/tags/labeler"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/altinity/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	"github.com/altinity/clickhouse-operator/pkg/model/chi/namer"
)

// AdjustPVC prepares PVC - labels and annotations
func (c *Creator) AdjustPVC(
	pvc *core.PersistentVolumeClaim,
	host *api.Host,
	template *api.VolumeClaimTemplate,
) *core.PersistentVolumeClaim {
	pvc.SetLabels(namer.Macro(host).Map(c.tagger.Label(labeler.LabelExistingPVC, pvc, host, template)))
	pvc.SetAnnotations(namer.Macro(host).Map(c.tagger.Annotate(annotator.AnnotateExistingPVC, pvc, host, template)))
	// And after the object is ready we can put version label
	labeler.MakeObjectVersion(&pvc.ObjectMeta, pvc)
	return pvc
}

// createPVC
func (c *Creator) CreatePVC(
	name string,
	namespace string,
	host *api.Host,
	spec *core.PersistentVolumeClaimSpec,
) *core.PersistentVolumeClaim {
	persistentVolumeClaim := core.PersistentVolumeClaim{
		TypeMeta: meta.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			// TODO
			//  this has to wait until proper disk inheritance procedure will be available
			// UPDATE
			//  we are close to proper disk inheritance
			// Right now we hit the following error:
			// "Forbidden: updates to statefulset spec for fields other than 'replicas', 'template', and 'updateStrategy' are forbidden"
			Labels:      namer.Macro(host).Map(c.tagger.Label(labeler.LabelNewPVC, host, false)),
			Annotations: namer.Macro(host).Map(c.tagger.Annotate(annotator.AnnotateNewPVC, host)),
		},
		// Append copy of PersistentVolumeClaimSpec
		Spec: *spec.DeepCopy(),
	}
	// TODO introduce normalization
	// Overwrite .Spec.VolumeMode
	volumeMode := core.PersistentVolumeFilesystem
	persistentVolumeClaim.Spec.VolumeMode = &volumeMode

	return &persistentVolumeClaim
}
