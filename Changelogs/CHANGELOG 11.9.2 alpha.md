# SysPulse Changelog

## Version 11.9.2 Alpha - Bug Fixes Release
*Release Date: July 19, 2025*

### 🐛 Bug Fixes

#### Process Filter Modal Improvements
- **Fixed widget navigation interference** - Widget navigation keys (C, M, D, N, P, G) no longer interfere when typing in search input fields (filtering processes)
- **Fixed Tab navigation in modals** - Tab key now properly navigates between form fields in process filter modal instead of triggering global widget cycling
- **Fixed filter state persistence** - Process filter modal now correctly remembers and displays the previously selected filter type and search term when reopened
- **Improved input field focus detection** - Enhanced logic to detect when user is typing in input fields to prevent unwanted global shortcut activation

#### Visual Consistency Improvements
- **Enhanced memory bar styling** - Memory bars now use consistent styling with colored usage portion and white empty portion using `░` character for better visual consistency
- **Simplified memory widget design** - Removed unnecessary color complexity from memory bars in favor of simple two-color design matching disk widget style
- **Updated disk bar display** - Disk bars now show next to the used/total GB information for better space utilization and visual consistency

#### Input Handling Enhancements
- **Smart input field detection** - Global input handler now properly detects when focus is on `tview.InputField` or `tview.Form` elements
- **Conditional shortcut processing** - Widget navigation shortcuts are now only processed when not actively typing in input fields
- **Preserved form functionality** - Tab navigation within forms works correctly while maintaining global Tab cycling for main interface
- **Eliminated redundant nil checks** - Removed unnecessary nil checks in type assertion patterns where the assertion itself provides nil safety

### 🔧 Technical Implementation

#### Input Capture Logic Enhancement
```go
// Enhanced input handler with smart field detection
currentFocused := d.App.GetFocus()
_, isInputField := currentFocused.(*tview.InputField)
_, isForm := currentFocused.(*tview.Form)

// Allow normal keyboard input when in an input field
shouldProcessGlobalKeys := isMainWidgetActive && !isInputField && !isForm
```

#### Modal State Management
- **Filter state preservation** - Process filter term and type are now properly stored and restored
- **Dropdown index calculation** - Correct mapping between filter types and dropdown indices
- **Search input restoration** - Previous search terms are populated when reopening filter modal
- **Event timing optimization** - SetChangedFunc now called after form setup to prevent interference

#### Widget Visual Improvements
- **Memory bar simplification** - Updated `getMemoryBar()` function to use simplified color scheme
- **Consistent progress bar styling** - Memory and disk widgets now share unified visual design language
- **Improved bar positioning** - Disk bars repositioned to appear inline with usage statistics

### 📋 Affected Components

#### Modal System
- `internal/services/UI/modals.go` - Enhanced process search modal with state persistence and improved input handling
- `internal/services/UI/layout.go` - Improved global input handling logic with smart field detection

#### Widget Systems
- `internal/services/memory/vmem.go` - Updated memory bar styling for consistency
- `internal/services/disk/usage&partitions.go` - Enhanced disk bar positioning and display
- `internal/utils/models.go` - Simplified memory model structure

### 🎯 User Experience Improvements

#### Filter Modal Usability
- **Intuitive typing experience** - Users can now type normally in search fields without unexpected widget switching
- **Persistent filter state** - No need to re-enter search terms or re-select filter types when reopening modal
- **Proper form navigation** - Tab key works as expected within forms for field-to-field navigation
- **Seamless modal interaction** - Input events are properly routed to prevent conflicts

#### Visual Enhancement
- **Unified progress bar styling** - Memory and disk widgets now share consistent visual design language
- **Clean color scheme** - Simplified color usage in progress bars for better readability
- **Professional appearance** - Consistent styling across all measurement widgets
- **Better space utilization** - Optimized layout for better information density

#### Keyboard Navigation
- **Context-aware shortcuts** - Global shortcuts only activate when appropriate, preserving typing experience
- **Consistent behavior** - Tab navigation works predictably in both main interface and modal contexts
- **Reduced user frustration** - Eliminated unexpected widget switching while typing

### 🔄 Migration Notes

#### Automatic Updates
- **No configuration changes required** - All improvements are automatically applied
- **Backward compatibility maintained** - Existing themes and configurations continue to work
- **No breaking changes** - All existing functionality preserved

#### Behavior Changes
- **Improved modal interaction** - Users will notice smoother typing experience in search fields
- **Enhanced form navigation** - Tab key behavior is now consistent with standard form conventions
- **Better state persistence** - Filter settings are remembered between modal opens
- **Consistent widget styling** - Memory and disk widgets now have unified appearance

### 🧪 Testing Recommendations

#### Filter Modal Testing
1. **Open process filter modal** (Press 'P' then 'F')
2. **Type search terms** - Verify letters like 'c', 'm', 'd' don't trigger widget navigation
3. **Use Tab key** - Confirm it navigates between form fields
4. **Apply filter and reopen** - Verify search term and filter type are restored
5. **Test global shortcuts** - Confirm they work normally in main interface

