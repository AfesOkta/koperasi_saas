package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const permCacheTTL = 10 * time.Minute

// PermissionCache handles Redis-backed permission caching per role version.
// Cache key format: perm:{roleID}:{roleVersion}
// Cached value: JSON array of permission name strings.
type PermissionCache struct {
	rdb *redis.Client
	db  *gorm.DB
}

// NewPermissionCache creates a new PermissionCache.
func NewPermissionCache(rdb *redis.Client, db *gorm.DB) *PermissionCache {
	return &PermissionCache{rdb: rdb, db: db}
}

// cacheKey returns the Redis key for a given role + version.
func (pc *PermissionCache) cacheKey(roleID uint, roleVersion int) string {
	return fmt.Sprintf("perm:%d:%d", roleID, roleVersion)
}

// GetPermissions returns the full permission set for a role, loading from DB on cache miss.
func (pc *PermissionCache) GetPermissions(ctx context.Context, roleID uint, roleVersion int) ([]string, error) {
	key := pc.cacheKey(roleID, roleVersion)

	// Try Redis first
	val, err := pc.rdb.Get(ctx, key).Result()
	if err == nil {
		var perms []string
		if jsonErr := json.Unmarshal([]byte(val), &perms); jsonErr == nil {
			return perms, nil
		}
	}

	// Cache miss — load from DB
	perms, err := pc.loadFromDB(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Store in Redis (fire-and-forget, don't block request)
	if jsonBytes, jsonErr := json.Marshal(perms); jsonErr == nil {
		if setErr := pc.rdb.Set(ctx, key, jsonBytes, permCacheTTL).Err(); setErr != nil {
			log.Printf("⚠️  Failed to cache permissions for role %d: %v", roleID, setErr)
		}
	}

	return perms, nil
}

// loadFromDB fetches all permission names for a role from the database.
func (pc *PermissionCache) loadFromDB(ctx context.Context, roleID uint) ([]string, error) {
	var names []string
	err := pc.db.WithContext(ctx).
		Table("permissions").
		Select("permissions.name").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Pluck("permissions.name", &names).Error
	return names, err
}

// HasPermission checks if a given permission name is in the cached set.
func (pc *PermissionCache) HasPermission(ctx context.Context, roleID uint, roleVersion int, permission string) bool {
	perms, err := pc.GetPermissions(ctx, roleID, roleVersion)
	if err != nil {
		log.Printf("⚠️  Permission check error for role %d: %v", roleID, err)
		return false
	}
	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}

// InvalidateRole removes cached permissions for a specific role version.
// Note: With version-based keys, old versions self-expire — this is for immediate invalidation.
func (pc *PermissionCache) InvalidateRole(ctx context.Context, roleID uint, roleVersion int) {
	key := pc.cacheKey(roleID, roleVersion)
	if err := pc.rdb.Del(ctx, key).Err(); err != nil {
		log.Printf("⚠️  Failed to invalidate cache for role %d v%d: %v", roleID, roleVersion, err)
	}
}
