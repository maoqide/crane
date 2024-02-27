package hpa

import (
	"fmt"

	"github.com/gocrane/crane/pkg/metricnaming"
	"github.com/gocrane/crane/pkg/recommendation/framework"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
)

// CheckDataProviders in PrePrepare phase, will create data source provider via your recommendation config.
func (rr *HPARecommender) CheckDataProviders(ctx *framework.RecommendationContext) error {
	return rr.ReplicasRecommender.CheckDataProviders(ctx)
}

func (rr *HPARecommender) CollectData(ctx *framework.RecommendationContext) error {
	err := rr.ReplicasRecommender.CollectData(ctx)

	resourceCpu := corev1.ResourceCPU
	labelSelector := labels.SelectorFromSet(ctx.Identity.Labels)
	caller := fmt.Sprintf(callerFormat, klog.KObj(ctx.Recommendation), ctx.Recommendation.UID)
	metricNamer := metricnaming.ResourceToWorkloadMetricNamer(ctx.Recommendation.Spec.TargetRef.DeepCopy(), &resourceCpu, labelSelector, caller)

	if err := metricNamer.Validate(); err != nil {
		return err
	}
	ctx.MetricNamer = metricNamer

	return err
}

func (rr *HPARecommender) PostProcessing(ctx *framework.RecommendationContext) error {
	return rr.ReplicasRecommender.PostProcessing(ctx)
}
