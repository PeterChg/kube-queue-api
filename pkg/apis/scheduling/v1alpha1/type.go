/*

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Queue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,name=metadata"`

	Spec   QueueSpec   `json:"spec,omitempty" protobuf:"bytes,2,name=spec"`
	Status QueueStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// QueueSpec defines the desired state of Queue
type QueueSpec struct {
	QueuePolicy       QueuePolicy `json:"queuePolicy,omitempty" protobuf:"bytes,1,opt,name=queuePolicy`
	Priority          *int32      `json:"priority,omitempty" protobuf:"varint,2,opt,name=priority"`
	PriorityClassName string      `json:"priorityClassName,omitempty" protobuf:"bytes,3,opt,name=priorityClassName"`
}

// QueueStatus defines the observed state of Queue
type QueueStatus struct {
	// TODO
}

// +k8s:openapi-gen=true
// QueuePolicy defines the queueing policy for the elements in the queue
type QueuePolicy string

const (
	QueuePolicyFIFO     QueuePolicy = "FIFO"
	QueuePolicyPriority QueuePolicy = "Priority"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QueueUnitList contains a list of QueueUnit
type QueueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Queue `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type QueueUnit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,name=metadata"`

	Spec   QueueUnitSpec   `json:"spec,omitempty" protobuf:"bytes,2,name=spec"`
	Status QueueUnitStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// QueueUnitSpec defines the desired state of QueueUnit
type QueueUnitSpec struct {
	ConsumerRef       *corev1.ObjectReference `json:"consumerRef,omitempty" protobuf:"bytes,1,opt,name=consumerRef"`
	Priority          *int32                  `json:"priority,omitempty" protobuf:"varint,2,opt,name=priority"`
	Queue             string                  `json:"queue,omitempty" protobuf:"bytes,3,opt,name=queue"`
	Resource          corev1.ResourceList     `json:"resource,omitempty" protobuf:"bytes,4,name=resource"`
	PriorityClassName string                  `json:"priorityClassName,omitempty" protobuf:"bytes,5,opt,name=priorityClassName"`
}

// QueueUnitStatus defines the observed state of QueueUnit
type QueueUnitStatus struct {
	Phase          QueueUnitPhase `json:"phase" protobuf:"bytes,1,name=phase"`
	Message        string         `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
	LastUpdateTime *metav1.Time   `json:"lastUpdateTime" protobuf:"bytes,3,name=lastUpdateTime"`
	Position       string         `json:"position" protobuf:"bytes,4,name=position"`
}

type QueueUnitPhase string

const (
	Enqueued     QueueUnitPhase = "Enqueued"
	Dequeued     QueueUnitPhase = "Dequeued"
	SchedReady   QueueUnitPhase = "Running"
	SchedSucceed QueueUnitPhase = "Succeed"
	SchedFailed  QueueUnitPhase = "Failed"
	Backoff      QueueUnitPhase = "TimeoutBackoff"
	JobNotFound  QueueUnitPhase = "related Job not found"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QueueUnitList contains a list of QueueUnit
type QueueUnitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QueueUnit `json:"items"`
}

// Suspend is a flag that instructs the job operator to suspend processing this job
const Suspend = "scheduling.x-k8s.io/suspend"

// Placement is the scheduling result of the scheduler
const Placement = "scheduling.x-k8s.io/placement"

// JobSuspended checks if a Job should be suspended by checking whether its annotation contains key Suspend and its
// value is set "true"
func JobSuspended(job metav1.Object) bool {
	const suspended = "true"
	annotations := job.GetAnnotations()
	if annotations != nil {
		if val, exist := annotations[Suspend]; exist {
			return val == suspended
		}
	}
	return false
}
