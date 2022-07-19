package remote

type Repository interface {
	SetRemote(address, username, password string, publicKey, privateKey []byte) error
	Remote() (Remote, error)
}
