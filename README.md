# GCEscape

This package provides series of collection, which escape from gc. It can provide much shorter gc pause when you have
lots of objects in your memory. In benchmark, the gc pause is 99+% less than using std collection. The more objects you
have, the more gc pause time you save

As the arena proposal is delayed, I think it's necessary to have a self-made gc escape collection.

Limited by my time, this repo still has lots of things to do, even the ut need to be enriched.

## how to use?

```shell
go get github.com/hyphennn/gcescape@latest
```
