package scanner

import "time"

// Canvas is the complete Developer Knowledge Canvas — a comprehensive,
// multi-lens representation of a codebase for rapid AI understanding.
type Canvas struct {
	SchemaVersion string           `json:"schema_version"`
	GeneratedAt   time.Time        `json:"generated_at"`
	Project       ProjectOverview  `json:"project"`
	Files         []FileCard       `json:"files"`
	Architecture  ArchitectureView `json:"architecture"`
	Dependencies  DependencyGraph  `json:"dependencies"`
	API           APIView          `json:"api"`
	Database      DatabaseView     `json:"database"`
	Security      SecurityView     `json:"security"`
	Git           GitSummary       `json:"git"`
	Notes         []CanvasNote     `json:"notes"`
	Groups        []CanvasGroup    `json:"groups"`
	Connections   []CanvasEdge     `json:"connections"`
	AI            AIContext        `json:"ai_context"`
	Lenses        LensesView       `json:"lenses"`
}

// ── Project ──────────────────────────────────────────────────────────

type ProjectOverview struct {
	Name            string   `json:"name"`
	RootPath        string   `json:"root_path"`
	Description     string   `json:"description"`
	Languages       []string `json:"languages"`
	Frameworks      []string `json:"frameworks"`
	PackageManagers []string `json:"package_managers"`
	Databases       []string `json:"databases"`
	Infra           []string `json:"infra"`
	TotalFiles      int      `json:"total_files"`
	TotalLOC        int      `json:"total_loc"`
	Monorepo        bool     `json:"monorepo"`
	Subprojects     int      `json:"subprojects"`
}

// ── File Cards ───────────────────────────────────────────────────────

type FileCard struct {
	Path         string     `json:"path"`
	Language     string     `json:"language"`
	LOC          int        `json:"loc"`
	Imports      []string   `json:"imports"`
	Exports      []string   `json:"exports"`
	Functions    []Function `json:"functions"`
	FileType     string     `json:"file_type"` // component, util, config, route, test, etc.
	LastModified string     `json:"last_modified"`
	Dependencies []string   `json:"dependencies"` // files this file imports
	Dependents   []string   `json:"dependents"`   // files that import this file
	Tests        []string   `json:"tests"`
	HasErrors    bool       `json:"has_errors"`
	ErrorCount   int        `json:"error_count"`
	Coverage     float64    `json:"coverage"` // 0-100, -1 if unknown
	HasTodos     bool       `json:"has_todos"`
	TODOList     []string   `json:"todo_list"`
	IsEntryPoint bool       `json:"is_entry_point"`
	IsDeadCode   bool       `json:"is_dead_code"`
	SecurityIssues []SecurityFinding `json:"security_issues"`
}

type Function struct {
	Name       string `json:"name"`
	Line       int    `json:"line"`
	IsExported bool   `json:"is_exported"`
	IsTest     bool   `json:"is_test"`
	Calls      []string `json:"calls"`
}

// ── Architecture ─────────────────────────────────────────────────────

type ArchitectureView struct {
	Layers         []ArchLayer   `json:"layers"`
	DataFlow       []DataFlow    `json:"data_flow"`
	EntryPoints    []EntryPoint  `json:"entry_points"`
	DirectoryTree  []DirNode     `json:"directory_tree"`
}

type ArchLayer struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Paths       []string `json:"paths"`
	Technologies []string `json:"technologies"`
}

type DataFlow struct {
	From string `json:"from"`
	To   string `json:"to"`
	Via  string `json:"via,omitempty"`
	Type string `json:"type"` // request, event, data, etc.
}

type EntryPoint struct {
	Path     string `json:"path"`
	Type     string `json:"type"` // cli, web, worker, test
	Command  string `json:"command,omitempty"`
	Port     int    `json:"port,omitempty"`
}

type DirNode struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	IsDir    bool      `json:"is_dir"`
	Children []DirNode `json:"children,omitempty"`
}

// ── Dependencies ─────────────────────────────────────────────────────

type DependencyGraph struct {
	Files          []DepNode       `json:"files"`
	Edges          []DepEdge       `json:"edges"`
	CircularGroups [][]string      `json:"circular_groups"`
	DeadFiles      []string        `json:"dead_files"`
	ImpactMap      map[string][]string `json:"impact_map"`
}

type DepNode struct {
	Path    string   `json:"path"`
	Imports int      `json:"imports"`
	Exports int      `json:"exports"`
	Depth   int      `json:"depth"`
	Layer   string   `json:"layer"`
}

type DepEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"` // import, re-export, dynamic
}

// ── API ──────────────────────────────────────────────────────────────

type APIView struct {
	Routes      []APIRoute    `json:"routes"`
	Controllers []string      `json:"controllers"`
	Middlewares []string      `json:"middlewares"`
	BasePath    string        `json:"base_path"`
}

type APIRoute struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Handler  string `json:"handler"`
	File     string `json:"file"`
	Auth     bool   `json:"auth"`
	RateLimited bool `json:"rate_limited"`
}

// ── Database ─────────────────────────────────────────────────────────

type DatabaseView struct {
	Models    []DBModel    `json:"models"`
	Relations []DBRelation `json:"relations"`
	Migrations []string    `json:"migrations"`
	ORM       string       `json:"orm,omitempty"`
}

