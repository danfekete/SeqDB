# [WIP] SeqDB

SeqDB is a distributed fault-tolerant *locking* sequential numbering database, which allows you to generate numbers sequentially without repetition and skips. 

The database is built with Go using [BoltDB](https://github.com/boltdb/bolt). It follows a plain text protocol inspired by memcache and redis.

## Architecture

The database is divided into *buckets*. There can be any number of buckets in the database but there must be at least one. Buckets may be divided based on your users or there can be one global bucket. 

Buckets are divided into *sequences*. Each sequence is named similar to a bucket, but there can't be two sequences in a bucket with the same name. The same sequence names in different buckets are allowed.

**Example**

You have a bucket named `First` and another bucket named `Second` . You can have a sequence named `default` in both the `First` and `Second` bucket, but there cannot be two of them in the `First` bucket.

## Distributed

SeqDB supports *shards* and *nodes* inside the shards. There can be any number of nodes or shard in the cluster. Shards must contain at least one node. In the cluster there must be at least one shard.

### Replication

The nodes inside a shard contain the same data. The nodes can be on the same server or on the other side of the world, SeqDB will make sure that the data is consistent in every node. If a node fails another node will become the master automatically. It is recommended to have at least 5 nodes in every shard, this was if two of them fail, the other three will be able to continue the service. 

For replication SeqDB uses the Raft consensus algorithm

### Sharding

Every shard may contain any number of nodes but at least one. While the data are consistent inside a shard, different shards contain different data.

### Service discovery

SeqDB can automatically register services on the same network without any configuration (besides setting the IP and port for the host and the shard ID)



## Locking

When working with a bucket, either reading or writing, the whole bucket is locked from concurrent access regardless of the sequence. 

## Protocol

### SET bucket sequence value

Directly set the sequence in the bucket

Returns `OK` if command successful



### GET bucket sequence

Return the current sequence



### INC bucket sequence

Increments the sequence by one and return the new value



### QUIT

Closes connection to the server

Returns `OK` if command successful