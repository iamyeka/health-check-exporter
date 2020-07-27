# health-check-exporter
Prometheus自定义exporter，基于Go语言开发，可搜集集群中所有Pod的健康检查耗时(单位毫秒)

目前只暴露一个Metric：**container_health_check_duration_millisecond**
```
# HELP container_health_check_duration_millisecond The time(millisecond) taken to invoke the health check interface
# TYPE container_health_check_duration_millisecond gauge
container_health_check_duration_millisecond{container_name="",namespace="kube-system",pod_name="coredns-67fc48b9d7-dtb9w"} 9.297547e+06
```

