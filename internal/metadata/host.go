package metadata

import (
	"fmt"
	"github.com/brocaar/chirpstack-gateway-bridge/internal/config"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	// "time"
)

func hostMetrics(conf config.Config) (map[string]string, error) {
	ethName := strings.ToLower(conf.MetaData.Host.Ifaces.Eth)   // ethernet interface
	wifiName := strings.ToLower(conf.MetaData.Host.Ifaces.Wlan) // 802.11 interface
	lteName := strings.ToLower(conf.MetaData.Host.Ifaces.Lte)   // LTE modem iface
	vpnName := strings.ToLower(conf.MetaData.Host.Ifaces.Vpn)   // ppp iface to CHR

	result := make(map[string]string)

	diskStat, err := disk.Usage("/")
	cpuPercentage, err := cpu.Percent(0, false)
	vmStat, err := mem.VirtualMemory()
	hostStat, err := host.Info()        // host or machine kernel, uptime, platform Info
	interfStat, err := net.Interfaces() // get interfaces MAC/hardware address
	ioStat, err := net.IOCounters(true) // get I/O counters

	if err != nil {
		return nil, err
	}
	// fmt.Println("config:", config.Config)
	result["br_ver"] = conf.General.Version
	result["disk"] = strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64)
	result["cpu"] = strconv.FormatFloat(cpuPercentage[0], 'f', 2, 64)
	result["ram"] = strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64)
	result["uptime"] = strconv.FormatUint(hostStat.Uptime, 10)

	ioMap := make(map[string]net.IOCountersStat)
	for _, io := range ioStat {
		ioMap[io.Name] = io
	}

	for _, interf := range interfStat {
		iface := ""

		switch strings.ToLower(interf.Name) {
		case ethName:
			iface = "eth"
		case wifiName:
			iface = "wlan"
		case lteName:
			iface = "lte"
		case vpnName:
			iface = "vpn"
		default:
			log.WithField("iface", interf.Name).Debug("the interface name wasn't mentioned in cfg. Skipped")
			continue
		}

		for k, v := range ifaceMetrics(interf, iface) {
			result[k] = v
		}

		if (ioMap[interf.Name] == net.IOCountersStat{}) {
			log.WithFields(log.Fields{
				"iface": interf.Name,
			}).Debug("no I/O stat for the interface")

			continue
		}
		for k, v := range netIOMetrics(ioMap[interf.Name], iface) {
			result[k] = v
		}
	}

	fields := log.Fields{}
	for k, v := range result {
		fields[k] = v
	}
	log.WithFields(fields).Debug("host metrics prepeared")

	return result, nil
}

// makes key-value pairs for interfaces we need to check
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

// makes key-value pairs of I/O counters for interfaces we need to check
func netIOMetrics(interf net.IOCountersStat, ifaceName string) map[string]string {
	result := make(map[string]string)

	result[fmt.Sprintf("%s_bytes_sent", ifaceName)] = strconv.FormatUint(interf.BytesSent, 10)
	result[fmt.Sprintf("%s_bytes_recv", ifaceName)] = strconv.FormatUint(interf.BytesRecv, 10)
	result[fmt.Sprintf("%s_packets_sent", ifaceName)] = strconv.FormatUint(interf.PacketsSent, 10)
	result[fmt.Sprintf("%s_packets_recv", ifaceName)] = strconv.FormatUint(interf.PacketsRecv, 10)
	result[fmt.Sprintf("%s_errin", ifaceName)] = strconv.FormatUint(interf.Errin, 10)
	result[fmt.Sprintf("%s_errout", ifaceName)] = strconv.FormatUint(interf.Errout, 10)
	result[fmt.Sprintf("%s_dropin", ifaceName)] = strconv.FormatUint(interf.Dropin, 10)
	result[fmt.Sprintf("%s_dropout", ifaceName)] = strconv.FormatUint(interf.Dropout, 10)
	result[fmt.Sprintf("%s_fifoin", ifaceName)] = strconv.FormatUint(interf.Fifoin, 10)
	result[fmt.Sprintf("%s_fifoout", ifaceName)] = strconv.FormatUint(interf.Fifoout, 10)

	return result
}
