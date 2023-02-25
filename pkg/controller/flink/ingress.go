package flink

import (
	"fmt"
	"regexp"

	flinkapp "github.com/lyft/flinkk8soperator/pkg/apis/app/v1beta1"
	"github.com/lyft/flinkk8soperator/pkg/controller/common"
	"github.com/lyft/flinkk8soperator/pkg/controller/config"
	"github.com/lyft/flinkk8soperator/pkg/controller/k8"
	networkingV1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const AppIngressName = "%s-%s"

var inputRegex = regexp.MustCompile(`{{[$]jobCluster}}`)

func ReplaceJobURL(value string, input string) string {
	return inputRegex.ReplaceAllString(value, input)
}

func GetFlinkUIIngressURL(jobName string) string {
	return ReplaceJobURL(config.GetConfig().FlinkIngressURLFormat, jobName)
}

func FetchJobManagerIngressCreateObj(app *flinkapp.FlinkApplication) *networkingV1.Ingress {
	podLabels := common.DuplicateMap(app.Labels)
	podLabels = common.CopyMap(podLabels, k8.GetAppLabel(app.Name))

	ingressMeta := v1.ObjectMeta{
		Name:      getJobManagerServiceName(app),
		Labels:    podLabels,
		Namespace: app.Namespace,
		OwnerReferences: []v1.OwnerReference{
			*v1.NewControllerRef(app, app.GroupVersionKind()),
		},
	}

	ingressServiceBackend := networkingV1.IngressServiceBackend{
		Name: getJobManagerServiceName(app),
		Port: networkingV1.ServiceBackendPort{
			Number: getUIPort(app),
		},
	}

	backend := networkingV1.IngressBackend{
		Service: &ingressServiceBackend,
	}
	pathType := networkingV1.PathType("ImplementationSpecific")
	ingressSpec := networkingV1.IngressSpec{
		Rules: []networkingV1.IngressRule{{
			Host: GetFlinkUIIngressURL(getIngressName(app)),
			IngressRuleValue: networkingV1.IngressRuleValue{
				HTTP: &networkingV1.HTTPIngressRuleValue{
					Paths: []networkingV1.HTTPIngressPath{{
						Backend:  backend,
						PathType: &pathType,
					}},
				},
			},
		}},
	}
	return &networkingV1.Ingress{
		ObjectMeta: ingressMeta,
		TypeMeta: v1.TypeMeta{
			APIVersion: networkingV1.SchemeGroupVersion.String(),
			Kind:       k8.Ingress,
		},
		Spec: ingressSpec,
	}

}

func getIngressName(app *flinkapp.FlinkApplication) string {
	if flinkapp.IsBlueGreenDeploymentMode(app.Spec.DeploymentMode) {
		return fmt.Sprintf(AppIngressName, app.Name, string(app.Status.UpdatingVersion))
	}
	return app.Name
}
