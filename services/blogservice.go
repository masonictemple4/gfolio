package services

import (
	"context"

	"github.com/masonictemple4/masonictempl/db"
	"github.com/masonictemple4/masonictempl/db/models"
)

type BlogService struct {
	Store *db.BlogStore
}

func NewBlogService() *BlogService {
	sDB := db.NewSqliteDB("blog.db", nil)
	store, err := db.NewBlogStore(db.WithDB(sDB))
	if err != nil {
		panic(err)
	}
	return &BlogService{Store: store}
}

func (b *BlogService) List(ctx context.Context) []models.Blog {
	blogs := make([]models.Blog, 0)
	order := "created_at desc"
	if err := b.Store.ListBlogs(&blogs, nil, nil, order, "Authors", "Tags"); err != nil {
		// TODO: Log the error
		return blogs
	}

	return blogs
}

func (b *BlogService) GetWithSlug(ctx context.Context, slug string, preloads ...string) (*models.Blog, error) {
	var blog models.Blog
	if err := b.Store.FindBySlug(&blog, slug, preloads...); err != nil {
		return nil, err
	}
	return &blog, nil
}

/*
var ErrUnknownIncrementType error = errors.New("unknown increment type")

func NewCount(log *slog.Logger, cs *db.CountStore) Count {
	return Count{
		Log:        log,
		CountStore: cs,
	}
}

type Count struct {
	Log        *slog.Logger
	CountStore *db.CountStore
}

func (cs Count) Increment(ctx context.Context, it IncrementType, sessionID string) (counts Counts, err error) {
	// Work out which operations to do.
	var global, session func(ctx context.Context, id string) (count int, err error)
	switch it {
	case IncrementTypeGlobal:
		global = cs.CountStore.Increment
		session = cs.CountStore.Get
	case IncrementTypeSession:
		global = cs.CountStore.Get
		session = cs.CountStore.Increment
	default:
		return counts, ErrUnknownIncrementType
	}

	// Run the operations in parallel.
	var wg sync.WaitGroup
	wg.Add(2)
	errs := make([]error, 2)
	go func() {
		defer wg.Done()
		counts.Global, errs[0] = global(ctx, "global")
	}()
	go func() {
		defer wg.Done()
		counts.Session, errs[1] = session(ctx, sessionID)
	}()
	wg.Wait()

	return counts, errors.Join(errs...)
}

func (cs Count) Get(ctx context.Context, sessionID string) (counts Counts, err error) {
	globalAndSessionCounts, err := cs.CountStore.BatchGet(ctx, "global", sessionID)
	if err != nil {
		err = fmt.Errorf("countservice: failed to get counts: %w", err)
		return
	}
	if len(globalAndSessionCounts) != 2 {
		err = fmt.Errorf("countservice: unexpected counts returned, expected 2, got %d", len(globalAndSessionCounts))
		return
	}
	counts.Global = globalAndSessionCounts[0]
	counts.Session = globalAndSessionCounts[1]
	return
}


*/
