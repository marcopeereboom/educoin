# educoin
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

**educoin is a toy blockchain implementation.**

This blockchain implementation is for education purposes only. It deliberately
leaves out complex details that make up a real blockchain. Code is simple
enough to be studied and understood in a few hours.

The lessons are from simple to more complex and are prefixed with a lesson
number (e.g. `0_blockchain`).

All code is written in go and uses the go unit test framework in order to
exercise it. For example:
```
$ cd 0_blockchain/
$ go test -v      
=== RUN   TestSuccess
--- PASS: TestSuccess (0.00s)
    blockchain_test.go:26: ================================================================================
    blockchain_test.go:27: Height           : 0
    blockchain_test.go:11: Timestamp        : 2018-09-30 17:17:53 -0500 CDT
    blockchain_test.go:12: PreviousBlockHash: 0000000000000000000000000000000000000000000000000000000000000000
    blockchain_test.go:13: Hash             : 6f8b654443d0ea969928d72e6293b692a6c1699cd568af50be285990ccded0c5
    blockchain_test.go:14: Data             : Decred is money!
    blockchain_test.go:33: Block valid      : true
    blockchain_test.go:26: ================================================================================
    blockchain_test.go:27: Height           : 1
    blockchain_test.go:11: Timestamp        : 2018-09-30 17:17:53 -0500 CDT
    blockchain_test.go:12: PreviousBlockHash: 6f8b654443d0ea969928d72e6293b692a6c1699cd568af50be285990ccded0c5
    blockchain_test.go:13: Hash             : b5fd8f68d96badfd04105a9938a48772b4ff19d4e8fbc9918c1f38eb22ba6735
    blockchain_test.go:14: Data             : Send 1 Decred to Alice
    blockchain_test.go:33: Block valid      : true
    blockchain_test.go:26: ================================================================================
    blockchain_test.go:27: Height           : 2
    blockchain_test.go:11: Timestamp        : 2018-09-30 17:17:53 -0500 CDT
    blockchain_test.go:12: PreviousBlockHash: b5fd8f68d96badfd04105a9938a48772b4ff19d4e8fbc9918c1f38eb22ba6735
    blockchain_test.go:13: Hash             : f5627bd22bc68ec22e3f32388d772b1b6806d4fa1bb13c6ba1195a54decbd7ec
    blockchain_test.go:14: Data             : Send 2 Decred to Bob
    blockchain_test.go:33: Block valid      : true
=== RUN   TestFailure
--- PASS: TestFailure (0.00s)
    blockchain_test.go:26: ================================================================================
    blockchain_test.go:27: Height           : 0
    blockchain_test.go:11: Timestamp        : 2018-09-30 17:17:53 -0500 CDT
    blockchain_test.go:12: PreviousBlockHash: 0000000000000000000000000000000000000000000000000000000000000000
    blockchain_test.go:13: Hash             : 6f8b654443d0ea969928d72e6293b692a6c1699cd568af50be285990ccded0c5
    blockchain_test.go:14: Data             : Decred is money!
    blockchain_test.go:33: Block valid      : true
    blockchain_test.go:26: ================================================================================
    blockchain_test.go:27: Height           : 1
    blockchain_test.go:11: Timestamp        : 2018-09-30 17:17:53 -0500 CDT
    blockchain_test.go:12: PreviousBlockHash: 6f8b654443d0ea969928d72e6293b692a6c1699cd568af50be285990ccded0c5
    blockchain_test.go:13: Hash             : b5fd8f68d96badfd04105a9938a48772b4ff19d4e8fbc9918c1f38eb22ba6735
    blockchain_test.go:14: Data             : Send 2 Decred to Alice
    blockchain_test.go:33: Block valid      : false
PASS
ok      github.com/marcopeereboom/educoin/0_blockchain  0.001s
```

Patches and comments are welcome!
