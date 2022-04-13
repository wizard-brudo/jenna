package main

type userInfo struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Id         int    `json:"id,string"`
	ProfilePic string `json:"profile_pic"`
}

type FacebookRequest struct {
	Entry []struct {
		Id        int `json:"id,string"`
		Messaging []struct {
			Message struct {
				Id          string `json:"mid"`
				Text        string `json:"text"`
				Attachments []struct {
					Type    string `json:"type"`
					Payload struct {
						Url string `json:"url"`
					} `json:"payload"`
				} `json:"attachments"`
			} `json:"message"`
			Delivery struct {
				Ids       []string `json:"mids"`
				Watermark int      `json:"watermark"`
			} `json:"delivery"`
			Recipient struct {
				Id int `json:"id,string"`
			} `json:"recipient"`
			Sender struct {
				Id int `json:"id,string"`
			} `json:"sender"`
			Timestamp int `json:"timestamp"`
		} `json:"messaging"`
		Time int `json:"time"`
	}
	Object string `json:"object"`
}

type FacebookResponse struct {
	RecipientId int    `json:"recipient_id,string"`
	MessageId   string `json:"message_id"`
}
type ResponseMessage struct {
	MessagingType string `json:"messaging_type"`
	Recipient     struct {
		Id int `json:"id,string"`
	} `json:"recipient"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}
