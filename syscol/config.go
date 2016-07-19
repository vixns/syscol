/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package syscol

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"os"
	"regexp"
	"time"
)

var Logger log.LoggerInterface

var Config *config = &config{
	FrameworkName:     "syscol",
	FrameworkRole:     "*",
	Cpus:              0.1,
	Mem:               64,
	ReportingInterval: 1 * time.Second,
	Transform:         "none",
	LogLevel:          "info",
	Principal:		   "",
	SecretFile:		   "",
}

var executorMask = regexp.MustCompile("executor.*")

type config struct {
	Api                string
	Master             string
	FrameworkName      string
	FrameworkRole      string
	User               string
	Cpus               float64
	Mem                float64
	ReportingInterval  time.Duration
	Executor           string
	ProducerProperties string
	Topic              string
	Transform          string // none, avro, proto
	SchemaRegistryUrl  string
	Namespace          string
	LogLevel           string
	Principal          string
	SecretFile         string
}

func (c *config) CanStart() bool {
	if c.Transform == TransformAvro && c.SchemaRegistryUrl == "" {
		return false
	}
	return c.ProducerProperties != "" && c.Topic != ""
}

func (c *config) Read(task *mesos.TaskInfo) {
	config := new(config)
	Logger.Debugf("Task data: %s", string(task.GetData()))
	err := json.Unmarshal(task.GetData(), config)
	if err != nil {
		Logger.Critical(err)
		os.Exit(1)
	}
	*c = *config
}

func (c *config) String() string {
	return fmt.Sprintf(`api:                 %s
master:              %s
framework name:      %s
framework role:      %s
user:                %s
cpus:                %.2f
mem:                 %.2f
reporting interval:  %s
executor:            %s
producer properties: %s
topic:               %s
transform:           %s
namespace:           %s
authentication principal:           %s
authentication secret file:           %s
log level:           %s
`, c.Api, c.Master, c.FrameworkName, c.FrameworkRole, c.User, c.Cpus, c.Mem, c.ReportingInterval,
		c.Executor, c.ProducerProperties, c.Topic, c.Transform, c.Namespace, c.Principal, c.SecretFile, c.LogLevel)
}

func InitLogging(level string) error {
	config := fmt.Sprintf(`<seelog minlevel="%s">
    <outputs formatid="main">
        <console />
    </outputs>

    <formats>
        <format id="main" format="%%Date/%%Time [%%LEVEL] %%Msg%%n"/>
    </formats>
</seelog>`, level)

	logger, err := log.LoggerFromConfigAsBytes([]byte(config))
	Config.LogLevel = level
	Logger = logger

	return err
}
