package quotify

// UsersResp is response item from slack
type UsersResp struct {
	Members []MemberResp `json:"members"`
}

// MemberResp is a single member response object
type MemberResp struct {
	ID   string `json:"id"`
	Name string `json:"real_name"`
}

// ChannelResp is channel response from slack
type ChannelResp struct {
	Messages []MessageResp `json:"messages"`
}

// MessageResp is a single slack message response
type MessageResp struct {
	Text string `json:"text"`
}

// Quote is a single quote with speaker and Speech
type Quote struct {
	Speaker string `json:"speaker"`
	Speech  string `json:"speech"`
}
