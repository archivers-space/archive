package archive

const qCollectionInsert = `
insert into collections 
  (id, created, updated, creator, title, schema, contents ) 
values ($1, $2, $3, $4, $5, $6, $7);`

const qCollectionUpdate = `
update collections 
set created=$2, updated=$3, creator=$4, title=$5, schema=$6, contents=$7 
where id = $1;`

const qCollectionById = `
select 
  id, created, updated, creator, title, schema, contents 
from collections 
where id = $1;`

const qCollectionDelete = `
delete from collections 
where id = $1;`

const qCollections = `
SELECT
  id, created, updated, creator, title, schema, contents
FROM collections 
ORDER BY created DESC 
LIMIT $1 OFFSET $2;`

const qMetadataLatest = `
select
  hash, time_stamp, key_id, subject, prev, meta 
from metadata 
where 
  key_id = $1 and 
  subject = $2 
order by time_stamp desc;`

const qMetadataForSubject = `
select
  hash, time_stamp, key_id, subject, prev, meta  
from metadata
where 
  subject = $1 and 
  deleted = false and 
  meta is not null;`

const qMetadataInsert = `
insert into metadata
  (hash, time_stamp, key_id, subject, prev, meta, deleted)
values 
  ($1, $2, $3, $4, $5, $6, false);`

const qPrimerById = `
select
  id, created, updated, short_title, title, description
from primers 
where 
  id = $1;`

const qPrimerInsert = `
insert into primers
  (id, created, updated, short_title, title, description)
values
  ($1, $2, $3, $4, $5, $6);`

const qPrimerUpdate = `
update primers set
  created = $2, updated = $3, short_title = $4, title = $5, description = $6
where
  id = $1;`

const qPrimerDelete = `
delete from primers 
where id = $1;`

const qPrimerSubprimers = `
select
  id, created, updated, title, description, url, primer_id, crawl, stale_duration,
  last_alert_sent, meta, stats
from subprimers
where 
  primer_id = $1;`

const qPrimersCrawling = `
select
  id, created, updated, short_title, title, description
from primers 
where
  crawl = true 
limit $1 offset $2;`

const qPrimersList = `
select
  id, created, updated, short_title, title, description 
from primers
order by created desc
limit $1 offset $2;`

const qSubprimerById = `
select
  id, created, updated, title, description, url, primer_id, crawl, stale_duration,
  last_alert_sent, meta, stats
from subprimers 
where 
  id = $1;`

const qSubprimerByUrl = `
select
  id, created, updated, title, description, url, primer_id, crawl, stale_duration,
  last_alert_sent, meta, stats
from subprimers 
where 
  url = $1;`

const qSubprimerInsert = `
insert into subprimers 
  (id, created, updated, title, description, url, primer_id, crawl, stale_duration,
   last_alert_sent, meta, stats) 
values 
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

const qSubprimerUpdate = `
update subprimers 
set 
  created = $2, updated = $3, title = $4, description = $5, url = $6, primer_id = $7, 
  crawl = $8, stale_duration = $9, last_alert_sent = $10, meta = $11, stats = $12
where
  id = $1;`

const qSubprimerDelete = `
delete from subprimers 
where 
  url = $1;`

const qSubprimerUrlCount = `
select count(1) 
from urls 
where 
  url ilike $1;`

const qSubprimerCrawlingUrls = `
select
  id, created, updated, title, description, url, primer_id, crawl, stale_duration,
  last_alert_sent, meta, stats
from subprimers 
where 
  crawl = true;`

const qSubprimerContentUrlCount = `
select count(1) 
from urls 
where 
  url ilike $1 and 
  content_sniff != 'text/html; charset=utf-8' 
  and hash != '';`

const qSubprimerContentWithMetadataCount = `
select count(1)
from urls 
where 
  urls.url ilike $1 and 
  urls.content_sniff != 'text/html; charset=utf-8' 
  and exists (select null from metadata where urls.hash = metadata.subject);`

const qSubprimerUndescribedContentUrls = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff, content_length, 
  title, id, headers_took, download_took, headers, meta, hash 
from urls 
where 
  url ilike $1
  and content_sniff != 'text/html; charset=utf-8'
  and last_get is not null
  -- confirm is not empty hash
  and hash != '1220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855'
  and not exists (select null from metadata where urls.hash = metadata.subject) 
limit $2 offset $3;`

const qSubprimerDescribedContentUrls = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff, content_length, 
  title, id, headers_took, download_took, headers, meta, hash 
from urls 
where 
  url ilike $1
  and content_sniff != 'text/html; charset=utf-8'
  and last_get is not null
  -- confirm is not empty hash
  and hash != '1220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855'
  and exists (select null from metadata where urls.hash = metadata.subject) 
limit $2 offset $3;`

const qSnapshotsByUrl = `
select
  url, created, status, duration, meta, hash
from snapshots 
where 
  url = $1;`

const qSnapshotInsert = `
insert into snapshots 
  (url, created, status, duration, meta, hash)
values 
  ($1, $2, $3, $4, $5, $6);`

const qUrlsSearch = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, title, id, headers_took, download_took, headers, meta, hash
from urls 
where 
  url ilike $1 
limit $2 offset $3;`

const qUrlsList = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, title, id, headers_took, download_took, headers, meta, hash
from urls 
order by created desc 
limit $1 offset $2;`

const qUrlsFetched = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, title, id, headers_took, download_took, headers, meta, hash
from urls 
where
  last_get is not null 
order by created desc
limit $1 offset $2;`

const qUrlsUnfetched = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, title, id, headers_took, download_took, headers, meta, hash
from urls
where 
  last_get is null 
order by created desc 
limit $1 offset $2;`

const qUrlsForHash = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, title, id, headers_took, download_took, headers, meta, hash
from urls
where 
  hash = $1;`

const qUrlByUrlString = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, title, id, headers_took, download_took, headers, meta, hash
from urls 
where
  url = $1;`

const qUrlInsert = `
insert into urls
  (url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, title, id, headers_took, download_took, headers, meta, hash)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);`

const qUrlUpdate = `
update urls 
set 
  created=$2, updated=$3, last_head=$4, last_get=$5, status=$6, content_type=$7, content_sniff=$8,
  content_length=$9, title=$10, id=$11, headers_took=$12, download_took=$13, headers=$14, meta=$15, hash=$16 
where 
  url = $1;`

const qUrlDelete = `
delete from urls 
where
  url = $1;`

const qUrlInboundLinkUrlStrings = `
select src 
from links 
where
  dst = $1;`

const qUrlOutboundLinkUrlStrings = `
select dst 
from links 
where
  src = $1;`

const qUrlDstLinks = `
select 
  urls.url, urls.created, urls.updated, last_head, last_get, status, content_type, content_sniff, 
  content_length, title, id, headers_took, download_took, headers, meta, hash 
from urls, links
where 
  links.src = $1 and 
  links.dst = urls.url;`

const qUrlSrcLinks = `
select
  urls.url, urls.created, urls.updated, last_head, last_get, status, content_type, content_sniff, 
  content_length, title, id, headers_took, download_took, headers, meta, hash 
from urls, links 
where 
  links.dst = $1 and 
  links.src = urls.url;`
