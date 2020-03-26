# Video Analytics Serving with EdgeX Integration

[Video Analytics Serving](https://github.com/intel/video-analytics-serving) (VAS) is designed to simplify the deployment and use of hardware optimized video analytics pipelines. It offers developers a simple way to create REStful APIs to start, stop, enumerate and customize pre-defined pipelines using either GStreamer or FFmpeg. Developers create pipeline templates using their framework of choice and Video Analytics Serving manages launching pipeline instances based on incoming requests.

This repository is a wrapper around [video analytics serving](https://github.com/intel/video-analytics-serving) that integrate with EdgeX for sending the inferenced data to EdgeX Device MQTT Service.


## Building and Running

Video Analytics Serving may be modified to co-exist in a container alongside other applications or can be built and run as a standalone service.
### Prerequisites
(1) Install [docker engine](https://docs.docker.com/install).  
(2) Install [docker compose](https://docs.docker.com/compose/install), if you plan to deploy through docker compose. Version 1.20+ is required.


### Building

To get started, build the service as a standalone component execute the following command

```bash
$ ./build.sh
```

### Running

After a successful build, run the service using the included script

```bash
$ ./run_basic.sh
```

This script issues a standard docker run command to launch the container. Volume mounting is used to publish the sample results to your host.


## Example Pipelines

Video Analytics Serving includes two [sample](pipelines) analytics pipelines for GStreamer. The GStreamer sample pipelines use [plugins](https://github.com/opencv/gst-video-analytics) for CNN model-based video analytics utilizing [Intel OpenVINO](https://software.intel.com/en-us/openvino-toolkit). When building with Docker, these plugins will be built and installed inside the Docker image. You can find documentation on the properties of these elements [here](https://github.com/opencv/gst-video-analytics/wiki/Elements).

### Without EdgeX Integration

|Pipeline| Description| Example Request| Example Detection|
|--------|------------|----------------|------------------|
|/pipelines/object_detection/1|Object Detection|curl localhost:8080/pipelines/object_detection/1 -X POST -H 'Content-Type: application/json' -d '{ "source": { "uri": "https://github.com/intel-iot-devkit/sample-videos/blob/master/bottle-detection.mp4?raw=true", "type": "uri" }, "destination": { "type": "file", "path": "/tmp/results.txt", "format":"stream"}}'|{"objects":[{"detection":{"bounding_box":{"x_max":0.9024832248687744,"x_min":0.7928286790847778,"y_max":0.8916158676147461,"y_min":0.3036915063858032},"confidence":0.6771504878997803,"label":"bottle","label_id":5},"h":212,"roi_type":"bottle","w":71,"x":507,"y":109}],"resolution":{"height":360,"width":640},"source":"https://github.com/intel-iot-devkit/sample-videos/blob/master/bottle-detection.mp4?raw=true","timestamp":39821229050}|
|/pipelines/face_detection/1|Face Detection|curl localhost:8080/pipelines/face_detection/1 -X POST -H 'Content-Type: application/json' -d '{ "source": { "uri": "https://github.com/intel-iot-devkit/sample-videos/blob/master/head-pose-face-detection-male.mp4?raw=true", "type": "uri" }, "destination": { "type": "file", "path": "/tmp/results1.txt", "format":"stream"}}'|{"objects":[{"detection":{"bounding_box":{"x_max":0.5856688022613525,"x_min":0.4449496269226074,"y_max":0.5292186737060547,"y_min":0.1539880633354187},"confidence":0.9999995231628418,"label":"face","label_id":1},"h":162,"roi_type":"face","w":108,"x":342,"y":67}],"resolution":{"height":432,"width":768},"source":"https://github.com/intel-iot-devkit/sample-videos/blob/master/head-pose-face-detection-male.mp4?raw=true","timestamp":133583333333}
|

#### Sample Request

With the service running, initiate a request to start a pipeline using the following commands.
```bash
$ curl localhost:8080/pipelines/object_detection/1 -X POST -H 'Content-Type: application/json' -d '{ "source": { "uri": "https://github.com/intel-iot-devkit/sample-videos/blob/master/bottle-detection.mp4?raw=true", "type": "uri" }, "destination": { "type": "file", "path": "/tmp/results.txt", "format":"stream"}}'

$ tail -f /tmp/results.txt
```
Note: /tmp/results.txt cannot exist prior to running the curl command. The pipeline will not overwrite existing files and therefore will not run if it exists.
### Sample Result
```json
{"objects":[{"detection":{"bounding_box":{"x_max":0.8810903429985046,"x_min":0.77934330701828,"y_max":0.8930767178535461,"y_min":0.3040514588356018},"confidence":0.5735679268836975,"label":"bottle","label_id":5},"h":213,"roi_type":"bottle","w":65,"x":499,"y":109}],"resolution":{"height":360,"width":640},"source":"https://github.com/intel-iot-devkit/sample-videos/blob/master/bottle-detection.mp4?raw=true","timestamp":972067039}
```

### With EdgeX Integration

For testing the Video Analytics Serving with EdgeX Integration, go to demo/ directory which contains the Docker Compose file with EdgeX Fuji release and VAS services. 



|Pipeline| Description| Example Request| Example Detection|
|--------|------------|----------------|------------------|
|/pipelines/object_detection/3|Object Detection|curl localhost:8080/pipelines/object_detection/3 -X POST -H 'Content-Type: application/json' -d '{ "source": { "uri": "https://github.com/intel-iot-devkit/sample-videos/blob/master/bottle-detection.mp4?raw=true", "type": "uri" }, "destination": { "type": "mqtt", "host": "mqtt:1883", "topic":"AnalyticsData"}}'|{"id":"1359c012-1cf7-4a8a-99a2-2ea95dae6f3e","created":1585170611617,"origin":1585170611165778379,"modified":1585170611617,"device":"MQTTVideoAnalyticservice","name":"analyticsdata","value":"{\"objects\":[{\"detection\":{\"bounding_box\":{\"x_max\":0.90288245677948,\"x_min\":0.7927788496017456,\"y_max\":0.89110267162323,\"y_min\":0.30375829339027405},\"confidence\":0.6906830668449402,\"label\":\"bottle\",\"label_id\":5},\"h\":212,\"roi_type\":\"bottle\",\"w\":71,\"x\":507,\"y\":109}],\"resolution\":{\"height\":360,\"width\":640},\"source\":\"https://github.com/intel-iot-devkit/sample-videos/blob/master/bottle-detection.mp4?raw=true\",\"timestamp\":39519553072}"}|
|/pipelines/face_detection/3|Face Detection|curl localhost:8080/pipelines/face_detection/3 -X POST -H 'Content-Type: application/json' -d '{ "source": { "uri": "https://github.com/intel-iot-devkit/sample-videos/blob/master/head-pose-face-detection-male.mp4?raw=true", "type": "uri" }, "destination": { "type": "mqtt", "host": "mqtt:1883", "topic":"AnalyticsData"}}'|{"id":"ee8757a2-5e46-403d-8461-4ef5d4f6f725","created":1585170799139,"origin":1585170798607191360,"modified":1585170799139,"device":"MQTTVideoAnalyticservice","name":"analyticsdata","value":"{\"objects\":[{\"detection\":{\"bounding_box\":{\"x_max\":0.729513943195343,\"x_min\":0.5531561970710754,\"y_max\":0.5413063764572144,\"y_min\":0.23096388578414917},\"confidence\":0.993606448173523,\"label\":\"face\",\"label_id\":1},\"h\":134,\"roi_type\":\"face\",\"w\":135,\"x\":425,\"y\":100}],\"resolution\":{\"height\":432,\"width\":768},\"source\":\"https://github.com/intel-iot-devkit/sample-videos/blob/master/head-pose-face-detection-male.mp4?raw=true\",\"timestamp\":21333333333}"}
|

#### Sample Request

With the service running, initiate a request to start a pipeline using the following commands.
```bash
$ cd demo/

$ sudo docker-compose up

$ curl localhost:8080/pipelines/object_detection/3 -X POST -H 'Content-Type: application/json' -d '{ "source": { "uri": "https://github.com/intel-iot-devkit/sample-videos/blob/master/bottle-detection.mp4?raw=true", "type": "uri" }, "destination": { "type": "mqtt", "host": "mqtt:1883", "topic":"AnalyticsData"}}'

$ curl http://localhost:48080/api/v1/reading/device/MQTTVideoAnalyticservice/10 | json_pp

```

### Sample Result
```json
{"id":"1359c012-1cf7-4a8a-99a2-2ea95dae6f3e","created":1585170611617,"origin":1585170611165778379,"modified":1585170611617,"device":"MQTTVideoAnalyticservice","name":"analyticsdata","value":"{\"objects\":[{\"detection\":{\"bounding_box\":{\"x_max\":0.90288245677948,\"x_min\":0.7927788496017456,\"y_max\":0.89110267162323,\"y_min\":0.30375829339027405},\"confidence\":0.6906830668449402,\"label\":\"bottle\",\"label_id\":5},\"h\":212,\"roi_type\":\"bottle\",\"w\":71,\"x\":507,\"y\":109}],\"resolution\":{\"height\":360,\"width\":640},\"source\":\"https://github.com/intel-iot-devkit/sample-videos/blob/master/bottle-detection.mp4?raw=true\",\"timestamp\":39519553072}"}
```
