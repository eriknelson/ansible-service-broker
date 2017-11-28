// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/api/distribution.proto

/*
Package distribution is a generated protocol buffer package.

It is generated from these files:
	google/api/distribution.proto

It has these top-level messages:
	Distribution
*/
package distribution

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import _ "github.com/golang/protobuf/ptypes/any"
import _ "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Distribution contains summary statistics for a population of values and,
// optionally, a histogram representing the distribution of those values across
// a specified set of histogram buckets.
//
// The summary statistics are the count, mean, sum of the squared deviation from
// the mean, the minimum, and the maximum of the set of population of values.
//
// The histogram is based on a sequence of buckets and gives a count of values
// that fall into each bucket.  The boundaries of the buckets are given either
// explicitly or by specifying parameters for a method of computing them
// (buckets of fixed width or buckets of exponentially increasing width).
//
// Although it is not forbidden, it is generally a bad idea to include
// non-finite values (infinities or NaNs) in the population of values, as this
// will render the `mean` and `sum_of_squared_deviation` fields meaningless.
type Distribution struct {
	// The number of values in the population. Must be non-negative.
	Count int64 `protobuf:"varint,1,opt,name=count" json:"count,omitempty"`
	// The arithmetic mean of the values in the population. If `count` is zero
	// then this field must be zero.
	Mean float64 `protobuf:"fixed64,2,opt,name=mean" json:"mean,omitempty"`
	// The sum of squared deviations from the mean of the values in the
	// population.  For values x_i this is:
	//
	//     Sum[i=1..n]((x_i - mean)^2)
	//
	// Knuth, "The Art of Computer Programming", Vol. 2, page 323, 3rd edition
	// describes Welford's method for accumulating this sum in one pass.
	//
	// If `count` is zero then this field must be zero.
	SumOfSquaredDeviation float64 `protobuf:"fixed64,3,opt,name=sum_of_squared_deviation,json=sumOfSquaredDeviation" json:"sum_of_squared_deviation,omitempty"`
	// If specified, contains the range of the population values. The field
	// must not be present if the `count` is zero.
	Range *Distribution_Range `protobuf:"bytes,4,opt,name=range" json:"range,omitempty"`
	// Defines the histogram bucket boundaries.
	BucketOptions *Distribution_BucketOptions `protobuf:"bytes,6,opt,name=bucket_options,json=bucketOptions" json:"bucket_options,omitempty"`
	// If `bucket_options` is given, then the sum of the values in `bucket_counts`
	// must equal the value in `count`.  If `bucket_options` is not given, no
	// `bucket_counts` fields may be given.
	//
	// Bucket counts are given in order under the numbering scheme described
	// above (the underflow bucket has number 0; the finite buckets, if any,
	// have numbers 1 through N-2; the overflow bucket has number N-1).
	//
	// The size of `bucket_counts` must be no greater than N as defined in
	// `bucket_options`.
	//
	// Any suffix of trailing zero bucket_count fields may be omitted.
	BucketCounts []int64 `protobuf:"varint,7,rep,packed,name=bucket_counts,json=bucketCounts" json:"bucket_counts,omitempty"`
}

func (m *Distribution) Reset()                    { *m = Distribution{} }
func (m *Distribution) String() string            { return proto.CompactTextString(m) }
func (*Distribution) ProtoMessage()               {}
func (*Distribution) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Distribution) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *Distribution) GetMean() float64 {
	if m != nil {
		return m.Mean
	}
	return 0
}

func (m *Distribution) GetSumOfSquaredDeviation() float64 {
	if m != nil {
		return m.SumOfSquaredDeviation
	}
	return 0
}

func (m *Distribution) GetRange() *Distribution_Range {
	if m != nil {
		return m.Range
	}
	return nil
}

func (m *Distribution) GetBucketOptions() *Distribution_BucketOptions {
	if m != nil {
		return m.BucketOptions
	}
	return nil
}

func (m *Distribution) GetBucketCounts() []int64 {
	if m != nil {
		return m.BucketCounts
	}
	return nil
}

