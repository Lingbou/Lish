package theme

// BuiltinThemes 存储所有内置主题
var BuiltinThemes = map[string]ColorScheme{
	"dark":            &DarkTheme{},
	"light":           &LightTheme{},
	"solarized-dark":  &SolarizedDarkTheme{},
	"solarized-light": &SolarizedLightTheme{},
	"gruvbox":         &GruvboxDarkTheme{},
	"dracula":         &DraculaTheme{},
	"nord":            &NordTheme{},
	"monokai":         &MonokaiProTheme{},
}

// =============================================================================
// 1. Dark Theme (默认)
// =============================================================================

type DarkTheme struct{}

func (t *DarkTheme) Primary() Color      { return NewColor("#00FF00", "bold") }
func (t *DarkTheme) Secondary() Color    { return NewColor("#AAAAAA") }
func (t *DarkTheme) Success() Color      { return NewColor("#00FF00") }
func (t *DarkTheme) Warning() Color      { return NewColor("#FFFF00") }
func (t *DarkTheme) Error() Color        { return NewColor("#FF5555") }
func (t *DarkTheme) Info() Color         { return NewColor("#5555FF") }
func (t *DarkTheme) Directory() Color    { return NewColor("#5555FF", "bold") }
func (t *DarkTheme) Executable() Color   { return NewColor("#00FF00", "bold") }
func (t *DarkTheme) Symlink() Color      { return NewColor("#00FFFF") }
func (t *DarkTheme) Archive() Color      { return NewColor("#FF00FF") }
func (t *DarkTheme) PromptUser() Color   { return NewColor("#00FF00", "bold") }
func (t *DarkTheme) PromptHost() Color   { return NewColor("#00FFFF") }
func (t *DarkTheme) PromptPath() Color   { return NewColor("#5555FF", "bold") }
func (t *DarkTheme) PromptGit() Color    { return NewColor("#FFFF00") }
func (t *DarkTheme) SyntaxCommand() Color  { return NewColor("#00FF00", "bold") }
func (t *DarkTheme) SyntaxArgument() Color { return NewColor("#AAAAAA") }
func (t *DarkTheme) SyntaxString() Color   { return NewColor("#FFFF00") }
func (t *DarkTheme) SyntaxVariable() Color { return NewColor("#00FFFF") }
func (t *DarkTheme) SyntaxOperator() Color { return NewColor("#FF00FF") }

// =============================================================================
// 2. Light Theme
// =============================================================================

type LightTheme struct{}

func (t *LightTheme) Primary() Color      { return NewColor("#008700", "bold") }
func (t *LightTheme) Secondary() Color    { return NewColor("#444444") }
func (t *LightTheme) Success() Color      { return NewColor("#008700") }
func (t *LightTheme) Warning() Color      { return NewColor("#AF5F00") }
func (t *LightTheme) Error() Color        { return NewColor("#D70000") }
func (t *LightTheme) Info() Color         { return NewColor("#0087FF") }
func (t *LightTheme) Directory() Color    { return NewColor("#0087FF", "bold") }
func (t *LightTheme) Executable() Color   { return NewColor("#008700", "bold") }
func (t *LightTheme) Symlink() Color      { return NewColor("#00AFAF") }
func (t *LightTheme) Archive() Color      { return NewColor("#AF00AF") }
func (t *LightTheme) PromptUser() Color   { return NewColor("#008700", "bold") }
func (t *LightTheme) PromptHost() Color   { return NewColor("#00AFAF") }
func (t *LightTheme) PromptPath() Color   { return NewColor("#0087FF", "bold") }
func (t *LightTheme) PromptGit() Color    { return NewColor("#AF5F00") }
func (t *LightTheme) SyntaxCommand() Color  { return NewColor("#008700", "bold") }
func (t *LightTheme) SyntaxArgument() Color { return NewColor("#444444") }
func (t *LightTheme) SyntaxString() Color   { return NewColor("#AF5F00") }
func (t *LightTheme) SyntaxVariable() Color { return NewColor("#00AFAF") }
func (t *LightTheme) SyntaxOperator() Color { return NewColor("#AF00AF") }

