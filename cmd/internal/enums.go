package internal

type MessageCommand string
type ResponseCommand string

const (
	CommandLoginButtom    MessageCommand = "login_buttom"
	CommandLogin          MessageCommand = "login"
	CommandRegisterButtom MessageCommand = "register_buttom"
	CommandRegister       MessageCommand = "register"
	CommandSolve          MessageCommand = "solve"
	CommandSelectChat     MessageCommand = "select_chat"
	CommandAddFriend      MessageCommand = "add_friend"
	CommandRemoveFriend   MessageCommand = "remove_friend"
	CommandSendMessage    MessageCommand = "send_message"
	CommandRefresh        MessageCommand = "refresh"
	CommandLogout         MessageCommand = "logout"
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
