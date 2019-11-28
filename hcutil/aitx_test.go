package hcutil

import "testing"

func TestVerifyMessage(t *testing.T) {
	type args struct {
		msg    string
		addr   Address
		sig    []byte
		pubKey []byte
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyMessage(tt.args.msg, tt.args.addr, tt.args.sig, tt.args.pubKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
