# cassdump2
cassandra2.2 dump to cql (experimental)

it uses the JSON export/import feature of > cassandra2.2

backport of cassdump3 for data only dumps.

table and type definition are missing/not correct

use https://github.com/ekle/cassdump3 for cassandra3

feel free to provide pull requests!

this tool is not intended for dumping big tables, more for developing and quickly reverting small databases.
it uses a simple `SELECT JSON * FROM ...` which does not work well with big tables, especially when you have many tombstones.