// The range of the population values.
type Distribution_Range struct {
	// The minimum of the population values.
	Min float64 `protobuf:"fixed64,1,opt,name=min" json:"min,omitempty"`
	// The maximum of the population values.
	Max float64 `protobuf:"fixed64,2,opt,name=max" json:"max,omitempty"`
}

func (m *Distribution_Range) Reset()                    { *m = Distribution_Range{} }
func (m *Distribution_Range) String() string            { return proto.CompactTextString(m) }
func (*Distribution_Range) ProtoMessage()               {}
func (*Distribution_Range) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *Distribution_Range) GetMin() float64 {
	if m != nil {
		return m.Min
	}
	return 0
}

func (m *Distribution_Range) GetMax() float64 {
	if m != nil {
		return m.Max
	}
	return 0
}

// A Distribution may optionally contain a histogram of the values in the
// population.  The histogram is given in `bucket_counts` as counts of values
// that fall into one of a sequence of non-overlapping buckets.  The sequence
// of buckets is described by `bucket_options`.
//
// A bucket specifies an inclusive lower bound and exclusive upper bound for
// the values that are counted for that bucket.  The upper bound of a bucket
// is strictly greater than the lower bound.
//
// The sequence of N buckets for a Distribution consists of an underflow
// bucket (number 0), zero or more finite buckets (number 1 through N - 2) and
// an overflow bucket (number N - 1).  The buckets are contiguous:  the lower
// bound of bucket i (i > 0) is the same as the upper bound of bucket i - 1.
// The buckets span the whole range of finite values: lower bound of the
// underflow bucket is -infinity and the upper bound of the overflow bucket is
// +infinity.  The finite buckets are so-called because both bounds are
// finite.
//
// `BucketOptions` describes bucket boundaries in one of three ways.  Two
// describe the boundaries by giving parameters for a formula to generate
// boundaries and one gives the bucket boundaries explicitly.
//
// If `bucket_boundaries` is not given, then no `bucket_counts` may be given.
type Distribution_BucketOptions struct {
	// Exactly one of these three fields must be set.
	//
	// Types that are valid to be assigned to Options:
	//	*Distribution_BucketOptions_LinearBuckets
	//	*Distribution_BucketOptions_ExponentialBuckets
	//	*Distribution_BucketOptions_ExplicitBuckets
	Options isDistribution_BucketOptions_Options `protobuf_oneof:"options"`
}