// =============================================================================
// 3. Solarized Dark
// =============================================================================

type SolarizedDarkTheme struct{}

func (t *SolarizedDarkTheme) Primary() Color      { return NewColor("#859900", "bold") }
func (t *SolarizedDarkTheme) Secondary() Color    { return NewColor("#93a1a1") }
func (t *SolarizedDarkTheme) Success() Color      { return NewColor("#859900") }
func (t *SolarizedDarkTheme) Warning() Color      { return NewColor("#b58900") }
func (t *SolarizedDarkTheme) Error() Color        { return NewColor("#dc322f") }
func (t *SolarizedDarkTheme) Info() Color         { return NewColor("#268bd2") }
func (t *SolarizedDarkTheme) Directory() Color    { return NewColor("#268bd2", "bold") }
func (t *SolarizedDarkTheme) Executable() Color   { return NewColor("#859900", "bold") }
func (t *SolarizedDarkTheme) Symlink() Color      { return NewColor("#2aa198") }
func (t *SolarizedDarkTheme) Archive() Color      { return NewColor("#d33682") }
func (t *SolarizedDarkTheme) PromptUser() Color   { return NewColor("#859900", "bold") }
func (t *SolarizedDarkTheme) PromptHost() Color   { return NewColor("#2aa198") }
func (t *SolarizedDarkTheme) PromptPath() Color   { return NewColor("#268bd2", "bold") }
func (t *SolarizedDarkTheme) PromptGit() Color    { return NewColor("#b58900") }
func (t *SolarizedDarkTheme) SyntaxCommand() Color  { return NewColor("#859900", "bold") }
func (t *SolarizedDarkTheme) SyntaxArgument() Color { return NewColor("#93a1a1") }
func (t *SolarizedDarkTheme) SyntaxString() Color   { return NewColor("#b58900") }
func (t *SolarizedDarkTheme) SyntaxVariable() Color { return NewColor("#2aa198") }
func (t *SolarizedDarkTheme) SyntaxOperator() Color { return NewColor("#d33682") }

// =============================================================================
// 4. Solarized Light
// =============================================================================

type SolarizedLightTheme struct{}

func (t *SolarizedLightTheme) Primary() Color      { return NewColor("#859900", "bold") }
func (t *SolarizedLightTheme) Secondary() Color    { return NewColor("#657b83") }
func (t *SolarizedLightTheme) Success() Color      { return NewColor("#859900") }
func (t *SolarizedLightTheme) Warning() Color      { return NewColor("#b58900") }
func (t *SolarizedLightTheme) Error() Color        { return NewColor("#dc322f") }
func (t *SolarizedLightTheme) Info() Color         { return NewColor("#268bd2") }
func (t *SolarizedLightTheme) Directory() Color    { return NewColor("#268bd2", "bold") }
func (t *SolarizedLightTheme) Executable() Color   { return NewColor("#859900", "bold") }
func (t *SolarizedLightTheme) Symlink() Color      { return NewColor("#2aa198") }
func (t *SolarizedLightTheme) Archive() Color      { return NewColor("#d33682") }
func (t *SolarizedLightTheme) PromptUser() Color   { return NewColor("#859900", "bold") }
func (t *SolarizedLightTheme) PromptHost() Color   { return NewColor("#2aa198") }
func (t *SolarizedLightTheme) PromptPath() Color   { return NewColor("#268bd2", "bold") }
func (t *SolarizedLightTheme) PromptGit() Color    { return NewColor("#b58900") }
func (t *SolarizedLightTheme) SyntaxCommand() Color  { return NewColor("#859900", "bold") }
func (t *SolarizedLightTheme) SyntaxArgument() Color { return NewColor("#657b83") }
func (t *SolarizedLightTheme) SyntaxString() Color   { return NewColor("#b58900") }
func (t *SolarizedLightTheme) SyntaxVariable() Color { return NewColor("#2aa198") }
func (t *SolarizedLightTheme) SyntaxOperator() Color { return NewColor("#d33682") }

