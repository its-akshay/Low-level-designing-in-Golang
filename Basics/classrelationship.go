package basics

import (
	"context"
	"fmt"
	"time"
)

// Association
// A uses B as part of its behavior. B exists independently.

type UserHandler struct {
	service UserService // reference, not ownership
	logger  Logger
}

// Aggregation
// A contains B, but B can exist without A.
// Player exists independently of Team
type Player struct {
	ID   string
	Name string
	Role string
}

// Team AGGREGATES Players — players can exist without the team
type Team struct {
	ID      string
	Name    string
	Players []*Player // aggregation: pointers to independent entities
}

// 3. Composition
// A owns B — B cannot meaningfully exist without A.
// OrderItem has no meaning outside an Order
type OrderItem struct {
	ProductID string
	Quantity  int
	UnitPrice float64
	// No separate ID — it's part of Order
}

func (i OrderItem) Subtotal() float64 {
	return float64(i.Quantity) * i.UnitPrice
}

// Order COMPOSES OrderItems — they live and die together
type Order struct {
	ID         string
	CustomerID string
	Items      []OrderItem // value type — owned, not shared
	Status     OrderStatus
	CreatedAt  time.Time
}

// DEPENDENCY
// A temporarily uses B — B is passed in as a parameter.
// EmailService DEPENDS ON Template — passed as parameter
type EmailService struct {
	smtpClient SMTPClient
}

// Template is a dependency — not stored in struct but just used in this call
func (s *EmailService) SendWelcome(ctx context.Context, user *User, tmpl *Template) error {
	body, err := tmpl.Render(map[string]any{"User": user})
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}
	return s.smtpClient.Send(ctx, user.Email(), "Welcome!", body)
}
