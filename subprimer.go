package archive

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pborman/uuid"
	"net/url"
	"time"
)

type Subprimer struct {
	Id            string                 `json:"id"`
	Created       time.Time              `json:"created"`
	Updated       time.Time              `json:"updated"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Url           string                 `json:"url"`
	PrimerId      string                 `json:"primerId"`
	Crawl         bool                   `json:"crawl"`
	StaleDuration time.Duration          `json:"staleDuration"`
	LastAlertSent *time.Time             `json:"lastAlertSent"`
	Meta          map[string]interface{} `json:"meta"`
	Stats         *SubprimerStats        `json:"stats"`
}

type SubprimerStats struct {
	UrlCount             int `json:"urlCount"`
	ContentUrlCount      int `json:"contentUrlCount"`
	ContentMetadataCount int `json:"contentMetadataCount"`
}

func (s *Subprimer) CalcStats(db sqlQueryExecable) error {
	urlCount, err := s.urlCount(db)
	if err != nil {
		return err
	}

	contentUrlCount, err := s.contentUrlCount(db)
	if err != nil {
		return err
	}

	metadataCount, err := s.contentWithMetadataCount(db)
	if err != nil {
		return err
	}

	s.Stats = &SubprimerStats{
		UrlCount:             urlCount,
		ContentUrlCount:      contentUrlCount,
		ContentMetadataCount: metadataCount,
	}

	// TODO - stop saving here & instead hook this up to some sort of cron task
	return s.Save(db)
}

func (s *Subprimer) urlCount(db sqlQueryable) (count int, err error) {
	err = db.QueryRow("select count(1) from urls where url ilike $1", "%"+s.Url+"%").Scan(&count)
	return
}

func (s *Subprimer) contentUrlCount(db sqlQueryable) (count int, err error) {
	err = db.QueryRow("select count(1) from urls where url ilike $1 and content_sniff != 'text/html; charset=utf-8' and hash != ''", "%"+s.Url+"%").Scan(&count)
	return
}

func (s *Subprimer) contentWithMetadataCount(db sqlQueryable) (count int, err error) {
	err = db.QueryRow("select count(1) from urls where urls.url ilike $1 and urls.content_sniff != 'text/html; charset=utf-8' and exists (select null from metadata where urls.hash = metadata.subject)", "%"+s.Url+"%").Scan(&count)
	return
}

// AsUrl retrieves the url that corresponds for the crawlUrl. If one doesn't exist & the url is saved,
// a new url is created
func (c *Subprimer) AsUrl(db sqlQueryExecable) (*Url, error) {
	// TODO - this assumes http protocol, make moar robust
	addr, err := url.Parse(fmt.Sprintf("http://%s", c.Url))
	if err != nil {
		return nil, err
	}

	u := &Url{Url: addr.String()}
	if err := u.Read(db); err != nil {
		if err == ErrNotFound {
			if err := u.Insert(db); err != nil {
				return u, err
			}
		} else {
			return nil, err
		}
	}

	return u, nil
}

// TODO - this currently doesn't check the status of metadata, gonna need to do that
// UndescribedContent returns a list of content-urls from this subprimer that need work.
func (s *Subprimer) UndescribedContent(db sqlQueryable, limit, offset int) ([]*Url, error) {
	rows, err := db.Query(QSubprimerUndescribedContent, "%"+s.Url+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make([]*Url, limit)
	i := 0
	for rows.Next() {
		u := &Url{}
		if err := u.UnmarshalSQL(rows); err != nil {
			return nil, err
		}
		urls[i] = u
		i++
	}

	return urls[:i], nil
}

// TODO - this currently doesn't check the status of metadata, gonna need to do that
// DescribedContent returns a list of content-urls from this subprimer that need work.
func (s *Subprimer) DescribedContent(db sqlQueryable, limit, offset int) ([]*Url, error) {
	rows, err := db.Query(QSubprimerDescribedContent, "%"+s.Url+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make([]*Url, limit)
	i := 0
	for rows.Next() {
		u := &Url{}
		if err := u.UnmarshalSQL(rows); err != nil {
			return nil, err
		}
		urls[i] = u
		i++
	}

	return urls[:i], nil
}

// func (s *Subprimer) Stats() {
// }

func (c *Subprimer) Read(db sqlQueryable) error {
	if c.Id != "" {
		row := db.QueryRow(fmt.Sprintf("select %s from subprimers where id = $1", subprimerCols()), c.Id)
		return c.UnmarshalSQL(row)
	} else if c.Url != "" {
		row := db.QueryRow(fmt.Sprintf("select %s from subprimers where url = $1", subprimerCols()), c.Url)
		return c.UnmarshalSQL(row)
	}
	return ErrNotFound
}

func (c *Subprimer) Save(db sqlQueryExecable) error {
	prev := &Subprimer{Url: c.Url}
	if err := prev.Read(db); err != nil {
		if err == ErrNotFound {
			c.Id = uuid.New()
			c.Created = time.Now().Round(time.Second)
			c.Updated = c.Created
			_, err := db.Exec(fmt.Sprintf("insert into subprimers (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", subprimerCols()), c.SQLArgs()...)
			return err
		} else {
			return err
		}
	} else {
		c.Updated = time.Now().Round(time.Second)
		_, err := db.Exec("update subprimers set created = $2, updated = $3, title = $4, description = $5, url = $6, primer_id = $7, crawl = $8, stale_duration = $9, last_alert_sent = $10, meta = $11, stats = $12 where id = $1", c.SQLArgs()...)
		return err
	}

	return nil
}

func (c *Subprimer) Delete(db sqlQueryExecable) error {
	_, err := db.Exec("delete from subprimers where url = $1", c.Url)
	return err
}

func (c *Subprimer) UnmarshalSQL(row sqlScannable) error {
	var (
		id, url, pId, title, description string
		created, updated                 time.Time
		lastAlert                        *time.Time
		stale                            int64
		crawl                            bool
		metaBytes, statsBytes            []byte
	)

	if err := row.Scan(&id, &created, &updated, &title, &description, &url, &pId, &crawl, &stale, &lastAlert, &metaBytes, &statsBytes); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	if lastAlert != nil {
		utc := lastAlert.In(time.UTC)
		lastAlert = &utc
	}

	var meta map[string]interface{}
	if metaBytes != nil {
		if err := json.Unmarshal(metaBytes, &meta); err != nil {
			return err
		}
	}

	stats := &SubprimerStats{}
	if statsBytes != nil {
		if err := json.Unmarshal(statsBytes, stats); err != nil {
			return err
		}
	}

	*c = Subprimer{
		Id:            id,
		Created:       created.In(time.UTC),
		Updated:       updated.In(time.UTC),
		Title:         title,
		Description:   description,
		Url:           url,
		PrimerId:      pId,
		Crawl:         crawl,
		StaleDuration: time.Duration(stale * 1000000),
		LastAlertSent: lastAlert,
		Meta:          meta,
		Stats:         stats,
	}

	return nil
}

func subprimerCols() string {
	return "id, created, updated, title, description, url, primer_id, crawl, stale_duration, last_alert_sent, meta, stats"
}

func (c *Subprimer) SQLArgs() []interface{} {
	date := c.LastAlertSent
	if date != nil {
		utc := date.In(time.UTC)
		date = &utc
	}

	metaBytes, err := json.Marshal(c.Meta)
	if err != nil {
		panic(err)
	}

	statBytes, err := json.Marshal(c.Stats)
	if err != nil {
		panic(err)
	}

	return []interface{}{
		c.Id,
		c.Created.In(time.UTC),
		c.Updated.In(time.UTC),
		c.Title,
		c.Description,
		c.Url,
		c.PrimerId,
		c.Crawl,
		c.StaleDuration / 1000000,
		date,
		metaBytes,
		statBytes,
	}
}
