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

All identifiers are defined under the package "tkrzw", which can be imported in source files of application programs as "github.com/estraier/tkrzw-go".

 import "github.com/estraier/tkrzw-go"

An instance of the struct "DBM" is used in order to handle a database.  You can store, delete, and retrieve records with the instance.  The result status of each operation is represented by an object of the struct "Status".  Iterator to access each record is implemented by the struct "Iterator".

The key and the value of the records are stored as byte arrays.  However, you can specify strings and other types which imlements the Stringer interface whereby the object is converted into a byte array.

Install the latest version of Tkrzw beforehand.  If you write the above import directive, the Go module for Tkrzw is installed implicitly when you build or run your program.  Go 1.14 or later is required to use this package.

The following code is a simple example to use a database, without checking errors.  Many methods accept both byte arrays and strings.  If strings are given, they are converted implicitly into byte arrays.

 package main

 import (
   "fmt"
   "github.com/estraier/tkrzw-go"
 )

 func main() {
   // Prepares the database.
   dbm := tkrzw.NewDBM()
   dbm.Open("casket.tkh", true, "truncate=true,num_buckets=100")

   // Sets records.
   // Keys and values are implicitly converted into bytes.
   dbm.Set("first", "hop", true)
   dbm.Set("second", "step", true)
   dbm.Set("third", "jump", true)

   // Retrieves record values as strings.
   fmt.Println(dbm.GetStrSimple("first", "*"))
   fmt.Println(dbm.GetStrSimple("second", "*"))
   fmt.Println(dbm.GetStrSimple("third", "*"))

   // Traverses records with a range over a channel.
   for record := range dbm.EachStr() {
     fmt.Println(record.Key, record.Value)
   }

   // Closes the database.
   dbm.Close()
 }

The following code is an advanced example where a so-called long transaction is done by the compare-and-exchange (aka compare-and-swap) idiom.  The example also shows how to use the iterator to access each record.

 package main

 import (
   "fmt"
   "github.com/estraier/tkrzw-go"
 )

 func main() {
   // Prepares the database.
   // The method OrDie causes panic if the status is not success.
   // You should write your own error handling in large scale programs.
   dbm := tkrzw.NewDBM()
   dbm.Open("casket.tkt", true, "truncate=true,num_buckets=100").OrDie()

   // Closes the database for sure and checks the error too.
   defer func() { dbm.Close().OrDie() }()

   // Two bank accounts for Bob and Alice.
   // Numeric values are converted into strings implicitly.
   dbm.Set("Bob", 1000, false).OrDie()
   dbm.Set("Alice", 3000, false).OrDie()

   // Function to do a money transfer atomically.
   transfer := func(src_key string, dest_key string, amount int64) *tkrzw.Status {
     // Gets the old values as numbers.
     old_src_value := tkrzw.ToInt(dbm.GetStrSimple(src_key, "0"))
     old_dest_value := tkrzw.ToInt(dbm.GetStrSimple(dest_key, "0"))

     // Calculates the new values.
     new_src_value := old_src_value - amount
     new_dest_value := old_dest_value + amount
     if new_src_value < 0 {
       return tkrzw.NewStatus(tkrzw.StatusApplicationError, "insufficient value")
     }

     // Prepares the pre-condition and the post-condition of the transaction.
     old_records := []tkrzw.KeyValueStrPair{
       {src_key, tkrzw.ToString(old_src_value)},
       {dest_key, tkrzw.ToString(old_dest_value)},
     }
     new_records := []tkrzw.KeyValueStrPair{
       {src_key, tkrzw.ToString(new_src_value)},
       {dest_key, tkrzw.ToString(new_dest_value)},
     }

     // Performs the transaction atomically.
     // This fails safely if other concurrent transactions break the pre-condition.
     return dbm.CompareExchangeMultiStr(old_records, new_records)
   }

   // Tries a transaction until it succeeds
   var status *tkrzw.Status
   for num_tries := 0; num_tries < 100; num_tries++ {
     status = transfer("Alice", "Bob", 500)
     if !status.Equals(tkrzw.StatusInfeasibleError) {
       break
     }
   }
   status.OrDie()

   // Traverses records in a primitive way.
   iter := dbm.MakeIterator()
   defer iter.Destruct()
   iter.First()
   for {
     key, value, status := iter.GetStr()
     if !status.IsOK() {
       break
     }
     fmt.Println(key, value)
     iter.Next()
   }
 }
*/
package tkrzw

// END OF FILE
