/*
Go Binding of Tkrzw

DBM (Database Manager) is a concept to store an associative array on a permanent storage.  In other words, DBM allows an application program to store key-value pairs in a file and reuse them later.  Each of keys and values is a string or a sequence of bytes.  A key must be unique within the database and a value is associated to it.  You can retrieve a stored record with its key very quickly.  Thanks to simple structure of DBM, its performance can be extremely high.

Tkrzw is a library implementing DBM with various algorithms.  It features high degrees of performance, concurrency, scalability and durability.  The following data structures are provided.

- HashDBM : File datatabase manager implementation based on hash table.
- TreeDBM : File datatabase manager implementation based on B+ tree.
- SkipDBM : File datatabase manager implementation based on skip list.
- TinyDBM : On-memory datatabase manager implementation based on hash table.
- BabyDBM : On-memory datatabase manager implementation based on B+ tree.
- CacheDBM : On-memory datatabase manager implementation with LRU deletion.
- StdHashDBM : On-memory DBM implementations using std::unordered_map.
- StdTreeDBM : On-memory DBM implementations using std::map.

Whereas Tkrzw is C++ library, this package provides its Go interface.  All above data structures are available via one adapter struct "DBM".  Read the homepage (http://dbmx.net/tkrzw/) for details.

DBM stores key-value pairs of strings.  Each string is represented as a byte array in Go.  Although you can also use methods with string arguments and return values, their internal representations are byte arrays.  Other types such as intergers can also be taken as parameters and they are converted into byte arrays implicitly.

All classes are defined under the package "tkrzw", which can be imported in source files of application programs as "github.com/estraier/tkrzw-go".

 import "github.com/estraier/tkrzw-go"

An instance of the struct "DBM" is used in order to handle a database.  You can store, delete, and retrieve records with the instance.  The result status of each operation is represented by an object of the struct "Status".  Iterator to access each record is implemented by the struct "Iterator".

The key and the value of the records are stored as byte arrays.  However, you can specify strings and other types which imlements the Stringer interface whereby the object is converted into a byte array.

Install the latest version of Tkrzw beforehand.  If you write the above import directive, the Go module for Tkrzw is installed implicitly when you build or run your program.  Go 1.14 or later is required to use this package.

The following code is a simple example to use a database, without checking errors.  Many methods accept both byte arrays and strings.  If strings are given, they are converted implicitly into byte arrays.

 import (
   "fmt"
   "github.com/estraier/tkrzw-go"
 )

 func main() {
   fmt.Println("Hello World")
 }
*/
package tkrzw