#### Widget Visual Testing
1. **Check memory bar styling** - Verify bars show colored usage with white empty portion
2. **Compare with disk bars** - Confirm consistent visual styling
3. **Test different usage levels** - Verify bars display correctly across usage ranges
4. **Check bar positioning** - Verify disk bars appear correctly next to usage statistics

### 🔧 Developer Notes

#### Input Event Flow
```
User Input → Global Handler → Field Detection → Route to Appropriate Handler
                          ↓
                   [Main Widget] → Process Global Shortcuts
                   [Form/Input] → Pass Through for Normal Typing
```

#### State Management Pattern
```go
// Filter state persistence pattern
if d.ProcessFilterActive {
    searchInput.SetText(d.ProcessFilterTerm)
    filterState.filterType = d.ProcessFilterType
}
```

#### Type Assertion Best Practice
```go
// Simplified type assertion without redundant nil check
_, isInputField := currentFocused.(*tview.InputField)
// Type assertion provides nil safety automatically
```

#### Widget Styling Consistency
```go
// Unified memory bar implementation
usedBar := strings.Repeat(utils.BAR, usedWidth)
emptyBar := strings.Repeat("░", barWidth-usedWidth)
return fmt.Sprintf("[%s]%s[-][%s]%s[-]", barColor, usedBar, 
    utils.GetColorFromName(d.Theme.Foreground), emptyBar)
```

---

### Summary

Version 11.9.2 Alpha focuses on fixing critical usability issues with the process filter modal and improving visual consistency across widgets. The main improvements include proper input handling that prevents global shortcuts from interfering with typing, persistent filter state that remembers user preferences, and unified widget styling for a more professional appearance. These changes significantly enhance the user experience while maintaining full backward compatibility.
- `internal/services/UI/modals.go` - Enhanced process search modal with state persistence
- `internal/services/UI/layout.go` - Improved global input handling logic

#### Memory Widget
- `internal/services/memory/vmem.go` - Updated memory bar styling to match disk widget pattern
- `internal/utils/models.go` - Simplified memory model structure

### 🎯 User Experience Improvements

#### Filter Modal Usability
- **Intuitive typing experience** - Users can now type normally in search fields without unexpected widget switching
- **Persistent filter state** - No need to re-enter search terms or re-select filter types when reopening modal
- **Proper form navigation** - Tab key works as expected within forms for field-to-field navigation

#### Keyboard Navigation
- **Context-aware shortcuts** - Global shortcuts only activate when appropriate, preserving typing experience
- **Consistent behavior** - Tab navigation works predictably in both main interface and modal contexts
- **Reduced user frustration** - Eliminated unexpected widget switching while typing

#### Visual Consistency
- **Unified progress bar styling** - Memory and disk widgets now share consistent visual design language
- **Clean color scheme** - Simplified color usage in progress bars for better readability
- **Professional appearance** - Consistent styling across all measurement widgets

### 🔄 Migration Notes

#### Automatic Updates
- **No configuration changes required** - All improvements are automatically applied
- **Backward compatibility maintained** - Existing themes and configurations continue to work
- **No breaking changes** - All existing functionality preserved

#### Behavior Changes
- **Improved modal interaction** - Users will notice smoother typing experience in search fields
- **Enhanced form navigation** - Tab key behavior is now consistent with standard form conventions
- **Better state persistence** - Filter settings are remembered between modal opens

### 🧪 Testing Recommendations

#### Filter Modal Testing
1. **Open process filter modal** (Press 'P' then 'F')
2. **Type search terms** - Verify letters like 'c', 'm', 'd' don't trigger widget navigation
3. **Use Tab key** - Confirm it navigates between form fields
4. **Apply filter and reopen** - Verify search term and filter type are restored
5. **Test global shortcuts** - Confirm they work normally in main interface

#### Memory Widget Testing
1. **Check memory bar styling** - Verify bars show colored usage with white empty portion
2. **Compare with disk bars** - Confirm consistent visual styling
3. **Test different memory usage levels** - Verify bars display correctly across usage ranges

### 🔧 Developer Notes

#### Input Event Flow
```
User Input → Global Handler → Field Detection → Route to Appropriate Handler
                          ↓
                   [Input Field] → Pass Through
                   [Main Widget] → Process Global Shortcuts
```

#### State Management Pattern
```go
// Filter state persistence pattern
if d.ProcessFilterActive {
    searchInput.SetText(d.ProcessFilterTerm)
    filterState.filterType = d.ProcessFilterType
}
```

#### Type Assertion Best Practice
```go
// Simplified type assertion without redundant nil check
_, isInputField := currentFocused.(*tview.InputField)
// Type assertion provides nil safety automatically
```

---

*This release focuses on improving user experience through better input handling and interface consistency. All changes are backward compatible and require no configuration updates.*
