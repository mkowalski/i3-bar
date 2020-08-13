package main

import (
	"time"

	"barista.run"
	"barista.run/modules/battery"
	"barista.run/modules/bluetooth"
	"barista.run/modules/clock"
	"barista.run/modules/cpuload"
	"barista.run/modules/meminfo"
	"barista.run/modules/netinfo"
	"barista.run/modules/volume"
	"barista.run/modules/volume/alsa"
	"barista.run/modules/wlan"

	//"github.com/glebtv/custom_barista/kbdlayout"
	"github.com/karampok/i3-bar/blocks"
	"github.com/karampok/i3-bar/xbacklight"
)

func main() {
	//ly := kbdxlayout.New().Output(blocks.Layout)
	br := xbacklight.New().Output(blocks.Brightness)
	snd := volume.New(alsa.DefaultMixer()).Output(blocks.Snd)
	bat := battery.All().Output(blocks.Bat)
	cpu := cpuload.New().Output(blocks.Cpu)
	mem := meminfo.New().Output(blocks.Mem)
	// temp := cputemp.New().Output(blocks.Temp)
	wi := wlan.Named("wlp2s0").Output(blocks.WLAN)
	nt1 := netinfo.Interface("wlp2s0").Output(blocks.Net)
	nt3 := netinfo.Interface("ppp0").Output(blocks.Net)
	nt4 := netinfo.Interface("tun0").Output(blocks.Net)
	ti := clock.Local().Output(time.Second, blocks.Clock)

	adapter, mac, _ := "hci0", "09:A5:C1:A6:5C:77", "bluez_sink.09_A5_C1_A6_5C_77.headset_head_unit"
	blD := bluetooth.Device(adapter, mac).Output(blocks.Blue)
	bl := bluetooth.DefaultAdapter().Output(blocks.Bluetooth)

	panic(barista.Run(
		blD, br, snd, bat, cpu, mem, wi, nt1, nt3, nt4, bl, ti,
	))

}
