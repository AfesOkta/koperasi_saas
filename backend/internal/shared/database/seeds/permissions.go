package seeds

import (
	"context"
	"log"

	iamModel "github.com/koperasi-gresik/backend/internal/modules/iam/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SystemPermissions defines the full set of permissions for the platform.
// Format: resource:action or resource:action:scope
var SystemPermissions = []iamModel.Permission{
	// IAM & Organization
	{Name: "user:read", Resource: "user", Action: "read", Scope: "any", Description: "View users in organization"},
	{Name: "user:create", Resource: "user", Action: "create", Scope: "any", Description: "Create new users"},
	{Name: "user:update", Resource: "user", Action: "update", Scope: "any", Description: "Edit user profiles"},
	{Name: "user:delete", Resource: "user", Action: "delete", Scope: "any", Description: "Remove users"},
	{Name: "role:read", Resource: "role", Action: "read", Scope: "any", Description: "View roles"},
	{Name: "role:manage", Resource: "role", Action: "manage", Scope: "any", Description: "Create/edit/delete roles and assign permissions"},
	{Name: "org:read", Resource: "org", Action: "read", Scope: "any", Description: "View organization info"},
	{Name: "org:update", Resource: "org", Action: "update", Scope: "any", Description: "Edit organization settings"},

	// Member
	{Name: "member:read", Resource: "member", Action: "read", Scope: "any", Description: "View member list and profiles"},
	{Name: "member:create", Resource: "member", Action: "create", Scope: "any", Description: "Register new members"},
	{Name: "member:update", Resource: "member", Action: "update", Scope: "any", Description: "Edit member data"},
	{Name: "member:delete", Resource: "member", Action: "delete", Scope: "any", Description: "Deactivate members"},
	{Name: "member:approve", Resource: "member", Action: "approve", Scope: "any", Description: "Approve member applications"},

	// Savings
	{Name: "savings:read", Resource: "savings", Action: "read", Scope: "any", Description: "View savings accounts"},
	{Name: "savings:deposit", Resource: "savings", Action: "deposit", Scope: "any", Description: "Record a deposit"},
	{Name: "savings:withdraw", Resource: "savings", Action: "withdraw", Scope: "any", Description: "Record a withdrawal"},
	{Name: "savings:close", Resource: "savings", Action: "close", Scope: "any", Description: "Close a savings account"},

	// Loans
	{Name: "loan:read", Resource: "loan", Action: "read", Scope: "any", Description: "View loan applications and schedules"},
	{Name: "loan:apply", Resource: "loan", Action: "apply", Scope: "any", Description: "Submit a loan application"},
	{Name: "loan:approve:any", Resource: "loan", Action: "approve", Scope: "any", Description: "Approve or reject any loan"},
	{Name: "loan:approve:own", Resource: "loan", Action: "approve", Scope: "own", Description: "Approve loans submitted by self"},
	{Name: "loan:disburse", Resource: "loan", Action: "disburse", Scope: "any", Description: "Disburse an approved loan"},
	{Name: "loan:payment", Resource: "loan", Action: "payment", Scope: "any", Description: "Record an installment payment"},
	{Name: "loan:waive", Resource: "loan", Action: "waive", Scope: "any", Description: "Waive penalties"},

	// Cash
	{Name: "cash:read", Resource: "cash", Action: "read", Scope: "any", Description: "View cash accounts and transactions"},
	{Name: "cash:in", Resource: "cash", Action: "in", Scope: "any", Description: "Record a cash-in transaction"},
	{Name: "cash:out", Resource: "cash", Action: "out", Scope: "any", Description: "Record a cash-out transaction"},

	// Accounting
	{Name: "accounting:read", Resource: "accounting", Action: "read", Scope: "any", Description: "View journals, ledger, reports"},
	{Name: "accounting:post", Resource: "accounting", Action: "post", Scope: "any", Description: "Create manual journal entries"},
	{Name: "accounting:void", Resource: "accounting", Action: "void", Scope: "any", Description: "Void or reverse a journal entry"},
	{Name: "accounting:close", Resource: "accounting", Action: "close", Scope: "any", Description: "Run EOD/EOM closing"},

	// Inventory
	{Name: "inventory:read", Resource: "inventory", Action: "read", Scope: "any", Description: "View products and stock"},
	{Name: "inventory:manage", Resource: "inventory", Action: "manage", Scope: "any", Description: "Create/update products, adjust stock"},

	// POS
	{Name: "pos:transact", Resource: "pos", Action: "transact", Scope: "any", Description: "Perform POS sales"},

	// Purchasing
	{Name: "purchasing:read", Resource: "purchasing", Action: "read", Scope: "any", Description: "View purchase orders"},
	{Name: "purchasing:manage", Resource: "purchasing", Action: "manage", Scope: "any", Description: "Create/approve purchase orders"},

	// Sales
	{Name: "sales:read", Resource: "sales", Action: "read", Scope: "any", Description: "View sales history"},

	// SHU
	{Name: "shu:read", Resource: "shu", Action: "read", Scope: "any", Description: "View SHU distributions"},
	{Name: "shu:calculate", Resource: "shu", Action: "calculate", Scope: "any", Description: "Run SHU calculation"},
	{Name: "shu:distribute", Resource: "shu", Action: "distribute", Scope: "any", Description: "Mark SHU as distributed"},

	// Reports
	{Name: "report:read", Resource: "report", Action: "read", Scope: "any", Description: "View all reports"},
	{Name: "report:export", Resource: "report", Action: "export", Scope: "any", Description: "Export reports (PDF/Excel)"},
}

// SeedPermissions upserts all system permissions into the database.
func SeedPermissions(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"resource", "action", "scope", "description"}),
		}).
		Create(&SystemPermissions).Error; err != nil {
		return err
	}
	log.Printf("✅ Seeded %d system permissions", len(SystemPermissions))
	return nil
}

