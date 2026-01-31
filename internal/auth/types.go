package auth

type ProvisionRequest struct {
	Version        int    `json:"version"`
	BootstrapToken string `json:"bootstrap_token"`
	CSR            string `json:"csr"` 
}
type ProvisionResponse struct {
	Certificate string `json:"certificate"` 
	CACert      string `json:"ca_cert"`     
	ExpiresAt   int64  `json:"expires_at"`  
	Serial      string `json:"serial"`      
}
