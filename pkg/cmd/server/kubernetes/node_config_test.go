package kubernetes

import (
	"reflect"
	"testing"
	"time"

	proxyoptions "k8s.io/kubernetes/cmd/kube-proxy/app/options"
	kubeletoptions "k8s.io/kubernetes/cmd/kubelet/app/options"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/componentconfig"
	"k8s.io/kubernetes/pkg/kubelet/rkt"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	"k8s.io/kubernetes/pkg/util"
	utilconfig "k8s.io/kubernetes/pkg/util/config"
	"k8s.io/kubernetes/pkg/util/diff"
)

func TestKubeletDefaults(t *testing.T) {
	defaults := kubeletoptions.NewKubeletServer()

	// This is a snapshot of the default config
	// If the default changes (new fields are added, or default values change), we want to know
	// Once we've reacted to the changes appropriately in BuildKubernetesNodeConfig(), update this expected default to match the new upstream defaults
	expectedDefaults := &kubeletoptions.KubeletServer{
		AuthPath:   util.NewStringFlag("/var/lib/kubelet/kubernetes_auth"),
		KubeConfig: util.NewStringFlag("/var/lib/kubelet/kubeconfig"),

		SystemReserved: utilconfig.ConfigurationMap{},
		KubeReserved:   utilconfig.ConfigurationMap{},
		KubeletConfiguration: componentconfig.KubeletConfiguration{
			Address:                     "0.0.0.0", // overridden
			AllowPrivileged:             false,     // overridden
			CAdvisorPort:                4194,      // disabled
			VolumeStatsAggPeriod:        unversioned.Duration{Duration: time.Minute},
			CertDirectory:               "/var/run/kubernetes",
			CgroupRoot:                  "",
			ClusterDNS:                  "", // overridden
			ClusterDomain:               "", // overridden
			ConfigureCBR0:               false,
			ContainerRuntime:            "docker",
			Containerized:               false, // overridden based on OPENSHIFT_CONTAINERIZED
			CPUCFSQuota:                 true,  // forced to true
			DockerExecHandlerName:       "native",
			EventBurst:                  10,
			EventRecordQPS:              5.0,
			EnableCustomMetrics:         false,
			EnableDebuggingHandlers:     true,
			EnableServer:                true,
			FileCheckFrequency:          unversioned.Duration{Duration: 20 * time.Second}, // overridden
			HealthzBindAddress:          "127.0.0.1",                                      // disabled
			HealthzPort:                 10248,                                            // disabled
			HostNetworkSources:          "*",                                              // overridden
			HostPIDSources:              "*",                                              // overridden
			HostIPCSources:              "*",                                              // overridden
			HTTPCheckFrequency:          unversioned.Duration{Duration: 20 * time.Second}, // disabled
			ImageMinimumGCAge:           unversioned.Duration{Duration: 2 * time.Minute},
			ImageGCHighThresholdPercent: 90,
			ImageGCLowThresholdPercent:  80,
			LowDiskSpaceThresholdMB:     256,
			MasterServiceNamespace:      "default",
			MaxContainerCount:           240,
			MaxPerPodContainerCount:     2,
			MaxOpenFiles:                1000000,
			MaxPods:                     110, // overridden
			MinimumGCAge:                unversioned.Duration{Duration: 1 * time.Minute},
			NetworkPluginDir:            "/usr/libexec/kubernetes/kubelet-plugins/net/exec/",
			NetworkPluginName:           "", // overridden
			NonMasqueradeCIDR:           "10.0.0.0/8",
			VolumePluginDir:             "/usr/libexec/kubernetes/kubelet-plugins/volume/exec/",
			NodeStatusUpdateFrequency:   unversioned.Duration{Duration: 10 * time.Second},
			NodeLabels:                  map[string]string{},
			OOMScoreAdj:                 -999,
			LockFilePath:                "",
			PodInfraContainerImage:      kubeletoptions.GetDefaultPodInfraContainerImage(), // overridden
			Port:                           10250, // overridden
			ReadOnlyPort:                   10255, // disabled
			RegisterNode:                   true,
			RegisterSchedulable:            true,
			RegistryBurst:                  10,
			RegistryPullQPS:                5.0,
			ResolverConfig:                 kubetypes.ResolvConfDefault,
			KubeletCgroups:                 "",
			RktAPIEndpoint:                 rkt.DefaultRktAPIServiceEndpoint,
			RktPath:                        "",
			RktStage1Image:                 "",
			RootDirectory:                  "/var/lib/kubelet", // overridden
			RuntimeCgroups:                 "",
			SerializeImagePulls:            true,
			StreamingConnectionIdleTimeout: unversioned.Duration{Duration: 4 * time.Hour},
			SyncFrequency:                  unversioned.Duration{Duration: 1 * time.Minute},
			SystemCgroups:                  "",
			TLSCertFile:                    "", // overridden to prevent cert generation
			TLSPrivateKeyFile:              "", // overridden to prevent cert generation
			ReconcileCIDR:                  true,
			KubeAPIQPS:                     5.0,
			KubeAPIBurst:                   10,
			ExperimentalFlannelOverlay:     false,
			OutOfDiskTransitionFrequency:   unversioned.Duration{Duration: 5 * time.Minute},
			HairpinMode:                    "promiscuous-bridge",
			BabysitDaemons:                 false,
			SeccompProfileRoot:             "/var/lib/kubelet/seccomp",
			CloudProvider:                  "auto-detect",
			RuntimeRequestTimeout:          unversioned.Duration{Duration: 2 * time.Minute},
			ContentType:                    "application/vnd.kubernetes.protobuf",
			EnableControllerAttachDetach:   true,

			EvictionPressureTransitionPeriod: unversioned.Duration{Duration: 5 * time.Minute},
		},
	}

	if !reflect.DeepEqual(defaults, expectedDefaults) {
		t.Logf("expected defaults, actual defaults: \n%s", diff.ObjectReflectDiff(expectedDefaults, defaults))
		t.Errorf("Got different defaults than expected, adjust in BuildKubernetesNodeConfig and update expectedDefaults")
	}
}

