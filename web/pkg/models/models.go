package models

import (
	"database/sql"
	"errors"
	"time"
)

// ErrType errors
var (
	ErrNoRecord = errors.New("models: no matching record found")
	// User tries to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// User tries to signup with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

//
// Fields
//

// GroupJSON JSON/JSONB fields
type GroupJSON struct {
	ID    int                    `db:"id"`
	Value map[string]interface{} `db:"value"`
}

// VarJSON JSON/JSONB fields
type VarJSON struct {
	ID    int                    `db:"id"`
	Value map[string]interface{} `db:"value"`
}

// Properties JSON/JSONB fields
type Properties map[string]interface{}

//
// Model
//

// User model
type User struct {
	ID             int
	UserName       string
	UserRole       string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// InventoryEnvironment model
type InventoryEnvironment struct {
	ID    int    `db:"id"`
	Value string `db:"value"`
}

// InventoryHost model
type InventoryHost struct {
	ID                int            `db:"id"`
	Hostname          string         `db:"hostname"`
	IP                string         `db:"ip"`
	Environment       string         `db:"environment"`
	MethodCheckConsul string         `db:"method_check_consul"`
	Approved          bool           `db:"approved"`
	Delete            bool           `db:"delete"`
	Status            int            `db:"status"`
	Description       sql.NullString `db:"description"`
	Properties        *Properties    `db:"properties"`
	Created           time.Time
	Updated           time.Time
	CountServices     int
}

// InventoryHVar model
type InventoryHVar struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Value string `db:"value"`
}

// InventoryHTag model
type InventoryHTag struct {
	ID    int    `db:"id"`
	Value string `db:"value"`
}

// InventoryGroup model
type InventoryGroup struct {
	ID          int    `db:"id"`
	Environment string `db:"environment"`
	ParentID    int    `db:"parent_id"`
	Value       string `db:"value"`
	Created     time.Time
	Updated     time.Time
}

// InventoryGVar model
type InventoryGVar struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Value string `db:"value"`
}

// InventoryService model (environments)
type InventoryService struct {
	ID                int `db:"id"`
	Host              InventoryHost
	Type              string         `db:"type"`
	Location          string         `db:"location"`
	Title             string         `db:"title"`
	Link1             string         `db:"link1"`
	TechName          string         `db:"techname"`
	Link2             string         `db:"link2"`
	Domain            string         `db:"domain"`
	Placement         string         `db:"placement"`
	PublicTime        string         `db:"publictime"`
	Team              string         `db:"team"`
	Resp              string         `db:"resp"`
	Value             string         `db:"value"`
	Port              int            `db:"port"`
	Approved          bool           `db:"approved"`
	Delete            bool           `db:"delete"`
	RegToConsul       bool           `db:"reg_to_consul"`
	StatusInConsul    bool           `db:"status_in_consul"`
	MethodCheckConsul string         `db:"method_check_consul"`
	Description       sql.NullString `db:"description"`
	Properties        *Properties    `db:"properties"`
	Created           time.Time
	Updated           time.Time
}

// ServerAgent model
type ServerAgent struct {
	ID             int    `db:"id"`
	Hostname       string `db:"hostname"`
	IP             string `db:"ip"`
	OSName         string `db:"os_name"`
	OSVendor       string `db:"os_vendor"`
	OSVersion      string `db:"os_version"`
	OSRelease      string `db:"os_release"`
	OSArchitecture string `db:"os_architecture"`
	KernelRelease  string `db:"kernel_release"`
	ModelCPU       string `db:"cpu_model"`
	NumCPU         string `db:"cpu_num"`
	RAM            string `db:"ram"`
	HDD            string `db:"hdd"`
	Uptime         string `db:"uptime"`
	ResolvConf     string `db:"resolv"`
	StorageSize    string `db:"storage_size"`
	TimeZone       string `db:"timezone"`
	MTU            int    `db:"mtu"`
	Delete         int    `db:"delete"`
	StatusHard     int    `db:"status_hard"`
	StatusResolv   int    `db:"status_resolv"`
	Created        time.Time
	Updated        time.Time
}

// PackageServer model
type PackageServer struct {
	Package       string `db:"package"`
	Version       string `db:"version"`
	StatusPackage string `db:"status_package"`
	Created       time.Time
	Updated       time.Time
}

// SysCtlServer model
type SysCtlServer struct {
	Name         string `db:"name"`
	Value        string `db:"value"`
	StatusSysCtl string `db:"status_sysctl"`
	Created      time.Time
	Updated      time.Time
}

// BaseTemplate model
type BaseTemplate struct {
	ID         int    `db:"id"`
	Type       string `db:"type"`
	Value      string `db:"value"`
	Parameters int
	Created    time.Time
	Updated    time.Time
}

// BaseTemplateItem model
type BaseTemplateItem struct {
	ID      int    `db:"id"`
	IDTpl   string `db:"idTpl"`
	Title   string `db:"title"`
	Value   string `db:"value"`
	Created time.Time
	Updated time.Time
}

// History model
type History struct {
	ID          int    `db:"id"`
	UserID      string `db:"user_id"`
	UserEmail   string
	EntityID    string `db:"entity_id"`
	Entity      string `db:"entity"`
	Action      string `db:"action"`
	Event       string `db:"event"`
	Description string `db:"description"`
	Created     time.Time
	Updated     time.Time
}
