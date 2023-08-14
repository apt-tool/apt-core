package enum

// Role of user controls the access of each user
type Role int

const (
	RoleAdmin Role = iota + 1
	RoleUser
	RoleViewer
)

func ConvertNumberToRole(input int) Role {
	switch input {
	case 1:
		return RoleAdmin
	case 2:
		return RoleUser
	case 3:
		return RoleViewer
	}

	return 0
}
