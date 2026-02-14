package ws

func SubAllMids() SubRequest {
	return SubRequest{
		Method: "subscribe",
		Subscription: map[string]string{
			"type": "allMids",
		},
	}
}

func SubUserFills(user string) SubRequest {
	return SubRequest{
		Method: "subscribe",
		Subscription: map[string]string{
			"type": "userFills",
			"user": user,
		},
	}
}

func SubUserFundings(user string) SubRequest {
	return SubRequest{
		Method: "subscribe",
		Subscription: map[string]string{
			"type": "userFundings",
			"user": user,
		},
	}
}

func SubOrderUpdates(user string) SubRequest {
	return SubRequest{
		Method: "subscribe",
		Subscription: map[string]string{
			"type": "orderUpdates",
			"user": user,
		},
	}
}
