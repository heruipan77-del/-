package devicemanager

type Config struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	DeviceID string `json:"deviceId"`
	Username string `json:"username"`
	Password string `json:"password"`
	Insecure bool   `json:"insecure"`
}

type StoragePool struct {
	ID                 string `json:"ID"`
	Name               string `json:"NAME"`
	HealthStatus       string `json:"HEALTHSTATUS"`
	UserTotalCapacity  string `json:"USERTOTALCAPACITY"`
	UserFreeCapacity   string `json:"USERFREECAPACITY"`
	DataSpace          string `json:"DATASPACE"`
	CompressedCapacity string `json:"COMPRESSEDCAPACITY"`
}

type LUN struct {
	ID            string `json:"ID"`
	Name          string `json:"NAME"`
	Capacity      string `json:"CAPACITY"`
	HealthStatus  string `json:"HEALTHSTATUS"`
	RunningStatus string `json:"RUNNINGSTATUS"`
}
