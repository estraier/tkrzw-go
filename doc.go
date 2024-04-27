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

The following types are mainly used.

- tkrzw.Status : Status of operations
- tkrzw.DBM : Polymorphic database manager
- tkrzw.Iterator : Iterator for each record
- tkrzw.Future : Future containing a status object and extra data
- tkrzw.AsyncDBM : Asynchronous database manager adapter
- tkrzw.File : Generic file implementation
- tkrzw.Index : Secondary index
- tkrzw.IndexIterator : Iterator for each record of the secondary index

An instance of the struct "DBM" is used in order to handle a database.  You can store, delete, and retrieve records with the instance.  The result status of each operation is represented by an object of the struct "Status".  Iterator to access each record is implemented by the struct "Iterator".

The key and the value of the records are stored as byte arrays.  However, you can specify strings and other types which imlements the Stringer interface whereby the object is converted into a byte array.

Install the latest version of Tkrzw beforehand.  If you write the above import directive and prepare the "go.mod" file, the Go module for Tkrzw is installed implicitly when you run "go get".  Go 1.14 or later is required to use this package.

The following code is a simple example to use a database, without checking errors.  Many methods accept both byte arrays and strings.  If strings are given, they are converted implicitly into byte arrays.

 package main

 import (
   "fmt"
   "github.com/estraier/tkrzw-go"
 )

 func main() {
   // Prepares the database.
   dbm := tkrzw.NewDBM()
   dbm.Open("casket.tkh", true,
     tkrzw.ParseParams("truncate=true,num_buckets=100"))

   // Sets records.
   // Keys and values are implicitly converted into bytes.
   dbm.Set("first", "hop", true)
   dbm.Set("second", "step", true)
   dbm.Set("third", "jump", true)

   // Retrieves record values as strings.
   fmt.Println(dbm.GetStrSimple("first", "*"))
   fmt.Println(dbm.GetStrSimple("second", "*"))
   fmt.Println(dbm.GetStrSimple("third", "*"))

   // Checks and deletes a record.
   if dbm.Check("first") {
     dbm.Remove("first")
   }

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
   dbm.Open("casket.tkt", true,
     tkrzw.ParseParams("truncate=true,num_buckets=100")).OrDie()

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

The following code is a typical example of the asynchronous API.  The AsyncDBM class manages a thread pool and handles database operations in the background in parallel.  Each Method of AsyncDBM returns a Future object to monitor the result.

 package main

 import (
   "fmt"
   "github.com/estraier/tkrzw-go"
 )

 func main() {
   // Prepares the database.
   dbm := tkrzw.NewDBM()
   dbm.Open("casket.tkt", true,
     tkrzw.ParseParams("truncate=true,num_buckets=100"))
   defer dbm.Close()

   // Prepares the asynchronous adapter with 4 worker threads.
   async := tkrzw.NewAsyncDBM(dbm, 4)
   defer async.Destruct()

   // Executes the Set method asynchronously.
   future := async.Set("hello", "world", true)
   // Does something in the foreground.
   fmt.Println("Setting a record")
   // Checks the result after awaiting the set operation.
   // The Get method releases the resource of the future.
   status := future.Get()
   if !status.IsOK() {
     fmt.Println("ERROR: " + status.String())
   }

   // Executes the get method asynchronously.
   future = async.Get("hello")
   // Does something in the foreground.
   fmt.Println("Getting a record")
   // Checks the result after awaiting the get operation.
   // The GetStr method releases the resource of the future.
   value, status := future.GetStr()
   if status.IsOK() {
     fmt.Println("VALUE: " + value)
   }
 }

The following code uses Process, ProcessMulti, and ProcessEach methods which take callback functions to process the record efficiently.  Process is useful to update a record atomically according to the current value.  ProcessEach is useful to access every record in the most efficient way.

 package main

 import (
   "fmt"
   "github.com/estraier/tkrzw-go"
   "regexp"
   "strings"
 )

 func main() {
   // Opens the database.
   dbm := tkrzw.NewDBM()
   dbm.Open("casket.tkh", true,
     tkrzw.ParseParams("truncate=true,num_buckets=100"))
   defer dbm.Close()

   // Sets records with lambda functions.
   dbm.Process("doc-1", func(k []byte, v []byte) interface{} {
     return "Tokyo is the capital city of Japan."
   }, true)
   dbm.Process("doc-2", func(k []byte, v []byte) interface{} {
     return "Is she living in Tokyo, Japan?"
   }, true)
   dbm.Process("doc-3", func(k []byte, v []byte) interface{} {
     return "She must leave Tokyo!"
   }, true)

   // Lowers record values.
   lower := func(k []byte, v []byte) interface{} {
     // If no matching record, nil is given as the value.
     if v == nil {
       return nil
     }
     // Sets the new value.
     // Note that the key and the value are a "bytes" object.
     return strings.ToLower(string(v))
   }
   dbm.Process("doc-1", lower, true)
   dbm.Process("doc-2", lower, true)
   dbm.Process("doc-3", lower, true)
   dbm.Process("non-existent", lower, true)

   // Adds multiple records at once.
   ops := []tkrzw.KeyProcPair{
     {"doc-4", func(k []byte, v []byte) interface{} { return "Tokyo Go!" }},
     {"doc-5", func(k []byte, v []byte) interface{} { return "Japan Go!" }},
   }
   dbm.ProcessMulti(ops, true)

   // Modifies multiple records at once.
   dbm.ProcessMulti([]tkrzw.KeyProcPair{{"doc-4", lower}, {"doc-5", lower}}, true)

   // Checks the whole content.
   // This uses an external iterator and is relavively slow.
   for record := range dbm.EachStr() {
     fmt.Println(record.Key, record.Value)
   }

   // Function for word counting.
   wordCounts := make(map[string]int)
   wordSplitter := regexp.MustCompile("\\W")
   wordCounter := func(key []byte, value []byte) interface{} {
     if key == nil {
       return nil
     }
     words := wordSplitter.Split(string(value), -1)
     for _, word := range words {
       if len(word) == 0 {
         continue
       }
       wordCounts[word] += 1
     }
     return nil
   }

   // The second parameter should be false if the value is not updated.
   dbm.ProcessEach(wordCounter, false)
   for word, count := range wordCounts {
     fmt.Printf("%s = %d\n", word, count)
   }

   // Returning RemoveBytes by the callbacks removes the record.
   dbm.Process("doc-1", func(k []byte, v []byte) interface{} {
     return tkrzw.RemoveBytes
   }, true)
   println(dbm.CountSimple())
   dbm.ProcessMulti([]tkrzw.KeyProcPair{
     {"doc-4", func(k []byte, v []byte) interface{} { return tkrzw.RemoveBytes }},
     {"doc-5", func(k []byte, v []byte) interface{} { return tkrzw.RemoveBytes }},
   }, true)
   println(dbm.CountSimple())
   dbm.ProcessEach(func(k []byte, v []byte) interface{} {
     return tkrzw.RemoveBytes
   }, true)
   println(dbm.CountSimple())
 }

The following code is an example to use a secondary index, which is useful to organize records by non-primary keys.

 package main

 import (
   "fmt"
   "github.com/estraier/tkrzw-go"
 )

 func main() {
   // Opens the index.
   index := tkrzw.NewIndex()
   index.Open("casket.tkt", true,
     tkrzw.ParseParams("truncate=true,num_buckets=100")).OrDie()
   defer func() { index.Close().OrDie() }()

   // Adds records to the index.
   // The key is a division name and the value is person name.
   index.Add("general", "anne").OrDie()
   index.Add("general", "matthew").OrDie()
   index.Add("general", "marilla").OrDie()
   index.Add("sales", "gilbert").OrDie()

   // Anne moves to the sales division.
   index.Remove("general", "anne").OrDie()
   index.Add("sales", "anne").OrDie()

   // Prints all members for each division.
   divisions := [] string{"general", "sales"}
   for _, division := range divisions {
     fmt.Printf("%s\n", division)
     members := index.GetValuesStr(division, 0)
     for _, member := range members {
       fmt.Printf(" -- %s\n", member)
     }
   }

   // Prints every record by iterator.
   iter := index.MakeIterator()
   iter.First()
   for {
     key, value, ok := iter.GetStr()
     if !ok {
       break
     }
     fmt.Printf("%s: %s\n", key, value)
     iter.Next()
   }
 }
*/
package tkrzw

// END OF FILE