type DBModel struct {
	Name       string            `json:"name"`
	File       string            `json:"file"`
	Fields     []DBField         `json:"fields"`
	TableName  string            `json:"table_name,omitempty"`
}

type DBField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsPK     bool   `json:"is_primary_key"`
	IsFK     bool   `json:"is_foreign_key"`
	Ref      string `json:"references,omitempty"`
	Nullable bool   `json:"nullable"`
}

type DBRelation struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Type   string `json:"type"` // has_one, has_many, belongs_to, many_to_many
	Via    string `json:"via,omitempty"`
}

// ── Security ─────────────────────────────────────────────────────────

type SecurityView struct {
	Findings    []SecurityFinding `json:"findings"`
	Score       int               `json:"score"`
	CriticalCnt int               `json:"critical_count"`
	HighCnt     int               `json:"high_count"`
}

type SecurityFinding struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Severity string `json:"severity"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	Fix      string `json:"fix,omitempty"`
}

// ── Git ──────────────────────────────────────────────────────────────

type GitSummary struct {
	TotalCommits   int              `json:"total_commits"`
	Branches       []string         `json:"branches"`
	Contributors   []Contributor    `json:"contributors"`
	LastCommit     string           `json:"last_commit"`
	HotFiles       []HotFile        `json:"hot_files"`
	ChurnMap       map[string]int   `json:"churn_map"`
}

type Contributor struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Count int    `json:"commit_count"`
}

type HotFile struct {
	Path   string `json:"path"`
	Changes int   `json:"changes"`
}

// ── Notes / Groups / Connections (User-editable Canvas Layer) ────────

type CanvasNote struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Color     string `json:"color"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	CreatedAt string `json:"created_at"`
	Tags      []string `json:"tags,omitempty"`
}

type CanvasGroup struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Members []string `json:"members"`
	Pinned  bool     `json:"pinned"`
}

type CanvasEdge struct {
	ID     string `json:"id"`
	From   string `json:"from"`
	To     string `json:"to"`
	Label  string `json:"label,omitempty"`
	Type   string `json:"type"` // auto, manual
	Dashed bool   `json:"dashed,omitempty"`
}

// ── AI Context (for rapid project understanding) ─────────────────────

type AIContext struct {
	Summary                  string             `json:"summary"`
	OneLiner                 string             `json:"one_liner"`
	Architecture             string             `json:"architecture"`
	KeyInsights              []string           `json:"key_insights"`
	Patterns                 []string           `json:"patterns"`
	Conventions              []string           `json:"conventions"`
	Gotchas                  []string           `json:"gotchas"`
	RecommendedReading       []string           `json:"recommended_reading"`
	EntryPoints              []EntryPoint       `json:"entry_points"`
	CriticalFiles            []string           `json:"critical_files"`
	FrequentChanges          []string           `json:"frequent_changes"`
	TechStack                []string           `json:"tech_stack"`
	ArchitectureDescription  string             `json:"architecture_description"`
	DangerZones              []string           `json:"danger_zones"`
	FAQs                     []FAQ              `json:"faqs"`
	Glossary                 []GlossaryTerm     `json:"glossary"`
}

type FAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type GlossaryTerm struct {
	Term        string `json:"term"`
	Definition  string `json:"definition"`
}

// ── Lenses ───────────────────────────────────────────────────────────

type LensesView struct {
	Technology  *TechnologyLens  `json:"technology,omitempty"`
	File        *FileLens        `json:"file,omitempty"`
	Function    *FunctionLens    `json:"function,omitempty"`
	API         *APILens         `json:"api,omitempty"`
	Database    *DatabaseLens    `json:"database,omitempty"`
	Security    *SecurityLens    `json:"security,omitempty"`
	Git         *GitLens         `json:"git,omitempty"`
	Runtime     *RuntimeLens     `json:"runtime,omitempty"`
	Architecture *ArchitectureLens `json:"architecture,omitempty"`
}

type TechnologyLens struct {
	Categories map[string][]string `json:"categories"`
}
type FileLens struct {
	ByLayer   map[string][]string `json:"by_layer"`
	ByType    map[string][]string `json:"by_type"`
}
type FunctionLens struct {
	CallGraph map[string][]string `json:"call_graph"`
}
type APILens struct {
	RoutesByGroup map[string][]APIRoute `json:"routes_by_group"`
}
type DatabaseLens struct {
	Tables []string `json:"tables"`
	Edges  []DBRelation `json:"edges"`
}
type SecurityLens struct {
	SeverityMap map[string][]SecurityFinding `json:"severity_map"`
}
type GitLens struct {
	ContributorMap map[string]int `json:"contributor_map"`
	ChurnMap       map[string]int `json:"churn_map"`
}
type RuntimeLens struct {
	Services []RuntimeService `json:"services"`
}
type RuntimeService struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Port    int    `json:"port,omitempty"`
	PID     int    `json:"pid,omitempty"`
}
type ArchitectureLens struct {
	AIArchitecture []ArchLayer `json:"ai_architecture"`
	FlowDiagram    string      `json:"flow_diagram,omitempty"`
}

// ── Constants ────────────────────────────────────────────────────────

const CanvasSchemaVersion = "1.0.0"
const CanvasSavePath = ".autodevs/canvas/canvas.json"
