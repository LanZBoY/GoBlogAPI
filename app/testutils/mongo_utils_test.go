package testutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAppendCallArgs(t *testing.T) {
	t.Run("withValues", func(t *testing.T) {
		fixed := []any{"a", 1}
		variadic := []int{2, 3}
		out := AppendCallArgs(fixed, variadic)
		assert.Equal(t, []any{"a", 1, 2, 3}, out)
	})

	t.Run("empty", func(t *testing.T) {
		fixed := []any{"a"}
		variadic := []string{}
		out := AppendCallArgs(fixed, variadic)
		assert.Equal(t, []any{"a"}, out)
	})
}

func TestMockCollection_CountDocuments(t *testing.T) {
	ctx := context.TODO()
	filter := bson.M{"f": 1}
	opts := options.Count()

	mc := new(MockCollection)
	mc.On("CountDocuments", ctx, filter, opts).Return(int64(5), nil)

	n, err := mc.CountDocuments(ctx, filter, opts)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), n)
	mc.AssertExpectations(t)
}

func TestMockCollection_Find(t *testing.T) {
	ctx := context.TODO()
	filter := bson.M{}
	opts := options.Find()
	cursor := new(MockCursor)

	mc := new(MockCollection)
	mc.On("Find", ctx, filter, opts).Return(cursor, nil)

	c, err := mc.Find(ctx, filter, opts)
	assert.NoError(t, err)
	assert.Same(t, cursor, c)
	mc.AssertExpectations(t)
}

func TestMockCollection_FindOne(t *testing.T) {
	ctx := context.TODO()
	filter := bson.M{"a": 1}
	opts := options.FindOne()
	sr := new(MockSingleResult)

	mc := new(MockCollection)
	mc.On("FindOne", ctx, filter, opts).Return(sr)

	r := mc.FindOne(ctx, filter, opts)
	assert.Same(t, sr, r)
	mc.AssertExpectations(t)
}

func TestMockCollection_UpdateOne(t *testing.T) {
	ctx := context.TODO()
	filter := bson.M{"x": 1}
	updateDoc := bson.M{"$set": bson.M{"x": 2}}
	opts := options.Update()
	res := &mongo.UpdateResult{MatchedCount: 1}

	mc := new(MockCollection)
	mc.On("UpdateOne", ctx, filter, updateDoc, opts).Return(res, nil)

	ur, err := mc.UpdateOne(ctx, filter, updateDoc, opts)
	assert.NoError(t, err)
	assert.Same(t, res, ur)
	mc.AssertExpectations(t)
}

func TestMockCollection_DeleteOne(t *testing.T) {
	ctx := context.TODO()
	filter := bson.M{"x": 1}
	opts := options.Delete()
	res := &mongo.DeleteResult{DeletedCount: 1}

	mc := new(MockCollection)
	mc.On("DeleteOne", ctx, filter, opts).Return(res, nil)

	dr, err := mc.DeleteOne(ctx, filter, opts)
	assert.NoError(t, err)
	assert.Same(t, res, dr)
	mc.AssertExpectations(t)
}

func TestMockCollection_InsertOne(t *testing.T) {
	ctx := context.TODO()
	data := bson.M{"a": 1}
	opts := options.InsertOne()
	res := &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}

	mc := new(MockCollection)
	mc.On("InsertOne", ctx, data, opts).Return(res, nil)

	ir, err := mc.InsertOne(ctx, data, opts)
	assert.NoError(t, err)
	assert.Same(t, res, ir)
	mc.AssertExpectations(t)
}

func TestMockCollection_Aggregate(t *testing.T) {
	ctx := context.TODO()
	pipeline := bson.A{bson.D{{"stage", 1}}}
	opts := options.Aggregate()
	cursor := new(MockCursor)

	mc := new(MockCollection)
	mc.On("Aggregate", ctx, pipeline, opts).Return(cursor, nil)

	c, err := mc.Aggregate(ctx, pipeline, opts)
	assert.NoError(t, err)
	assert.Same(t, cursor, c)
	mc.AssertExpectations(t)
}

func TestMockCursorMethods(t *testing.T) {
	cursor := new(MockCursor)
	ctx := context.TODO()
	var out []bson.M
	cursor.On("All", ctx, &out).Return(nil)
	cursor.On("Next", ctx).Return(true)
	cursor.On("Decode", &out).Return(nil)
	cursor.On("Close", ctx).Return(nil)

	assert.NoError(t, cursor.All(ctx, &out))
	assert.True(t, cursor.Next(ctx))
	assert.NoError(t, cursor.Decode(&out))
	assert.NoError(t, cursor.Close(ctx))
	cursor.AssertExpectations(t)
}

func TestMockSingleResult_Decode(t *testing.T) {
	msr := new(MockSingleResult)
	var v bson.M
	msr.On("Decode", &v).Return(nil)

	assert.NoError(t, msr.Decode(&v))
	msr.AssertExpectations(t)
}

func TestMockMongoUtils_CountDocumentWithPipeline(t *testing.T) {
	mu := new(MockMongoUtils)
	ctx := context.TODO()
	aggr := new(MockCollection)
	pipeline := bson.A{}
	mu.On("CountDocumentWithPipeline", ctx, aggr, pipeline).Return(int64(8), nil)

	total, err := mu.CountDocumentWithPipeline(ctx, aggr, pipeline)
	assert.NoError(t, err)
	assert.Equal(t, int64(8), total)
	mu.AssertExpectations(t)
}

func TestMockObjectIDCreator_ObjectIDFromHex(t *testing.T) {
	moc := new(MockObjectIDCreator)
	oid := primitive.NewObjectID()
	s := oid.Hex()
	moc.On("ObjectIDFromHex", s).Return(oid, nil)

	out, err := moc.ObjectIDFromHex(s)
	assert.NoError(t, err)
	assert.Equal(t, oid, out)
	moc.AssertExpectations(t)
}

func TestMockPasswordUtils(t *testing.T) {
	mp := new(MockPasswordUtils)
	mp.On("GenerateSalt", 8).Return("salt", nil)
	mp.On("HashPassword", "pwd", "salt").Return("hashed", nil)
	mp.On("VerifyPassword", "hashed", "pwd", "salt").Return(true)

	salt, err := mp.GenerateSalt(8)
	assert.NoError(t, err)
	assert.Equal(t, "salt", salt)

	hashed, err := mp.HashPassword("pwd", "salt")
	assert.NoError(t, err)
	assert.Equal(t, "hashed", hashed)

	assert.True(t, mp.VerifyPassword("hashed", "pwd", "salt"))
	mp.AssertExpectations(t)
}
