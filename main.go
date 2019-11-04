package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/think-free/mqttclient"
)

// Answer the inverter answer
type Answer struct {
	Type        string        `json:"type"`
	SN          string        `json:"SN"`
	Ver         string        `json:"ver"`
	Data        []float64     `json:"Data"`
	Information []interface{} `json:"Information"`
}

func main() {

	inverter := flag.String("inverter", "5.8.8.8", "Inverter IP")
	cname := flag.String("name", "solax", "Mqtt client name")
	broker := flag.String("broker", "broker", "Mqtt broker ip")
	topic := flag.String("topic", "solax", "Mqtt base topic")
	flag.Parse()

	cli := mqttclient.NewMqttClient("Device_"+*cname, *broker)
	cli.Connect()
	cli.SendHB(*topic + "/hb")

	var name [110]string
	name[0] = "pv1_current"
	name[1] = "pv2_current"
	name[2] = "pv1_voltage"
	name[3] = "pv2_voltage"
	name[4] = "output_current_phase_1"
	name[5] = "network_voltage_phase_1"
	name[6] = "ac_power"
	name[7] = "inverter_temperature"
	name[8] = "today_energy"
	name[9] = "total_energy"
	name[10] = "exported_power"
	name[11] = "pv1_power"
	name[12] = "pv2_power"
	name[13] = "battery_voltage"
	name[14] = "battery_current"
	name[15] = "battery_power"
	name[16] = "battery_temperature"
	name[21] = "battery_remaining_capacity"
	name[41] = "total_feed_in_energy"
	name[42] = "total_consumption"
	name[43] = "power_now_phase_1"
	name[44] = "power_now_phase_2"
	name[45] = "power_now_phase_3"
	name[46] = "output_current_phase_2"
	name[47] = "output_current_phase_3"
	name[48] = "network_voltage_phase_2"
	name[49] = "network_voltage_phase_3"
	name[50] = "grid_frequency_phase_1"
	name[51] = "grid_frequency_phase_2"
	name[52] = "grid_frequency_phase_3"
	name[53] = "eps_voltage"
	name[54] = "eps_current"
	name[55] = "eps_power"
	name[56] = "eps_frequency"

	url := "http://" + *inverter + "/?optType=ReadRealTimeData"

	for {
		resp, err := http.Post(url, "", nil)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {

			log.Println(err)
		} else {

			var ans Answer
			json.Unmarshal(body, &ans)

			fmt.Println("")
			for i := 0; i < len(ans.Data); i++ {

				if name[i] != "" {
					log.Println(name[i], ":", ans.Data[i])
					cli.PublishMessage(*topic+"/"+name[i], ans.Data[i])
				}
			}
		}

		time.Sleep(time.Second * 5)
	}
}
