package timestamp

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"time"
)

// Timestamp 同时支持JSON和BSON 时间格式转换为时间戳格式
type Timestamp time.Time

func Now() Timestamp {
	return Timestamp(time.Now())
}

func (t Timestamp) ToTime() time.Time {
	return time.Time(t)
}

func (t Timestamp) MarshalBinary() ([]byte, error) {
	timestamp := time.Time(t).UnixMilli()
	return []byte(fmt.Sprint(timestamp)), nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return t.MarshalBinary()
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	//*t, err = timestamp.Parse(`"`+RFC3339+`"`, string(data))
	return err
}

//MarshalBSONValue marshal bson value
func (t Timestamp) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.TypeDateTime, bsoncore.AppendTime(nil, time.Time(t)), nil
}

// UnmarshalBSON unmarshal bson
func (t *Timestamp) UnmarshalBSON(data []byte) (err error) {
	readTime, _, _ := bsoncore.ReadTime(data)
	*t = Timestamp(readTime)
	return nil
}

//func (t *Timestamp) UnmarshalBSONValue(vType bsontype.Type, value []byte) error {
//
//	dt, rem, ok := bsoncore.ReadDateTime(value)
//
//	dateTime := bsonx.DateTime(dt)
//	fmt.Println(dateTime, rem, ok)
//	return nil
//
//}