func TestProxyConfig(t *testing.T) {
	// This is a snapshot of the default config
	// If the default changes (new fields are added, or default values change), we want to know
	// Once we've reacted to the changes appropriately in buildKubeProxyConfig(), update this expected default to match the new upstream defaults
	oomScoreAdj := int32(-999)
	ipTablesMasqueratebit := int32(14)
	expectedDefaultConfig := &proxyoptions.ProxyServerConfig{
		KubeProxyConfiguration: componentconfig.KubeProxyConfiguration{
			BindAddress:        "0.0.0.0",
			ClusterCIDR:        "",
			HealthzPort:        10249,         // disabled
			HealthzBindAddress: "127.0.0.1",   // disabled
			OOMScoreAdj:        &oomScoreAdj,  // disabled
			ResourceContainer:  "/kube-proxy", // disabled
			IPTablesSyncPeriod: unversioned.Duration{Duration: 30 * time.Second},
			// from k8s.io/kubernetes/cmd/kube-proxy/app/options/options.go
			// defaults to 14.
			IPTablesMasqueradeBit:          &ipTablesMasqueratebit,
			UDPIdleTimeout:                 unversioned.Duration{Duration: 250 * time.Millisecond},
			ConntrackMax:                   256 * 1024,                                          // 4x default (64k)
			ConntrackTCPEstablishedTimeout: unversioned.Duration{Duration: 86400 * time.Second}, // 1 day (1/5 default)
		},
		ConfigSyncPeriod: 15 * time.Minute,
		KubeAPIQPS:       5.0,
		KubeAPIBurst:     10,
		ContentType:      "application/vnd.kubernetes.protobuf",
	}

	actualDefaultConfig := proxyoptions.NewProxyConfig()

	if !reflect.DeepEqual(expectedDefaultConfig, actualDefaultConfig) {
		t.Errorf("Default kube proxy config has changed. Adjust buildKubeProxyConfig() as needed to disable or make use of additions.")
		t.Logf("Difference %s", diff.ObjectReflectDiff(expectedDefaultConfig, actualDefaultConfig))
	}

}