// =============================================================================
// 5. Gruvbox Dark
// =============================================================================

type GruvboxDarkTheme struct{}

func (t *GruvboxDarkTheme) Primary() Color      { return NewColor("#b8bb26", "bold") }
func (t *GruvboxDarkTheme) Secondary() Color    { return NewColor("#ebdbb2") }
func (t *GruvboxDarkTheme) Success() Color      { return NewColor("#b8bb26") }
func (t *GruvboxDarkTheme) Warning() Color      { return NewColor("#fabd2f") }
func (t *GruvboxDarkTheme) Error() Color        { return NewColor("#fb4934") }
func (t *GruvboxDarkTheme) Info() Color         { return NewColor("#83a598") }
func (t *GruvboxDarkTheme) Directory() Color    { return NewColor("#83a598", "bold") }
func (t *GruvboxDarkTheme) Executable() Color   { return NewColor("#b8bb26", "bold") }
func (t *GruvboxDarkTheme) Symlink() Color      { return NewColor("#8ec07c") }
func (t *GruvboxDarkTheme) Archive() Color      { return NewColor("#d3869b") }
func (t *GruvboxDarkTheme) PromptUser() Color   { return NewColor("#b8bb26", "bold") }
func (t *GruvboxDarkTheme) PromptHost() Color   { return NewColor("#8ec07c") }
func (t *GruvboxDarkTheme) PromptPath() Color   { return NewColor("#83a598", "bold") }
func (t *GruvboxDarkTheme) PromptGit() Color    { return NewColor("#fabd2f") }
func (t *GruvboxDarkTheme) SyntaxCommand() Color  { return NewColor("#b8bb26", "bold") }
func (t *GruvboxDarkTheme) SyntaxArgument() Color { return NewColor("#ebdbb2") }
func (t *GruvboxDarkTheme) SyntaxString() Color   { return NewColor("#fabd2f") }
func (t *GruvboxDarkTheme) SyntaxVariable() Color { return NewColor("#8ec07c") }
func (t *GruvboxDarkTheme) SyntaxOperator() Color { return NewColor("#d3869b") }

// =============================================================================
// 6. Dracula
// =============================================================================

type DraculaTheme struct{}

func (t *DraculaTheme) Primary() Color      { return NewColor("#50fa7b", "bold") }
func (t *DraculaTheme) Secondary() Color    { return NewColor("#f8f8f2") }
func (t *DraculaTheme) Success() Color      { return NewColor("#50fa7b") }
func (t *DraculaTheme) Warning() Color      { return NewColor("#f1fa8c") }
func (t *DraculaTheme) Error() Color        { return NewColor("#ff5555") }
func (t *DraculaTheme) Info() Color         { return NewColor("#8be9fd") }
func (t *DraculaTheme) Directory() Color    { return NewColor("#bd93f9", "bold") }
func (t *DraculaTheme) Executable() Color   { return NewColor("#50fa7b", "bold") }
func (t *DraculaTheme) Symlink() Color      { return NewColor("#8be9fd") }
func (t *DraculaTheme) Archive() Color      { return NewColor("#ff79c6") }
func (t *DraculaTheme) PromptUser() Color   { return NewColor("#50fa7b", "bold") }
func (t *DraculaTheme) PromptHost() Color   { return NewColor("#8be9fd") }
func (t *DraculaTheme) PromptPath() Color   { return NewColor("#bd93f9", "bold") }
func (t *DraculaTheme) PromptGit() Color    { return NewColor("#f1fa8c") }
func (t *DraculaTheme) SyntaxCommand() Color  { return NewColor("#50fa7b", "bold") }
func (t *DraculaTheme) SyntaxArgument() Color { return NewColor("#f8f8f2") }
func (t *DraculaTheme) SyntaxString() Color   { return NewColor("#f1fa8c") }
func (t *DraculaTheme) SyntaxVariable() Color { return NewColor("#8be9fd") }
func (t *DraculaTheme) SyntaxOperator() Color { return NewColor("#ff79c6") }

