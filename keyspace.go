package main

type Keyspace struct {
	Keyspace_name    string
	Durable_writes   bool
	Strategy_class   string
	Strategy_options string
}

func getKeyspaces(include []string, exclude []string) []Keyspace {
	iter := con.Query(`SELECT
				keyspace_name,
				durable_writes,
				strategy_class,
				strategy_options
			   FROM
				system.schema_keyspaces`).Iter()
	var keyspaces []Keyspace
	var keyspace Keyspace
ksloop:
	for iter.Scan(&keyspace.Keyspace_name, &keyspace.Durable_writes, &keyspace.Strategy_class, &keyspace.Strategy_options) {
		if len(include) > 0 {
			keep := false
			for _, i := range include {
				if i == keyspace.Keyspace_name {
					keep = true
				}
			}
			if !keep {
				continue ksloop
			}
		}
		if len(exclude) > 0 {
			for _, e := range exclude {
				if e == keyspace.Keyspace_name {
					continue ksloop
				}
			}
		}
		keyspaces = append(keyspaces, keyspace)
	}
	FatalIfError("getKeyspaces", iter.Close())
	return keyspaces
}
