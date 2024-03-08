package define

type ServerMode string

const (
	ServerModeDebug   ServerMode = "debug"
	ServerModeTest    ServerMode = "test"
	ServerModeRelease ServerMode = "release"
)

type ClientType string

const (
	ClientTypeWeapp ClientType = "weapp"
)

type Gender int

const (
	GenderMale   Gender = 1
	GenderFemale Gender = 2
)

func (t Gender) Label() string {
	switch t {
	case GenderMale:
		return "男"
	case GenderFemale:
		return "女"
	default:
		return "UNKNOWN"
	}
}

type StorageClassType string

const (
	// 标准存储（Standard）
	StorageStandard StorageClassType = "standard"

	// 低频访问（Infrequent Access）
	StorageIA StorageClassType = "ia"

	// 归档存储（Archive）
	StorageArchive StorageClassType = "archive"

	// 冷归档存储（Cold Archive）
	StorageColdArchive StorageClassType = "cold_archive"
)
