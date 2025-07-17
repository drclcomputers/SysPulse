# SysPulse Themes

This directory contains various pre-configured themes for SysPulse with different color schemes, layouts, and widget configurations. All themes are configured with `updatetime=1` for fast updates.

## Available Themes

### 1. **Cyberpunk** (`cyberpunk.json`)
- **Background**: Black
- **Primary Colors**: Cyan and Magenta
- **Style**: Futuristic neon aesthetic
- **Layout**: 4x4 grid with most widgets enabled
- **Features**: High contrast cyan/magenta color scheme, process tree disabled for cleaner look

### 2. **Matrix** (`matrix.json`)
- **Background**: Black
- **Primary Colors**: Green and Lime
- **Style**: Classic "Matrix" movie aesthetic
- **Layout**: 3x4 grid, simplified layout
- **Features**: Minimal widget set (CPU, Memory, Disk, Network, GPU, Process), clean green theme

### 3. **Ocean** (`ocean.json`)
- **Background**: Navy
- **Primary Colors**: Blue and Aqua
- **Style**: Deep ocean/underwater theme
- **Layout**: 4x3 grid with spacing
- **Features**: Soothing blue tones, white text on navy background

### 4. **Sunset** (`sunset.json`)
- **Background**: Maroon
- **Primary Colors**: Yellow and Orange
- **Style**: Warm sunset colors
- **Layout**: 3x5 grid, wide layout
- **Features**: All widgets enabled except battery, warm color palette

### 5. **Monochrome** (`monochrome.json`)
- **Background**: Black
- **Primary Colors**: White, Silver, and Gray
- **Style**: Clean black & white aesthetic
- **Layout**: 2x6 grid, horizontal layout
- **Features**: Professional look, minimal colors, wide process widget

### 6. **Neon** (`neon.json`)
- **Background**: Black
- **Primary Colors**: Purple and Pink
- **Style**: Neon nightclub aesthetic
- **Layout**: 6x2 grid, vertical layout
- **Features**: All widgets enabled in vertical arrangement, vibrant purple/pink theme

### 7. **Forest** (`forest.json`)
- **Background**: Olive
- **Primary Colors**: Yellow and Lime
- **Style**: Nature/forest theme
- **Layout**: 4x4 grid with spacing
- **Features**: Earth tones, spaced layout for better readability

### 8. **Fire** (`fire.json`)
- **Background**: Black
- **Primary Colors**: Red and Orange
- **Style**: Intense fire/heat theme
- **Layout**: 3x3 grid, compact
- **Features**: Aggressive red/orange colors, minimal widget set for performance focus

## How to Use Themes

### Method 1: Copy to config.json
```bash
# Copy your preferred theme to the main config
cp themes/cyberpunk.json config.json
```

### Method 2: Backup and Replace
```bash
# Backup current config
cp config.json config.backup.json

# Use new theme
cp themes/matrix.json config.json
```

### Method 3: Manual Configuration
Edit `config.json` and copy the contents from your preferred theme file.

## Theme Customization

Each theme file contains these main sections:

- **Global Colors**: `background`, `foreground`, `altforeground`
- **Component Colors**: `cpu`, `memory`, `network`, `disk`, `gpu` bar colors
- **Layout Configuration**: Widget positions, sizes, and enabled states
- **Widget Colors**: Individual `border_color` and `foreground_color` for each widget
- **Settings**: `updatetime`, `processsort`, `export` configuration

## Creating Custom Themes

To create your own theme:

1. Copy an existing theme file
2. Modify the colors and layout to your preference
3. Available colors: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`, `orange`, `purple`, `pink`, `lime`, `teal`, `aqua`, `navy`, `gray`, `silver`, `maroon`, `olive`
4. Adjust widget positions using `row`, `column`, `rowSpan`, `colSpan`
5. Enable/disable widgets with the `enabled` field

## Performance Notes

- All themes use `updatetime=1` for maximum responsiveness
- Themes with fewer enabled widgets will use less system resources
- Larger grid layouts (more columns/rows) may require larger terminal windows
- Some widgets like `process_tree` and `disk_io` are disabled in certain themes for better performance

## Color Scheme Examples

### High Contrast Themes
- **Cyberpunk**: Cyan/Magenta on Black
- **Matrix**: Green/Lime on Black
- **Fire**: Red/Orange on Black

### Professional Themes
- **Ocean**: Blue/Aqua on Navy
- **Monochrome**: White/Silver/Gray on Black

### Warm Themes
- **Sunset**: Yellow/Orange on Maroon
- **Forest**: Yellow/Lime on Olive

### Vibrant Themes
- **Neon**: Purple/Pink on Black

Choose the theme that best fits your terminal environment and personal preferences!
