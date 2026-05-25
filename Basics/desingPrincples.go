package basics

import "context"

// 1. DRY — Don't Repeat Yourself
// Intuition
// Every piece of *knowledge* should have a *single authoritative representation* in the system.
// It's not just about copy-pasting code — it's about not having the same *decision* encoded in multiple places.
// Like putting a a validation into one simple common function and call it in handlers

// ## 2. KISS — Keep It Simple, Stupid
// ### Intuition
// The simplest solution that correctly solves the problem is usually best. Complexity is a liability — it's harder to test, debug, and modify.

// ## 3. YAGNI — You Aren't Gonna Need It
// ### Intuition
// Don't build features you *might* need. Build what's needed *now*. Premature generalization is as bad as premature optimization.

// badcode:
// ❌ BAD: Building multi-tenant support before it's required
type TenantAwareUserService struct {
	repo        UserRepository
	tenantRepo  TenantRepository
	tenantCtx   TenantContext
	tenantCache TenantCache
	// ... 5 more tenant-related dependencies
}

func (s *TenantAwareUserService) GetUser(ctx context.Context, tenantID, userID string) (*User, error) {
	tenant, _ := s.tenantRepo.Find(ctx, tenantID)
	// complex tenant isolation logic that isn't needed yet
}



// goodcode
// ✅ GOOD: Build for today's needs; design for extension
type UserService struct {
	repo UserRepository
}

// Simple for now. When multi-tenancy is needed, add tenantID to context
// or extend the repository interface. Don't over-engineer upfront.
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repo.FindByID(ctx, id)
}
