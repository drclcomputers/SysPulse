# SysPulse Changelog

## Version 11.9.4 Alpha - Modal Enhancement Release
*Release Date: July 21, 2025*

### ✨ New Features

#### Enhanced Modal System
- **Scrollable CPU information modal** - CPU details are now displayed in a scrollable TextView for better readability with large amounts of processor information
- **Scrollable disk information modal** - Disk partition details can now be scrolled through for systems with many partitions and extensive filesystem data
- **Scrollable disk I/O modal** - Disk I/O statistics are now presented in a scrollable format for better navigation through detailed performance metrics
- **Improved terminal compatibility** - Enhanced border fallback system for better visibility on basic terminals like Ubuntu Server and SSH sessions
- **Multiple widgets selection keys** - Instead of TAB and Shift+TAB, you can now use the Left and Right Arrow Keys to navigate around.

### 🔧 Technical Implementation

#### Modal Information System
- **Unified scrollable modals** - All information modals now use consistent scrollable TextView implementation instead of basic modal dialogs
- **Better space utilization** - Large amounts of system information no longer overflow modal boundaries or get truncated
- **Enhanced navigation** - Users can scroll through detailed system information using arrow keys and standard navigation controls

#### Terminal Compatibility Enhancement
- **ASCII border fallbacks** - Improved detection and fallback for terminals that don't support Unicode box drawing characters
- **SSH session detection** - Better handling of remote terminal sessions with limited display capabilities
- **Terminal type detection** - Enhanced detection of basic terminal types (linux, vt100, vt102, vt220, ansi, dumb) for appropriate fallback rendering

### 🎯 User Experience Improvements

#### Information Access
- **No more truncated data** - All CPU specifications, disk partition details, and I/O statistics are now fully accessible through scrolling
- **Better readability** - Large information displays are properly formatted within scrollable containers with consistent styling
- **Consistent interaction** - All modals now follow the same scrollable pattern for familiar user experience across widgets

#### Terminal Support
- **Ubuntu Server compatibility** - Widgets now display visible borders on Ubuntu Server and other terminal-only distributions
- **Enhanced focus indication** - Improved visual feedback when navigating between widgets on systems with limited terminal capabilities
- **Cross-platform consistency** - Unified appearance across different terminal emulators and operating systems

### 🔄 Migration Notes

#### Automatic Updates
- **No configuration changes required** - All improvements are automatically applied without user intervention
- **Backward compatibility maintained** - Existing themes and configurations continue to work without modification
- **Enhanced functionality** - Users will immediately benefit from improved modal readability and better terminal support

#### Behavior Changes
- **Improved modal interaction** - Information modals now use scrollable text views instead of static modal dialogs
- **Better border visibility** - Terminals that don't support Unicode will automatically use ASCII borders with bold attributes
- **Enhanced navigation** - Modal content can be scrolled and navigated more intuitively

### 📋 Affected Components

#### Modal System
- `internal/services/UI/widgets.go` - Updated CPU, disk, and disk I/O modals to use scrollable TextView components
- `internal/utils/util.go` - Enhanced terminal capability detection and border fallback functions

#### Terminal Compatibility
- `internal/services/UI/widgets.go` - Applied `utils.SetBorderStyle()` to ProcessTree and Battery widgets
- `internal/plugins/example.go` - Updated plugin widgets to use fallback border styles
- `internal/plugins/docker.go` - Enhanced Docker plugin widgets with terminal compatibility
- `internal/plugins/integration.go` - Applied fallback styling to plugin containers

### 🧪 Testing Recommendations

#### Modal Testing
1. **Open system information modals** (Press 'I' on CPU, Disk, or any widget)
2. **Test scrolling functionality** - Use arrow keys to navigate through long content
3. **Verify content completeness** - Ensure all information is accessible through scrolling
4. **Test on different terminals** - Verify modals work correctly across various terminal types

#### Terminal Compatibility Testing
1. **Test on Ubuntu Server** - Verify widget borders are visible in terminal-only environments
2. **SSH connection testing** - Confirm borders display correctly when accessing via SSH
3. **Basic terminal testing** - Test with terminals like `linux`, `vt100` to verify ASCII fallbacks
4. **Focus navigation** - Ensure widget focus is clearly visible with fallback borders

---

*This release focuses on improving information accessibility through enhanced modal scrolling and ensuring compatibility across diverse terminal environments. All changes are backward compatible and require no configuration updates.*

**Made with ❤️ by drclcomputers**

*For issues, feature requests, and contributions, visit our [GitHub repository](https://github.com/drclcomputers/syspulse)*