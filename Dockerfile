FROM golang:1.9
WORKDIR /go/src/
RUN go get github.com/eclipse/paho.mqtt.golang
COPY  ./edgexwrapper /go/src/edgexwrapper
RUN cd edgexwrapper/ && go build device.go && cd ..

FROM video_analytics_serving_gstreamer
WORKDIR /home/video-analytics
RUN mkdir edgexwrapper
COPY --from=0 /go/src/edgexwrapper/device /home/video-analytics/edgexwrapper/
COPY ./pipelines /home/video-analytics/pipelines
COPY   docker-entrypoint.sh /home/video-analytics/
ENTRYPOINT ["./docker-entrypoint.sh"]
