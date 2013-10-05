#Directions & Examples
First, make sure every machine in our distributed system is running the server program. Ths server program is located at 'logServer/server.go'. Each machine should also have Golang installed to run the server.

Make sure every machine in the distributed system has a metadata.txt file with any relevant information about the machine written in the file. This file should be located in the same directory as the 'server.go' program.

Now, on any one machine, to query for certain key-value pairs in log files, we must run the program 'logClient/grep_client.go'. To work properly, the client program uses a 'masterlist.txt' file located in the same directory. This file contains a list of the IP addresses of each machine in the system.

TO run properly, the client program asks for two arguments, the first argument being the pattern matching for the key and the second argument for the value.
Ex: 'go run grep_client.go <keyPattern> <valuePattern>'

Note: to ignore any particular pattern in the key or value, use a wildcard statement in the respective argument. Don't use '.*' as the '.' and the '*' symbol in the first position of an argument breaks all other arguments in Golang. As a replacement, use '^.*' as the wildcard statement, as this works well with our program.

After running the client program with the arguments, the client will query the servers, and the response will be printed out on console.

Example queries=
```
'go run grep_client.go hello world'
'go run grep_client.go ^ERROR$ ^.*Hi.*$
```
#Underlying Architecture

Our distributed logging system allows you to run system grep calls through all the machines in our system and returns the relevant matches back to the user. Grep is executed on each machine's log separately. The machines are independent in the sense that they do not wait for each other to complete. The client receives the response as soon as a machine completes executing grep on its log. As the machines are not waiting for each other there is no blocking in the execution process.

The machine where we are querying from has a masterlist.txt that contains the ip addresses of all the machines we will query including itself. When we want to add a new machine we have to update the masterlist.txt. Similarly, when we want to remove a machine from the system we just delete the ip address of the machine from masterlist.txt.

Each machine in our distributed system has a metadata.txt file that contains it's ip address, name of the machine and the log file it should operate on. The machine also has a logServer that runs a grep command on this log file upon request from our logClient.

In order to make sure that our system is fault-tolerant we ignore machines that are down. Therefore, as soon as we observe that a machine in our system is down the whole system does not break down. In the worst case (where all machines are down), we return results from the machine that is invoking the grep_client.go.

#Does our system work?

Apart from testing the functionality manually, we created unit tests to make sure our system works. Our test generates a test log on a remote machine on the fly and runs a grep on it to check if we are retrieving results that matches with our expected results. The unit test considers different forms of data for example, rare, frequent and somewhat frequent strings. 

#Average query latency
We tested our system which contained 4 different machines containing a log file of size 137 MB each. This is how the results look like
###Rare (Query results appear in individual log files 3.85%)

1. 624.354182 ms

2. 623.701059 ms 

3. 614.587465 ms

4. 615.16865  ms

5. 633.258762 ms

Average latency: 622.2140236 ms

###Frequent (Query results appera in individual log files 19.23%)

1. 9.465099 s

2. 9.436086196 s

3. 9.506596702 s

4. 9.401853494 s

5. 9.480838906 s

Average latency: 9.4580948596 s

#Authors
###Irtefa, Mohd
###Lee, Stephen
