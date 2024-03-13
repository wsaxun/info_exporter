# 功能描述
1. 自定义 exporter demo
2. 建立 gorutinue 定时采集信息
3. 自定义 collect 暴露指标

# build
docker build -t info_exporter:v1.0.0 .

# run
docker run -itd --name info_exporter -p 8080:8080 -e MFYENV=dev  -e MFYKEY=mykey -e MFYCYCLE=120  info_exporter:v1.0.0