// =============================================================================
// 7. Nord
// =============================================================================

type NordTheme struct{}

func (t *NordTheme) Primary() Color      { return NewColor("#a3be8c", "bold") }
func (t *NordTheme) Secondary() Color    { return NewColor("#d8dee9") }
func (t *NordTheme) Success() Color      { return NewColor("#a3be8c") }
func (t *NordTheme) Warning() Color      { return NewColor("#ebcb8b") }
func (t *NordTheme) Error() Color        { return NewColor("#bf616a") }
func (t *NordTheme) Info() Color         { return NewColor("#81a1c1") }
func (t *NordTheme) Directory() Color    { return NewColor("#81a1c1", "bold") }
func (t *NordTheme) Executable() Color   { return NewColor("#a3be8c", "bold") }
func (t *NordTheme) Symlink() Color      { return NewColor("#88c0d0") }
func (t *NordTheme) Archive() Color      { return NewColor("#b48ead") }
func (t *NordTheme) PromptUser() Color   { return NewColor("#a3be8c", "bold") }
func (t *NordTheme) PromptHost() Color   { return NewColor("#88c0d0") }
func (t *NordTheme) PromptPath() Color   { return NewColor("#81a1c1", "bold") }
func (t *NordTheme) PromptGit() Color    { return NewColor("#ebcb8b") }
func (t *NordTheme) SyntaxCommand() Color  { return NewColor("#a3be8c", "bold") }
func (t *NordTheme) SyntaxArgument() Color { return NewColor("#d8dee9") }
func (t *NordTheme) SyntaxString() Color   { return NewColor("#ebcb8b") }
func (t *NordTheme) SyntaxVariable() Color { return NewColor("#88c0d0") }
func (t *NordTheme) SyntaxOperator() Color { return NewColor("#b48ead") }

// =============================================================================
// 8. Monokai Pro
// =============================================================================

type MonokaiProTheme struct{}

func (t *MonokaiProTheme) Primary() Color      { return NewColor("#a9dc76", "bold") }
func (t *MonokaiProTheme) Secondary() Color    { return NewColor("#fcfcfa") }
func (t *MonokaiProTheme) Success() Color      { return NewColor("#a9dc76") }
func (t *MonokaiProTheme) Warning() Color      { return NewColor("#ffd866") }
func (t *MonokaiProTheme) Error() Color        { return NewColor("#ff6188") }
func (t *MonokaiProTheme) Info() Color         { return NewColor("#78dce8") }
func (t *MonokaiProTheme) Directory() Color    { return NewColor("#ab9df2", "bold") }
func (t *MonokaiProTheme) Executable() Color   { return NewColor("#a9dc76", "bold") }
func (t *MonokaiProTheme) Symlink() Color      { return NewColor("#78dce8") }
func (t *MonokaiProTheme) Archive() Color      { return NewColor("#ff6188") }
func (t *MonokaiProTheme) PromptUser() Color   { return NewColor("#a9dc76", "bold") }
func (t *MonokaiProTheme) PromptHost() Color   { return NewColor("#78dce8") }
func (t *MonokaiProTheme) PromptPath() Color   { return NewColor("#ab9df2", "bold") }
func (t *MonokaiProTheme) PromptGit() Color    { return NewColor("#ffd866") }
func (t *MonokaiProTheme) SyntaxCommand() Color  { return NewColor("#a9dc76", "bold") }
func (t *MonokaiProTheme) SyntaxArgument() Color { return NewColor("#fcfcfa") }
func (t *MonokaiProTheme) SyntaxString() Color   { return NewColor("#ffd866") }
func (t *MonokaiProTheme) SyntaxVariable() Color { return NewColor("#78dce8") }
func (t *MonokaiProTheme) SyntaxOperator() Color { return NewColor("#ff6188") }

