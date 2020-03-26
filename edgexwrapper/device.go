// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"net/url"
)

const (
	brokerUrl  = "mqtt"
	brokerPort = 1883
	username   = "admin"
	password   = "public"
)

func main() {
	go runCommandHandler()
	select {}
}

// runCommandHandler use to test receiving commands from the device service and responded back for get/set commands.
//
// Use a REST client to send a command to the service like:
// http://localhost:49982/api/v1/devices/{device id}>/message - use POST on this one with
// {"message":"some text"} in body http://localhost:49982/api/v1/devices/<device id>/ping - use GET
// http://localhost:49982/api/v1/devices/<device id>/randnum - use GET
//
// If command micro service is running, the same can be performed through command to device service
// like this http://localhost:48082/api/v1/device/<device id>/command/<command id>
//
// Requires the Device Service, Command, Core Data, Metadata and Mongo to all be running
func runCommandHandler() {
	var mqttClientId = "VideoAnalyticSubscriber"
	var qos = 0
	var topic = "AnalyticsData"

	uri := &url.URL{
		Scheme: "tcp",
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
		User:   url.UserPassword(username, password),
	}

	client, err := createMqttClient(mqttClientId, uri)
	defer client.Disconnect(5000)
	if err != nil {
		fmt.Println(err)
	}

	token := client.Subscribe(topic, byte(qos), onCommandReceivedFromBroker)
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	select {}
}

func onCommandReceivedFromBroker(client mqtt.Client, message mqtt.Message) {


var mqttClientId = "IncomingDataPublisher"
	var qos = byte(0)
	var topic = "DataTopic"

	uri := &url.URL{
		Scheme: "tcp",
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
		User:   url.UserPassword(username, password),
	}

	client, err := createMqttClient(mqttClientId, uri)
	defer client.Disconnect(5000)
	if err != nil {
		fmt.Println(err)
	}

	var data = make(map[string]interface{})
	data["name"] = "MQTTVideoAnalyticservice"
	data["cmd"] = "analyticsdata"
	data["method"] = "get"
     
   jsonStr := string(message.Payload()[:])

	data["analyticsdata"] = jsonStr
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		client.Publish(topic, qos, false, jsonData)

		fmt.Println(fmt.Sprintf("Send response: %v", string(jsonData)))
}

func sendTestData(response map[string]interface{}) {
	var mqttClientId = "ResponsePublisher"
	var qos = byte(0)
	var topic = "ResponseTopic"

	uri := &url.URL{
		Scheme: "tcp",
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
		User:   url.UserPassword(username, password),
	}

	client, err := createMqttClient(mqttClientId, uri)
	defer client.Disconnect(5000)
	if err != nil {
		fmt.Println(err)
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	client.Publish(topic, qos, false, jsonData)

	fmt.Println(fmt.Sprintf("Send response: %v", string(jsonData)))
}

func createMqttClient(clientID string, uri *url.URL) (mqtt.Client, error) {
	fmt.Println(fmt.Sprintf("Create MQTT client and connection: uri=%v clientID=%v ", uri.String(), clientID))
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s", uri.Scheme, uri.Host))
	opts.SetClientID(clientID)
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetConnectionLostHandler(func(client mqtt.Client, e error) {
		fmt.Println(fmt.Sprintf("Connection lost : %v", e))
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			fmt.Println(fmt.Sprintf("Reconnection failed : %v", e))
		} else {
			fmt.Println(fmt.Sprintf("Reconnection sucessful : %v", e))
		}
	})

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return client, token.Error()
	}

	return client, nil
}