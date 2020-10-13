package models

// Keys contains all environment variables.
type Keys struct {
	MONGO_URI             string
	MONGO_DB_NAME         string
	MAILGUN_PRIVATE_KEY   string
	MAILGUN_DOMAIN        string
	MAILGUN_API_BASE      string
	SMTP_SERVER			  string
	SMTP_PORT             int
	SMTP_USER             string
	SMTP_PASSWORD         string
	PORT                  string
	SECRET                string
	CLIENT_URL            string
	MONTHLY_REQUEST_LIMIT int
	IS_SELFHOSTED         bool
}