func (m *Distribution_BucketOptions) Reset()                    { *m = Distribution_BucketOptions{} }
func (m *Distribution_BucketOptions) String() string            { return proto.CompactTextString(m) }
func (*Distribution_BucketOptions) ProtoMessage()               {}
func (*Distribution_BucketOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

type isDistribution_BucketOptions_Options interface {
	isDistribution_BucketOptions_Options()
}

type Distribution_BucketOptions_LinearBuckets struct {
	LinearBuckets *Distribution_BucketOptions_Linear `protobuf:"bytes,1,opt,name=linear_buckets,json=linearBuckets,oneof"`
}
type Distribution_BucketOptions_ExponentialBuckets struct {
	ExponentialBuckets *Distribution_BucketOptions_Exponential `protobuf:"bytes,2,opt,name=exponential_buckets,json=exponentialBuckets,oneof"`
}
type Distribution_BucketOptions_ExplicitBuckets struct {
	ExplicitBuckets *Distribution_BucketOptions_Explicit `protobuf:"bytes,3,opt,name=explicit_buckets,json=explicitBuckets,oneof"`
}

func (*Distribution_BucketOptions_LinearBuckets) isDistribution_BucketOptions_Options()      {}
func (*Distribution_BucketOptions_ExponentialBuckets) isDistribution_BucketOptions_Options() {}
func (*Distribution_BucketOptions_ExplicitBuckets) isDistribution_BucketOptions_Options()    {}

func (m *Distribution_BucketOptions) GetOptions() isDistribution_BucketOptions_Options {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *Distribution_BucketOptions) GetLinearBuckets() *Distribution_BucketOptions_Linear {
	if x, ok := m.GetOptions().(*Distribution_BucketOptions_LinearBuckets); ok {
		return x.LinearBuckets
	}
	return nil
}

func (m *Distribution_BucketOptions) GetExponentialBuckets() *Distribution_BucketOptions_Exponential {
	if x, ok := m.GetOptions().(*Distribution_BucketOptions_ExponentialBuckets); ok {
		return x.ExponentialBuckets
	}
	return nil
}

func (m *Distribution_BucketOptions) GetExplicitBuckets() *Distribution_BucketOptions_Explicit {
	if x, ok := m.GetOptions().(*Distribution_BucketOptions_ExplicitBuckets); ok {
		return x.ExplicitBuckets
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Distribution_BucketOptions) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Distribution_BucketOptions_OneofMarshaler, _Distribution_BucketOptions_OneofUnmarshaler, _Distribution_BucketOptions_OneofSizer, []interface{}{
		(*Distribution_BucketOptions_LinearBuckets)(nil),
		(*Distribution_BucketOptions_ExponentialBuckets)(nil),
		(*Distribution_BucketOptions_ExplicitBuckets)(nil),
	}
}

func _Distribution_BucketOptions_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Distribution_BucketOptions)
	// options
	switch x := m.Options.(type) {
	case *Distribution_BucketOptions_LinearBuckets:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.LinearBuckets); err != nil {
			return err
		}
	case *Distribution_BucketOptions_ExponentialBuckets:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ExponentialBuckets); err != nil {
			return err
		}
	case *Distribution_BucketOptions_ExplicitBuckets:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ExplicitBuckets); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Distribution_BucketOptions.Options has unexpected type %T", x)
	}
	return nil
}

func _Distribution_BucketOptions_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Distribution_BucketOptions)
	switch tag {
	case 1: // options.linear_buckets
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Distribution_BucketOptions_Linear)
		err := b.DecodeMessage(msg)
		m.Options = &Distribution_BucketOptions_LinearBuckets{msg}
		return true, err
	case 2: // options.exponential_buckets
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Distribution_BucketOptions_Exponential)
		err := b.DecodeMessage(msg)
		m.Options = &Distribution_BucketOptions_ExponentialBuckets{msg}
		return true, err
	case 3: // options.explicit_buckets
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Distribution_BucketOptions_Explicit)
		err := b.DecodeMessage(msg)
		m.Options = &Distribution_BucketOptions_ExplicitBuckets{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Distribution_BucketOptions_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Distribution_BucketOptions)
	// options
	switch x := m.Options.(type) {
	case *Distribution_BucketOptions_LinearBuckets:
		s := proto.Size(x.LinearBuckets)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Distribution_BucketOptions_ExponentialBuckets:
		s := proto.Size(x.ExponentialBuckets)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Distribution_BucketOptions_ExplicitBuckets:
		s := proto.Size(x.ExplicitBuckets)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Specify a sequence of buckets that all have the same width (except
// overflow and underflow).  Each bucket represents a constant absolute
// uncertainty on the specific value in the bucket.
//
// Defines `num_finite_buckets + 2` (= N) buckets with these boundaries for
// bucket `i`:
//
//    Upper bound (0 <= i < N-1):     offset + (width * i).
//    Lower bound (1 <= i < N):       offset + (width * (i - 1)).
type Distribution_BucketOptions_Linear struct {
	// Must be greater than 0.
	NumFiniteBuckets int32 `protobuf:"varint,1,opt,name=num_finite_buckets,json=numFiniteBuckets" json:"num_finite_buckets,omitempty"`
	// Must be greater than 0.
	Width float64 `protobuf:"fixed64,2,opt,name=width" json:"width,omitempty"`
	// Lower bound of the first bucket.
	Offset float64 `protobuf:"fixed64,3,opt,name=offset" json:"offset,omitempty"`
}

func (m *Distribution_BucketOptions_Linear) Reset()         { *m = Distribution_BucketOptions_Linear{} }
func (m *Distribution_BucketOptions_Linear) String() string { return proto.CompactTextString(m) }
func (*Distribution_BucketOptions_Linear) ProtoMessage()    {}
func (*Distribution_BucketOptions_Linear) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{0, 1, 0}
}

func (m *Distribution_BucketOptions_Linear) GetNumFiniteBuckets() int32 {
	if m != nil {
		return m.NumFiniteBuckets
	}
	return 0
}

func (m *Distribution_BucketOptions_Linear) GetWidth() float64 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *Distribution_BucketOptions_Linear) GetOffset() float64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

