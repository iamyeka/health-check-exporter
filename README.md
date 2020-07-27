# health-check-exporter
Prometheus自定义exporter，基于Go语言开发，可采集集群中所有Pod的健康检查耗时(单位：毫秒)

Core Feature:
- 对所有命名空间下的Pod发起livenessProbe指针调用，返回调用耗时，当无法访问时，返回-1

提供启动参数
- web.listen-port 暴露的Http端口，默认80
- web.telemetry-path 访问指标的路径，默认/metrics

目前只暴露一个metric：**container_health_check_duration_millisecond**
```
# HELP container_health_check_duration_millisecond The time(millisecond) taken to invoke the health check interface
# TYPE container_health_check_duration_millisecond gauge
container_health_check_duration_millisecond{container_name="",namespace="kube-system",pod_name="coredns-67fc48b9d7-dtb9w"} 9.297547e+06
```

通过kubectl工具快速部署体验
```shell script
kubectl create deploy health-check-exporter --image=registry.cn-hangzhou.aliyuncs.com/clo/health-check-exporter:1.0
```