package secrets

type Secret struct {
	Username            string `json:"username"`
	Password            string `json:"password"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	Name                string `json:"dbname"`
	Engine              string `json:"engine"`
	DBClusterIdentifier string `json:"dbInstanceIdentifier"`
}

type AzureAPISecret struct {
	SubscriptionKey string `json:"subscriptionKey"`
	Endpoint        string `json:"endpoint"`
}

type S3CredentialSecret struct {
	S3AccessKey string `json:"s3AccessKey"`
	S3SecretKey string `json:"s3SecretKey"`
}

type AzureSFTPSecret struct {
	AwsSftpHost          string `json:"AwsSftpHost"`
	AwsSftpHostPort      uint   `json:"AwsSftpHostPort,string"`
	AwsSftpHostAccount   string `json:"AwsSftpHostAccount"`
	AwsSftpHostKey       string `json:"AwsSftpHostKey"`
	AwsSftpHostKeyPhrase string `json:"AwsSftpHostKeyPhrase"`
}
