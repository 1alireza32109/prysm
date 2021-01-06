package attestations

import (
	"testing"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/shared/aggregation"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
)

func TestAggregateAttestations_MaxCover_NewMaxCover(t *testing.T) {
	type args struct {
		atts []*ethpb.Attestation
	}
	tests := []struct {
		name      string
		args      args
		want      *aggregation.MaxCoverProblem
		wantedErr string
	}{
		{
			name: "nil attestations",
			args: args{
				atts: nil,
			},
			wantedErr: ErrInvalidAttestationCount.Error(),
		},
		{
			name: "no attestations",
			args: args{
				atts: []*ethpb.Attestation{},
			},
			wantedErr: ErrInvalidAttestationCount.Error(),
		},
		{
			name: "attestations of different bitlist length",
			args: args{
				atts: []*ethpb.Attestation{
					{AggregationBits: bitfield.NewBitlist(64)},
					{AggregationBits: bitfield.NewBitlist(128)},
				},
			},
			wantedErr: aggregation.ErrBitsDifferentLen.Error(),
		},
		{
			name: "single attestation",
			args: args{
				atts: []*ethpb.Attestation{
					{AggregationBits: bitfield.Bitlist{0b00001010, 0b1}},
				},
			},
			want: &aggregation.MaxCoverProblem{
				Candidates: aggregation.MaxCoverCandidates{
					aggregation.NewMaxCoverCandidate(0, &bitfield.Bitlist{0b00001010, 0b1}),
				},
			},
		},
		{
			name: "multiple attestations",
			args: args{
				atts: []*ethpb.Attestation{
					{AggregationBits: bitfield.Bitlist{0b00001010, 0b1}},
					{AggregationBits: bitfield.Bitlist{0b00101010, 0b1}},
					{AggregationBits: bitfield.Bitlist{0b11111010, 0b1}},
					{AggregationBits: bitfield.Bitlist{0b00000010, 0b1}},
					{AggregationBits: bitfield.Bitlist{0b00000001, 0b1}},
				},
			},
			want: &aggregation.MaxCoverProblem{
				Candidates: aggregation.MaxCoverCandidates{
					aggregation.NewMaxCoverCandidate(0, &bitfield.Bitlist{0b00001010, 0b1}),
					aggregation.NewMaxCoverCandidate(1, &bitfield.Bitlist{0b00101010, 0b1}),
					aggregation.NewMaxCoverCandidate(2, &bitfield.Bitlist{0b11111010, 0b1}),
					aggregation.NewMaxCoverCandidate(3, &bitfield.Bitlist{0b00000010, 0b1}),
					aggregation.NewMaxCoverCandidate(4, &bitfield.Bitlist{0b00000001, 0b1}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMaxCover(tt.args.atts)
			if tt.wantedErr != "" {
				assert.ErrorContains(t, tt.wantedErr, err)
			} else {
				assert.NoError(t, err)
				assert.DeepEqual(t, tt.want, got)
			}
		})
	}
}

func TestAggregateAttestations_MaxCover_AttList_validate(t *testing.T) {
	tests := []struct {
		name      string
		atts      attList
		wantedErr string
	}{
		{
			name:      "nil list",
			atts:      nil,
			wantedErr: "nil list",
		},
		{
			name:      "empty list",
			atts:      attList{},
			wantedErr: "empty list",
		},
		{
			name:      "first bitlist is nil",
			atts:      attList{&ethpb.Attestation{}},
			wantedErr: "bitlist cannot be nil or empty",
		},
		{
			name: "non first bitlist is nil",
			atts: attList{
				&ethpb.Attestation{AggregationBits: bitfield.NewBitlist(64)},
				&ethpb.Attestation{},
			},
			wantedErr: "bitlist cannot be nil or empty",
		},
		{
			name: "first bitlist is empty",
			atts: attList{
				&ethpb.Attestation{AggregationBits: bitfield.Bitlist{}},
			},
			wantedErr: "bitlist cannot be nil or empty",
		},
		{
			name: "non first bitlist is empty",
			atts: attList{
				&ethpb.Attestation{AggregationBits: bitfield.NewBitlist(64)},
				&ethpb.Attestation{AggregationBits: bitfield.Bitlist{}},
			},
			wantedErr: "bitlist cannot be nil or empty",
		},
		{
			name: "bitlists of non equal length",
			atts: attList{
				&ethpb.Attestation{AggregationBits: bitfield.NewBitlist(64)},
				&ethpb.Attestation{AggregationBits: bitfield.NewBitlist(64)},
				&ethpb.Attestation{AggregationBits: bitfield.NewBitlist(63)},
				&ethpb.Attestation{AggregationBits: bitfield.NewBitlist(64)},
			},
			wantedErr: "bitlists of different length",
		},
		{
			name: "valid bitlists",
			atts: attList{
				&ethpb.Attestation{AggregationBits: bitfield.NewBitlist(64)},
				&ethpb.Attestation{AggregationBits: bitfield.NewBitlist(64)},
				&ethpb.Attestation{AggregationBits: bitfield.NewBitlist(64)},
				&ethpb.Attestation{AggregationBits: bitfield.NewBitlist(64)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.atts.validate()
			if tt.wantedErr != "" {
				assert.ErrorContains(t, tt.wantedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
