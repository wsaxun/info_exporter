# build
docker build -t info_exporter:v1.0.0 .

# run
docker run -itd --name info_exporter -p 8080:8080 -e MFYENV=dev  -e MFYKEY=mykey -e MFYCYCLE=120  info_exporter:v1.0.0
