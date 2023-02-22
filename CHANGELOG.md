# v4.0.2

- (internal) Update postgres adapter

# v4.0.1

- (bug) Fix slice filter function for postgres

# v4.0.0

- (bug) Cleanup tests (postgres)
- (bug) Fix connection concurrency issues (postgres)
- (bug) Fix connection transaction issues (postgres)
- (bug) Use connection pool (postgres)
- (bug) Create table if not exists on queries and upserts (postgres)

# v4.0.0-alpha.2

- (bug) correct version tagging

# v4.0.0-alpha.1

- (bc) Rework filters (In / NotIn / InArray / NotInArray)
- (bc) Rename `FieldTypeString` to `FieldTypeText`
- (feature) Add postgres support

# v3.1.0

- (feature) Add NotEquals filter
- (feature) Add InArray filter
- (feature) Add NotInArray filter

# v3.0.1

- (bug) Remove module v2 import paths

# v3.0.0

- (bc) Add config struct and default timeouts

# v2.0.2

- (feature) Add database delete operaton for mongodb

# v2.0.1

- (bug) Make operator types public

# v2.0.0

- (bc) All new application programming interface
- (bc) Drops support for Meilisearch
- (feature) Adds support for index CRUD operations

# v1.2.0

- (feature) Add support for collection listing

# v1.1.2

- (bug) Dropping collection as part of rename now works

# v1.1.1

- (feature) Collection renaming for mongodb
- (feature) Collection dropping for mongodb
- (internal) Same as v1.1.0

# v1.1.0

- (feature) Collection renaming for mongodb
- (feature) Collection dropping for mongodb

# v1.0.0

- (bc) Move code out of src folder for clean package imports
- (feature) Transaction support for mongodb

# v0.4.1

- (bug) Fix package naming

# v0.4.0

- (feature) Allow sorting in queries