// Specify a sequence of buckets that have a width that is proportional to
// the value of the lower bound.  Each bucket represents a constant relative
// uncertainty on a specific value in the bucket.
//
// Defines `num_finite_buckets + 2` (= N) buckets with these boundaries for
// bucket i:
//
//    Upper bound (0 <= i < N-1):     scale * (growth_factor ^ i).
//    Lower bound (1 <= i < N):       scale * (growth_factor ^ (i - 1)).
type Distribution_BucketOptions_Exponential struct {
	// Must be greater than 0.
	NumFiniteBuckets int32 `protobuf:"varint,1,opt,name=num_finite_buckets,json=numFiniteBuckets" json:"num_finite_buckets,omitempty"`
	// Must be greater than 1.
	GrowthFactor float64 `protobuf:"fixed64,2,opt,name=growth_factor,json=growthFactor" json:"growth_factor,omitempty"`
	// Must be greater than 0.
	Scale float64 `protobuf:"fixed64,3,opt,name=scale" json:"scale,omitempty"`
}

func (m *Distribution_BucketOptions_Exponential) Reset() {
	*m = Distribution_BucketOptions_Exponential{}
}
func (m *Distribution_BucketOptions_Exponential) String() string { return proto.CompactTextString(m) }
func (*Distribution_BucketOptions_Exponential) ProtoMessage()    {}
func (*Distribution_BucketOptions_Exponential) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{0, 1, 1}
}

func (m *Distribution_BucketOptions_Exponential) GetNumFiniteBuckets() int32 {
	if m != nil {
		return m.NumFiniteBuckets
	}
	return 0
}

func (m *Distribution_BucketOptions_Exponential) GetGrowthFactor() float64 {
	if m != nil {
		return m.GrowthFactor
	}
	return 0
}

func (m *Distribution_BucketOptions_Exponential) GetScale() float64 {
	if m != nil {
		return m.Scale
	}
	return 0
}

// A set of buckets with arbitrary widths.
//
// Defines `size(bounds) + 1` (= N) buckets with these boundaries for
// bucket i:
//
//    Upper bound (0 <= i < N-1):     bounds[i]
//    Lower bound (1 <= i < N);       bounds[i - 1]
//
// There must be at least one element in `bounds`.  If `bounds` has only one
// element, there are no finite buckets, and that single element is the
// common boundary of the overflow and underflow buckets.
type Distribution_BucketOptions_Explicit struct {
	// The values must be monotonically increasing.
	Bounds []float64 `protobuf:"fixed64,1,rep,packed,name=bounds" json:"bounds,omitempty"`
}

func (m *Distribution_BucketOptions_Explicit) Reset()         { *m = Distribution_BucketOptions_Explicit{} }
func (m *Distribution_BucketOptions_Explicit) String() string { return proto.CompactTextString(m) }
func (*Distribution_BucketOptions_Explicit) ProtoMessage()    {}
func (*Distribution_BucketOptions_Explicit) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{0, 1, 2}
}

func (m *Distribution_BucketOptions_Explicit) GetBounds() []float64 {
	if m != nil {
		return m.Bounds
	}
	return nil
}

