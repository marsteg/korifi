package repositories

import (
	"context"
	"fmt"
	"slices"
	"time"

	"code.cloudfoundry.org/korifi/api/authorization"
	apierrors "code.cloudfoundry.org/korifi/api/errors"
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/tools"

	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	BuildResourceType = "Build"
)

type BuildRecord struct {
	GUID            string
	SpaceGUID       string
	State           string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	StagingErrorMsg string
	StagingMemoryMB int
	StagingDiskMB   int
	Lifecycle       Lifecycle
	PackageGUID     string
	DropletGUID     string
	AppGUID         string
	Labels          map[string]string
	Annotations     map[string]string
	ImageRef        string
}

func (r BuildRecord) Relationships() map[string]string {
	return map[string]string{
		"app": r.AppGUID,
	}
}

type BuildRepo struct {
	klient Klient
}

func NewBuildRepo(
	klient Klient,
) *BuildRepo {
	return &BuildRepo{
		klient: klient,
	}
}

func (b *BuildRepo) GetBuild(ctx context.Context, authInfo authorization.Info, buildGUID string) (BuildRecord, error) {
	build := korifiv1alpha1.CFBuild{
		ObjectMeta: metav1.ObjectMeta{
			Name: buildGUID,
		},
	}
	if err := b.klient.Get(ctx, &build); err != nil {
		return BuildRecord{}, fmt.Errorf("failed to get build: %w", apierrors.FromK8sError(err, BuildResourceType))
	}

	return b.cfBuildToBuildRecord(build), nil
}

func (b *BuildRepo) GetLatestBuildByAppGUID(ctx context.Context, authInfo authorization.Info, spaceGUID string, appGUID string) (BuildRecord, error) {
	buildList := &korifiv1alpha1.CFBuildList{}
	_, err := b.klient.List(ctx, buildList,
		InNamespace(spaceGUID),
		WithLabel(korifiv1alpha1.CFAppGUIDLabelKey, appGUID),
		WithOrdering("created_at"),
	)
	if err != nil {
		return BuildRecord{}, apierrors.FromK8sError(err, BuildResourceType)
	}

	if len(buildList.Items) == 0 {
		return BuildRecord{}, apierrors.NewNotFoundError(fmt.Errorf("builds for app %q in space %q not found", appGUID, spaceGUID), BuildResourceType)
	}

	return b.cfBuildToBuildRecord(buildList.Items[0]), nil
}

func (b *BuildRepo) cfBuildToBuildRecord(cfBuild korifiv1alpha1.CFBuild) BuildRecord {
	toReturn := BuildRecord{
		GUID:            cfBuild.Name,
		SpaceGUID:       cfBuild.Namespace,
		State:           tools.IfZero(cfBuild.Status.State, korifiv1alpha1.BuildStateStaging),
		CreatedAt:       cfBuild.CreationTimestamp.Time,
		UpdatedAt:       getLastUpdatedTime(&cfBuild),
		StagingErrorMsg: "",
		StagingMemoryMB: cfBuild.Spec.StagingMemoryMB,
		StagingDiskMB:   cfBuild.Spec.StagingDiskMB,
		Lifecycle: Lifecycle{
			Type: string(cfBuild.Spec.Lifecycle.Type),
			Data: LifecycleData{
				Buildpacks: []string{},
				Stack:      cfBuild.Spec.Lifecycle.Data.Stack,
			},
		},
		PackageGUID: cfBuild.Spec.PackageRef.Name,
		DropletGUID: "",
		AppGUID:     cfBuild.Spec.AppRef.Name,
		Labels:      cfBuild.Labels,
		Annotations: cfBuild.Annotations,
	}

	if cfBuild.Spec.Lifecycle.Type == "docker" {
		toReturn.Lifecycle.Data = LifecycleData{}
	}

	if cfBuild.Spec.Lifecycle.Data.Buildpacks != nil {
		toReturn.Lifecycle.Data.Buildpacks = cfBuild.Spec.Lifecycle.Data.Buildpacks
	}

	switch cfBuild.Status.State {
	case korifiv1alpha1.BuildStateStaged:
		toReturn.DropletGUID = cfBuild.Name
	case korifiv1alpha1.BuildStateFailed:
		toReturn.StagingErrorMsg = "Unknown error"
		conditionStatus := meta.FindStatusCondition(cfBuild.Status.Conditions, korifiv1alpha1.SucceededConditionType)
		if conditionStatus != nil {
			toReturn.StagingErrorMsg = conditionStatus.Message
		}
	}

	return toReturn
}

func (b *BuildRepo) CreateBuild(ctx context.Context, authInfo authorization.Info, message CreateBuildMessage) (BuildRecord, error) {
	cfBuild := message.toCFBuild()
	if err := b.klient.Create(ctx, &cfBuild); err != nil {
		return BuildRecord{}, apierrors.FromK8sError(err, BuildResourceType)
	}

	return b.cfBuildToBuildRecord(cfBuild), nil
}

func (b *BuildRepo) ListBuilds(ctx context.Context, authInfo authorization.Info, message ListBuildsMessage) (ListResult[BuildRecord], error) {
	buildList := &korifiv1alpha1.CFBuildList{}
	pageInfo, err := b.klient.List(ctx, buildList, message.toListOptions()...)
	if err != nil {
		return ListResult[BuildRecord]{}, fmt.Errorf("failed to list builds: %w", apierrors.FromK8sError(err, BuildResourceType))
	}

	return ListResult[BuildRecord]{
		Records:  slices.Collect(it.Map(slices.Values(buildList.Items), b.cfBuildToBuildRecord)),
		PageInfo: pageInfo,
	}, nil
}

type CreateBuildMessage struct {
	AppGUID         string
	PackageGUID     string
	SpaceGUID       string
	StagingMemoryMB int
	StagingDiskMB   int
	Lifecycle       Lifecycle
	Labels          map[string]string
	Annotations     map[string]string
}

type ListBuildsMessage struct {
	PackageGUIDs []string
	AppGUIDs     []string
	States       []string
	OrderBy      string
	Pagination   Pagination
}

func (m *ListBuildsMessage) toListOptions() []ListOption {
	return []ListOption{
		WithLabelIn(korifiv1alpha1.CFPackageGUIDLabelKey, m.PackageGUIDs),
		WithLabelIn(korifiv1alpha1.CFAppGUIDLabelKey, m.AppGUIDs),
		WithLabelIn(korifiv1alpha1.CFBuildStateLabelKey, m.States),
		WithPaging(m.Pagination),
		WithOrdering(m.OrderBy),
	}
}

func (m CreateBuildMessage) toCFBuild() korifiv1alpha1.CFBuild {
	return korifiv1alpha1.CFBuild{
		ObjectMeta: metav1.ObjectMeta{
			Name:        uuid.NewString(),
			Namespace:   m.SpaceGUID,
			Labels:      m.Labels,
			Annotations: m.Annotations,
		},
		Spec: korifiv1alpha1.CFBuildSpec{
			PackageRef: corev1.LocalObjectReference{
				Name: m.PackageGUID,
			},
			AppRef: corev1.LocalObjectReference{
				Name: m.AppGUID,
			},
			StagingMemoryMB: m.StagingMemoryMB,
			StagingDiskMB:   m.StagingDiskMB,
			Lifecycle: korifiv1alpha1.Lifecycle{
				Type: korifiv1alpha1.LifecycleType(m.Lifecycle.Type),
				Data: korifiv1alpha1.LifecycleData{
					Buildpacks: m.Lifecycle.Data.Buildpacks,
					Stack:      m.Lifecycle.Data.Stack,
				},
			},
		},
	}
}
