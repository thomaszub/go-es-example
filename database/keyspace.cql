CREATE KEYSPACE IF NOT EXISTS {{.Keyspace}}
WITH REPLICATION = {
    'class': 'SimpleStrategy',
    'replication_factor': 1 
};
