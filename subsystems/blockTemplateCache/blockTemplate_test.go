package blockTemplateCache

import (
	"encoding/hex"
	"github.com/Snipa22/go-tari-grpc-lib/v3/tari_generated"
	"reflect"
	"testing"
)

func Test_hashDevTest(t *testing.T) {
	hexString := "fc0bc2f078c716b6166dd497f5fd3cab376eb1665f5f8ee149ffb0893961b69d"
	data, _ := hex.DecodeString(hexString)
	tests := []struct {
		name string
		want []byte
	}{
		{name: "Verify we get the same result as Tari", want: data},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hashDevTest(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hashDevTest() = %x, want %v", got, tt.want)
			}
		})
	}
}

/*
	func TestHashBlockHeader(t *testing.T) {
		type args struct {
			header *tari_generated.BlockHeader
		}
		heightZeroResponse, _ := hex.DecodeString("55f5e1fb50c7146f54bd171fc8bc1eb658bff07b7befe5d140ddccffe5aff41c")
		fmt.Println(heightZeroResponse)
		tests := []struct {
			name    string
			args    args
			want    []byte
			wantErr bool
		}{
			{
				name: "Verify Tari result with a static block template - Height 0",
				args: args{&tari_generated.BlockHeader{
					Hash:              nil,
					Version:           2,
					Height:            0,
					PrevHash:          make([]byte, 32),
					Timestamp:         946688461,
					OutputMr:          make([]byte, 32),
					BlockOutputMr:     make([]byte, 32),
					KernelMr:          make([]byte, 32),
					InputMr:           make([]byte, 32),
					TotalKernelOffset: make([]byte, 32),
					Nonce:             4,
					Pow:               &tari_generated.ProofOfWork{PowAlgo: 1},
					KernelMmrSize:     0,
					OutputMmrSize:     0,
					TotalScriptOffset: make([]byte, 32),
					ValidatorNodeMr:   make([]byte, 32),
					ValidatorNodeSize: 0,
				}},
				want:    heightZeroResponse,
				wantErr: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := HashBlockHeader(tt.args.header)
				if (err != nil) != tt.wantErr {
					t.Errorf("HashBlockHeader() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("HashBlockHeader() got = %v, want %v", got, tt.want)
				}
			})
		}
	}
*/
func TestGetHeaderDiff(t *testing.T) {
	type args struct {
		header     *tari_generated.BlockHeader
		miningHash []byte
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "Verify Tari result with a static block template - Height 0",
			args: args{&tari_generated.BlockHeader{
				Hash:              nil,
				Version:           2,
				Height:            0,
				PrevHash:          make([]byte, 32),
				Timestamp:         946688461,
				OutputMr:          make([]byte, 32),
				BlockOutputMr:     make([]byte, 32),
				KernelMr:          make([]byte, 32),
				InputMr:           make([]byte, 32),
				TotalKernelOffset: make([]byte, 32),
				Nonce:             4,
				Pow:               &tari_generated.ProofOfWork{PowAlgo: 1},
				KernelMmrSize:     0,
				OutputMmrSize:     0,
				TotalScriptOffset: make([]byte, 32),
				ValidatorNodeMr:   make([]byte, 32),
				ValidatorNodeSize: 0,
			}, []byte{85, 245, 225, 251, 80, 199, 20, 111, 84, 189, 23, 31, 200, 188, 30, 182, 88, 191, 240, 123, 123, 239, 229, 209, 64, 221, 204, 255, 229, 175, 244, 28}},
			want:    4,
			wantErr: false,
		},
		{
			name: "Generate known tari result w/ target verification from daemon",
			args: args{&tari_generated.BlockHeader{
				Hash:              nil,
				Version:           2,
				Height:            0,
				PrevHash:          make([]byte, 32),
				Timestamp:         946688461,
				OutputMr:          make([]byte, 32),
				BlockOutputMr:     make([]byte, 32),
				KernelMr:          make([]byte, 32),
				InputMr:           make([]byte, 32),
				TotalKernelOffset: make([]byte, 32),
				Nonce:             12238891555,
				Pow:               &tari_generated.ProofOfWork{PowAlgo: 1},
				KernelMmrSize:     0,
				OutputMmrSize:     0,
				TotalScriptOffset: make([]byte, 32),
				ValidatorNodeMr:   make([]byte, 32),
				ValidatorNodeSize: 0,
			}, []byte{153, 120, 61, 169, 149, 9, 184, 83, 201, 15, 60, 16, 190, 168, 164, 224, 5, 241, 68, 141, 124, 192, 137, 75, 255, 62, 93, 82, 170, 71, 65, 82}},
			want:    13645861923,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHeaderDiff(tt.args.header, tt.args.miningHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHeaderDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetHeaderDiff() got = %v, want %v", got, tt.want)
			}
		})
	}
}
