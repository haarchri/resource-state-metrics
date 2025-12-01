package external

import (
	"io"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	metricsstore "k8s.io/kube-state-metrics/v2/pkg/metrics_store"
)

// collectors defines behaviors to implement custom Go-based collectors for the "main" instance.
type gvkr struct {
	schema.GroupVersionKind
	schema.GroupVersionResource
}
type collectors interface {
	BuildCollector(kubeconfig string) *metricsstore.MetricsStore
	GVKR() gvkr
	Register()
	Name() string
}

type collectorsType struct {
	kubeconfig      string
	collectors      []collectors
	builtCollectors []*metricsstore.MetricsStore
	enabled         map[string]bool
}

func (ct *collectorsType) SetKubeConfig(kubeconfig string) *collectorsType {
	ct.kubeconfig = kubeconfig

	return ct
}

func (ct *collectorsType) SetEnabled(enabledList string) *collectorsType {
	if enabledList == "" {
		// Empty string means no collectors enabled, set nil map
		ct.enabled = nil

		return ct
	}
	ct.enabled = make(map[string]bool)
	for _, name := range strings.Split(enabledList, ",") {
		name = strings.TrimSpace(name)
		if name != "" {
			ct.enabled[name] = true
		}
	}

	return ct
}

func (ct *collectorsType) Register(c collectors) {
	ct.collectors = append(ct.collectors, c)
}

func (ct *collectorsType) Build() {
	for _, c := range ct.collectors {
		// If enabled is nil, it means SetEnabled was called with empty string = disable all
		// If enabled is not nil, only build collectors that are in the enabled map
		if ct.enabled != nil && ct.enabled[c.Name()] {
			ct.builtCollectors = append(ct.builtCollectors, c.BuildCollector(ct.kubeconfig))
		}
	}
}

func (ct *collectorsType) Write(w io.Writer) {
	for _, c := range ct.builtCollectors {
		mw := metricsstore.NewMetricsWriter(c)
		_ = mw.WriteAll(w)
	}
}

var availableCollectors = []collectors{
	&clusterResourceQuotaCollector{},
}

var collectorsInstance = &collectorsType{
	collectors: availableCollectors,
}

//nolint:revive
func CollectorsGetter() *collectorsType {
	return collectorsInstance
}
