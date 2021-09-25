gbkr-com/id3
===

An implementation of the ID3 decision tree algorithm, which learns from CSV
conformant data.

The code is organised as:

* `views.go` provides an interface and implementations for ID3 to inspect CSV data
* `decisions.go` defines the internal representation of the decision tree, including writing and reading that tree as JSON
* `learn.go` is the ID3 algorithm itself.