// rolePermissionMatrix maps system role names to their permission names.
var rolePermissionMatrix = map[string][]string{
	"admin": {
		"user:read", "user:create", "user:update", "user:delete",
		"role:read", "role:manage",
		"org:read", "org:update",
		"member:read", "member:create", "member:update", "member:delete", "member:approve",
		"savings:read", "savings:deposit", "savings:withdraw", "savings:close",
		"loan:read", "loan:apply", "loan:approve:any", "loan:approve:own", "loan:disburse", "loan:payment", "loan:waive",
		"cash:read", "cash:in", "cash:out",
		"accounting:read", "accounting:post", "accounting:void", "accounting:close",
		"inventory:read", "inventory:manage",
		"pos:transact",
		"purchasing:read", "purchasing:manage",
		"sales:read",
		"shu:read", "shu:calculate", "shu:distribute",
		"report:read", "report:export",
	},
	"manager": {
		"user:read",
		"role:read",
		"org:read",
		"member:read", "member:create", "member:update", "member:approve",
		"savings:read", "savings:deposit", "savings:withdraw", "savings:close",
		"loan:read", "loan:apply", "loan:approve:any", "loan:approve:own", "loan:disburse", "loan:payment", "loan:waive",
		"cash:read", "cash:in", "cash:out",
		"accounting:read", "accounting:post",
		"inventory:read", "inventory:manage",
		"pos:transact",
		"purchasing:read", "purchasing:manage",
		"sales:read",
		"shu:read", "shu:calculate",
		"report:read", "report:export",
	},
	"teller": {
		"member:read",
		"savings:read", "savings:deposit", "savings:withdraw",
		"loan:read", "loan:apply", "loan:approve:own", "loan:payment",
		"cash:read", "cash:in", "cash:out",
		"inventory:read",
		"pos:transact",
		"purchasing:read",
		"sales:read",
	},
	"staff": {
		"member:read",
		"savings:read",
		"loan:read",
		"cash:read",
		"accounting:read",
		"inventory:read",
		"purchasing:read",
		"sales:read",
		"shu:read",
		"report:read",
	},
	"auditor": {
		"savings:read",
		"loan:read",
		"cash:read",
		"accounting:read",
		"shu:read",
		"report:read", "report:export",
	},
}

// SeedSystemRoles creates the 5 system roles for a given organization and assigns permissions.
// Safe to call on every org creation — uses upsert logic.
func SeedSystemRoles(ctx context.Context, db *gorm.DB, orgID uint) error {
	// Build permission name → ID map
	var allPerms []iamModel.Permission
	if err := db.WithContext(ctx).Find(&allPerms).Error; err != nil {
		return err
	}
	permMap := make(map[string]uint, len(allPerms))
	for _, p := range allPerms {
		permMap[p.Name] = p.ID
	}

	for roleName, permNames := range rolePermissionMatrix {
		// Upsert the role
		role := iamModel.Role{}
		result := db.WithContext(ctx).
			Where(iamModel.Role{Name: roleName}).
			Where("organization_id = ?", orgID).
			FirstOrCreate(&role, iamModel.Role{
				Name:        roleName,
				Description: "System role: " + roleName,
				IsSystem:    true,
				Version:     1,
			})
		if result.Error != nil {
			return result.Error
		}
		role.OrganizationID = orgID

		// Build permissions slice
		var perms []iamModel.Permission
		for _, pName := range permNames {
			if id, ok := permMap[pName]; ok {
				perms = append(perms, iamModel.Permission{ID: id})
			}
		}

		// Replace associations
		if err := db.WithContext(ctx).Model(&role).Association("Permissions").Replace(perms); err != nil {
			return err
		}
	}

	log.Printf("✅ Seeded 5 system roles for org %d", orgID)
	return nil
}
