package storagemanager

import (
	"github.com/libopenstorage/cloudops"
	"github.com/libopenstorage/cloudops/pkg/storagedistribution"
	"github.com/libopenstorage/cloudops/unsupported"
)

type csiStorageManager struct {
	cloudops.StorageManager
	decisionMatrix *cloudops.StorageDecisionMatrix
}

// NewCSIStorageManager returns an CSI implementation for Storage Management
func NewCSIStorageManager(
	decisionMatrix cloudops.StorageDecisionMatrix,
) (cloudops.StorageManager, error) {
	return newCSIStorageManager(decisionMatrix)
}

// newCSIStorageManager returns an CSI implementation for Storage Management
func newCSIStorageManager(
	decisionMatrix cloudops.StorageDecisionMatrix,
) (cloudops.StorageManager, error) {
	return &csiStorageManager{
		StorageManager: unsupported.NewUnsupportedStorageManager(),
		decisionMatrix: &decisionMatrix}, nil
}

func (a *csiStorageManager) GetStorageDistribution(
	request *cloudops.StorageDistributionRequest,
) (*cloudops.StorageDistributionResponse, error) {
	response := &cloudops.StorageDistributionResponse{}
	for _, userRequest := range request.UserStorageSpec {
		// for for request, find how many instances per zone needs to have storage
		// and the storage spec for each of them
		instStorage, instancesPerZone, _, err :=
			storagedistribution.GetStorageDistributionForPool(
				a.decisionMatrix,
				userRequest,
				request.InstancesPerZone,
				request.ZoneCount,
			)
		if err != nil {
			return nil, err
		}
		response.InstanceStorage = append(
			response.InstanceStorage,
			&cloudops.StoragePoolSpec{
				DriveCapacityGiB: instStorage.DriveCapacityGiB,
				DriveType:        instStorage.DriveType,
				InstancesPerZone: instancesPerZone,
				DriveCount:       instStorage.DriveCount,
			},
		)
	}
	return response, nil
}

func (a *csiStorageManager) RecommendStoragePoolUpdate(
	request *cloudops.StoragePoolUpdateRequest) (*cloudops.StoragePoolUpdateResponse, error) {
	resp, _, err := storagedistribution.GetStorageUpdateConfig(request, a.decisionMatrix)
	return resp, err
}

func init() {
	cloudops.RegisterStorageManager(cloudops.CSI, newCSIStorageManager)
}
