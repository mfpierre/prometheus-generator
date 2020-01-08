# prometheus-generator

Generate predictively prometheus metrics

## usage

`./prometheus-generator --counters 10 --gauges 10`

## kubernetes sample deployment with Datadog autodiscovery

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus-generator
spec:
  selector:
    matchLabels:
      app: prometheus-generator
  replicas: 1
  template:
    metadata:
      labels:
        app: prometheus-generator
      name: prometheus-generator
      annotations:
        ad.datadoghq.com/prometheus-generator.check_names: '["prometheus"]'
        ad.datadoghq.com/prometheus-generator.init_configs: '[{}]'
        ad.datadoghq.com/prometheus-generator.instances: '[{"prometheus_url": "http://%%host%%:8080/metrics","namespace": "prom_gen","metrics": ["*"]}]'
    spec:
      containers:
      - image: mfpierre/prometheus-generator
        imagePullPolicy: Always
        name: prometheus-generator
        resources:
          limits:
            memory: 100Mi
        command: ["./prometheus-generator"]
        args: ["--counters","5000","--gauges","5000"]
        ports:
        - containerPort: 8080
```
