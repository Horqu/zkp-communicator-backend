package internal

type AppView int
type MessageCommand string
type ResponseCommand string

const (
	ViewMain AppView = iota
	ViewRegister
	ViewLogin
	ViewResolver
	ViewLogged
	ViewLoading
	ViewError
)
const (
	MessageLoginButtom    MessageCommand = "login_buttom"
	MessageLogin          MessageCommand = "login"
	MessageRegisterButtom MessageCommand = "register_buttom"
	MessageRegister       MessageCommand = "register"
	MessageSolve          MessageCommand = "solve"
	MessageSelectChat     MessageCommand = "select_chat"
	MessageAddFriend      MessageCommand = "add_friend"
	MessageRemoveFriend   MessageCommand = "remove_friend"
	MessageSendMessage    MessageCommand = "send_message"
	MessageRefresh        MessageCommand = "refresh"
	MessageLogout         MessageCommand = "logout"
)

const (
	ResponseLoginPage       ResponseCommand = "login_page"
	ResponseLoginSuccess    ResponseCommand = "login_success"
	ResponseLoginError      ResponseCommand = "login_error"
	ResponseRegisterPage    ResponseCommand = "register_page"
	ResponseRegisterSuccess ResponseCommand = "register_success"
	ResponseRegisterError   ResponseCommand = "register_error"
	ResponseSolveSuccess    ResponseCommand = "solve_success"
	ResponseSolveError      ResponseCommand = "solve_error"
	ResponseSelectChat      ResponseCommand = "select_chat"
	ResponseAddFriend       ResponseCommand = "add_friend"
	ResponseRemoveFriend    ResponseCommand = "remove_friend"
	ResponseSendMessage     ResponseCommand = "send_message"
	ResponseRefresh         ResponseCommand = "refresh"
	ResponseLogout          ResponseCommand = "logout"
	ResponseError           ResponseCommand = "error"
)
