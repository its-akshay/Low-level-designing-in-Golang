// ======================================================
// INTERFACE = CONTRACT / ABSTRACTION
// ======================================================

// NotificationSender defines WHAT behavior we need.
// Any struct that has a Send() method with this exact
// signature automatically satisfies this interface.
type NotificationSender interface {

    // Send sends a notification.
    //
    // ctx -> used for cancellation, timeout, request lifecycle
    // notification -> contains notification data
    // returns error if sending fails
    Send(ctx context.Context, notification Notification) error
}

// ======================================================
// SHARED NOTIFICATION DATA
// ======================================================

// Notification represents generic notification data.
// This struct is shared between Email, SMS, Push etc.
type Notification struct {

    // Receiver of notification
    To string

    // Subject/title of notification
    // Mostly useful for Email
    Subject string

    // Actual content/message
    Body string

    // Type of notification
    // Example: Email, SMS, Push
    Type NotificationType
}

// ======================================================
// EMAIL IMPLEMENTATION
// ======================================================

// EmailSender contains all email-specific configuration.
//
// This is the CONCRETE implementation of NotificationSender.
type EmailSender struct {

    // SMTP server host
    smtpHost string

    // SMTP server port
    smtpPort int

    // Email credentials
    username string
    password string

    // SMTP client used to send emails
    client *smtp.Client
}

// Send implements NotificationSender interface for Email.
//
// Since EmailSender has:
// Send(context.Context, Notification) error
//
// it automatically satisfies NotificationSender.
func (e *EmailSender) Send(ctx context.Context, n Notification) error {

    // Create new email message
    msg := mail.NewMessage()

    // Set recipient email
    msg.SetHeader("To", n.To)

    // Set email subject
    msg.SetHeader("Subject", n.Subject)

    // Set email body
    // "text/html" means body supports HTML
    msg.SetBody("text/html", n.Body)

    // Actual SMTP sending logic.
    //
    // OrderService doesn't know anything about SMTP.
    // This complexity is hidden inside EmailSender.
    return e.client.SendMail(n.To, msg)
}

// ======================================================
// SMS IMPLEMENTATION
// ======================================================

// SMSSender contains SMS-specific dependencies/config.
//
// Another concrete implementation of NotificationSender.
type SMSSender struct {

    // Twilio SDK/API client
    twilioClient *twilio.Client

    // Sender phone number
    fromNumber string
}

// Send implements NotificationSender for SMS.
func (s *SMSSender) Send(ctx context.Context, n Notification) error {

    // Send SMS using Twilio API.
    //
    // SMS doesn't care about Subject.
    // Only phone number + message body matter.
    _, err := s.twilioClient.SendSMS(
        s.fromNumber, // sender
        n.To,         // receiver
        n.Body,       // SMS content
    )

    return err
}

// ======================================================
// BUSINESS LOGIC / CONSUMER
// ======================================================

// OrderService handles order-related business logic.
//
// IMPORTANT:
// It depends on ABSTRACTION (NotificationSender)
// instead of concrete implementations like:
//
// ❌ EmailSender
// ❌ SMSSender
//
// This creates loose coupling.
type OrderService struct {

    // Any notifier can be injected here:
    // - EmailSender
    // - SMSSender
    // - PushSender
    // - MockSender (for testing)
    notifier NotificationSender
}

// ConfirmOrder confirms order and sends notification.
func (s *OrderService) ConfirmOrder(
    ctx context.Context,
    order *Order,
) error {

    // ==========================================
    // Business logic
    // ==========================================

    // Example:
    // validate order
    // update DB
    // mark as confirmed
    // etc...

    // ==========================================
    // Send notification
    // ==========================================

    // OrderService DOES NOT KNOW:
    // - how email works
    // - how SMS works
    // - which API is used
    //
    // It only knows:
    // "something can Send()"

    return s.notifier.Send(ctx, Notification{

        // Receiver email/phone
        To: order.CustomerEmail,

        // Notification title
        Subject: "Order Confirmed",

        // Dynamic message body
        Body: fmt.Sprintf(
            "Your order #%s is confirmed!",
            order.ID,
        ),
    })
}