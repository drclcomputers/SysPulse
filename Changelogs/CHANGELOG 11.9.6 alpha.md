# SysPulse Changelog - Version 11.9.6 alpha

## 🚀 Performance Improvements

### System Call Optimization
- **Process Information Caching**
  - Added intelligent caching system for process information
  - Reduced repeated syscalls for frequently accessed process data
  - Configurable cache TTL (Time To Live) settings
  - Improved process details view performance

### New Performance Features
- **Syscall Batching System**
  - Implemented new SyscallBatcher utility
  - Groups related system calls to reduce kernel transitions
  - Configurable batch sizes and intervals
  - Optimized multi-threaded execution of batched operations

### Configuration Enhancements
- **Added Performance Configuration Section**
  ```json
  "performance": {
    "process_cache_ttl": 2,
    "full_scan_interval": 10,
    "syscall_batch_size": 100,
    "process_update_interval": 2
  }
  ```
  - Customizable process cache duration
  - Adjustable full system scan intervals
  - Configurable syscall batch sizes
  - Tunable process update frequency

### Code Structure Improvements
- **New Caching Infrastructure**
  - Added `ProcessCache` type for efficient process data storage
  - Implemented thread-safe cache operations
  - Added automatic cache invalidation
  - Reduced memory footprint for process monitoring

### Technical Details
- Reduced syscall frequency through intelligent caching
- Implemented concurrent batch processing of system calls
- Added performance monitoring and optimization utilities
- Enhanced process information gathering efficiency

## 🔧 How to Configure
Performance settings can be adjusted in `config.json` under the new `performance` section to fine-tune system resource usage according to your needs.

## 💡 Notes
- Default cache TTL is set to 2 seconds
- Full system scans occur every 10 seconds by default
- Process updates are batched in groups of 100 by default
- All performance settings are configurable through the configuration file

## Summary
- Some performance benefits. CPU usage has been brought down by 1-2%, while ram increased with 0.5-1MB.