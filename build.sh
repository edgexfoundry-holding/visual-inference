#!/bin/bash -e

VIDEO_ANALYTICS_SERVICE_REPO="https://github.com/intel/video-analytics-serving.git"

#sudo docker build ${VIDEO_ANALYTICS_SERVICE_REPO}#0.2.3-alpha -f Dockerfile.gst.base -t video_analytics_serving_gstreamer_base 

#sudo docker build ${VIDEO_ANALYTICS_SERVICE_REPO}#0.2.3-alpha -f Dockerfile.gst -t video_analytics_serving_gstreamer 

sudo docker build -t video_analytics_serving_gstreamer_edgex:0.2.2 .
