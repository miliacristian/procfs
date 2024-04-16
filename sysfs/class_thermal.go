// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build linux
// +build linux

package sysfs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/miliacristian/procfs/internal/util"

	"fmt"
	"io/fs"

)

// ClassThermalZoneStats contains info from files in /sys/class/thermal/thermal_zone<zone>
// for a single <zone>.
// https://www.kernel.org/doc/Documentation/thermal/sysfs-api.txt
type ClassThermalZoneStats struct {
	Name    string  // The name of the zone from the directory structure.
	Type    string  // The type of thermal zone.
	Temp    int64   // Temperature in millidegree Celsius.
	Policy  string  // One of the various thermal governors used for a particular zone.
	Mode    *bool   // Optional: One of the predefined values in [enabled, disabled].
	Passive *uint64 // Optional: millidegrees Celsius. (0 for disabled, > 1000 for enabled+value)
}

// ClassThermalZoneStats returns Thermal Zone metrics for all zones.
func (fso FS) ClassThermalZoneStats() ([]ClassThermalZoneStats, error) {
	zones, err := filepath.Glob(fso.sys.Path("class/thermal/thermal_zone[0-9]*"))
	//comment
	//comment 2
	fmt.Println("called function ClassThermalZoneStats")
	if err != nil {
	    fmt.Println("message1 error function ClassThermalZoneStats")
	    if errors.As(err, new(*fs.PathError)){
            fmt.Println("message2 error function ClassThermalZoneStats")
		}
		return nil, err
	}

	fmt.Println("len zone:",len(zones))
	stats := make([]ClassThermalZoneStats, 0, len(zones))
	fmt.Println("stats before:",stats)
	for _, zone := range zones {
	    fmt.Println("for zone:",zone)
		zoneStats, err := parseClassThermalZone(zone)
		fmt.Println("zoneStats:",zoneStats)
		if err != nil {
		    fmt.Println("message2.1 error function ClassThermalZoneStats")
			if errors.Is(err, syscall.ENODATA) {
				continue
			}
			fmt.Println("message3 error function ClassThermalZoneStats")
            if errors.As(err, new(*fs.PathError)){
                fmt.Println("message4 error function ClassThermalZoneStats")
            }
			return nil, err
		}
		zoneStats.Name = strings.TrimPrefix(filepath.Base(zone), "thermal_zone")
		stats = append(stats, zoneStats)
	}
	fmt.Println("stats after:",stats)
	return stats, nil
}

func parseClassThermalZone(zone string) (ClassThermalZoneStats, error) {
	// Required attributes.
	zoneType, err := util.SysReadFile(filepath.Join(zone, "type"))
	if err != nil {
	    fmt.Println("message5 error function ClassThermalZoneStats")
	    if errors.As(err, new(*fs.PathError)){
            fmt.Println("message6 error function ClassThermalZoneStats")
		}
		return ClassThermalZoneStats{}, err
	}
	zonePolicy, err := util.SysReadFile(filepath.Join(zone, "policy"))
	if err != nil {
	    fmt.Println("message7 error function ClassThermalZoneStats")
	    if errors.As(err, new(*fs.PathError)){
            fmt.Println("message8 error function ClassThermalZoneStats")
		}
		return ClassThermalZoneStats{}, err
	}
	zoneTemp, err := util.ReadIntFromFile(filepath.Join(zone, "temp"))
	if err != nil {
	    fmt.Println("message9 error function ClassThermalZoneStats")
	    if errors.As(err, new(*fs.PathError)){
            fmt.Println("message10 error function ClassThermalZoneStats")
		}
		return ClassThermalZoneStats{}, err
	}

	// Optional attributes.
	mode, err := util.SysReadFile(filepath.Join(zone, "mode"))
	if err != nil && !os.IsNotExist(err) && !os.IsPermission(err) {
	    fmt.Println("message11 error function ClassThermalZoneStats")
	    if errors.As(err, new(*fs.PathError)){
            fmt.Println("message12 error function ClassThermalZoneStats")
		}
		return ClassThermalZoneStats{}, err
	}
	zoneMode := util.ParseBool(mode)

	var zonePassive *uint64
	passive, err := util.ReadUintFromFile(filepath.Join(zone, "passive"))
	if os.IsNotExist(err) || os.IsPermission(err) {
		zonePassive = nil
	} else if err != nil {
	    fmt.Println("message13 error function ClassThermalZoneStats")
	    if errors.As(err, new(*fs.PathError)){
            fmt.Println("message14 error function ClassThermalZoneStats")
		}
		return ClassThermalZoneStats{}, err
	} else {
		zonePassive = &passive
	}

	return ClassThermalZoneStats{
		Type:    zoneType,
		Policy:  zonePolicy,
		Temp:    zoneTemp,
		Mode:    zoneMode,
		Passive: zonePassive,
	}, nil
}
