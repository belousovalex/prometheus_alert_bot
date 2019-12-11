package core

type Subscriber struct {
	patterns *Patterns
	sender   Sender
}

type AlertManagerReceiver struct {
	subscribers []Subscriber
}

func (r *AlertManagerReceiver) Subscribe(sender Sender, patterns *Patterns) {
	r.subscribers = append(r.subscribers, Subscriber{
		patterns: patterns,
		sender:   sender,
	})
}

func (r *AlertManagerReceiver) Receive(messages []*PrometheusAlert) {
	for _, subscriber := range r.subscribers {
		var needSend []*PrometheusAlert
		for _, message := range messages {
			if r.IsNeedSend(message, subscriber.patterns) {
				needSend = append(needSend, message)
			}
		}
		if len(needSend) > 0 {
			go subscriber.sender.Send(needSend)
		}

	}
}

func (r *AlertManagerReceiver) IsNeedSend(message *PrometheusAlert, patterns *Patterns) bool {
	if patterns == nil {
		return true
	}
	if len(patterns.Labels) > 0 {
		for labelName, reg := range patterns.Labels {
			val, ok := message.Labels[labelName]
			if !ok {
				return false
			}
			if !reg.Match([]byte(val)) {
				return false
			}
		}
	}
	return true
}
