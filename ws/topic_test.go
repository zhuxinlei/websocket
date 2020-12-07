package ws

import (
	"server-tokenhouse-ws/config"
	"testing"
)

func TestTopic_IsValid(t *testing.T) {
	config.LoadConfig("../config.yaml")
	InitValidTopicTrie()

	tests := []struct {
		name string
		t    Topic
		want bool
	}{
		//{
		//	name:"",
		//	t: Topic(""),
		//	want: false,
		//},
		//{
		//	name:"ctc.siesusd.kline",
		//	t: Topic("ctc.siesusd.kline"),
		//	want: false,
		//},
		//{
		//	name:"ctc.siesusd.kline.2",
		//	t: Topic("ctc.siesusd.kline.2"),
		//	want: false,
		//},
		//{
		//	name:"ctc.siesusd.order.1",
		//	t: Topic("ctc.siesusd.order.1"),
		//	want: false,
		//},
		{
			name: "ctc.siesusd.order",
			t:    Topic("ctc.siesusd.order"),
			want: true,
		},
		{
			name: "ctc.siesusd.trade",
			t:    Topic("ctc.siesusd.trade"),
			want: true,
		},
		{
			name: "ctc.siesusd.kline.1",
			t:    Topic("ctc.siesusd.kline.1"),
			want: true,
		},
		{
			name: "ctc.siesusd.kline.10080",
			t:    Topic("ctc.siesusd.kline.10080"),
			want: true,
		},
		{
			name: "ctc.siesusd.kline.10080.1",
			t:    Topic("ctc.siesusd.kline.10080.1"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
