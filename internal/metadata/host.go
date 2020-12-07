package metadata

import (
	"fmt"
	"github.com/brocaar/chirpstack-gateway-bridge/internal/config"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"strconv"
	"strings"
	// "time"
)

func hostMetrics(conf config.Config) (map[string]string, error) {
	ethName := strings.ToLower(conf.MetaData.Host.Ifaces.Eth)   // "enp0s25"
	wifiName := strings.ToLower(conf.MetaData.Host.Ifaces.Wlan) // "wlp3s0"
	lteName := strings.ToLower(conf.MetaData.Host.Ifaces.Lte)   // "ppp0" // LTE modem iface
	vpnName := strings.ToLower(conf.MetaData.Host.Ifaces.Vpn)   // "ppp1" // ppp iface to CHR

	result := make(map[string]string)

	diskStat, err := disk.Usage("/")
	cpuPercentage, err := cpu.Percent(0, false)
	vmStat, err := mem.VirtualMemory()
	hostStat, err := host.Info()        // host or machine kernel, uptime, platform Info
	interfStat, err := net.Interfaces() // get interfaces MAC/hardware address
	if err != nil {
		return nil, err
	}

	result["br_ver"] = conf.General.Version
	result["disk"] = strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64)
	result["cpu"] = strconv.FormatFloat(cpuPercentage[0], 'f', 2, 64)
	result["ram"] = strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64)
	result["uptime"] = strconv.FormatUint(hostStat.Uptime, 10)

	for _, interf := range interfStat {
		switch strings.ToLower(interf.Name) {

		case ethName:
			for k, v := range ifaceMetrics(interf, "eth") {
				result[k] = v
			}

		case wifiName:
			for k, v := range ifaceMetrics(interf, "wlan") {
				result[k] = v
			}

		case lteName:
			for k, v := range ifaceMetrics(interf, "lte") {
				result[k] = v
			}

		case vpnName:
			for k, v := range ifaceMetrics(interf, "vpn") {
				result[k] = v
			}
		}
	}

	return result, nil
}

// makes key-value pares for interfaces we need to check
func ifaceMetrics(interf net.InterfaceStat, ifaceName string) map[string]string {
	result := make(map[string]string)
	if interf.HardwareAddr != "" {
		result[fmt.Sprintf("%s_mac", ifaceName)] = interf.HardwareAddr
	}

	if len(interf.Addrs) > 0 {
		ip := ""
		for idx, addr := range interf.Addrs {
			if idx > 0 {
				ip = fmt.Sprintf("%s, %s", ip, addr.Addr) // addr.String()
			} else {
				ip = addr.Addr
			}
		}
		result[fmt.Sprintf("%s_ip", ifaceName)] = ip
	}

	if len(interf.Flags) > 0 {
		result[fmt.Sprintf("%s_flags", ifaceName)] = strings.Join(interf.Flags, ", ")
	}
	return result
}
