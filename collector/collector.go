package collector

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sync"
)

type Metrics struct {
	metrics   map[string]*prometheus.Desc
	mutex     sync.Mutex
	clientset *kubernetes.Clientset
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
			"container_monitor_check_duration_millisecond": newGlobalMetric("container_monitor_check_duration_millisecond", "The time(millisecond) taken to invoke the health check interface", []string{"namespace", "container_name", "pod_name"}),
		},
		clientset: clientset,
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

	pods, err := c.clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, item := range pods.Items {
		meta := item.ObjectMeta
		podName := meta.Name
		labels := meta.Labels
		containerName := labels["app"]

		ch <- prometheus.MustNewConstMetric(c.metrics["container_monitor_check_duration_millisecond"], prometheus.GaugeValue, 0, "kube-system", containerName, podName)
	}

}
