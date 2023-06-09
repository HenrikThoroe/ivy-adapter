package sys

import (
	"crypto/sha256"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jaypipes/ghw"
	"github.com/klauspost/cpuid/v2"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// DeviceIdentifier contains fields that uniquely identify a device.
type DeviceIdentifier struct {
	Name string `json:"name"` // The name of the device.
	ID   string `json:"id"`   // The ID of the device.
}

// CPU contains information about the CPU(s) of a device.
type CPU struct {
	Cores        int      `json:"cores,omitempty"`        // The number of cores of the CPU.
	Model        string   `json:"model"`                  // The model of the CPU.
	Vendor       string   `json:"vendor"`                 // The vendor of the CPU.
	Threads      int      `json:"threads,omitempty"`      // The number of threads of the CPU.
	Capabilities []string `json:"capabilities,omitempty"` // The capabilities of the CPU.
}

// GPU contains information about the GPU(s) of a device.
type GPU struct {
	Model  string `json:"model"`            // The model of the GPU.
	Vendor string `json:"vendor"`           // The vendor of the GPU.
	Memory int    `json:"memory,omitempty"` // The memory of the GPU in bytes.
}

// Device contains information about a device.
type Device struct {
	Cpu    []CPU  `json:"cpu"`              // The CPU(s) of the device.
	Memory int    `json:"memory,omitempty"` // The memory of the device in bytes.
	Gpu    []GPU  `json:"gpu"`              // The GPU(s) of the device.
	Model  string `json:"model"`            // The model of the device.
	OS     string `json:"os"`               // The operating system of the device.
	Arch   string `json:"arch"`             // The architecture of the device.
}

// DeviceInfo returns information about the device.
func DeviceInfo() (Device, DeviceIdentifier) {
	h := Device{
		Cpu: make([]CPU, 0),
		Gpu: make([]GPU, 0),
	}
	id := DeviceIdentifier{}

	loadGPU(&h)
	loadCPU(&h)
	loadMemory(&h)
	loadHost(&h, &id.Name, &id.ID)

	return h, id
}

func loadGPU(h *Device) {
	if gpu, err := ghw.GPU(); err == nil && gpu != nil {
		for _, card := range gpu.GraphicsCards {
			if card == nil {
				continue
			}

			if info := card.DeviceInfo; info != nil {
				gpu := GPU{}

				if info.Product != nil {
					gpu.Model = info.Product.Name
				}

				if info.Vendor != nil {
					gpu.Vendor = info.Vendor.Name
				}

				if info.Node != nil && info.Node.Memory != nil {
					gpu.Memory = int(info.Node.Memory.TotalPhysicalBytes)
				}

				h.Gpu = append(h.Gpu, gpu)
			}
		}
	}
}

func loadCPU(h *Device) {
	h.Cpu = append(h.Cpu, CPU{
		Cores:        cpuid.CPU.PhysicalCores,
		Model:        cpuid.CPU.BrandName,
		Vendor:       cpuid.CPU.VendorString,
		Threads:      cpuid.CPU.LogicalCores,
		Capabilities: cpuid.CPU.FeatureSet(),
	})
}

func loadMemory(h *Device) {
	if mem, err := ghw.Memory(); err == nil && mem != nil {
		h.Memory = int(mem.TotalUsableBytes)
		return
	}

	if mem, err := mem.VirtualMemory(); err == nil && mem != nil {
		h.Memory = int(mem.Available)
		return
	}
}

func loadHost(h *Device, name *string, id *string) {
	h.OS = runtime.GOOS
	h.Arch = runtime.GOARCH

	if hn, err := os.Hostname(); err == nil && hn != "" {
		*name = hn
	}

	if info, err := ghw.Product(); err == nil && info != nil {
		h.Model = info.Name

		if *name == "" {
			*name = info.Name + " (" + info.Vendor + ")"
		}

		if info.UUID != "" {
			hash := sha256.New()
			hash.Write([]byte(info.UUID))
			*id = strings.ToLower(fmt.Sprintf("%x", hash.Sum(nil)))
		}

		return
	}

	if info, err := host.Info(); err == nil && info != nil {
		h.Model = info.PlatformFamily + " " + info.PlatformVersion + " (" + info.OS + ")"

		if *name == "" {
			*name = info.Hostname
		}

		if info.HostID != "" {
			hash := sha256.New()
			hash.Write([]byte(info.HostID))
			*id = strings.ToLower(fmt.Sprintf("%x", hash.Sum(nil)))
		}

		return
	}
}
