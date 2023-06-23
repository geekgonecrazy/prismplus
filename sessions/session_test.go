package sessions

import (
	"testing"
	"time"

	"github.com/geekgonecrazy/rtmp-lib/av"
)

func TestSession_RelayPacket(t *testing.T) {
	type args struct {
		p av.Packet
	}
	tests := []struct {
		name                  string
		session               Session
		args                  args
		wantedPacketTime      time.Duration
		wantedBufferDuration  time.Duration
		wantedDiscrepancySize time.Duration
	}{
		{
			name: "Receive a packet",
			session: Session{
				bufferedDuration: time.Duration(0),
			},
			args: args{
				p: av.Packet{
					Time: time.Duration(500 * time.Millisecond),
				},
			},
			wantedPacketTime:      time.Duration(500 * time.Millisecond),
			wantedBufferDuration:  time.Duration(500 * time.Millisecond),
			wantedDiscrepancySize: time.Duration(0),
		},

		{
			name: "1st Reconnect 1/3 Packet - 500ms timing",
			session: Session{
				previousIncomingDuration: time.Duration(5 * time.Minute),
				bufferedDuration:         time.Duration(5 * time.Minute),
			},
			args: args{
				p: av.Packet{
					Time: time.Duration(500 * time.Millisecond),
				},
			},
			wantedPacketTime:      time.Duration(5*time.Minute) + time.Duration(500*time.Millisecond),
			wantedBufferDuration:  time.Duration(5*time.Minute) + time.Duration(500*time.Millisecond),
			wantedDiscrepancySize: time.Duration(5 * time.Minute),
		},

		{
			name: "1st Reconnect 2/3 Packet - 600ms timing",
			session: Session{
				previousIncomingDuration: time.Duration(5*time.Minute) + time.Duration(500*time.Millisecond),
				bufferedDuration:         time.Duration(5*time.Minute) + time.Duration(500*time.Millisecond),
				discrepancySize:          time.Duration(5 * time.Minute),
			},
			args: args{
				p: av.Packet{
					Time: time.Duration(600 * time.Millisecond),
				},
			},
			wantedPacketTime:      time.Duration(5*time.Minute) + time.Duration(600*time.Millisecond),
			wantedBufferDuration:  time.Duration(5*time.Minute) + time.Duration(600*time.Millisecond),
			wantedDiscrepancySize: time.Duration(5 * time.Minute),
		},

		{
			name: "1st Reconnect 3/3 Packet - 2m700ms timing",
			session: Session{
				previousIncomingDuration: time.Duration(5*time.Minute) + time.Duration(600*time.Millisecond),
				bufferedDuration:         time.Duration(5*time.Minute) + time.Duration(600*time.Millisecond),
				discrepancySize:          time.Duration(5 * time.Minute),
			},
			args: args{
				p: av.Packet{
					Time: time.Duration(2*time.Minute) + time.Duration(700*time.Millisecond),
				},
			},
			wantedPacketTime:      time.Duration(7*time.Minute) + time.Duration(700*time.Millisecond),
			wantedBufferDuration:  time.Duration(7*time.Minute) + time.Duration(700*time.Millisecond),
			wantedDiscrepancySize: time.Duration(5 * time.Minute),
		},

		{
			name: "2nd Reconnect 1/3 Packet - 20ms timing",
			session: Session{
				previousIncomingDuration: time.Duration(7*time.Minute) + time.Duration(700*time.Millisecond),
				bufferedDuration:         time.Duration(7*time.Minute) + time.Duration(700*time.Millisecond),
				discrepancySize:          time.Duration(7*time.Minute) + time.Duration(700*time.Millisecond),
			},
			args: args{
				p: av.Packet{
					Time: time.Duration(20 * time.Millisecond),
				},
			},
			wantedPacketTime:      time.Duration(7*time.Minute) + time.Duration(720*time.Millisecond),
			wantedBufferDuration:  time.Duration(7*time.Minute) + time.Duration(720*time.Millisecond),
			wantedDiscrepancySize: time.Duration(7*time.Minute) + time.Duration(700*time.Millisecond),
		},

		{
			name: "2nd Reconnect 2/3 Packet - 30ms timing",
			session: Session{
				previousIncomingDuration: time.Duration(7*time.Minute) + time.Duration(720*time.Millisecond),
				bufferedDuration:         time.Duration(7*time.Minute) + time.Duration(720*time.Millisecond),
				discrepancySize:          time.Duration(7*time.Minute) + time.Duration(700*time.Millisecond),
			},
			args: args{
				p: av.Packet{
					Time: time.Duration(30 * time.Millisecond),
				},
			},
			wantedPacketTime:      time.Duration(7*time.Minute) + time.Duration(730*time.Millisecond),
			wantedBufferDuration:  time.Duration(7*time.Minute) + time.Duration(730*time.Millisecond),
			wantedDiscrepancySize: time.Duration(7*time.Minute) + time.Duration(700*time.Millisecond),
		},

		{
			name: "2nd Reconnect 3/3 Packet - 40ms timing",
			session: Session{
				previousIncomingDuration: time.Duration(7*time.Minute) + time.Duration(730*time.Millisecond),
				bufferedDuration:         time.Duration(7*time.Minute) + time.Duration(730*time.Millisecond),
				discrepancySize:          time.Duration(7*time.Minute) + time.Duration(700*time.Millisecond),
			},
			args: args{
				p: av.Packet{
					Time: time.Duration(40 * time.Millisecond),
				},
			},
			wantedPacketTime:      time.Duration(7*time.Minute) + time.Duration(740*time.Millisecond),
			wantedBufferDuration:  time.Duration(7*time.Minute) + time.Duration(740*time.Millisecond),
			wantedDiscrepancySize: time.Duration(7*time.Minute) + time.Duration(700*time.Millisecond),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.session.buffer = make(chan av.Packet, 10)

			tt.session.RelayPacket(tt.args.p)

			result := <-tt.session.buffer

			if tt.session.discrepancySize != tt.wantedDiscrepancySize {
				t.Errorf("s.RelayPacket() got %s discrepancy size; wanted %s discrepancy size", tt.session.discrepancySize, tt.wantedDiscrepancySize)
			}

			if tt.session.bufferedDuration != tt.wantedBufferDuration {
				t.Errorf("s.RelayPacket() got %s buffered duration; wanted %s buffered duration", tt.session.bufferedDuration, tt.wantedBufferDuration)
			}

			if result.Time != tt.wantedPacketTime {
				t.Errorf("s.RelayPacket() got %s packet time duration; wanted %s packet time duration", result.Time, tt.wantedPacketTime)
			}

		})
	}
}
