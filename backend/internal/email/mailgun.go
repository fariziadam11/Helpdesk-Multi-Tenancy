package email

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MailgunClient handles email sending via Mailgun API
type MailgunClient struct {
	domain string
	apiKey string
	sender string
	client *http.Client
}

// NewMailgunClient creates a new Mailgun email client
func NewMailgunClient(domain, apiKey, sender string) *MailgunClient {
	return &MailgunClient{
		domain: domain,
		apiKey: apiKey,
		sender: sender,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// EmailRequest represents an email to be sent
type EmailRequest struct {
	To      string
	Subject string
	HTML    string
}

// SendEmail sends an email via Mailgun API
func (m *MailgunClient) SendEmail(req EmailRequest) error {
	// Prepare form data
	data := url.Values{}
	data.Set("from", m.sender)
	data.Set("to", req.To)
	data.Set("subject", req.Subject)
	data.Set("html", req.HTML)

	// Create HTTP request
	apiURL := fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", m.domain)
	httpReq, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set basic auth
	auth := base64.StdEncoding.EncodeToString([]byte("api:" + m.apiKey))
	httpReq.Header.Set("Authorization", "Basic "+auth)

	// Send request
	resp, err := m.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mailgun API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email
func (m *MailgunClient) SendPasswordResetEmail(to, resetLink string) error {
	html := m.getPasswordResetTemplate(resetLink)

	return m.SendEmail(EmailRequest{
		To:      to,
		Subject: "Reset Password Anda - Werk Ticketing",
		HTML:    html,
	})
}

// getPasswordResetTemplate returns the HTML template for password reset email
func (m *MailgunClient) getPasswordResetTemplate(resetLink string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		body {
			font-family: Arial, sans-serif;
			line-height: 1.6;
			color: #333;
			margin: 0;
			padding: 0;
			background-color: #f4f4f4;
		}
		.container {
			max-width: 600px;
			margin: 20px auto;
			background-color: #ffffff;
			border-radius: 8px;
			overflow: hidden;
			box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
		}
		.header {
			background-color: #6929C4;
			color: #ffffff;
			padding: 30px 20px;
			text-align: center;
		}
		.header h1 {
			margin: 0;
			font-size: 24px;
		}
		.content {
			padding: 30px 20px;
		}
		.content p {
			margin: 0 0 15px 0;
		}
		.button {
			display: inline-block;
			background-color: #6929C4;
			color: #ffffff !important;
			padding: 12px 30px;
			text-decoration: none;
			border-radius: 4px;
			margin: 20px 0;
			font-weight: 500;
		}
		.button:hover {
			background-color: #8A3FFC;
		}
		.footer {
			background-color: #f4f4f4;
			padding: 20px;
			text-align: center;
			font-size: 12px;
			color: #666;
		}
		.warning {
			background-color: #fff3cd;
			border-left: 4px solid #ffc107;
			padding: 12px;
			margin: 20px 0;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🔐 Reset Password</h1>
		</div>
		<div class="content">
			<p>Halo,</p>
			<p>Kami menerima permintaan untuk mereset password akun Werk Ticketing Anda.</p>
			<p>Klik tombol di bawah ini untuk membuat password baru:</p>
			<p style="text-align: center;">
				<a href="%s" class="button">Reset Password</a>
			</p>
			<div class="warning">
				<strong>⚠️ Penting:</strong>
				<ul style="margin: 5px 0; padding-left: 20px;">
					<li>Link ini akan kadaluarsa dalam <strong>1 jam</strong></li>
					<li>Link hanya dapat digunakan <strong>satu kali</strong></li>
				</ul>
			</div>
			<p>Jika Anda tidak meminta reset password, abaikan email ini. Password Anda tidak akan berubah.</p>
			<p style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #e0e0e0;">
				Salam,<br>
				<strong>Tim Werk Ticketing</strong>
			</p>
		</div>
		<div class="footer">
			<p>Email ini dikirim secara otomatis. Mohon tidak membalas email ini.</p>
			<p>&copy; 2026 Werk Ticketing. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`, resetLink)
}
