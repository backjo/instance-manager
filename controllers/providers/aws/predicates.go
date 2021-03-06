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

package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/eks"
)

func IsUsingLaunchConfiguration(group *autoscaling.Group) bool {
	if group.LaunchConfigurationName != nil {
		return true
	}
	return false
}

func IsUsingLaunchTemplate(group *autoscaling.Group) bool {
	if group.LaunchTemplate != nil {
		if group.LaunchTemplate.LaunchTemplateName != nil {
			return true
		}
	}
	return false
}

func IsNodeGroupInConditionState(key string, condition string) bool {
	conditionStates := map[string]ManagedNodeGroupReconcileState{
		"CREATING":      ManagedNodeGroupOngoingState,
		"UPDATING":      ManagedNodeGroupOngoingState,
		"DELETING":      ManagedNodeGroupOngoingState,
		"ACTIVE":        ManagedNodeGroupFiniteState,
		"DEGRADED":      ManagedNodeGroupFiniteState,
		"CREATE_FAILED": ManagedNodeGroupUnrecoverableError,
		"DELETE_FAILED": ManagedNodeGroupUnrecoverableDeleteError,
	}
	state := conditionStates[key]

	switch condition {
	case "OngoingState":
		return state.OngoingState
	case "FiniteState":
		return state.FiniteState
	case "UnrecoverableError":
		return state.UnrecoverableError
	case "UnrecoverableDeleteError":
		return state.UnrecoverableDeleteError
	default:
		return false
	}
}

func IsProfileInConditionState(key string, condition string) bool {

	conditionStates := map[string]CloudResourceReconcileState{
		aws.StringValue(nil):                 FiniteDeleted,
		eks.FargateProfileStatusCreating:     OngoingState,
		eks.FargateProfileStatusActive:       FiniteState,
		eks.FargateProfileStatusDeleting:     OngoingState,
		eks.FargateProfileStatusCreateFailed: UpdateRecoverableError,
		eks.FargateProfileStatusDeleteFailed: UnrecoverableDeleteError,
	}
	state := conditionStates[key]
	switch condition {
	case "OngoingState":
		return state.OngoingState
	case "FiniteState":
		return state.FiniteState
	case "FiniteDeleted":
		return state.FiniteDeleted
	case "UpdateRecoverableError":
		return state.UpdateRecoverableError
	case "UnrecoverableError":
		return state.UnrecoverableError
	case "UnrecoverableDeleteError":
		return state.UnrecoverableDeleteError
	default:
		return false
	}
}

func IsUsingMixedInstances(group *autoscaling.Group) bool {
	if group.MixedInstancesPolicy != nil {
		return true
	}
	return false
}
