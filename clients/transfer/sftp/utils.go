package sftp

import "os"

// CreateFlags returns the commonly specified flags when uploading file to remote path.
func CreateFlags() int {
	return os.O_RDWR | os.O_CREATE | os.O_TRUNC
}

// CreateAwsTransferFlags returns the commonly specified flags when uploading file to AWS transfer SFTP.
func CreateAwsTransferFlags() int {
	return os.O_WRONLY | os.O_CREATE | os.O_TRUNC
}
