** Ручной запуск контейнера prometheus **
```docker run -d \
  --name prometheus \
  --restart always \
  --network host \
  --user "0:0" \
  -v /etc/prometheus:/etc/prometheus \
  -v prometheus_data:/prometheus \
  prom/prometheus \
  --config.file=/etc/prometheus/prometheus.yml \
  --storage.tsdb.path=/prometheus```