package db

import (
	"errors"

	"github.com/masonictemple4/masonictempl/db/models"
	"gorm.io/gorm"
)

var (
	ErrNoDB = errors.New("db: no database provided")
)

type BlogStore struct {
	db *gorm.DB
}

type OptionsFn func(*BlogStore)

/*
TODO: Was going to set this up with a pattern like:

	func WithSQLite(name string) func(*BlogStore) {
		return func(bs *BlogStore) {
			bs.db = NewSqliteDB(name, nil)
		}
	}

However, that would require a new options fn type
to then pass into the NewBlogStore function. This
felt cleaner.

Usage:
store := NewSqliteDB("blog.db", nil)
newStore := NewBlogStore(store)
*/
func WithDB(gDB *gorm.DB) func(*BlogStore) {
	return func(bs *BlogStore) {
		bs.db = gDB
	}
}

func NewBlogStore(opts ...OptionsFn) (store *BlogStore, err error) {
	for _, o := range opts {
		o(store)
	}

	if store.db == nil {
		err = ErrNoDB
		store = nil
	}

	err = store.migrate()

	return
}

func (bs *BlogStore) Close() error {
	db, err := bs.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (bs *BlogStore) migrate() error {
	return bs.db.AutoMigrate(
		&models.Tag{},
		&models.Media{},
		&models.User{},
		&models.Comment{},
		&models.Blog{},
	)
}

func (bs *BlogStore) ListBlogs(b *[]models.Blog, query map[string]any, limits map[string]int, preloads ...string) error {
	tx := bs.db.Model(b)

	if len(preloads) > 0 {
		for _, preload := range preloads {
			tx = tx.Preload(preload)
		}
	}

	if limits != nil {
		limit, limitOk := limits["limit"]
		if limitOk {
			tx = tx.Limit(limit)
		}
		offset, offsetOk := limits["offset"]
		if offsetOk {
			tx = tx.Offset(offset)
		}
	}

	if query != nil {
		for k, v := range query {
			tx = tx.Where(k, v)
		}
	}

	return tx.Find(b).Error
}

/*
type OptionsFunc func(*CountStore)

func WithClient(client *dynamodb.Client) func(*CountStore) {
	return func(ms *CountStore) {
		ms.db = client
	}
}

func NewCountStore(tableName, region string, options ...OptionsFunc) (s *CountStore, err error) {
	s = &CountStore{
		tableName: tableName,
	}
	for _, o := range options {
		o(s)
	}
	if s.db == nil {
		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
		if err != nil {
			return s, err
		}
		s.db = dynamodb.NewFromConfig(cfg)
	}
	return
}

type CountStore struct {
	db        *dynamodb.Client
	tableName string
}

func stripEmpty(strings []string) (op []string) {
	for _, s := range strings {
		if s != "" {
			op = append(op, s)
		}
	}
	return
}

type countRecord struct {
	PK    string `dynamodbav:"_pk"`
	Count int    `dynamodbav:"count"`
}

func (s CountStore) BatchGet(ctx context.Context, ids ...string) (counts []int, err error) {
	nonEmptyIDs := stripEmpty(ids)
	if len(nonEmptyIDs) == 0 {
		return nil, nil
	}

	// Make DynamoDB keys.
	ris := make(map[string]types.KeysAndAttributes)
	for _, id := range nonEmptyIDs {
		ri := ris[s.tableName]
		ri.Keys = append(ris[s.tableName].Keys, map[string]types.AttributeValue{
			"_pk": &types.AttributeValueMemberS{
				Value: id,
			},
		})
		ri.ConsistentRead = aws.Bool(true)
		ris[s.tableName] = ri
	}

	// Execute the batch request.
	var batchResponses []map[string]types.AttributeValue

	// DynamoDB might not process everything, so we need a loop.
	var unprocessedAttempts int
	for {
		var bgio *dynamodb.BatchGetItemOutput
		bgio, err = s.db.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: ris,
		})
		if err != nil {
			return
		}
		for _, responses := range bgio.Responses {
			batchResponses = append(batchResponses, responses...)
		}
		if len(bgio.UnprocessedKeys) > 0 {
			ris = bgio.UnprocessedKeys
			unprocessedAttempts++
			if unprocessedAttempts > 3 {
				err = fmt.Errorf("countstore: exceeded three attempts to get all counts")
				return
			}
			continue
		}
		break
	}

	// Process the responses into structs.
	crs := []countRecord{}
	err = attributevalue.UnmarshalListOfMaps(batchResponses, &crs)
	if err != nil {
		err = fmt.Errorf("countstore: failed to unmarshal result of BatchGet: %w", err)
		return
	}

	// Match up the inputs to the records.
	idToCount := make(map[string]int, len(ids))
	for _, cr := range crs {
		idToCount[cr.PK] = cr.Count
	}

	// Create the output in the right order.
	// Missing values are defaulted to zero.
	for _, id := range ids {
		counts = append(counts, idToCount[id])
	}

	return
}

func (s CountStore) Get(ctx context.Context, id string) (count int, err error) {
	if id == "" {
		return
	}
	gio, err := s.db.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"_pk": &types.AttributeValueMemberS{
				Value: id,
			},
		},
		TableName:      &s.tableName,
		ConsistentRead: aws.Bool(true),
	})
	if err != nil || gio.Item == nil {
		return
	}

	var cr countRecord
	err = attributevalue.UnmarshalMap(gio.Item, &cr)
	if err != nil {
		return 0, fmt.Errorf("countstore: failed to process result of Get: %w", err)
	}
	count = cr.Count

	return
}

func (s CountStore) Increment(ctx context.Context, id string) (count int, err error) {
	if id == "" {
		return
	}
	uio, err := s.db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"_pk": &types.AttributeValueMemberS{
				Value: id,
			},
		},
		TableName:        &s.tableName,
		UpdateExpression: aws.String("SET #c = if_not_exists(#c, :zero) + :one"),
		ExpressionAttributeNames: map[string]string{
			"#c": "count",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":zero": &types.AttributeValueMemberN{Value: "0"},
			":one":  &types.AttributeValueMemberN{Value: "1"},
		},
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		return
	}

	// Parse the response.
	var cr countRecord
	err = attributevalue.UnmarshalMap(uio.Attributes, &cr)
	if err != nil {
		return 0, fmt.Errorf("countstore: failed to process result of Increment: %w", err)
	}
	count = cr.Count

	return
}
*/
