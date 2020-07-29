package models

/*User is a struct containing userinfo:
* userName -- nickname;
* fullName -- name and surname;
* password;
* motto;
 */
type User struct {
	ID       string
	UserName string `json:"user_name"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Motto    string `json:"motto"`
}
