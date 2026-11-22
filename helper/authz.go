package helper

type PermissionSet struct {
	rules map[string]map[string]bool
	raw   map[string][]string
}

func NewPermissionSet(raw map[string][]string) *PermissionSet {
	rules := make(map[string]map[string]bool, len(raw))
	for role, perms := range raw {
		rules[role] = make(map[string]bool, len(perms))
		for _, p := range perms {
			rules[role][p] = true
		}
	}
	return &PermissionSet{
		rules: rules,
		raw:   raw,
	}
}

func (p *PermissionSet) Can(role, permission string) bool {
	if p == nil || p.rules == nil {
		return false
	}
	perms, exists := p.rules[role]
	if !exists {
		return false
	}
	return perms[permission] || perms["*"]
}

func (p *PermissionSet) PermissionsOf(role string) []string {
	if p == nil || p.raw == nil {
		return nil
	}
	return p.raw[role]
}

func HasAnyRole(userRole string, allowed ...string) bool {
	for _, r := range allowed {
		if userRole == r {
			return true
		}
	}
	return false
}

func DefaultClinicPermissions() map[string][]string {
	return map[string][]string{
		"pasien": {
			"profile:read_self", "profile:update_self",
			"clinic:read",
			"appointment:create", "appointment:read_self", "appointment:cancel_self",
			"medical_record:read",
		},
		"dokter": {
			"profile:read_self", "profile:update_self",
			"clinic:read",
			"queue:read_doctor",
			"medical_record:create", "medical_record:add_prescription", "medical_record:read",
		},
		"admin": {
			"profile:read_self", "profile:update_self",
			"clinic:read", "clinic:manage",
			"appointment:read_all", "appointment:update_status",
		},
	}
}
