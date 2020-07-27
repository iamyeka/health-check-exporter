package collector

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Metrics struct {
	metrics    map[string]*prometheus.Desc
	mutex      sync.Mutex
	clientset  *kubernetes.Clientset
	httpClient *http.Client
}

func newGlobalMetric(metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(metricName, docString, labels, nil)
}

func NewMetrics() *Metrics {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			"container_health_check_duration_millisecond": newGlobalMetric("container_health_check_duration_millisecond", "The time(millisecond) taken to invoke the health check interface", []string{"namespace", "container_name", "pod_name"}),
		},
		clientset:  clientset,
		httpClient: &http.Client{Timeout: 3 * time.Second},
	}
}

/**
 * 接口：Describe
 * 功能：传递结构体中的指标描述符到channel
 */
func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

/**
 * 接口：Collect
 * 功能：抓取最新的数据，传递给channel
 */
func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // 加锁
	defer c.mutex.Unlock()

	pods, err := c.clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	items := pods.Items
	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		tmp := item
		go healthCheck(&tmp, c, ch, &wg)
	}

	wg.Wait()
}

func healthCheck(pod *coreV1.Pod, c *Metrics, ch chan<- prometheus.Metric, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	meta := pod.ObjectMeta
	spec := pod.Spec
	status := pod.Status
	podName := meta.Name
	labels := meta.Labels
	containerName := labels["app"]

	livenessProbe := spec.Containers[0].LivenessProbe

	if livenessProbe != nil && livenessProbe.HTTPGet != nil {
		podIP := status.PodIP
		httpGet := livenessProbe.HTTPGet

		start := time.Now()

		var scheme string
		if coreV1.URISchemeHTTP == httpGet.Scheme {
			scheme = "http://"
		} else {
			scheme = "https://"
		}

		resp, err := c.httpClient.Get(scheme + podIP + ":" + strconv.Itoa(int(httpGet.Port.IntVal)) + httpGet.Path)

		var duration time.Duration
		if err != nil {
			duration = -1
		} else {
			duration = time.Since(start)
		}

		if resp != nil {
			defer resp.Body.Close()
		}

		ch <- prometheus.MustNewConstMetric(c.metrics["container_health_check_duration_millisecond"], prometheus.GaugeValue, float64(duration), meta.Namespace, containerName, podName)

	}

}
