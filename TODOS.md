TODOS
=
- [ ] **important** BEFORE you start mining a block, verify it with the current chain
- [x] amounts should be `float`s not `int`s
- [x] more of a SWE thing but blocks should explode into transactions, which the lbm processes as records
- [x] (small) numtransactions as a field for the block is probably not necessary since it's embedded in len(transactions) anyway
- [x] wallets are just public keys lol
- [x] **important**: need to store struct pointers in places like maps to make sure we can properly share memory between routines.
    - further on this, EVEN with interfaces in types (such as map val types) we should store pointers to structs
- [x] make a dedicated field in the block for the miner transaction instead of checking for one miner transaction
- [x] naive reward calculation (constant)
- [ ] transaction fees (should not include the miner wallet), instead we should somehow get this from the block itself
- [x] miner should sign its minertx
- [ ] (maybe?) make the chain talk to the lcm which signals to the miner instead of this buffered chan business
- [ ] **important**: eventually we need to make sure someone doesn't just spam mining transactions lol, no block should have mining transactions in the map of transactions
- [ ] approximate timestamps must be created (>= median of prev. 11 blocks and <= network time + 2 hours)
- [ ] (small) not sure, but should headblocks store block hashes or just a pointer to the block itself? would save some error potential
- [x] transaction ids (`txid`s) need to be orderable for consistent serialization for hashing (change from uint64 to HashType) - just compare in 4 uint64 chunks
- [ ] think about how block orphaning will happen (on a time interval, once a chain grows much longer than others)
- [x] (small) transactions have timestamps 
- [ ] (stretch) serializing / deserializing engine for all structs
    - this functionality could automatically take marshalable structs that you can register protocol numbers for
- [x] concurrency (everywhere)
- [ ] (moderate) transaction load generator with random values (that are both valid and invalid)
- [ ] (extension) what to do if someone sends you an almost valid chain 
- [ ] (super extension) lewatch web client that you can use to watch the chain growing in real time
- [ ] transactions need to spit out miner wallets so that we can verify the block

# Notes
- We can have events being broadcasted when blocks get added or removed, then respond to those events with orphaning, etc. (could either use a chan model or callback)
    - my thinking is that initially we can just never orphan stuff but then eventually orphan when growing two longer. It's important to be able to mempool stuff as a miner
