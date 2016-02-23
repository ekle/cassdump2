package main

import "sort"

type Table struct {
	Keyspace_name          string
	Table_name             string
	Bloom_filter_fp_chance float64
	Caching                string
	Comment                string
	/*	Compaction                  map[string]string
		Compression                 map[string]string
		Crc_check_chance            float64
		Dclocal_read_repair_chance  float64
		Default_time_to_live        int32
		Extensions                  map[string][]byte
		Flags                       []string
		Gc_grace_seconds            int32
		Id                          gocql.UUID
		Max_index_interval          int32
		Memtable_flush_period_in_ms int32
		Min_index_interval          int32
		Read_repair_chance          float64
		Speculative_retry           string*/
	Columns []Column
}

type Column struct {
	Keyspace_name   string
	Table_name      string
	Column_name     string
	Componont_index int
	Index_name      string
	Index_options   string
	Index_type      string
	Type            string
	Validator       string
}

type ColumnByPosition []Column

func (a ColumnByPosition) Len() int           { return len(a) }
func (a ColumnByPosition) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ColumnByPosition) Less(i, j int) bool { return a[i].Componont_index < a[j].Componont_index }

func getTables(keyspace Keyspace) []Table {
	// TODO: add include and exlude filters
	iter := con.Query(`SELECT
				keyspace_name,
				columnfamily_name,
				bloom_filter_fp_chance,
				caching,
				comment`+ /*,
		compaction,
		compression,
		crc_check_chance,
		dclocal_read_repair_chance,
		default_time_to_live,
		extensions,
		flags,
		gc_grace_seconds,
		id,
		max_index_interval,
		memtable_flush_period_in_ms,
		min_index_interval,
		read_repair_chance,
		speculative_retry*/`
			FROM system.schema_columnfamilies
			WHERE keyspace_name = ?`, keyspace.Keyspace_name).Iter()
	var tables []Table
	var t Table
	for iter.Scan(
		&t.Keyspace_name,
		&t.Table_name,
		&t.Bloom_filter_fp_chance,
		&t.Caching,
		&t.Comment) {
		/*
			&t.Compaction,
			&t.Compression,
			&t.Crc_check_chance,
			&t.Dclocal_read_repair_chance,
			&t.Default_time_to_live,
			&t.Extensions,
			&t.Flags,
			&t.Gc_grace_seconds,
			&t.Id,
			&t.Max_index_interval,
			&t.Memtable_flush_period_in_ms,
			&t.Min_index_interval,
			&t.Read_repair_chance,
			&t.Speculative_retry
			) {*/
		// ---
		tables = append(tables, t)
	}
	FatalIfError("getTables", iter.Close())
	for k, table := range tables {
		iter := con.Query(`SELECT
					keyspace_name,
					columnfamily_name,
					column_name,
					component_index,
					index_name,
					index_options,
					index_type,
					type,
					validator
				FROM system.schema_columns
				WHERE keyspace_name = ? AND columnfamily_name = ?`, keyspace.Keyspace_name, table.Table_name).Iter()
		var c Column
		for iter.Scan(&c.Keyspace_name, &c.Table_name, &c.Column_name, &c.Componont_index, &c.Index_name, &c.Index_options, &c.Index_type, &c.Type, &c.Validator) {
			table.Columns = append(table.Columns, c)
		}
		FatalIfError("getTables.getColumns", iter.Close())
		sort.Sort(ColumnByPosition(table.Columns))
		tables[k] = table
	}
	return tables
}
