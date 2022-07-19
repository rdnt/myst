package remote

type Remote struct {
	Address    string
	Username   string
	Password   string
	PublicKey  []byte
	PrivateKey []byte
}
