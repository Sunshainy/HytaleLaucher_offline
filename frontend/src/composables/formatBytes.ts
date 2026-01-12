/**
 * Formats a byte count into a human-readable string.
 * @param bytes - The number of bytes to format
 * @returns A formatted string like "1.5 MB"
 */
export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * Formats a bytes-per-second value into a human-readable speed string.
 * @param bps - The speed in bytes per second
 * @returns A formatted string like "1.5 MB/s"
 */
export function formatSpeed(bps: number): string {
  return formatBytes(bps) + '/s'
}

/**
 * Formats a duration in seconds into a human-readable string.
 * @param seconds - The duration in seconds
 * @returns A formatted string like "2h 30m"
 */
export function formatDuration(seconds: number): string {
  if (seconds < 60) {
    return `${Math.round(seconds)}s`
  }
  if (seconds < 3600) {
    const minutes = Math.floor(seconds / 60)
    const secs = Math.round(seconds % 60)
    return secs > 0 ? `${minutes}m ${secs}s` : `${minutes}m`
  }
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return minutes > 0 ? `${hours}h ${minutes}m` : `${hours}h`
}

export default formatBytes
