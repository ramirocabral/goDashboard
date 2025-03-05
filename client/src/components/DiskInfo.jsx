"use client"

import { useState, useEffect } from "react"
import { useWebSocket } from "../contexts/WebSocketContext"

const DiskInfo = () => {
  const { disksData, fetchStaticData } = useWebSocket()
  const [staticDisksInfo, setStaticDisksInfo] = useState(null)

  useEffect(() => {
    const fetchDisksInfo = async () => {
      const data = await fetchStaticData("disks/info")
      if (data) {
        setStaticDisksInfo(data)
      }
    }

    fetchDisksInfo()
  }, [fetchStaticData])

  // Combine real-time and static data
  const combinedData = {
    ...staticDisksInfo,
    ...disksData,
  }

  // Format bytes to human-readable format
  const formatBytes = (bytes, decimals = 2) => {
    if (!bytes) return "0 B"

    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ["B", "KB", "MB", "GB", "TB"]

    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${Number.parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
  }

  // Calculate usage percentage
  const calculatePercentage = (used, total) => {
    if (!used || !total) return 0
    return ((used / total) * 100).toFixed(2)
  }

  // Get disks from combined data or use placeholder
  const disks = combinedData?.disks || [
    {
      device: "/dev/sda",
      mount_point: "/",
      fs_type: "ext4",
      total_space: 1000000000000,
      used_space: 250000000000,
      available_space: 750000000000,
      inodes_total: 61054976,
      inodes_used: 3276800,
      inodes_available: 57778176,
    },
    {
      device: "/dev/sdb",
      mount_point: "/data",
      fs_type: "xfs",
      total_space: 2000000000000,
      used_space: 500000000000,
      available_space: 1500000000000,
      inodes_total: 122109952,
      inodes_used: 1638400,
      inodes_available: 120471552,
    },
  ]

  return (
    <div className="space-y-6">
      <div className="card">
        <div className="card-header">
          <div className="card-title">Disk Information</div>
          <div className="card-description">Storage devices and mount points</div>
        </div>
        <div className="card-content">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left pb-2">Device</th>
                  <th className="text-left pb-2">Mount Point</th>
                  <th className="text-left pb-2">File System</th>
                  <th className="text-left pb-2">Size</th>
                  <th className="text-left pb-2">Used</th>
                  <th className="text-left pb-2">Available</th>
                  <th className="text-left pb-2">Usage</th>
                </tr>
              </thead>
              <tbody>
                {disks.map((disk, index) => (
                  <tr key={index} className="border-b border-border">
                    <td className="py-2">{disk.device}</td>
                    <td className="py-2">{disk.mount_point}</td>
                    <td className="py-2">{disk.fs_type}</td>
                    <td className="py-2">{formatBytes(disk.total_space)}</td>
                    <td className="py-2">{formatBytes(disk.used_space)}</td>
                    <td className="py-2">{formatBytes(disk.available_space)}</td>
                    <td className="py-2">
                      <div className="flex items-center">
                        <div className="w-24 bg-secondary rounded-full h-2 mr-2">
                          <div
                            className="bg-primary h-2 rounded-full"
                            style={{ width: `${calculatePercentage(disk.used_space, disk.total_space)}%` }}
                          ></div>
                        </div>
                        <span>{calculatePercentage(disk.used_space, disk.total_space)}%</span>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <div className="card-title">Inode Usage</div>
          <div className="card-description">File system inode information</div>
        </div>
        <div className="card-content">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left pb-2">Mount Point</th>
                  <th className="text-left pb-2">Total Inodes</th>
                  <th className="text-left pb-2">Used Inodes</th>
                  <th className="text-left pb-2">Available Inodes</th>
                  <th className="text-left pb-2">Usage</th>
                </tr>
              </thead>
              <tbody>
                {disks.map((disk, index) => (
                  <tr key={index} className="border-b border-border">
                    <td className="py-2">{disk.mount_point}</td>
                    <td className="py-2">{disk.inodes_total?.toLocaleString() || "N/A"}</td>
                    <td className="py-2">{disk.inodes_used?.toLocaleString() || "N/A"}</td>
                    <td className="py-2">{disk.inodes_available?.toLocaleString() || "N/A"}</td>
                    <td className="py-2">
                      <div className="flex items-center">
                        <div className="w-24 bg-secondary rounded-full h-2 mr-2">
                          <div
                            className="bg-primary h-2 rounded-full"
                            style={{ width: `${calculatePercentage(disk.inodes_used, disk.inodes_total)}%` }}
                          ></div>
                        </div>
                        <span>{calculatePercentage(disk.inodes_used, disk.inodes_total)}%</span>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <div className="card-title">Disk Partitions</div>
          <div className="card-description">Physical disk partitions</div>
        </div>
        <div className="card-content">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left pb-2">Device</th>
                  <th className="text-left pb-2">Type</th>
                  <th className="text-left pb-2">Size</th>
                  <th className="text-left pb-2">Mount Point</th>
                  <th className="text-left pb-2">UUID</th>
                </tr>
              </thead>
              <tbody>
                {(
                  combinedData?.partitions || [
                    {
                      device: "/dev/sda1",
                      type: "ext4",
                      size: "512 MB",
                      mount_point: "/boot",
                      uuid: "123e4567-e89b-12d3-a456-426614174000",
                    },
                    {
                      device: "/dev/sda2",
                      type: "swap",
                      size: "8 GB",
                      mount_point: "[SWAP]",
                      uuid: "123e4567-e89b-12d3-a456-426614174001",
                    },
                    {
                      device: "/dev/sda3",
                      type: "ext4",
                      size: "991.5 GB",
                      mount_point: "/",
                      uuid: "123e4567-e89b-12d3-a456-426614174002",
                    },
                    {
                      device: "/dev/sdb1",
                      type: "xfs",
                      size: "2 TB",
                      mount_point: "/data",
                      uuid: "123e4567-e89b-12d3-a456-426614174003",
                    },
                  ]
                ).map((partition, index) => (
                  <tr key={index} className="border-b border-border">
                    <td className="py-2">{partition.device}</td>
                    <td className="py-2">{partition.type}</td>
                    <td className="py-2">{partition.size}</td>
                    <td className="py-2">{partition.mount_point}</td>
                    <td className="py-2">{partition.uuid}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  )
}

export default DiskInfo

