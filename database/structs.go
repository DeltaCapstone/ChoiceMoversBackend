package DB

type Customer struct {
	ID           int
	UserName     string
	PasswordHash []byte
	//FristName    string
	//LastName     string
	Email        string
	PhonePrimary string
	PhoneOther   []string
}

type Employee struct {
	ID           int
	UserName     string
	PasswordHash []byte
	//FristName    string
	//LastName     string
	Email        string
	PhonePrimary string
	PhoneOther   []string
	EmployeeType string
}