func init() {
	proto.RegisterType((*Distribution)(nil), "google.api.Distribution")
	proto.RegisterType((*Distribution_Range)(nil), "google.api.Distribution.Range")
	proto.RegisterType((*Distribution_BucketOptions)(nil), "google.api.Distribution.BucketOptions")
	proto.RegisterType((*Distribution_BucketOptions_Linear)(nil), "google.api.Distribution.BucketOptions.Linear")
	proto.RegisterType((*Distribution_BucketOptions_Exponential)(nil), "google.api.Distribution.BucketOptions.Exponential")
	proto.RegisterType((*Distribution_BucketOptions_Explicit)(nil), "google.api.Distribution.BucketOptions.Explicit")
}

func init() { proto.RegisterFile("google/api/distribution.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 544 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x5d, 0x6f, 0xd3, 0x30,
	0x14, 0x5d, 0x96, 0xb5, 0x85, 0xdb, 0x0f, 0x8a, 0x19, 0x28, 0x44, 0x7c, 0x54, 0x9b, 0x84, 0x2a,
	0x01, 0x89, 0x54, 0x90, 0x78, 0xe0, 0xad, 0x1b, 0x53, 0x1f, 0x40, 0x9b, 0x8c, 0xc4, 0x03, 0x42,
	0x8a, 0x9c, 0xc4, 0xc9, 0x0c, 0x89, 0x1d, 0x62, 0x67, 0x2b, 0xef, 0xfc, 0x29, 0xfe, 0x1d, 0x8a,
	0xed, 0x6e, 0x19, 0x08, 0xa9, 0xbc, 0xf9, 0xde, 0x73, 0x7c, 0xce, 0xb9, 0x57, 0x71, 0xe0, 0x71,
	0x2e, 0x44, 0x5e, 0xd0, 0x90, 0x54, 0x2c, 0x4c, 0x99, 0x54, 0x35, 0x8b, 0x1b, 0xc5, 0x04, 0x0f,
	0xaa, 0x5a, 0x28, 0x81, 0xc0, 0xc0, 0x01, 0xa9, 0x98, 0xff, 0xa8, 0x43, 0x25, 0x9c, 0x0b, 0x45,
	0x5a, 0xa2, 0x34, 0x4c, 0xff, 0xa1, 0x45, 0x75, 0x15, 0x37, 0x59, 0x48, 0xf8, 0x0f, 0x0b, 0x3d,
	0xfd, 0x13, 0x52, 0xac, 0xa4, 0x52, 0x91, 0xb2, 0x32, 0x84, 0x83, 0x9f, 0x03, 0x18, 0x1d, 0x77,
	0xcc, 0xd1, 0x3e, 0xf4, 0x12, 0xd1, 0x70, 0xe5, 0x39, 0x33, 0x67, 0xee, 0x62, 0x53, 0x20, 0x04,
	0x7b, 0x25, 0x25, 0xdc, 0xdb, 0x9d, 0x39, 0x73, 0x07, 0xeb, 0x33, 0x7a, 0x03, 0x9e, 0x6c, 0xca,
	0x48, 0x64, 0x91, 0xfc, 0xde, 0x90, 0x9a, 0xa6, 0x51, 0x4a, 0x2f, 0x98, 0x4e, 0xe6, 0xb9, 0x9a,
	0x77, 0x5f, 0x36, 0xe5, 0x69, 0xf6, 0xd1, 0xa0, 0xc7, 0x1b, 0x10, 0xbd, 0x86, 0x5e, 0x4d, 0x78,
	0x4e, 0xbd, 0xbd, 0x99, 0x33, 0x1f, 0x2e, 0x9e, 0x04, 0xd7, 0x93, 0x06, 0xdd, 0x2c, 0x01, 0x6e,
	0x59, 0xd8, 0x90, 0xd1, 0x07, 0x98, 0xc4, 0x4d, 0xf2, 0x8d, 0xaa, 0x48, 0x54, 0x7a, 0x7a, 0xaf,
	0xaf, 0xaf, 0x3f, 0xfb, 0xe7, 0xf5, 0xa5, 0xa6, 0x9f, 0x1a, 0x36, 0x1e, 0xc7, 0xdd, 0x12, 0x1d,
	0x82, 0x6d, 0x44, 0x7a, 0x42, 0xe9, 0x0d, 0x66, 0xee, 0xdc, 0xc5, 0x23, 0xd3, 0x3c, 0xd2, 0x3d,
	0xff, 0x39, 0xf4, 0x74, 0x06, 0x34, 0x05, 0xb7, 0x64, 0x5c, 0xef, 0xc4, 0xc1, 0xed, 0x51, 0x77,
	0xc8, 0xda, 0x2e, 0xa4, 0x3d, 0xfa, 0xbf, 0xf6, 0x60, 0x7c, 0xc3, 0x12, 0x7d, 0x82, 0x49, 0xc1,
	0x38, 0x25, 0x75, 0x64, 0x54, 0xa5, 0x16, 0x18, 0x2e, 0x5e, 0x6e, 0x17, 0x39, 0x78, 0xaf, 0x2f,
	0xaf, 0x76, 0xf0, 0xd8, 0xc8, 0x18, 0x54, 0x22, 0x0a, 0xf7, 0xe8, 0xba, 0x12, 0x9c, 0x72, 0xc5,
	0x48, 0x71, 0x25, 0xbe, 0xab, 0xc5, 0x17, 0x5b, 0x8a, 0xbf, 0xbb, 0x56, 0x58, 0xed, 0x60, 0xd4,
	0x11, 0xdc, 0xd8, 0x7c, 0x81, 0x29, 0x5d, 0x57, 0x05, 0x4b, 0x98, 0xba, 0xf2, 0x70, 0xb5, 0x47,
	0xb8, 0xbd, 0x87, 0xbe, 0xbe, 0xda, 0xc1, 0x77, 0x36, 0x52, 0x56, 0xdd, 0x4f, 0xa1, 0x6f, 0xe6,
	0x43, 0x2f, 0x00, 0xf1, 0xa6, 0x8c, 0x32, 0xc6, 0x99, 0xa2, 0x37, 0x56, 0xd5, 0xc3, 0x53, 0xde,
	0x94, 0x27, 0x1a, 0xd8, 0xa4, 0xda, 0x87, 0xde, 0x25, 0x4b, 0xd5, 0xb9, 0x5d, 0xbd, 0x29, 0xd0,
	0x03, 0xe8, 0x8b, 0x2c, 0x93, 0x54, 0xd9, 0x4f, 0xcf, 0x56, 0xfe, 0x05, 0x0c, 0x3b, 0x83, 0xfe,
	0xa7, 0xd5, 0x21, 0x8c, 0xf3, 0x5a, 0x5c, 0xaa, 0xf3, 0x28, 0x23, 0x89, 0x12, 0xb5, 0xb5, 0x1c,
	0x99, 0xe6, 0x89, 0xee, 0xb5, 0x79, 0x64, 0x42, 0x0a, 0x6a, 0x8d, 0x4d, 0xe1, 0x1f, 0xc0, 0xad,
	0xcd, 0xf0, 0x6d, 0xb6, 0x58, 0x34, 0x3c, 0x6d, 0x8d, 0xdc, 0x36, 0x9b, 0xa9, 0x96, 0xb7, 0x61,
	0x60, 0x3f, 0xe5, 0xe5, 0x57, 0x98, 0x24, 0xa2, 0xec, 0x6c, 0x75, 0x79, 0xb7, 0xbb, 0xd6, 0xb3,
	0xf6, 0xad, 0x9e, 0x39, 0x9f, 0x8f, 0x2c, 0x21, 0x17, 0x05, 0xe1, 0x79, 0x20, 0xea, 0x3c, 0xcc,
	0x29, 0xd7, 0x2f, 0x39, 0x34, 0x10, 0xa9, 0x98, 0xfc, 0xeb, 0x8f, 0xf2, 0xb6, 0x5b, 0xc4, 0x7d,
	0xcd, 0x7f, 0xf5, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x88, 0x8e, 0xc5, 0x4b, 0x80, 0x04, 0x00, 0x00,
}
