// Original content of update_operation.go

package main

import (
    "context"
    "fmt"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/apimachinery/pkg/api/errors"
    "github.com/kyma-project/kyma-environment-broker/api/v2beta1"
    "github.com/kyma-project/kyma-environment-broker/common/broker"
)

// UpdateInstance updates the instance with the given parameters.
func (s *Service) UpdateInstance(ctx context.Context, instance *v2beta1.Instance, params *broker.UpdateInstanceParams) (*v2beta1.Instance, error) {
    // Check if the input parameters are the same as the current instance parameters
    if reflect.DeepEqual(instance.Spec, params.Spec) && reflect.DeepEqual(instance.Status, params.Status) {
        // If nothing has changed, return the current instance with a status code indicating no update was performed
        return instance, nil
    }

    // Perform the update operation
    updatedInstance, err := s.clientset.KymaV2beta1().Instances().Update(ctx, instance, metav1.UpdateOptions{})
    if err != nil {
        if errors.IsNotFound(err) {
            return nil, fmt.Errorf("instance not found: %v", err)
        }
        return nil, err
    }

    // If the update was successful and nothing changed, return an error indicating no update was performed
    if reflect.DeepEqual(updatedInstance.Spec, instance.Spec) && reflect.DeepEqual(updatedInstance.Status, instance.Status) {
        return nil, fmt.Errorf("no changes detected after update, returning error to indicate no update was performed")
    }

    // Return the updated instance
    return updatedInstance, nil
}
