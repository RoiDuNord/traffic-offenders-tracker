package repo

type Repo interface {
	Connect(conn string) error
	Insert(key string, value []byte) (int, error)
	Close() error
}

type RepoService struct {
	repo Repo
}

func New(repo Repo) *RepoService {
	return &RepoService{repo: repo}
}

func (rm *RepoService) InsertMessage(key string, value []byte) (int, error) {
	return rm.repo.Insert(key, value)
}

func (rm *RepoService) Close() error {
	return rm.repo.Close()
}
