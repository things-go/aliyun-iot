package mock

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/go-core-package/lib/logger"

	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
)

// just for test
const (
	ProductKey    = "a1NHTWLlMny"
	ProductSecret = "qv30SFtpGf3tSBfP"
	DeviceName    = "mygateway"
	DeviceSecret  = "6ed85dfb2ec5cd3746104cc3b2e0b188"
)

const (
	SensorProductKey    = "a15aMYCIe4I"
	SensorProductSecret = "fkrQaPSraQTHcXbQ"
	SensorDeviceName    = "mysensor"
	SensorDeviceSecret  = "1a5a0ffd6b06e402740067cbbc7de2ee"
)

var MetaTetrad = infra.MetaTetrad{
	ProductKey,
	ProductSecret,
	DeviceName,
	DeviceSecret,
}

var MetaTriad = infra.MetaTriad{
	ProductKey,
	DeviceName,
	DeviceSecret,
}

var SensorTerad = infra.MetaTetrad{
	"a15aMYCIe4I",
	"fkrQaPSraQTHcXbQ",
	"mysensor",
	"",
}

var SensorTriad = infra.MetaTriad{
	"a15aMYCIe4I",
	"mysensor",
	"1a5a0ffd6b06e402740067cbbc7de2ee",
}

type Instance struct {
	GatewayVersion string  `json:"gateway_version"`
	SystemInfo     string  `json:"system_info"`
	CpuUsage       float32 `json:"cpu_usage"`
	MemoryUsage    float32 `json:"memory_usage"`
	MemoryTotal    float64 `json:"memory_total"`
	MemoryFree     float64 `json:"memory_free"`
	CpuCoreNumber  int32   `json:"cpu_core_number"`
	DiskUsage      float32 `json:"disk_usage"`
	LightSwitch    int     `json:"light_switch"`
}

type Alarm struct {
	High int `json:"high"`
}

type SensorInstance struct {
	CurrentTemperature float64 `json:"CurrentTemperature"`
	CurrentHumidity    float64 `json:"CurrentHumidity"`
}

func Init(triad infra.MetaTriad) *aiot.MQTTClient {
	signs, err := sign.Generate(triad, infra.CloudRegionDomain{Region: infra.CloudRegionShangHai}, sign.WithTimestamp())
	if err != nil {
		panic(err)
	}

	opts := mqtt.NewClientOptions().
		AddBroker(signs.Addr).
		SetClientID(signs.ClientIDWithExt()).
		SetUsername(signs.UserName).
		SetPassword(signs.Password).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(cli mqtt.Client) {
			log.Println("mqtt client connection success")
		}).
		SetConnectionLostHandler(func(cli mqtt.Client, err error) {
			log.Println("mqtt client connection lost, ", err)
		})

	mqc := mqtt.NewClient(opts)
	mqc.Connect().Wait()

	client := aiot.NewWithMQTT(
		triad,
		mqc,
		aiot.WithEnableNTP(),
		aiot.WithEnableDesired(),
		aiot.WithEnableDiag(),
		aiot.WithEnableGateway(),
		aiot.WithCallback(mockCb{}),
		aiot.WithGwCallback(mockCb{}),
		aiot.WithLogger(logger.New(log.New(os.Stdout, "mqtt --> ", log.LstdFlags), logger.WithEnable(true))),
	)

	if err = client.Connect(); err != nil {
		panic(err)
	}
	return client
}

func InstanceValue() Instance {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	return Instance{
		GatewayVersion: "v0.1.0",
		SystemInfo:     "test gateway",
		CpuUsage:       float32(rand.Intn(1000)) / 10,
		MemoryUsage:    float32(stats.HeapInuse / stats.HeapSys),
		MemoryTotal:    float64(stats.HeapSys),
		MemoryFree:     float64(stats.HeapIdle),
		CpuCoreNumber:  int32(runtime.NumCPU()),
		DiskUsage:      float32(rand.Intn(1000)) / 10,
		LightSwitch:    rand.Intn(2),
	}
}

func SensorInstanceValue() SensorInstance {
	return SensorInstance{
		float64(rand.Intn(160) - 40),
		float64(rand.Intn(100)),
	}
}

func ThingEventPropertyPost(client *aiot.MQTTClient) {
	for {
		err := client.LinkThingEventPropertyPost(ProductKey, DeviceName, InstanceValue(), 3*time.Second)
		if err != nil {
			log.Printf("error: %#v", err)
		}
		time.Sleep(time.Second * 30)
	}
}

type mockCb struct {
	aiot.NopCb
	aiot.NopGwCb
}

func (sf mockCb) RRPCRequest(c *aiot.Client, messageID, productKey, deviceName string, payload []byte) error {
	req := &aiot.Request{}
	if err := json.Unmarshal(payload, req); err != nil {
		return err
	}
	c.Log.Debugf("rrpc.resopnse.%s", messageID)
	c.Log.Debugf("%+v", req)
	return c.RRPCResponse(productKey, deviceName, messageID, aiot.Response{
		ID:   req.ID,
		Code: infra.CodeSuccess,
		Data: "{}",
	})
}

func (sf mockCb) ThingTopoChange(c *aiot.Client, params aiot.TopoChangeParams) error {
	c.Log.Debugf("%+v", params)
	return nil
}
