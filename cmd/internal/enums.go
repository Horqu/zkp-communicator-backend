package internal

type MessageCommand string
type ResponseCommand string

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
	ResponseLoginPage        ResponseCommand = "login_page"
	ResponseLoginSuccess     ResponseCommand = "login_success"
	ResponseLoginError       ResponseCommand = "login_error"
	ResponseRegisterPage     ResponseCommand = "register_page"
	ResponseRegisterSuccess  ResponseCommand = "register_success"
	ResponseRegisterError    ResponseCommand = "register_error"
	ResponseSchnorrChallenge ResponseCommand = "schnorr_challenge"
	ResponseFFSChallenge     ResponseCommand = "ffs_challenge"
	ResponseSigmaChallenge   ResponseCommand = "sigma_challenge"
	ResponseSolveSuccess     ResponseCommand = "solve_success"
	ResponseSolveError       ResponseCommand = "solve_error"
	ResponseSelectChat       ResponseCommand = "select_chat"
	ResponseAddFriend        ResponseCommand = "add_friend"
	ResponseRemoveFriend     ResponseCommand = "remove_friend"
	ResponseSendMessage      ResponseCommand = "send_message"
	ResponseRefresh          ResponseCommand = "refresh"
	ResponseLogout           ResponseCommand = "logout"
	ResponseError            ResponseCommand = "error"
)
