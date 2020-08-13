package blocks

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/base/watchers/netlink"
	"barista.run/colors"
	"barista.run/format"
	"barista.run/modules/battery"
	"barista.run/modules/bluetooth"
	"barista.run/modules/cpuload"
	"barista.run/modules/meminfo"
	"barista.run/modules/netinfo"
	"barista.run/modules/volume"
	"barista.run/modules/wlan"
	"barista.run/outputs"
	"barista.run/pango"
	"barista.run/pango/icons/material"
	"github.com/glebtv/custom_barista/kbdlayout"
	"github.com/martinlindhe/unit"
)

var spacer = pango.Text("   ").XXSmall()

func home(path string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(usr.HomeDir, path)
}

func init() {
	colors.LoadFromMap(map[string]string{
		"good":     "#6d6",
		"degraded": "#dd6",
		"bad":      "#d66",
		"dim-icon": "#777",
	})

	material.Load(home("git/github.com/material-design-icons"))
}

// Clock ...
func Clock(now time.Time) bar.Output {
	return outputs.Pango(
		pango.Icon("material-today"),
		pango.Icon("material-access-time"),
		now.Format("Mon 2 Jan "),
		now.Format("15:04"),
	).OnClick(click.RunLeft("gsimplecal")).Color(colors.Scheme("dim-icon"))
}

// Bat ...
func Bat(i battery.Info) bar.Output {
	if i.Status == battery.Disconnected || i.Status == battery.Unknown {
		return outputs.Textf("%v", i.Status).Urgent(true)
	}

	disp := pango.Textf("%d%%", i.RemainingPct())
	icon := pango.Icon("material-battery-std")
	cl := colors.Scheme("dim-icon")

	if i.Status == battery.Charging {
		icon = pango.Icon("material-battery-charging-full")
	}

	var urgent bool
	switch {
	case i.RemainingPct() <= 10:
		urgent = true
	case i.RemainingPct() <= 20:
		cl = colors.Scheme("bad")
	case i.RemainingPct() <= 30:
		cl = colors.Scheme("degraded")
	}
	return outputs.Pango(icon, disp).Color(cl).Urgent(urgent)
}

// Snd ...
func Snd(v volume.Volume) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-volume-down")

	if v.Pct() > 66 {
		ic = pango.Icon("material-volume-up")
		// cl = colors.Scheme("degraded")
	}
	if v.Mute || v.Pct() == 0 {
		ic = pango.Icon("material-volume-off")
		// cl = colors.Scheme("dim-icon")
	}

	return outputs.Pango(ic, spacer, pango.Textf("%2d%%", v.Pct())).Color(cl)
}

// Brightness ...
func Brightness(i int) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-brightness-medium")

	if i > 50 {
		ic = pango.Icon("material-brightness-high")
	}
	return outputs.Pango(ic, spacer, fmt.Sprintf("%2.0d%%", i)).Color(cl)
}

// Layout ...
func Layout(m *kbdlayout.Module, i kbdlayout.Info) bar.Output {
	la := strings.ToLower(i.Layout)
	ic := pango.Icon("material-language")
	c := colors.Scheme("dim-icon")
	if la != "us" {
		c = colors.Scheme("degraded")
	}
	return outputs.Pango(ic, spacer, fmt.Sprintf("%s", la)).OnClick(m.Click).Color(c)
}

// Net ...
func Net(s netinfo.State) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-settings-ethernet")

	if len(s.IPs) >= 1 {
		disp := pango.Textf(fmt.Sprintf("%s:%v", s.Name, s.IPs[0]))
		return outputs.Pango(ic, spacer, disp).Color(cl)
	}
	return nil
}

// Yubi ...
func Yubi(x bool, y bool) bar.Output {
	if x {
		return outputs.Textf("g").Background(colors.Scheme("dim-icon")).MinWidth(200)
	}
	if y {
		return outputs.Textf("s").Background(colors.Scheme("dim-icon")).MinWidth(200)
	}
	return nil
}

// WLAN ...
func WLAN(i wlan.Info) bar.Output {
	disp := pango.Textf(fmt.Sprintf("%s", i.SSID))
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-signal-wifi-4-bar")

	if i.State == netlink.Down || !i.Enabled() || i.Connecting() || !i.Connected() {
		return nil
	}

	return outputs.Pango(ic, spacer, disp).Color(cl)
}

// Bluetooth ...
func Bluetooth(s bluetooth.AdapterInfo) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-bluetooth")

	if !s.Powered {
		return nil
		// cl = colors.Scheme("dim-icon")
	}
	return outputs.Pango(ic).Color(cl)
}

// Blue ...
func Blue(i bluetooth.DeviceInfo) bar.Output {
	dp := pango.Textf(fmt.Sprintf("%s", i.Alias))
	cl := colors.Scheme("good")
	ic := pango.Icon("material-bluetooth-connected")

	if !i.Connected {
		return nil
	}
	return outputs.Pango(ic, spacer, dp).Color(cl)
}

// Mem ...
func Mem(i meminfo.Info) bar.Output {
	cl := colors.Scheme("dim-icon")
	ic := pango.Icon("material-memory")
	dp := pango.Textf(`%s/%s`,
		format.IBytesize(i["MemTotal"]-i["MemAvailable"]),
		format.IBytesize(i["MemTotal"]))

	// switch {
	// case i["MemAvailable"] < 15000000.0:
	// 	cl = colors.Scheme("degraded")
	// case i["MemAvailable"] < 0.33:
	// 	cl = colors.Scheme("bad")
	// }
	return outputs.Pango(ic, spacer, dp).Color(cl)
}

// Cpu ...
func Cpu(i cpuload.LoadAvg) bar.Output {
	dp := pango.Textf(fmt.Sprintf("%0.2f", i.Min1()))
	ic := pango.Icon("material-memory")
	cl := colors.Scheme("dim-icon")

	switch {
	case i.Min1() >= 5.0:
		cl = colors.Scheme("degraded")
	case i.Min1() >= 7.0:
		cl = colors.Scheme("bad")
	}

	return outputs.
		Pango(ic, spacer, dp).Color(cl)
}

// Temp ...
func Temp(i unit.Temperature) bar.Output {
	dp := pango.Textf(fmt.Sprintf("%.0f C", i.Celsius()))
	ic := pango.Icon("material-memory")
	cl := colors.Scheme("dim-icon")

	return outputs.
		Pango(ic, spacer, dp).Color(cl)
}

// Snd2 ...
// func Snd2(v volume.Volume) bar.Output {
// 	cl := colors.Scheme("good")
// 	ic := pango.Icon("material-volume-up")
// 	pct := v.Pct()
// 	if v.Mute {
// 		ic = pango.Icon("material-volume-off")
// 		cl = colors.Scheme("bad")
// 	}

// 	return outputs.
// 		Pango(ic, spacer, pango.Textf("%2d%%", pct)).Color(cl)
// }
