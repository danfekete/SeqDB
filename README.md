# SeqDB

SeqDB is a *locking* sequential numbering database, which allows you to generate numbers sequentially without repetition and skips. 

The database is built with Go using [BoltDB](https://github.com/boltdb/bolt). It follows a plain text protocol inspired by memcache and redis.

## Architecture

The database is divided into *buckets*. There can be any number of buckets in the database but there must be at least one. Buckets may be divided based on your users or there can be one global bucket. 

Buckets are divided into *sequences*. Each sequence is named similar to a bucket, but there can't be two sequences in a bucket with the same name. The same sequence names in different buckets are allowed.

**Example**

You have a bucket named `First` and another bucket named `Second` . You can have a sequence named `default` in both the `First` and `Second` bucket, but there cannot be two of them in the `First` bucket.

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

